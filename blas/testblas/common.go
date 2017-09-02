// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testblas

import (
	"math"
	"math/cmplx"
	"testing"

	"gonum.org/v1/gonum/blas"
)

// throwPanic will throw unexpected panics if true, or will just report them as errors if false
const throwPanic = true

func dTolEqual(a, b float64) bool {
	if math.IsNaN(a) && math.IsNaN(b) {
		return true
	}
	if a == b {
		return true
	}
	m := math.Max(math.Abs(a), math.Abs(b))
	if m > 1 {
		a /= m
		b /= m
	}
	if math.Abs(a-b) < 1e-14 {
		return true
	}
	return false
}

func dSliceTolEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !dTolEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func dStridedSliceTolEqual(n int, a []float64, inca int, b []float64, incb int) bool {
	ia := 0
	ib := 0
	if inca <= 0 {
		ia = -(n - 1) * inca
	}
	if incb <= 0 {
		ib = -(n - 1) * incb
	}
	for i := 0; i < n; i++ {
		if !dTolEqual(a[ia], b[ib]) {
			return false
		}
		ia += inca
		ib += incb
	}
	return true
}

func dSliceEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !dTolEqual(a[i], b[i]) {
			return false
		}
	}
	return true
}

func dCopyTwoTmp(x, xTmp, y, yTmp []float64) {
	if len(x) != len(xTmp) {
		panic("x size mismatch")
	}
	if len(y) != len(yTmp) {
		panic("y size mismatch")
	}
	copy(xTmp, x)
	copy(yTmp, y)
}

// returns true if the function panics
func panics(f func()) (b bool) {
	defer func() {
		err := recover()
		if err != nil {
			b = true
		}
	}()
	f()
	return
}

func testpanics(f func(), name string, t *testing.T) {
	b := panics(f)
	if !b {
		t.Errorf("%v should panic and does not", name)
	}
}

func sliceOfSliceCopy(a [][]float64) [][]float64 {
	n := make([][]float64, len(a))
	for i := range a {
		n[i] = make([]float64, len(a[i]))
		copy(n[i], a[i])
	}
	return n
}

func sliceCopy(a []float64) []float64 {
	n := make([]float64, len(a))
	copy(n, a)
	return n
}

func flatten(a [][]float64) []float64 {
	if len(a) == 0 {
		return nil
	}
	m := len(a)
	n := len(a[0])
	s := make([]float64, m*n)
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			s[i*n+j] = a[i][j]
		}
	}
	return s
}

func unflatten(a []float64, m, n int) [][]float64 {
	s := make([][]float64, m)
	for i := 0; i < m; i++ {
		s[i] = make([]float64, n)
		for j := 0; j < n; j++ {
			s[i][j] = a[i*n+j]
		}
	}
	return s
}

// flattenTriangular turns the upper or lower triangle of a dense slice of slice
// into a single slice with packed storage. a must be a square matrix.
func flattenTriangular(a [][]float64, ul blas.Uplo) []float64 {
	m := len(a)
	aFlat := make([]float64, m*(m+1)/2)
	var k int
	if ul == blas.Upper {
		for i := 0; i < m; i++ {
			k += copy(aFlat[k:], a[i][i:])
		}
		return aFlat
	}
	for i := 0; i < m; i++ {
		k += copy(aFlat[k:], a[i][:i+1])
	}
	return aFlat
}

// flattenBanded turns a dense banded slice of slice into the compact banded matrix format
func flattenBanded(a [][]float64, ku, kl int) []float64 {
	m := len(a)
	n := len(a[0])
	if ku < 0 || kl < 0 {
		panic("testblas: negative band length")
	}
	nRows := m
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

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func allPairs(x, y []int) [][2]int {
	var p [][2]int
	for _, v0 := range x {
		for _, v1 := range y {
			p = append(p, [2]int{v0, v1})
		}
	}
	return p
}

func zsame(x, y []complex128) bool {
	if len(x) != len(y) {
		return false
	}
	for i, v := range x {
		w := y[i]
		if v != w && !math.IsNaN(real(v)) && !math.IsNaN(imag(v)) && !math.IsNaN(real(w)) && !math.IsNaN(imag(w)) {
			return false
		}
	}
	return true
}

func makeZVector(data []complex128, inc int) []complex128 {
	if inc == 0 {
		panic("bad test")
	}
	if len(data) == 0 {
		return nil
	}
	inc = abs(inc)
	x := make([]complex128, (len(data)-1)*inc+1)
	for i := range x {
		x[i] = cmplx.NaN()
	}
	for i, v := range data {
		x[i*inc] = v
	}
	return x
}

func makeZGeneral(data []complex128, m, n int, ld int) []complex128 {
	if m < 0 || n < 0 {
		panic("bad test")
	}
	if len(data) != m*n {
		panic("bad test")
	}
	if ld < max(1, n) {
		panic("bad test")
	}
	if len(data) == 0 {
		return nil
	}
	a := make([]complex128, (m-1)*ld+n)
	for i := range a {
		a[i] = cmplx.NaN()
	}
	for i := 0; i < m; i++ {
		copy(a[i*ld:i*ld+n], data[i*n:i*n+n])
	}
	return a
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
