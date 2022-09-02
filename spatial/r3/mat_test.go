// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3

import (
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/num/quat"
)

func TestMatAdd(t *testing.T) {
	const tol = 1e-16
	rnd := rand.New(rand.NewSource(1))
	for tc := 0; tc < 20; tc++ {
		a := randomMat(rnd)
		b := randomMat(rnd)
		var (
			want mat.Dense
			got  Mat
		)
		want.Add(a, b)
		got.Add(a, b)
		if !mat.EqualApprox(&got, &want, tol) {
			t.Errorf("unexpected result for matrix add:\ngot:\n%v\nwant:\n%v", mat.Formatted(&got), mat.Formatted(&want))
		}
	}
}

func TestMatSub(t *testing.T) {
	const tol = 1e-16
	rnd := rand.New(rand.NewSource(1))
	for tc := 0; tc < 20; tc++ {
		a := randomMat(rnd)
		b := randomMat(rnd)
		var (
			want mat.Dense
			got  Mat
		)
		want.Sub(a, b)
		got.Sub(a, b)
		if !mat.EqualApprox(&got, &want, tol) {
			t.Errorf("unexpected result for matrix subtract:\ngot:\n%v\nwant:\n%v", mat.Formatted(&got), mat.Formatted(&want))
		}
	}
}

func TestMatMul(t *testing.T) {
	const tol = 1e-14
	rnd := rand.New(rand.NewSource(1))
	for tc := 0; tc < 20; tc++ {
		a := randomMat(rnd)
		b := randomMat(rnd)
		var (
			want mat.Dense
			got  Mat
		)
		want.Mul(a, b)
		got.Mul(a, b)
		if !mat.EqualApprox(&got, &want, tol) {
			t.Errorf("unexpected result for matrix multiply:\ngot:\n%v\nwant:\n%v", mat.Formatted(&got), mat.Formatted(&want))
		}
	}
}

func TestMatScale(t *testing.T) {
	const tol = 1e-16
	rnd := rand.New(rand.NewSource(1))
	for tc := 0; tc < 20; tc++ {
		v := rnd.Float64()
		a := randomMat(rnd)
		var (
			want mat.Dense
			got  Mat
		)
		want.Scale(v, a)
		got.Scale(v, a)
		if !mat.EqualApprox(&got, &want, tol) {
			t.Errorf("unexpected result for matrix scale:\ngot:\n%v\nwant:\n%v", mat.Formatted(&got), mat.Formatted(&want))
		}
	}
}

func TestMatCloneFrom(t *testing.T) {
	const tol = 1e-16
	rnd := rand.New(rand.NewSource(1))
	for tc := 0; tc < 20; tc++ {
		want := randomMat(rnd)
		got := NewMat(nil)
		got.CloneFrom(want)
		if !mat.EqualApprox(got, want, tol) {
			t.Errorf("unexpected result from CloneFrom:\ngot:\n%v\nwant:\n%v", mat.Formatted(got), mat.Formatted(want))
		}
	}
}

func TestSkew(t *testing.T) {
	const tol = 1e-16
	rnd := rand.New(rand.NewSource(1))
	for tc := 0; tc < 20; tc++ {
		sk := NewMat(nil)
		v1 := randomVec(rnd)
		v2 := randomVec(rnd)
		sk.Skew(v1)
		want := Cross(v1, v2)
		got := sk.MulVec(v2)
		if d := Sub(want, got); Dot(d, d) > tol {
			t.Errorf("r3.Cross(v1,v2) does not agree with r3.Skew(v1)*v2: got:%v want:%v", got, want)
		}
	}
}

func TestTranspose(t *testing.T) {
	const tol = 1e-16
	rnd := rand.New(rand.NewSource(1))
	for tc := 0; tc < 20; tc++ {
		d := mat.NewDense(3, 3, nil)
		m := randomMat(rnd)
		d.CloneFrom(m)
		mt := m.T()
		dt := d.T()
		if !mat.Equal(mt, dt) {
			t.Errorf("Dense.T() not equal to r3.Mat.T():\ngot:\n%v\nwant:\n%v", mat.Formatted(mt), mat.Formatted(dt))
		}
		vd := mat.NewVecDense(3, nil)
		v := randomVec(rnd)
		vd.SetVec(0, v.X)
		vd.SetVec(1, v.Y)
		vd.SetVec(2, v.Z)
		vd.MulVec(dt, vd)
		want := Vec{X: vd.AtVec(0), Y: vd.AtVec(1), Z: vd.AtVec(2)}
		got := m.MulVecTrans(v)
		if d := Sub(want, got); Dot(d, d) > tol {
			t.Errorf("VecDense.MulVec(dense.T()) not agree with r3.Mat.MulVec(r3.Vec): got:%v want:%v", got, want)
		}
	}
}

func randomMat(rnd *rand.Rand) *Mat {
	m := Mat{new(array)}
	for iv := 0; iv < 9; iv++ {
		i := iv / 3
		j := iv % 3
		m.Set(i, j, (rnd.Float64()-0.5)*20)
	}
	return &m
}

func randomVec(rnd *rand.Rand) (v Vec) {
	v.X = (rnd.Float64() - 0.5) * 20
	v.Y = (rnd.Float64() - 0.5) * 20
	v.Z = (rnd.Float64() - 0.5) * 20
	return v
}

func TestDet(t *testing.T) {
	const tol = 1e-11
	rnd := rand.New(rand.NewSource(1))
	for tc := 0; tc < 20; tc++ {
		m := randomMat(rnd)
		got := m.Det()
		want := mat.Det(m)
		if math.Abs(got-want) > tol {
			t.Errorf("r3.Mat.Det() not equal to mat.Det(). got %f, want %f", got, want)
		}
	}
}

func TestOuter(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for tc := 0; tc < 20; tc++ {
		alpha := rnd.Float64()
		d := mat.NewDense(3, 3, nil)
		n := NewMat(nil)
		v1 := randomVec(rnd)
		v2 := randomVec(rnd)
		d1 := mat.NewVecDense(3, []float64{v1.X, v1.Y, v1.Z})
		d2 := mat.NewVecDense(3, []float64{v2.X, v2.Y, v2.Z})
		d.Outer(alpha, d1, d2)
		n.Outer(alpha, v1, v2)
		if !mat.Equal(d, n) {
			t.Error("matrices not equal")
		}
	}
}

func TestRotationMat(t *testing.T) {
	const tol = 1e-14
	rnd := rand.New(rand.NewSource(1))

	for tc := 0; tc < 20; tc++ {
		// Generate a random unit quaternion.
		q := quat.Number{Real: rnd.NormFloat64(), Imag: rnd.NormFloat64(), Jmag: rnd.NormFloat64(), Kmag: rnd.NormFloat64()}
		q = quat.Scale(1/quat.Abs(q), q)

		// Convert it to a rotation matrix R.
		r := Rotation(q).Mat()

		// Check that the matrix has the determinant approximately equal to 1.
		diff := math.Abs(r.Det() - 1)
		if diff > tol {
			t.Errorf("case %d: unexpected determinant of R; got=%f, want=1 (diff=%v)", tc, r.Det(), diff)
			continue
		}

		// Generate a random point.
		v := Vec{X: rnd.NormFloat64(), Y: rnd.NormFloat64(), Z: rnd.NormFloat64()}
		// Rotate it using the formula q*v*conj(q).
		want := Rotation(q).Rotate(v)
		// Rotate it using the rotation matrix R.
		got := r.MulVec(v)
		// Check that |got-want| is small.
		diff = Norm(Sub(got, want))
		if diff > tol {
			t.Errorf("case %d: unexpected result; got=%f, want=%f, (diff=%v)", tc, got, want, diff)
		}
	}
}

func BenchmarkQuat(b *testing.B) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < b.N; i++ {
		q := quat.Number{Real: rnd.Float64(), Imag: rnd.Float64(), Jmag: rnd.Float64(), Kmag: rnd.Float64()}
		if Rotation(q).Mat() == nil {
			b.Fatal("nil return")
		}
	}
}

var scalarFields = []struct {
	// This is the scalar field function.
	field    func(p Vec) float64
	gradient func(p Vec) Vec
	hessian  func(p Vec) *Mat
}{
	{
		field: func(p Vec) float64 {
			// 4*x^3 + 5*y^2 + 3*z^4
			z2 := p.Z * p.Z
			return 4*p.X*p.X*p.X + 5*p.Y*p.Y + 3*z2*z2
		},
		gradient: func(p Vec) Vec {
			return Vec{X: 12 * p.X * p.X, Y: 10 * p.Y, Z: 12 * p.Z * p.Z * p.Z}
		},
		hessian: func(p Vec) *Mat {
			return NewMat([]float64{
				24 * p.X, 0, 0,
				0, 10, 0,
				0, 0, 36 * p.Z * p.Z,
			})
		},
	},
	{
		field: func(p Vec) float64 {
			// cos(x) * sin(z) * y^4
			y2 := p.Y * p.Y
			return math.Cos(p.X) * math.Sin(p.Z) * y2 * y2
		},
		gradient: func(p Vec) Vec {
			y3 := p.Y * p.Y * p.Y
			y4 := p.Y * y3
			sx, cx := math.Sincos(p.X)
			sz, cz := math.Sincos(p.Z)
			return Vec{X: -y4 * sx * sz, Y: 4 * y3 * cx * sz, Z: y4 * cx * cz}
		},
		hessian: func(p Vec) *Mat {
			y3 := p.Y * p.Y * p.Y
			y4 := y3 * p.Y
			sx, cx := math.Sincos(p.X)
			sz, cz := math.Sincos(p.Z)
			return NewMat([]float64{
				-y4 * cx * sz, -4 * y3 * sx * sz, -y4 * sx * cz,
				-4 * y3 * sx * sz, 12 * p.Y * p.Y * cx * sz, 4 * y3 * cx * cz,
				-y4 * sx * cz, 4 * y3 * cx * cz, -y4 * cx * sz,
			})
		},
	},
}

func TestMatHessian(t *testing.T) {
	const (
		tol = 1e-5
		h   = 8e-4
	)
	step := Vec{X: h, Y: h, Z: h}
	rnd := rand.New(rand.NewSource(1))
	for _, test := range scalarFields {
		for i := 0; i < 30; i++ {
			p := randomVec(rnd)
			got := NewMat(nil)
			got.Hessian(p, step, test.field)
			want := test.hessian(p)
			if !mat.EqualApprox(got, want, tol) {
				t.Errorf("matrices not equal within tol\ngot:  %v\nwant:  %v",
					mat.Formatted(got), mat.Formatted(want))
			}
		}
	}
}

func TestMatJacobian(t *testing.T) {
	const (
		tol = 1e-5
		h   = 8e-4
	)
	step := Vec{X: h, Y: h, Z: h}
	rnd := rand.New(rand.NewSource(1))
	for _, test := range vectorFields {
		for i := 0; i < 1; i++ {
			p := randomVec(rnd)
			got := NewMat(nil)
			got.Jacobian(p, step, test.field)
			want := test.jacobian(p)
			if !mat.EqualApprox(got, want, tol) {
				t.Errorf("matrices not equal within tol\ngot:  %v\nwant:  %v",
					mat.Formatted(got), mat.Formatted(want))
			}
		}
	}
}
