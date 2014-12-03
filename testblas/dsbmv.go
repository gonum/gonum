package testblas

import (
	"math"
	"testing"

	"github.com/gonum/blas"
)

type Dsbmver interface {
	Dsbmv(ul blas.Uplo, n, k int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int)
}

func DsbmvTest(t *testing.T, blasser Dsbmver) {
	for i, test := range []struct {
		ul    blas.Uplo
		n     int
		k     int
		alpha float64
		beta  float64
		a     [][]float64
		x     []float64
		y     []float64

		ans []float64
	}{
		{
			ul:    blas.Upper,
			n:     4,
			k:     2,
			alpha: 2,
			beta:  3,
			a: [][]float64{
				{7, 8, 2, 0},
				{0, 8, 2, -3},
				{0, 0, 3, 6},
				{0, 0, 0, 9},
			},
			x:   []float64{1, 2, 3, 4},
			y:   []float64{-1, -2, -3, -4},
			ans: []float64{55, 30, 69, 84},
		},
		{
			ul:    blas.Lower,
			n:     4,
			k:     2,
			alpha: 2,
			beta:  3,
			a: [][]float64{
				{7, 0, 0, 0},
				{8, 8, 0, 0},
				{2, 2, 3, 0},
				{0, -3, 6, 9},
			},
			x:   []float64{1, 2, 3, 4},
			y:   []float64{-1, -2, -3, -4},
			ans: []float64{55, 30, 69, 84},
		},
	} {
		extra := 0
		var aFlat []float64
		if test.ul == blas.Upper {
			aFlat = flattenBanded(test.a, test.k, 0)
		} else {
			aFlat = flattenBanded(test.a, 0, test.k)
		}
		incTest := func(incX, incY, extra int) {
			xnew := makeIncremented(test.x, incX, extra)
			ynew := makeIncremented(test.y, incY, extra)
			ans := makeIncremented(test.ans, incY, extra)
			blasser.Dsbmv(test.ul, test.n, test.k, test.alpha, aFlat, test.k+1, xnew, incX, test.beta, ynew, incY)
			if !dSliceTolEqual(ans, ynew) {
				t.Errorf("Case %v: Want %v, got %v", i, ans, ynew)
			}
		}
		incTest(1, 1, extra)
		incTest(1, 3, extra)
		incTest(1, -3, extra)
		incTest(2, 3, extra)
		incTest(2, -3, extra)
		incTest(3, 2, extra)
		incTest(-3, 2, extra)
	}
}

// flattenBanded turns a dense banded slice of slice into the compact banded matrix format
func flattenBanded(a [][]float64, ku, kl int) []float64 {
	m := len(a)
	n := len(a[0])
	if ku < 0 || kl < 0 {
		panic("testblas: negative band length")
	}
	// banded size is minimum of m and n because otherwise just have a bunch of zeros
	nRows := m
	if m < n {
		nRows = n
	}
	nCols := (ku + kl + 1)
	aflat := make([]float64, nRows*nCols)
	for i := range aflat {
		aflat[i] = math.NaN()
	}
	// loop over the rows, and then the bands
	// elements in the ith row stay in the ith row
	// order in bands is kept
	for i := 0; i < nRows; i++ {
		min := -kl
		if i-kl < 0 {
			min = -i
		}
		max := ku
		if i+ku >= n {
			max = n - i - 1
		}
		for j := min; j <= max; j++ {
			col := kl + j
			aflat[i*nCols+col] = a[i][i+j]
		}
	}
	return aflat
}

// makeIncremented takes a slice with inc == 1 and makes an incremented version
// and adds extra values on the end
func makeIncremented(x []float64, inc int, extra int) []float64 {
	if inc == 0 {
		panic("zero inc")
	}
	absinc := inc
	if absinc < 0 {
		absinc = -inc
	}
	xcopy := make([]float64, len(x))
	if inc > 0 {
		copy(xcopy, x)
	} else {
		for i := 0; i < len(x); i++ {
			xcopy[i] = x[len(x)-i-1]
		}
	}
	// don't use NaN because it makes comparison hard
	// Do use a weird unique value for easier debugging
	counter := 100.0
	var xnew []float64
	for i, v := range xcopy {
		xnew = append(xnew, v)
		if i != len(x)-1 {
			for j := 0; j < absinc-1; j++ {
				xnew = append(xnew, counter)
				counter++
			}
		}
	}
	for i := 0; i < extra; i++ {
		xnew = append(xnew, counter)
		counter++
	}
	return xnew
}
