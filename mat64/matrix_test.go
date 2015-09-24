// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"fmt"
	"math"
	"testing"

	"gopkg.in/check.v1"
)

// Tests
func Test(t *testing.T) { check.TestingT(t) }

type S struct{}

var _ = check.Suite(&S{})

func leaksPanic(fn func()) (panicked bool) {
	defer func() {
		r := recover()
		panicked = r != nil
	}()
	Maybe(fn)
	return
}

func panics(fn func()) (panicked bool, message string) {
	defer func() {
		r := recover()
		panicked = r != nil
		message = fmt.Sprint(r)
	}()
	fn()
	return
}

func flatten(f [][]float64) (r, c int, d []float64) {
	r = len(f)
	if r == 0 {
		panic("bad test: no row")
	}
	c = len(f[0])
	d = make([]float64, 0, r*c)
	for _, row := range f {
		if len(row) != c {
			panic("bad test: ragged input")
		}
		d = append(d, row...)
	}
	return r, c, d
}

func unflatten(r, c int, d []float64) [][]float64 {
	m := make([][]float64, r)
	for i := 0; i < r; i++ {
		m[i] = d[i*c : (i+1)*c]
	}
	return m
}

func eye() *Dense {
	return NewDense(3, 3, []float64{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
	})
}

func (s *S) TestCol(c *check.C) {
	f := func(a Matrix) interface{} {
		_, c := a.Dims()
		ans := make([][]float64, c)
		for j := range ans {
			ans[j] = Col(nil, j, a)
		}
		return ans
	}
	denseComparison := func(a *Dense) interface{} {
		_, c := a.Dims()
		ans := make([][]float64, c)
		for j := range ans {
			ans[j] = Col(nil, j, a)
		}
		return ans
	}
	testOneInputFunc(c, "Col", f, denseComparison, sameAnswerF64SliceOfSlice, isAnyType, isAnySize)
	f = func(a Matrix) interface{} {
		r, c := a.Dims()
		ans := make([][]float64, c)
		for j := range ans {
			ans[j] = make([]float64, r)
			Col(ans[j], j, a)
		}
		return ans
	}
	testOneInputFunc(c, "Col", f, denseComparison, sameAnswerF64SliceOfSlice, isAnyType, isAnySize)
}

func (s *S) TestRow(c *check.C) {
	f := func(a Matrix) interface{} {
		r, _ := a.Dims()
		ans := make([][]float64, r)
		for i := range ans {
			ans[i] = Row(nil, i, a)
		}
		return ans
	}
	denseComparison := func(a *Dense) interface{} {
		r, _ := a.Dims()
		ans := make([][]float64, r)
		for i := range ans {
			ans[i] = Row(nil, i, a)
		}
		return ans
	}
	testOneInputFunc(c, "Row", f, denseComparison, sameAnswerF64SliceOfSlice, isAnyType, isAnySize)
	f = func(a Matrix) interface{} {
		r, c := a.Dims()
		ans := make([][]float64, r)
		for i := range ans {
			ans[i] = make([]float64, c)
			Row(ans[i], i, a)
		}
		return ans
	}
	testOneInputFunc(c, "Row", f, denseComparison, sameAnswerF64SliceOfSlice, isAnyType, isAnySize)
}

func (s *S) TestDot(c *check.C) {
	f := func(a, b Matrix) interface{} {
		return Dot(a, b)
	}
	denseComparison := func(a, b *Dense) interface{} {
		return Dot(a, b)
	}
	testTwoInputFunc(c, "Dot", f, denseComparison, sameAnswerFloatApprox, legalTypesAll, legalSizeSameRectangular)
}

func (s *S) TestEqual(c *check.C) {
	f := func(a, b Matrix) interface{} {
		return Equal(a, b)
	}
	denseComparison := func(a, b *Dense) interface{} {
		return Equal(a, b)
	}
	testTwoInputFunc(c, "Equal", f, denseComparison, sameAnswerBool, legalTypesAll, isAnySize2)
}

func (s *S) TestMax(c *check.C) {
	// A direct test of Max with *Dense arguments is in TestNewDense.
	f := func(a Matrix) interface{} {
		return Max(a)
	}
	denseComparison := func(a *Dense) interface{} {
		return Max(a)
	}
	testOneInputFunc(c, "Max", f, denseComparison, sameAnswerFloat, isAnyType, isAnySize)
}

func (s *S) TestMin(c *check.C) {
	// A direct test of Min with *Dense arguments is in TestNewDense.
	f := func(a Matrix) interface{} {
		return Min(a)
	}
	denseComparison := func(a *Dense) interface{} {
		return Min(a)
	}
	testOneInputFunc(c, "Min", f, denseComparison, sameAnswerFloat, isAnyType, isAnySize)
}

func (s *S) TestMaybe(c *check.C) {
	for i, test := range []struct {
		fn     func()
		panics bool
	}{
		{
			func() {},
			false,
		},
		{
			func() { panic("panic") },
			true,
		},
		{
			func() { panic(Error{"panic"}) },
			false,
		},
	} {
		c.Check(leaksPanic(test.fn), check.Equals, test.panics, check.Commentf("Test %d", i))
	}
}

func (s *S) TestNorm(c *check.C) {
	for i, test := range []struct {
		a    [][]float64
		ord  float64
		norm float64
	}{
		{
			a:    [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}},
			ord:  1,
			norm: 30,
		},
		{
			a:    [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}},
			ord:  2,
			norm: 25.495097567963924,
		},
		{
			a:    [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}},
			ord:  inf,
			norm: 33,
		},
		{
			a:    [][]float64{{1, -2, -2}, {-4, 5, 6}},
			ord:  1,
			norm: 8,
		},
		{
			a:    [][]float64{{1, -2, -2}, {-4, 5, 6}},
			ord:  inf,
			norm: 15,
		},
	} {
		a := NewDense(flatten(test.a))
		if math.Abs(Norm(a, test.ord)-test.norm) > 1e-14 {
			c.Errorf("Mismatch test %d: %v norm = %f", i, test.a, test.norm)
		}
	}

	f := func(a Matrix) interface{} {
		return Norm(a, 1)
	}
	denseComparison := func(a *Dense) interface{} {
		return Norm(a, 1)
	}
	testOneInputFunc(c, "Norm_1", f, denseComparison, sameAnswerFloatApprox, isAnyType, isAnySize)

	f = func(a Matrix) interface{} {
		return Norm(a, 2)
	}
	denseComparison = func(a *Dense) interface{} {
		return Norm(a, 2)
	}
	testOneInputFunc(c, "Norm_2", f, denseComparison, sameAnswerFloatApprox, isAnyType, isAnySize)

	f = func(a Matrix) interface{} {
		return Norm(a, math.Inf(1))
	}
	denseComparison = func(a *Dense) interface{} {
		return Norm(a, math.Inf(1))
	}
	testOneInputFunc(c, "Norm_inf", f, denseComparison, sameAnswerFloatApprox, isAnyType, isAnySize)
}

func (s *S) TestSum(c *check.C) {
	f := func(a Matrix) interface{} {
		return Sum(a)
	}
	denseComparison := func(a *Dense) interface{} {
		return Sum(a)
	}
	testOneInputFunc(c, "Sum", f, denseComparison, sameAnswerFloatApprox, isAnyType, isAnySize)
}

func (s *S) TestTrace(c *check.C) {
	for _, test := range []struct {
		a     *Dense
		trace float64
	}{
		{
			a:     NewDense(3, 3, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}),
			trace: 15,
		},
	} {
		trace := Trace(test.a)
		if trace != test.trace {
			c.Errorf("Trace mismatch. Want %v, got %v", test.trace, trace)
		}
	}
	f := func(a Matrix) interface{} {
		return Trace(a)
	}
	denseComparison := func(a *Dense) interface{} {
		return Trace(a)
	}
	testOneInputFunc(c, "Trace", f, denseComparison, sameAnswerFloat, isAnyType, isSquare)
}
