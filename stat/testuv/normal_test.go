package testuv

import (
	"math"
	"testing"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

func getNormalValues() []float64 {
	normal := distuv.Normal{Mu: 0, Sigma: 1, Src: rand.New(rand.NewSource(1))}
	size := 100
	nums := make([]float64, size)
	for i := range nums {
		nums[i] = normal.Rand()
	}
	return nums
}

func TestNormalSkewTest(t *testing.T) {
	// Even though the skew test and kurtosis test are a copy from the scipy's the results
	// are not exactly the same because the skew/kurtosis corrections are not exactly the same
	// but the statistic value is close enough.
	for i, test := range []struct {
		values   []float64
		expected float64
	}{
		{
			values:   []float64{0, 1, 2, 3, 4, 5, 6, 7, 8},
			expected: 1.018464355396213,
		},
		{
			values:   []float64{2, 8, 0, 4, 1, 9, 9, 0},
			expected: 0.5562585562766172,
		},
		{
			values:   []float64{1, 2, 3, 4, 5, 6, 7, 8000},
			expected: 4.319816401673864,
		},
		{
			values:   []float64{100, 100, 100, 100, 100, 100, 100, 101},
			expected: 4.319820025201098,
		},
	} {
		z := NormalSkewTest(test.values)
		if math.Abs(z-test.expected) > 1e-7 {
			t.Errorf("NormalSkewTest mismatch case %d. Expected %v, Found %v", i, test.expected, z)
		}
	}
}

func TestNormalKurtosisTest(t *testing.T) {
	for i, test := range []struct {
		values   []float64
		expected float64
	}{
		{
			values:   []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
			expected: -1.6867033202073243,
		},
		{
			values:   getNormalValues(),
			expected: -1.437284826863465,
		},
	} {
		z := NormalKurtosisTest(test.values)
		if math.Abs(z-test.expected) > 1e-7 {
			t.Errorf("NormalKurtosisTest mismatch case %d. Expected %v, Found %v", i, test.expected, z)
		}
	}
}

func TestNormalTest(t *testing.T) {
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
		values   []float64
		expected float64
	}{
		{
			values:   []float64{0, 1, 2, 3, 4, 5, 6, 7, 8},
			expected: 1.7509245653153176,
		},
		{
			values:   []float64{2, 8, 0, 4, 1, 9, 9, 0},
			expected: 11.454757293481551,
		},
		{
			values:   []float64{1, 2, 3, 4, 5, 6, 7, 8000},
			expected: 40.53534243515444,
		},
		{
			values:   []float64{100, 100, 100, 100, 100, 100, 100, 101},
			expected: 40.53539760601764,
		},
		{
			values:   []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19},
			expected: 3.9272951079276743,
		},
		{
			values:   getNormalValues(),
			expected: 2.083696972397116,
		},
	} {
		z := NormalTest(test.values)
		if math.Abs(z-test.expected) > 1e-7 {
			t.Errorf("NormalTest mismatch case %d. Expected %v, Found %v", i, test.expected, z)
		}
	}
}
