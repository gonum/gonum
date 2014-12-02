package testblas

import (
	"math"
	"testing"

	"github.com/gonum/blas"
)

type Dtbmver interface {
	Dtbmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n, k int, a []float64, lda int, x []float64, incX int)
}

func DtbmvTest(t *testing.T, blasser Dtbmver) {
	for i, test := range []struct {
		ul  blas.Uplo
		tA  blas.Transpose
		d   blas.Diag
		n   int
		k   int
		a   [][]float64
		x   []float64
		ans []float64
	}{
		{
			ul: blas.Upper,
			tA: blas.NoTrans,
			d:  blas.Unit,
			n:  3,
			k:  1,
			a: [][]float64{
				{1, 2, 0},
				{0, 1, 4},
				{0, 0, 1},
			},
			x:   []float64{2, 3, 4},
			ans: []float64{8, 19, 4},
		},
		{
			ul: blas.Upper,
			tA: blas.NoTrans,
			d:  blas.NonUnit,
			n:  5,
			k:  1,
			a: [][]float64{
				{1, 3, 0, 0, 0},
				{0, 6, 7, 0, 0},
				{0, 0, 2, 1, 0},
				{0, 0, 0, 12, 3},
				{0, 0, 0, 0, -1},
			},
			x:   []float64{1, 2, 3, 4, 5},
			ans: []float64{7, 33, 10, 63, -5},
		},
		{
			ul: blas.Lower,
			tA: blas.NoTrans,
			d:  blas.NonUnit,
			n:  5,
			k:  1,
			a: [][]float64{
				{7, 0, 0, 0, 0},
				{3, 6, 0, 0, 0},
				{0, 7, 2, 0, 0},
				{0, 0, 1, 12, 0},
				{0, 0, 0, 3, -1},
			},
			x:   []float64{1, 2, 3, 4, 5},
			ans: []float64{7, 15, 20, 51, 7},
		},
		{
			ul: blas.Upper,
			tA: blas.Trans,
			d:  blas.NonUnit,
			n:  5,
			k:  2,
			a: [][]float64{
				{7, 3, 9, 0, 0},
				{0, 6, 7, 10, 0},
				{0, 0, 2, 1, 11},
				{0, 0, 0, 12, 3},
				{0, 0, 0, 0, -1},
			},
			x:   []float64{1, 2, 3, 4, 5},
			ans: []float64{7, 15, 29, 71, 40},
		},
		{
			ul: blas.Lower,
			tA: blas.Trans,
			d:  blas.NonUnit,
			n:  5,
			k:  2,
			a: [][]float64{
				{7, 0, 0, 0, 0},
				{3, 6, 0, 0, 0},
				{9, 7, 2, 0, 0},
				{0, 10, 1, 12, 0},
				{0, 0, 11, 3, -1},
			},
			x:   []float64{1, 2, 3, 4, 5},
			ans: []float64{40, 73, 65, 63, -5},
		},
	} {
		extra := 0
		var aFlat []float64
		if test.ul == blas.Upper {
			aFlat = flattenBanded(test.a, test.k, 0)
		} else {
			aFlat = flattenBanded(test.a, 0, test.k)
		}
		incTest := func(incX, extra int) {
			xnew := makeIncremented(test.x, incX, extra)
			ans := makeIncremented(test.ans, incX, extra)
			lda := test.k + 1
			blasser.Dtbmv(test.ul, test.tA, test.d, test.n, test.k, aFlat, lda, xnew, incX)
			if !dSliceTolEqual(ans, xnew) {
				t.Errorf("Case %v, Inc %v: Want %v, got %v", i, incX, ans, xnew)
			}
		}
		incTest(1, extra)
		incTest(3, extra)
		incTest(-2, extra)
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
