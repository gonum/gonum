// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

package c64

func AxpyUnitary(alpha complex64, x, y []complex64)

func AxpyUnitaryTo(dst []complex64, alpha complex64, x, y []complex64)

func AxpyInc(alpha complex64, x, y []complex64, n, incX, incY, ix, iy uintptr)

func AxpyIncTo(dst []complex64, incDst, idst uintptr, alpha complex64, x, y []complex64, n, incX, incY, ix, iy uintptr)
