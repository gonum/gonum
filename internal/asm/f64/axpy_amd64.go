// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !noasm && !gccgo && !safe
// +build !noasm,!gccgo,!safe

package f64

import "golang.org/x/sys/cpu"

// Function pointers for runtime CPU feature detection
var (
	axpyUnitaryImpl      func(alpha float64, x, y []float64)
	axpyUnitaryToImpl    func(dst []float64, alpha float64, x, y []float64)
)

// Assembly functions
func axpyUnitarySSE2Asm(alpha float64, x, y []float64)
func axpyUnitaryToSSE2Asm(dst []float64, alpha float64, x, y []float64)
func AxpyUnitaryAVX2(alpha float64, x, y []float64)
func AxpyUnitaryFMA(alpha float64, x, y []float64)
func axpyUnitaryToAVX2(dst []float64, alpha float64, x, y []float64)

func init() {
	if cpu.X86.HasAVX2 && cpu.X86.HasFMA {
		axpyUnitaryImpl = AxpyUnitaryFMA
		// axpyUnitaryToImpl = axpyUnitaryToAVX2 // TODO: implement AVX2 version
		axpyUnitaryToImpl = axpyUnitaryToSSE2Asm
	} else if cpu.X86.HasAVX2 {
		axpyUnitaryImpl = AxpyUnitaryAVX2
		axpyUnitaryToImpl = axpyUnitaryToSSE2Asm
	} else {
		axpyUnitaryImpl = axpyUnitarySSE2Asm
		axpyUnitaryToImpl = axpyUnitaryToSSE2Asm
	}
}

// AxpyUnitary is
//
//	for i, v := range x {
//		y[i] += alpha * v
//	}
func AxpyUnitary(alpha float64, x, y []float64) {
	axpyUnitaryImpl(alpha, x, y)
}

// AxpyUnitaryTo is
//
//	for i, v := range x {
//		dst[i] = alpha*v + y[i]
//	}
func AxpyUnitaryTo(dst []float64, alpha float64, x, y []float64) {
	axpyUnitaryToImpl(dst, alpha, x, y)
}

