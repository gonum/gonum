package wilcoxon

import (
	"gonum.org/v1/gonum/stat/distuv"
	"math"
	"sort"
)

func ensureDataConformance(x []float64, y []float64) {
	if len(x) == 0 {
		panic("x array length is zero")
	}
	if len(y) == 0 {
		panic("y array length is zero")
	}
	if len(x) != len(y) {
		panic("dimension mismatch in ensure data conformance")
	}
}

func CalculateDifferences(x []float64, y []float64, z []float64) {
	ensureDataConformance(x, y)
	ensureDataConformance(x, z)
	for i := 0; i < len(x); i++ {
		z[i] = x[i] - y[i]
	}
}

func CalculateAbsDifference(x []float64, y []float64, z []float64) {
	ensureDataConformance(x, y)
	ensureDataConformance(x, z)
	for i := 0; i < len(x); i++ {
		z[i] = math.Abs(y[i] - x[i])
	}
}

func Rank(in []float64, out []float64) ([]float64, int) {
	ensureDataConformance(in, out)

	var tieAdjust = 0
	type pairStruct struct {
		Position int
		Value    float64
	}
	var lengthOfIn = len(in)
	var pairStructArr = make([]pairStruct, lengthOfIn)
	for i := 0; i < lengthOfIn; i++ {
		pairStructArr[i] = pairStruct{
			Position: i,
			Value:    in[i],
		}
	}
	sort.Slice(pairStructArr, func(i, j int) bool {
		return pairStructArr[i].Value < pairStructArr[j].Value
	})

	var mpOut = make(map[float64]int, len(out))   // tracks the rank
	var mpCount = make(map[float64]int, len(out)) // tracks the count
	var rank = 1
	for i := 0; i < lengthOfIn; i++ {
		var val = pairStructArr[i].Value
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

	var tmpRank = make(map[float64]struct{})
	for i := 0; i < lengthOfIn; i++ {
		var pos = pairStructArr[i].Position
		var val = pairStructArr[i].Value
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
	return out, tieAdjust
}

func WilCoxonSignedRankTest(x []float64, y []float64) (float64, int, int) {
	ensureDataConformance(x, y)

	var z = make([]float64, len(x))
	var absZ = make([]float64, len(x))
	CalculateDifferences(x, y, z)
	CalculateAbsDifference(x, y, absZ)

	var ranks = make([]float64, len(x))
	var tieAdj = 0
	ranks, tieAdj = Rank(absZ, ranks)
	var tmpLenOfX = len(x)

	var WPlus = 0.0
	var WMinus = 0.0
	for i := 0; i < len(x); i++ {
		if z[i] > 0 {
			WPlus += ranks[i]
		} else if z[i] == 0 {
			tmpLenOfX--
		}
	}
	WMinus = (float64(tmpLenOfX*(tmpLenOfX+1)) / 2.0) - WPlus
	return math.Max(WMinus, WPlus), tmpLenOfX, tieAdj
}

func CalculateExactPValue(Wmax float64, N int) float64 {
	var m = 1 << N

	largerRankSums := 0

	for i := 0; i < m; i++ {
		rankSum := 0

		// Generate all possible rank sums
		for j := 0; j < N; j++ {

			// (i >> j) & 1 extract i's j-th bit from the right
			if ((i >> j) & 1) == 1 {
				rankSum += j + 1
			}
		}

		if float64(rankSum) >= Wmax {
			largerRankSums++
		}
	}

	/*
	 * largerRankSums / m gives the one-sided p-value, so it's multiplied
	 * with 2 to get the two-sided p-value
	 */
	return 2 * (float64(largerRankSums) / float64(m))
}

func CalculateAsymptoticPValue(Wmin float64, NZ int, tieAdj int) float64 {

	// n should be number of non zeros

	ES := float64(NZ*(NZ+1)) / 4.0

	VarS := ES * (float64(2*NZ + 1)) / 6.0

	VarS -= float64(tieAdj) / 48.0

	// - 0.5 is a continuity correction
	z := (Wmin - ES - 0.5) / math.Sqrt(VarS)

	standardNormal := distuv.UnitNormal

	return 2 * standardNormal.CDF(z)
}

func WilcoxonSignedRankTest(x []float64, y []float64, exactPValue bool) float64 {
	ensureDataConformance(x, y)
	Wmax, N, tieAdj := WilCoxonSignedRankTest(x, y)
	if exactPValue && N > 30 {
		panic("number too large")
	}
	if exactPValue {
		return CalculateExactPValue(Wmax, N)
	} else {
		Wmin := (float64(N*(N+1)) / 2.0) - Wmax
		return CalculateAsymptoticPValue(Wmin, N, tieAdj)
	}
}
