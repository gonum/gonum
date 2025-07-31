// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build noasm || gccgo || safe || !amd64
// +build noasm gccgo safe !amd64

package f64

import "unsafe"

// HasGemmKernel indicates whether optimized GEMM kernels are available
var HasGemmKernel = false

// GemmKernel4x4 computes C[4x4] += A[4xK] * B[Kx4]
// This is the fallback implementation for platforms without assembly support
func GemmKernel4x4(a *float64, b *float64, c *float64, k int, lda, ldb, ldc int) {
	// Create slices from pointers
	aSlice := (*[1 << 30]float64)(unsafe.Pointer(a))[:4*k:4*k]
	bSlice := (*[1 << 30]float64)(unsafe.Pointer(b))[:k*ldb:k*ldb]
	cSlice := (*[1 << 30]float64)(unsafe.Pointer(c))[:4*ldc:4*ldc]
	
	// Compute C[4x4] += A[4xK] * B[Kx4]
	for i := 0; i < 4; i++ {
		for j := 0; j < 4; j++ {
			sum := cSlice[i*ldc+j]
			for kk := 0; kk < k; kk++ {
				sum += aSlice[i*lda+kk] * bSlice[kk*ldb+j]
			}
			cSlice[i*ldc+j] = sum
		}
	}
}

// GemmKernel8x8 computes C[8x8] += A[8xK] * B[Kx8]
// Fallback implementation
func GemmKernel8x8(a *float64, b *float64, c *float64, k int, lda, ldb, ldc int) {
	// Use 4x4 kernels as fallback
	GemmKernel4x4(a, b, c, k, lda, ldb, ldc)
	GemmKernel4x4(a, (*float64)(unsafe.Pointer(uintptr(unsafe.Pointer(b))+4*8)), 
		(*float64)(unsafe.Pointer(uintptr(unsafe.Pointer(c))+4*8)), k, lda, ldb, ldc)
	GemmKernel4x4((*float64)(unsafe.Pointer(uintptr(unsafe.Pointer(a))+4*lda*8)), b, 
		(*float64)(unsafe.Pointer(uintptr(unsafe.Pointer(c))+4*ldc*8)), k, lda, ldb, ldc)
	GemmKernel4x4((*float64)(unsafe.Pointer(uintptr(unsafe.Pointer(a))+4*lda*8)), 
		(*float64)(unsafe.Pointer(uintptr(unsafe.Pointer(b))+4*8)), 
		(*float64)(unsafe.Pointer(uintptr(unsafe.Pointer(c))+4*ldc*8+4*8)), k, lda, ldb, ldc)
}