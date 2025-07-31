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

func TestScalUnitaryAVX2(t *testing.T) {
	if !cpu.X86.HasAVX2 {
		t.Skip("AVX2 not available")
	}
	
	tests := []struct {
		alpha float64
		x     []float64
		want  []float64
	}{
		// Simple cases
		{2, []float64{1}, []float64{2}},
		{3, []float64{1, 2}, []float64{3, 6}},
		{-2, []float64{1, 2, 3}, []float64{-2, -4, -6}},
		{0.5, []float64{2, 4, 6, 8}, []float64{1, 2, 3, 4}},
		
		// Test different lengths to exercise all code paths
		{2, []float64{1, 2, 3, 4, 5}, []float64{2, 4, 6, 8, 10}},
		{3, []float64{1, 1, 1, 1, 1, 1, 1, 1}, []float64{3, 3, 3, 3, 3, 3, 3, 3}},
		
		// 16 elements to test main loop
		{
			2,
			[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
			[]float64{2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32},
		},
		
		// 17 elements to test tail handling
		{
			0.5,
			[]float64{2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34},
			[]float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17},
		},
	}
	
	for i, test := range tests {
		x := make([]float64, len(test.x))
		copy(x, test.x)
		ScalUnitaryAVX2(test.alpha, x)
		for j, v := range x {
			if math.Abs(v-test.want[j]) > 1e-14 {
				t.Errorf("test %d: ScalUnitaryAVX2(%v, %v)[%d] = %v, want %v", 
					i, test.alpha, test.x, j, v, test.want[j])
			}
		}
	}
}

// Benchmark comparing implementations
func BenchmarkScalUnitaryAVX2(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000, 100000}
	
	for _, size := range sizes {
		x := make([]float64, size)
		
		// Initialize with some data
		for i := range x {
			x[i] = float64(i)
		}
		
		b.Run("Size"+fmt.Sprintf("%d", size), func(b *testing.B) {
			b.Run("PureGo", func(b *testing.B) {
				b.SetBytes(int64(size * 8))
				xcopy := make([]float64, len(x))
				for i := 0; i < b.N; i++ {
					copy(xcopy, x)
					// Pure Go implementation
					for j := range xcopy {
						xcopy[j] *= 2.5
					}
				}
			})
			
			b.Run("SSE2", func(b *testing.B) {
				b.SetBytes(int64(size * 8))
				xcopy := make([]float64, len(x))
				for i := 0; i < b.N; i++ {
					copy(xcopy, x)
					scalUnitarySSE2(2.5, xcopy)
				}
			})
			
			if cpu.X86.HasAVX2 {
				b.Run("AVX2", func(b *testing.B) {
					b.SetBytes(int64(size * 8))
					xcopy := make([]float64, len(x))
					for i := 0; i < b.N; i++ {
						copy(xcopy, x)
						ScalUnitaryAVX2(2.5, xcopy)
					}
				})
			}
		})
	}
}