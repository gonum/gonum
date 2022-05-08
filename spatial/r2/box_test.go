// Copyright ©2022 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package r2

import (
	"testing"

	"golang.org/x/exp/rand"
)

func TestBoxContains(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 200; i++ {
		b := randomBox(rnd)
		for j := 0; j < 10; j++ {
			contained := b.random(rnd)
			if !b.Contains(contained) {
				t.Error("bounding box should contain Vec")
			}
		}
		uncontained := [4]Vec{
			Add(b.Max, Vec{1, 0}),
			Add(b.Max, Vec{0, 1}),
			Sub(b.Min, Vec{1, 0}),
			Sub(b.Min, Vec{0, 1}),
		}
		for _, unc := range uncontained {
			if b.Contains(unc) {
				t.Error("box should not contain vec")
			}
		}
	}
}

func TestBoxUnion(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 200; i++ {
		b1 := randomBox(rnd)
		b2 := randomBox(rnd)
		u := b1.Union(b2)
		for j := 0; j < 10; j++ {
			contained := b1.random(rnd)
			if !u.Contains(contained) {
				t.Error("union should contain b1's Vec")
			}
			contained = b2.random(rnd)
			if !u.Contains(contained) {
				t.Error("union should contain b2's Vec")
			}
		}
		uncontained := [4]Vec{
			Add(maxElem(b1.Max, b2.Max), Vec{1, 0}),
			Add(maxElem(b1.Max, b2.Max), Vec{0, 1}),
			Sub(minElem(b1.Min, b2.Min), Vec{1, 0}),
			Sub(minElem(b1.Min, b2.Min), Vec{0, 1}),
		}
		for _, unc := range uncontained {
			if !b1.Contains(unc) && !b2.Contains(unc) && u.Contains(unc) {
				t.Error("union should not contain Vec")
			}
		}
	}
}

func TestBoxCenter(t *testing.T) {
	const tol = 1e-11
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 300; i++ {
		b := randomBox(rnd)
		center := b.Center()
		size := b.Size()
		newBox := centeredBox(center, size)
		if b.Empty() {
			t.Fatal("random box result must be well formed")
		}
		if !vecApproxEqual(b.Min, newBox.Min, tol) {
			t.Errorf("min values of box not equal. got %g, expected %g", newBox.Min, b.Min)
		}
		if !vecApproxEqual(b.Max, newBox.Max, tol) {
			t.Errorf("max values of box not equal. got %g, expected %g", newBox.Max, b.Max)
		}
	}
}

func TestBoxAdd(t *testing.T) {
	const tol = 1e-14
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 12; i++ {
		b := randomBox(rnd)
		v := randomVec(rnd)
		got := b.Add(v)
		want := Box{Min: Add(b.Min, v), Max: Add(b.Max, v)}
		if !vecApproxEqual(got.Min, want.Min, tol) {
			t.Error("box min incorrect result")
		}
		if !vecApproxEqual(got.Max, want.Max, tol) {
			t.Error("box max incorrect result")
		}
	}
}

func TestBoxScale(t *testing.T) {
	const tol = 1e-11
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 300; i++ {
		b := randomBox(rnd)
		size := b.Size()

		scaler := absElem(randomVec(rnd))
		scaled := b.Scale(scaler)
		gotScaler := divElem(scaled.Size(), size)
		if !vecApproxEqual(scaler, gotScaler, tol) {
			t.Errorf("got scaled %g, expected %g", gotScaler, scaler)
		}
		center := b.Center()
		scaledCenter := scaled.Center()
		if !vecApproxEqual(center, scaledCenter, tol) {
			t.Error("scale modified center")
		}
	}
}

func TestBoxEmpty(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 300; i++ {
		v := absElem(randomVec(rnd))
		b := randomBox(rnd)
		min := b.Min
		max := b.Max
		if !(Box{Min: min, Max: min}).Empty() {
			t.Error("Box{min,min} should be empty")
		}
		if !(Box{Min: max, Max: max}).Empty() {
			t.Error("Box{max,max} should be empty")
		}
		bmm := Box{Min: min, Max: Sub(min, v)}
		if !bmm.Empty() {
			t.Error("Box{min,min-v} should be empty")
		} else if bmm.Canon().Empty() {
			t.Error("Canonical box of Box{min,min-v} is not empty")
		}
		bMM := Box{Min: Add(max, v), Max: max}
		if !bMM.Empty() {
			t.Error("Box{max+v,max} should be empty")
		} else if bmm.Canon().Empty() {
			t.Error("Canonical box of Box{max+v,max} is not empty")
		}
	}
}
func TestBoxCanon(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 300; i++ {
		b := randomBox(rnd)
		badBox := Box{Min: b.Max, Max: b.Min}
		canon := badBox.Canon()
		if canon != b {
			t.Error("swapped box canon should be equal to original box")
		}
	}
}

func TestBoxVertices(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < 300; i++ {
		b := randomBox(rnd)
		gots := b.Vertices()
		wants := goldenVertices(b)
		if len(gots) != len(wants) {
			t.Fatalf("bad length of vertices. expect 4, got %d", len(gots))
		}
		for j, want := range wants {
			got := gots[j]
			if !vecEqual(want, got) {
				t.Errorf("%dth vertex not equal", j)
			}
		}
	}
}

// randomBox returns a random valid bounding Box.
func randomBox(rnd *rand.Rand) Box {
	spatialScale := randomRange(0, 2000)
	boxScale := randomRange(0.01, 1000)
	return centeredBox(Scale(spatialScale, randomVec(rnd)), Scale(boxScale, absElem(randomVec(rnd))))
}

// Random returns a random point within the Box.
// used to facilitate testing
func (b Box) random(rnd *rand.Rand) Vec {
	return Vec{
		X: randomRange(b.Min.X, b.Max.X),
		Y: randomRange(b.Min.Y, b.Max.Y),
	}
}

// randomRange returns a random float64 [a,b)
func randomRange(a, b float64) float64 {
	return a + (b-a)*rand.Float64()
}

func goldenVertices(a Box) []Vec {
	return []Vec{
		0: a.Min,
		1: {a.Max.X, a.Min.Y},
		2: a.Max,
		3: {a.Min.X, a.Max.Y},
	}
}
