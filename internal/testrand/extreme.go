// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testrand

import (
	"math"
	"math/rand/v2"
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
		return extremeFloat64Exp[e.rnd.IntN(len(extremeFloat64Exp))]
	case e.nan():
		return math.NaN()
	}

	return e.rnd.ExpFloat64()
}

// Float32 returns a pseudo-random float32 in range [0.0, 1.0).
func (e *extreme) Float32() float32 {
	switch {
	case e.p():
		return extremeFloat32Unit[e.rnd.IntN(len(extremeFloat32Unit))]
	case e.nan():
		return float32(math.NaN())
	}

	return e.rnd.Float32()
}

// Float64 returns a pseudo-random float64 in range [0.0, 1.0).
func (e *extreme) Float64() float64 {
	switch {
	case e.p():
		return extremeFloat64Unit[e.rnd.IntN(len(extremeFloat64Unit))]
	case e.nan():
		return math.NaN()
	}

	return e.rnd.Float64()
}

// Int returns a non-negative pseudo-random int.
func (e *extreme) Int() int {
	if e.p() {
		return extremeInt[e.rnd.IntN(len(extremeInt))]
	}
	return e.rnd.Int()
}

// Int32 returns a non-negative pseudo-random int32.
func (e *extreme) Int32() int32 {
	if e.p() {
		return extremeInt31[e.rnd.IntN(len(extremeInt31))]
	}
	return e.rnd.Int32()
}

// Int32N returns a non-negative pseudo-random int32 from range [0, n).
func (e *extreme) Int32N(n int32) int32 {
	if e.p() {
		switch rand.IntN(4) {
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
	return e.rnd.Int32N(n)
}

// Int64 returns a non-negative pseudo-random int64.
func (e *extreme) Int64() int64 {
	if e.p() {
		return extremeInt63[e.rnd.IntN(len(extremeInt63))]
	}
	return e.rnd.Int64()
}

// Int64N returns a non-negative pseudo-random int from range [0, n).
func (e *extreme) Int64N(n int64) int64 {
	if e.p() {
		switch rand.IntN(4) {
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
	return e.rnd.Int64N(n)
}

// IntN returns a non-negative pseudo-random int from range [0, n).
func (e *extreme) IntN(n int) int {
	if e.p() {
		switch rand.IntN(4) {
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
	return e.rnd.IntN(n)
}

// NormFloat64 returns a normally distributed pseudo-random float64 in range [-math.MaxFloat64, math.MaxFloat64].
func (e *extreme) NormFloat64() float64 {
	switch {
	case e.p():
		return extremeFloat64Norm[e.rnd.IntN(len(extremeFloat64Norm))]
	case e.nan():
		return math.NaN()
	}

	return e.rnd.NormFloat64()
}

// Uint32 returns a pseudo-random uint32.
func (e *extreme) Uint32() uint32 {
	if e.p() {
		return extremeUint32[e.rnd.IntN(len(extremeUint32))]
	}
	return e.rnd.Uint32()
}

// Uint64 returns a pseudo-random uint64.
func (e *extreme) Uint64() uint64 {
	if e.p() {
		return extremeUint64[e.rnd.IntN(len(extremeUint64))]
	}
	return e.rnd.Uint64()
}

// Uint64N returns a pseudo-random uint64 from range [0, n).
func (e *extreme) Uint64N(n uint64) uint64 {
	if e.p() {
		switch rand.IntN(4) {
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
	return e.rnd.Uint64N(n)
}
