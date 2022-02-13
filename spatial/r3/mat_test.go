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
		v1 := randomVec(rnd)
		v2 := randomVec(rnd)
		sk := Skew(v1)
		want := Cross(v1, v2)
		got := sk.MulVec(v2)
		if d := want.Sub(got); d.Dot(d) > tol {
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
		if d := want.Sub(got); d.Dot(d) > tol {
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

func TestRotationFromQuat(t *testing.T) {
	const tol = 1e-11
	rnd := rand.New(rand.NewSource(1))
	var backing [9]float64 // reuse memory.
	for tc := 0; tc < 20; tc++ {
		q := quat.Number{Real: rnd.Float64(), Imag: rnd.Float64(), Jmag: rnd.Float64(), Kmag: rnd.Float64()}
		qabs := quat.Abs(q)
		q = quat.Scale(1/qabs, q)
		m := NewMat(backing[:])
		m.RotationFromQuat(q)
		w, x, y, z := q.Real, q.Imag, q.Jmag, q.Kmag
		x2, y2, z2 := x*x, y*y, z*z
		norm := math.Sqrt(w*w + x2 + y2 + z2)
		_ = norm
		expect := NewMat([]float64{ // From https://en.wikipedia.org/wiki/Rotation_matrix#Quaternion
			1 - 2*y2 - 2*z2, 2*x*y - 2*z*w, 2*x*z + 2*y*w,
			2*x*y + 2*z*w, 1 - 2*x2 - 2*z2, 2*y*z - 2*x*w,
			2*x*z - 2*y*w, 2*y*z + 2*x*w, 1 - 2*x2 - 2*y2,
		})
		if !mat.EqualApprox(m, expect, tol) {
			t.Errorf("Out of tolerance.")
		}
		det := m.Det()
		if math.Abs(det-1) > tol {
			t.Errorf("determinant expected to be 1. got %f", det)
		}
	}
}

func BenchmarkQuat(b *testing.B) {
	rnd := rand.New(rand.NewSource(1))
	m := NewMat(nil)
	for i := 0; i < b.N; i++ {
		q := quat.Number{Real: rnd.Float64(), Imag: rnd.Float64(), Jmag: rnd.Float64(), Kmag: rnd.Float64()}
		m.RotationFromQuat(q)
	}
}
