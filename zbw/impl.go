// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package zbw is a wrapper for a blas.Complex128 implementation.
package zbw

import "github.com/gonum/blas"

var impl blas.Complex128

func Register(i blas.Complex128) {
	impl = i
}
