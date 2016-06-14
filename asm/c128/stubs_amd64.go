// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

package c128

func AxpyUnitary(alpha complex128, x, y []complex128)

func AxpyUnitaryTo(dst []complex128, alpha complex128, x, y []complex128)

func AxpyInc(alpha complex128, x, y []complex128, n, incX, incY, ix, iy uintptr)

func AxpyIncTo(dst []complex128, incDst, idst uintptr, alpha complex128, x, y []complex128, n, incX, incY, ix, iy uintptr)
