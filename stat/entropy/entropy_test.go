// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package entropy

import (
	"math"
	"testing"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
)

func TestShannon(t *testing.T) {
	t.Log("Testing Entropy")
	p1 := mat.NewVecDense(4, []float64{0.5, 0.5, 0.5, 0.5})

	if r := Shannon(p1); math.Abs(r-1.386294) > 0.001 {
		t.Errorf("Entropy of four state uniform distribution should be 1.386294 but it is %f", r)
	}

	p2 := mat.NewVecDense(4, []float64{1.0, 0.0, 0.0, 0.0})

	if r := Shannon(p2); r != 0.0 {
		t.Errorf("Entropy of deterministic distribution should be 0.0 but it is %f", r)
	}
}

func TestEntropyChaoShen(t *testing.T) {
	t.Log("Testing Chao-Shen Entropy")
	r := 0.0
	for i := 0; i < 100; i++ {
		h := make([]int, 5000, 5000)
		for j := 0; j < 5000; j++ {
			h[j] = int(rand.Intn(100))
		}
		r += ChaoShen(h)
	}

	r /= 100.0

	if math.Abs(r-4.595091) > 0.1 {
		t.Errorf("Entropy should be 4.595091 and not %f", r)
	}

}

func TestEntropyMLBC(t *testing.T) {
	t.Log("Testing Maximum Likelihood Bias Corrected")
	r := 0.0
	for i := 0; i < 100; i++ {
		h := make([]int, 5000, 5000)
		for j := 0; j < 5000; j++ {
			h[j] = int(rand.Int63n(100))
		}
		r += MLBC(h)
	}

	r /= 100.0

	if math.Abs(r-4.604982) > 0.1 {
		t.Errorf("Entropy should be 4.604982 and not %f", r)
	}
}

func TestEntropyHorvitzThompson(t *testing.T) {
	t.Log("Testing Horvitz-Thompson Base 2")
	r := 0.0
	for i := 0; i < 100; i++ {
		h := make([]int, 5000, 5000)
		for j := 0; j < 5000; j++ {
			h[j] = int(rand.Int63n(100))
		}
		r += HorvitzThompson(h)
	}

	r /= 100.0

	if math.Abs(r-4.595032) > 0.1 {
		t.Errorf("Entropy should be 4.595032 and not %f", r)
	}
}
