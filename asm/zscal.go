// Generated code do not edit. Run `go generate`.

// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package asm

func ZscalUnitary(alpha complex128, x []complex128) {
	for i := range x {
		x[i] *= alpha
	}
}

func ZscalUnitaryTo(dst []complex128, alpha complex128, x []complex128) {
	for i, v := range x {
		dst[i] = alpha * v
	}
}
