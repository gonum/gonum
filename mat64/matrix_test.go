// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"fmt"
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
