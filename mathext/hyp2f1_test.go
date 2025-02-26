// Tests in this file are taken from R's hypergeo package.
// https://github.com/cran/hypergeo
package mathext

import (
	"flag"
	"fmt"
	"log"
	"math"
	"testing"
)

var tests = []struct {
	x float64
	y float64
}{
	{x: 0.28, y: 1.3531156987873853569937},
	{x: -0.79, y: 0.5773356740314405932679},
	{x: 0.56, y: 2.1085704049533617876477},
	{x: -2.13, y: 0.3352446571148822718200},
	{x: -0.43, y: 0.7150355048137748692483},
	{x: -1.23, y: 0.4670987707934830535095},
}

func TestHyp2f1(t *testing.T) {
	t.Parallel()
	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			y := h(1.21, 1.443, 1.88, test.x)
			if d := y - test.y; math.Abs(d) > 1e-12 {
				t.Errorf("%f %f %f %f", test.x, y, test.y, d)
			}
		})
	}
}

func TestHyp2f1_15_1_15(t *testing.T) {
	t.Parallel()

	// eqn15_1_15a is the right hand side of equation 15.1.15 of
	// M. Abramowitz and I. A. Stegun 1965. Handbook of Mathematical Functions, New York: Dover.
	eqn15_1_15a := func(a, z float64) float64 {
		return h(a, 1-a, 3./2, math.Pow(math.Sin(z), 2))
	}
	eqn15_1_15b := func(a, z float64) float64 {
		return math.Sin((2*a-1)*z) / ((2*a - 1) * math.Sin(z))
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			// Ignore z=-2.13, since both R's hypergeo and Maple don't handle this case either.
			if test.x == -2.13 {
				return
			}

			a := eqn15_1_15a(0.2, test.x)
			b := eqn15_1_15b(0.2, test.x)
			if d := a - b; math.Abs(d) > 1e-6 {
				t.Errorf("%f %f %f %f", test.x, a, b, d)
			}
		})
	}
}

func TestHyp2f1_15_2_10(t *testing.T) {
	t.Parallel()

	eqn15_2_10 := func(a, b, c, z float64) float64 {
		return (c-a)*h(a-1, b, c, z) + (2*a-c-a*z+b*z)*h(a, b, c, z) + a*(z-1)*h(a+1, b, c, z)
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			v := eqn15_2_10(0.1, 0.44, 0.611, test.x)
			if math.Abs(v) > 1e-6 {
				t.Errorf("%f %f", test.x, v)
			}
		})
	}
}

func TestHyp2f1_15_1(t *testing.T) {
	t.Parallel()

	type equation struct {
		name string
		f    func(float64) float64
	}

	equations := []equation{
		{
			name: "15_1_3",
			f: func(z float64) float64 {
				lhs := h(1, 1, 2, z)
				rhs := -math.Log(1-z) / z
				return math.Abs(rhs - lhs)
			},
		},
		{
			name: "15_1_5",
			f: func(z float64) float64 {
				lhs := h(1./2, 1, 3./2, -z*z)
				rhs := math.Atan(z) / z
				return math.Abs(rhs - lhs)
			},
		},
		{
			name: "15_1_7a",
			f: func(z float64) float64 {
				lhs := h(1./2, 1./2, 3./2, -z*z)
				rhs := math.Sqrt(1+z*z) * h(1, 1, 3./2, -z*z)
				return math.Abs(rhs - lhs)
			},
		},
		{
			name: "15_1_7b",
			f: func(z float64) float64 {
				lhs := h(1./2, 1./2, 3./2, -z*z)
				rhs := math.Log(z+math.Sqrt(1+z*z)) / z
				return math.Abs(rhs - lhs)
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()
			for _, eqn := range equations {
				if d := eqn.f(test.x); d > 1e-10 {
					t.Errorf("%s %f %f", eqn.name, test.x, d)
				}
			}
		})
	}
}

func TestHyp2f1_15_1_zz(t *testing.T) {
	t.Parallel()

	type equation struct {
		name string
		f    func(float64) float64
	}

	equations := []equation{
		{
			name: "15_1_4",
			f: func(z float64) float64 {
				lhs := h(1./2, 1, 3./2, z*z)
				rhs := 0.5 * math.Log((1+z)/(1-z)) / z
				return math.Abs(rhs - lhs)
			},
		},
		{
			name: "15_1_6a",
			f: func(z float64) float64 {
				lhs := h(1./2, 1./2, 3./2, z*z)
				rhs := math.Sqrt(1-z*z) * h(1, 1, 3./2, z*z)
				return math.Abs(rhs - lhs)
			},
		},
		{
			name: "15_1_6b",
			f: func(z float64) float64 {
				lhs := h(1./2, 1./2, 3./2, z*z)
				rhs := math.Asin(z) / z
				return math.Abs(rhs - lhs)
			},
		},
	}

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			t.Parallel()

			if test.x*test.x > 1 {
				return
			}

			for _, eqn := range equations {
				if d := eqn.f(test.x); d > 1e-10 {
					t.Errorf("%s %f %f", eqn.name, test.x, d)
				}
			}
		})
	}
}

func TestHyp2f1_Igor_Kojanov(t *testing.T) {
	t.Parallel()

	var y float64
	y = h(1, 2, 3, 0)
	if d := y - 1; d != 0 {
		t.Errorf("%f %f", y, d)
	}

	y = h(1, 1.64, 2.64, -0.1111)
	if d := y - 0.9361003540660249866434; math.Abs(d) > 1e-15 {
		t.Errorf("%f %f", y, d)
	}
}

func TestHyp2f1_John_Ormerod(t *testing.T) {
	t.Parallel()

	y := h(5.25, 1, 6.5, 0.501)
	if d := y - 1.70239432012007391092082702795; math.Abs(d) > 1e-10 {
		t.Errorf("%f %f", y, d)
	}
}

func h(a, b, c, x float64) float64 {
	y, err := Hyp2f1(a, b, c, x)
	if err != nil {
		panic(err)
	}
	return y
}

func TestMain(m *testing.M) {
	flag.Parse()
	log.SetFlags(log.Lmicroseconds | log.Llongfile | log.LstdFlags)

	m.Run()
}
