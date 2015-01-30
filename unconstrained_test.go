// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package optimize

import (
	"fmt"
	"math"
	"reflect"
	"testing"

	"github.com/gonum/floats"
)

var negInf = math.Inf(-1)

// The Fletcher-Powell helical valley function
// Dim = 3
// X0 = [-1, 0, 0]
// OptX = [1, 0, 0]
// OptF = 0
type HelicalValley struct{}

func (HelicalValley) F(x []float64) float64 {
	θ := 0.5 * math.Atan2(x[1], x[0]) / math.Pi
	r := math.Hypot(x[0], x[1])

	f1 := 10 * (x[2] - 10*θ)
	f2 := 10 * (r - 1)
	f3 := x[2]

	return f1*f1 + f2*f2 + f3*f3
}

func (HelicalValley) Df(x, grad []float64) {
	θ := 0.5 * math.Atan2(x[1], x[0]) / math.Pi
	r := math.Hypot(x[0], x[1])
	s := x[2] - 10*θ
	t := 5 * s / r / r / math.Pi

	grad[0] = 200 * (x[0] - x[0]/r + x[1]*t)
	grad[1] = 200 * (x[1] - x[1]/r - x[0]*t)
	grad[2] = 2 * (x[2] + 100*s)
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
		sum += f * f
	}
	return sum
}

func (BiggsEXP2) Df(x, grad []float64) {
	for i := range grad {
		grad[i] = 0
	}
	for i := 1; i <= 10; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z)
		f := math.Exp(-x[0]*z) - 5*math.Exp(-x[1]*z) - y

		dfdx0 := -z * math.Exp(-x[0]*z)
		dfdx1 := 5 * z * math.Exp(-x[1]*z)

		grad[0] += 2 * f * dfdx0
		grad[1] += 2 * f * dfdx1
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
		sum += f * f
	}
	return sum
}

func (BiggsEXP3) Df(x, grad []float64) {
	for i := range grad {
		grad[i] = 0
	}
	for i := 1; i <= 10; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z)
		f := math.Exp(-x[0]*z) - x[2]*math.Exp(-x[1]*z) - y

		dfdx0 := -z * math.Exp(-x[0]*z)
		dfdx1 := x[2] * z * math.Exp(-x[1]*z)
		dfdx2 := -math.Exp(-x[1] * z)

		grad[0] += 2 * f * dfdx0
		grad[1] += 2 * f * dfdx1
		grad[2] += 2 * f * dfdx2
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
		sum += f * f
	}
	return sum
}

func (BiggsEXP4) Df(x, grad []float64) {
	for i := range grad {
		grad[i] = 0
	}
	for i := 1; i <= 10; i++ {
		z := float64(i) / 10
		y := math.Exp(-z) - 5*math.Exp(-10*z)
		f := x[2]*math.Exp(-x[0]*z) - x[3]*math.Exp(-x[1]*z) - y

		dfdx0 := -z * x[2] * math.Exp(-x[0]*z)
		dfdx1 := z * x[3] * math.Exp(-x[1]*z)
		dfdx2 := math.Exp(-x[0] * z)
		dfdx3 := -math.Exp(-x[1] * z)

		grad[0] += 2 * f * dfdx0
		grad[1] += 2 * f * dfdx1
		grad[2] += 2 * f * dfdx2
		grad[3] += 2 * f * dfdx3
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
		sum += f * f
	}
	return sum
}

func (BiggsEXP5) Df(x, grad []float64) {
	for i := range grad {
		grad[i] = 0
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

		grad[0] += 2 * f * dfdx0
		grad[1] += 2 * f * dfdx1
		grad[2] += 2 * f * dfdx2
		grad[3] += 2 * f * dfdx3
		grad[4] += 2 * f * dfdx4
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
		sum += f * f
	}
	return sum
}

func (BiggsEXP6) Df(x, grad []float64) {
	for i := range grad {
		grad[i] = 0
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

		grad[0] += 2 * f * dfdx0
		grad[1] += 2 * f * dfdx1
		grad[2] += 2 * f * dfdx2
		grad[3] += 2 * f * dfdx3
		grad[4] += 2 * f * dfdx4
		grad[5] += 2 * f * dfdx5
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
		c := 0.5 * float64(8-i)
		b := c - x[2]
		d := b * b
		e := math.Exp(-0.5 * x[1] * d)
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
		c := 0.5 * float64(8-i)
		b := c - x[2]
		d := b * b
		e := math.Exp(-0.5 * x[1] * d)
		f := x[0]*e - g.y(i)

		grad[0] += 2 * f * e
		grad[1] -= f * e * d * x[0]
		grad[2] += 2 * f * e * x[0] * x[1] * b
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
	return f1*f1 + f2*f2
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
		sum += f * f
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
// X0 = [..., (n-i)/n, ...], i=1,...,n,
// OptX = [1, ..., 1]
// OptF = 0
type VariablyDimensioned struct{}

func (v VariablyDimensioned) F(x []float64) (sum float64) {
	for i := 0; i < len(x); i++ {
		t := x[i] - 1
		sum += t * t
	}
	s := 0.0
	for i := 0; i < len(x); i++ {
		s += float64(i+1) * (x[i] - 1)
	}
	s *= s
	sum += s
	s *= s
	sum += s
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
		grad[i] = 2 * ((x[i] - 1) + s*float64(i+1)*(1+2*s*s))
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
		s2 *= s2
		t := s1 - s2 - 1
		sum += t * t
	}
	t := x[0] * x[0]
	sum += t
	t = x[1] - t - 1
	sum += t * t

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

		t := s1 - s2*s2 - 1
		for j := 0; j < len(x); j++ {
			grad[j] += 2 * t * math.Pow(c, float64(j-1)) * (float64(j) - 2*s2*c)
		}
	}
	t := x[1] - x[0]*x[0] - 1
	grad[0] += 2 * (1 - 2*t) * x[0]
	grad[1] += 2 * t
}

// Penalty function #1
// J. More, B.S. Garbow, K.E. Hillstrom, Testing unconstrained optimization software.
// ACM Trans.Math. Softw. 7 (1981), 17-41.
// Dim = n
// X0 = [1, ..., n]
// For Dim = 4:
// OptF = 2.2499775...e-5
// For Dim = 10:
// OptF = 7.0876515...e-5
type Penalty1 struct{}

func (Penalty1) F(x []float64) (sum float64) {
	for i := 0; i < len(x); i++ {
		t := x[i] - 1
		sum += t * t
	}
	sum *= 1e-5

	s := 0.0
	for i := 0; i < len(x); i++ {
		s += x[i] * x[i]
	}
	sum += (s - 0.25) * (s - 0.25)

	return sum
}

func (Penalty1) Df(x, grad []float64) {
	s := 0.0
	for i := 0; i < len(x); i++ {
		s += x[i] * x[i]
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

	s := -1.0
	for i := 0; i < dim; i++ {
		s += float64(dim-i) * x[i] * x[i]
	}
	for i := 1; i < dim; i++ {
		yi := math.Exp(float64(i+1)/10) + math.Exp(float64(i)/10)
		f := math.Exp(x[i]/10) + math.Exp(x[i-1]/10) - yi
		sum += f * f
	}
	for i := 1; i < dim; i++ {
		f := math.Exp(x[i]/10) - math.Exp(-1.0/10)
		sum += f * f
	}
	sum *= 1e-5

	sum += (x[0] - 0.2) * (x[0] - 0.2)
	sum += s * s

	return sum
}

func (Penalty2) Df(x, grad []float64) {
	dim := len(x)

	s := 0.0
	for i := 0; i < dim; i++ {
		s += float64(dim-i) * x[i] * x[i]
	}
	s--

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
// OptX = [1e6, 2e-6]
// OptF = 0
type Brown struct{}

func (Brown) F(x []float64) float64 {
	f1 := x[0] - 1e6
	f2 := x[1] - 2e-6
	f3 := x[0]*x[1] - 2
	return f1*f1 + f2*f2 + f3*f3
}

func (Brown) Df(x, grad []float64) {
	f1 := x[0] - 1e6
	f2 := x[1] - 2e-6
	f3 := x[0]*x[1] - 2

	grad[0] = 2*f1 + 2*f3*x[1]
	grad[1] = 2*f2 + 2*f3*x[0]
}

// The Brown and Dennis function
// J. More, B.S. Garbow, K.E. Hillstrom, Testing unconstrained optimization software.
// ACM Trans.Math. Softw. 7 (1981), 17-41.
// Dim = 4
// X0 = [25, 5, -5, -1]
// OptF = 85822.20162635628
type BrownDennis struct{}

func (BrownDennis) F(x []float64) (sum float64) {
	for i := 0; i < 20; i++ {
		c := float64(i+1) / 5
		d1 := x[0] + c*x[1] - math.Exp(c)
		d2 := x[2] + x[3]*math.Sin(c) - math.Cos(c)
		f := d1*d1 + d2*d2
		sum += f * f
	}
	return sum
}

func (BrownDennis) Df(x, grad []float64) {
	for i := range grad {
		grad[i] = 0
	}

	for i := 0; i < 20; i++ {
		c := float64(i+1) / 5
		d1 := x[0] + c*x[1] - math.Exp(c)
		d2 := x[2] + x[3]*math.Sin(c) - math.Cos(c)
		f := d1*d1 + d2*d2
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

		sum += f * f
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

		grad[0] += 2 * f * math.Exp(-e) * math.Pow(d, x[2]) / x[0] / x[0]
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
		sum += f * f
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

// The Extended Rosenbrock function implemented with F+Df pair of methods.
// Very difficult to minimize if the starting point is far from the minimum.
// Dim = n
// X0 = [-1.2, 1] for Dim = 2
// OptF = 0 (global)
// OptX = [1, ..., 1]
type Rosenbrock struct{}

func (Rosenbrock) F(x []float64) (sum float64) {
	for i := 0; i < len(x)-1; i++ {
		a := 1 - x[i]
		b := x[i+1] - x[i]*x[i]
		sum += a*a + 100*b*b
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
		grad[i] -= 400 * (x[i+1] - x[i]*x[i]) * x[i]
	}
	for i := 1; i < dim; i++ {
		grad[i] += 200 * (x[i] - x[i-1]*x[i-1])
	}
}

// The Extended Rosenbrock function implemented with F+FDf pair of methods.
type RosenbrockFDf struct{}

func (RosenbrockFDf) F(x []float64) (sum float64) {
	for i := 0; i < len(x)-1; i++ {
		a := 1 - x[i]
		b := x[i+1] - x[i]*x[i]
		sum += a*a + 100*b*b
	}
	return sum
}

func (f RosenbrockFDf) FDf(x, grad []float64) float64 {
	dim := len(x)
	for i := range grad {
		grad[i] = 0
	}

	for i := 0; i < dim-1; i++ {
		grad[i] -= 2 * (1 - x[i])
		grad[i] -= 400 * (x[i+1] - x[i]*x[i]) * x[i]
	}
	for i := 1; i < dim; i++ {
		grad[i] += 200 * (x[i] - x[i-1]*x[i-1])
	}

	return f.F(x)
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
		t := x[i+1] - 2*x[i+2]
		f3 := t * t
		t = x[i] - x[i+3]
		f4 := t * t
		sum += f1*f1 + 5*f2*f2 + f3*f3 + 10*f4*f4
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
		t1 := x[i+1] - 2*x[i+2]
		f3 := t1 * t1
		t2 := x[i] - x[i+3]
		f4 := t2 * t2

		grad[i] = 2*f1 + 40*f4*t2
		grad[i+1] = 20*f1 + 4*f3*t1
		grad[i+2] = 10*f2 - 8*f3*t1
		grad[i+3] = -10*f2 - 40*f4*t2
	}
}

// The Beale function
// Dim = 2
// X0 = [1, 1]
// OptF = 0
// OptX = [3, 0.5]
type Beale struct{}

func (Beale) F(x []float64) float64 {
	f1 := 1.5 - x[0]*(1-x[1])
	f2 := 2.25 - x[0]*(1-x[1]*x[1])
	f3 := 2.625 - x[0]*(1-x[1]*x[1]*x[1])
	return f1*f1 + f2*f2 + f3*f3
}

func (Beale) Df(x, grad []float64) {
	t1 := 1 - x[1]
	t2 := 1 - x[1]*x[1]
	t3 := 1 - x[1]*x[1]*x[1]

	f1 := 1.5 - x[0]*t1
	f2 := 2.25 - x[0]*t2
	f3 := 2.625 - x[0]*t3

	grad[0] = -2 * (f1*t1 + f2*t2 + f3*t3)
	grad[1] = 2 * x[0] * (f1 + 2*f2*x[1] + 3*f3*x[1]*x[1])
}

// The Wood function
// Dim = 4
// X0 = [-3, -1, -3, -1]
// OptF = 0
// OptX = [1, 1, 1, 1]
type Wood struct{}

func (Wood) F(x []float64) (sum float64) {
	f1 := x[1] - x[0]*x[0]
	f2 := 1 - x[0]
	f3 := x[3] - x[2]*x[2]
	f4 := 1 - x[2]
	f5 := x[1] + x[3] - 2
	f6 := x[2] - x[3]

	sum = 100*f1*f1 + f2*f2 + 90*f3*f3
	sum += f4*f4 + 10*f5*f5 + 0.1*f6*f6
	return sum
}

func (Wood) Df(x, grad []float64) {
	f1 := x[1] - x[0]*x[0]
	f2 := 1 - x[0]
	f3 := x[3] - x[2]*x[2]
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

type UnconstrainedTest struct {
	// f is the function that is being minimized.
	f Function
	// x is the initial guess.
	x []float64
	// optVal is the value of f at a minimum.
	optVal float64

	// optLoc is the location of the minimum. If it is not known, optLoc is nil.
	optLoc []float64
	// gradTol is the absolute gradient tolerance for the test. If gradTol == 0,
	// the default tolerance 1e-12 will be used.
	gradTol float64
	// tol is the tolerance for checking the accuracy of result.F. If tol == 0,
	// the default tolerance 1e-5 will be used.
	tol float64
}

func (t UnconstrainedTest) String() string {
	dim := len(t.x)
	if dim <= 10 {
		// Print the initial and optimal X only for small-dimensional problems.
		return fmt.Sprintf("F: %v\nDim: %v\nGradientAbsTol: %v\nInitial X: %v\nWant X: %v\nWant F(X): %v",
			reflect.TypeOf(t.f), dim, t.gradTol, t.x, t.optLoc, t.optVal)
	}
	return fmt.Sprintf("F: %v\nDim: %v\nGradientAbsTol: %v\nWant F(X): %v",
		reflect.TypeOf(t.f), dim, t.gradTol, t.optVal)
}

var gradientDescentTests = []UnconstrainedTest{
	{
		f:       Rosenbrock{},
		x:       []float64{-1.2, 1},
		optVal:  0,
		optLoc:  []float64{1, 1},
		gradTol: defaultGradientAbsTol,
	},
	{
		f:      Rosenbrock{},
		x:      []float64{-1.2, 1},
		optVal: 0,
		optLoc: []float64{1, 1},
	},
	{
		f:       Rosenbrock{},
		x:       []float64{-120, 100, 50},
		optVal:  0,
		optLoc:  []float64{1, 1, 1},
		gradTol: defaultGradientAbsTol,
	},
	{
		f:      RosenbrockFDf{},
		x:      []float64{-1.2, 1},
		optVal: 0,
		optLoc: []float64{1, 1},
	},
	{
		f:       RosenbrockFDf{},
		x:       []float64{-120, 100, 50},
		optVal:  0,
		optLoc:  []float64{1, 1, 1},
		gradTol: defaultGradientAbsTol,
	},
	{
		f:      HelicalValley{},
		x:      []float64{-1, 0, 0},
		optVal: 0,
		optLoc: []float64{1, 0, 0},
	},
	{
		f:      BiggsEXP4{},
		x:      []float64{1, 2, 1, 1},
		optVal: 0,
		optLoc: []float64{1, 10, 1, 5},
	},
}

var cgTests = []UnconstrainedTest{
	{
		f:      Rosenbrock{},
		x:      []float64{-1200000, 1000000},
		optVal: 0,
		optLoc: []float64{1, 1},
	},
	{
		f:       BiggsEXP6{},
		x:       []float64{1, 2, 1, 1, 1, 1},
		optVal:  0.005655649925,
		gradTol: 1e-6,
	},
	{
		f:       Gaussian{},
		x:       []float64{0.4, 1, 0},
		optVal:  1.12793e-8,
		optLoc:  []float64{0.3989561, 1.0000191, 0},
		gradTol: 1e-9,
	},
	generateVariablyDimensioned(1000, 1e-9),
	generateVariablyDimensioned(10000, 1e-8),
	{
		f:       Watson{},
		x:       []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		optVal:  0,
		optLoc:  []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		gradTol: 1e-6,
	},
	{
		f:       Penalty1{},
		x:       []float64{1, 2, 3, 4},
		optVal:  2.2499775e-5,
		gradTol: 1e-10,
	},
	{
		f:       Penalty1{},
		x:       []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		optVal:  7.0876515e-5,
		gradTol: 1e-8,
	},
	{
		f:       Penalty2{},
		x:       []float64{0.5, 0.5, 0.5, 0.5},
		optVal:  9.37629e-6,
		gradTol: 1e-6,
	},
	{
		f:       Penalty2{},
		x:       []float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5},
		optVal:  2.93660e-4,
		gradTol: 1e-6,
	},
	{
		f:       ExtendedPowell{},
		x:       []float64{3, -1, 0, 1},
		optVal:  0,
		optLoc:  []float64{0, 0, 0, 0},
		gradTol: 1e-6,
	},
	{
		f:      Beale{},
		x:      []float64{1, 1},
		optVal: 0,
		optLoc: []float64{3, 0.5},
	},
	{
		f:      Wood{},
		x:      []float64{-3, -1, -3, -1},
		optVal: 0,
		optLoc: []float64{1, 1, 1, 1},
	},
	{
		f:      Trigonometric{},
		x:      []float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1},
		optVal: 0,
		optLoc: []float64{
			0.0429645656464827, 0.043976286943042, 0.0450933996565844, 0.0463389157816542, 0.0477443839560646,
			0.0493547352444078, 0.0512373449259557, 0.195209463715277, 0.164977664720328, 0.0601485854398078,
		},
		gradTol: 1e-8,
	},
}

var newtonTests = []UnconstrainedTest{
	{
		f:      Rosenbrock{},
		x:      []float64{10, 10, 10, 10},
		optVal: 0,
		optLoc: []float64{1, 1, 1, 1},
	},
	{
		f:      Rosenbrock{},
		x:      []float64{-12000, 10000},
		optVal: 0,
		optLoc: []float64{1, 1},
	},
	{
		f:       Gaussian{},
		x:       []float64{0.4, 1, 0},
		optVal:  1.12793e-8,
		optLoc:  []float64{0.3989561, 1.0000191, 0},
		gradTol: 1e-11,
	},
	{
		f:      Powell{},
		x:      []float64{0, 1},
		optVal: 0,
		optLoc: []float64{1.09815933e-5, 9.10614674},
	},
	{
		f:      Box{},
		x:      []float64{0, 10, 20},
		optVal: 0,
		optLoc: []float64{1, 10, 1},
	},
	generateVariablyDimensioned(10, 0),
	generateVariablyDimensioned(100, 0),
	{
		f:       Watson{},
		x:       []float64{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		optVal:  0,
		optLoc:  []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
		gradTol: 1e-8,
	},
	{
		f:      Penalty1{},
		x:      []float64{1, 2, 3, 4},
		optVal: 2.2499775e-5,
	},
	{
		f:       Penalty2{},
		x:       []float64{0.5, 0.5, 0.5, 0.5},
		optVal:  9.37629e-6,
		gradTol: 1e-11,
	},
	{
		f:      Brown{},
		x:      []float64{1, 1},
		optVal: 0,
		optLoc: []float64{1e6, 2e-6},
	},
	{
		f:      GulfRD{},
		x:      []float64{5, 2.5, 0.15},
		optVal: 0,
		optLoc: []float64{50, 25, 1.5},
	},
	{
		f:      ExtendedPowell{},
		x:      []float64{3, -1, 0, 1},
		optVal: 0,
		optLoc: []float64{0, 0, 0, 0},
	},
	{
		f:      Beale{},
		x:      []float64{1, 1},
		optVal: 0,
		optLoc: []float64{3, 0.5},
	},
	{
		f:      Wood{},
		x:      []float64{-3, -1, -3, -1},
		optVal: 0,
		optLoc: []float64{1, 1, 1, 1},
	},
}

var bfgsTests = []UnconstrainedTest{
	{
		f:       BiggsEXP6{},
		x:       []float64{1, 2, 1, 1, 1, 1},
		optVal:  0.005655649925,
		gradTol: 1e-10,
	},
	{
		f:       Penalty1{},
		x:       []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		optVal:  7.0876515e-5,
		gradTol: 1e-10,
	},
	{
		f:       Penalty2{},
		x:       []float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5},
		optVal:  2.93660e-4,
		gradTol: 1e-10,
	},
	{
		f:       BrownDennis{},
		x:       []float64{25, 5, -5, -1},
		optVal:  85822.20162635628,
		gradTol: 1e-5,
	},
	{
		f:      Trigonometric{},
		x:      []float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1},
		optVal: 0,
		optLoc: []float64{
			0.0429645656464827, 0.043976286943042, 0.0450933996565844, 0.0463389157816542, 0.0477443839560646,
			0.0493547352444078, 0.0512373449259557, 0.195209463715277, 0.164977664720328, 0.0601485854398078,
		},
		gradTol: 1e-11,
	},
}

var lbfgsTests = []UnconstrainedTest{
	{
		f:      Rosenbrock{},
		x:      []float64{-1200000, 1000000},
		optVal: 0,
		optLoc: []float64{1, 1},
	},
	{
		f:       BiggsEXP6{},
		x:       []float64{1, 2, 1, 1, 1, 1},
		optVal:  0.005655649925,
		gradTol: 1e-8,
	},
	generateVariablyDimensioned(1000, 1e-10),
	generateVariablyDimensioned(10000, 1e-8),
	{
		f:       Penalty1{},
		x:       []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		optVal:  7.0876515e-5,
		gradTol: 1e-11,
	},
	{
		f:       Penalty2{},
		x:       []float64{0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5, 0.5},
		optVal:  2.93660e-4,
		gradTol: 1e-9,
	},
	{
		f:       BrownDennis{},
		x:       []float64{25, 5, -5, -1},
		optVal:  85822.20162635628,
		gradTol: 1e-4, // This is the best LBFGS can currently do.
	},
	{
		f:      Trigonometric{},
		x:      []float64{0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1, 0.1},
		optVal: 0,
		optLoc: []float64{
			0.0429645656464827, 0.043976286943042, 0.0450933996565844, 0.0463389157816542, 0.0477443839560646,
			0.0493547352444078, 0.0512373449259557, 0.195209463715277, 0.164977664720328, 0.0601485854398078,
		},
		gradTol: 1e-9,
	},
}

func generateVariablyDimensioned(dim int, gradTol float64) UnconstrainedTest {
	x := make([]float64, dim)
	for i := range x {
		x[i] = float64(dim-i-1) / float64(dim)
	}
	optLoc := make([]float64, dim)
	for i := range optLoc {
		optLoc[i] = 1
	}
	return UnconstrainedTest{
		f:       VariablyDimensioned{},
		x:       x,
		optVal:  0,
		optLoc:  optLoc,
		gradTol: gradTol,
	}
}

func TestMinimize(t *testing.T) {
	// TODO: When method is nil, Local chooses the method automatically. At
	// present, it always chooses BFGS (or panics if the function does not
	// implement Df() or FDf()). For now, run this test with the simplest set
	// of problems and revisit this later when more methods are added.
	testMinimize(t, gradientDescentTests, nil)
}

func TestGradientDescent(t *testing.T) {
	testMinimize(t, gradientDescentTests, &GradientDescent{})
}

func TestGradientDescentBacktracking(t *testing.T) {
	testMinimize(t, gradientDescentTests, &GradientDescent{
		LinesearchMethod: &Backtracking{
			FunConst: 0.1,
		},
	})
}

func TestGradientDescentBisection(t *testing.T) {
	testMinimize(t, gradientDescentTests, &GradientDescent{
		LinesearchMethod: &Bisection{},
	})
}

func TestCG(t *testing.T) {
	var tests []UnconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, cgTests...)
	testMinimize(t, tests, &CG{})
}

func TestFletcherReevesQuadStep(t *testing.T) {
	var tests []UnconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, cgTests...)
	testMinimize(t, tests, &CG{
		Variant:     &FletcherReeves{},
		InitialStep: &QuadraticStepSize{},
	})
}

func TestFletcherReevesFirstOrderStep(t *testing.T) {
	var tests []UnconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, cgTests...)
	testMinimize(t, tests, &CG{
		Variant:     &FletcherReeves{},
		InitialStep: &FirstOrderStepSize{},
	})
}

func TestHestenesStiefelQuadStep(t *testing.T) {
	var tests []UnconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, cgTests...)
	testMinimize(t, tests, &CG{
		Variant:     &HestenesStiefel{},
		InitialStep: &QuadraticStepSize{},
	})
}

func TestHestenesStiefelFirstOrderStep(t *testing.T) {
	var tests []UnconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, cgTests...)
	testMinimize(t, tests, &CG{
		Variant:     &HestenesStiefel{},
		InitialStep: &FirstOrderStepSize{},
	})
}

func TestPolakRibiereQuadStep(t *testing.T) {
	var tests []UnconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, cgTests...)
	testMinimize(t, tests, &CG{
		Variant:     &PolakRibierePolyak{},
		InitialStep: &QuadraticStepSize{},
	})
}

func TestPolakRibiereFirstOrderStep(t *testing.T) {
	var tests []UnconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, cgTests...)
	testMinimize(t, tests, &CG{
		Variant:     &PolakRibierePolyak{},
		InitialStep: &FirstOrderStepSize{},
	})
}

func TestHagerZhangQuadStep(t *testing.T) {
	var tests []UnconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, cgTests...)
	testMinimize(t, tests, &CG{
		Variant:     &HagerZhang{},
		InitialStep: &QuadraticStepSize{},
	})
}

func TestHagerZhangFirstOrderStep(t *testing.T) {
	var tests []UnconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, cgTests...)
	testMinimize(t, tests, &CG{
		Variant:     &HagerZhang{},
		InitialStep: &FirstOrderStepSize{},
	})
}

func TestBFGS(t *testing.T) {
	var tests []UnconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, newtonTests...)
	tests = append(tests, bfgsTests...)
	testMinimize(t, tests, &BFGS{})
}

func TestLBFGS(t *testing.T) {
	var tests []UnconstrainedTest
	tests = append(tests, gradientDescentTests...)
	tests = append(tests, newtonTests...)
	tests = append(tests, lbfgsTests...)
	testMinimize(t, tests, &LBFGS{})
}

func testMinimize(t *testing.T, tests []UnconstrainedTest, method Method) {
	for _, test := range tests {
		settings := &Settings{
			FunctionAbsTol: math.Inf(-1),
		}
		if test.gradTol == 0 {
			test.gradTol = 1e-12
		}
		settings.GradientAbsTol = test.gradTol
		if test.tol == 0 {
			test.tol = 1e-5
		}

		result, err := Local(test.f, test.x, settings, method)
		if err != nil {
			t.Errorf("error finding minimum (%v) for:\n%v", err, test)
			continue
		}

		if result == nil {
			t.Errorf("nil result without error for:\n%v", test)
			continue
		}

		// Check that the optimum function value is as expected.
		if math.Abs(result.F-test.optVal) > test.tol {
			t.Errorf("Minimum not found, exited with status: %v. Want: %v, Got: %v for:\n%v",
				result.Status, test.optVal, result.F, test)
			continue
		}

		funcs, funcInfo := getFunctionInfo(test.f)

		// Evaluate the norm of the gradient at the found optimum location.
		var optF, optNorm float64
		if funcInfo.IsFunctionGradient {
			g := make([]float64, len(test.x))
			optF = funcs.gradFunc.FDf(result.X, g)
			optNorm = floats.Norm(g, math.Inf(1))
		} else {
			optF = funcs.function.F(result.X)
			if funcInfo.IsGradient {
				g := make([]float64, len(test.x))
				funcs.gradient.Df(result.X, g)
				optNorm = floats.Norm(g, math.Inf(1))
			}
		}

		// Check that the function value at the found optimum location is
		// equal to result.F
		if optF != result.F {
			t.Errorf("Function value at the optimum location %v not equal to the returned value %v for:\n%v",
				optF, result.F, test)
		}

		// Check that the norm of the gradient at the found optimum location is
		// smaller than the tolerance.
		if optNorm >= settings.GradientAbsTol {
			t.Errorf("Norm of the gradient at the optimum location %v not smaller than tolerance %v for:\n%v",
				optNorm, settings.GradientAbsTol, test)
		}

		// We are going to restart the solution using a fixed starting gradient
		// and value, so evaluate them.
		settings.UseInitialData = true
		if funcInfo.IsFunctionGradient {
			settings.InitialGradient = resize(settings.InitialGradient, len(test.x))
			settings.InitialFunctionValue = funcs.gradFunc.FDf(test.x, settings.InitialGradient)
		} else {
			settings.InitialFunctionValue = funcs.function.F(test.x)
			if funcInfo.IsGradient {
				settings.InitialGradient = resize(settings.InitialGradient, len(test.x))
				funcs.gradient.Df(test.x, settings.InitialGradient)
			}
		}

		// Rerun the test again to make sure that it gets the same answer with
		// the same starting condition. Moreover, we are using the initial data
		// in settings.InitialFunctionValue and settings.InitialGradient.
		result2, err2 := Local(test.f, test.x, settings, method)
		if err2 != nil {
			t.Errorf("error finding minimum second time (%v) for:\n%v", err2, test)
			continue
		}

		if result2 == nil {
			t.Errorf("second time nil result without error for:\n%v", test)
			continue
		}

		// At the moment all the optimizers are deterministic, so check that we
		// get _exactly_ the same answer second time as well.
		if result.F != result2.F {
			t.Errorf("Different minimum second time. First: %v, Second: %v for:\n%v",
				result.F, result2.F, test)
		}

		// Check that providing initial data reduces the number of function
		// and/or gradient calls exactly by one.
		if funcInfo.IsFunctionGradient {
			if result.FunctionGradientEvals != result2.FunctionGradientEvals+1 {
				t.Errorf("Providing initial data does not reduce the number of function/gradient calls for:\n%v", test)
				continue
			}
		} else {
			if result.FunctionEvals != result2.FunctionEvals+1 {
				t.Errorf("Providing initial data does not reduce the number of functions calls for:\n%v", test)
				continue
			}
			if funcInfo.IsGradient {
				if result.GradientEvals != result2.GradientEvals+1 {
					t.Errorf("Providing initial data does not reduce the number of gradient calls for:\n%v", test)
					continue
				}
			}
		}

		// TODO: Enable this test once the optimizers reliably handle
		// minimization from a (nearly) optimal location.
		// if test.optLoc != nil {
		// 	settings.UseInitialData = false
		// 	// Try starting the optimizer from a (nearly) optimum location given by test.optLoc
		// 	_, err3 := Local(test.f, test.optLoc, settings, method)
		// 	if err3 != nil {
		// 		t.Errorf("error finding minimum from a (nearly) optimum location (%v) for:\n%v", err3, test)
		// 	}
		// }
	}
}
