// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testuv

import (
	"math"
	"testing"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat/distuv"
)

const tol = 1e-7

func rnorm(n int) []float64 {
	normal := distuv.Normal{Mu: 0, Sigma: 1, Src: rand.New(rand.NewSource(1))}
	nums := make([]float64, n)
	for i := range nums {
		nums[i] = normal.Rand()
	}
	return nums
}

func TestNormalSkew(t *testing.T) {
	// Even though the skew test and kurtosis test are a copy from the scipy's the results
	// are not exactly the same because the skew/kurtosis corrections are not exactly the same
	// but the statistic value is close enough.
	for i, test := range []struct {
		values []float64
		want   float64
	}{
		{
			values: []float64{0, 1, 2, 3, 4, 5, 6, 7, 8},
			want:   1.018464355396213,
		},
		{
			values: []float64{2, 8, 0, 4, 1, 9, 9, 0},
			want:   0.5562585562766172,
		},
		{
			values: []float64{1, 2, 3, 4, 5, 6, 7, 8000},
			want:   4.319816401673864,
		},
		{
			values: []float64{100, 100, 100, 100, 100, 100, 100, 101},
			want:   4.319820025201098,
		},
	} {
		z := NormalSkew(test.values)
		if floats.EqualWithinAbsOrRel(z, test.want, tol, tol) {
			t.Errorf("NormalSkew mismatch case %d. Expected %v, Found %v", i, test.want, z)
		}
	}
}

func TestNormalKurtosis(t *testing.T) {
	for i, test := range []struct {
		values []float64
		want   float64
	}{
		{
			values: []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
			want:   -1.6867033202073243,
		},
		{
			values: rnorm(100),
			want:   -1.437284826863465,
		},
	} {
		z := NormalKurtosis(test.values)
		if math.Abs(z-test.want) > 1e-7 {
			t.Errorf("NormalKurtosis mismatch case %d. Expected %v, Found %v", i, test.want, z)
		}
	}
}

func TestNormal(t *testing.T) {
	// normal test with scipy yield similar results:
	// ```
	// from scipy.stats import normaltest
	// print(normaltest(list(range(9))).statistic)
	// print(normaltest([2, 8, 0, 4, 1, 9, 9, 0]).statistic)
	// print(normaltest([1, 2, 3, 4, 5, 6, 7, 8000]).statistic)
	// print(normaltest([100, 100, 100, 100, 100, 100, 100, 101]).statistic)
	// print(normaltest(list(range(20))).statistic)
	// ```

	for i, test := range []struct {
		values []float64
		want   float64
	}{
		{
			values: []float64{0, 1, 2, 3, 4, 5, 6, 7, 8},
			want:   1.7509245653153176,
		},
		{
			values: []float64{2, 8, 0, 4, 1, 9, 9, 0},
			want:   11.454757293481551,
		},
		{
			values: []float64{1, 2, 3, 4, 5, 6, 7, 8000},
			want:   40.53534243515444,
		},
		{
			values: []float64{100, 100, 100, 100, 100, 100, 100, 101},
			want:   40.53539760601764,
		},
		{
			values: []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
			want:   3.9272951079276743,
		},
		{
			values: rnorm(100),
			want:   2.083696972397116,
		},
	} {
		z := Normal(test.values)
		if floats.EqualWithinAbsOrRel(z, test.want, tol, tol) {
			t.Errorf("Normal mismatch case %d. Expected %v, Found %v", i, test.want, z)
		}
	}
}
