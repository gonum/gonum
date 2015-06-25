// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package native

import "github.com/gonum/lapack"

type Implementation struct{}

var _ lapack.Float64 = Implementation{}

const (
	badUplo = "lapack: illegal triangle"
	nLT0    = "lapack: n < 0"
	badLdA  = "lapack: index of a out of range"
)

func blockSize() int {
	return 64
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
