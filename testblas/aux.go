package testblas

import (
	"math"
	"testing"
)

// throwPanic will throw unexpected panics if true, or will just report them as errors if false
const throwPanic = true

func dTolEqual(a, b float64) bool {
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

func dSliceEqual(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !(a[i] == b[i]) {
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
	for i, val := range x {
		xTmp[i] = val
	}
	for i, val := range y {
		yTmp[i] = val
	}
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
