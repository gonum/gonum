// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !noasm && !gccgo && !safe
// +build !noasm,!gccgo,!safe

package f64

import "golang.org/x/sys/cpu"

// Function pointers for runtime CPU feature detection
var (
	gemvNImpl  func(m, n uintptr, alpha float64, a []float64, lda uintptr, x []float64, incX uintptr, beta float64, y []float64, incY uintptr)
	hasGemvAVX2 bool
)

// Assembly functions
func gemvNSSE2(m, n uintptr, alpha float64, a []float64, lda uintptr, x []float64, incX uintptr, beta float64, y []float64, incY uintptr)
func GemvNAVX2(m, n uintptr, alpha float64, a []float64, lda uintptr, x []float64, incX uintptr, beta float64, y []float64, incY uintptr)

func init() {
	if cpu.X86.HasAVX2 && cpu.X86.HasFMA {
		gemvNImpl = GemvNAVX2
		hasGemvAVX2 = true
	} else {
		gemvNImpl = gemvNSSE2
		hasGemvAVX2 = false
	}
}

// GemvN computes
//
//	y = alpha * A * x + beta * y
//
// where A is an m×n dense matrix, x and y are vectors, and alpha and beta are scalars.
func GemvN(m, n uintptr, alpha float64, a []float64, lda uintptr, x []float64, incX uintptr, beta float64, y []float64, incY uintptr) {
	gemvNImpl(m, n, alpha, a, lda, x, incX, beta, y, incY)
}

// Forward declaration for existing assembly implementation
func gemvN(m, n uintptr, alpha float64, a []float64, lda uintptr, x []float64, incX uintptr, beta float64, y []float64, incY uintptr)

// Export hasGemvAVX2 for testing
var HasGemvAVX2 = &hasGemvAVX2