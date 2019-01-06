// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"testing"

	"gonum.org/v1/gonum/blas/testblas"
)

func TestZgemm(t *testing.T) {
	testblas.ZgemmTest(t, impl)
}

func TestZsyrk(t *testing.T) {
	testblas.ZsyrkTest(t, impl)
}
