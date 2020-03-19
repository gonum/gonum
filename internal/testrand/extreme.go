// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testrand

import (
	"math"

	"golang.org/x/exp/rand"
)

// extreme is a pseudo-random number generator that has high probability of returning extreme values.
type extreme struct {
	probability    float64
	nanProbability float64

	rnd Rand
}

// newExtreme creates a new extreme pseudo-random generator.
// p is the probability of returning an extreme value.
// nan is the probability of returning a NaN.
func newExtreme(p, nan float64, rnd Rand) *extreme {
	return &extreme{p, nan, rnd}
}

// Perm returns a permutation of integers [0, n).
func (e *extreme) Perm(n int) []int { return e.rnd.Perm(n) }

// Read generates len(p) pseudo-random bytes.
func (e *extreme) Read(p []byte) (n int, err error) { return e.rnd.Read(p) }

// Seed reseeds the pseudo-random generator.
func (e *extreme) Seed(seed uint64) { e.rnd.Seed(seed) }

// Shuffle shuffles n items using the swap callback.
func (e *extreme) Shuffle(n int, swap func(i, j int)) { e.rnd.Shuffle(n, swap) }

// p returns true when the generator should output an extreme value.
func (e *extreme) p() bool {
	if e.probability <= 0 {
		return false
	}
	return e.rnd.Float64() < e.probability
}

// nan returns true when the generator should output nan.
func (e *extreme) nan() bool {
	if e.nanProbability <= 0 {
		return false
	}
	return e.rnd.Float64() < e.nanProbability
}

// ExpFloat64 returns an exponentialy distributed pseudo-random float64 in range (0, math.MaxFloat64].
func (e *extreme) ExpFloat64() float64 {
	switch {
	case e.p():
		return extremeFloat64Exp[e.rnd.Intn(len(extremeFloat64Exp))]
	case e.nan():
		return math.NaN()
	}

	return e.rnd.ExpFloat64()
}

// Float32 returns a pseudo-random float32 in range [0.0, 1.0).
func (e *extreme) Float32() float32 {
	switch {
	case e.p():
		return extremeFloat32Unit[e.rnd.Intn(len(extremeFloat32Unit))]
	case e.nan():
		return float32(math.NaN())
	}

	return e.rnd.Float32()
}

// Float64 returns a pseudo-random float64 in range [0.0, 1.0).
func (e *extreme) Float64() float64 {
	switch {
	case e.p():
		return extremeFloat64Unit[e.rnd.Intn(len(extremeFloat64Unit))]
	case e.nan():
		return math.NaN()
	}

	return e.Float64()
}

// Int returns a non-negative pseudo-random int.
func (e *extreme) Int() int {
	if e.p() {
		return extremeInt[e.rnd.Intn(len(extremeInt))]
	}
	return e.Int()
}

// Int31 returns a non-negative pseudo-random int32.
func (e *extreme) Int31() int32 {
	if e.p() {
		return extremeInt31[e.rnd.Intn(len(extremeInt31))]
	}
	return e.Int31()
}

// Int31n returns a non-negative pseudo-random int32 from range [0, n).
func (e *extreme) Int31n(n int32) int32 {
	if e.p() {
		switch rand.Intn(4) {
		case 0:
			return 0
		case 1:
			return 1
		case 2:
			return n / 2
		case 3:
			return n - 1
		}
	}
	return e.Int31n(n)
}

// Int63 returns a non-negative pseudo-random int64.
func (e *extreme) Int63() int64 {
	if e.p() {
		return extremeInt63[e.rnd.Intn(len(extremeInt63))]
	}
	return e.Int63()
}

// Int63n returns a non-negative pseudo-random int from range [0, n).
func (e *extreme) Int63n(n int64) int64 {
	if e.p() {
		switch rand.Intn(4) {
		case 0:
			return 0
		case 1:
			return 1
		case 2:
			return n / 2
		case 3:
			return n - 1
		}
	}
	return e.Int63n(n)
}

// Int returns a non-negative pseudo-random int from range [0, n).
func (e *extreme) Intn(n int) int {
	if e.p() {
		switch rand.Intn(4) {
		case 0:
			return 0
		case 1:
			return 1
		case 2:
			return n / 2
		case 3:
			return n - 1
		}
	}
	return e.Intn(n)
}

// NormFloat64 returns a normally distributed pseudo-random float64 in range [-math.MaxFloat64, math.MaxFloat64].
func (e *extreme) NormFloat64() float64 {
	switch {
	case e.p():
		return extremeFloat64Norm[e.rnd.Intn(len(extremeFloat64Norm))]
	case e.nan():
		return math.NaN()
	}

	return e.NormFloat64()
}

// Uint32 returns a pseudo-random uint32.
func (e *extreme) Uint32() uint32 {
	if e.p() {
		return extremeUint32[e.rnd.Intn(len(extremeUint32))]
	}
	return e.Uint32()
}

// Uint32 returns a pseudo-random uint64.
func (e *extreme) Uint64() uint64 {
	if e.p() {
		return extremeUint64[e.rnd.Intn(len(extremeUint64))]
	}
	return e.Uint64()
}

// Uint64n returns a pseudo-random uint64 from range [0, n).
func (e *extreme) Uint64n(n uint64) uint64 {
	if e.p() {
		switch rand.Intn(4) {
		case 0:
			return 0
		case 1:
			return 1
		case 2:
			return n / 2
		case 3:
			return n - 1
		}
	}
	return e.Uint64n(n)
}
