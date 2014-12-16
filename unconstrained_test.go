// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"math"
	"testing"

	"github.com/gonum/blas/goblas"
	"github.com/gonum/floats"
	"github.com/gonum/matrix/mat64"
)

var negInf = math.Inf(-1)

func init() {
	mat64.Register(goblas.Blas{})
}

// The Fletcher-Powell helical valley function
// Dim = 3
// X0 = [-1, 0, 0]
// OptX = [1, 0, 0]
// OptF = 0
type HelicalValley struct{}

func (HelicalValley) F(x []float64) float64 {
	θ := 0.5 * math.Atan2(x[1], x[0]) / math.Pi
	r := math.Sqrt(math.Pow(x[0], 2) + math.Pow(x[1], 2))

	f1 := 10 * (x[2] - 10*θ)
	f2 := 10 * (r - 1)
	f3 := x[2]

	return math.Pow(f1, 2) + math.Pow(f2, 2) + math.Pow(f3, 2)
}

func (HelicalValley) Df(x, g []float64) {
	θ := 0.5 * math.Atan2(x[1], x[0]) / math.Pi
	r := math.Sqrt(math.Pow(x[0], 2) + math.Pow(x[1], 2))
	s := x[2] - 10*θ
	t := 5 * s / math.Pow(r, 2) / math.Pi

	g[0] = 200 * (x[0] - x[0]/r + x[1]*t)
	g[1] = 200 * (x[1] - x[1]/r - x[0]*t)
	g[2] = 2 * (x[2] + 100*s)
}

// The Biggs' EXP2 function
// M.C. Biggs, Minimization algorithms making use of non-quadratic properties
// of the objective function. J. Inst. Maths Applics 8 (1971), 315-327.
// Dim = 2
// X0 = [1, 2]
// OptX = [1, 10]
// OptF = 0
type BiggsEXP2 struct{}

func (BiggsEXP2) F(x []float64) (sum float64) {
	for i := 1; i <= 10; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z)
		f := math.Exp(-x[0]*z) - 5*math.Exp(-x[1]*z) - y
		sum += math.Pow(f, 2)
	}
	return sum
}

func (BiggsEXP2) Df(x, g []float64) {
	for i := 0; i < len(g); i++ {
		g[i] = 0
	}
	for i := 1; i <= 10; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z)
		f := math.Exp(-x[0]*z) - 5*math.Exp(-x[1]*z) - y

		dfdx0 := -z * math.Exp(-x[0]*z)
		dfdx1 := 5 * z * math.Exp(-x[1]*z)

		g[0] += 2 * f * dfdx0
		g[1] += 2 * f * dfdx1
	}
}

// The Biggs' EXP3 function
// M.C. Biggs, Minimization algorithms making use of non-quadratic properties
// of the objective function. J. Inst. Maths Applics 8 (1971), 315-327.
// Dim = 3
// X0 = [1, 2, 1]
// OptX = [1, 10, 5]
// OptF = 0
type BiggsEXP3 struct{}

func (BiggsEXP3) F(x []float64) (sum float64) {
	for i := 1; i <= 10; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z)
		f := math.Exp(-x[0]*z) - x[2]*math.Exp(-x[1]*z) - y
		sum += math.Pow(f, 2)
	}
	return sum
}

func (BiggsEXP3) Df(x, g []float64) {
	for i := 0; i < len(g); i++ {
		g[i] = 0
	}
	for i := 1; i <= 10; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z)
		f := math.Exp(-x[0]*z) - x[2]*math.Exp(-x[1]*z) - y

		dfdx0 := -z * math.Exp(-x[0]*z)
		dfdx1 := x[2] * z * math.Exp(-x[1]*z)
		dfdx2 := -math.Exp(-x[1] * z)

		g[0] += 2 * f * dfdx0
		g[1] += 2 * f * dfdx1
		g[2] += 2 * f * dfdx2
	}
}

// The Biggs' EXP4 function
// M.C. Biggs, Minimization algorithms making use of non-quadratic properties
// of the objective function. J. Inst. Maths Applics 8 (1971), 315-327.
// Dim = 4
// X0 = [1, 2, 1, 1]
// OptX = [1, 10, 1, 5]
// OptF = 0
type BiggsEXP4 struct{}

func (BiggsEXP4) F(x []float64) (sum float64) {
	for i := 1; i <= 10; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z)
		f := x[2]*math.Exp(-x[0]*z) - x[3]*math.Exp(-x[1]*z) - y
		sum += math.Pow(f, 2)
	}
	return sum
}

func (BiggsEXP4) Df(x, g []float64) {
	for i := 0; i < len(g); i++ {
		g[i] = 0
	}
	for i := 1; i <= 10; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z)
		f := x[2]*math.Exp(-x[0]*z) - x[3]*math.Exp(-x[1]*z) - y

		dfdx0 := -z * x[2] * math.Exp(-x[0]*z)
		dfdx1 := z * x[3] * math.Exp(-x[1]*z)
		dfdx2 := math.Exp(-x[0] * z)
		dfdx3 := -math.Exp(-x[1] * z)

		g[0] += 2 * f * dfdx0
		g[1] += 2 * f * dfdx1
		g[2] += 2 * f * dfdx2
		g[3] += 2 * f * dfdx3
	}
}

// The Biggs' EXP5 function
// M.C. Biggs, Minimization algorithms making use of non-quadratic properties
// of the objective function. J. Inst. Maths Applics 8 (1971), 315-327.
// Dim = 5
// X0 = [1, 2, 1, 1, 1]
// OptX = [1, 10, 1, 5, 4]
// OptF = 0
type BiggsEXP5 struct{}

func (BiggsEXP5) F(x []float64) (sum float64) {
	for i := 1; i <= 11; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z) + 3*math.Exp(-4*z)
		f := x[2]*math.Exp(-x[0]*z) - x[3]*math.Exp(-x[1]*z) + 3*math.Exp(-x[4]*z) - y
		sum += math.Pow(f, 2)
	}
	return sum
}

func (BiggsEXP5) Df(x, g []float64) {
	for i := 0; i < len(g); i++ {
		g[i] = 0
	}
	for i := 1; i <= 11; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z) + 3*math.Exp(-4*z)
		f := x[2]*math.Exp(-x[0]*z) - x[3]*math.Exp(-x[1]*z) + 3*math.Exp(-x[4]*z) - y

		dfdx0 := -z * x[2] * math.Exp(-x[0]*z)
		dfdx1 := z * x[3] * math.Exp(-x[1]*z)
		dfdx2 := math.Exp(-x[0] * z)
		dfdx3 := -math.Exp(-x[1] * z)
		dfdx4 := -3 * z * math.Exp(-x[4]*z)

		g[0] += 2 * f * dfdx0
		g[1] += 2 * f * dfdx1
		g[2] += 2 * f * dfdx2
		g[3] += 2 * f * dfdx3
		g[4] += 2 * f * dfdx4
	}
}

// The Biggs' EXP6 function
// M.C. Biggs, Minimization algorithms making use of non-quadratic properties
// of the objective function. J. Inst. Maths Applics 8 (1971), 315-327.
// Dim = 6
// X0 = [1, 2, 1, 1, 1, 1]
// OptX = [1, 10, 1, 5, 4, 3]
// OptF = 0
// OptF = 0.005655649925...
type BiggsEXP6 struct{}

func (BiggsEXP6) F(x []float64) (sum float64) {
	for i := 1; i <= 13; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z) + 3*math.Exp(-4*z)
		f := x[2]*math.Exp(-x[0]*z) - x[3]*math.Exp(-x[1]*z) + x[5]*math.Exp(-x[4]*z) - y
		sum += math.Pow(f, 2)
	}
	return sum
}

func (BiggsEXP6) Df(x, g []float64) {
	for i := 0; i < len(g); i++ {
		g[i] = 0
	}
	for i := 1; i <= 13; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z) + 3*math.Exp(-4*z)
		f := x[2]*math.Exp(-x[0]*z) - x[3]*math.Exp(-x[1]*z) + x[5]*math.Exp(-x[4]*z) - y

		dfdx0 := -z * x[2] * math.Exp(-x[0]*z)
		dfdx1 := z * x[3] * math.Exp(-x[1]*z)
		dfdx2 := math.Exp(-x[0] * z)
		dfdx3 := -math.Exp(-x[1] * z)
		dfdx4 := -z * x[5] * math.Exp(-x[4]*z)
		dfdx5 := math.Exp(-x[4] * z)

		g[0] += 2 * f * dfdx0
		g[1] += 2 * f * dfdx1
		g[2] += 2 * f * dfdx2
		g[3] += 2 * f * dfdx3
		g[4] += 2 * f * dfdx4
		g[5] += 2 * f * dfdx5
	}
}

// Gaussian function
// J. More, B.S. Garbow, K.E. Hillstrom, Testing unconstrained optimization software.
// ACM Trans.Math. Softw. 7 (1981), 17-41.
// Dim = 3
// X0 = [0.4, 1, 0]
// OptX = [0.3989561..., 1.0000191..., 0]
// OptF = 1.12793...e-8
type Gaussian struct{}

func (Gaussian) y(i int) (yi float64) {
	switch i {
	case 1, 15:
		yi = 0.0009
	case 2, 14:
		yi = 0.0044
	case 3, 13:
		yi = 0.0175
	case 4, 12:
		yi = 0.0540
	case 5, 11:
		yi = 0.1295
	case 6, 10:
		yi = 0.2420
	case 7, 9:
		yi = 0.3521
	case 8:
		yi = 0.3989
	}
	return yi
}

func (g Gaussian) F(x []float64) (sum float64) {
	for i := 1; i <= 15; i++ {
		c := float64(8-i) / 2
		d := math.Pow(c-x[2], 2)
		e := math.Exp(-x[1] * d / 2)
		f := x[0]*e - g.y(i)
		sum += f * f
	}
	return sum
}

func (g Gaussian) Df(x, grad []float64) {
	grad[0] = 0
	grad[1] = 0
	grad[2] = 0
	for i := 1; i <= 15; i++ {
		c := float64(8-i) / 2
		d := math.Pow(c-x[2], 2)
		e := math.Exp(-x[1] * d / 2)
		f := x[0]*e - g.y(i)

		grad[0] += 2 * f * e
		grad[1] -= f * e * d * x[0]
		grad[2] += 2 * f * e * x[0] * x[1] * (c - x[2])
	}
}

// The Powell's badly scaled function
// J. More, B.S. Garbow, K.E. Hillstrom, Testing unconstrained optimization software.
// ACM Trans.Math. Softw. 7 (1981), 17-41.
// Dim = 2
// X0 = [0, 1]
// OptX = [1.09815933...e-5, 9.10614674...]
// OptF = 0
type Powell struct{}

func (Powell) F(x []float64) float64 {
	f1 := 1e4*x[0]*x[1] - 1
	f2 := math.Exp(-x[0]) + math.Exp(-x[1]) - 1.0001
	return math.Pow(f1, 2) + math.Pow(f2, 2)
}

func (Powell) Df(x, grad []float64) {
	f1 := 1e4*x[0]*x[1] - 1
	f2 := math.Exp(-x[0]) + math.Exp(-x[1]) - 1.0001

	grad[0] = 2 * (1e4*f1*x[1] - f2*math.Exp(-x[0]))
	grad[1] = 2 * (1e4*f1*x[0] - f2*math.Exp(-x[1]))
}

// The Box' three-dimensional function
// J. More, B.S. Garbow, K.E. Hillstrom, Testing unconstrained optimization software.
// ACM Trans.Math. Softw. 7 (1981), 17-41.
// Dim = 3
// X0 = [0, 10, 20]
// OptX = [1, 10, 1], [10, 1, -1], [a, a, 0]
// OptF = 0
type Box struct{}

func (Box) F(x []float64) (sum float64) {
	for i := 1; i <= 10; i++ {
		c := float64(i) / 10
		y := math.Exp(-c) - math.Exp(10*c)
		f := math.Exp(-c*x[0]) - math.Exp(-c*x[1]) - x[2]*y
		sum += math.Pow(f, 2)
	}
	return sum
}

func (Box) Df(x, grad []float64) {
	grad[0] = 0
	grad[1] = 0
	grad[2] = 0

	for i := 1; i <= 10; i++ {
		c := float64(i) / 10
		y := math.Exp(-c) - math.Exp(10*c)
		f := math.Exp(-c*x[0]) - math.Exp(-c*x[1]) - x[2]*y

		grad[0] += -2 * f * c * math.Exp(-c*x[0])
		grad[1] += -2 * f * c * math.Exp(-c*x[1])
		grad[2] += -2 * f * y
	}
}

// Variably dimensioned function
// J. More, B.S. Garbow, K.E. Hillstrom, Testing unconstrained optimization software.
// ACM Trans.Math. Softw. 7 (1981), 17-41.
// Dim = n
// X0 = [..., (n-i)/n, ...]
// OptX = [1, ..., 1]
// OptF = 0
type VariablyDimensioned struct{}

func (v VariablyDimensioned) F(x []float64) (sum float64) {
	for i := 0; i < len(x); i++ {
		sum += math.Pow(x[i]-1, 2)
	}
	s := 0.0
	for i := 0; i < len(x); i++ {
		s += float64(i+1) * (x[i] - 1)
	}
	sum += math.Pow(s, 2) + math.Pow(s, 4)
	return sum
}

func (v VariablyDimensioned) Df(x, grad []float64) {
	for i := 0; i < len(grad); i++ {
		grad[i] = 0
	}
	s := 0.0
	for i := 0; i < len(x); i++ {
		s += float64(i+1) * (x[i] - 1)
	}
	for i := 0; i < len(grad); i++ {
		grad[i] = 2 * ((x[i] - 1) + s*float64(i+1)*(1+2*math.Pow(s, 2)))
	}
}

// The Watson function
// J. More, B.S. Garbow, K.E. Hillstrom, Testing unconstrained optimization software.
// ACM Trans.Math. Softw. 7 (1981), 17-41.
// For Dim = 9, the problem of minimizing the Watson function is very ill conditioned.
// Dim = n, 2 <= n <= 31
// X0 = [0, ..., 0]
// OptX = [1, ..., 1]
// For Dim = 6 also:
// OptX = [-0.015725, 1.012435, -0.232992, 1.260430, -1.513729, 0.992996]
// OptF = 2.287687...e-3
// For Dim = 9 also:
// OptX = [-0.000015, 0.999790, 0.014764, 0.146342, 1.000821, -2.617731, 4.104403, -3.143612, 1.052627]
// OptF = 1.39976...e-6
// For Dim = 12 also:
// OptF = 4.72238...e-10
type Watson struct{}

func (Watson) F(x []float64) (sum float64) {
	for i := 1; i <= 29; i++ {
		c := float64(i) / 29

		s1 := 0.0
		for j := 1; j < len(x); j++ {
			s1 += float64(j) * x[j] * math.Pow(c, float64(j-1))
		}

		s2 := 0.0
		for j := 0; j < len(x); j++ {
			s2 += x[j] * math.Pow(c, float64(j))
		}
		s2 = math.Pow(s2, 2)

		sum += math.Pow(s1-s2-1, 2)
	}
	sum += math.Pow(x[0], 2)
	sum += math.Pow(x[1]-math.Pow(x[0], 2)-1, 2)
	return sum
}

func (Watson) Df(x, grad []float64) {
	for i := 0; i < len(grad); i++ {
		grad[i] = 0
	}

	for i := 1; i <= 29; i++ {
		c := float64(i) / 29

		s1 := 0.0
		for j := 1; j < len(x); j++ {
			s1 += float64(j) * x[j] * math.Pow(c, float64(j-1))
		}

		s2 := 0.0
		for j := 0; j < len(x); j++ {
			s2 += x[j] * math.Pow(c, float64(j))
		}

		t := s1 - math.Pow(s2, 2) - 1
		for j := 0; j < len(x); j++ {
			grad[j] += 2 * t * math.Pow(c, float64(j-1)) * (float64(j) - 2*s2*c)
		}
	}
	t := x[1] - math.Pow(x[0], 2) - 1
	grad[0] += 2 * (1 - 2*t) * x[0]
	grad[1] += 2 * t
}

// Penalty function #1
// J. More, B.S. Garbow, K.E. Hillstrom, Testing unconstrained optimization software.
// ACM Trans.Math. Softw. 7 (1981), 17-41.
// Dim = n
// X0 = [..., i-1, i, i+1...]
// For Dim = 4:
// OptF = 2.2499775...e-5
// For Dim = 10:
// OptF = 7.0876515...e-5
type Penalty1 struct{}

func (Penalty1) F(x []float64) (sum float64) {
	for i := 0; i < len(x); i++ {
		sum += math.Pow(x[i]-1, 2)
	}
	sum *= 1e-5

	s := 0.0
	for i := 0; i < len(x); i++ {
		s += math.Pow(x[i], 2)
	}
	sum += math.Pow(s-0.25, 2)

	return sum
}

func (Penalty1) Df(x, grad []float64) {
	s := 0.0
	for i := 0; i < len(x); i++ {
		s += math.Pow(x[i], 2)
	}
	s -= 0.25

	for i := 0; i < len(grad); i++ {
		grad[i] = 2 * (2*s*x[i] + 1e-5*(x[i]-1))
	}
}

// Penalty function #2
// J. More, B.S. Garbow, K.E. Hillstrom, Testing unconstrained optimization software.
// ACM Trans.Math. Softw. 7 (1981), 17-41.
// Dim = n
// X0 = [0.5, ..., 0.5]
// For Dim = 4:
// OptF = 9.37629...e-6
// For Dim = 10:
// OptF = 2.93660...e-4
type Penalty2 struct{}

func (Penalty2) F(x []float64) (sum float64) {
	dim := len(x)

	s := 0.0
	for i := 0; i < dim; i++ {
		s += float64(dim-i) * math.Pow(x[i], 2)
	}
	s -= 1

	for i := 1; i < dim; i++ {
		yi := math.Exp(float64(i+1)/10) + math.Exp(float64(i)/10)
		f := math.Exp(x[i]/10) + math.Exp(x[i-1]/10) - yi
		sum += math.Pow(f, 2)
	}
	for i := 1; i < dim; i++ {
		f := math.Exp(x[i]/10) - math.Exp(-1.0/10)
		sum += math.Pow(f, 2)
	}
	sum *= 1e-5

	sum += math.Pow(x[0]-0.2, 2)
	sum += math.Pow(s, 2)

	return sum
}

func (Penalty2) Df(x, grad []float64) {
	dim := len(x)

	s := 0.0
	for i := 0; i < dim; i++ {
		s += float64(dim-i) * math.Pow(x[i], 2)
	}
	s -= 1

	for i := 0; i < dim; i++ {
		grad[i] = 4 * s * float64(dim-i) * x[i]
	}
	for i := 1; i < dim; i++ {
		yi := math.Exp(float64(i+1)/10) + math.Exp(float64(i)/10)
		f := math.Exp(x[i]/10) + math.Exp(x[i-1]/10) - yi
		grad[i] += 1e-5 * f * math.Exp(x[i]/10) / 5
		grad[i-1] += 1e-5 * f * math.Exp(x[i-1]/10) / 5
	}
	for i := 1; i < dim; i++ {
		f := math.Exp(x[i]/10) - math.Exp(-1.0/10)
		grad[i] += 1e-5 * f * math.Exp(x[i]/10) / 5
	}
	grad[0] += 2 * (x[0] - 0.2)
}

// The Brown's badly scaled function
// J. More, B.S. Garbow, K.E. Hillstrom, Testing unconstrained optimization software.
// ACM Trans.Math. Softw. 7 (1981), 17-41.
// Dim = 2
// X0 = [1, 1]
// OptX = [1e6, 2*1e-6]
// OptF = 0
type Brown struct{}

func (Brown) F(x []float64) (sum float64) {
	sum = math.Pow(x[0]-1e6, 2)
	sum += math.Pow(x[1]-2e-6, 2)
	sum += math.Pow(x[0]*x[1]-2, 2)
	return sum
}

func (Brown) Df(x, g []float64) {
	f1 := x[0] - 1e6
	f2 := x[1] - 2e-6
	f3 := x[0]*x[1] - 2

	g[0] = 2*f1 + 2*f3*x[1]
	g[1] = 2*f2 + 2*f3*x[0]
}

// The Brown and Dennis function
// J. More, B.S. Garbow, K.E. Hillstrom, Testing unconstrained optimization software.
// ACM Trans.Math. Softw. 7 (1981), 17-41.
// Dim = 4
// X0 = [25, 5, -5, -1]
// OptF = 85822.2
type BrownDennis struct{}

func (BrownDennis) F(x []float64) (sum float64) {
	for i := 0; i < 20; i++ {
		c := float64(i+1) / 5
		d1 := x[0] + c*x[1] - math.Exp(c)
		d2 := x[2] + x[3]*math.Sin(c) - math.Cos(c)
		f := math.Pow(d1, 2) + math.Pow(d2, 2)
		sum += math.Pow(f, 2)
	}
	return sum
}

func (BrownDennis) Df(x, grad []float64) {
	for i := 0; i < len(grad); i++ {
		grad[i] = 0
	}

	for i := 0; i < 20; i++ {
		c := float64(i+1) / 5
		d1 := x[0] + c*x[1] - math.Exp(c)
		d2 := x[2] + x[3]*math.Sin(c) - math.Cos(c)
		f := math.Pow(d1, 2) + math.Pow(d2, 2)
		grad[0] += 4 * f * d1
		grad[1] += 4 * f * d1 * c
		grad[2] += 4 * f * d2
		grad[3] += 4 * f * d2 * math.Sin(c)
	}
}

// The Gulf R&D function
// J. More, B.S. Garbow, K.E. Hillstrom, Testing unconstrained optimization software.
// ACM Trans.Math. Softw. 7 (1981), 17-41.
// Dim = 3
// X0 = [5, 2.5, 0.15]
// OptF = 0 (global)
// OptX = [50, 25, 1.5]
//
// OptF = 0.038 (local)
// OptX = [99.89537833835886, 60.61453903025014, 9.16124389236433]
// This local minimum is very flat and the minimizer is surrounded by a
// "plateau", where the gradient is zero everywhere and the function equals
// 0.0385.
//
// OptF = 0.038 (local)
// OptX = [201.662589489426, 60.616331504682, 10.224891158489]
type GulfRD struct{}

func (GulfRD) F(x []float64) (sum float64) {
	for i := 0; i < 100; i++ {
		c := float64(i+1) / 100
		yi := 25 + math.Pow(-50*math.Log(c), 2.0/3.0)
		d := math.Abs(yi - x[1])
		e := math.Pow(d, x[2]) / x[0]
		f := math.Exp(-e) - c

		sum += math.Pow(f, 2)
	}
	return sum
}

func (GulfRD) Df(x, grad []float64) {
	for i := 0; i < len(grad); i++ {
		grad[i] = 0
	}

	for i := 0; i < 100; i++ {
		c := float64(i+1) / 100
		yi := 25 + math.Pow(-50*math.Log(c), 2.0/3.0)
		d := math.Abs(yi - x[1])
		e := math.Pow(d, x[2]) / x[0]
		f := math.Exp(-e) - c

		grad[0] += 2 * f * math.Exp(-e) * math.Pow(d, x[2]) / math.Pow(x[0], 2)
		grad[1] += 2 * f * math.Exp(-e) * e * x[2] / d
		grad[2] -= 2 * f * math.Exp(-e) * e * math.Log(d)
	}
}

// The Trigonometric function
// J. More, B.S. Garbow, K.E. Hillstrom, Testing unconstrained optimization software.
// ACM Trans.Math. Softw. 7 (1981), 17-41.
// Dim = n
// X0 = [1/n, ..., 1/n]
//
// OptF = 0 (global)
// OptX = [0.0429645656464827, 0.043976286943042, 0.0450933996565844, 0.0463389157816542, 0.0477443839560646,
//         0.0493547352444078, 0.0512373449259557, 0.195209463715277, 0.164977664720328, 0.0601485854398078]
//
// OptF = 2.79506e-5 (local)
// OptX = [0.0551509, 0.0568408, 0.0587639, 0.0609906, 0.0636263,
//         0.0668432, 0.208162, 0.164363, 0.0850068, 0.0914314]
type Trigonometric struct{}

func (Trigonometric) F(x []float64) (sum float64) {
	dim := len(x)

	s1 := 0.0
	for j := 0; j < dim; j++ {
		s1 += math.Cos(x[j])
	}
	for i := 0; i < dim; i++ {
		f := float64(dim+i) - float64(i)*math.Cos(x[i]) - math.Sin(x[i]) - s1
		sum += math.Pow(f, 2)
	}

	return sum
}

func (Trigonometric) Df(x, grad []float64) {
	dim := len(x)
	for i := 0; i < dim; i++ {
		grad[i] = 0
	}

	s1 := 0.0
	for j := 0; j < dim; j++ {
		s1 += math.Cos(x[j])
	}

	s2 := 0.0
	for i := 0; i < dim; i++ {
		f := float64(dim+i) - float64(i)*math.Cos(x[i]) - math.Sin(x[i]) - s1
		s2 += f
		grad[i] = 2 * f * (float64(i)*math.Sin(x[i]) - math.Cos(x[i]))
	}
	for i := 0; i < dim; i++ {
		grad[i] += 2 * s2 * math.Sin(x[i])
	}
}

// The extended Rosenbrock function
// Very difficult to minimize if the starting point is far from the minimum.
// Dim = n
// X0 = [-1.2, 1] for Dim = 2
// OptF = 0 (global)
// OptX = [1, ..., 1]
type Rosenbrock struct{}

func (Rosenbrock) F(x []float64) (sum float64) {
	for i := 0; i < len(x)-1; i++ {
		sum += math.Pow(1-x[i], 2) + 100*math.Pow(x[i+1]-math.Pow(x[i], 2), 2)
	}
	return sum
}

func (Rosenbrock) Df(x, grad []float64) {
	dim := len(x)
	for i := range grad {
		grad[i] = 0
	}

	for i := 0; i < dim-1; i++ {
		grad[i] -= 2 * (1 - x[i])
		grad[i] -= 400 * (x[i+1] - math.Pow(x[i], 2)) * x[i]
	}
	for i := 1; i < dim; i++ {
		grad[i] += 200 * (x[i] - math.Pow(x[i-1], 2))
	}
}

// The Extended Powell singular function
// Dim = n multiple of 4
// X0 = [3, -1, 0, 1, 3, -1, 0, 1, ...]
// OptF = 0
// OptX = [0, ..., 0]
type ExtendedPowell struct{}

func (ExtendedPowell) F(x []float64) (sum float64) {
	dim := len(x)
	if dim%4 != 0 {
		panic("dimension of the problem must be a multiple of 4")
	}

	for i := 0; i < dim; i += 4 {
		f1 := x[i] + 10*x[i+1]
		f2 := x[i+2] - x[i+3]
		f3 := math.Pow(x[i+1]-2*x[i+2], 2)
		f4 := math.Pow(x[i]-x[i+3], 2)
		sum += math.Pow(f1, 2) + 5*math.Pow(f2, 2) + math.Pow(f3, 2) + 10*math.Pow(f4, 2)
	}
	return sum
}

func (ExtendedPowell) Df(x, grad []float64) {
	dim := len(x)
	if dim%4 != 0 {
		panic("dimension of the problem must be a multiple of 4")
	}

	for i := 0; i < dim; i += 4 {
		f1 := x[i] + 10*x[i+1]
		f2 := x[i+2] - x[i+3]
		f3 := math.Pow(x[i+1]-2*x[i+2], 2)
		f4 := math.Pow(x[i]-x[i+3], 2)

		grad[i] = 2*f1 + 40*f4*(x[i]-x[i+3])
		grad[i+1] = 20*f1 + 4*f3*(x[i+1]-2*x[i+2])
		grad[i+2] = 10*f2 - 8*f3*(x[i+1]-2*x[i+2])
		grad[i+3] = -10*f2 - 40*f4*(x[i]-x[i+3])
	}
}

// The Beale function
// Dim = 2
// X0 = [1, 1]
// OptF = 0
// OptX = [3, 0.5]
type Beale struct{}

func (Beale) F(x []float64) (sum float64) {
	f1 := 1.5 - x[0]*(1-x[1])
	f2 := 2.25 - x[0]*(1-math.Pow(x[1], 2))
	f3 := 2.625 - x[0]*(1-math.Pow(x[1], 3))
	return math.Pow(f1, 2) + math.Pow(f2, 2) + math.Pow(f3, 2)
}

func (Beale) Df(x, grad []float64) {
	f1 := 1.5 - x[0]*(1-x[1])
	f2 := 2.25 - x[0]*(1-math.Pow(x[1], 2))
	f3 := 2.625 - x[0]*(1-math.Pow(x[1], 3))

	grad[0] = -2 * (f1*(1-x[1]) + f2*(1-math.Pow(x[1], 2)) + f3*(1-math.Pow(x[1], 3)))
	grad[1] = 2 * x[0] * (f1 + 2*f2*x[1] + 3*f3*math.Pow(x[1], 2))
}

// The Wood function
// Dim = 4
// X0 = [-3, -1, -3, -1]
// OptF = 0
// OptX = [1, 1, 1, 1]
type Wood struct{}

func (Wood) F(x []float64) (sum float64) {
	f1 := x[1] - math.Pow(x[0], 2)
	f2 := 1 - x[0]
	f3 := x[3] - math.Pow(x[2], 2)
	f4 := 1 - x[2]
	f5 := x[1] + x[3] - 2
	f6 := x[2] - x[3]

	sum = 100*math.Pow(f1, 2) + math.Pow(f2, 2) + 90*math.Pow(f3, 2)
	sum += math.Pow(f4, 2) + 10*math.Pow(f5, 2) + 0.1*math.Pow(f6, 2)
	return sum
}

func (Wood) Df(x, grad []float64) {
	f1 := x[1] - math.Pow(x[0], 2)
	f2 := 1 - x[0]
	f3 := x[3] - math.Pow(x[2], 2)
	f4 := 1 - x[2]
	f5 := x[1] + x[3] - 2
	f6 := x[2] - x[3]

	grad[0] = -2 * (200*f1*x[0] + f2)
	grad[1] = 2 * (100*f1 + 10*f5)
	grad[2] = 2 * (-180*f3*x[2] - f4 + 0.1*f6)
	grad[3] = 2 * (90*f3 + 10*f5 - 0.1*f6)
}

// The linear function
type Linear struct{}

func (Linear) F(x []float64) float64 {
	return floats.Sum(x)
}

func (Linear) Df(x, grad []float64) {
	for i := range grad {
		grad[i] = 1
	}
}

func TestMinimize(t *testing.T) {
	testMinimize(t, nil)
}

func TestGradientDescent(t *testing.T) {
	testMinimize(t, &GradientDescent{})
}

func TestGradientDescentBacktracking(t *testing.T) {
	testMinimize(t, &GradientDescent{
		LinesearchMethod: &Backtracking{
			FunConst: 0.1,
		},
	})
}

func TestGradientDescentBisection(t *testing.T) {
	testMinimize(t, &GradientDescent{
		LinesearchMethod: &Bisection{},
	})
}

func TestBFGS(t *testing.T) {
	testMinimize(t, &BFGS{})
}

func TestLBFGS(t *testing.T) {
	testMinimize(t, &LBFGS{})
}

func testMinimize(t *testing.T, method Method) {
	// This should be replaced with a more general testing framework with
	// a plugable method

	for i, test := range []struct {
		F Function
		X []float64

		OptVal float64
		OptLoc []float64

		Tol      float64
		Settings *Settings
	}{
		{
			F:      Rosenbrock{2},
			X:      []float64{15, 10},
			OptVal: 0,
			OptLoc: []float64{1, 1},
			Tol:    1e-4,

			Settings: DefaultSettings(),
		},
		{
			F:      Rosenbrock{4},
			X:      []float64{-150, 100, 5, -6},
			OptVal: 0,
			OptLoc: []float64{1, 1, 1, 1},
			Tol:    1e-4,

			Settings: &Settings{
				FunctionAbsTol: math.Inf(-1),
				GradientAbsTol: 1e-12,
			},
		},
		{
			F:      Rosenbrock{2},
			X:      []float64{15, 10},
			OptVal: 0,
			OptLoc: []float64{1, 1},
			Tol:    1e-4,

			Settings: &Settings{
				FunctionAbsTol: math.Inf(-1),
				GradientAbsTol: 1e-12,
			},
		},
		{
			F:      Rosenbrock{2},
			X:      []float64{-1.2, 1},
			OptVal: 0,
			OptLoc: []float64{1, 1},
			Tol:    1e-4,

			Settings: &Settings{
				FunctionAbsTol: math.Inf(-1),
				GradientAbsTol: 1e-3,
			},
		},
		/*
			// TODO: Turn this on when we have an adaptive linsearch method.
			// Gradient descent with backtracking will basically never finish
			{
				F:      Linear{8},
				X:      []float64{9, 8, 7, 6, 5, 4, 3, 2},
				OptVal: negInf,
				OptLoc: []float64{negInf, negInf, negInf, negInf, negInf, negInf, negInf, negInf},

				Settings: &Settings{
					FunctionAbsTol: math.Inf(-1),
				},
			},
		*/
	} {
		test.Settings.Recorder = nil
		result, err := Local(test.F, test.X, test.Settings, method)
		if err != nil {
			t.Errorf("error finding minimum: %v", err.Error())
			continue
		}
		// fmt.Println("%#v\n", result) // for debugging
		// TODO: Better tests
		if math.Abs(result.F-test.OptVal) > test.Tol {
			t.Errorf("Case %v: Minimum not found, exited with status: %v. Want: %v, Got: %v", i, result.Status, test.OptVal, result.F)
			continue
		}
		if result == nil {
			t.Errorf("Case %v: nil result without error", i)
			continue
		}

		// rerun it again to ensure it gets the same answer with the same starting
		// condition
		result2, err2 := Local(test.F, test.X, test.Settings, method)
		if err2 != nil {
			t.Errorf("error finding minimum second time: %v", err2.Error())
			continue
		}
		if result2 == nil {
			t.Errorf("Case %v: nil result without error", i)
			continue
		}
		/*
			// For debugging purposes, can't use DeepEqual naively becaus of NaNs
			// kill the runtime before the check, because those don't need to be equal
			result.Runtime = 0
			result2.Runtime = 0
			if !reflect.DeepEqual(result, result2) {
				t.Error(eqString)
				continue
			}
		*/
	}
}
