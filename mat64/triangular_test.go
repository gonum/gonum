package mat64

import (
	"math"
	"math/rand"
	"reflect"

	"github.com/gonum/blas"
	"github.com/gonum/blas/blas64"
	"gopkg.in/check.v1"
)

func (s *S) TestNewTriangular(c *check.C) {
	for i, test := range []struct {
		data  []float64
		n     int
		upper bool
		mat   *TriDense
	}{
		{
			data: []float64{
				1, 2, 3,
				4, 5, 6,
				7, 8, 9,
			},
			n:     3,
			upper: true,
			mat: &TriDense{blas64.Triangular{
				N:      3,
				Stride: 3,
				Uplo:   blas.Upper,
				Data:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
				Diag:   blas.NonUnit,
			}},
		},
	} {
		tri := NewTriDense(test.n, test.upper, test.data)
		rows, cols := tri.Dims()

		if rows != test.n {
			c.Errorf("unexpected number of rows for test %d: got: %d want: %d", i, rows, test.n)
		}
		if cols != test.n {
			c.Errorf("unexpected number of cols for test %d: got: %d want: %d", i, cols, test.n)
		}
		if !reflect.DeepEqual(tri, test.mat) {
			c.Errorf("unexpected data slice for test %d: got: %v want: %v", i, tri, test.mat)
		}
	}

	for _, upper := range []bool{false, true} {
		panicked, message := panics(func() { NewTriDense(3, upper, []float64{1, 2}) })
		if !panicked || message != ErrShape.Error() {
			c.Errorf("expected panic for invalid data slice length for upper=%t", upper)
		}
	}
}

func (s *S) TestTriAtSet(c *check.C) {
	tri := &TriDense{blas64.Triangular{
		N:      3,
		Stride: 3,
		Uplo:   blas.Upper,
		Data:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9},
		Diag:   blas.NonUnit,
	}}

	rows, cols := tri.Dims()

	// Check At out of bounds
	for _, row := range []int{-1, rows, rows + 1} {
		panicked, message := panics(func() { tri.At(row, 0) })
		if !panicked || message != ErrRowAccess.Error() {
			c.Errorf("expected panic for invalid row access N=%d r=%d", rows, row)
		}
	}
	for _, col := range []int{-1, cols, cols + 1} {
		panicked, message := panics(func() { tri.At(0, col) })
		if !panicked || message != ErrColAccess.Error() {
			c.Errorf("expected panic for invalid column access N=%d c=%d", cols, col)
		}
	}

	// Check Set out of bounds
	for _, row := range []int{-1, rows, rows + 1} {
		panicked, message := panics(func() { tri.SetTri(row, 0, 1.2) })
		if !panicked || message != ErrRowAccess.Error() {
			c.Errorf("expected panic for invalid row access N=%d r=%d", rows, row)
		}
	}
	for _, col := range []int{-1, cols, cols + 1} {
		panicked, message := panics(func() { tri.SetTri(0, col, 1.2) })
		if !panicked || message != ErrColAccess.Error() {
			c.Errorf("expected panic for invalid column access N=%d c=%d", cols, col)
		}
	}

	for _, st := range []struct {
		row, col int
		uplo     blas.Uplo
	}{
		{row: 2, col: 1, uplo: blas.Upper},
		{row: 1, col: 2, uplo: blas.Lower},
	} {
		tri.mat.Uplo = st.uplo
		panicked, message := panics(func() { tri.SetTri(st.row, st.col, 1.2) })
		if !panicked || message != ErrTriangleSet.Error() {
			c.Errorf("expected panic for %+v", st)
		}
	}

	for _, st := range []struct {
		row, col  int
		uplo      blas.Uplo
		orig, new float64
	}{
		{row: 2, col: 1, uplo: blas.Lower, orig: 8, new: 15},
		{row: 1, col: 2, uplo: blas.Upper, orig: 6, new: 15},
	} {
		tri.mat.Uplo = st.uplo
		if e := tri.At(st.row, st.col); e != st.orig {
			c.Errorf("unexpected value for At(%d, %d): got: %v want: %v", st.row, st.col, e, st.orig)
		}
		tri.SetTri(st.row, st.col, st.new)
		if e := tri.At(st.row, st.col); e != st.new {
			c.Errorf("unexpected value for At(%d, %d) after SetTri(%[1]d, %d, %v): got: %v want: %[3]v", st.row, st.col, st.new, e)
		}
	}
}

func (s *S) TestTriDenseCopy(c *check.C) {
	for i := 0; i < 100; i++ {
		size := rand.Intn(100)
		r, err := randDense(size, 0.9, rand.NormFloat64)
		if size == 0 {
			if err != ErrZeroLength {
				c.Fatalf("expected error %v: got: %v", ErrZeroLength, err)
			}
			continue
		}
		if err != nil {
			c.Fatalf("unexpected error: %v", err)
		}

		u := NewTriDense(size, true, nil)
		l := NewTriDense(size, false, nil)

		for _, typ := range []Matrix{r, (*basicVectorer)(r), (*basicMatrix)(r)} {
			for j := range u.mat.Data {
				u.mat.Data[j] = math.NaN()
				l.mat.Data[j] = math.NaN()
			}
			u.Copy(typ)
			l.Copy(typ)
			for m := 0; m < size; m++ {
				for n := 0; n < size; n++ {
					want := typ.At(m, n)
					switch {
					case m < n: // Upper triangular matrix.
						if got := u.At(m, n); got != want {
							c.Errorf("unexpected upper value for At(%d, %d) for test %d: got: %v want: %v", m, n, i, got, want)
						}
					case m == n: // Diagonal matrix.
						if got := u.At(m, n); got != want {
							c.Errorf("unexpected upper value for At(%d, %d) for test %d: got: %v want: %v", m, n, i, got, want)
						}
						if got := l.At(m, n); got != want {
							c.Errorf("unexpected diagonal value for At(%d, %d) for test %d: got: %v want: %v", m, n, i, got, want)
						}
					case m < n: // Lower triangular matrix.
						if got := l.At(m, n); got != want {
							c.Errorf("unexpected lower value for At(%d, %d) for test %d: got: %v want: %v", m, n, i, got, want)
						}
					}
				}
			}
		}
	}
}

func (s *S) TestTriTriDenseCopy(c *check.C) {
	for i := 0; i < 100; i++ {
		size := rand.Intn(100)
		r, err := randDense(size, 1, rand.NormFloat64)
		if size == 0 {
			if err != ErrZeroLength {
				c.Fatalf("expected error %v: got: %v", ErrZeroLength, err)
			}
			continue
		}
		if err != nil {
			c.Fatalf("unexpected error: %v", err)
		}

		ur := NewTriDense(size, true, nil)
		lr := NewTriDense(size, false, nil)

		ur.Copy(r)
		lr.Copy(r)

		u := NewTriDense(size, true, nil)
		u.Copy(ur)
		if !equal(u, ur) {
			c.Fatal("unexpected result for U triangle copy of U triangle: not equal")
		}

		l := NewTriDense(size, false, nil)
		l.Copy(lr)
		if !equal(l, lr) {
			c.Fatal("unexpected result for L triangle copy of L triangle: not equal")
		}

		zero(u.mat.Data)
		u.Copy(lr)
		if !isDiagonal(u) {
			c.Fatal("unexpected result for U triangle copy of L triangle: off diagonal non-zero element")
		}
		if !equalDiagonal(u, lr) {
			c.Fatal("unexpected result for U triangle copy of L triangle: diagonal not equal")
		}

		zero(l.mat.Data)
		l.Copy(ur)
		if !isDiagonal(l) {
			c.Fatal("unexpected result for L triangle copy of U triangle: off diagonal non-zero element")
		}
		if !equalDiagonal(l, ur) {
			c.Fatal("unexpected result for L triangle copy of U triangle: diagonal not equal")
		}
	}
}
