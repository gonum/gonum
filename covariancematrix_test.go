// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"testing"

	"github.com/gonum/blas/cblas"
	"github.com/gonum/floats"
	"github.com/gonum/matrix/mat64"
)

func init() {
	mat64.Register(cblas.Blas{})
}

func TestCovarianceMatrix(t *testing.T) {
	for i, test := range []struct {
		mat  mat64.Matrix
		r, c int
		x    []float64
	}{
		{
			mat: mat64.NewDense(5, 2, []float64{
				-2, -4,
				-1, 2,
				0, 0,
				1, -2,
				2, 4,
			}),
			r: 2,
			c: 2,
			x: []float64{
				2.5, 3,
				3, 10,
			},
		},
	} {
		c := CovarianceMatrix(test.mat).RawMatrix()
		if c.Rows != test.r {
			t.Errorf("%d: expected rows %d, found %d", i, test.r, c.Rows)
		}
		if c.Cols != test.c {
			t.Errorf("%d: expected cols %d, found %d", i, test.c, c.Cols)
		}
		if !floats.Equal(test.x, c.Data) {
			t.Errorf("%d: expected data %#q, found %#q", i, test.x, c.Data)
		}
	}
}
