// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"math/cmplx"
	"testing"

	"golang.org/x/exp/rand"
)

func TestCDenseNewAtSet(t *testing.T) {
	t.Parallel()
	for cas, test := range []struct {
		a          []complex128
		rows, cols int
	}{
		{
			a: []complex128{0, 0, 0,
				0, 0, 0,
				0, 0, 0},
			rows: 3,
			cols: 3,
		},
	} {
		aCopy := make([]complex128, len(test.a))
		copy(aCopy, test.a)
		mZero := NewCDense(test.rows, test.cols, nil)
		rows, cols := mZero.Dims()
		if rows != test.rows {
			t.Errorf("unexpected number of rows for test %d: got: %d want: %d", cas, rows, test.rows)
		}
		if cols != test.cols {
			t.Errorf("unexpected number of cols for test %d: got: %d want: %d", cas, cols, test.cols)
		}
		m := NewCDense(test.rows, test.cols, aCopy)
		for i := 0; i < test.rows; i++ {
			for j := 0; j < test.cols; j++ {
				v := m.At(i, j)
				idx := i*test.rows + j
				if v != test.a[idx] {
					t.Errorf("unexpected get value for test %d at i=%d, j=%d: got: %v, want: %v", cas, i, j, v, test.a[idx])
				}
				add := complex(float64(i+1), float64(j+1))
				m.Set(i, j, v+add)
				if m.At(i, j) != test.a[idx]+add {
					t.Errorf("unexpected set value for test %d at i=%d, j=%d: got: %v, want: %v", cas, i, j, v, test.a[idx]+add)
				}
			}
		}
	}
}

func TestCDenseConjElem(t *testing.T) {
	t.Parallel()

	rnd := rand.New(rand.NewSource(1))

	for r := 1; r <= 8; r++ {
		for c := 1; c <= 8; c++ {
			const (
				empty = iota
				fit
				sliced
				self
			)
			for _, dst := range []int{empty, fit, sliced, self} {
				const (
					noTrans = iota
					trans
					conjTrans
					bothHT
					bothTH
				)
				for _, src := range []int{noTrans, trans, conjTrans, bothHT, bothTH} {
					d := NewCDense(r, c, nil)
					for i := 0; i < r; i++ {
						for j := 0; j < c; j++ {
							d.Set(i, j, complex(rnd.NormFloat64(), rnd.NormFloat64()))
						}
					}

					var (
						a  CMatrix
						op string
					)
					switch src {
					case noTrans:
						a = d
					case trans:
						r, c = c, r
						a = d.T()
						op = ".T"
					case conjTrans:
						r, c = c, r
						a = d.H()
						op = ".H"
					case bothHT:
						a = d.H().T()
						op = ".H.T"
					case bothTH:
						a = d.T().H()
						op = ".T.H"
					default:
						panic("invalid src op")
					}
					aCopy := NewCDense(r, c, nil)
					aCopy.Copy(a)

					var got *CDense
					switch dst {
					case empty:
						got = &CDense{}
					case fit:
						got = NewCDense(r, c, nil)
					case sliced:
						got = NewCDense(r*2, c*2, nil).Slice(1, r+1, 1, c+1).(*CDense)
					case self:
						if r != c && (src == conjTrans || src == trans) {
							continue
						}
						got = d
					default:
						panic("invalid dst size")
					}

					got.Conj(a)

					for i := 0; i < r; i++ {
						for j := 0; j < c; j++ {
							if got.At(i, j) != cmplx.Conj(aCopy.At(i, j)) {
								t.Errorf("unexpected results a%s[%d, %d] for r=%d c=%d %v != %v",
									op, i, j, r, c, got.At(i, j), cmplx.Conj(a.At(i, j)),
								)
							}
						}
					}
				}
			}
		}
	}
}

func TestCDenseGrow(t *testing.T) {
	t.Parallel()
	m := &CDense{}
	m = m.Grow(10, 10).(*CDense)
	rows, cols := m.Dims()
	capRows, capCols := m.Caps()
	if rows != 10 {
		t.Errorf("unexpected value for rows: got: %d want: 10", rows)
	}
	if cols != 10 {
		t.Errorf("unexpected value for cols: got: %d want: 10", cols)
	}
	if capRows != 10 {
		t.Errorf("unexpected value for capRows: got: %d want: 10", capRows)
	}
	if capCols != 10 {
		t.Errorf("unexpected value for capCols: got: %d want: 10", capCols)
	}

	// Test grow within caps is in-place.
	m.Set(1, 1, 1)
	v := m.Slice(1, 5, 1, 5).(*CDense)
	if v.At(0, 0) != m.At(1, 1) {
		t.Errorf("unexpected viewed element value: got: %v want: %v", v.At(0, 0), m.At(1, 1))
	}
	v = v.Grow(5, 5).(*CDense)
	if !CEqual(v, m.Slice(1, 10, 1, 10)) {
		t.Error("unexpected view value after grow")
	}

	// Test grow bigger than caps copies.
	v = v.Grow(5, 5).(*CDense)
	if !CEqual(v.Slice(0, 9, 0, 9), m.Slice(1, 10, 1, 10)) {
		t.Error("unexpected mismatched common view value after grow")
	}
	v.Set(0, 0, 0)
	if CEqual(v.Slice(0, 9, 0, 9), m.Slice(1, 10, 1, 10)) {
		t.Error("unexpected matching view value after grow past capacity")
	}

	// Test grow uses existing data slice when matrix is zero size.
	v.Reset()
	p, l := &v.mat.Data[:1][0], cap(v.mat.Data)
	*p = 1 // This element is at position (-1, -1) relative to v and so should not be visible.
	v = v.Grow(5, 5).(*CDense)
	if &v.mat.Data[:1][0] != p {
		t.Error("grow unexpectedly copied slice within cap limit")
	}
	if cap(v.mat.Data) != l {
		t.Errorf("unexpected change in data slice capacity: got: %d want: %d", cap(v.mat.Data), l)
	}
	if v.At(0, 0) != 0 {
		t.Errorf("unexpected value for At(0, 0): got: %v want: 0", v.At(0, 0))
	}
}

func TestCDenseAdd(t *testing.T) {
	type args struct {
		a CMatrix
		b CMatrix
	}
	tests := []struct {
		name string
		args args
		want *CDense
	}{
		{
			name: "1 + 1 ?= 2",
			args: args{
				a: NewCDense(3, 3, []complex128{
					1.0 + 0.0i, 0.0 + 0.0i, 0.0 + 0.0i,
					0.0 + 0.0i, 1.0 + 0.0i, 0.0 + 0.0i,
					0.0 + 0.0i, 0.0 + 0.0i, 1.0 + 0.0i,
				}),
				b: NewCDense(3, 3, []complex128{
					1.0 + 0.0i, 0.0 + 0.0i, 0.0 + 0.0i,
					0.0 + 0.0i, 1.0 + 0.0i, 0.0 + 0.0i,
					0.0 + 0.0i, 0.0 + 0.0i, 1.0 + 0.0i,
				}),
			},
			want: NewCDense(3, 3, []complex128{
				2.0 + 0.0i, 0.0 + 0.0i, 0.0 + 0.0i,
				0.0 + 0.0i, 2.0 + 0.0i, 0.0 + 0.0i,
				0.0 + 0.0i, 0.0 + 0.0i, 2.0 + 0.0i,
			}),
		},
		{
			name: "Random 5x5 matrices sum",
			args: args{
				a: NewCDense(5, 5, []complex128{
					1.34 + 0.56i, 0.78 + 0.90i, 0.12 + 0.34i, 0.56 + 0.78i, 0.90 + 0.12i,
					0.56 + 0.78i, 0.90 + 0.12i, 0.34 + 0.56i, 0.78 + 0.90i, 0.12 + 0.34i,
					0.78 + 0.90i, 0.12 + 0.34i, 0.56 + 0.78i, 0.90 + 0.12i, 0.34 + 0.56i,
					0.90 + 0.12i, 0.34 + 0.56i, 0.78 + 0.90i, 0.12 + 0.34i, 0.56 + 0.78i,
					0.12 + 0.34i, 0.56 + 0.78i, 0.90 + 0.12i, 0.34 + 0.56i, 0.78 + 0.90i,
				}),
				b: NewCDense(5, 5, []complex128{
					1.12 + 0.23i, 0.45 + 0.67i, 0.78 + 0.89i, 0.12 + 0.34i, 0.56 + 0.78i,
					0.34 + 0.56i, 0.78 + 0.90i, 0.12 + 0.34i, 0.56 + 0.78i, 0.90 + 0.12i,
					0.56 + 0.78i, 0.90 + 0.12i, 0.34 + 0.56i, 0.78 + 0.90i, 0.12 + 0.34i,
					0.78 + 0.90i, 0.12 + 0.34i, 0.56 + 0.78i, 0.90 + 0.12i, 0.34 + 0.56i,
					0.90 + 0.12i, 0.34 + 0.56i, 0.78 + 0.90i, 0.12 + 0.34i, 0.56 + 0.78i,
				}),
			},
			want: NewCDense(5, 5, []complex128{
				2.46 + 0.79i, 1.23 + 1.57i, 0.9 + 1.23i, 0.68 + 1.12i, 1.46 + 0.9i,
				0.9 + 1.34i, 1.68 + 1.02i, 0.46 + 0.9i, 1.34 + 1.68i, 1.02 + 0.46i,
				1.34 + 1.68i, 1.02 + 0.46i, 0.9 + 1.34i, 1.68 + 1.02i, 0.46 + 0.9i,
				1.68 + 1.02i, 0.46 + 0.9i, 1.34 + 1.68i, 1.02 + 0.46i, 0.9 + 1.34i,
				1.02 + 0.46i, 0.9 + 1.34i, 1.68 + 1.02i, 0.46 + 0.9i, 1.34 + 1.68i,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var receiver CDense
			receiver.Add(tt.args.a, tt.args.b)
			if !CEqualApprox(&receiver, tt.want, 1e-14) {
				t.Errorf("CDense Addition fail. Wanted: %v, got: %v", tt.want, receiver)
			}
		})
	}
}

func TestCDenseSub(t *testing.T) {
	type args struct {
		a CMatrix
		b CMatrix
	}
	tests := []struct {
		name string
		args args
		want *CDense
	}{
		{
			name: "1 - 1 ?= 0",
			args: args{
				a: NewCDense(3, 3, []complex128{
					1.0 + 0.0i, 0.0 + 0.0i, 0.0 + 0.0i,
					0.0 + 0.0i, 1.0 + 0.0i, 0.0 + 0.0i,
					0.0 + 0.0i, 0.0 + 0.0i, 1.0 + 0.0i,
				}),
				b: NewCDense(3, 3, []complex128{
					1.0 + 0.0i, 0.0 + 0.0i, 0.0 + 0.0i,
					0.0 + 0.0i, 1.0 + 0.0i, 0.0 + 0.0i,
					0.0 + 0.0i, 0.0 + 0.0i, 1.0 + 0.0i,
				}),
			},
			want: NewCDense(3, 3, []complex128{
				0.0 + 0.0i, 0.0 + 0.0i, 0.0 + 0.0i,
				0.0 + 0.0i, 0.0 + 0.0i, 0.0 + 0.0i,
				0.0 + 0.0i, 0.0 + 0.0i, 0.0 + 0.0i,
			}),
		},
		{
			name: "Random 5x5 matrices sum",
			args: args{
				a: NewCDense(5, 5, []complex128{
					1.34 + 0.56i, 0.78 + 0.90i, 0.12 + 0.34i, 0.56 + 0.78i, 0.90 + 0.12i,
					0.56 + 0.78i, 0.90 + 0.12i, 0.34 + 0.56i, 0.78 + 0.90i, 0.12 + 0.34i,
					0.78 + 0.90i, 0.12 + 0.34i, 0.56 + 0.78i, 0.90 + 0.12i, 0.34 + 0.56i,
					0.90 + 0.12i, 0.34 + 0.56i, 0.78 + 0.90i, 0.12 + 0.34i, 0.56 + 0.78i,
					0.12 + 0.34i, 0.56 + 0.78i, 0.90 + 0.12i, 0.34 + 0.56i, 0.78 + 0.90i,
				}),
				b: NewCDense(5, 5, []complex128{
					-1.12 - 0.23i, -0.45 - 0.67i, -0.78 - 0.89i, -0.12 - 0.34i, -0.56 - 0.78i,
					-0.34 - 0.56i, -0.78 - 0.90i, -0.12 - 0.34i, -0.56 - 0.78i, -0.90 - 0.12i,
					-0.56 - 0.78i, -0.90 - 0.12i, -0.34 - 0.56i, -0.78 - 0.90i, -0.12 - 0.34i,
					-0.78 - 0.90i, -0.12 - 0.34i, -0.56 - 0.78i, -0.90 - 0.12i, -0.34 - 0.56i,
					-0.90 - 0.12i, -0.34 - 0.56i, -0.78 - 0.90i, -0.12 - 0.34i, -0.56 - 0.78i,
				}),
			},
			want: NewCDense(5, 5, []complex128{
				2.46 + 0.79i, 1.23 + 1.57i, 0.9 + 1.23i, 0.68 + 1.12i, 1.46 + 0.9i,
				0.9 + 1.34i, 1.68 + 1.02i, 0.46 + 0.9i, 1.34 + 1.68i, 1.02 + 0.46i,
				1.34 + 1.68i, 1.02 + 0.46i, 0.9 + 1.34i, 1.68 + 1.02i, 0.46 + 0.9i,
				1.68 + 1.02i, 0.46 + 0.9i, 1.34 + 1.68i, 1.02 + 0.46i, 0.9 + 1.34i,
				1.02 + 0.46i, 0.9 + 1.34i, 1.68 + 1.02i, 0.46 + 0.9i, 1.34 + 1.68i,
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var receiver CDense
			receiver.Sub(tt.args.a, tt.args.b)
			if !CEqualApprox(&receiver, tt.want, 1e-14) {
				t.Errorf("CDense Subtraction fail. Wanted: %v, got: %v", tt.want, receiver)
			}
		})
	}
}
