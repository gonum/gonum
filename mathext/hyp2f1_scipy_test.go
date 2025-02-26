package mathext

import (
	"fmt"
	"math"
	"testing"
)

func TestHyp2f1_Miscellaneous(t *testing.T) {
	t.Parallel()

	// Taken from test_hyp2f1.py test_miscellaneous.
	y := h(1.3, -0.2, 0.3, -2.1)
	expected := 1.8202169687521206
	if d := y - expected; math.Abs(d) > 5e-15 {
		t.Errorf("%f %f %f", y, expected, d)
	}
}

func TestHyp2f1_SomeRealPoints(t *testing.T) {
	t.Parallel()

	// Taken from test_mpmath.py test_hyp2f1_real_some_points.
	tests := []struct {
		a        float64
		b        float64
		c        float64
		z        float64
		expected float64
	}{
		{a: 1. / 3, b: 2. / 3, c: 5. / 6, z: 27. / 32, expected: 1.6},
		{a: 1. / 4, b: 1. / 2, c: 3. / 4, z: 80. / 81, expected: 1.8},
		{a: 0.7235, b: -1, c: -5, z: 0.3, expected: 1.04341},
		{a: 0.25, b: 1. / 3, c: 2, z: 0.999, expected: 1.0682644949603062},
		{a: 0.25, b: 1. / 3, c: 2, z: -1, expected: 0.9665658449252437},
		{a: 2, b: 3, c: 5, z: 0.99, expected: 27.699347904322664},
		{a: 3. / 2, b: -0.5, c: 3, z: 0.99, expected: 0.6840303684391167},
		{a: 2, b: 2.5, c: -3.25, z: 0.999, expected: 2.183739328012162e+26},
		{a: -8, b: 18.016500331508873, c: 10.805295997850628, z: 0.90875647507000001, expected: -3.566216341442061e-09},
		{a: -10, b: 900, c: -10.5, z: 0.99, expected: 2.5101757354622962e+22},
		{a: -10, b: 900, c: 10.5, z: 0.99, expected: 5.5748237303615776e+17},
		{a: -1, b: 2, c: 1, z: -1, expected: 3},
		{a: 0.5, b: 1 - 270.5, c: 1.5, z: 0.999 * 0.999, expected: 0.053963052503373715},
	}
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			y := h(test.a, test.b, test.c, test.z)
			if d := (y - test.expected) / test.expected; math.Abs(d) > 1e-10 {
				t.Errorf("%f %f %f", y, test.expected, d)
			}
		})
	}
}
