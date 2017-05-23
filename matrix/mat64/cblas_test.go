// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build cblas

package mat64

import (
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/blas/cgo"
)

func init() {
	blas64.Use(cgo.Implementation{})
}
