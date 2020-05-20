package interp

import (
	"math"
	"testing"
)

func TestNewConstInterpolator1D(t *testing.T) {
	const value float64 = 42.0
	i1d := NewConstInterpolator1D(value)
	if i1d.start() != math.Inf(-1) {
		t.Errorf("unexpected start() value: got: %g want: %g", i1d.start(), math.Inf(-1))
	}
	if i1d.end() != math.Inf(1) {
		t.Errorf("unexpected end() value: got: %g want: %g", i1d.end(), math.Inf(1))
	}
	var xs = [...]float64{math.Inf(-1), -11, 0.4, 1e9, math.Inf(1)}
	for _, x := range xs {
		y := i1d.eval(x)
		if y != value {
			t.Errorf("unexpected value of evaluate(%g): got: %g want: %g", x, y, value)
		}
	}
}
