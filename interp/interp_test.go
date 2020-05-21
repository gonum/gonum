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
		x          float64
		expected_i int
		expected_x float64
	}
	param_sets := [...]params{{0, 0, 0}, {0.3, 0, 0}, {1, 1, 1}, {1.5, 1, 1}, {2, 2, 2}}
	for _, param := range param_sets {
		i, x := find_segment(xs, param.x)
		if i != param.expected_i || x != param.expected_x {
			t.Errorf("unexpected value of find_segment(xs, %g): got %d, %g want: %d, %g", param.x, i, x, param.expected_i, param.expected_x)
		}
	}
	panic_xs := [...]float64{-0.5, 2.1}
	expected_messages := [...]string{
		"interp: x value -0.5 below lower bound 0",
		"interp: x value 2.1 above upper bound 2",
	}
	for i, x := range panic_xs {
		panicked, message := panics(func() { find_segment(xs, x) })
		if !panicked || message != expected_messages[i] {
			t.Errorf("expected panic with message '%s' for evaluating at invalid x: %g", expected_messages[i], x)
		}
	}
}

func BenchmarkFindSegment(b *testing.B) {
	xs := []float64{0, 1.5, 3, 4.5, 6, 7.5, 9, 12, 13.5, 16.5}
	for i := 0; i < b.N; i++ {
		find_segment(xs, 0)
		find_segment(xs, 16.5)
		find_segment(xs, 8.25)
		find_segment(xs, 4.125)
		find_segment(xs, 13.6)
		find_segment(xs, 13.5)
		find_segment(xs, 6)
		find_segment(xs, 4.5)
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
	type panic_params struct {
		xs               []float64
		ys               []float64
		expected_message string
	}
	panic_param_sets := [...]panic_params{
		{xs, []float64{-0.5, 1.5}, "xs and ys have different lengths"},
		{[]float64{0.3}, []float64{0}, "too few points for interpolation"},
		{[]float64{0.3, 0.3}, []float64{0, 0}, "x values not strictly increasing"},
		{[]float64{0.3, -0.3}, []float64{0, 0}, "x values not strictly increasing"},
	}
	for _, params := range panic_param_sets {
		panicked, message := panics(func() { NewLinearInterpolator1D(params.xs, params.ys) })
		expected_message := fmt.Sprintf("interp: %s", params.expected_message)
		if !panicked || message != expected_message {
			t.Errorf("expected panic for xs: %v and ys: %v with message: %s", params.xs, params.ys, expected_message)
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
		x          float64
		expected_y float64
	}
	param_sets := [...]params{{0.1, -0.3}, {0.5, 0.5}, {0.8, 1.1}, {1.2, 1.4}}
	const tolerance float64 = 1e-15
	for _, params := range param_sets {
		y := i1d.eval(params.x)
		if math.Abs(y-params.expected_y) > tolerance {
			t.Errorf("unexpected eval(%g) value: got: %g want: %g with tolerance: %g", params.x, y, params.expected_y, tolerance)
		}
	}
}
