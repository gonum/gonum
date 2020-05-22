package interp

import (
	"fmt"
	"math"
	"testing"
)

func panics(fn func()) (panicked bool, message string) {
	defer func() {
		r := recover()
		panicked = r != nil
		message = fmt.Sprint(r)
	}()
	fn()
	return
}

func TestNewConstInterpolator1D(t *testing.T) {
	t.Parallel()
	const value float64 = 42.0
	i1d := NewConstInterpolator1D(value)
	if i1d.begin() != math.Inf(-1) {
		t.Errorf("unexpected begin() value: got: %g want: %g", i1d.begin(), math.Inf(-1))
	}
	if i1d.end() != math.Inf(1) {
		t.Errorf("unexpected end() value: got: %g want: %g", i1d.end(), math.Inf(1))
	}
}

func TestConstInterpolator1DEval(t *testing.T) {
	t.Parallel()
	const value float64 = 42.0
	i1d := NewConstInterpolator1D(value)
	xs := [...]float64{math.Inf(-1), -11, 0.4, 1e9, math.Inf(1)}
	for _, x := range xs {
		y := i1d.eval(x)
		if y != value {
			t.Errorf("unexpected eval(%g) value: got: %g want: %g", x, y, value)
		}
	}
}

func TestFindSegment(t *testing.T) {
	t.Parallel()
	xs := []float64{0, 1, 2}
	type params struct {
		x         float64
		expectedI int
		expectedX float64
	}
	paramSets := [...]params{{0, 0, 0}, {0.3, 0, 0}, {1, 1, 1}, {1.5, 1, 1}, {2, 2, 2}}
	for _, param := range paramSets {
		i, x := findSegment(xs, param.x)
		if i != param.expectedI || x != param.expectedX {
			t.Errorf("unexpected value of findSegment(xs, %g): got %d, %g want: %d, %g", param.x, i, x, param.expectedI, param.expectedX)
		}
	}
	panicXs := [...]float64{-0.5, 2.1}
	expectedMessages := [...]string{
		"interp: x value -0.5 below lower bound 0",
		"interp: x value 2.1 above upper bound 2",
	}
	for i, x := range panicXs {
		panicked, message := panics(func() { findSegment(xs, x) })
		if !panicked || message != expectedMessages[i] {
			t.Errorf("expected panic with message '%s' for evaluating at invalid x: %g", expectedMessages[i], x)
		}
	}
}

func BenchmarkFindSegment(b *testing.B) {
	xs := []float64{0, 1.5, 3, 4.5, 6, 7.5, 9, 12, 13.5, 16.5}
	for i := 0; i < b.N; i++ {
		findSegment(xs, 0)
		findSegment(xs, 16.5)
		findSegment(xs, 8.25)
		findSegment(xs, 4.125)
		findSegment(xs, 13.6)
		findSegment(xs, 13.5)
		findSegment(xs, 6)
		findSegment(xs, 4.5)
	}
}

func TestNewLinearInterpolator1D(t *testing.T) {
	t.Parallel()
	xs := []float64{0, 1, 2}
	i1d := NewLinearInterpolator1D(xs, []float64{-0.5, 1.5, 1})
	if xs[0] != i1d.begin() {
		t.Errorf("unexpected begin() value: got %g: want: %g", i1d.begin(), xs[0])
	}
	if xs[2] != i1d.end() {
		t.Errorf("unexpected end() value: got %g: want: %g", i1d.end(), xs[2])
	}
	type panicParams struct {
		xs              []float64
		ys              []float64
		expectedMessage string
	}
	panicParamSets := [...]panicParams{
		{xs, []float64{-0.5, 1.5}, "xs and ys have different lengths"},
		{[]float64{0.3}, []float64{0}, "too few points for interpolation"},
		{[]float64{0.3, 0.3}, []float64{0, 0}, "x values not strictly increasing"},
		{[]float64{0.3, -0.3}, []float64{0, 0}, "x values not strictly increasing"},
	}
	for _, params := range panicParamSets {
		panicked, message := panics(func() { NewLinearInterpolator1D(params.xs, params.ys) })
		expectedMessage := fmt.Sprintf("interp: %s", params.expectedMessage)
		if !panicked || message != expectedMessage {
			t.Errorf("expected panic for xs: %v and ys: %v with message: %s", params.xs, params.ys, expectedMessage)
		}
	}
}

func TestLinearInterpolator1DEval(t *testing.T) {
	t.Parallel()
	xs := []float64{0, 1, 2}
	ys := []float64{-0.5, 1.5, 1}
	i1d := NewLinearInterpolator1D(xs, ys)
	for i, x := range xs {
		y := i1d.eval(x)
		if y != ys[i] {
			t.Errorf("unexpected eval(%g) value: got: %g want: %g", x, y, x)
		}
	}
	type params struct {
		x         float64
		expectedY float64
	}
	paramSets := [...]params{{0.1, -0.3}, {0.5, 0.5}, {0.8, 1.1}, {1.2, 1.4}}
	const tolerance float64 = 1e-15
	for _, params := range paramSets {
		y := i1d.eval(params.x)
		if math.Abs(y-params.expectedY) > tolerance {
			t.Errorf("unexpected eval(%g) value: got: %g want: %g with tolerance: %g", params.x, y, params.expectedY, tolerance)
		}
	}
}
