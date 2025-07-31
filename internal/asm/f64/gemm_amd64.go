// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !noasm && !gccgo && !safe
// +build !noasm,!gccgo,!safe

package f64

import (
	"golang.org/x/sys/cpu"
	"unsafe"
)

// HasGemmKernel indicates whether optimized GEMM kernels are available
var HasGemmKernel = cpu.X86.HasAVX2 && cpu.X86.HasFMA

// GEMM kernel configuration
const (
	// Micro-kernel dimensions
	MR = 4 // Number of rows in micro-kernel
	NR = 4 // Number of columns in micro-kernel
	
	// Larger kernel dimensions
	MR8 = 8 // 8x8 kernel rows
	NR8 = 8 // 8x8 kernel columns
)

// Assembly functions
func gemmKernel4x4AVX2(a *float64, b *float64, c *float64, k int, lda, ldb, ldc int)
func gemmKernel8x8AVX2(a *float64, b *float64, c *float64, k int, lda, ldb, ldc int)

// GemmKernel4x4 computes C[4x4] += A[4xK] * B[Kx4]
// This is the core micro-kernel for GEMM operations
func GemmKernel4x4(a *float64, b *float64, c *float64, k int, lda, ldb, ldc int) {
	if cpu.X86.HasAVX2 && cpu.X86.HasFMA {
		gemmKernel4x4AVX2(a, b, c, k, lda, ldb, ldc)
	} else {
		// Fallback to scalar implementation
		gemmKernel4x4Scalar(a, b, c, k, lda, ldb, ldc)
	}
}

// Scalar fallback implementation
func gemmKernel4x4Scalar(a *float64, b *float64, c *float64, k int, lda, ldb, ldc int) {
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

// GemmMicroKernel is the main entry point for blocked GEMM
// It processes an MRxNR block of C using the appropriate kernel
func GemmMicroKernel(m, n, k int, alpha float64,
	a []float64, lda int,
	b []float64, ldb int,
	c []float64, ldc int) {
	
	// For now, only handle 4x4 blocks
	if m == MR && n == NR && alpha == 1.0 {
		GemmKernel4x4(&a[0], &b[0], &c[0], k, lda, ldb, ldc)
		return
	}
	
	// Fallback to reference implementation for other cases
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			sum := 0.0
			for kk := 0; kk < k; kk++ {
				sum += a[i*lda+kk] * b[kk*ldb+j]
			}
			c[i*ldc+j] += alpha * sum
		}
	}
}

// GemmKernel8x8 computes C[8x8] += A[8xK] * B[Kx8]
// This is the larger micro-kernel for GEMM operations
func GemmKernel8x8(a *float64, b *float64, c *float64, k int, lda, ldb, ldc int) {
	if cpu.X86.HasAVX2 && cpu.X86.HasFMA {
		gemmKernel8x8AVX2(a, b, c, k, lda, ldb, ldc)
	} else {
		// Fallback to scalar implementation
		aSlice := (*[1 << 30]float64)(unsafe.Pointer(a))[:8*k:8*k]
		bSlice := (*[1 << 30]float64)(unsafe.Pointer(b))[:k*ldb:k*ldb]
		cSlice := (*[1 << 30]float64)(unsafe.Pointer(c))[:8*ldc:8*ldc]
		
		for i := 0; i < 8; i++ {
			for j := 0; j < 8; j++ {
				sum := cSlice[i*ldc+j]
				for kk := 0; kk < k; kk++ {
					sum += aSlice[i*lda+kk] * bSlice[kk*ldb+j]
				}
				cSlice[i*ldc+j] = sum
			}
		}
	}
}