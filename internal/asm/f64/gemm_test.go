// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package f64

import (
	"fmt"
	"math"
	"testing"
)

// TestGemmKernel4x4 tests the 4x4 GEMM micro-kernel
func TestGemmKernel4x4(t *testing.T) {
	// Test case: 4x3 * 3x4 = 4x4
	k := 3
	
	// A matrix (4x3)
	a := []float64{
		1, 2, 3,
		4, 5, 6,
		7, 8, 9,
		10, 11, 12,
	}
	
	// B matrix (3x4)
	b := []float64{
		1, 0, 0, 1,
		0, 1, 0, 1,
		0, 0, 1, 1,
	}
	
	// C matrix (4x4) - initialized to zero
	c := make([]float64, 16)
	
	// Expected result
	expected := []float64{
		1, 2, 3, 6,    // [1 2 3] * [[1 0 0 1] [0 1 0 1] [0 0 1 1]]
		4, 5, 6, 15,   // [4 5 6] * B
		7, 8, 9, 24,   // [7 8 9] * B
		10, 11, 12, 33, // [10 11 12] * B
	}
	
	// Call the kernel
	GemmKernel4x4(&a[0], &b[0], &c[0], k, 3, 4, 4)
	
	// Check results
	for i := 0; i < 16; i++ {
		if math.Abs(c[i]-expected[i]) > 1e-14 {
			t.Errorf("GemmKernel4x4: c[%d] = %f, want %f", i, c[i], expected[i])
		}
	}
}

// TestGemmKernel4x4Accumulate tests that the kernel accumulates into C
func TestGemmKernel4x4Accumulate(t *testing.T) {
	k := 2
	
	a := []float64{
		1, 2,
		3, 4,
		5, 6,
		7, 8,
	}
	
	b := []float64{
		1, 0, 0, 0,
		0, 1, 0, 0,
	}
	
	// C matrix with initial values
	c := []float64{
		10, 10, 10, 10,
		20, 20, 20, 20,
		30, 30, 30, 30,
		40, 40, 40, 40,
	}
	
	// Expected: C += A*B
	expected := []float64{
		11, 12, 10, 10,  // [10 10 10 10] + [1 2 0 0]
		23, 24, 20, 20,  // [20 20 20 20] + [3 4 0 0]
		35, 36, 30, 30,  // [30 30 30 30] + [5 6 0 0]
		47, 48, 40, 40,  // [40 40 40 40] + [7 8 0 0]
	}
	
	GemmKernel4x4(&a[0], &b[0], &c[0], k, 2, 4, 4)
	
	for i := 0; i < 16; i++ {
		if math.Abs(c[i]-expected[i]) > 1e-14 {
			t.Errorf("GemmKernel4x4 accumulate: c[%d] = %f, want %f", i, c[i], expected[i])
		}
	}
}

// TestGemmKernel8x8 tests the 8x8 GEMM micro-kernel
func TestGemmKernel8x8(t *testing.T) {
	// Test case: 8x4 * 4x8 = 8x8
	k := 4
	
	// A matrix (8x4)
	a := []float64{
		1, 2, 3, 4,
		5, 6, 7, 8,
		9, 10, 11, 12,
		13, 14, 15, 16,
		17, 18, 19, 20,
		21, 22, 23, 24,
		25, 26, 27, 28,
		29, 30, 31, 32,
	}
	
	// B matrix (4x8) - identity-like for easy verification
	b := []float64{
		1, 0, 0, 0, 1, 0, 0, 0,
		0, 1, 0, 0, 0, 1, 0, 0,
		0, 0, 1, 0, 0, 0, 1, 0,
		0, 0, 0, 1, 0, 0, 0, 1,
	}
	
	// C matrix (8x8) - initialized to zero
	c := make([]float64, 64)
	
	// Call the kernel
	GemmKernel8x8(&a[0], &b[0], &c[0], k, 4, 8, 8)
	
	// Expected: each row of A appears twice in C
	// Row 0: [1,2,3,4, 1,2,3,4]
	// Row 1: [5,6,7,8, 5,6,7,8]
	// etc.
	for i := 0; i < 8; i++ {
		for j := 0; j < 4; j++ {
			expected := a[i*4+j]
			if math.Abs(c[i*8+j]-expected) > 1e-14 {
				t.Errorf("GemmKernel8x8: c[%d,%d] = %f, want %f", i, j, c[i*8+j], expected)
			}
			if math.Abs(c[i*8+j+4]-expected) > 1e-14 {
				t.Errorf("GemmKernel8x8: c[%d,%d] = %f, want %f", i, j+4, c[i*8+j+4], expected)
			}
		}
	}
}

// BenchmarkGemmKernel4x4 benchmarks the 4x4 kernel
func BenchmarkGemmKernel4x4(b *testing.B) {
	sizes := []int{4, 8, 16, 32, 64, 128}
	
	for _, k := range sizes {
		b.Run("K"+fmt.Sprintf("%d", k), func(b *testing.B) {
			// Create matrices
			a := make([]float64, 4*k)
			bm := make([]float64, k*4)
			c := make([]float64, 16)
			
			// Initialize with some data
			for i := range a {
				a[i] = float64(i)
			}
			for i := range bm {
				bm[i] = float64(i)
			}
			
			b.ResetTimer()
			b.SetBytes(int64(4*k + k*4 + 16) * 8) // Total memory touched
			
			for i := 0; i < b.N; i++ {
				GemmKernel4x4(&a[0], &bm[0], &c[0], k, k, 4, 4)
			}
		})
	}
}

// BenchmarkGemmKernel8x8 benchmarks the 8x8 kernel
func BenchmarkGemmKernel8x8(b *testing.B) {
	sizes := []int{8, 16, 32, 64, 128, 256}
	
	for _, k := range sizes {
		b.Run("K"+fmt.Sprintf("%d", k), func(b *testing.B) {
			// Create matrices
			a := make([]float64, 8*k)
			bm := make([]float64, k*8)
			c := make([]float64, 64)
			
			// Initialize with some data
			for i := range a {
				a[i] = float64(i)
			}
			for i := range bm {
				bm[i] = float64(i)
			}
			
			b.ResetTimer()
			b.SetBytes(int64(8*k + k*8 + 64) * 8) // Total memory touched
			
			for i := 0; i < b.N; i++ {
				GemmKernel8x8(&a[0], &bm[0], &c[0], k, k, 8, 8)
			}
		})
	}
}