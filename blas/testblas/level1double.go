// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testblas

import (
	"fmt"
	"math"
	"testing"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/floats/scalar"
)

type DoubleOneVectorCase struct {
	Name       string
	X          []float64
	Incx       int
	N          int
	Panic      bool
	Dasum      float64
	Dnrm2      float64
	Idamax     int
	DscalCases []DScalCase
}

type DScalCase struct {
	Alpha float64
	Ans   []float64
	Name  string
}

var DoubleOneVectorCases = []DoubleOneVectorCase{
	{
		Name:   "AllPositive",
		X:      []float64{6, 5, 4, 2, 6},
		Incx:   1,
		N:      5,
		Panic:  false,
		Dasum:  23,
		Dnrm2:  10.81665382639196787935766380241148783875388972153573863813135,
		Idamax: 0,
		DscalCases: []DScalCase{
			{
				Alpha: 0,
				Ans:   []float64{0, 0, 0, 0, 0},
			},
			{
				Alpha: 1,
				Ans:   []float64{6, 5, 4, 2, 6},
			},
			{
				Alpha: -2,
				Ans:   []float64{-12, -10, -8, -4, -12},
			},
		},
	},
	{
		Name:   "LeadingZero",
		X:      []float64{0, 1},
		Incx:   1,
		N:      2,
		Panic:  false,
		Dasum:  1,
		Dnrm2:  1,
		Idamax: 1,
		DscalCases: []DScalCase{
			{
				Alpha: 0,
				Ans:   []float64{0, 0},
			},
			{
				Alpha: 1,
				Ans:   []float64{0, 1},
			},
			{
				Alpha: -2,
				Ans:   []float64{0, -2},
			},
		},
	},
	{
		Name:   "MaxInMiddle",
		X:      []float64{6, 5, 9, 0, 6},
		Incx:   1,
		N:      5,
		Panic:  false,
		Dasum:  26,
		Dnrm2:  13.34166406412633371248943627250846646911846482744007727141318,
		Idamax: 2,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{-12, -10, -18, 0, -12},
			},
		},
	},
	{
		Name:   "MaxAtEnd",
		X:      []float64{6, 5, -9, 0, 10},
		Incx:   1,
		N:      5,
		Panic:  false,
		Dasum:  30,
		Dnrm2:  15.55634918610404553681857596630667886426639062914642880494347,
		Idamax: 4,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{-12, -10, 18, 0, -20},
			},
		},
	},
	{
		Name:   "AllNegative",
		X:      []float64{-6, -5, -4, -2, -6},
		Incx:   1,
		N:      5,
		Panic:  false,
		Dasum:  23,
		Dnrm2:  10.81665382639196787935766380241148783875388972153573863813135,
		Idamax: 0,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{12, 10, 8, 4, 12},
			},
		},
	},
	{
		Name:   "AllMixed",
		X:      []float64{-6, 5, 4, -2, -6},
		Incx:   1,
		N:      5,
		Panic:  false,
		Dasum:  23,
		Dnrm2:  10.81665382639196787935766380241148783875388972153573863813135,
		Idamax: 0,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{12, -10, -8, 4, 12},
			},
		},
	},
	{
		Name:   "ZeroN",
		X:      []float64{-6, 5, 4, -2, -6},
		Incx:   1,
		N:      0,
		Panic:  false,
		Dasum:  0,
		Dnrm2:  0,
		Idamax: -1,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{-6, 5, 4, -2, -6},
			},
		},
	},
	{
		Name:   "OneN",
		X:      []float64{-6, 5, 4, -2, -6},
		Incx:   1,
		N:      1,
		Panic:  false,
		Dasum:  6,
		Dnrm2:  6,
		Idamax: 0,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{12, 5, 4, -2, -6},
			},
		},
	},
	{
		Name:   "PositiveExactInc",
		X:      []float64{-6, 5, 10, -2, -5},
		Incx:   2,
		N:      3,
		Panic:  false,
		Dasum:  21,
		Dnrm2:  12.68857754044952038019377274608948979173952662752515253090272,
		Idamax: 1,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{12, 5, -20, -2, 10},
			},
		},
	},
	{
		Name:   "PositiveOffInc",
		X:      []float64{-6, 5, 4, -2, -6, 8, 10, 11},
		Incx:   3,
		N:      3,
		Panic:  false,
		Dasum:  18,
		Dnrm2:  11.83215956619923208513465658312323409683100246158868064575943,
		Idamax: 2,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{12, 5, 4, 4, -6, 8, -20, 11},
			},
		},
	},
	{
		Name:   "PositiveShortInc",
		X:      []float64{-6, 5, 4, -2, -6, 8, 10, 11},
		Incx:   3,
		N:      2,
		Panic:  false,
		Dasum:  8,
		Dnrm2:  6.324555320336758663997787088865437067439110278650433653715009,
		Idamax: 0,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{12, 5, 4, 4, -6, 8, 10, 11},
			},
		},
	},
	{
		Name:   "NegativeInc",
		X:      []float64{-6, 5, 4, -2, -6},
		Incx:   -1,
		N:      5,
		Panic:  false,
		Dasum:  0,
		Dnrm2:  0,
		Idamax: -1,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{-6, 5, 4, -2, -6},
			},
		},
	},
	{
		Name:   "NegativeExactInc",
		X:      []float64{-6, 5, 4, -2, -6},
		Incx:   -2,
		N:      3,
		Panic:  false,
		Dasum:  0,
		Dnrm2:  0,
		Idamax: -1,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{-6, 5, 4, -2, -6},
			},
		},
	},
	{
		Name:   "NegativeOffInc",
		X:      []float64{-6, 5, 4, -2, -6, 8, 10, 11},
		Incx:   -3,
		N:      2,
		Panic:  false,
		Dasum:  0,
		Dnrm2:  0,
		Idamax: -1,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{-6, 5, 4, -2, -6, 8, 10, 11},
			},
		},
	},
	{
		Name:   "NegativeShortInc",
		X:      []float64{-6, 5, 4, -2, -6, 8, 10, 11},
		Incx:   -3,
		N:      2,
		Panic:  false,
		Dasum:  0,
		Dnrm2:  0,
		Idamax: -1,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{-6, 5, 4, -2, -6, 8, 10, 11},
			},
		},
	},
	{
		Name:  "NegativeN",
		X:     []float64{-6, 5, 4, -2, -6},
		Incx:  2,
		N:     -5,
		Panic: true,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{-6, 5, 4, -2, -6},
			},
		},
	},
	{
		Name:  "ZeroInc",
		X:     []float64{-6, 5, 4, -2, -6},
		Incx:  0,
		N:     5,
		Panic: true,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{-6, 5, 4, -2, -6},
			},
		},
	},
	{
		Name:  "OutOfBounds",
		X:     []float64{-6, 5, 4, -2, -6},
		Incx:  2,
		N:     6,
		Panic: true,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{-6, 5, 4, -2, -6},
			},
		},
	},
	{
		Name:   "NegativeOutOfBounds",
		X:      []float64{-6, 5, 4, -2, -6},
		Incx:   -2,
		N:      6,
		Panic:  false,
		Dasum:  0,
		Dnrm2:  0,
		Idamax: -1,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{-6, 5, 4, -2, -6},
			},
		},
	},
	{
		Name:   "NaN",
		X:      []float64{math.NaN(), 2.0},
		Incx:   1,
		N:      2,
		Panic:  false,
		Dasum:  math.NaN(),
		Dnrm2:  math.NaN(),
		Idamax: 0,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{math.NaN(), -4.0},
			},
			{
				Alpha: 0,
				Ans:   []float64{0, 0},
			},
		},
	},
	{
		Name:   "NaNInc",
		X:      []float64{math.NaN(), math.NaN(), 2.0},
		Incx:   2,
		N:      2,
		Panic:  false,
		Dasum:  math.NaN(),
		Dnrm2:  math.NaN(),
		Idamax: 0,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{math.NaN(), math.NaN(), -4.0},
			},
			{
				Alpha: 0,
				Ans:   []float64{0, math.NaN(), 0},
			},
		},
	},
	{
		Name:   "Empty",
		X:      []float64{},
		Incx:   1,
		N:      0,
		Panic:  false,
		Dasum:  0,
		Dnrm2:  0,
		Idamax: -1,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{},
			},
			{
				Alpha: 0,
				Ans:   []float64{},
			},
		},
	},
	{
		Name:   "EmptyZeroInc",
		X:      []float64{},
		Incx:   0,
		N:      0,
		Panic:  true,
		Dasum:  0,
		Dnrm2:  0,
		Idamax: -1,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{},
			},
			{
				Alpha: 0,
				Ans:   []float64{},
			},
		},
	},
	{
		Name:   "EmptyReverse",
		X:      []float64{},
		Incx:   -1,
		N:      0,
		Panic:  false,
		Dasum:  0,
		Dnrm2:  0,
		Idamax: -1,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{},
			},
			{
				Alpha: 0,
				Ans:   []float64{},
			},
		},
	},
	{
		Name:   "MultiInf",
		X:      []float64{5, math.Inf(1), math.Inf(-1), 8, 9},
		Incx:   1,
		N:      5,
		Panic:  false,
		Dasum:  math.Inf(1),
		Dnrm2:  math.Inf(1),
		Idamax: 1,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{-10, math.Inf(-1), math.Inf(1), -16, -18},
			},
			{
				Alpha: 0,
				Ans:   []float64{0, 0, 0, 0, 0},
			},
		},
	},
	{
		Name:   "NaNInf",
		X:      []float64{5, math.NaN(), math.Inf(-1), 8, 9},
		Incx:   1,
		N:      5,
		Panic:  false,
		Dasum:  math.NaN(),
		Dnrm2:  math.NaN(),
		Idamax: 2,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{-10, math.NaN(), math.Inf(1), -16, -18},
			},
			{
				Alpha: 0,
				Ans:   []float64{0, 0, 0, 0, 0},
			},
		},
	},
	{
		Name:   "InfNaN",
		X:      []float64{5, math.Inf(1), math.NaN(), 8, 9},
		Incx:   1,
		N:      5,
		Panic:  false,
		Dasum:  math.NaN(),
		Dnrm2:  math.NaN(),
		Idamax: 1,
		DscalCases: []DScalCase{
			{
				Alpha: -2,
				Ans:   []float64{-10, math.Inf(-1), math.NaN(), -16, -18},
			},
			{
				Alpha: 0,
				Ans:   []float64{0, 0, 0, 0, 0},
			},
		},
	},
}

type DoubleTwoVectorCase struct {
	Name  string
	X     []float64
	Y     []float64
	XTmp  []float64
	YTmp  []float64
	Incx  int
	Incy  int
	N     int
	Panic bool
	// For Daxpy
	DaxpyCases []DaxpyCase
	DdotAns    float64
	DswapAns   DTwoVecAnswer
	DcopyAns   DTwoVecAnswer
	DrotCases  []DrotCase
	DrotmCases []DrotmCase
}

type DaxpyCase struct {
	Alpha float64
	Ans   []float64
}

type DrotCase struct {
	C    float64
	S    float64
	XAns []float64
	YAns []float64
}

type DrotmCase struct {
	P    blas.DrotmParams
	XAns []float64
	YAns []float64
	Name string
}

type DTwoVecAnswer struct {
	X []float64
	Y []float64
}

var DoubleTwoVectorCases = []DoubleTwoVectorCase{
	{
		Name:  "UnitaryInc",
		X:     []float64{10, 15, -6, 3, 14, 7},
		Y:     []float64{8, -2, 4, 7, 6, -3},
		XTmp:  []float64{0, 0, 0, 0, 0, 0},
		YTmp:  []float64{0, 0, 0, 0, 0, 0},
		Incx:  1,
		Incy:  1,
		N:     6,
		Panic: false,
		DaxpyCases: []DaxpyCase{
			{
				Alpha: 1,
				Ans:   []float64{18, 13, -2, 10, 20, 4},
			},
			{
				Alpha: 2,
				Ans:   []float64{28, 28, -8, 13, 34, 11},
			},
			{
				Alpha: -3,
				Ans:   []float64{-22, -47, 22, -2, -36, -24},
			},
			{
				Alpha: 0,
				Ans:   []float64{8, -2, 4, 7, 6, -3},
			},
		},
		DdotAns: 110,
		DswapAns: DTwoVecAnswer{
			X: []float64{8, -2, 4, 7, 6, -3},
			Y: []float64{10, 15, -6, 3, 14, 7},
		},
		DcopyAns: DTwoVecAnswer{
			X: []float64{10, 15, -6, 3, 14, 7},
			Y: []float64{10, 15, -6, 3, 14, 7},
		},
		DrotCases: []DrotCase{
			{
				C:    math.Cos(0),
				S:    math.Sin(0),
				XAns: []float64{10, 15, -6, 3, 14, 7},
				YAns: []float64{8, -2, 4, 7, 6, -3},
			},
			{
				C:    math.Cos(25 * math.Pi / 180),
				S:    math.Sin(25 * math.Pi / 180),
				XAns: []float64{12.444023964292095, 12.749380282068351, -3.7473736752571014, 5.677251193294846, 15.224018588957296, 5.076299724034451},
				YAns: []float64{3.024279678886205, -8.151889500183792, 6.160940718590796, 5.076299724034451, -0.4788089421498931, -5.677251193294846},
			},
			{
				C:    math.Cos(0.5 * math.Pi),
				S:    math.Sin(0.5 * math.Pi),
				XAns: []float64{8, -2, 4, 7, 6, -3},
				YAns: []float64{-10, -15, 6, -3, -14, -7},
			},
			{
				C:    math.Cos(math.Pi),
				S:    math.Sin(math.Pi),
				XAns: []float64{-10, -15, 6, -3, -14, -7},
				YAns: []float64{-8, 2, -4, -7, -6, 3},
			},
		},
		DrotmCases: []DrotmCase{
			{
				P: blas.DrotmParams{
					Flag: blas.Identity,
					H:    [4]float64{0.9, 0.1, -0.1, 0.5},
				},
				XAns: []float64{10, 15, -6, 3, 14, 7},
				YAns: []float64{8, -2, 4, 7, 6, -3},
				Name: "Neg2Flag",
			},
			{
				P: blas.DrotmParams{
					Flag: blas.Rescaling,
					H:    [4]float64{0.9, 0.1, -0.1, 0.5},
				},
				XAns: []float64{8.2, 13.7, -5.8, 2, 12, 6.6},
				YAns: []float64{5, 0.5, 1.4, 3.8, 4.4, -0.8},
				Name: "Neg1Flag",
			},
			{
				P: blas.DrotmParams{
					Flag: blas.OffDiagonal,
					H:    [4]float64{1, 0.1, -0.1, 1},
				},
				XAns: []float64{9.2, 15.2, -6.4, 2.3, 13.4, 7.3},
				YAns: []float64{9, -0.5, 3.4, 7.3, 7.4, -2.3},
				Name: "ZeroFlag",
			},
			{
				P: blas.DrotmParams{
					Flag: blas.Diagonal,
					H:    [4]float64{0.5, -1, 1, 0.7},
				},
				XAns: []float64{13, 5.5, 1, 8.5, 13, 0.5},
				YAns: []float64{-4.4, -16.4, 8.8, 1.9, -9.8, -9.1},
				Name: "OneFlag",
			},
		},
	},
	{
		Name:  "UnitaryIncLong",
		X:     []float64{10, 15, -6, 3, 14, 7, 8, -9, 10},
		Y:     []float64{8, -2, 4, 7, 6, -3, 7, -6},
		XTmp:  []float64{0, 0, 0, 0, 0, 0, 0, 0, 0},
		YTmp:  []float64{0, 0, 0, 0, 0, 0, 0, 0},
		Incx:  1,
		Incy:  1,
		N:     6,
		Panic: false,
		DaxpyCases: []DaxpyCase{
			{
				Alpha: 1,
				Ans:   []float64{18, 13, -2, 10, 20, 4, 7, -6},
			},
			{
				Alpha: 2,
				Ans:   []float64{28, 28, -8, 13, 34, 11, 7, -6},
			},
			{
				Alpha: -3,
				Ans:   []float64{-22, -47, 22, -2, -36, -24, 7, -6},
			},
			{
				Alpha: 0,
				Ans:   []float64{8, -2, 4, 7, 6, -3, 7, -6},
			},
		},
		DdotAns: 110,
		DswapAns: DTwoVecAnswer{
			X: []float64{8, -2, 4, 7, 6, -3, 8, -9, 10},
			Y: []float64{10, 15, -6, 3, 14, 7, 7, -6},
		},
		DcopyAns: DTwoVecAnswer{
			X: []float64{10, 15, -6, 3, 14, 7, 8, -9, 10},
			Y: []float64{10, 15, -6, 3, 14, 7, 7, -6},
		},
		DrotCases: []DrotCase{
			{
				C:    math.Cos(0),
				S:    math.Sin(0),
				XAns: []float64{10, 15, -6, 3, 14, 7, 8, -9, 10},
				YAns: []float64{8, -2, 4, 7, 6, -3, 7, -6},
			},
			{
				C:    math.Cos(25 * math.Pi / 180),
				S:    math.Sin(25 * math.Pi / 180),
				XAns: []float64{12.444023964292095, 12.749380282068351, -3.7473736752571014, 5.677251193294846, 15.224018588957296, 5.076299724034451, 8, -9, 10},
				YAns: []float64{3.024279678886205, -8.151889500183792, 6.160940718590796, 5.076299724034451, -0.4788089421498931, -5.677251193294846, 7, -6},
			},
			{
				C:    math.Cos(0.5 * math.Pi),
				S:    math.Sin(0.5 * math.Pi),
				XAns: []float64{8, -2, 4, 7, 6, -3, 8, -9, 10},
				YAns: []float64{-10, -15, 6, -3, -14, -7, 7, -6},
			},
			{
				C:    math.Cos(math.Pi),
				S:    math.Sin(math.Pi),
				XAns: []float64{-10, -15, 6, -3, -14, -7, 8, -9, 10},
				YAns: []float64{-8, 2, -4, -7, -6, 3, 7, -6},
			},
		},
		DrotmCases: []DrotmCase{
			{
				P: blas.DrotmParams{
					Flag: blas.Identity,
					H:    [4]float64{0.9, 0.1, -0.1, 0.5},
				},
				XAns: []float64{10, 15, -6, 3, 14, 7, 8, -9, 10},
				YAns: []float64{8, -2, 4, 7, 6, -3, 7, -6},
				Name: "Neg2Flag",
			},
			{
				P: blas.DrotmParams{
					Flag: blas.Rescaling,
					H:    [4]float64{0.9, 0.1, -0.1, 0.5},
				},
				XAns: []float64{8.2, 13.7, -5.8, 2, 12, 6.6, 8, -9, 10},
				YAns: []float64{5, 0.5, 1.4, 3.8, 4.4, -0.8, 7, -6},
				Name: "Neg1Flag",
			},
			{
				P: blas.DrotmParams{
					Flag: blas.OffDiagonal,
					H:    [4]float64{1, 0.1, -0.1, 1},
				},
				XAns: []float64{9.2, 15.2, -6.4, 2.3, 13.4, 7.3, 8, -9, 10},
				YAns: []float64{9, -0.5, 3.4, 7.3, 7.4, -2.3, 7, -6},
				Name: "ZeroFlag",
			},
			{
				P: blas.DrotmParams{
					Flag: blas.Diagonal,
					H:    [4]float64{0.5, -1, 1, 0.7},
				},
				XAns: []float64{13, 5.5, 1, 8.5, 13, 0.5, 8, -9, 10},
				YAns: []float64{-4.4, -16.4, 8.8, 1.9, -9.8, -9.1, 7, -6},
				Name: "OneFlag",
			},
		},
	},
	{
		Name:  "PositiveInc",
		X:     []float64{10, 15, -6, 3, 14, 7},
		Y:     []float64{8, -2, 4, 7, 6, -3, -4, 10},
		XTmp:  []float64{0, 0, 0, 0, 0, 0},
		YTmp:  []float64{0, 0, 0, 0, 0, 0, 0, 0},
		Incx:  2,
		Incy:  3,
		N:     3,
		Panic: false,
		DaxpyCases: []DaxpyCase{
			{
				Alpha: 2,
				Ans:   []float64{28, -2, 4, -5, 6, -3, 24, 10},
			},
		},
		DdotAns: -18,
		DswapAns: DTwoVecAnswer{
			X: []float64{8, 15, 7, 3, -4, 7},
			Y: []float64{10, -2, 4, -6, 6, -3, 14, 10},
		},
		DcopyAns: DTwoVecAnswer{
			X: []float64{10, 15, -6, 3, 14, 7},
			Y: []float64{10, -2, 4, -6, 6, -3, 14, 10},
		},
		DrotCases: []DrotCase{
			{
				C:    math.Cos(25 * math.Pi / 180),
				S:    math.Sin(25 * math.Pi / 180),
				XAns: []float64{12.444023964292095, 15, -2.479518890035003, 3, 10.997835971550302, 7},
				YAns: []float64{3.024279678886205, -2, 4, 8.879864079700745, 6, -3, -9.541886812516392, 10},
			},
		},
		DrotmCases: []DrotmCase{
			{
				P: blas.DrotmParams{
					Flag: blas.Rescaling,
					H:    [4]float64{0.9, 0.1, -0.1, 0.5},
				},
				XAns: []float64{8.2, 15, -6.1, 3, 13, 7},
				YAns: []float64{5, -2, 4, 2.9, 6, -3, -0.6, 10},
			},
			{
				P: blas.DrotmParams{
					Flag: blas.OffDiagonal,
					H:    [4]float64{1, 0.1, -0.1, 1},
				},
				XAns: []float64{9.2, 15, -6.7, 3, 14.4, 7},
				YAns: []float64{9, -2, 4, 6.4, 6, -3, -2.6, 10},
			},
			{
				P: blas.DrotmParams{
					Flag: blas.Diagonal,
					H:    [4]float64{0.5, -1, 1, 0.7},
				},
				XAns: []float64{13, 15, 4, 3, 3, 7},
				YAns: []float64{-4.4, -2, 4, 10.9, 6, -3, -16.8, 10},
			},
		},
	},
	{
		Name:  "NegativeInc",
		X:     []float64{10, 15, -6, 3, 14, 7},
		Y:     []float64{8, -2, 4, 7, 6, -3, -4, 10},
		XTmp:  []float64{0, 0, 0, 0, 0, 0},
		YTmp:  []float64{0, 0, 0, 0, 0, 0, 0, 0},
		Incx:  -2,
		Incy:  -3,
		N:     3,
		Panic: false,
		DaxpyCases: []DaxpyCase{
			{
				Alpha: 2,
				Ans:   []float64{28, -2, 4, -5, 6, -3, 24, 10},
			},
		},
		DdotAns: -18,
		DswapAns: DTwoVecAnswer{
			X: []float64{8, 15, 7, 3, -4, 7},
			Y: []float64{10, -2, 4, -6, 6, -3, 14, 10},
		},
		DcopyAns: DTwoVecAnswer{
			X: []float64{10, 15, -6, 3, 14, 7},
			Y: []float64{10, -2, 4, -6, 6, -3, 14, 10},
		},
		DrotCases: []DrotCase{
			{
				C:    math.Cos(25 * math.Pi / 180),
				S:    math.Sin(25 * math.Pi / 180),
				XAns: []float64{12.444023964292095, 15, -2.479518890035003, 3, 10.997835971550302, 7},
				YAns: []float64{3.024279678886205, -2, 4, 8.879864079700745, 6, -3, -9.541886812516392, 10},
			},
		},
		DrotmCases: []DrotmCase{
			{
				P: blas.DrotmParams{
					Flag: blas.Rescaling,
					H:    [4]float64{0.9, 0.1, -0.1, 0.5},
				},
				XAns: []float64{8.2, 15, -6.1, 3, 13, 7},
				YAns: []float64{5, -2, 4, 2.9, 6, -3, -0.6, 10},
			},
			{
				P: blas.DrotmParams{
					Flag: blas.OffDiagonal,
					H:    [4]float64{1, 0.1, -0.1, 1},
				},
				XAns: []float64{9.2, 15, -6.7, 3, 14.4, 7},
				YAns: []float64{9, -2, 4, 6.4, 6, -3, -2.6, 10},
			},
			{
				P: blas.DrotmParams{
					Flag: blas.Diagonal,
					H:    [4]float64{0.5, -1, 1, 0.7},
				},
				XAns: []float64{13, 15, 4, 3, 3, 7},
				YAns: []float64{-4.4, -2, 4, 10.9, 6, -3, -16.8, 10},
			},
		},
	},
	{
		Name:  "MixedInc1",
		X:     []float64{10, 15, -6, 3, 14, 7},
		Y:     []float64{8, -2, 4, 7, 6, -3, -4, 10},
		XTmp:  []float64{0, 0, 0, 0, 0, 0},
		YTmp:  []float64{0, 0, 0, 0, 0, 0, 0, 0},
		Incx:  2,
		Incy:  -3,
		N:     3,
		Panic: false,
		DaxpyCases: []DaxpyCase{
			{
				Alpha: 2,
				Ans:   []float64{36, -2, 4, -5, 6, -3, 16, 10},
			},
		},
		DdotAns: 30,
		DswapAns: DTwoVecAnswer{
			X: []float64{-4, 15, 7, 3, 8, 7},
			Y: []float64{14, -2, 4, -6, 6, -3, 10, 10},
		},
		DcopyAns: DTwoVecAnswer{
			X: []float64{10, 15, -6, 3, 14, 7},
			Y: []float64{14, -2, 4, -6, 6, -3, 10, 10},
		},
		DrotCases: []DrotCase{
			{
				C:    math.Cos(25 * math.Pi / 180),
				S:    math.Sin(25 * math.Pi / 180),
				XAns: []float64{7.372604823403701, 15, -2.479518890035003, 3, 16.069255112438693, 7},
				YAns: []float64{1.333806631923407, -2, 4, 8.879864079700745, 6, -3, -7.851413765553595, 10},
			},
		},
		DrotmCases: []DrotmCase{
			{
				P: blas.DrotmParams{
					Flag: blas.Rescaling,
					H:    [4]float64{0.9, 0.1, -0.1, 0.5},
				},
				XAns: []float64{9.4, 15, -6.1, 3, 11.8, 7},
				YAns: []float64{5.4, -2, 4, 2.9, 6, -3, -1, 10},
			},
			{
				P: blas.DrotmParams{
					Flag: blas.OffDiagonal,
					H:    [4]float64{1, 0.1, -0.1, 1},
				},
				XAns: []float64{10.4, 15, -6.7, 3, 13.2, 7},
				YAns: []float64{9.4, -2, 4, 6.4, 6, -3, -3, 10},
			},
			{
				P: blas.DrotmParams{
					Flag: blas.Diagonal,
					H:    [4]float64{0.5, -1, 1, 0.7},
				},
				XAns: []float64{1, 15, 4, 3, 15, 7},
				YAns: []float64{-8.4, -2, 4, 10.9, 6, -3, -12.8, 10},
			},
		},
	},
	{
		Name:  "MixedInc2",
		X:     []float64{10, 15, -6, 3, 14, 7},
		Y:     []float64{8, -2, 4, 7, 6, -3, -4, 10},
		XTmp:  []float64{0, 0, 0, 0, 0, 0},
		YTmp:  []float64{0, 0, 0, 0, 0, 0, 0, 0},
		Incx:  -2,
		Incy:  3,
		N:     3,
		Panic: false,
		DaxpyCases: []DaxpyCase{
			{
				Alpha: 2,
				Ans:   []float64{36, -2, 4, -5, 6, -3, 16, 10},
			},
		},
		DdotAns: 30,
		DswapAns: DTwoVecAnswer{
			X: []float64{-4, 15, 7, 3, 8, 7},
			Y: []float64{14, -2, 4, -6, 6, -3, 10, 10},
		},
		DcopyAns: DTwoVecAnswer{
			X: []float64{10, 15, -6, 3, 14, 7},
			Y: []float64{14, -2, 4, -6, 6, -3, 10, 10},
		},
		DrotCases: []DrotCase{
			{
				C:    math.Cos(25 * math.Pi / 180),
				S:    math.Sin(25 * math.Pi / 180),
				XAns: []float64{7.372604823403701, 15, -2.479518890035003, 3, 16.069255112438693, 7},
				YAns: []float64{1.333806631923407, -2, 4, 8.879864079700745, 6, -3, -7.851413765553595, 10},
			},
		},
		DrotmCases: []DrotmCase{
			{
				P: blas.DrotmParams{
					Flag: blas.Rescaling,
					H:    [4]float64{0.9, 0.1, -0.1, 0.5},
				},
				XAns: []float64{9.4, 15, -6.1, 3, 11.8, 7},
				YAns: []float64{5.4, -2, 4, 2.9, 6, -3, -1, 10},
			},
			{
				P: blas.DrotmParams{
					Flag: blas.OffDiagonal,
					H:    [4]float64{1, 0.1, -0.1, 1},
				},
				XAns: []float64{10.4, 15, -6.7, 3, 13.2, 7},
				YAns: []float64{9.4, -2, 4, 6.4, 6, -3, -3, 10},
			},
			{
				P: blas.DrotmParams{
					Flag: blas.Diagonal,
					H:    [4]float64{0.5, -1, 1, 0.7},
				},
				XAns: []float64{1, 15, 4, 3, 15, 7},
				YAns: []float64{-8.4, -2, 4, 10.9, 6, -3, -12.8, 10},
			},
		},
	},
	{
		Name:  "ZeroN",
		X:     []float64{10, 15, -6, 3, 14, 7},
		Y:     []float64{8, -2, 4, 7, 6, -3, -4, 10},
		XTmp:  []float64{0, 0, 0, 0, 0, 0},
		YTmp:  []float64{0, 0, 0, 0, 0, 0, 0, 0},
		Incx:  -2,
		Incy:  3,
		N:     0,
		Panic: false,
		DaxpyCases: []DaxpyCase{
			{
				Alpha: 2,
				Ans:   []float64{8, -2, 4, 7, 6, -3, -4, 10},
			},
		},
		DswapAns: DTwoVecAnswer{
			X: []float64{10, 15, -6, 3, 14, 7},
			Y: []float64{8, -2, 4, 7, 6, -3, -4, 10},
		},
		DcopyAns: DTwoVecAnswer{
			X: []float64{10, 15, -6, 3, 14, 7},
			Y: []float64{8, -2, 4, 7, 6, -3, -4, 10},
		},
		DrotCases: []DrotCase{
			{
				C:    math.Cos(25 * math.Pi / 180),
				S:    math.Sin(25 * math.Pi / 180),
				XAns: []float64{10, 15, -6, 3, 14, 7},
				YAns: []float64{8, -2, 4, 7, 6, -3, -4, 10},
			},
		},
		DrotmCases: []DrotmCase{
			{
				P: blas.DrotmParams{
					Flag: blas.Rescaling,
					H:    [4]float64{0.9, 0.1, -0.1, 0.5},
				},
				XAns: []float64{10, 15, -6, 3, 14, 7},
				YAns: []float64{8, -2, 4, 7, 6, -3, -4, 10},
			},
		},
	},
	{
		Name:  "NegativeN",
		X:     []float64{10, 15, -6, 3, 14, 7},
		Y:     []float64{8, -2, 4, 7, 6, -3, -4, 10},
		XTmp:  []float64{0, 0, 0, 0, 0, 0},
		YTmp:  []float64{0, 0, 0, 0, 0, 0, 0, 0},
		Incx:  -2,
		Incy:  3,
		N:     -3,
		Panic: true,
		DaxpyCases: []DaxpyCase{
			{
				Alpha: 2,
				Ans:   []float64{36, -2, 4, -5, 6, -3, 16, 10},
			},
		},
		DrotCases: []DrotCase{
			{
				C:    math.Cos(25 * math.Pi / 180),
				S:    math.Sin(25 * math.Pi / 180),
				XAns: []float64{10, 15, -6, 3, 14, 7},
				YAns: []float64{8, -2, 4, 7, 6, -3, -4, 10},
			},
		},
		DrotmCases: []DrotmCase{
			{
				P: blas.DrotmParams{
					Flag: blas.Rescaling,
					H:    [4]float64{0.9, 0.1, -0.1, 0.5},
				},
				XAns: []float64{8.2, 13.7, -5.8, 2, 12, 6.6},
				YAns: []float64{5, 0.5, 1.4, 3.8, 4.4, -0.8},
			},
		},
	},
	{
		Name:  "ZeroIncX",
		X:     []float64{10, 15, -6, 3, 14, 7},
		Y:     []float64{8, -2, 4, 7, 6, -3, -4, 10},
		XTmp:  []float64{0, 0, 0, 0, 0, 0},
		YTmp:  []float64{0, 0, 0, 0, 0, 0, 0, 0},
		Incx:  0,
		Incy:  3,
		N:     2,
		Panic: true,
		DaxpyCases: []DaxpyCase{
			{
				Alpha: 2,
				Ans:   []float64{36, -2, 4, -5, 6, -3, 16, 10},
			},
		},
		DrotCases: []DrotCase{
			{
				C:    math.Cos(25 * math.Pi / 180),
				S:    math.Sin(25 * math.Pi / 180),
				XAns: []float64{10, 15, -6, 3, 14, 7},
				YAns: []float64{8, -2, 4, 7, 6, -3, -4, 10},
			},
		},
		DrotmCases: []DrotmCase{
			{
				P: blas.DrotmParams{
					Flag: blas.Rescaling,
					H:    [4]float64{0.9, 0.1, -0.1, 0.5},
				},
				XAns: []float64{8.2, 13.7, -5.8, 2, 12, 6.6},
				YAns: []float64{5, 0.5, 1.4, 3.8, 4.4, -0.8},
			},
		},
	},
	{
		Name:  "ZeroIncY",
		X:     []float64{10, 15, -6, 3, 14, 7},
		Y:     []float64{8, -2, 4, 7, 6, -3, -4, 10},
		XTmp:  []float64{0, 0, 0, 0, 0, 0},
		YTmp:  []float64{0, 0, 0, 0, 0, 0, 0, 0},
		Incx:  1,
		Incy:  0,
		N:     2,
		Panic: true,
		DaxpyCases: []DaxpyCase{
			{
				Alpha: 2,
				Ans:   []float64{36, -2, 4, -5, 6, -3, 16, 10},
			},
		},
		DrotCases: []DrotCase{
			{
				C:    math.Cos(25 * math.Pi / 180),
				S:    math.Sin(25 * math.Pi / 180),
				XAns: []float64{10, 15, -6, 3, 14, 7},
				YAns: []float64{8, -2, 4, 7, 6, -3, -4, 10},
			},
		},
		DrotmCases: []DrotmCase{
			{
				P: blas.DrotmParams{
					Flag: blas.Rescaling,
					H:    [4]float64{0.9, 0.1, -0.1, 0.5},
				},
				XAns: []float64{8.2, 13.7, -5.8, 2, 12, 6.6},
				YAns: []float64{5, 0.5, 1.4, 3.8, 4.4, -0.8},
			},
		},
	},
	{
		Name:  "OutOfBoundsX",
		X:     []float64{10, 15, -6, 3, 14, 7},
		Y:     []float64{8, -2, 4, 7, 6, -3, -4, 10},
		XTmp:  []float64{0, 0, 0, 0, 0, 0},
		YTmp:  []float64{0, 0, 0, 0, 0, 0, 0, 0},
		Incx:  8,
		Incy:  2,
		N:     2,
		Panic: true,
		DaxpyCases: []DaxpyCase{
			{
				Alpha: 2,
				Ans:   []float64{36, -2, 4, -5, 6, -3, 16, 10},
			},
		},
		DrotCases: []DrotCase{
			{
				C:    math.Cos(25 * math.Pi / 180),
				S:    math.Sin(25 * math.Pi / 180),
				XAns: []float64{10, 15, -6, 3, 14, 7},
				YAns: []float64{8, -2, 4, 7, 6, -3, -4, 10},
			},
		},
		DrotmCases: []DrotmCase{
			{
				P: blas.DrotmParams{
					Flag: blas.Rescaling,
					H:    [4]float64{0.9, 0.1, -0.1, 0.5},
				},
				XAns: []float64{8.2, 13.7, -5.8, 2, 12, 6.6},
				YAns: []float64{5, 0.5, 1.4, 3.8, 4.4, -0.8},
			},
		},
	},
	{
		Name:  "OutOfBoundsY",
		X:     []float64{10, 15, -6, 3, 14, 7},
		Y:     []float64{8, -2, 4, 7, 6, -3, -4, 10},
		XTmp:  []float64{0, 0, 0, 0, 0, 0},
		YTmp:  []float64{0, 0, 0, 0, 0, 0, 0, 0},
		Incx:  2,
		Incy:  8,
		N:     2,
		Panic: true,
		DaxpyCases: []DaxpyCase{
			{
				Alpha: 2,
				Ans:   []float64{36, -2, 4, -5, 6, -3, 16, 10},
			},
		},
		DrotCases: []DrotCase{
			{
				C:    math.Cos(25 * math.Pi / 180),
				S:    math.Sin(25 * math.Pi / 180),
				XAns: []float64{10, 15, -6, 3, 14, 7},
				YAns: []float64{8, -2, 4, 7, 6, -3, -4, 10},
			},
		},
		DrotmCases: []DrotmCase{
			{
				P: blas.DrotmParams{
					Flag: blas.Rescaling,
					H:    [4]float64{0.9, 0.1, -0.1, 0.5},
				},
				XAns: []float64{10, 15, -6, 3, 14, 7},
				YAns: []float64{8, -2, 4, 7, 6, -3, -4, 10},
			},
		},
	},
	{
		Name:  "Empty",
		X:     []float64{},
		Y:     []float64{},
		Incx:  1,
		Incy:  1,
		N:     0,
		Panic: false,
		DaxpyCases: []DaxpyCase{
			{
				Alpha: 2,
				Ans:   []float64{},
			},
		},
		DrotCases: []DrotCase{
			{
				C:    math.Cos(25 * math.Pi / 180),
				S:    math.Sin(25 * math.Pi / 180),
				XAns: []float64{},
				YAns: []float64{},
			},
		},
		DrotmCases: []DrotmCase{
			{
				P: blas.DrotmParams{
					Flag: blas.Rescaling,
					H:    [4]float64{0.9, 0.1, -0.1, 0.5},
				},
				XAns: []float64{},
				YAns: []float64{},
			},
		},
	},
	{
		Name:  "EmptyZeroIncX",
		X:     []float64{},
		Y:     []float64{},
		Incx:  0,
		Incy:  1,
		N:     0,
		Panic: true,
		DaxpyCases: []DaxpyCase{
			{
				Alpha: 2,
				Ans:   []float64{},
			},
		},
		DrotCases: []DrotCase{
			{
				C:    math.Cos(25 * math.Pi / 180),
				S:    math.Sin(25 * math.Pi / 180),
				XAns: []float64{},
				YAns: []float64{},
			},
		},
		DrotmCases: []DrotmCase{
			{
				P: blas.DrotmParams{
					Flag: blas.Rescaling,
					H:    [4]float64{0.9, 0.1, -0.1, 0.5},
				},
				XAns: []float64{},
				YAns: []float64{},
			},
		},
	},
	{
		Name:  "EmptyZeroIncY",
		X:     []float64{},
		Y:     []float64{},
		Incx:  1,
		Incy:  0,
		N:     0,
		Panic: true,
		DaxpyCases: []DaxpyCase{
			{
				Alpha: 2,
				Ans:   []float64{},
			},
		},
		DrotCases: []DrotCase{
			{
				C:    math.Cos(25 * math.Pi / 180),
				S:    math.Sin(25 * math.Pi / 180),
				XAns: []float64{},
				YAns: []float64{},
			},
		},
		DrotmCases: []DrotmCase{
			{
				P: blas.DrotmParams{
					Flag: blas.Rescaling,
					H:    [4]float64{0.9, 0.1, -0.1, 0.5},
				},
				XAns: []float64{},
				YAns: []float64{},
			},
		},
	},
	{
		Name:  "EmptyReverse",
		X:     []float64{},
		Y:     []float64{},
		Incx:  -1,
		Incy:  -1,
		N:     0,
		Panic: false,
		DaxpyCases: []DaxpyCase{
			{
				Alpha: 2,
				Ans:   []float64{},
			},
		},
		DrotCases: []DrotCase{
			{
				C:    math.Cos(25 * math.Pi / 180),
				S:    math.Sin(25 * math.Pi / 180),
				XAns: []float64{},
				YAns: []float64{},
			},
		},
		DrotmCases: []DrotmCase{
			{
				P: blas.DrotmParams{
					Flag: blas.Rescaling,
					H:    [4]float64{0.9, 0.1, -0.1, 0.5},
				},
				XAns: []float64{},
				YAns: []float64{},
			},
		},
	},
}

type Ddotter interface {
	Ddot(n int, x []float64, incX int, y []float64, incY int) float64
}

func DdotTest(t *testing.T, d Ddotter) {
	ddot := d.Ddot
	for _, c := range DoubleTwoVectorCases {
		dCopyTwoTmp(c.X, c.XTmp, c.Y, c.YTmp)
		if c.Panic {
			f := func() { ddot(c.N, c.XTmp, c.Incx, c.YTmp, c.Incy) }
			testpanics(f, c.Name, t)
			continue
		}
		dot := ddot(c.N, c.XTmp, c.Incx, c.YTmp, c.Incy)
		if !dTolEqual(dot, c.DdotAns) {
			t.Errorf("ddot: mismatch %v: expected %v, found %v", c.Name, c.DdotAns, dot)
		}
	}

	// check it works for 16-byte unaligned slices
	x := []float64{1, 1, 1, 1, 1}
	if n := ddot(4, x[:4], 1, x[1:], 1); n != 4 {
		t.Errorf("ddot: mismatch Unaligned: expected %v, found %v", 4, n)
	}
	if n := ddot(2, x[:4], 2, x[1:], 2); n != 2 {
		t.Errorf("ddot: mismatch Unaligned: expected %v, found %v", 2, n)
	}
	if n := ddot(2, x[:4], 3, x[1:], 3); n != 2 {
		t.Errorf("ddot: mismatch Unaligned: expected %v, found %v", 2, n)
	}
}

type Dnrm2er interface {
	Dnrm2(n int, x []float64, incX int) float64
}

func Dnrm2Test(t *testing.T, blasser Dnrm2er) {
	dnrm2 := blasser.Dnrm2
	for _, c := range DoubleOneVectorCases {
		if c.Panic {
			f := func() { dnrm2(c.N, c.X, c.Incx) }
			testpanics(f, c.Name, t)
			continue
		}
		v := dnrm2(c.N, c.X, c.Incx)
		if !dTolEqual(v, c.Dnrm2) {
			t.Errorf("dnrm2: mismatch %v: expected %v, found %v", c.Name, c.Dnrm2, v)
		}
	}
}

type Dasumer interface {
	Dasum(n int, x []float64, incX int) float64
}

func DasumTest(t *testing.T, blasser Dasumer) {
	dasum := blasser.Dasum
	for _, c := range DoubleOneVectorCases {
		if c.Panic {
			f := func() { dasum(c.N, c.X, c.Incx) }
			testpanics(f, c.Name, t)
			continue
		}
		v := dasum(c.N, c.X, c.Incx)
		if !dTolEqual(v, c.Dasum) {
			t.Errorf("dasum: mismatch %v: expected %v, found %v", c.Name, c.Dasum, v)
		}
	}
}

type Idamaxer interface {
	Idamax(n int, x []float64, incX int) int
}

func IdamaxTest(t *testing.T, blasser Idamaxer) {
	idamax := blasser.Idamax
	for _, c := range DoubleOneVectorCases {
		if c.Panic {
			f := func() { idamax(c.N, c.X, c.Incx) }
			testpanics(f, c.Name, t)
			continue
		}
		v := idamax(c.N, c.X, c.Incx)
		if v != c.Idamax {
			s := fmt.Sprintf("idamax: mismatch %v: expected %v, found %v", c.Name, c.Idamax, v)
			if floats.HasNaN(c.X) {
				t.Log(s)
			} else {
				t.Error(s)
			}
		}
	}
}

type Dswapper interface {
	Dswap(n int, x []float64, incX int, y []float64, incY int)
}

func DswapTest(t *testing.T, d Dswapper) {
	dswap := d.Dswap
	for _, c := range DoubleTwoVectorCases {
		dCopyTwoTmp(c.X, c.XTmp, c.Y, c.YTmp)
		if c.Panic {
			f := func() { dswap(c.N, c.XTmp, c.Incx, c.YTmp, c.Incy) }
			testpanics(f, c.Name, t)
			continue
		}
		dswap(c.N, c.XTmp, c.Incx, c.YTmp, c.Incy)
		if !dSliceTolEqual(c.XTmp, c.DswapAns.X) {
			t.Errorf("dswap: x mismatch %v: expected %v, found %v", c.Name, c.DswapAns.X, c.XTmp)
		}
		if !dSliceTolEqual(c.YTmp, c.DswapAns.Y) {
			t.Errorf("dswap: y mismatch %v: expected %v, found %v", c.Name, c.DswapAns.Y, c.YTmp)
		}
	}
}

type Dcopier interface {
	Dcopy(n int, x []float64, incX int, y []float64, incY int)
}

func DcopyTest(t *testing.T, d Dcopier) {
	dcopy := d.Dcopy
	for _, c := range DoubleTwoVectorCases {
		dCopyTwoTmp(c.X, c.XTmp, c.Y, c.YTmp)
		if c.Panic {
			f := func() { dcopy(c.N, c.XTmp, c.Incx, c.YTmp, c.Incy) }
			testpanics(f, c.Name, t)
			continue
		}
		dcopy(c.N, c.XTmp, c.Incx, c.YTmp, c.Incy)
		if !dSliceTolEqual(c.XTmp, c.DcopyAns.X) {
			t.Errorf("dswap: x mismatch %v: expected %v, found %v", c.Name, c.DcopyAns.X, c.XTmp)
		}
		if !dSliceTolEqual(c.YTmp, c.DcopyAns.Y) {
			t.Errorf("dswap: y mismatch %v: expected %v, found %v", c.Name, c.DcopyAns.Y, c.YTmp)
		}
	}
}

type Daxpyer interface {
	Daxpy(n int, alpha float64, x []float64, incX int, y []float64, incY int)
}

func DaxpyTest(t *testing.T, d Daxpyer) {
	daxpy := d.Daxpy
	for _, c := range DoubleTwoVectorCases {
		for _, kind := range c.DaxpyCases {
			dCopyTwoTmp(c.X, c.XTmp, c.Y, c.YTmp)
			if c.Panic {
				f := func() { daxpy(c.N, kind.Alpha, c.XTmp, c.Incx, c.YTmp, c.Incy) }
				testpanics(f, c.Name, t)
				continue
			}
			daxpy(c.N, kind.Alpha, c.XTmp, c.Incx, c.YTmp, c.Incy)
			if !dSliceTolEqual(c.YTmp, kind.Ans) {
				t.Errorf("daxpy: mismatch %v: expected %v, found %v", c.Name, kind.Ans, c.YTmp)
			}
		}
	}
}

type DrotgTestStruct struct {
	Name       string
	A, B       float64
	C, S, R, Z float64
}

var DrotgTests = []DrotgTestStruct{
	{
		Name: "ZeroAB",
		C:    1,
	},
	{
		Name: "PosA_ZeroB",
		A:    0.5,
		C:    1,
		R:    0.5,
	},
	{
		Name: "NegA_ZeroB",
		A:    -4.6,
		C:    1,
		R:    -4.6,
	},
	{
		Name: "ZeroA_PosB",
		B:    3,
		S:    1,
		R:    3,
		Z:    1,
	},
	{
		Name: "ZeroA_NegB",
		B:    -0.3,
		S:    1,
		R:    -0.3,
		Z:    1,
	},
	{
		Name: "PosA_PosB_AGTB",
		A:    5,
		B:    0.3,
		C:    0.99820484546577868593549038000,
		S:    0.05989229072794672115612942280,
		R:    5.00899191454727744602429072688,
		Z:    0.05989229072794672115612942280,
	},
	{
		Name: "PosA_PosB_ALTB",
		A:    3,
		B:    4,
		C:    3.0 / 5,
		S:    4.0 / 5,
		R:    5,
		Z:    5.0 / 3.0,
	},

	{
		Name: "PosA_NegB_AGTB",
		A:    2.6,
		B:    -0.9,
		C:    0.94498607344025815971847507095,
		S:    -0.32711056388316628605639521686,
		R:    2.751363298439520872718790879655,
		Z:    -0.3271105638831662860563952168,
	},
	{
		Name: "PosA_NegB_ALTB",
		A:    2.6,
		B:    -2.9,
		C:    -0.6675450157520258540548049558,
		S:    0.7445694406464903756765132200,
		R:    -3.8948684188300893100043812234,
		Z:    1 / -0.6675450157520258540548049558,
	},
	{
		Name: "NegA_PosB_AGTB",
		A:    -11.4,
		B:    10.3,
		C:    0.7419981952497362418487847947,
		S:    -0.6704018781642353764072353847,
		R:    -15.363918770938617534070671122,
		Z:    -0.6704018781642353764072353847,
	},
	{
		Name: "NegA_PosB_ALTB",
		A:    -1.4,
		B:    10.3,
		C:    -0.1346838895922121112404717523,
		S:    0.9908886162855605326977564640,
		R:    10.394710193170370442523552032,
		Z:    1 / -0.1346838895922121112404717523,
	},
	{
		Name: "NegA_NegB_AGTB",
		A:    -11.4,
		B:    10.3,
		C:    0.7419981952497362418487847947,
		S:    -0.6704018781642353764072353847,
		R:    -15.363918770938617534070671122,
		Z:    -0.6704018781642353764072353847,
	},
	{
		Name: "NegA_NegB_ALTB",
		A:    -1.4,
		B:    -10.3,
		C:    0.1346838895922121112404717523,
		S:    0.9908886162855605326977564640,
		R:    -10.394710193170370442523552032,
		Z:    1 / 0.1346838895922121112404717523,
	},
}

type Drotger interface {
	Drotg(a, b float64) (c, s, r, z float64)
}

func DrotgTest(t *testing.T, d Drotger, skipExtreme bool) {
	drotg := d.Drotg
	for _, test := range DrotgTests {
		c, s, r, z := drotg(test.A, test.B)
		if !dTolEqual(c, test.C) {
			t.Errorf("drotg: c mismatch %v: expected %v, found %v", test.Name, test.C, c)
		}
		if !dTolEqual(s, test.S) {
			t.Errorf("drotg: s mismatch %v: expected %v, found %v", test.Name, test.S, s)
		}
		if !dTolEqual(r, test.R) {
			t.Errorf("drotg: r mismatch %v: expected %v, found %v", test.Name, test.R, r)
		}
		if !dTolEqual(z, test.Z) {
			t.Errorf("drotg: z mismatch %v: expected %v, found %v", test.Name, test.Z, z)
		}
	}

	const (
		ulp    = 0x1p-52
		safmin = 0x1p-1022
		safmax = 1 / safmin
		tol    = 20 * ulp
	)
	values := []float64{
		-1 / ulp,
		-1,
		-1.0 / 3,
		-ulp,
		0,
		ulp,
		1.0 / 3,
		1,
		1 / ulp,
		math.Inf(-1),
		math.Inf(1),
		math.NaN(),
	}
	if !skipExtreme {
		values = append(values, -safmax, -safmin, safmin, safmax)
	}
	for _, f := range values {
		for _, g := range values {
			name := fmt.Sprintf("Case f=%v,g=%v", f, g)

			// Generate a plane rotation so that
			//  [ cs sn] * [f] = [r]
			//  [-sn cs]   [g] = [0]
			// where cs*cs + sn*sn = 1.
			cs, sn, r, _ := drotg(f, g)

			switch {
			case math.IsNaN(f) || math.IsNaN(g):
				if !math.IsNaN(r) {
					t.Errorf("%v: unexpected r=%v; want NaN", name, r)
				}
			case math.IsInf(f, 0) || math.IsInf(g, 0):
				if !math.IsNaN(r) && !math.IsInf(r, 0) {
					t.Errorf("%v: unexpected r=%v; want NaN or Inf", name, r)
				}
			default:
				d := math.Max(math.Abs(f), math.Abs(g))
				d = math.Min(math.Max(safmin, d), safmax)
				fs := f / d
				gs := g / d
				rs := r / d

				// Check that cs*f + sn*g = r.
				rnorm := math.Abs(rs)
				if rnorm == 0 {
					rnorm = math.Max(math.Abs(fs), math.Abs(gs))
					if rnorm == 0 {
						rnorm = 1
					}
				}
				resid := math.Abs(rs-(cs*fs+sn*gs)) / rnorm
				if resid > tol {
					t.Errorf("%v: cs*f + sn*g != r; resid=%v", name, resid)
				}

				// Check that -sn*f + cs*g = 0.
				resid = math.Abs(-sn*fs + cs*gs)
				if resid > tol {
					t.Errorf("%v: -sn*f + cs*g != 0; resid=%v", name, resid)
				}

				// Check that cs*cs + sn*sn = 1.
				resid = math.Abs(1 - (cs*cs + sn*sn))
				if resid > tol {
					t.Errorf("%v: cs*cs + sn*sn != 1; resid=%v", name, resid)
				}

				// Check that cs is non-negative.
				if math.Abs(f) > math.Abs(g) {
					if cs < 0 {
						t.Errorf("%v: cs is negative; cs=%v", name, cs)
					}
				} else {
					if cs*math.Copysign(1, f)*math.Copysign(1, g) < 0 {
						t.Errorf("%v: sign of cs doesn't match sign of f and g; cs=%v, sign(f)=%v, sign(g)=%v",
							name, cs, math.Copysign(1, f), math.Copysign(1, g))
					}
				}
			}
		}
	}
}

type Drotmger interface {
	Drotmg(d1, d2, x1, y1 float64) (p blas.DrotmParams, rd1, rd2, rx1 float64)
	Drotmer
}

func DrotmgTest(t *testing.T, d Drotmger) {
	const tol = 1e-14

	var tests = []struct {
		name                   string
		d1, d2, x1, y1         float64
		d1Want, d2Want, x1Want float64
		pWant                  blas.DrotmParams
	}{
		{
			name:   "NegD1",
			d1:     -4,
			d2:     6,
			x1:     8,
			y1:     -4,
			d1Want: 0,
			d2Want: 0,
			x1Want: 0,
			pWant: blas.DrotmParams{
				Flag: blas.Rescaling,
				H:    [4]float64{0, 0, 0, 0},
			},
		},
		{
			name:   "ZeroD2",
			d1:     4,
			d2:     0,
			x1:     8,
			y1:     -5,
			d1Want: 4,
			d2Want: 0,
			x1Want: 8,
			pWant: blas.DrotmParams{
				Flag: blas.Identity,
				H:    [4]float64{0, 0, 0, 0},
			},
		},
		{
			name:   "ZeroY1",
			d1:     4,
			d2:     -6,
			x1:     8,
			y1:     0,
			d1Want: 4,
			d2Want: -6,
			x1Want: 8,
			pWant: blas.DrotmParams{
				Flag: blas.Identity,
				H:    [4]float64{0, 0, 0, 0},
			},
		},
		{
			name:   "NegQ2_and_AQ1_LT_AQ2",
			d1:     8,
			d2:     -6,
			x1:     4,
			y1:     8,
			d1Want: 0,
			d2Want: 0,
			x1Want: 0,
			pWant: blas.DrotmParams{
				Flag: blas.Rescaling,
				H:    [4]float64{0, 0, 0, 0},
			},
		},
		{
			name:   "ZeroD1",
			d1:     0,
			d2:     2,
			x1:     8,
			y1:     4,
			d1Want: 2,
			d2Want: 0,
			x1Want: 4,
			pWant: blas.DrotmParams{
				Flag: blas.Diagonal,
				H:    [4]float64{0, 0, 0, 0},
			},
		},
		{
			name:   "AbsQ1_GT_AbsQU__D2_Pos",
			d1:     2,
			d2:     3,
			x1:     8,
			y1:     5,
			d1Want: 1.2610837438423645,
			d2Want: 1.8916256157635467,
			x1Want: 12.6875,
			pWant: blas.DrotmParams{
				Flag: blas.OffDiagonal,
				H:    [4]float64{0, -0.625, 0.9375, 0},
			},
		},
		{
			name:   "AbsQ1_GT_AbsQU__D2_Neg",
			d1:     2,
			d2:     -3,
			x1:     8,
			y1:     5,
			d1Want: 4.830188679245283,
			d2Want: -7.245283018867925,
			x1Want: 3.3125,
			pWant: blas.DrotmParams{
				Flag: blas.OffDiagonal,
				H:    [4]float64{0, -0.625, -0.9375, 0},
			},
		},
		{
			name:   "AbsQ1_LT_AbsQU__D2_Pos",
			d1:     2,
			d2:     3,
			x1:     5,
			y1:     8,
			d1Want: 2.3801652892561984,
			d2Want: 1.586776859504132,
			x1Want: 121.0 / 12,
			pWant: blas.DrotmParams{
				Flag: blas.Diagonal,
				H:    [4]float64{5.0 / 12, 0, 0, 0.625},
			},
		},
		{
			name:   "D1=D2_X1=X2",
			d1:     2,
			d2:     2,
			x1:     8,
			y1:     8,
			d1Want: 1,
			d2Want: 1,
			x1Want: 16,
			pWant: blas.DrotmParams{
				Flag: blas.Diagonal,
				H:    [4]float64{1, 0, 0, 1},
			},
		},
		{
			name:   "RD1_Big_RD2_Big_Flag_0",
			d1:     1600000000,
			d2:     800000000,
			x1:     8,
			y1:     7,
			d1Want: 68.96627824858757,
			d2Want: 34.483139124293785,
			x1Want: 45312,
			pWant: blas.DrotmParams{
				Flag: blas.Rescaling,
				H:    [4]float64{4096, -3584, 1792, 4096},
			},
		},
		{
			name:   "RD1_Big_RD2_Big_Flag_1",
			d1:     800000000,
			d2:     1600000000,
			x1:     8,
			y1:     7,
			d1Want: 57.6914092640818,
			d2Want: 28.8457046320409,
			x1Want: 47396.57142857142,
			pWant: blas.DrotmParams{
				Flag: blas.Rescaling,
				H:    [4]float64{2340.5714285714284, -4096, 4096, 4681.142857142857},
			},
		},
		{
			name:   "RD1_Big_RD2_Med_Flag_0",
			d1:     20000000,
			d2:     2,
			x1:     8,
			y1:     8,
			d1Want: 1.1920927762985347,
			d2Want: 1.9999998000000199,
			x1Want: 32768.0032768,
			pWant: blas.DrotmParams{
				Flag: blas.Rescaling,
				H:    [4]float64{4096, -1, 0.0004096, 1},
			},
		},
		{
			name:   "RD1_Big_RD2_Med_Flag_1",
			d1:     2,
			d2:     20000000000,
			x1:     8,
			y1:     80000000000,
			d1Want: 1192.0928955078125,
			d2Want: 2,
			x1Want: 3.2768e+14,
			pWant: blas.DrotmParams{
				Flag: blas.Rescaling,
				H:    [4]float64{4.096e-17, -1, 4096, 1e-10},
			},
		},
		{
			name:   "D1_Big_D2_Small_Flag_1",
			d1:     0.000000014,
			d2:     2000000000,
			x1:     0.000008,
			y1:     8000000,
			d1Want: 119.20928955078125,
			d2Want: 0.234881024,
			x1Want: 3.2768e+10,
			pWant: blas.DrotmParams{
				Flag: blas.Rescaling,
				H:    [4]float64{2.8671999999999997e-26, -0.000244140625, 4096, 2.44140625e-16},
			},
		},
		{
			name:   "RD1_Med_RD2_Big_Flag_0",
			d1:     2,
			d2:     20000000000,
			x1:     80000000,
			y1:     8,
			d1Want: 1.9998000199980002,
			d2Want: 1191.9736981379988,
			x1Want: 8.0008e+07,
			pWant: blas.DrotmParams{
				Flag: blas.Rescaling,
				H:    [4]float64{1, -0.0004096, 1000, 4096},
			},
		},
		{
			name:   "D1_Med_D2_Big_Flag_1",
			d1:     20000000000,
			d2:     0.4,
			x1:     80000000,
			y1:     80000000000000000,
			d1Want: 0.39999998000000103,
			d2Want: 1192.092835903171,
			x1Want: 8.0000004e+16,
			pWant: blas.DrotmParams{
				Flag: blas.Rescaling,
				H:    [4]float64{50, -4096, 1, 4.096e-06},
			},
		},
		{
			name:   "RD1_Med_RD2_Small_Flag_0",
			d1:     1.2,
			d2:     0.000000000045,
			x1:     2.7,
			y1:     8,
			d1Want: 1.1999999996049382,
			d2Want: 0.0007549747197514486,
			x1Want: 2.700000000888889,
			pWant: blas.DrotmParams{
				Flag: blas.Rescaling,
				H:    [4]float64{1, -0.0007233796296296296, 1.1111111111111111e-10, 0.000244140625},
			},
		},
		{
			name:   "RD1_Med_RD2_Small_Flag_1",
			d1:     1.2,
			d2:     0.000000000045,
			x1:     2.7,
			y1:     80000000000,
			d1Want: 0.0007549747199770676,
			d2Want: 1.19999999996355,
			x1Want: 1.9531250000593264e+07,
			pWant: blas.DrotmParams{
				Flag: blas.Rescaling,
				H:    [4]float64{0.0002197265625, -1, 0.000244140625, 3.375e-11},
			},
		},
		{
			name:   "D1_Small_D2_Big_Flag_1",
			d1:     120000000000000000,
			d2:     0.000000000012345,
			x1:     0.08,
			y1:     8000000000000,
			d1Want: 0.00010502490698765249,
			d2Want: 216.1836123957717,
			x1Want: 3.8516669198055897e+09,
			pWant: blas.DrotmParams{
				Flag: blas.Rescaling,
				H:    [4]float64{2.3731773997569866e+10, -1.6777216e+07, 0.000244140625, 1.6777216e-07},
			},
		},
		{
			name:   "RD1_Small_RD2_Med_Flag_0",
			d1:     0.0000000002,
			d2:     20,
			x1:     0.8,
			y1:     0.000000008,
			d1Want: 0.003355409645903541,
			d2Want: 19.99980000199998,
			x1Want: 0.000195314453125,
			pWant: blas.DrotmParams{
				Flag: blas.Rescaling,
				H:    [4]float64{0.000244140625, -1e-08, 0.24414062499999997, 1},
			},
		},
		{
			name:   "RD1_Small_RD2_Med_Flag_1",
			d1:     0.02,
			d2:     0.000000000004,
			x1:     0.008,
			y1:     8000000,
			d1Want: 6.710886366445568e-05,
			d2Want: 0.019999999900000003,
			x1Want: 1953.125009765625,
			pWant: blas.DrotmParams{
				Flag: blas.Rescaling,
				H:    [4]float64{0.0012207031250000002, -1, 0.000244140625, 1e-09},
			},
		},
		{
			// Values consistent with the low precision output posted at the  OpenBLAS issue.
			// See https://github.com/xianyi/OpenBLAS/issues/1452.
			name:   "OpenBLAS#1452",
			d1:     5.9e-8,
			d2:     5.960464e-8,
			x1:     1,
			y1:     150,
			d1Want: 0.9999559282289687,
			d2Want: 0.9898121986058326,
			x1Want: 0.03662270484346241,
			pWant: blas.DrotmParams{
				Flag: blas.Rescaling,
				H:    [4]float64{1.6110934624105326e-06, -0.000244140625, 0.000244140625, 1.6276041666666668e-06},
			},
		},

		// Test cases with values taken from BLAS/TESTING/dblat1.f in the reference.
		{
			name:   "netlib/BLAS/TESTING#1",
			d1:     0.1,
			d2:     0.3,
			x1:     1.2,
			y1:     0.2,
			d1Want: 12 / 130.0,
			d2Want: 36 / 130.0,
			x1Want: 1.3,
			pWant: blas.DrotmParams{
				Flag: blas.OffDiagonal,
				H:    [4]float64{0, -1 / 6.0, 0.5, 0},
			},
		},
		{
			name:   "netlib/BLAS/TESTING#2",
			d1:     0.7,
			d2:     0.2,
			x1:     0.6,
			y1:     4.2,
			d1Want: 14 / 75.0,
			d2Want: 49 / 75.0,
			x1Want: 4.5,
			pWant: blas.DrotmParams{
				Flag: blas.Diagonal,
				H:    [4]float64{0.5, 0, 0, 1 / 7.0},
			},
		},
		{
			name:   "netlib/BLAS/TESTING#3",
			d1:     0,
			d2:     0,
			x1:     0,
			y1:     0,
			d1Want: 0,
			d2Want: 0,
			x1Want: 0,
			pWant: blas.DrotmParams{
				Flag: blas.Identity,
				H:    [4]float64{0, 0, 0, 0},
			},
		},
		{
			name:   "netlib/BLAS/TESTING#4",
			d1:     4,
			d2:     -1,
			x1:     2,
			y1:     4,
			d1Want: 0,
			d2Want: 0,
			x1Want: 0,
			pWant: blas.DrotmParams{
				Flag: blas.Rescaling,
				H:    [4]float64{0, 0, 0, 0},
			},
		},
		{
			name:   "netlib/BLAS/TESTING#5",
			d1:     6e-10,
			d2:     2e-2,
			x1:     1e5,
			y1:     10,
			d1Want: 45e-11 * 4096 * 4096,
			d2Want: 15e-3,
			x1Want: 4e5 / (3 * 4096),
			pWant: blas.DrotmParams{
				Flag: blas.Rescaling,
				H:    [4]float64{1 / 4096.0, -1e-4, 1e4 / (3 * 4096), 1},
			},
		},
		{
			name:   "netlib/BLAS/TESTING#6",
			d1:     4e10,
			d2:     2e-2,
			x1:     1e-5,
			y1:     10,
			d1Want: 4e10 / (1.5 * 4096 * 4096),
			d2Want: 2e-2 / 1.5,
			x1Want: 6.144e-2,
			pWant: blas.DrotmParams{
				Flag: blas.Rescaling,
				H:    [4]float64{4096, -1e6, 5e-7 * 4096, 1},
			},
		},
		{
			name:   "netlib/BLAS/TESTING#7",
			d1:     2e-10,
			d2:     4e-2,
			x1:     1e5,
			y1:     10,
			d1Want: 4 / 150.0,
			d2Want: (2e-10 / 1.5) * 4096 * 4096,
			x1Want: 15,
			pWant: blas.DrotmParams{
				Flag: blas.Rescaling,
				H:    [4]float64{5e-5, -1 / 4096.0, 1, 1e4 / 4096},
			},
		},
		{
			name:   "netlib/BLAS/TESTING#8",
			d1:     2e10,
			d2:     4e-2,
			x1:     1e-5,
			y1:     10,
			d1Want: 4 / 150.0,
			d2Want: 2e10 / (1.5 * 4096 * 4096),
			x1Want: 15,
			pWant: blas.DrotmParams{
				Flag: blas.Rescaling,
				H:    [4]float64{5e5, -4096, 1, 4.096e-3},
			},
		},
	}

	for _, test := range tests {
		p, d1, d2, x1 := d.Drotmg(test.d1, test.d2, test.x1, test.y1)

		if p.Flag != test.pWant.Flag {
			t.Errorf("%v: unexpected flag, want %v, got %v", test.name, test.pWant.Flag, p.Flag)
		}
		for i, hi := range p.H {
			if !scalar.EqualWithinAbsOrRel(test.pWant.H[i], hi, tol, tol) {
				t.Errorf("%v: unexpected H, want %v, got %v", test.name, test.pWant.H, p.H)
				break
			}
		}
		if !scalar.EqualWithinAbsOrRel(d1, test.d1Want, tol, tol) {
			t.Errorf("%v: unexpected d1, want %v, got %v", test.name, test.d1Want, d1)
		}
		if !scalar.EqualWithinAbsOrRel(d2, test.d2Want, tol, tol) {
			t.Errorf("%v: unexpected d2, want %v, got %v", test.name, test.d2Want, d2)
		}
		if !scalar.EqualWithinAbsOrRel(x1, test.x1Want, tol, tol) {
			t.Errorf("%v: unexpected x1, want %v, got %v", test.name, test.x1Want, x1)
		}

		// Drotmg routines compute the components of a modified Givens transformation
		// matrix H that zeros the y-component of the resulting vector,
		//
		//  [x1; 0] := H[x1 sqrt(d1); y1 sqrt(d2)].
		//
		// Drotm performs a modified Givens rotation of points in the plane,
		//
		//  [x1; y1] := H[x1; y1].
		y := []float64{test.y1}
		d.Drotm(1, []float64{test.x1}, 1, y, 1, p)
		if d2 >= 0 {
			y[0] *= math.Sqrt(d2)
		}
		if !scalar.EqualWithinAbs(y[0], 0, 1e-10) {
			t.Errorf("%v: y1 not zeroed out, got %v", test.name, y[0])
		}
	}
}

type Droter interface {
	Drot(n int, x []float64, incX int, y []float64, incY int, c, s float64)
}

func DrotTest(t *testing.T, d Droter) {
	drot := d.Drot
	for _, c := range DoubleTwoVectorCases {
		for _, kind := range c.DrotCases {
			dCopyTwoTmp(c.X, c.XTmp, c.Y, c.YTmp)
			if c.Panic {
				f := func() { drot(c.N, c.XTmp, c.Incx, c.YTmp, c.Incy, kind.C, kind.S) }
				testpanics(f, c.Name, t)
				continue
			}
			drot(c.N, c.XTmp, c.Incx, c.YTmp, c.Incy, kind.C, kind.S)
			if !dSliceTolEqual(c.XTmp, kind.XAns) {
				t.Errorf("drot: x mismatch %v: expected %v, found %v", c.Name, kind.XAns, c.XTmp)
			}
			if !dSliceTolEqual(c.YTmp, kind.YAns) {
				t.Errorf("drot: y mismatch %v: expected %v, found %v", c.Name, kind.YAns, c.YTmp)
			}
		}
	}
}

type Drotmer interface {
	Drotm(n int, x []float64, incX int, y []float64, incY int, p blas.DrotmParams)
}

func DrotmTest(t *testing.T, d Drotmer) {
	drotm := d.Drotm
	for _, c := range DoubleTwoVectorCases {
		for _, kind := range c.DrotmCases {
			dCopyTwoTmp(c.X, c.XTmp, c.Y, c.YTmp)
			if c.Panic {
				f := func() { drotm(c.N, c.XTmp, c.Incx, c.YTmp, c.Incy, kind.P) }
				testpanics(f, c.Name+", "+kind.Name, t)
				continue
			}
			drotm(c.N, c.XTmp, c.Incx, c.YTmp, c.Incy, kind.P)
			if !dSliceTolEqual(c.XTmp, kind.XAns) {
				t.Errorf("drotm: mismatch %v: expected %v, found %v", c.Name, kind.XAns, c.XTmp)
			}
			if !dSliceTolEqual(c.YTmp, kind.YAns) {
				t.Errorf("drotm: mismatch %v: expected %v, found %v", c.Name, kind.YAns, c.YTmp)
			}
		}
	}
}

type Dscaler interface {
	Dscal(n int, alpha float64, x []float64, incX int)
}

func DscalTest(t *testing.T, blasser Dscaler) {
	dscal := blasser.Dscal
	for _, c := range DoubleOneVectorCases {
		for _, kind := range c.DscalCases {
			xTmp := make([]float64, len(c.X))
			copy(xTmp, c.X)
			if c.Panic {
				f := func() { dscal(c.N, kind.Alpha, xTmp, c.Incx) }
				testpanics(f, c.Name, t)
				continue
			}
			dscal(c.N, kind.Alpha, xTmp, c.Incx)
			if !dSliceTolEqual(xTmp, kind.Ans) {
				t.Errorf("dscal: mismatch %v, %v: expected %v, found %v", c.Name, kind.Name, kind.Ans, xTmp)
			}
		}
	}
}
