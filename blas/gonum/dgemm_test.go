// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"fmt"
	"testing"
	"gonum.org/v1/gonum/blas"
)

func BenchmarkDgemmOptimized(b *testing.B) {
	sizes := []struct {
		m, n, k int
	}{
		{10, 10, 10},
		{20, 20, 20},
		{50, 50, 50},
		{100, 100, 100},
		{200, 200, 200},
		{500, 500, 500},
	}
	
	impl := Implementation{}
	
	for _, size := range sizes {
		// Allocate matrices
		a := make([]float64, size.m*size.k)
		bb := make([]float64, size.k*size.n)
		c := make([]float64, size.m*size.n)
		
		// Initialize with some data
		for i := range a {
			a[i] = float64(i)
		}
		for i := range bb {
			bb[i] = float64(i)
		}
		
		b.Run(fmt.Sprintf("Size%dx%dx%d", size.m, size.n, size.k), func(b *testing.B) {
			b.SetBytes(int64(size.m*size.k + size.k*size.n + size.m*size.n) * 8)
			for i := 0; i < b.N; i++ {
				// Reset C
				for j := range c {
					c[j] = 0
				}
				impl.Dgemm(blas.NoTrans, blas.NoTrans, size.m, size.n, size.k, 
					1.0, a, size.k, bb, size.n, 0.0, c, size.n)
			}
		})
	}
}