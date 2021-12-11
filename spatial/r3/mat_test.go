// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3

import (
	"math"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/mat"
)

func TestMatScale(t *testing.T) {
	const tol = 1e-12
	rnd := rand.New(rand.NewSource(1))
	for tc := 0; tc < 20; tc++ {
		v := rnd.Float64()
		a := randomMat(rnd)
		gotmat := NewMat(nil)
		gotmat.Scale(v, a)
		for iv := range a.data {
			i := iv / 3
			j := iv % 3
			expect := v * a.At(i, j)
			got := gotmat.At(i, j)
			if math.Abs(got-expect) > tol {
				t.Errorf(
					"case %d: got=%v, want=%v",
					tc, got, expect)
			}
		}
	}
}

func TestMatCloneFrom(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for tc := 0; tc < 20; tc++ {
		a := randomMat(rnd)
		gotmat := NewMat(nil)
		gotmat.CloneFrom(a)
		if !mat.Equal(a, gotmat) {
			t.Error("Clonefrom fail")
		}
	}
}

func TestSkew(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for tc := 0; tc < 20; tc++ {
		v1 := randomVec(rnd)
		v2 := randomVec(rnd)
		sk := Skew(v1)
		got := sk.MulVec(v2)
		expect := Cross(v1, v2)
		if got != expect {
			t.Error("r3.Cross(v1,v2) not match with r3.Skew(v1)*v2")
		}
	}
}

func TestTranspose(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for tc := 0; tc < 20; tc++ {
		d := mat.NewDense(3, 3, nil)
		m := randomMat(rnd)
		d.CloneFrom(m)
		mt := m.T()
		dt := d.T()
		if !mat.Equal(mt, dt) {
			t.Error("Dense.T() not equal to r3.Mat.T()")
		}
		vd := mat.NewVecDense(3, nil)
		v := randomVec(rnd)
		vd.SetVec(0, v.X)
		vd.SetVec(1, v.Y)
		vd.SetVec(2, v.Z)
		got := m.MulVecTrans(v)
		vd.MulVec(dt, vd)
		if vd.AtVec(0) != got.X || vd.AtVec(1) != got.Y || vd.AtVec(2) != got.Z {
			t.Error("VecDense.MulVec(dense.T()) not equal to r3.Mat.MulVec(r3.Vec)")
		}
	}
}

func randomMat(rnd *rand.Rand) *Mat {
	m := Mat{data: new([3][3]float64)}
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
