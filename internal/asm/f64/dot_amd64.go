// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build !noasm && !gccgo && !safe
// +build !noasm,!gccgo,!safe

package f64

import "golang.org/x/sys/cpu"

// Function pointers for runtime CPU feature detection
var (
	dotUnitaryImpl func(x, y []float64) float64
	hasDotAVX2     bool
)

// Assembly functions - rename to avoid conflicts
func dotUnitarySSE2Asm(x, y []float64) float64
func DotUnitaryAVX2(x, y []float64) float64

func init() {
	if cpu.X86.HasAVX2 && cpu.X86.HasFMA {
		dotUnitaryImpl = DotUnitaryAVX2
		hasDotAVX2 = true
	} else {
		dotUnitaryImpl = dotUnitarySSE2Asm
		hasDotAVX2 = false
	}
}

// DotUnitary is
//
//	for i, v := range x {
//		sum += y[i] * v
//	}
//	return sum
func DotUnitary(x, y []float64) float64 {
	return dotUnitaryImpl(x, y)
}

// Export hasDotAVX2 for testing
var HasDotAVX2 = &hasDotAVX2