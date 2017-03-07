// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmv

import (
	"log"
	"math/rand"
	"testing"

	"github.com/gonum/matrix/mat64"
)

func BenchmarkMarginalNormal10(b *testing.B) {
	sz := 10
	rnd := rand.New(rand.NewSource(1))
	normal := randomNormal(sz, rnd)
	_ = normal.CovarianceMatrix(nil) // pre-compute sigma
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		marg, ok := normal.MarginalNormal([]int{1}, nil)
		if !ok {
			b.Error("bad test")
		}
		_ = marg
	}
}

func BenchmarkMarginalNormalReset10(b *testing.B) {
	sz := 10
	rnd := rand.New(rand.NewSource(1))
	normal := randomNormal(sz, rnd)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		marg, ok := normal.MarginalNormal([]int{1}, nil)
		if !ok {
			b.Error("bad test")
		}
		_ = marg
	}
}

func BenchmarkMarginalNormalSingle10(b *testing.B) {
	sz := 10
	rnd := rand.New(rand.NewSource(1))
	normal := randomNormal(sz, rnd)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		marg := normal.MarginalNormalSingle(1, nil)
		_ = marg
	}
}

func randomNormal(sz int, rnd *rand.Rand) *Normal {
	mu := make([]float64, sz)
	for i := range mu {
		mu[i] = rnd.Float64()
	}
	data := make([]float64, sz*sz)
	for i := range data {
		data[i] = rnd.Float64()
	}
	dM := mat64.NewDense(sz, sz, data)
	var sigma mat64.SymDense
	sigma.SymOuterK(1, dM)

	normal, ok := NewNormal(mu, &sigma, nil)
	if !ok {
		log.Fatal("bad test, not pos def")
	}
	return normal
}
