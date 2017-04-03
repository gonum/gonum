// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

package f32

func AxpyUnitary(alpha float32, x, y []float32)

func AxpyUnitaryTo(dst []float32, alpha float32, x, y []float32)

func AxpyInc(alpha float32, x, y []float32, n, incX, incY, ix, iy uintptr)

func AxpyIncTo(dst []float32, incDst, idst uintptr, alpha float32, x, y []float32, n, incX, incY, ix, iy uintptr)
