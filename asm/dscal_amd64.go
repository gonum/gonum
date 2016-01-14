// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

package asm

func DscalUnitary(alpha float64, x []float64)

func DscalUnitaryTo(dst []float64, alpha float64, x []float64)

func DscalInc(alpha float64, x []float64, n, incX, ix uintptr)

func DscalIncTo(dst []float64, incDst, idst uintptr, alpha float64, x []float64, n, incX, ix uintptr)
