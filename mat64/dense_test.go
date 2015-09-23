// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"math"
	"math/rand"
	"testing"

	"github.com/gonum/blas/blas64"
	"github.com/gonum/floats"
	"gopkg.in/check.v1"
)

func asDense(d *Dense) Matrix {
	return d
}
func asBasicMatrix(d *Dense) Matrix {
	return (*basicMatrix)(d)
}
func asBasicVectorer(d *Dense) Matrix {
	return (*basicVectorer)(d)
}

func asBasicSymmetric(s *SymDense) Matrix {
	return (*basicSymmetric)(s)
}

func asBasicTriangular(t *TriDense) Triangular {
	return (*basicTriangular)(t)
}

func (s *S) TestNewDense(c *check.C) {
	for i, test := range []struct {
		a          []float64
		rows, cols int
		min, max   float64
		fro        float64
		mat        *Dense
	}{
		{
			[]float64{
				0, 0, 0,
				0, 0, 0,
				0, 0, 0,
			},
			3, 3,
			0, 0,
			0,
			&Dense{
				mat: blas64.General{
					Rows: 3, Cols: 3,
					Stride: 3,
					Data:   []float64{0, 0, 0, 0, 0, 0, 0, 0, 0},
				},
				capRows: 3, capCols: 3,
			},
		},
		{
			[]float64{
				1, 1, 1,
				1, 1, 1,
				1, 1, 1,
			},
			3, 3,
			1, 1,
			3,
			&Dense{
				mat: blas64.General{
					Rows: 3, Cols: 3,
					Stride: 3,
					Data:   []float64{1, 1, 1, 1, 1, 1, 1, 1, 1},
				},
				capRows: 3, capCols: 3,
			},
		},
		{
			[]float64{
				1, 0, 0,
				0, 1, 0,
				0, 0, 1,
			},
			3, 3,
			0, 1,
			1.7320508075688772,
			&Dense{
				mat: blas64.General{
					Rows: 3, Cols: 3,
					Stride: 3,
					Data:   []float64{1, 0, 0, 0, 1, 0, 0, 0, 1},
				},
				capRows: 3, capCols: 3,
			},
		},
		{
			[]float64{
				-1, 0, 0,
				0, -1, 0,
				0, 0, -1,
			},
			3, 3,
			-1, 0,
			1.7320508075688772,
			&Dense{
				mat: blas64.General{
					Rows: 3, Cols: 3,
					Stride: 3,
					Data:   []float64{-1, 0, 0, 0, -1, 0, 0, 0, -1},
				},
				capRows: 3, capCols: 3,
			},
		},
		{
			[]float64{
				1, 2, 3,
				4, 5, 6,
			},
			2, 3,
			1, 6,
			9.539392014169458,
			&Dense{
				mat: blas64.General{
					Rows: 2, Cols: 3,
					Stride: 3,
					Data:   []float64{1, 2, 3, 4, 5, 6},
				},
				capRows: 2, capCols: 3,
			},
		},
		{
			[]float64{
				1, 2,
				3, 4,
				5, 6,
			},
			3, 2,
			1, 6,
			9.539392014169458,
			&Dense{
				mat: blas64.General{
					Rows: 3, Cols: 2,
					Stride: 2,
					Data:   []float64{1, 2, 3, 4, 5, 6},
				},
				capRows: 3, capCols: 2,
			},
		},
	} {
		m := NewDense(test.rows, test.cols, test.a)
		rows, cols := m.Dims()
		c.Check(rows, check.Equals, test.rows, check.Commentf("Test %d", i))
		c.Check(cols, check.Equals, test.cols, check.Commentf("Test %d", i))
		c.Check(Min(m), check.Equals, test.min, check.Commentf("Test %d", i))
		c.Check(Max(m), check.Equals, test.max, check.Commentf("Test %d", i))
		c.Check(m.Norm(0), check.Equals, test.fro, check.Commentf("Test %d", i))
		c.Check(m, check.DeepEquals, test.mat, check.Commentf("Test %d", i))
		c.Check(Equal(m, test.mat), check.Equals, true, check.Commentf("Test %d", i))
	}
}

func (s *S) TestAtSet(c *check.C) {
	for test, af := range [][][]float64{
		{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, // even
		{{1, 2}, {4, 5}, {7, 8}},          // wide
		{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, //skinny
	} {
		m := NewDense(flatten(af))
		rows, cols := m.Dims()
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				c.Check(m.At(i, j), check.Equals, af[i][j], check.Commentf("At test %d", test))

				v := float64(i * j)
				m.Set(i, j, v)
				c.Check(m.At(i, j), check.Equals, v, check.Commentf("Set test %d", test))
			}
		}
		// Check access out of bounds fails
		c.Check(func() { m.At(rows, 0) }, check.PanicMatches, ErrRowAccess.Error(), check.Commentf("Test %d", test))
		c.Check(func() { m.At(rows+1, 0) }, check.PanicMatches, ErrRowAccess.Error(), check.Commentf("Test %d", test))
		c.Check(func() { m.At(0, cols) }, check.PanicMatches, ErrColAccess.Error(), check.Commentf("Test %d", test))
		c.Check(func() { m.At(0, cols+1) }, check.PanicMatches, ErrColAccess.Error(), check.Commentf("Test %d", test))
		c.Check(func() { m.At(-1, 0) }, check.PanicMatches, ErrRowAccess.Error(), check.Commentf("Test %d", test))
		c.Check(func() { m.At(0, -1) }, check.PanicMatches, ErrColAccess.Error(), check.Commentf("Test %d", test))

		// Check access out of bounds fails
		c.Check(func() { m.Set(rows, 0, 1.2) }, check.PanicMatches, ErrRowAccess.Error(), check.Commentf("Test %d", test))
		c.Check(func() { m.Set(rows+1, 0, 1.2) }, check.PanicMatches, ErrRowAccess.Error(), check.Commentf("Test %d", test))
		c.Check(func() { m.Set(0, cols, 1.2) }, check.PanicMatches, ErrColAccess.Error(), check.Commentf("Test %d", test))
		c.Check(func() { m.Set(0, cols+1, 1.2) }, check.PanicMatches, ErrColAccess.Error(), check.Commentf("Test %d", test))
		c.Check(func() { m.Set(-1, 0, 1.2) }, check.PanicMatches, ErrRowAccess.Error(), check.Commentf("Test %d", test))
		c.Check(func() { m.Set(0, -1, 1.2) }, check.PanicMatches, ErrColAccess.Error(), check.Commentf("Test %d", test))
	}
}

func (s *S) TestRowCol(c *check.C) {
	for i, af := range [][][]float64{
		{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
		{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}},
		{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}},
	} {
		a := NewDense(flatten(af))
		for ri, row := range af {
			c.Check(Row(nil, ri, a), check.DeepEquals, row, check.Commentf("Test %d", i))
		}
		for ci := range af[0] {
			col := make([]float64, a.mat.Rows)
			for j := range col {
				col[j] = float64(ci + 1 + j*a.mat.Cols)
			}
			c.Check(Col(nil, ci, a), check.DeepEquals, col, check.Commentf("Test %d", i))
		}
	}
}

func (s *S) TestSetRowColumn(c *check.C) {
	for _, as := range [][][]float64{
		{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
		{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}},
		{{1, 2, 3, 4}, {5, 6, 7, 8}, {9, 10, 11, 12}},
	} {
		for ri, row := range as {
			a := NewDense(flatten(as))
			t := &Dense{}
			t.Clone(a)
			a.SetRow(ri, make([]float64, a.mat.Cols))
			t.Sub(t, a)
			c.Check(t.Norm(0), check.Equals, floats.Norm(row, 2))
		}

		for ci := range as[0] {
			a := NewDense(flatten(as))
			t := &Dense{}
			t.Clone(a)
			a.SetCol(ci, make([]float64, a.mat.Rows))
			col := make([]float64, a.mat.Rows)
			for j := range col {
				col[j] = float64(ci + 1 + j*a.mat.Cols)
			}
			t.Sub(t, a)
			c.Check(t.Norm(0), check.Equals, floats.Norm(col, 2))
		}
	}
}

func (s *S) TestRowColView(c *check.C) {
	for _, test := range []struct {
		mat [][]float64
	}{
		{
			mat: [][]float64{
				{1, 2, 3, 4, 5},
				{6, 7, 8, 9, 10},
				{11, 12, 13, 14, 15},
				{16, 17, 18, 19, 20},
				{21, 22, 23, 24, 25},
			},
		},
		{
			mat: [][]float64{
				{1, 2, 3, 4},
				{6, 7, 8, 9},
				{11, 12, 13, 14},
				{16, 17, 18, 19},
				{21, 22, 23, 24},
			},
		},
		{
			mat: [][]float64{
				{1, 2, 3, 4, 5},
				{6, 7, 8, 9, 10},
				{11, 12, 13, 14, 15},
				{16, 17, 18, 19, 20},
			},
		},
	} {
		// This over cautious approach to building a matrix data
		// slice is to ensure that changes to flatten in the future
		// do not mask a regression to the issue identified in
		// gonum/matrix#110.
		rows, cols, flat := flatten(test.mat)
		m := NewDense(rows, cols, flat[:len(flat):len(flat)])

		c.Check(func() { m.RowView(-1) }, check.PanicMatches, ErrRowAccess.Error())
		c.Check(func() { m.RowView(rows) }, check.PanicMatches, ErrRowAccess.Error())
		c.Check(func() { m.ColView(-1) }, check.PanicMatches, ErrColAccess.Error())
		c.Check(func() { m.ColView(cols) }, check.PanicMatches, ErrColAccess.Error())

		for i := 0; i < rows; i++ {
			vr := m.RowView(i)
			c.Check(vr.Len(), check.Equals, cols)
			for j := 0; j < cols; j++ {
				c.Check(vr.At(j, 0), check.Equals, test.mat[i][j])
			}
		}
		for j := 0; j < cols; j++ {
			vr := m.ColView(j)
			c.Check(vr.Len(), check.Equals, rows)
			for i := 0; i < rows; i++ {
				c.Check(vr.At(i, 0), check.Equals, test.mat[i][j])
			}
		}
		m = m.View(1, 1, rows-2, cols-2).(*Dense)
		for i := 1; i < rows-1; i++ {
			vr := m.RowView(i - 1)
			c.Check(vr.Len(), check.Equals, cols-2)
			for j := 1; j < cols-1; j++ {
				c.Check(vr.At(j-1, 0), check.Equals, test.mat[i][j])
			}
		}
		for j := 1; j < cols-1; j++ {
			vr := m.ColView(j - 1)
			c.Check(vr.Len(), check.Equals, rows-2)
			for i := 1; i < rows-1; i++ {
				c.Check(vr.At(i-1, 0), check.Equals, test.mat[i][j])
			}
		}
	}
}

func (s *S) TestGrow(c *check.C) {
	m := &Dense{}
	m = m.Grow(10, 10).(*Dense)
	rows, cols := m.Dims()
	capRows, capCols := m.Caps()
	c.Check(rows, check.Equals, 10)
	c.Check(cols, check.Equals, 10)
	c.Check(capRows, check.Equals, 10)
	c.Check(capCols, check.Equals, 10)

	// Test grow within caps is in-place.
	m.Set(1, 1, 1)
	v := m.View(1, 1, 4, 4).(*Dense)
	c.Check(v.At(0, 0), check.Equals, m.At(1, 1))
	v = v.Grow(5, 5).(*Dense)
	c.Check(Equal(v, m.View(1, 1, 9, 9)), check.Equals, true)

	// Test grow bigger than caps copies.
	v = v.Grow(5, 5).(*Dense)
	c.Check(Equal(v.View(0, 0, 9, 9), m.View(1, 1, 9, 9)), check.Equals, true)
	v.Set(0, 0, 0)
	c.Check(Equal(v.View(0, 0, 9, 9), m.View(1, 1, 9, 9)), check.Equals, false)

	// Test grow uses existing data slice when matrix is zero size.
	v.Reset()
	p, l := &v.mat.Data[:1][0], cap(v.mat.Data)
	*p = 1
	v = v.Grow(5, 5).(*Dense)
	c.Check(&v.mat.Data[:1][0], check.Equals, p)
	c.Check(cap(v.mat.Data), check.Equals, l)
	c.Check(v.At(0, 0), check.Equals, 0.)
}

func (s *S) TestAdd(c *check.C) {
	for i, test := range []struct {
		a, b, r [][]float64
	}{
		{
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
		},
		{
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
			[][]float64{{2, 2, 2}, {2, 2, 2}, {2, 2, 2}},
		},
		{
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
			[][]float64{{2, 0, 0}, {0, 2, 0}, {0, 0, 2}},
		},
		{
			[][]float64{{-1, 0, 0}, {0, -1, 0}, {0, 0, -1}},
			[][]float64{{-1, 0, 0}, {0, -1, 0}, {0, 0, -1}},
			[][]float64{{-2, 0, 0}, {0, -2, 0}, {0, 0, -2}},
		},
		{
			[][]float64{{1, 2, 3}, {4, 5, 6}},
			[][]float64{{1, 2, 3}, {4, 5, 6}},
			[][]float64{{2, 4, 6}, {8, 10, 12}},
		},
	} {
		a := NewDense(flatten(test.a))
		b := NewDense(flatten(test.b))
		r := NewDense(flatten(test.r))

		var temp Dense
		temp.Add(a, b)
		c.Check(Equal(&temp, r), check.Equals, true, check.Commentf("Test %d: %v Add %v expect %v got %v",
			i, test.a, test.b, test.r, unflatten(temp.mat.Rows, temp.mat.Cols, temp.mat.Data)))

		zero(temp.mat.Data)
		temp.Add(a, b)
		c.Check(Equal(&temp, r), check.Equals, true, check.Commentf("Test %d: %v Add %v expect %v got %v",
			i, test.a, test.b, test.r, unflatten(temp.mat.Rows, temp.mat.Cols, temp.mat.Data)))

		// These probably warrant a better check and failure. They should never happen in the wild though.
		temp.mat.Data = nil
		c.Check(func() { temp.Add(a, b) }, check.PanicMatches, "runtime error: index out of range", check.Commentf("Test %d", i))

		a.Add(a, b)
		c.Check(Equal(a, r), check.Equals, true, check.Commentf("Test %d: %v Add %v expect %v got %v",
			i, test.a, test.b, test.r, unflatten(a.mat.Rows, a.mat.Cols, a.mat.Data)))
	}

	method := func(receiver, a, b Matrix) {
		rd := receiver.(Adder)
		rd.Add(a, b)
	}
	denseComparison := func(receiver, a, b *Dense) {
		receiver.Add(a, b)
	}
	testTwoInput(c, "Add", &Dense{}, method, denseComparison, legalTypesAll, legalSizeSameRectangular)
}

func (s *S) TestSub(c *check.C) {
	for i, test := range []struct {
		a, b, r [][]float64
	}{
		{
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
		},
		{
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
		},
		{
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
		},
		{
			[][]float64{{-1, 0, 0}, {0, -1, 0}, {0, 0, -1}},
			[][]float64{{-1, 0, 0}, {0, -1, 0}, {0, 0, -1}},
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
		},
		{
			[][]float64{{1, 2, 3}, {4, 5, 6}},
			[][]float64{{1, 2, 3}, {4, 5, 6}},
			[][]float64{{0, 0, 0}, {0, 0, 0}},
		},
	} {
		a := NewDense(flatten(test.a))
		b := NewDense(flatten(test.b))
		r := NewDense(flatten(test.r))

		var temp Dense
		temp.Sub(a, b)
		c.Check(Equal(&temp, r), check.Equals, true, check.Commentf("Test %d: %v Sub %v expect %v got %v",
			i, test.a, test.b, test.r, unflatten(temp.mat.Rows, temp.mat.Cols, temp.mat.Data)))

		zero(temp.mat.Data)
		temp.Sub(a, b)
		c.Check(Equal(&temp, r), check.Equals, true, check.Commentf("Test %d: %v Sub %v expect %v got %v",
			i, test.a, test.b, test.r, unflatten(temp.mat.Rows, temp.mat.Cols, temp.mat.Data)))

		// These probably warrant a better check and failure. They should never happen in the wild though.
		temp.mat.Data = nil
		c.Check(func() { temp.Sub(a, b) }, check.PanicMatches, "runtime error: index out of range", check.Commentf("Test %d", i))

		a.Sub(a, b)
		c.Check(Equal(a, r), check.Equals, true, check.Commentf("Test %d: %v Sub %v expect %v got %v",
			i, test.a, test.b, test.r, unflatten(a.mat.Rows, a.mat.Cols, a.mat.Data)))
	}
}

func (s *S) TestMulElem(c *check.C) {
	for i, test := range []struct {
		a, b, r [][]float64
	}{
		{
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
		},
		{
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
		},
		{
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
		},
		{
			[][]float64{{-1, 0, 0}, {0, -1, 0}, {0, 0, -1}},
			[][]float64{{-1, 0, 0}, {0, -1, 0}, {0, 0, -1}},
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
		},
		{
			[][]float64{{1, 2, 3}, {4, 5, 6}},
			[][]float64{{1, 2, 3}, {4, 5, 6}},
			[][]float64{{1, 4, 9}, {16, 25, 36}},
		},
	} {
		a := NewDense(flatten(test.a))
		b := NewDense(flatten(test.b))
		r := NewDense(flatten(test.r))

		var temp Dense
		temp.MulElem(a, b)
		c.Check(Equal(&temp, r), check.Equals, true, check.Commentf("Test %d: %v MulElem %v expect %v got %v",
			i, test.a, test.b, test.r, unflatten(temp.mat.Rows, temp.mat.Cols, temp.mat.Data)))

		zero(temp.mat.Data)
		temp.MulElem(a, b)
		c.Check(Equal(&temp, r), check.Equals, true, check.Commentf("Test %d: %v MulElem %v expect %v got %v",
			i, test.a, test.b, test.r, unflatten(temp.mat.Rows, temp.mat.Cols, temp.mat.Data)))

		// These probably warrant a better check and failure. They should never happen in the wild though.
		temp.mat.Data = nil
		c.Check(func() { temp.MulElem(a, b) }, check.PanicMatches, "runtime error: index out of range", check.Commentf("Test %d", i))

		a.MulElem(a, b)
		c.Check(Equal(a, r), check.Equals, true, check.Commentf("Test %d: %v MulElem %v expect %v got %v",
			i, test.a, test.b, test.r, unflatten(a.mat.Rows, a.mat.Cols, a.mat.Data)))
	}
}

// A comparison that treats NaNs as equal, for testing.
func (m *Dense) same(b Matrix) bool {
	br, bc := b.Dims()
	if br != m.mat.Rows || bc != m.mat.Cols {
		return false
	}
	for r := 0; r < br; r++ {
		for c := 0; c < bc; c++ {
			if av, bv := m.At(r, c), b.At(r, c); av != bv && !(math.IsNaN(av) && math.IsNaN(bv)) {
				return false
			}
		}
	}
	return true
}

func (s *S) TestDivElem(c *check.C) {
	for i, test := range []struct {
		a, b, r [][]float64
	}{
		{
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			[][]float64{{math.Inf(1), math.NaN(), math.NaN()}, {math.NaN(), math.Inf(1), math.NaN()}, {math.NaN(), math.NaN(), math.Inf(1)}},
		},
		{
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
		},
		{
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
			[][]float64{{1, math.NaN(), math.NaN()}, {math.NaN(), 1, math.NaN()}, {math.NaN(), math.NaN(), 1}},
		},
		{
			[][]float64{{-1, 0, 0}, {0, -1, 0}, {0, 0, -1}},
			[][]float64{{-1, 0, 0}, {0, -1, 0}, {0, 0, -1}},
			[][]float64{{1, math.NaN(), math.NaN()}, {math.NaN(), 1, math.NaN()}, {math.NaN(), math.NaN(), 1}},
		},
		{
			[][]float64{{1, 2, 3}, {4, 5, 6}},
			[][]float64{{1, 2, 3}, {4, 5, 6}},
			[][]float64{{1, 1, 1}, {1, 1, 1}},
		},
	} {
		a := NewDense(flatten(test.a))
		b := NewDense(flatten(test.b))
		r := NewDense(flatten(test.r))

		var temp Dense
		temp.DivElem(a, b)
		c.Check(temp.same(r), check.Equals, true, check.Commentf("Test %d: %v DivElem %v expect %v got %v",
			i, test.a, test.b, test.r, unflatten(temp.mat.Rows, temp.mat.Cols, temp.mat.Data)))

		zero(temp.mat.Data)
		temp.DivElem(a, b)
		c.Check(temp.same(r), check.Equals, true, check.Commentf("Test %d: %v DivElem %v expect %v got %v",
			i, test.a, test.b, test.r, unflatten(temp.mat.Rows, temp.mat.Cols, temp.mat.Data)))

		// These probably warrant a better check and failure. They should never happen in the wild though.
		temp.mat.Data = nil
		c.Check(func() { temp.DivElem(a, b) }, check.PanicMatches, "runtime error: index out of range", check.Commentf("Test %d", i))

		a.DivElem(a, b)
		c.Check(a.same(r), check.Equals, true, check.Commentf("Test %d: %v DivElem %v expect %v got %v",
			i, test.a, test.b, test.r, unflatten(a.mat.Rows, a.mat.Cols, a.mat.Data)))
	}
}

func (s *S) TestMul(c *check.C) {
	for i, test := range []struct {
		a, b, r [][]float64
	}{
		{
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
		},
		{
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
			[][]float64{{3, 3, 3}, {3, 3, 3}, {3, 3, 3}},
		},
		{
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
		},
		{
			[][]float64{{-1, 0, 0}, {0, -1, 0}, {0, 0, -1}},
			[][]float64{{-1, 0, 0}, {0, -1, 0}, {0, 0, -1}},
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
		},
		{
			[][]float64{{1, 2, 3}, {4, 5, 6}},
			[][]float64{{1, 2}, {3, 4}, {5, 6}},
			[][]float64{{22, 28}, {49, 64}},
		},
		{
			[][]float64{{0, 1, 1}, {0, 1, 1}, {0, 1, 1}},
			[][]float64{{0, 1, 1}, {0, 1, 1}, {0, 1, 1}},
			[][]float64{{0, 2, 2}, {0, 2, 2}, {0, 2, 2}},
		},
	} {
		a := NewDense(flatten(test.a))
		b := NewDense(flatten(test.b))
		r := NewDense(flatten(test.r))

		var temp Dense
		temp.Mul(a, b)
		c.Check(Equal(&temp, r), check.Equals, true, check.Commentf("Test %d: %v Mul %v expect %v got %v",
			i, test.a, test.b, test.r, unflatten(temp.mat.Rows, temp.mat.Cols, temp.mat.Data)))

		zero(temp.mat.Data)
		temp.Mul(a, b)
		c.Check(Equal(&temp, r), check.Equals, true, check.Commentf("Test %d: %v Mul %v expect %v got %v",
			i, test.a, test.b, test.r, unflatten(a.mat.Rows, a.mat.Cols, a.mat.Data)))

		// These probably warrant a better check and failure. They should never happen in the wild though.
		temp.mat.Data = nil
		c.Check(func() { temp.Mul(a, b) }, check.PanicMatches, "blas: index of c out of range", check.Commentf("Test %d", i))
	}

	method := func(receiver, a, b Matrix) {
		rd := receiver.(Muler)
		rd.Mul(a, b)
	}
	denseComparison := func(receiver, a, b *Dense) {
		receiver.Mul(a, b)
	}
	legalSizeMul := func(ar, ac, br, bc int) bool {
		return ac == br
	}
	testTwoInput(c, "Mul", &Dense{}, method, denseComparison, legalTypesAll, legalSizeMul)
}

func randDense(size int, rho float64, rnd func() float64) (*Dense, error) {
	if size == 0 {
		return nil, ErrZeroLength
	}
	d := &Dense{
		mat: blas64.General{
			Rows: size, Cols: size, Stride: size,
			Data: make([]float64, size*size),
		},
		capRows: size, capCols: size,
	}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if rand.Float64() < rho {
				d.Set(i, j, rnd())
			}
		}
	}
	return d, nil
}

func (s *S) TestExp(c *check.C) {
	for i, t := range []struct {
		a    [][]float64
		want [][]float64
		mod  func(*Dense)
	}{
		{
			a:    [][]float64{{-49, 24}, {-64, 31}},
			want: [][]float64{{-0.7357587581474017, 0.5518190996594223}, {-1.4715175990917921, 1.103638240717339}},
		},
		{
			a:    [][]float64{{-49, 24}, {-64, 31}},
			want: [][]float64{{-0.7357587581474017, 0.5518190996594223}, {-1.4715175990917921, 1.103638240717339}},
			mod: func(a *Dense) {
				d := make([]float64, 100)
				for i := range d {
					d[i] = math.NaN()
				}
				*a = *NewDense(10, 10, d).View(1, 1, 2, 2).(*Dense)
			},
		},
		{
			a:    [][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
			want: [][]float64{{2.71828182845905, 0, 0}, {0, 2.71828182845905, 0}, {0, 0, 2.71828182845905}},
		},
	} {
		var got Dense
		if t.mod != nil {
			t.mod(&got)
		}
		got.Exp(NewDense(flatten(t.a)))
		c.Check(EqualApprox(&got, NewDense(flatten(t.want)), 1e-12), check.Equals, true, check.Commentf("Test %d", i))
	}
}

func (s *S) TestPow(c *check.C) {
	for i, t := range []struct {
		a    [][]float64
		n    int
		mod  func(*Dense)
		want [][]float64
	}{
		{
			a:    [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			n:    0,
			want: [][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
		},
		{
			a:    [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			n:    0,
			want: [][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
			mod: func(a *Dense) {
				d := make([]float64, 100)
				for i := range d {
					d[i] = math.NaN()
				}
				*a = *NewDense(10, 10, d).View(1, 1, 3, 3).(*Dense)
			},
		},
		{
			a:    [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			n:    1,
			want: [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
		},
		{
			a:    [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			n:    1,
			want: [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			mod: func(a *Dense) {
				d := make([]float64, 100)
				for i := range d {
					d[i] = math.NaN()
				}
				*a = *NewDense(10, 10, d).View(1, 1, 3, 3).(*Dense)
			},
		},
		{
			a:    [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			n:    2,
			want: [][]float64{{30, 36, 42}, {66, 81, 96}, {102, 126, 150}},
		},
		{
			a:    [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			n:    2,
			want: [][]float64{{30, 36, 42}, {66, 81, 96}, {102, 126, 150}},
			mod: func(a *Dense) {
				d := make([]float64, 100)
				for i := range d {
					d[i] = math.NaN()
				}
				*a = *NewDense(10, 10, d).View(1, 1, 3, 3).(*Dense)
			},
		},
		{
			a:    [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			n:    3,
			want: [][]float64{{468, 576, 684}, {1062, 1305, 1548}, {1656, 2034, 2412}},
		},
		{
			a:    [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			n:    3,
			want: [][]float64{{468, 576, 684}, {1062, 1305, 1548}, {1656, 2034, 2412}},
			mod: func(a *Dense) {
				d := make([]float64, 100)
				for i := range d {
					d[i] = math.NaN()
				}
				*a = *NewDense(10, 10, d).View(1, 1, 3, 3).(*Dense)
			},
		},
	} {
		var got Dense
		if t.mod != nil {
			t.mod(&got)
		}
		got.Pow(NewDense(flatten(t.a)), t.n)
		c.Check(Equal(&got, NewDense(flatten(t.want))), check.Equals, true, check.Commentf("Test %d", i))
	}
}

func (s *S) TestPowN(c *check.C) {
	for i, t := range []struct {
		a    [][]float64
		mod  func(*Dense)
		want [][]float64
	}{
		{
			a: [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
		},
		{
			a: [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}},
			mod: func(a *Dense) {
				d := make([]float64, 100)
				for i := range d {
					d[i] = math.NaN()
				}
				*a = *NewDense(10, 10, d).View(1, 1, 3, 3).(*Dense)
			},
		},
	} {
		for n := 1; n <= 14; n++ {
			var got, want Dense
			if t.mod != nil {
				t.mod(&got)
			}
			got.Pow(NewDense(flatten(t.a)), n)
			want.iterativePow(NewDense(flatten(t.a)), n)
			c.Check(Equal(&got, &want), check.Equals, true, check.Commentf("Test %d", i))
		}
	}
}

func (m *Dense) iterativePow(a Matrix, n int) {
	m.Clone(a)
	for i := 1; i < n; i++ {
		m.Mul(m, a)
	}
}

func (s *S) TestCloneT(c *check.C) {
	for i, test := range []struct {
		a, t [][]float64
	}{
		{
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
		},
		{
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
		},
		{
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
		},
		{
			[][]float64{{-1, 0, 0}, {0, -1, 0}, {0, 0, -1}},
			[][]float64{{-1, 0, 0}, {0, -1, 0}, {0, 0, -1}},
		},
		{
			[][]float64{{1, 2, 3}, {4, 5, 6}},
			[][]float64{{1, 4}, {2, 5}, {3, 6}},
		},
	} {
		a := NewDense(flatten(test.a))
		t := NewDense(flatten(test.t))

		var r, rr Dense

		r.Clone(a.T())
		c.Check(Equal(&r, t), check.Equals, true, check.Commentf("Test %d: %v transpose = %v", i, test.a, test.t))

		rr.Clone(r.T())
		c.Check(Equal(&rr, a), check.Equals, true, check.Commentf("Test %d: %v transpose = %v", i, test.a, test.t))

		zero(r.mat.Data)
		r.Clone(a.T())
		c.Check(Equal(&r, t), check.Equals, true, check.Commentf("Test %d: %v transpose = %v", i, test.a, test.t))

		zero(rr.mat.Data)
		rr.Clone(r.T())
		c.Check(Equal(&rr, a), check.Equals, true, check.Commentf("Test %d: %v transpose = %v", i, test.a, test.t))
	}
}

func (s *S) TestCopyT(c *check.C) {
	for i, test := range []struct {
		a, t [][]float64
	}{
		{
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
		},
		{
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
		},
		{
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
		},
		{
			[][]float64{{-1, 0, 0}, {0, -1, 0}, {0, 0, -1}},
			[][]float64{{-1, 0, 0}, {0, -1, 0}, {0, 0, -1}},
		},
		{
			[][]float64{{1, 2, 3}, {4, 5, 6}},
			[][]float64{{1, 4}, {2, 5}, {3, 6}},
		},
	} {
		a := NewDense(flatten(test.a))
		t := NewDense(flatten(test.t))

		ar, ac := a.Dims()
		r := NewDense(ac, ar, nil)
		rr := NewDense(ar, ac, nil)

		r.Copy(a.T())
		c.Check(Equal(r, t), check.Equals, true, check.Commentf("Test %d: %v transpose = %v", i, test.a, test.t))

		rr.Copy(r.T())
		c.Check(Equal(rr, a), check.Equals, true, check.Commentf("Test %d: %v transpose = %v", i, test.a, test.t))

		zero(r.mat.Data)
		r.Copy(a.T())
		c.Check(Equal(r, t), check.Equals, true, check.Commentf("Test %d: %v transpose = %v", i, test.a, test.t))

		zero(rr.mat.Data)
		rr.Copy(r.T())
		c.Check(Equal(rr, a), check.Equals, true, check.Commentf("Test %d: %v transpose = %v", i, test.a, test.t))
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
			ord:  0,
			norm: 25.49509756796392,
		},
		{
			a:    [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}},
			ord:  1,
			norm: 30,
		},
		{
			a:    [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}},
			ord:  -1,
			norm: 22,
		},
		{
			a:    [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}},
			ord:  2,
			norm: 25.46240743603639,
		},
		{
			a:    [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}},
			ord:  -2,
			norm: 9.013990486603544e-16,
		},
		{
			a:    [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}},
			ord:  inf,
			norm: 33,
		},
		{
			a:    [][]float64{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}, {10, 11, 12}},
			ord:  -inf,
			norm: 6,
		},
		{
			a:    [][]float64{{1, -2, -2}, {-4, 5, 6}},
			ord:  1,
			norm: 8,
		},
		{
			a:    [][]float64{{1, -2, -2}, {-4, 5, 6}},
			ord:  -1,
			norm: 5,
		},
		{
			a:    [][]float64{{1, -2, -2}, {-4, 5, 6}},
			ord:  inf,
			norm: 15,
		},
		{
			a:    [][]float64{{1, -2, -2}, {-4, 5, 6}},
			ord:  -inf,
			norm: 5,
		},
	} {
		a := NewDense(flatten(test.a))
		c.Check(a.Norm(test.ord), check.Equals, test.norm, check.Commentf("Test %d: %v norm = %f", i, test.a, test.norm))
	}
}

func identity(r, c int, v float64) float64 { return v }

func (s *S) TestApply(c *check.C) {
	for i, test := range []struct {
		a, t [][]float64
		fn   func(r, c int, v float64) float64
	}{
		{
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			identity,
		},
		{
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
			identity,
		},
		{
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
			identity,
		},
		{
			[][]float64{{-1, 0, 0}, {0, -1, 0}, {0, 0, -1}},
			[][]float64{{-1, 0, 0}, {0, -1, 0}, {0, 0, -1}},
			identity,
		},
		{
			[][]float64{{1, 2, 3}, {4, 5, 6}},
			[][]float64{{1, 2, 3}, {4, 5, 6}},
			identity,
		},
		{
			[][]float64{{1, 2, 3}, {4, 5, 6}},
			[][]float64{{2, 4, 6}, {8, 10, 12}},
			func(r, c int, v float64) float64 { return v * 2 },
		},
		{
			[][]float64{{1, 2, 3}, {4, 5, 6}},
			[][]float64{{0, 2, 0}, {0, 5, 0}},
			func(r, c int, v float64) float64 {
				if c == 1 {
					return v
				}
				return 0
			},
		},
		{
			[][]float64{{1, 2, 3}, {4, 5, 6}},
			[][]float64{{0, 0, 0}, {4, 5, 6}},
			func(r, c int, v float64) float64 {
				if r == 1 {
					return v
				}
				return 0
			},
		},
	} {
		a := NewDense(flatten(test.a))
		t := NewDense(flatten(test.t))

		var r Dense

		r.Apply(test.fn, a)
		c.Check(Equal(&r, t), check.Equals, true, check.Commentf("Test %d: obtained %v expect: %v", i, r.mat.Data, t.mat.Data))

		a.Apply(test.fn, a)
		c.Check(Equal(a, t), check.Equals, true, check.Commentf("Test %d: obtained %v expect: %v", i, a.mat.Data, t.mat.Data))
	}
}

func (s *S) TestClone(c *check.C) {
	for i, test := range []struct {
		a    [][]float64
		i, j int
		v    float64
	}{
		{
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			1, 1,
			1,
		},
		{
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
			0, 0,
			0,
		},
	} {
		a := NewDense(flatten(test.a))
		b := *a
		a.Clone(a)
		a.Set(test.i, test.j, test.v)

		c.Check(Equal(&b, a), check.Equals, false, check.Commentf("Test %d: %v cloned and altered = %v", i, a, &b))
	}
}

// TODO(kortschak) Roll this into testOneInput when it exists.
func (s *S) TestCopyPanic(c *check.C) {
	for _, a := range []*Dense{
		&Dense{},
		&Dense{mat: blas64.General{Rows: 1}},
		&Dense{mat: blas64.General{Cols: 1}},
	} {
		var rows, cols int
		m := NewDense(1, 1, nil)
		panicked, _ := panics(func() { rows, cols = m.Copy(a) })
		c.Check(panicked, check.Equals, false)
		c.Check(rows, check.Equals, 0)
		c.Check(cols, check.Equals, 0)
	}
}

func (s *S) TestStack(c *check.C) {
	for i, test := range []struct {
		a, b, e [][]float64
	}{
		{
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
		},
		{
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}, {1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
		},
		{
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
			[][]float64{{0, 1, 0}, {0, 0, 1}, {1, 0, 0}},
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}, {0, 1, 0}, {0, 0, 1}, {1, 0, 0}},
		},
	} {
		a := NewDense(flatten(test.a))
		b := NewDense(flatten(test.b))

		var s Dense
		s.Stack(a, b)

		c.Check(Equal(&s, NewDense(flatten(test.e))), check.Equals, true, check.Commentf("Test %d: %v stack %v = %v", i, a, b, s))
	}
}

func (s *S) TestAugment(c *check.C) {
	for i, test := range []struct {
		a, b, e [][]float64
	}{
		{
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			[][]float64{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
			[][]float64{{0, 0, 0, 0, 0, 0}, {0, 0, 0, 0, 0, 0}, {0, 0, 0, 0, 0, 0}},
		},
		{
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
			[][]float64{{1, 1, 1}, {1, 1, 1}, {1, 1, 1}},
			[][]float64{{1, 1, 1, 1, 1, 1}, {1, 1, 1, 1, 1, 1}, {1, 1, 1, 1, 1, 1}},
		},
		{
			[][]float64{{1, 0, 0}, {0, 1, 0}, {0, 0, 1}},
			[][]float64{{0, 1, 0}, {0, 0, 1}, {1, 0, 0}},
			[][]float64{{1, 0, 0, 0, 1, 0}, {0, 1, 0, 0, 0, 1}, {0, 0, 1, 1, 0, 0}},
		},
	} {
		a := NewDense(flatten(test.a))
		b := NewDense(flatten(test.b))

		var s Dense
		s.Augment(a, b)

		c.Check(Equal(&s, NewDense(flatten(test.e))), check.Equals, true, check.Commentf("Test %d: %v stack %v = %v", i, a, b, s))
	}
}

func (s *S) TestRankOne(c *check.C) {
	for i, test := range []struct {
		x     []float64
		y     []float64
		m     [][]float64
		alpha float64
	}{
		{
			x:     []float64{5},
			y:     []float64{10},
			m:     [][]float64{{2}},
			alpha: -3,
		},
		{
			x:     []float64{5, 6, 1},
			y:     []float64{10},
			m:     [][]float64{{2}, {-3}, {5}},
			alpha: -3,
		},

		{
			x:     []float64{5},
			y:     []float64{10, 15, 8},
			m:     [][]float64{{2, -3, 5}},
			alpha: -3,
		},
		{
			x: []float64{1, 5},
			y: []float64{10, 15},
			m: [][]float64{
				{2, -3},
				{4, -1},
			},
			alpha: -3,
		},
		{
			x: []float64{2, 3, 9},
			y: []float64{8, 9},
			m: [][]float64{
				{2, 3},
				{4, 5},
				{6, 7},
			},
			alpha: -3,
		},
		{
			x: []float64{2, 3},
			y: []float64{8, 9, 9},
			m: [][]float64{
				{2, 3, 6},
				{4, 5, 7},
			},
			alpha: -3,
		},
	} {
		want := &Dense{}
		xm := NewDense(len(test.x), 1, test.x)
		ym := NewDense(1, len(test.y), test.y)

		want.Mul(xm, ym)
		want.Scale(test.alpha, want)
		want.Add(want, NewDense(flatten(test.m)))

		a := NewDense(flatten(test.m))
		m := &Dense{}
		// Check with a new matrix
		m.RankOne(a, test.alpha, NewVector(len(test.x), test.x), NewVector(len(test.y), test.y))
		c.Check(Equal(m, want), check.Equals, true, check.Commentf("Test %v. Want %v, Got %v", i, want, m))
		// Check with the same matrix
		a.RankOne(a, test.alpha, NewVector(len(test.x), test.x), NewVector(len(test.y), test.y))
		c.Check(Equal(a, want), check.Equals, true, check.Commentf("Test %v. Want %v, Got %v", i, want, m))
	}
}

func (s *S) TestOuter(c *check.C) {
	for i, test := range []struct {
		x []float64
		y []float64
	}{
		{
			x: []float64{5},
			y: []float64{10},
		},
		{
			x: []float64{5, 6, 1},
			y: []float64{10},
		},

		{
			x: []float64{5},
			y: []float64{10, 15, 8},
		},
		{
			x: []float64{1, 5},
			y: []float64{10, 15},
		},
		{
			x: []float64{2, 3, 9},
			y: []float64{8, 9},
		},
		{
			x: []float64{2, 3},
			y: []float64{8, 9, 9},
		},
	} {
		want := &Dense{}
		xm := NewDense(len(test.x), 1, test.x)
		ym := NewDense(1, len(test.y), test.y)

		want.Mul(xm, ym)

		var m Dense
		// Check with a new matrix
		m.Outer(NewVector(len(test.x), test.x), NewVector(len(test.y), test.y))
		c.Check(Equal(&m, want), check.Equals, true, check.Commentf("Test %v. Want %v, Got %v", i, want, m))
		// Check with the same matrix
		m.Outer(NewVector(len(test.x), test.x), NewVector(len(test.y), test.y))
		c.Check(Equal(&m, want), check.Equals, true, check.Commentf("Test %v. Want %v, Got %v", i, want, m))
	}
}

var (
	wd *Dense
)

func BenchmarkMulDense100Half(b *testing.B)        { denseMulBench(b, 100, 0.5) }
func BenchmarkMulDense100Tenth(b *testing.B)       { denseMulBench(b, 100, 0.1) }
func BenchmarkMulDense1000Half(b *testing.B)       { denseMulBench(b, 1000, 0.5) }
func BenchmarkMulDense1000Tenth(b *testing.B)      { denseMulBench(b, 1000, 0.1) }
func BenchmarkMulDense1000Hundredth(b *testing.B)  { denseMulBench(b, 1000, 0.01) }
func BenchmarkMulDense1000Thousandth(b *testing.B) { denseMulBench(b, 1000, 0.001) }
func denseMulBench(b *testing.B, size int, rho float64) {
	b.StopTimer()
	a, _ := randDense(size, rho, rand.NormFloat64)
	d, _ := randDense(size, rho, rand.NormFloat64)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var n Dense
		n.Mul(a, d)
		wd = &n
	}
}

func BenchmarkPreMulDense100Half(b *testing.B)        { densePreMulBench(b, 100, 0.5) }
func BenchmarkPreMulDense100Tenth(b *testing.B)       { densePreMulBench(b, 100, 0.1) }
func BenchmarkPreMulDense1000Half(b *testing.B)       { densePreMulBench(b, 1000, 0.5) }
func BenchmarkPreMulDense1000Tenth(b *testing.B)      { densePreMulBench(b, 1000, 0.1) }
func BenchmarkPreMulDense1000Hundredth(b *testing.B)  { densePreMulBench(b, 1000, 0.01) }
func BenchmarkPreMulDense1000Thousandth(b *testing.B) { densePreMulBench(b, 1000, 0.001) }
func densePreMulBench(b *testing.B, size int, rho float64) {
	b.StopTimer()
	a, _ := randDense(size, rho, rand.NormFloat64)
	d, _ := randDense(size, rho, rand.NormFloat64)
	wd = NewDense(size, size, nil)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		wd.Mul(a, d)
	}
}

func BenchmarkRow10(b *testing.B)   { rowBench(b, 10) }
func BenchmarkRow100(b *testing.B)  { rowBench(b, 100) }
func BenchmarkRow1000(b *testing.B) { rowBench(b, 1000) }

func rowBench(b *testing.B, size int) {
	a, _ := randDense(size, 1, rand.NormFloat64)
	_, c := a.Dims()
	dst := make([]float64, c)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Row(dst, 0, a)
	}
}

func BenchmarkExp10(b *testing.B)   { expBench(b, 10) }
func BenchmarkExp100(b *testing.B)  { expBench(b, 100) }
func BenchmarkExp1000(b *testing.B) { expBench(b, 1000) }

func expBench(b *testing.B, size int) {
	a, _ := randDense(size, 1, rand.NormFloat64)

	b.ResetTimer()
	var m Dense
	for i := 0; i < b.N; i++ {
		m.Exp(a)
	}
}

func BenchmarkPow10_3(b *testing.B)   { powBench(b, 10, 3) }
func BenchmarkPow100_3(b *testing.B)  { powBench(b, 100, 3) }
func BenchmarkPow1000_3(b *testing.B) { powBench(b, 1000, 3) }
func BenchmarkPow10_4(b *testing.B)   { powBench(b, 10, 4) }
func BenchmarkPow100_4(b *testing.B)  { powBench(b, 100, 4) }
func BenchmarkPow1000_4(b *testing.B) { powBench(b, 1000, 4) }
func BenchmarkPow10_5(b *testing.B)   { powBench(b, 10, 5) }
func BenchmarkPow100_5(b *testing.B)  { powBench(b, 100, 5) }
func BenchmarkPow1000_5(b *testing.B) { powBench(b, 1000, 5) }
func BenchmarkPow10_6(b *testing.B)   { powBench(b, 10, 6) }
func BenchmarkPow100_6(b *testing.B)  { powBench(b, 100, 6) }
func BenchmarkPow1000_6(b *testing.B) { powBench(b, 1000, 6) }
func BenchmarkPow10_7(b *testing.B)   { powBench(b, 10, 7) }
func BenchmarkPow100_7(b *testing.B)  { powBench(b, 100, 7) }
func BenchmarkPow1000_7(b *testing.B) { powBench(b, 1000, 7) }
func BenchmarkPow10_8(b *testing.B)   { powBench(b, 10, 8) }
func BenchmarkPow100_8(b *testing.B)  { powBench(b, 100, 8) }
func BenchmarkPow1000_8(b *testing.B) { powBench(b, 1000, 8) }
func BenchmarkPow10_9(b *testing.B)   { powBench(b, 10, 9) }
func BenchmarkPow100_9(b *testing.B)  { powBench(b, 100, 9) }
func BenchmarkPow1000_9(b *testing.B) { powBench(b, 1000, 9) }

func powBench(b *testing.B, size, n int) {
	a, _ := randDense(size, 1, rand.NormFloat64)

	b.ResetTimer()
	var m Dense
	for i := 0; i < b.N; i++ {
		m.Pow(a, n)
	}
}

func BenchmarkMulTransDense100Half(b *testing.B)        { denseMulTransBench(b, 100, 0.5) }
func BenchmarkMulTransDense100Tenth(b *testing.B)       { denseMulTransBench(b, 100, 0.1) }
func BenchmarkMulTransDense1000Half(b *testing.B)       { denseMulTransBench(b, 1000, 0.5) }
func BenchmarkMulTransDense1000Tenth(b *testing.B)      { denseMulTransBench(b, 1000, 0.1) }
func BenchmarkMulTransDense1000Hundredth(b *testing.B)  { denseMulTransBench(b, 1000, 0.01) }
func BenchmarkMulTransDense1000Thousandth(b *testing.B) { denseMulTransBench(b, 1000, 0.001) }
func denseMulTransBench(b *testing.B, size int, rho float64) {
	b.StopTimer()
	a, _ := randDense(size, rho, rand.NormFloat64)
	d, _ := randDense(size, rho, rand.NormFloat64)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var n Dense
		n.Mul(a, d.T())
		wd = &n
	}
}

func BenchmarkMulTransDenseSym100Half(b *testing.B)        { denseMulTransSymBench(b, 100, 0.5) }
func BenchmarkMulTransDenseSym100Tenth(b *testing.B)       { denseMulTransSymBench(b, 100, 0.1) }
func BenchmarkMulTransDenseSym1000Half(b *testing.B)       { denseMulTransSymBench(b, 1000, 0.5) }
func BenchmarkMulTransDenseSym1000Tenth(b *testing.B)      { denseMulTransSymBench(b, 1000, 0.1) }
func BenchmarkMulTransDenseSym1000Hundredth(b *testing.B)  { denseMulTransSymBench(b, 1000, 0.01) }
func BenchmarkMulTransDenseSym1000Thousandth(b *testing.B) { denseMulTransSymBench(b, 1000, 0.001) }
func denseMulTransSymBench(b *testing.B, size int, rho float64) {
	b.StopTimer()
	a, _ := randDense(size, rho, rand.NormFloat64)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		var n Dense
		n.Mul(a, a.T())
		wd = &n
	}
}
