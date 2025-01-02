// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package testrand provides random generation and flags for testing.
package testrand

import (
	"flag"
	"math/rand/v2"
)

var (
	seedFlag    = flag.Uint64("testrand.seed", 1, "random seed for tests (0=randomize)")
	extremeFlag = flag.Float64("testrand.extreme", 0, "probability of returning extreme values")
	nanFlag     = flag.Float64("testrand.nan", 0, "probability of returning nan")
)

// TB is an interface that corresponds to a subset of *testing.T and *testing.B.
type TB interface {
	Logf(format string, args ...interface{})
}

// Source corresponds to the interface in golang.org/x/exp/rand.Source.
type Source = rand.Source

// Rand corresponds to math/rand/v2.Rand.
type Rand interface {
	ExpFloat64() float64
	Float32() float32
	Float64() float64
	Int() int
	Int32() int32
	Int32N(n int32) int32
	Int64() int64
	Int64N(n int64) int64
	IntN(n int) int
	NormFloat64() float64
	Perm(n int) []int
	Shuffle(n int, swap func(i, j int))
	Uint32() uint32
	Uint64() uint64
	Uint64N(n uint64) uint64
}

// New returns a new random number generator using the global flags.
func New(tb TB) Rand {
	seed := *seedFlag
	if seed == 0 {
		seed = rand.Uint64()
	}

	// Don't log the default case.
	if seed == 1 && *extremeFlag == 0 && *nanFlag == 0 {
		base := rand.New(rand.NewPCG(seed, seed))
		return base
	}

	tb.Logf("seed=%d, prob=%.2f, nan=%.2f", seed, *extremeFlag, *nanFlag)

	base := rand.New(rand.NewPCG(seed, seed))
	if *extremeFlag <= 0 && *nanFlag <= 0 {
		return base
	}

	return newExtreme(*extremeFlag, *nanFlag, base)
}

// NewSource returns a new source for random numbers.
func NewSource(tb TB) Source {
	seed := *seedFlag
	if seed == 0 {
		seed = rand.Uint64()
	}

	// Don't log the default case.
	if seed == 1 {
		return rand.NewPCG(seed, seed)
	}

	tb.Logf("seed %d", seed)
	return rand.NewPCG(seed, seed)
}
