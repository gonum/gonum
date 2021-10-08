package optimize_test

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"gonum.org/v1/gonum/optimize"
)

func ExampleNelderMead() {
	p := optimize.Problem{
		Func: func(x []float64) float64 {
			return math.Pow(x[0]-(4.0/5.0), 2)
		},
	}

	withoutBounds := &optimize.NelderMead{}
	result1, _ := optimize.Minimize(p, []float64{2}, &optimize.Settings{}, withoutBounds)
	fmt.Printf("without bounds: %.1f\n", result1.X[0])

	withBounds := &optimize.NelderMead{
		Bounds: []optimize.Bound{
			{Min: 1, Max: 3},
		},
	}
	result2, _ := optimize.Minimize(p, []float64{2}, &optimize.Settings{}, withBounds)
	fmt.Printf("with bounds: %.1f\n", result2.X[0])

	// Output:
	// without bounds: 0.8
	// with bounds: 1.0
}

func TestMinimize_NelderMead(t *testing.T) {
	isBoundsErr := func(t *testing.T, err error) {
		t.Helper()
		if err == nil || !strings.Contains(err.Error(), "incorrect number of bounds") {
			t.Errorf("it should return a bounds length error, got %#v", err)
		}
	}
	doesNotOccur := func(t *testing.T, err error) {
		t.Helper()
		if err != nil {
			t.Errorf("did not expect an error got %#v", err)
		}
	}

	tests := []struct {
		Name      string
		Bounds    []optimize.Bound
		X         []float64
		ExpectErr func(t *testing.T, err error)
	}{
		{
			Name:      "no bounds set",
			Bounds:    nil,
			X:         []float64{0},
			ExpectErr: doesNotOccur,
		},
		{
			Name:      "non nil bounds",
			Bounds:    []optimize.Bound{},
			X:         []float64{0},
			ExpectErr: isBoundsErr,
		},
		{
			Name:      "zero len bounds with non zero len x",
			Bounds:    []optimize.Bound{},
			X:         []float64{0, 0},
			ExpectErr: isBoundsErr,
		},
		{
			Name:      "equal len x and len bounds",
			Bounds:    []optimize.Bound{{}, {}},
			X:         []float64{0, 0},
			ExpectErr: doesNotOccur,
		},
		{
			Name:      "not equal len x and len bounds",
			Bounds:    []optimize.Bound{{}, {}, {}},
			X:         []float64{0},
			ExpectErr: isBoundsErr,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			p := optimize.Problem{
				Func: func(x []float64) float64 {
					return math.Pow(x[0]-(4.0/5.0), 2)
				},
			}

			withBounds := &optimize.NelderMead{
				Bounds: test.Bounds,
			}

			_, err := optimize.Minimize(p, test.X, &optimize.Settings{}, withBounds)
			test.ExpectErr(t, err)
		})
	}
}
