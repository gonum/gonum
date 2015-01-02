// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package zbw is a wrapper for a blas.Float64 implementation.
package dbw

import "github.com/gonum/blas"

var impl blas.Float64

func Register(i blas.Float64) {
	impl = i
}
