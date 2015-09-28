// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"github.com/gonum/floats"
)

func TestNewSymmetric(t *testing.T) {
	for i, test := range []struct {
		data []float64
		n    int
		mat  *SymDense
	}{
		{
			data: []float64{
				1, 2, 3,
				4, 5, 6,
				7, 8, 9,
			},
			n: 3,
			mat: &SymDense{blas64.Symmetric{
				N:      3,
				Stride: 3,
				Uplo:   blas.Upper,
				Data:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
			}},
		},
	} {
		sym := NewSymDense(test.n, test.data)
		rows, cols := sym.Dims()

		if rows != test.n {
			t.Errorf("unexpected number of rows for test %d: got: %d want: %d", i, rows, test.n)
		}
		if cols != test.n {
			t.Errorf("unexpected number of cols for test %d: got: %d want: %d", i, cols, test.n)
		}
		if !reflect.DeepEqual(sym, test.mat) {
			t.Errorf("unexpected data slice for test %d: got: %v want: %v", i, sym, test.mat)
		}

		m := NewDense(test.n, test.n, test.data)
		if !reflect.DeepEqual(sym.mat.Data, m.mat.Data) {
			t.Errorf("unexpected data slice mismatch for test %d: got: %v want: %v", i, sym.mat.Data, m.mat.Data)
		}
	}

	panicked, message := panics(func() { NewSymDense(3, []float64{1, 2}) })
	if !panicked || message != ErrShape.Error() {
		t.Error("expected panic for invalid data slice length")
	}
}

func TestSymAtSet(t *testing.T) {
	sym := &SymDense{blas64.Symmetric{
		N:      3,
		Stride: 3,
		Uplo:   blas.Upper,
		Data:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
	}}
	rows, cols := sym.Dims()

	// Check At out of bounds
	for _, row := range []int{-1, rows, rows + 1} {
		panicked, message := panics(func() { sym.At(row, 0) })
		if !panicked || message != ErrRowAccess.Error() {
			t.Errorf("expected panic for invalid row access N=%d r=%d", rows, row)
		}
	}
	for _, col := range []int{-1, cols, cols + 1} {
		panicked, message := panics(func() { sym.At(0, col) })
		if !panicked || message != ErrColAccess.Error() {
			t.Errorf("expected panic for invalid column access N=%d c=%d", cols, col)
		}
	}

	// Check Set out of bounds
	for _, row := range []int{-1, rows, rows + 1} {
		panicked, message := panics(func() { sym.SetSym(row, 0, 1.2) })
		if !panicked || message != ErrRowAccess.Error() {
			t.Errorf("expected panic for invalid row access N=%d r=%d", rows, row)
		}
	}
	for _, col := range []int{-1, cols, cols + 1} {
		panicked, message := panics(func() { sym.SetSym(0, col, 1.2) })
		if !panicked || message != ErrColAccess.Error() {
			t.Errorf("expected panic for invalid column access N=%d c=%d", cols, col)
		}
	}

	for _, st := range []struct {
		row, col  int
		orig, new float64
	}{
		{row: 1, col: 2, orig: 6, new: 15},
		{row: 2, col: 1, orig: 15, new: 12},
	} {
		if e := sym.At(st.row, st.col); e != st.orig {
			t.Errorf("unexpected value for At(%d, %d): got: %v want: %v", st.row, st.col, e, st.orig)
		}
		if e := sym.At(st.col, st.row); e != st.orig {
			t.Errorf("unexpected value for At(%d, %d): got: %v want: %v", st.col, st.row, e, st.orig)
		}
		sym.SetSym(st.row, st.col, st.new)
		if e := sym.At(st.row, st.col); e != st.new {
			t.Errorf("unexpected value for At(%d, %d) after SetSym(%[1]d, %[2]d, %[4]v): got: %[3]v want: %v", st.row, st.col, e, st.new)
		}
		if e := sym.At(st.col, st.row); e != st.new {
			t.Errorf("unexpected value for At(%d, %d) after SetSym(%[2]d, %[1]d, %[4]v): got: %[3]v want: %v", st.col, st.row, e, st.new)
		}
	}
}

func TestSymAdd(t *testing.T) {
	for _, test := range []struct {
		n int
	}{
		{n: 1},
		{n: 2},
		{n: 3},
		{n: 4},
		{n: 5},
		{n: 10},
	} {
		n := test.n
		a := NewSymDense(n, nil)
		for i := range a.mat.Data {
			a.mat.Data[i] = rand.Float64()
		}
		b := NewSymDense(n, nil)
		for i := range a.mat.Data {
			b.mat.Data[i] = rand.Float64()
		}
		var m Dense
		m.Add(a, b)

		// Check with new receiver
		var s SymDense
		s.AddSym(a, b)
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				want := m.At(i, j)
				if got := s.At(i, j); got != want {
					t.Errorf("unexpected value for At(%d, %d): got: %v want: %v", i, j, got, want)
				}
			}
		}

		// Check with equal receiver
		s.CopySym(a)
		s.AddSym(&s, b)
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				want := m.At(i, j)
				if got := s.At(i, j); got != want {
					t.Errorf("unexpected value for At(%d, %d): got: %v want: %v", i, j, got, want)
				}
			}
		}
	}

	method := func(receiver, a, b Matrix) {
		type addSymer interface {
			AddSym(a, b Symmetric)
		}
		rd := receiver.(addSymer)
		rd.AddSym(a.(Symmetric), b.(Symmetric))
	}
	denseComparison := func(receiver, a, b *Dense) {
		receiver.Add(a, b)
	}
	testTwoInput(t, "AddSym", &SymDense{}, method, denseComparison, legalTypesSym, legalSizeSameSquare, 1e-14)
}

func TestCopy(t *testing.T) {
	for _, test := range []struct {
		n int
	}{
		{n: 1},
		{n: 2},
		{n: 3},
		{n: 4},
		{n: 5},
		{n: 10},
	} {
		n := test.n
		a := NewSymDense(n, nil)
		for i := range a.mat.Data {
			a.mat.Data[i] = rand.Float64()
		}
		s := NewSymDense(n, nil)
		s.CopySym(a)
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				want := a.At(i, j)
				if got := s.At(i, j); got != want {
					t.Errorf("unexpected value for At(%d, %d): got: %v want: %v", i, j, got, want)
				}
			}
		}
	}
}

// TODO(kortschak) Roll this into testOneInput when it exists.
// https://github.com/gonum/matrix/issues/171
func TestSymCopyPanic(t *testing.T) {
	var (
		a SymDense
		n int
	)
	m := NewSymDense(1, nil)
	panicked, message := panics(func() { n = m.CopySym(&a) })
	if panicked {
		t.Errorf("unexpected panic: %v", message)
	}
	if n != 0 {
		t.Errorf("unexpected n: got: %d want: 0", n)
	}
}

func TestSymRankOne(t *testing.T) {
	for _, test := range []struct {
		n int
	}{
		{n: 1},
		{n: 2},
		{n: 3},
		{n: 4},
		{n: 5},
		{n: 10},
	} {
		n := test.n
		alpha := 2.0
		a := NewSymDense(n, nil)
		for i := range a.mat.Data {
			a.mat.Data[i] = rand.Float64()
		}
		x := make([]float64, n)
		for i := range x {
			x[i] = rand.Float64()
		}

		xMat := NewDense(n, 1, x)
		var m Dense
		m.Mul(xMat, xMat.T())
		m.Scale(alpha, &m)
		m.Add(&m, a)

		// Check with new receiver
		s := NewSymDense(n, nil)
		s.SymRankOne(a, alpha, NewVector(len(x), x))
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				want := m.At(i, j)
				if got := s.At(i, j); got != want {
					t.Errorf("unexpected value for At(%d, %d): got: %v want: %v", i, j, got, want)
				}
			}
		}

		// Check with reused receiver
		copy(s.mat.Data, a.mat.Data)
		s.SymRankOne(s, alpha, NewVector(len(x), x))
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				want := m.At(i, j)
				if got := s.At(i, j); got != want {
					t.Errorf("unexpected value for At(%d, %d): got: %v want: %v", i, j, got, want)
				}
			}
		}
	}
}

func TestRankTwo(t *testing.T) {
	for _, test := range []struct {
		n int
	}{
		{n: 1},
		{n: 2},
		{n: 3},
		{n: 4},
		{n: 5},
		{n: 10},
	} {
		n := test.n
		alpha := 2.0
		a := NewSymDense(n, nil)
		for i := range a.mat.Data {
			a.mat.Data[i] = rand.Float64()
		}
		x := make([]float64, n)
		y := make([]float64, n)
		for i := range x {
			x[i] = rand.Float64()
			y[i] = rand.Float64()
		}

		xMat := NewDense(n, 1, x)
		yMat := NewDense(n, 1, y)
		var m Dense
		m.Mul(xMat, yMat.T())
		var tmp Dense
		tmp.Mul(yMat, xMat.T())
		m.Add(&m, &tmp)
		m.Scale(alpha, &m)
		m.Add(&m, a)

		// Check with new receiver
		s := NewSymDense(n, nil)
		s.RankTwo(a, alpha, NewVector(len(x), x), NewVector(len(y), y))
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				if !floats.EqualWithinAbsOrRel(s.At(i, j), m.At(i, j), 1e-14, 1e-14) {
					t.Errorf("unexpected element value at (%d,%d): got: %f want: %f", i, j, m.At(i, j), s.At(i, j))
				}
			}
		}

		// Check with reused receiver
		copy(s.mat.Data, a.mat.Data)
		s.RankTwo(s, alpha, NewVector(len(x), x), NewVector(len(y), y))
		for i := 0; i < n; i++ {
			for j := i; j < n; j++ {
				if !floats.EqualWithinAbsOrRel(s.At(i, j), m.At(i, j), 1e-14, 1e-14) {
					t.Errorf("unexpected element value at (%d,%d): got: %f want: %f", i, j, m.At(i, j), s.At(i, j))
				}
			}
		}
	}
}

func TestSymRankK(t *testing.T) {
	alpha := 3.0
	method := func(receiver, a, b Matrix) {
		type SymRankKer interface {
			SymRankK(a Symmetric, alpha float64, x Matrix)
		}
		rd := receiver.(SymRankKer)
		rd.SymRankK(a.(Symmetric), alpha, b)
	}
	denseComparison := func(receiver, a, b *Dense) {
		var tmp Dense
		tmp.Mul(b, b.T())
		tmp.Scale(alpha, &tmp)
		receiver.Add(a, &tmp)
	}
	legalTypes := func(a, b Matrix) bool {
		_, ok := a.(Symmetric)
		return ok
	}
	legalSize := func(ar, ac, br, bc int) bool {
		if ar != ac {
			return false
		}
		return br == ar
	}
	testTwoInput(t, "SymRankK", &SymDense{}, method, denseComparison, legalTypes, legalSize, 1e-14)
}

func TestScaleSym(t *testing.T) {
	f := 3.0
	method := func(receiver, a Matrix) {
		type ScaleSymer interface {
			ScaleSym(f float64, a Symmetric)
		}
		rd := receiver.(ScaleSymer)
		rd.ScaleSym(f, a.(Symmetric))
	}
	denseComparison := func(receiver, a *Dense) {
		receiver.Scale(f, a)
	}
	testOneInput(t, "ScaleSym", &SymDense{}, method, denseComparison, legalTypeSym, isSquare, 1e-14)
}
