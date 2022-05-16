// Copyright ©2022 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r3

import (
	"math"
	"testing"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
)

func TestAffineZeroValue(t *testing.T) {
	const tol = 1e-12
	var eye Affine
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 60; i++ {
		m := randomAffine(rnd)
		got := eye.Mul(m)
		if !affineEquals(got, m, tol) {
			t.Error("zero value is not identity matrix or bad Mul")
		}
		gotZ := m.Mul(zeroAffine)
		if !affineEquals(gotZ, zeroAffine, tol) {
			t.Error("zeroAffine multiplied by random matrix should be zeroAffine")
		}
	}
}

func TestAffineInverse(t *testing.T) {
	const tol = 1e-12
	var eye Affine

	if !affineEquals(eye, eye.Inv(), tol) {
		t.Error("identity inverse must be identity")
	}
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 60; i++ {
		a := randomAffine(rnd)
		inv := a.Inv()
		got := a.Mul(inv)
		got2 := inv.Mul(a)
		if !affineEquals(got, eye, tol) || !affineEquals(got2, eye, tol) {
			t.Error("inv(A) * A should be identity")
		}
	}
}

func TestAffineTranslate(t *testing.T) {
	const tol = 1e-12
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 60; i++ {
		var a Affine
		translate := randomVec(rnd)
		a = a.AddTranslation(translate)
		vec := randomVec(rnd)
		transformed := a.Transform(vec)
		gotTranslate := Sub(transformed, vec)
		if !vecApproxEqual(gotTranslate, translate, tol) {
			t.Error("Translate failed")
		}
	}
}

func TestAffineScale(t *testing.T) {
	const tol = 1e-12
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		var a Affine
		vec := randomVec(rnd)
		scaling := randomVec(rnd)
		// Scale about origin.
		a = a.Scale(Vec{}, scaling)
		got := a.Transform(vec)
		if math.Abs(got.X/vec.X-scaling.X) > tol {
			t.Error("bad x scaling", got.X, vec.X, scaling.X)
		}
		if math.Abs(got.Y/vec.Y-scaling.Y) > tol {
			t.Error("bad y scaling", got.Y, vec.Y, scaling.Y)
		}
		if math.Abs(got.Z/vec.Z-scaling.Z) > tol {
			t.Error("bad z scaling", got.Z, vec.Z, scaling.Z)
		}
	}
}

func TestNewAffine(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		T1 := randomAffine(rnd)
		dat1 := T1.sliceCopy()
		T2 := NewAffine(dat1)
		dat2 := T2.sliceCopy()
		if len(dat1) != 16 || len(dat1) != len(dat2) {
			t.Fatal("bad slice length")
		}
		for j := range dat1 {
			if dat1[j] != dat2[j] {
				t.Error("bad slice data storage from SliceCopy or NewAffine")
			}
		}
	}
}

func TestAffineTranspose(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		T1 := randomAffine(rnd)
		Ttr := T1.transpose()
		Ttr2 := Ttr.transpose()
		if T1 != Ttr2 {
			t.Error("double transpose not return same matrix")
		}
		data := T1.sliceCopy()
		dataT := Ttr.sliceCopy()
		for i := 0; i < 4; i++ {
			for j := 0; j < 4; j++ {
				idx := i*4 + j
				idxT := j*4 + i
				if dataT[idxT] != data[idx] {
					t.Errorf("transpose data at index (%d,%d) mismatch", i, j)
				}
			}
		}
	}
}

func TestAffineMulAssociative(t *testing.T) {
	const tol = 1e-12
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 100; i++ {
		T1 := randomAffine(rnd)
		T2 := randomAffine(rnd)
		v := randomVec(rnd)
		want := T1.Transform(T2.Transform(v))
		got := T1.Mul(T2).Transform(v)
		if !vecApproxEqual(want, got, tol) {
			t.Errorf("Affine Mul should be associative, want %.4g, got %.4g", want, got)
		}
	}
}
func BenchmarkAffineInverse(b *testing.B) {
	b.StopTimer()
	rnd := rand.New(rand.NewSource(1))
	Ts := make([]Affine, b.N)
	for i := 0; i < b.N; i++ {
		Ts[i] = randomAffine(rnd)
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Ts[i] = Ts[i].Inv()
	}
}

func BenchmarkMatAffineInverse(b *testing.B) {
	b.StopTimer()
	rnd := rand.New(rand.NewSource(1))
	Ts := make([]*mat.Dense, b.N)
	for i := 0; i < b.N; i++ {
		Ts[i] = mat.NewDense(4, 4, randomAffine(rnd).sliceCopy())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		Ts[i].Inverse(Ts[i])
	}
}

func randomAffine(rnd *rand.Rand) Affine {
	var m Affine
	rot := NewRotation(math.Pi*(rnd.Float64()-0.5), randomVec(rnd))
	m = composeAffine(randomVec(rnd), rot)
	return m
}

// equals tests the equality of the Affines to within a tolerance.
func affineEquals(a, b Affine, tolerance float64) bool {
	return math.Abs(a.d00-b.d00) < tolerance &&
		math.Abs(a.x01-b.x01) < tolerance &&
		math.Abs(a.x02-b.x02) < tolerance &&
		math.Abs(a.x03-b.x03) < tolerance &&
		math.Abs(a.x10-b.x10) < tolerance &&
		math.Abs(a.d11-b.d11) < tolerance &&
		math.Abs(a.x12-b.x12) < tolerance &&
		math.Abs(a.x13-b.x13) < tolerance &&
		math.Abs(a.x20-b.x20) < tolerance &&
		math.Abs(a.x21-b.x21) < tolerance &&
		math.Abs(a.d22-b.d22) < tolerance &&
		math.Abs(a.x23-b.x23) < tolerance &&
		math.Abs(a.x30-b.x30) < tolerance &&
		math.Abs(a.x31-b.x31) < tolerance &&
		math.Abs(a.x32-b.x32) < tolerance &&
		math.Abs(a.d33-b.d33) < tolerance
}

// composeAffine creates a new affine transform for a given translation, scaling
// vector scale and quaternion rotation.
// The identity transform is constructed with
//  composeAffine(Vec{}, Warp{}, Rotation{})
func composeAffine(translate Vec, q Rotation) Affine {
	R := makeAffineRotation(q)
	R.x03 = translate.X
	R.x13 = translate.Y
	R.x23 = translate.Z
	return R
}
