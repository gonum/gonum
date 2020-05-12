// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package test

import (
	"math"
	"sort"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat/distuv"
)

const badLengthMismatch = "test: slice length mismatch"
const badLengthIsZero = "test: slice length is zero"
const badNumberTooLarge = "test: number too large"

type pairStruct struct {
	Position int
	Value    float64
}

func ensureDataConformance(x, y []float64) {
	if len(x) == 0 || len(y) == 0 {
		panic(badLengthIsZero)
	}
	if len(x) != len(y) {
		panic(badLengthMismatch)
	}
}

func calculateAbsDifference(dst, x, y []float64) {
	for index, xVal := range x {
		dst[index] = math.Abs(y[index] - xVal)
	}
}

func rank(out, in []float64) int {
	tieAdjust := 0
	lengthOfIn := len(in)
	pairStructArr := make([]pairStruct, lengthOfIn)
	for index, inVal := range in {
		pairStructArr[index] = pairStruct{
			Position: index,
			Value:    inVal,
		}
	}
	sort.Slice(pairStructArr, func(i, j int) bool {
		return pairStructArr[i].Value < pairStructArr[j].Value
	})

	mpOut := make(map[float64]int, len(out))   // tracks the rank
	mpCount := make(map[float64]int, len(out)) // tracks the count
	tmpRank := make(map[float64]struct{})      // tracks the tied ranks
	rank := 1
	for i := 0; i < lengthOfIn; i++ {
		val := pairStructArr[i].Value
		if val != 0 {
			if _, ok := mpOut[val]; !ok {
				mpOut[val] = rank
			} else {
				mpOut[val] += rank
			}
			if _, ok := mpCount[val]; !ok {
				mpCount[val] = 1
			} else {
				mpCount[val] += 1
			}
			rank++
		}
	}

	var pos int
	var val float64
	for i := 0; i < lengthOfIn; i++ {
		pos = pairStructArr[i].Position
		val = pairStructArr[i].Value
		if val != 0 {
			out[pos] = float64(mpOut[val]) / float64(mpCount[val])
		}
		if mpCount[val] > 1 {
			if _, ok := tmpRank[val]; !ok {
				tmpRank[val] = struct{}{}
			}
		}
	}
	for val := range tmpRank {
		tieAdjust += (mpCount[val] * mpCount[val] * mpCount[val]) - (mpCount[val])

	}
	return tieAdjust
}

func wilCoxonSignedRankTest(x, y []float64) (float64, int, int) {
	z := make([]float64, len(x))
	absZ := make([]float64, len(x))
	floats.SubTo(z, x, y)
	calculateAbsDifference(absZ, x, y)

	ranks := make([]float64, len(x))
	tieAdj := rank(ranks, absZ)
	tmpLenOfX := len(x)

	WPlus := 0.0
	WMinus := 0.0
	for index, rank := range ranks {
		if z[index] > 0 {
			WPlus += rank
		} else if z[index] == 0 {
			tmpLenOfX--
		}
	}
	WMinus = (float64(tmpLenOfX*(tmpLenOfX+1)) / 2.0) - WPlus
	return math.Max(WMinus, WPlus), tmpLenOfX, tieAdj
}

func calculateExactPValue(Wmax float64, NZ int) float64 {
	m := 1 << NZ
	largerRankSums := 0

	for i := 0; i < m; i++ {
		rankSum := 0
		// Generate all possible rank sums
		for j := 0; j < NZ; j++ {
			// (i >> j) & 1 extract i's j-th bit from the right
			if ((i >> j) & 1) == 1 {
				rankSum += j + 1
			}
		}
		if float64(rankSum) >= Wmax {
			largerRankSums++
		}
	}

	// largerRankSums / m gives the one-sided p-value, so it's multiplied
	// with 2 to get the two-sided p-value
	return 2 * (float64(largerRankSums) / float64(m))
}

func calculateAsymptoticPValue(Wmin float64, NZ int, tieAdj int) float64 {
	// n should be number of non-zero absolute difference pairs
	ES := float64(NZ*(NZ+1)) / 4.0
	VarS := (ES * (float64(2*NZ + 1)) / 6.0) - (float64(tieAdj) / 48.0)

	// - 0.5 is a continuity correction
	z := (Wmin - ES - 0.5) / math.Sqrt(VarS)

	standardNormal := distuv.UnitNormal
	return 2 * standardNormal.CDF(z)
}

/*
	WilcoxonSignedRankTest implements Wilcoxon signed rank test (https://en.wikipedia.org/wiki/Wilcoxon_signed-rank_test)

	The Wilcoxon signed-rank test tests the null hypothesis that two related paired samples come from the same distribution. In particular, it tests whether the distribution of the differences x - y is symmetric about zero. It is a non-parametric version of the paired T-test.

	Parameters:
		x:
			The first set of measurements.
		y:
			The second set of measurements.
		exactPValue:
			Exact P value computation is expensive, it is only available for the measurements with dimension less than 30.
			Exact P Value as true for measurements greater than 30 panics.

    Returns:
		The two sided p-value for the test
*/
func WilcoxonSignedRankTest(x []float64, y []float64, exactPValue bool) float64 {
	ensureDataConformance(x, y)
	Wmax, N, tieAdj := wilCoxonSignedRankTest(x, y)
	if exactPValue && N > 30 {
		panic(badNumberTooLarge)
	}
	if exactPValue {
		return calculateExactPValue(Wmax, N)
	} else {
		Wmin := (float64(N*(N+1)) / 2.0) - Wmax
		return calculateAsymptoticPValue(Wmin, N, tieAdj)
	}
}
