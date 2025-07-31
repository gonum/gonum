// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package f64

import (
	"fmt"
	"math"
	"testing"
	"golang.org/x/sys/cpu"
)

func TestGemvNAVX2(t *testing.T) {
	if !cpu.X86.HasAVX2 || !cpu.X86.HasFMA {
		t.Skip("AVX2/FMA not available")
	}
	
	tests := []struct {
		m, n   int
		alpha  float64
		a      []float64
		lda    int
		x      []float64
		incX   int
		beta   float64
		y      []float64
		incY   int
		want   []float64
	}{
		// Simple 2x2 case
		{
			m: 2, n: 2, alpha: 1, beta: 0,
			a:    []float64{1, 2, 3, 4},
			lda:  2,
			x:    []float64{1, 1},
			incX: 1,
			y:    []float64{0, 0},
			incY: 1,
			want: []float64{3, 7}, // [1 2] * [1] = [3], [3 4] * [1] = [7]
		},
		// 3x3 case with alpha and beta
		{
			m: 3, n: 3, alpha: 2, beta: 1,
			a:    []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			lda:  3,
			x:    []float64{1, 0, 1},
			incX: 1,
			y:    []float64{1, 1, 1},
			incY: 1,
			want: []float64{9, 21, 33}, // 2*([1 2 3]*[1 0 1]') + 1*[1] = 2*4 + 1 = 9
		},
		// 4x4 case to test main loop
		{
			m: 4, n: 4, alpha: 1, beta: 0,
			a: []float64{
				1, 0, 0, 0,
				0, 1, 0, 0,
				0, 0, 1, 0,
				0, 0, 0, 1,
			},
			lda:  4,
			x:    []float64{1, 2, 3, 4},
			incX: 1,
			y:    []float64{0, 0, 0, 0},
			incY: 1,
			want: []float64{1, 2, 3, 4}, // Identity matrix
		},
		// 5x3 case to test remainder handling
		{
			m: 5, n: 3, alpha: 1, beta: 0,
			a: []float64{
				1, 1, 1,
				2, 2, 2,
				3, 3, 3,
				4, 4, 4,
				5, 5, 5,
			},
			lda:  3,
			x:    []float64{1, 1, 1},
			incX: 1,
			y:    []float64{0, 0, 0, 0, 0},
			incY: 1,
			want: []float64{3, 6, 9, 12, 15},
		},
	}
	
	for i, test := range tests {
		y := make([]float64, len(test.y))
		copy(y, test.y)
		
		GemvNAVX2(uintptr(test.m), uintptr(test.n), test.alpha, test.a, uintptr(test.lda),
			test.x, uintptr(test.incX), test.beta, y, uintptr(test.incY))
		
		for j, v := range y {
			if math.Abs(v-test.want[j]) > 1e-14 {
				t.Errorf("test %d: GemvNAVX2 y[%d] = %v, want %v", i, j, v, test.want[j])
			}
		}
	}
}

// Benchmark comparing implementations
func BenchmarkGemvN(b *testing.B) {
	sizes := []struct{ m, n int }{
		{10, 10},
		{100, 100},
		{1000, 100},
		{100, 1000},
		{1000, 1000},
	}
	
	for _, size := range sizes {
		a := make([]float64, size.m*size.n)
		x := make([]float64, size.n)
		y := make([]float64, size.m)
		
		// Initialize with some data
		for i := range a {
			a[i] = float64(i % 10)
		}
		for i := range x {
			x[i] = float64(i)
		}
		
		name := fmt.Sprintf("M%dN%d", size.m, size.n)
		b.Run(name, func(b *testing.B) {
			b.Run("PureGo", func(b *testing.B) {
				b.SetBytes(int64((size.m*size.n + size.m + size.n) * 8))
				ycopy := make([]float64, len(y))
				for i := 0; i < b.N; i++ {
					copy(ycopy, y)
					// Pure Go implementation
					for i := 0; i < size.m; i++ {
						sum := 0.0
						for j := 0; j < size.n; j++ {
							sum += a[i*size.n+j] * x[j]
						}
						ycopy[i] = 2.0*sum + 1.0*ycopy[i]
					}
				}
			})
			
			b.Run("SSE2", func(b *testing.B) {
				b.SetBytes(int64((size.m*size.n + size.m + size.n) * 8))
				ycopy := make([]float64, len(y))
				for i := 0; i < b.N; i++ {
					copy(ycopy, y)
					gemvNSSE2(uintptr(size.m), uintptr(size.n), 2.0, a, uintptr(size.n),
						x, 1, 1.0, ycopy, 1)
				}
			})
			
			// Always benchmark AVX2 directly
			b.Run("AVX2", func(b *testing.B) {
				if !cpu.X86.HasAVX2 || !cpu.X86.HasFMA {
					b.Skip("AVX2/FMA not available")
				}
				b.SetBytes(int64((size.m*size.n + size.m + size.n) * 8))
				ycopy := make([]float64, len(y))
				for i := 0; i < b.N; i++ {
					copy(ycopy, y)
					GemvNAVX2(uintptr(size.m), uintptr(size.n), 2.0, a, uintptr(size.n),
						x, 1, 1.0, ycopy, 1)
				}
			})
		})
	}
}