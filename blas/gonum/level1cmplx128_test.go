// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"testing"

	"gonum.org/v1/gonum/blas/testblas"
)

func TestZaxpy(t *testing.T) {
	testblas.ZaxpyTest(t, impl)
}

func TestZcopy(t *testing.T) {
	testblas.ZcopyTest(t, impl)
}

func TestZdotc(t *testing.T) {
	testblas.ZdotcTest(t, impl)
}

func TestZdotu(t *testing.T) {
	testblas.ZdotuTest(t, impl)
}
