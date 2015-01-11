// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !amd64

package asm

func DdotUnitary(x []float64, y []float64) (sum float64) {
	for i, v := range x {
		sum += y[i] * v
	}
	return
}
