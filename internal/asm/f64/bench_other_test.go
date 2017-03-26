// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package f64

import (
	"math"
	"testing"
)

func benchAbsSum(f func(x []float64) float64, sz int, t *testing.B) {
	dst := y[:sz]
	for i := 0; i < t.N; i++ {
		f(dst)
	}
}

var naiveAbsSum = func(x []float64) (sum float64) {
	for _, v := range x {
		sum += math.Abs(v)
	}
	return sum
}

func BenchmarkAbsSum1(t *testing.B)      { benchAbsSum(AbsSum, 1, t) }
func BenchmarkAbsSum2(t *testing.B)      { benchAbsSum(AbsSum, 2, t) }
func BenchmarkAbsSum3(t *testing.B)      { benchAbsSum(AbsSum, 3, t) }
func BenchmarkAbsSum4(t *testing.B)      { benchAbsSum(AbsSum, 4, t) }
func BenchmarkAbsSum5(t *testing.B)      { benchAbsSum(AbsSum, 5, t) }
func BenchmarkAbsSum10(t *testing.B)     { benchAbsSum(AbsSum, 10, t) }
func BenchmarkAbsSum100(t *testing.B)    { benchAbsSum(AbsSum, 100, t) }
func BenchmarkAbsSum1000(t *testing.B)   { benchAbsSum(AbsSum, 1000, t) }
func BenchmarkAbsSum10000(t *testing.B)  { benchAbsSum(AbsSum, 10000, t) }
func BenchmarkAbsSum100000(t *testing.B) { benchAbsSum(AbsSum, 100000, t) }
func BenchmarkAbsSum500000(t *testing.B) { benchAbsSum(AbsSum, 500000, t) }

func BenchmarkLAbsSum1(t *testing.B)      { benchAbsSum(naiveAbsSum, 1, t) }
func BenchmarkLAbsSum2(t *testing.B)      { benchAbsSum(naiveAbsSum, 2, t) }
func BenchmarkLAbsSum3(t *testing.B)      { benchAbsSum(naiveAbsSum, 3, t) }
func BenchmarkLAbsSum4(t *testing.B)      { benchAbsSum(naiveAbsSum, 4, t) }
func BenchmarkLAbsSum5(t *testing.B)      { benchAbsSum(naiveAbsSum, 5, t) }
func BenchmarkLAbsSum10(t *testing.B)     { benchAbsSum(naiveAbsSum, 10, t) }
func BenchmarkLAbsSum100(t *testing.B)    { benchAbsSum(naiveAbsSum, 100, t) }
func BenchmarkLAbsSum1000(t *testing.B)   { benchAbsSum(naiveAbsSum, 1000, t) }
func BenchmarkLAbsSum10000(t *testing.B)  { benchAbsSum(naiveAbsSum, 10000, t) }
func BenchmarkLAbsSum100000(t *testing.B) { benchAbsSum(naiveAbsSum, 100000, t) }
func BenchmarkLAbsSum500000(t *testing.B) { benchAbsSum(naiveAbsSum, 500000, t) }

func benchAbsSumInc(t *testing.B, ln, inc int, f func(x []float64, n, incX int) float64) {
	for i := 0; i < t.N; i++ {
		f(x, ln, inc)
	}
}

var naiveAbsSumInc = func(x []float64, n, incX int) (sum float64) {
	for i := 0; i < n*incX; i += incX {
		sum += math.Abs(x[i])
	}
	return sum
}

func BenchmarkF64AbsSumIncN1Inc1(b *testing.B) { benchAbsSumInc(b, 1, 1, AbsSumInc) }

func BenchmarkF64AbsSumIncN2Inc1(b *testing.B)  { benchAbsSumInc(b, 2, 1, AbsSumInc) }
func BenchmarkF64AbsSumIncN2Inc2(b *testing.B)  { benchAbsSumInc(b, 2, 2, AbsSumInc) }
func BenchmarkF64AbsSumIncN2Inc4(b *testing.B)  { benchAbsSumInc(b, 2, 4, AbsSumInc) }
func BenchmarkF64AbsSumIncN2Inc10(b *testing.B) { benchAbsSumInc(b, 2, 10, AbsSumInc) }

func BenchmarkF64AbsSumIncN3Inc1(b *testing.B)  { benchAbsSumInc(b, 3, 1, AbsSumInc) }
func BenchmarkF64AbsSumIncN3Inc2(b *testing.B)  { benchAbsSumInc(b, 3, 2, AbsSumInc) }
func BenchmarkF64AbsSumIncN3Inc4(b *testing.B)  { benchAbsSumInc(b, 3, 4, AbsSumInc) }
func BenchmarkF64AbsSumIncN3Inc10(b *testing.B) { benchAbsSumInc(b, 3, 10, AbsSumInc) }

func BenchmarkF64AbsSumIncN4Inc1(b *testing.B)  { benchAbsSumInc(b, 4, 1, AbsSumInc) }
func BenchmarkF64AbsSumIncN4Inc2(b *testing.B)  { benchAbsSumInc(b, 4, 2, AbsSumInc) }
func BenchmarkF64AbsSumIncN4Inc4(b *testing.B)  { benchAbsSumInc(b, 4, 4, AbsSumInc) }
func BenchmarkF64AbsSumIncN4Inc10(b *testing.B) { benchAbsSumInc(b, 4, 10, AbsSumInc) }

func BenchmarkF64AbsSumIncN10Inc1(b *testing.B)  { benchAbsSumInc(b, 10, 1, AbsSumInc) }
func BenchmarkF64AbsSumIncN10Inc2(b *testing.B)  { benchAbsSumInc(b, 10, 2, AbsSumInc) }
func BenchmarkF64AbsSumIncN10Inc4(b *testing.B)  { benchAbsSumInc(b, 10, 4, AbsSumInc) }
func BenchmarkF64AbsSumIncN10Inc10(b *testing.B) { benchAbsSumInc(b, 10, 10, AbsSumInc) }

func BenchmarkF64AbsSumIncN1000Inc1(b *testing.B)  { benchAbsSumInc(b, 1000, 1, AbsSumInc) }
func BenchmarkF64AbsSumIncN1000Inc2(b *testing.B)  { benchAbsSumInc(b, 1000, 2, AbsSumInc) }
func BenchmarkF64AbsSumIncN1000Inc4(b *testing.B)  { benchAbsSumInc(b, 1000, 4, AbsSumInc) }
func BenchmarkF64AbsSumIncN1000Inc10(b *testing.B) { benchAbsSumInc(b, 1000, 10, AbsSumInc) }

func BenchmarkF64AbsSumIncN100000Inc1(b *testing.B)  { benchAbsSumInc(b, 100000, 1, AbsSumInc) }
func BenchmarkF64AbsSumIncN100000Inc2(b *testing.B)  { benchAbsSumInc(b, 100000, 2, AbsSumInc) }
func BenchmarkF64AbsSumIncN100000Inc4(b *testing.B)  { benchAbsSumInc(b, 100000, 4, AbsSumInc) }
func BenchmarkF64AbsSumIncN100000Inc10(b *testing.B) { benchAbsSumInc(b, 100000, 10, AbsSumInc) }

func BenchmarkLF64AbsSumIncN1Inc1(b *testing.B) { benchAbsSumInc(b, 1, 1, naiveAbsSumInc) }

func BenchmarkLF64AbsSumIncN2Inc1(b *testing.B)  { benchAbsSumInc(b, 2, 1, naiveAbsSumInc) }
func BenchmarkLF64AbsSumIncN2Inc2(b *testing.B)  { benchAbsSumInc(b, 2, 2, naiveAbsSumInc) }
func BenchmarkLF64AbsSumIncN2Inc4(b *testing.B)  { benchAbsSumInc(b, 2, 4, naiveAbsSumInc) }
func BenchmarkLF64AbsSumIncN2Inc10(b *testing.B) { benchAbsSumInc(b, 2, 10, naiveAbsSumInc) }

func BenchmarkLF64AbsSumIncN3Inc1(b *testing.B)  { benchAbsSumInc(b, 3, 1, naiveAbsSumInc) }
func BenchmarkLF64AbsSumIncN3Inc2(b *testing.B)  { benchAbsSumInc(b, 3, 2, naiveAbsSumInc) }
func BenchmarkLF64AbsSumIncN3Inc4(b *testing.B)  { benchAbsSumInc(b, 3, 4, naiveAbsSumInc) }
func BenchmarkLF64AbsSumIncN3Inc10(b *testing.B) { benchAbsSumInc(b, 3, 10, naiveAbsSumInc) }

func BenchmarkLF64AbsSumIncN4Inc1(b *testing.B)  { benchAbsSumInc(b, 4, 1, naiveAbsSumInc) }
func BenchmarkLF64AbsSumIncN4Inc2(b *testing.B)  { benchAbsSumInc(b, 4, 2, naiveAbsSumInc) }
func BenchmarkLF64AbsSumIncN4Inc4(b *testing.B)  { benchAbsSumInc(b, 4, 4, naiveAbsSumInc) }
func BenchmarkLF64AbsSumIncN4Inc10(b *testing.B) { benchAbsSumInc(b, 4, 10, naiveAbsSumInc) }

func BenchmarkLF64AbsSumIncN10Inc1(b *testing.B)  { benchAbsSumInc(b, 10, 1, naiveAbsSumInc) }
func BenchmarkLF64AbsSumIncN10Inc2(b *testing.B)  { benchAbsSumInc(b, 10, 2, naiveAbsSumInc) }
func BenchmarkLF64AbsSumIncN10Inc4(b *testing.B)  { benchAbsSumInc(b, 10, 4, naiveAbsSumInc) }
func BenchmarkLF64AbsSumIncN10Inc10(b *testing.B) { benchAbsSumInc(b, 10, 10, naiveAbsSumInc) }

func BenchmarkLF64AbsSumIncN1000Inc1(b *testing.B)  { benchAbsSumInc(b, 1000, 1, naiveAbsSumInc) }
func BenchmarkLF64AbsSumIncN1000Inc2(b *testing.B)  { benchAbsSumInc(b, 1000, 2, naiveAbsSumInc) }
func BenchmarkLF64AbsSumIncN1000Inc4(b *testing.B)  { benchAbsSumInc(b, 1000, 4, naiveAbsSumInc) }
func BenchmarkLF64AbsSumIncN1000Inc10(b *testing.B) { benchAbsSumInc(b, 1000, 10, naiveAbsSumInc) }

func BenchmarkLF64AbsSumIncN100000Inc1(b *testing.B)  { benchAbsSumInc(b, 100000, 1, naiveAbsSumInc) }
func BenchmarkLF64AbsSumIncN100000Inc2(b *testing.B)  { benchAbsSumInc(b, 100000, 2, naiveAbsSumInc) }
func BenchmarkLF64AbsSumIncN100000Inc4(b *testing.B)  { benchAbsSumInc(b, 100000, 4, naiveAbsSumInc) }
func BenchmarkLF64AbsSumIncN100000Inc10(b *testing.B) { benchAbsSumInc(b, 100000, 10, naiveAbsSumInc) }

func benchAdd(f func(dst, s []float64), sz int, t *testing.B) {
	dst, s := y[:sz], x[:sz]
	for i := 0; i < t.N; i++ {
		f(dst, s)
	}
}

var naiveAdd = func(dst, s []float64) {
	for i, v := range s {
		dst[i] += v
	}
}

func BenchmarkAdd1(t *testing.B)      { benchAdd(Add, 1, t) }
func BenchmarkAdd2(t *testing.B)      { benchAdd(Add, 2, t) }
func BenchmarkAdd3(t *testing.B)      { benchAdd(Add, 3, t) }
func BenchmarkAdd4(t *testing.B)      { benchAdd(Add, 4, t) }
func BenchmarkAdd5(t *testing.B)      { benchAdd(Add, 5, t) }
func BenchmarkAdd10(t *testing.B)     { benchAdd(Add, 10, t) }
func BenchmarkAdd100(t *testing.B)    { benchAdd(Add, 100, t) }
func BenchmarkAdd1000(t *testing.B)   { benchAdd(Add, 1000, t) }
func BenchmarkAdd10000(t *testing.B)  { benchAdd(Add, 10000, t) }
func BenchmarkAdd100000(t *testing.B) { benchAdd(Add, 100000, t) }
func BenchmarkAdd500000(t *testing.B) { benchAdd(Add, 500000, t) }

func BenchmarkLAdd1(t *testing.B)      { benchAdd(naiveAdd, 1, t) }
func BenchmarkLAdd2(t *testing.B)      { benchAdd(naiveAdd, 2, t) }
func BenchmarkLAdd3(t *testing.B)      { benchAdd(naiveAdd, 3, t) }
func BenchmarkLAdd4(t *testing.B)      { benchAdd(naiveAdd, 4, t) }
func BenchmarkLAdd5(t *testing.B)      { benchAdd(naiveAdd, 5, t) }
func BenchmarkLAdd10(t *testing.B)     { benchAdd(naiveAdd, 10, t) }
func BenchmarkLAdd100(t *testing.B)    { benchAdd(naiveAdd, 100, t) }
func BenchmarkLAdd1000(t *testing.B)   { benchAdd(naiveAdd, 1000, t) }
func BenchmarkLAdd10000(t *testing.B)  { benchAdd(naiveAdd, 10000, t) }
func BenchmarkLAdd100000(t *testing.B) { benchAdd(naiveAdd, 100000, t) }
func BenchmarkLAdd500000(t *testing.B) { benchAdd(naiveAdd, 500000, t) }

func benchAddConst(f func(a float64, x []float64), sz int, t *testing.B) {
	a, x := 1., x[:sz]
	for i := 0; i < t.N; i++ {
		f(a, x)
	}
}

var naiveAddConst = func(a float64, x []float64) {
	for i := range x {
		x[i] += a
	}
}

func BenchmarkAddConst1(t *testing.B)      { benchAddConst(AddConst, 1, t) }
func BenchmarkAddConst2(t *testing.B)      { benchAddConst(AddConst, 2, t) }
func BenchmarkAddConst3(t *testing.B)      { benchAddConst(AddConst, 3, t) }
func BenchmarkAddConst4(t *testing.B)      { benchAddConst(AddConst, 4, t) }
func BenchmarkAddConst5(t *testing.B)      { benchAddConst(AddConst, 5, t) }
func BenchmarkAddConst10(t *testing.B)     { benchAddConst(AddConst, 10, t) }
func BenchmarkAddConst100(t *testing.B)    { benchAddConst(AddConst, 100, t) }
func BenchmarkAddConst1000(t *testing.B)   { benchAddConst(AddConst, 1000, t) }
func BenchmarkAddConst10000(t *testing.B)  { benchAddConst(AddConst, 10000, t) }
func BenchmarkAddConst100000(t *testing.B) { benchAddConst(AddConst, 100000, t) }
func BenchmarkAddConst500000(t *testing.B) { benchAddConst(AddConst, 500000, t) }

func BenchmarkLAddConst1(t *testing.B)      { benchAddConst(naiveAddConst, 1, t) }
func BenchmarkLAddConst2(t *testing.B)      { benchAddConst(naiveAddConst, 2, t) }
func BenchmarkLAddConst3(t *testing.B)      { benchAddConst(naiveAddConst, 3, t) }
func BenchmarkLAddConst4(t *testing.B)      { benchAddConst(naiveAddConst, 4, t) }
func BenchmarkLAddConst5(t *testing.B)      { benchAddConst(naiveAddConst, 5, t) }
func BenchmarkLAddConst10(t *testing.B)     { benchAddConst(naiveAddConst, 10, t) }
func BenchmarkLAddConst100(t *testing.B)    { benchAddConst(naiveAddConst, 100, t) }
func BenchmarkLAddConst1000(t *testing.B)   { benchAddConst(naiveAddConst, 1000, t) }
func BenchmarkLAddConst10000(t *testing.B)  { benchAddConst(naiveAddConst, 10000, t) }
func BenchmarkLAddConst100000(t *testing.B) { benchAddConst(naiveAddConst, 100000, t) }
func BenchmarkLAddConst500000(t *testing.B) { benchAddConst(naiveAddConst, 500000, t) }

func benchCumSum(f func(a, b []float64) []float64, sz int, t *testing.B) {
	a, b := x[:sz], y[:sz]
	for i := 0; i < t.N; i++ {
		f(a, b)
	}
}

var naiveCumSum = func(dst, s []float64) []float64 {
	if len(s) == 0 {
		return dst
	}
	dst[0] = s[0]
	for i, v := range s[1:] {
		dst[i+1] = dst[i] + v
	}
	return dst
}

func BenchmarkCumSum1(t *testing.B)      { benchCumSum(CumSum, 1, t) }
func BenchmarkCumSum2(t *testing.B)      { benchCumSum(CumSum, 2, t) }
func BenchmarkCumSum3(t *testing.B)      { benchCumSum(CumSum, 3, t) }
func BenchmarkCumSum4(t *testing.B)      { benchCumSum(CumSum, 4, t) }
func BenchmarkCumSum5(t *testing.B)      { benchCumSum(CumSum, 5, t) }
func BenchmarkCumSum10(t *testing.B)     { benchCumSum(CumSum, 10, t) }
func BenchmarkCumSum100(t *testing.B)    { benchCumSum(CumSum, 100, t) }
func BenchmarkCumSum1000(t *testing.B)   { benchCumSum(CumSum, 1000, t) }
func BenchmarkCumSum10000(t *testing.B)  { benchCumSum(CumSum, 10000, t) }
func BenchmarkCumSum100000(t *testing.B) { benchCumSum(CumSum, 100000, t) }
func BenchmarkCumSum500000(t *testing.B) { benchCumSum(CumSum, 500000, t) }

func BenchmarkLCumSum1(t *testing.B)      { benchCumSum(naiveCumSum, 1, t) }
func BenchmarkLCumSum2(t *testing.B)      { benchCumSum(naiveCumSum, 2, t) }
func BenchmarkLCumSum3(t *testing.B)      { benchCumSum(naiveCumSum, 3, t) }
func BenchmarkLCumSum4(t *testing.B)      { benchCumSum(naiveCumSum, 4, t) }
func BenchmarkLCumSum5(t *testing.B)      { benchCumSum(naiveCumSum, 5, t) }
func BenchmarkLCumSum10(t *testing.B)     { benchCumSum(naiveCumSum, 10, t) }
func BenchmarkLCumSum100(t *testing.B)    { benchCumSum(naiveCumSum, 100, t) }
func BenchmarkLCumSum1000(t *testing.B)   { benchCumSum(naiveCumSum, 1000, t) }
func BenchmarkLCumSum10000(t *testing.B)  { benchCumSum(naiveCumSum, 10000, t) }
func BenchmarkLCumSum100000(t *testing.B) { benchCumSum(naiveCumSum, 100000, t) }
func BenchmarkLCumSum500000(t *testing.B) { benchCumSum(naiveCumSum, 500000, t) }

func benchCumProd(f func(a, b []float64) []float64, sz int, t *testing.B) {
	a, b := x[:sz], y[:sz]
	for i := 0; i < t.N; i++ {
		f(a, b)
	}
}

var naiveCumProd = func(dst, s []float64) []float64 {
	if len(s) == 0 {
		return dst
	}
	dst[0] = s[0]
	for i, v := range s[1:] {
		dst[i+1] = dst[i] + v
	}
	return dst
}

func BenchmarkCumProd1(t *testing.B)      { benchCumProd(CumProd, 1, t) }
func BenchmarkCumProd2(t *testing.B)      { benchCumProd(CumProd, 2, t) }
func BenchmarkCumProd3(t *testing.B)      { benchCumProd(CumProd, 3, t) }
func BenchmarkCumProd4(t *testing.B)      { benchCumProd(CumProd, 4, t) }
func BenchmarkCumProd5(t *testing.B)      { benchCumProd(CumProd, 5, t) }
func BenchmarkCumProd10(t *testing.B)     { benchCumProd(CumProd, 10, t) }
func BenchmarkCumProd100(t *testing.B)    { benchCumProd(CumProd, 100, t) }
func BenchmarkCumProd1000(t *testing.B)   { benchCumProd(CumProd, 1000, t) }
func BenchmarkCumProd10000(t *testing.B)  { benchCumProd(CumProd, 10000, t) }
func BenchmarkCumProd100000(t *testing.B) { benchCumProd(CumProd, 100000, t) }
func BenchmarkCumProd500000(t *testing.B) { benchCumProd(CumProd, 500000, t) }

func BenchmarkLCumProd1(t *testing.B)      { benchCumProd(naiveCumProd, 1, t) }
func BenchmarkLCumProd2(t *testing.B)      { benchCumProd(naiveCumProd, 2, t) }
func BenchmarkLCumProd3(t *testing.B)      { benchCumProd(naiveCumProd, 3, t) }
func BenchmarkLCumProd4(t *testing.B)      { benchCumProd(naiveCumProd, 4, t) }
func BenchmarkLCumProd5(t *testing.B)      { benchCumProd(naiveCumProd, 5, t) }
func BenchmarkLCumProd10(t *testing.B)     { benchCumProd(naiveCumProd, 10, t) }
func BenchmarkLCumProd100(t *testing.B)    { benchCumProd(naiveCumProd, 100, t) }
func BenchmarkLCumProd1000(t *testing.B)   { benchCumProd(naiveCumProd, 1000, t) }
func BenchmarkLCumProd10000(t *testing.B)  { benchCumProd(naiveCumProd, 10000, t) }
func BenchmarkLCumProd100000(t *testing.B) { benchCumProd(naiveCumProd, 100000, t) }
func BenchmarkLCumProd500000(t *testing.B) { benchCumProd(naiveCumProd, 500000, t) }

func benchDiv(f func(a, b []float64), sz int, t *testing.B) {
	a, b := x[:sz], y[:sz]
	for i := 0; i < t.N; i++ {
		f(a, b)
	}
}

var naiveDiv = func(a, b []float64) {
	for i, v := range b {
		a[i] /= v
	}
}

func BenchmarkDiv1(t *testing.B)      { benchDiv(Div, 1, t) }
func BenchmarkDiv2(t *testing.B)      { benchDiv(Div, 2, t) }
func BenchmarkDiv3(t *testing.B)      { benchDiv(Div, 3, t) }
func BenchmarkDiv4(t *testing.B)      { benchDiv(Div, 4, t) }
func BenchmarkDiv5(t *testing.B)      { benchDiv(Div, 5, t) }
func BenchmarkDiv10(t *testing.B)     { benchDiv(Div, 10, t) }
func BenchmarkDiv100(t *testing.B)    { benchDiv(Div, 100, t) }
func BenchmarkDiv1000(t *testing.B)   { benchDiv(Div, 1000, t) }
func BenchmarkDiv10000(t *testing.B)  { benchDiv(Div, 10000, t) }
func BenchmarkDiv100000(t *testing.B) { benchDiv(Div, 100000, t) }
func BenchmarkDiv500000(t *testing.B) { benchDiv(Div, 500000, t) }

func BenchmarkLDiv1(t *testing.B)      { benchDiv(naiveDiv, 1, t) }
func BenchmarkLDiv2(t *testing.B)      { benchDiv(naiveDiv, 2, t) }
func BenchmarkLDiv3(t *testing.B)      { benchDiv(naiveDiv, 3, t) }
func BenchmarkLDiv4(t *testing.B)      { benchDiv(naiveDiv, 4, t) }
func BenchmarkLDiv5(t *testing.B)      { benchDiv(naiveDiv, 5, t) }
func BenchmarkLDiv10(t *testing.B)     { benchDiv(naiveDiv, 10, t) }
func BenchmarkLDiv100(t *testing.B)    { benchDiv(naiveDiv, 100, t) }
func BenchmarkLDiv1000(t *testing.B)   { benchDiv(naiveDiv, 1000, t) }
func BenchmarkLDiv10000(t *testing.B)  { benchDiv(naiveDiv, 10000, t) }
func BenchmarkLDiv100000(t *testing.B) { benchDiv(naiveDiv, 100000, t) }
func BenchmarkLDiv500000(t *testing.B) { benchDiv(naiveDiv, 500000, t) }

func benchDivTo(f func(dst, a, b []float64) []float64, sz int, t *testing.B) {
	dst, a, b := z[:sz], x[:sz], y[:sz]
	for i := 0; i < t.N; i++ {
		f(dst, a, b)
	}
}

var naiveDivTo = func(dst, s, t []float64) []float64 {
	for i, v := range s {
		dst[i] = v / t[i]
	}
	return dst
}

func BenchmarkDivTo1(t *testing.B)      { benchDivTo(DivTo, 1, t) }
func BenchmarkDivTo2(t *testing.B)      { benchDivTo(DivTo, 2, t) }
func BenchmarkDivTo3(t *testing.B)      { benchDivTo(DivTo, 3, t) }
func BenchmarkDivTo4(t *testing.B)      { benchDivTo(DivTo, 4, t) }
func BenchmarkDivTo5(t *testing.B)      { benchDivTo(DivTo, 5, t) }
func BenchmarkDivTo10(t *testing.B)     { benchDivTo(DivTo, 10, t) }
func BenchmarkDivTo100(t *testing.B)    { benchDivTo(DivTo, 100, t) }
func BenchmarkDivTo1000(t *testing.B)   { benchDivTo(DivTo, 1000, t) }
func BenchmarkDivTo10000(t *testing.B)  { benchDivTo(DivTo, 10000, t) }
func BenchmarkDivTo100000(t *testing.B) { benchDivTo(DivTo, 100000, t) }
func BenchmarkDivTo500000(t *testing.B) { benchDivTo(DivTo, 500000, t) }

func BenchmarkLDivTo1(t *testing.B)      { benchDivTo(naiveDivTo, 1, t) }
func BenchmarkLDivTo2(t *testing.B)      { benchDivTo(naiveDivTo, 2, t) }
func BenchmarkLDivTo3(t *testing.B)      { benchDivTo(naiveDivTo, 3, t) }
func BenchmarkLDivTo4(t *testing.B)      { benchDivTo(naiveDivTo, 4, t) }
func BenchmarkLDivTo5(t *testing.B)      { benchDivTo(naiveDivTo, 5, t) }
func BenchmarkLDivTo10(t *testing.B)     { benchDivTo(naiveDivTo, 10, t) }
func BenchmarkLDivTo100(t *testing.B)    { benchDivTo(naiveDivTo, 100, t) }
func BenchmarkLDivTo1000(t *testing.B)   { benchDivTo(naiveDivTo, 1000, t) }
func BenchmarkLDivTo10000(t *testing.B)  { benchDivTo(naiveDivTo, 10000, t) }
func BenchmarkLDivTo100000(t *testing.B) { benchDivTo(naiveDivTo, 100000, t) }
func BenchmarkLDivTo500000(t *testing.B) { benchDivTo(naiveDivTo, 500000, t) }

func benchL1Norm(f func(a, b []float64) float64, sz int, t *testing.B) {
	a, b := x[:sz], y[:sz]
	for i := 0; i < t.N; i++ {
		f(a, b)
	}
}

var naiveL1Norm = func(s, t []float64) float64 {
	var norm float64
	for i, v := range s {
		norm += math.Abs(t[i] - v)
	}
	return norm
}

func BenchmarkL1Norm1(t *testing.B)      { benchL1Norm(L1Norm, 1, t) }
func BenchmarkL1Norm2(t *testing.B)      { benchL1Norm(L1Norm, 2, t) }
func BenchmarkL1Norm3(t *testing.B)      { benchL1Norm(L1Norm, 3, t) }
func BenchmarkL1Norm4(t *testing.B)      { benchL1Norm(L1Norm, 4, t) }
func BenchmarkL1Norm5(t *testing.B)      { benchL1Norm(L1Norm, 5, t) }
func BenchmarkL1Norm10(t *testing.B)     { benchL1Norm(L1Norm, 10, t) }
func BenchmarkL1Norm100(t *testing.B)    { benchL1Norm(L1Norm, 100, t) }
func BenchmarkL1Norm1000(t *testing.B)   { benchL1Norm(L1Norm, 1000, t) }
func BenchmarkL1Norm10000(t *testing.B)  { benchL1Norm(L1Norm, 10000, t) }
func BenchmarkL1Norm100000(t *testing.B) { benchL1Norm(L1Norm, 100000, t) }
func BenchmarkL1Norm500000(t *testing.B) { benchL1Norm(L1Norm, 500000, t) }

func BenchmarkLL1Norm1(t *testing.B)      { benchL1Norm(naiveL1Norm, 1, t) }
func BenchmarkLL1Norm2(t *testing.B)      { benchL1Norm(naiveL1Norm, 2, t) }
func BenchmarkLL1Norm3(t *testing.B)      { benchL1Norm(naiveL1Norm, 3, t) }
func BenchmarkLL1Norm4(t *testing.B)      { benchL1Norm(naiveL1Norm, 4, t) }
func BenchmarkLL1Norm5(t *testing.B)      { benchL1Norm(naiveL1Norm, 5, t) }
func BenchmarkLL1Norm10(t *testing.B)     { benchL1Norm(naiveL1Norm, 10, t) }
func BenchmarkLL1Norm100(t *testing.B)    { benchL1Norm(naiveL1Norm, 100, t) }
func BenchmarkLL1Norm1000(t *testing.B)   { benchL1Norm(naiveL1Norm, 1000, t) }
func BenchmarkLL1Norm10000(t *testing.B)  { benchL1Norm(naiveL1Norm, 10000, t) }
func BenchmarkLL1Norm100000(t *testing.B) { benchL1Norm(naiveL1Norm, 100000, t) }
func BenchmarkLL1Norm500000(t *testing.B) { benchL1Norm(naiveL1Norm, 500000, t) }

func benchLinfNorm(f func(a, b []float64) float64, sz int, t *testing.B) {
	a, b := x[:sz], y[:sz]
	for i := 0; i < t.N; i++ {
		f(a, b)
	}
}

var naiveLinfNorm = func(s, t []float64) float64 {
	var norm float64
	if len(s) == 0 {
		return 0
	}
	norm = math.Abs(t[0] - s[0])
	for i, v := range s[1:] {
		absDiff := math.Abs(t[i+1] - v)
		if absDiff > norm || math.IsNaN(norm) {
			norm = absDiff
		}
	}
	return norm
}

func BenchmarkLinfNorm1(t *testing.B)      { benchLinfNorm(LinfNorm, 1, t) }
func BenchmarkLinfNorm2(t *testing.B)      { benchLinfNorm(LinfNorm, 2, t) }
func BenchmarkLinfNorm3(t *testing.B)      { benchLinfNorm(LinfNorm, 3, t) }
func BenchmarkLinfNorm4(t *testing.B)      { benchLinfNorm(LinfNorm, 4, t) }
func BenchmarkLinfNorm5(t *testing.B)      { benchLinfNorm(LinfNorm, 5, t) }
func BenchmarkLinfNorm10(t *testing.B)     { benchLinfNorm(LinfNorm, 10, t) }
func BenchmarkLinfNorm100(t *testing.B)    { benchLinfNorm(LinfNorm, 100, t) }
func BenchmarkLinfNorm1000(t *testing.B)   { benchLinfNorm(LinfNorm, 1000, t) }
func BenchmarkLinfNorm10000(t *testing.B)  { benchLinfNorm(LinfNorm, 10000, t) }
func BenchmarkLinfNorm100000(t *testing.B) { benchLinfNorm(LinfNorm, 100000, t) }
func BenchmarkLinfNorm500000(t *testing.B) { benchLinfNorm(LinfNorm, 500000, t) }

func BenchmarkLLinfNorm1(t *testing.B)      { benchLinfNorm(naiveLinfNorm, 1, t) }
func BenchmarkLLinfNorm2(t *testing.B)      { benchLinfNorm(naiveLinfNorm, 2, t) }
func BenchmarkLLinfNorm3(t *testing.B)      { benchLinfNorm(naiveLinfNorm, 3, t) }
func BenchmarkLLinfNorm4(t *testing.B)      { benchLinfNorm(naiveLinfNorm, 4, t) }
func BenchmarkLLinfNorm5(t *testing.B)      { benchLinfNorm(naiveLinfNorm, 5, t) }
func BenchmarkLLinfNorm10(t *testing.B)     { benchLinfNorm(naiveLinfNorm, 10, t) }
func BenchmarkLLinfNorm100(t *testing.B)    { benchLinfNorm(naiveLinfNorm, 100, t) }
func BenchmarkLLinfNorm1000(t *testing.B)   { benchLinfNorm(naiveLinfNorm, 1000, t) }
func BenchmarkLLinfNorm10000(t *testing.B)  { benchLinfNorm(naiveLinfNorm, 10000, t) }
func BenchmarkLLinfNorm100000(t *testing.B) { benchLinfNorm(naiveLinfNorm, 100000, t) }
func BenchmarkLLinfNorm500000(t *testing.B) { benchLinfNorm(naiveLinfNorm, 500000, t) }
