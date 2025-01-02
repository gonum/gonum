// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"fmt"
	"math/rand/v2"
	"reflect"
	"testing"

	"gonum.org/v1/gonum/lapack/lapack64"
)

func TestNewTridiag(t *testing.T) {
	for i, test := range []struct {
		n         int
		dl, d, du []float64
		panics    bool
		want      *Tridiag
		dense     *Dense
	}{
		{
			n:      1,
			dl:     nil,
			d:      []float64{1.2},
			du:     nil,
			panics: false,
			want: &Tridiag{
				mat: lapack64.Tridiagonal{
					N:  1,
					DL: nil,
					D:  []float64{1.2},
					DU: nil,
				},
			},
			dense: NewDense(1, 1, []float64{1.2}),
		},
		{
			n:      1,
			dl:     []float64{},
			d:      []float64{1.2},
			du:     []float64{},
			panics: false,
			want: &Tridiag{
				mat: lapack64.Tridiagonal{
					N:  1,
					DL: []float64{},
					D:  []float64{1.2},
					DU: []float64{},
				},
			},
			dense: NewDense(1, 1, []float64{1.2}),
		},
		{
			n:      4,
			dl:     []float64{1.2, 2.3, 3.4},
			d:      []float64{4.5, 5.6, 6.7, 7.8},
			du:     []float64{8.9, 9.0, 0.1},
			panics: false,
			want: &Tridiag{
				mat: lapack64.Tridiagonal{
					N:  4,
					DL: []float64{1.2, 2.3, 3.4},
					D:  []float64{4.5, 5.6, 6.7, 7.8},
					DU: []float64{8.9, 9.0, 0.1},
				},
			},
			dense: NewDense(4, 4, []float64{
				4.5, 8.9, 0, 0,
				1.2, 5.6, 9.0, 0,
				0, 2.3, 6.7, 0.1,
				0, 0, 3.4, 7.8,
			}),
		},
		{
			n:      4,
			dl:     nil,
			d:      nil,
			du:     nil,
			panics: false,
			want: &Tridiag{
				mat: lapack64.Tridiagonal{
					N:  4,
					DL: []float64{0, 0, 0},
					D:  []float64{0, 0, 0, 0},
					DU: []float64{0, 0, 0},
				},
			},
			dense: NewDense(4, 4, nil),
		},
		{
			n:      -1,
			panics: true,
		},
		{
			n:      0,
			panics: true,
		},
		{
			n:      1,
			dl:     []float64{1.2},
			d:      nil,
			du:     nil,
			panics: true,
		},
		{
			n:      1,
			dl:     nil,
			d:      []float64{1.2, 2.3},
			du:     nil,
			panics: true,
		},
		{
			n:      1,
			dl:     []float64{},
			d:      nil,
			du:     []float64{},
			panics: true,
		},
		{
			n:      4,
			dl:     []float64{1.2},
			d:      nil,
			du:     nil,
			panics: true,
		},
		{
			n:      4,
			dl:     []float64{1.2, 2.3, 3.4},
			d:      []float64{4.5, 5.6, 6.7, 7.8, 1.2},
			du:     []float64{8.9, 9.0, 0.1},
			panics: true,
		},
	} {
		var a *Tridiag
		panicked, msg := panics(func() {
			a = NewTridiag(test.n, test.dl, test.d, test.du)
		})
		if panicked {
			if !test.panics {
				t.Errorf("Case %d: unexpected panic: %s", i, msg)
			}
			continue
		}
		if test.panics {
			t.Errorf("Case %d: expected panic", i)
			continue
		}

		r, c := a.Dims()
		if r != test.n {
			t.Errorf("Case %d: unexpected number of rows: got=%d want=%d", i, r, test.n)
		}
		if c != test.n {
			t.Errorf("Case %d: unexpected number of columns: got=%d want=%d", i, c, test.n)
		}

		kl, ku := a.Bandwidth()
		if kl != 1 || ku != 1 {
			t.Errorf("Case %d: unexpected bandwidth: got=%d,%d want=1,1", i, kl, ku)
		}

		if !reflect.DeepEqual(a, test.want) {
			t.Errorf("Case %d: unexpected value via reflect: got=%v, want=%v", i, a, test.want)
		}
		if !Equal(a, test.want) {
			t.Errorf("Case %d: unexpected value via mat.Equal: got=%v, want=%v", i, a, test.want)
		}
		if !Equal(a, test.dense) {
			t.Errorf("Case %d: unexpected value via mat.Equal(Tridiag,Dense):\ngot:\n% v\nwant:\n% v", i, Formatted(a), Formatted(test.dense))
		}
	}
}

func TestTridiagAtSet(t *testing.T) {
	t.Parallel()
	for _, n := range []int{1, 2, 3, 4, 7, 10} {
		tri, ref := newTestTridiag(n)

		name := fmt.Sprintf("Case n=%v", n)

		// Check At explicitly with all valid indices.
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if tri.At(i, j) != ref.At(i, j) {
					t.Errorf("%v: unexpected value for At(%d,%d): got %v, want %v",
						name, i, j, tri.At(i, j), ref.At(i, j))
				}
			}
		}
		// Check At via a call to Equal.
		if !Equal(tri, ref) {
			t.Errorf("%v: unexpected value:\ngot: % v\nwant:% v",
				name, Formatted(tri, Prefix("     ")), Formatted(ref, Prefix("     ")))
		}

		// Check At out of bounds.
		for _, i := range []int{-1, n, n + 1} {
			for j := 0; j < n; j++ {
				panicked, message := panics(func() { tri.At(i, j) })
				if !panicked || message != ErrRowAccess.Error() {
					t.Errorf("%v: expected panic for invalid row access at (%d,%d)", name, i, j)
				}
			}
		}
		for _, j := range []int{-1, n, n + 1} {
			for i := 0; i < n; i++ {
				panicked, message := panics(func() { tri.At(i, j) })
				if !panicked || message != ErrColAccess.Error() {
					t.Errorf("%v: expected panic for invalid column access at (%d,%d)", name, i, j)
				}
			}
		}

		// Check SetBand out of bounds.
		for _, i := range []int{-1, n, n + 1} {
			for j := 0; j < n; j++ {
				panicked, message := panics(func() { tri.SetBand(i, j, 1.2) })
				if !panicked || message != ErrRowAccess.Error() {
					t.Errorf("%v: expected panic for invalid row access at (%d,%d)", name, i, j)
				}
			}
		}
		for _, j := range []int{-1, n, n + 1} {
			for i := 0; i < n; i++ {
				panicked, message := panics(func() { tri.SetBand(i, j, 1.2) })
				if !panicked || message != ErrColAccess.Error() {
					t.Errorf("%v: expected panic for invalid column access at (%d,%d)", name, i, j)
				}
			}
		}
		for i := 0; i < n; i++ {
			for j := 0; j <= i-2; j++ {
				panicked, message := panics(func() { tri.SetBand(i, j, 1.2) })
				if !panicked || message != ErrBandSet.Error() {
					t.Errorf("%v: expected panic for invalid access at (%d,%d)", name, i, j)
				}
			}
			for j := i + 2; j < n; j++ {
				panicked, message := panics(func() { tri.SetBand(i, j, 1.2) })
				if !panicked || message != ErrBandSet.Error() {
					t.Errorf("%v: expected panic for invalid access at (%d,%d)", name, i, j)
				}
			}
		}

		// Check SetBand within bandwidth.
		for i := 0; i < n; i++ {
			for j := max(0, i-1); j <= min(i+1, n-1); j++ {
				want := float64(i*n + j + 100)
				tri.SetBand(i, j, want)
				if got := tri.At(i, j); got != want {
					t.Errorf("%v: unexpected value at (%d,%d) after SetBand: got %v, want %v", name, i, j, got, want)
				}
			}
		}
	}
}

func newTestTridiag(n int) (*Tridiag, *Dense) {
	var dl, d, du []float64
	d = make([]float64, n)
	if n > 1 {
		dl = make([]float64, n-1)
		du = make([]float64, n-1)
	}
	for i := range d {
		d[i] = float64(i*n + i + 1)
	}
	for j := range dl {
		i := j + 1
		dl[j] = float64(i*n + j + 1)
	}
	for i := range du {
		j := i + 1
		du[i] = float64(i*n + j + 1)
	}
	dense := make([]float64, n*n)
	for i := 0; i < n; i++ {
		for j := max(0, i-1); j <= min(i+1, n-1); j++ {
			dense[i*n+j] = float64(i*n + j + 1)
		}
	}
	return NewTridiag(n, dl, d, du), NewDense(n, n, dense)
}

func TestTridiagReset(t *testing.T) {
	t.Parallel()
	for _, n := range []int{1, 2, 3, 4, 7, 10} {
		a, _ := newTestTridiag(n)
		if a.IsEmpty() {
			t.Errorf("Case n=%d: matrix is empty", n)
		}
		a.Reset()
		if !a.IsEmpty() {
			t.Errorf("Case n=%d: matrix is not empty after Reset", n)
		}
	}
}

func TestTridiagDiagView(t *testing.T) {
	t.Parallel()
	for _, n := range []int{1, 2, 3, 4, 7, 10} {
		a, _ := newTestTridiag(n)
		testDiagView(t, n, a)
	}
}

func TestTridiagZero(t *testing.T) {
	t.Parallel()
	for _, n := range []int{1, 2, 3, 4, 7, 10} {
		a, _ := newTestTridiag(n)
		a.Zero()
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				if a.At(i, j) != 0 {
					t.Errorf("Case n=%d: unexpected non-zero at (%d,%d): got %f", n, i, j, a.At(i, j))
				}
			}
		}
	}
}

func TestTridiagSolveTo(t *testing.T) {
	t.Parallel()

	const tol = 1e-13

	rnd := rand.New(rand.NewPCG(1, 1))
	random := func(n int) []float64 {
		d := make([]float64, n)
		for i := range d {
			d[i] = rnd.NormFloat64()
		}
		return d
	}

	for _, n := range []int{1, 2, 3, 4, 7, 10} {
		a := NewTridiag(n, random(n-1), random(n), random(n-1))
		var aDense Dense
		aDense.CloneFrom(a)
		for _, trans := range []bool{false, true} {
			for _, nrhs := range []int{1, 2, 5} {
				const (
					denseB = iota
					rawB
					basicB
				)
				for _, bType := range []int{denseB, rawB, basicB} {
					const (
						emptyDst = iota
						shapedDst
						bIsDst
					)
					for _, dstType := range []int{emptyDst, shapedDst, bIsDst} {
						if dstType == bIsDst && bType != denseB {
							continue
						}

						var b Matrix
						switch bType {
						case denseB:
							b = NewDense(n, nrhs, random(n*nrhs))
						case rawB:
							b = &rawMatrix{asBasicMatrix(NewDense(n, nrhs, random(n*nrhs)))}
						case basicB:
							b = asBasicMatrix(NewDense(n, nrhs, random(n*nrhs)))
						default:
							panic("bad bType")
						}

						var dst *Dense
						switch dstType {
						case emptyDst:
							dst = new(Dense)
						case shapedDst:
							dst = NewDense(n, nrhs, random(n*nrhs))
						case bIsDst:
							dst = b.(*Dense)
						default:
							panic("bad dstType")
						}

						name := fmt.Sprintf("n=%d,nrhs=%d,trans=%t,dstType=%d,bType=%d", n, nrhs, trans, dstType, bType)

						var want Dense
						var err error
						if !trans {
							err = want.Solve(&aDense, b)
						} else {
							err = want.Solve(aDense.T(), b)
						}
						if err != nil {
							t.Fatalf("%v: unexpected failure when computing reference solution: %v", name, err)
						}

						err = a.SolveTo(dst, trans, b)
						if err != nil {
							t.Fatalf("%v: unexpected failure from Tridiag.SolveTo: %v", name, err)
						}

						var diff Dense
						diff.Sub(dst, &want)
						if resid := Norm(&diff, 1); resid > tol*float64(n) {
							t.Errorf("%v: unexpected result; resid=%v, want<=%v", name, resid, tol*float64(n))
						}
					}
				}
			}
		}
	}
}

func TestTridiagSolveVecTo(t *testing.T) {
	t.Parallel()

	const tol = 1e-13

	rnd := rand.New(rand.NewPCG(1, 1))
	random := func(n int) []float64 {
		d := make([]float64, n)
		for i := range d {
			d[i] = rnd.NormFloat64()
		}
		return d
	}

	for _, n := range []int{1, 2, 3, 4, 7, 10} {
		a := NewTridiag(n, random(n-1), random(n), random(n-1))
		var aDense Dense
		aDense.CloneFrom(a)
		for _, trans := range []bool{false, true} {
			const (
				denseB = iota
				rawB
				basicB
			)
			for _, bType := range []int{denseB, rawB, basicB} {
				const (
					emptyDst = iota
					shapedDst
					bIsDst
				)
				for _, dstType := range []int{emptyDst, shapedDst, bIsDst} {
					if dstType == bIsDst && bType != denseB {
						continue
					}

					var b Vector
					switch bType {
					case denseB:
						b = NewVecDense(n, random(n))
					case rawB:
						b = &rawVector{asBasicVector(NewVecDense(n, random(n)))}
					case basicB:
						b = asBasicVector(NewVecDense(n, random(n)))
					default:
						panic("bad bType")
					}

					var dst *VecDense
					switch dstType {
					case emptyDst:
						dst = new(VecDense)
					case shapedDst:
						dst = NewVecDense(n, random(n))
					case bIsDst:
						dst = b.(*VecDense)
					default:
						panic("bad dstType")
					}

					name := fmt.Sprintf("n=%d,trans=%t,dstType=%d,bType=%d", n, trans, dstType, bType)

					var want VecDense
					var err error
					if !trans {
						err = want.SolveVec(&aDense, b)
					} else {
						err = want.SolveVec(aDense.T(), b)
					}
					if err != nil {
						t.Fatalf("%v: unexpected failure when computing reference solution: %v", name, err)
					}

					err = a.SolveVecTo(dst, trans, b)
					if err != nil {
						t.Fatalf("%v: unexpected failure from Tridiag.SolveTo: %v", name, err)
					}

					var diff Dense
					diff.Sub(dst, &want)
					if resid := Norm(&diff, 1); resid > tol*float64(n) {
						t.Errorf("%v: unexpected result; resid=%v, want<=%v", name, resid, tol*float64(n))
					}
				}
			}
		}
	}
}
