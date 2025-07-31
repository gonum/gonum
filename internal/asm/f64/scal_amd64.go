// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !noasm && !gccgo && !safe
// +build !noasm,!gccgo,!safe

package f64

import "golang.org/x/sys/cpu"

// Function pointers for runtime CPU feature detection
var (
	scalUnitaryImpl func(alpha float64, x []float64)
	hasScalAVX2     bool
)

// Assembly functions - rename to avoid conflicts
func scalUnitarySSE2(alpha float64, x []float64)
func ScalUnitaryAVX2(alpha float64, x []float64)

func init() {
	if cpu.X86.HasAVX2 {
		scalUnitaryImpl = ScalUnitaryAVX2
		hasScalAVX2 = true
	} else {
		scalUnitaryImpl = scalUnitarySSE2
		hasScalAVX2 = false
	}
}

// ScalUnitary is
//
//	for i := range x {
//		x[i] *= alpha
//	}
func ScalUnitary(alpha float64, x []float64) {
	scalUnitaryImpl(alpha, x)
}

// Forward declaration for existing assembly implementation
func scalUnitary(alpha float64, x []float64)

// Export hasScalAVX2 for testing
var HasScalAVX2 = &hasScalAVX2