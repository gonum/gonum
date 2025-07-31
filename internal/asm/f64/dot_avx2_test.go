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

func TestDotUnitaryAVX2(t *testing.T) {
	if !cpu.X86.HasAVX2 || !cpu.X86.HasFMA {
		t.Skip("AVX2/FMA not available")
	}
	
	tests := []struct {
		x, y []float64
		want float64
	}{
		// Simple cases
		{[]float64{1}, []float64{2}, 2},
		{[]float64{1, 2}, []float64{3, 4}, 11},
		{[]float64{1, 2, 3}, []float64{4, 5, 6}, 32},
		{[]float64{1, 2, 3, 4}, []float64{5, 6, 7, 8}, 70},
		
		// Test different lengths to exercise all code paths
		{[]float64{1, 2, 3, 4, 5}, []float64{1, 1, 1, 1, 1}, 15},
		{[]float64{2, 2, 2, 2, 2, 2, 2, 2}, []float64{3, 3, 3, 3, 3, 3, 3, 3}, 48},
		
		// 16 elements to test main loop
		{
			[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			[]float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			136, // sum of 1 to 16
		},
		
		// 17 elements to test tail handling
		{
			[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
			[]float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
			153, // sum of 1 to 17
		},
	}
	
	for i, test := range tests {
		got := DotUnitaryAVX2(test.x, test.y)
		if math.Abs(got-test.want) > 1e-14 {
			t.Errorf("test %d: DotUnitaryAVX2(%v, %v) = %v, want %v", 
				i, test.x, test.y, got, test.want)
		}
	}
}

// Benchmark comparing implementations
func BenchmarkDotUnitary(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000, 100000}
	
	for _, size := range sizes {
		x := make([]float64, size)
		y := make([]float64, size)
		
		// Initialize with some data
		for i := range x {
			x[i] = float64(i)
			y[i] = float64(i + 1)
		}
		
		b.Run("Size"+fmt.Sprintf("%d", size), func(b *testing.B) {
			b.Run("PureGo", func(b *testing.B) {
				b.SetBytes(int64(size * 8 * 2))
				var sum float64
				for i := 0; i < b.N; i++ {
					// Pure Go implementation
					sum = 0
					for j, v := range x {
						sum += y[j] * v
					}
				}
				_ = sum
			})
			
			b.Run("SSE2", func(b *testing.B) {
				b.SetBytes(int64(size * 8 * 2))
				var sum float64
				for i := 0; i < b.N; i++ {
					sum = dotUnitarySSE2Asm(x, y)
				}
				_ = sum
			})
			
			if cpu.X86.HasAVX2 && cpu.X86.HasFMA {
				b.Run("AVX2", func(b *testing.B) {
					b.SetBytes(int64(size * 8 * 2))
					var sum float64
					for i := 0; i < b.N; i++ {
						sum = DotUnitaryAVX2(x, y)
					}
					_ = sum
				})
			}
		})
	}
}