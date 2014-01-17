// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"fmt"
	check "launchpad.net/gocheck"
	"math"
)

func isLowerTriangular(a *Dense) bool {
	rows, cols := a.Dims()
	for r := 0; r < rows; r++ {
		for c := r + 1; c < cols; c++ {
			if math.Abs(a.At(r, c)) > 1e-14 {
				return false
			}
		}
	}
	return true
}

func (s *S) TestLQD(c *check.C) {
	for _, test := range []struct {
		a    [][]float64
		name string
	}{
		{
			name: "Square",
			a: [][]float64{
				{1.3, 2.4, 8.9},
				{-2.6, 8.7, 9.1},
				{5.6, 5.8, 2.1},
			},
		},
		{
			name: "Skinny",
			a: [][]float64{
				{1.3, 2.4, 8.9},
				{-2.6, 8.7, 9.1},
				{5.6, 5.8, 2.1},
				{19.4, 5.2, -26.1},
			},
		},
		{
			name: "Id",
			a: [][]float64{
				{1, 0, 0},
				{0, 1, 0},
				{0, 0, 1},
			},
		},
	} {

		a := NewDense(flatten(test.a))
		at := new(Dense)
		at.TCopy(a)

		qf := QR(a)

		rows, cols := a.Dims()

		lq := LQ(at)
		l := lq.L()
		lt := NewDense(rows, cols, nil)
		ltview := new(Dense)
		*ltview = *lt
		ltview.View(0, 0, cols, cols)
		ltview.TCopy(l)
		lq.ApplyQ(lt, true)

		qf.QR.TCopy(qf.QR)
		row := qf.QR.RowView(0)
		fmt.Println(row, lq.LQ.RowView(0))

		fmt.Println(lt, a)

		c.Check(a.EqualsApprox(lt, 1e-13), check.Equals, true, check.Commentf("Test %v: Q*R != A", test.name))
		c.Check(isLowerTriangular(l), check.Equals, true,
			check.Commentf("Test %v: L not lower triangular", test.name))
	}
}
