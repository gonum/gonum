// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package testfunc holds various functions to test autofd features.
package testfunc // import "gonum.org/v1/gonum/diff/autofd/internal/testfunc"

import "math"

func F1(x float64) float64 {
	return x * x
}

func F2(y float64) float64 {
	return y * y
}

func F3(x float64) float64 {
	return 2 * x * x
}

func F4(x float64) float64 {
	return 2 / (x * x)
}

func F5(x float64) float64 {
	return 2 / (x * -x)
}

func F6(x float64) float64 {
	return 2 + x - x
}

func F7(x float64) float64 {
	return math.Cos(2 * math.Pi * x)
}

func F8(x float64) float64 {
	return math.Exp(x) / math.Sqrt(math.Pow(math.Sin(x), 3)+math.Pow(math.Cos(x), 3))
}

type T1 struct{}

func (T1) F(x float64) float64 {
	return 2*x + 3*x*x + 4*math.Pow(x, 3)
}

type T2 = T1

func ErrF1(x, y float64) float64 {
	return x + y
}

func ErrF2(x float32) float64 {
	return float64(x)
}

func ErrF3(x float64) (float64, float64) {
	return x, x * x
}

func ErrF4(x float64) float32 {
	return float32(x)
}

func ErrF5(x float64) (o float64) {
	return
}

func ErrF6(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func ErrF7(x float64) float64 {
	switch {
	case x < 0:
		return -x
	default:

		return x
	}
}

func ErrF8(x float64) float64 {
	for i := 0; i < 10; i++ {
		if i == 5 {
			return -x
		}
	}
	return x
}

type ErrT1 struct {
	F float64
}
