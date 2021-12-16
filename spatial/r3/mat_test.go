// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3

import (
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/mat"
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
