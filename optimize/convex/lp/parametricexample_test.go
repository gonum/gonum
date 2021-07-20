// Copyright ©2016 The Gonum authors. all rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lp_test

import (
	"fmt"
	"log"
	"math"
	"sort"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/optimize/convex/lp"
	"gonum.org/v1/gonum/stat"
)

type IntFloat64Pair struct {
	x []int
	y []float64
}

func (p IntFloat64Pair) Len() int {
	return len(p.x)
}

func (p IntFloat64Pair) Swap(i, j int) {
	p.x[i], p.x[j] = p.x[j], p.x[i]
	p.y[i], p.y[j] = p.y[j], p.y[i]
}

func (p IntFloat64Pair) Less(i, j int) bool {
	if math.Abs(p.y[i]) > math.Abs(p.y[j]) {
		return true
	}
	return false
}

func ExampleParametric() {
	rnd := rand.New(rand.NewSource(0))
	c := []float64{-1, -2, 0, 0}
	a := mat.NewDense(2, 4, []float64{-1, 2, 1, 0, 3, 1, 0, 1})
	b := []float64{4, 9}

	opt, x, _, err := lp.Parametric(c, a, b, 0, nil, true, rnd)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("opt: %v\n", opt)
	fmt.Printf("x: %v\n", x)
	// Output:
	// opt: -8
	// x: [2 3 0 0]
}

func ExampleParametricWithTrace() {
	// The Dantzig selector as discussed in Candes, E., & Tao, T. (2007).
	// "The Dantzig selector: Statistical estimation when p is much larger than n".
	// The annals of Statistics, 2313-2351.
	const n, m, d, s = 150, 100, 250, 8
	rnd := rand.New(rand.NewSource(0))

	idxdata := make([]int, s)
	inIdx := make(map[int]struct{})
	for i := 0; i < s; i++ {
		idx := rnd.Intn(d)
		_, ok := inIdx[idx]
		for ok {
			idx = rnd.Intn(d)
			_, ok = inIdx[idx]
		}
		idxdata[i] = idx
		inIdx[idx] = struct{}{}
	}

	betadata := make([]float64, d)
	nonzero := make([]float64, s)
	for i, v := range idxdata {
		if a, b := rnd.Float64(), 1+rnd.NormFloat64(); a < 0.5 {
			nonzero[i] = b
			betadata[v] = b
		} else {
			nonzero[i] = -b
			betadata[v] = -b
		}
	}
	beta := mat.NewVecDense(d, betadata)

	xdata := make([]float64, n*d)
	for i := range xdata {
		xdata[i] = rnd.NormFloat64()
	}
	X := mat.NewDense(n, d, xdata)

	// standardize columns
	coltmp := make([]float64, n)
	for j := 0; j < d; j++ {
		mat.Col(coltmp, j, X)
		mean, std := stat.MeanStdDev(coltmp, nil)
		tau := 1 / std
		for i, v := range coltmp {
			coltmp[i] = tau * (v - mean)
		}
		X.SetCol(j, coltmp)
	}
	trainX := X.Slice(0, m, 0, d)
	validateX := X.Slice(m, n, 0, d)

	epsdata := make([]float64, n)
	for i := range epsdata {
		epsdata[i] = rnd.NormFloat64()
	}
	eps := mat.NewVecDense(n, epsdata)

	ydata := make([]float64, n)
	y := mat.NewVecDense(n, ydata)
	trainY := mat.NewVecDense(m, ydata[:m])
	validateY := mat.NewVecDense(n-m, ydata[m:n])
	y.MulVec(X, beta)
	y.AddVec(y, eps)

	b := make([]float64, 2*d)
	XTy := mat.NewVecDense(d, b[0:d])
	XTy.MulVec(trainX.T(), trainY)
	XTy.ScaleVec(1.0/m, XTy)

	adata := make([]float64, 8*d*d)
	a := mat.NewDense(2*d, 4*d, adata)
	a11 := a.Slice(0, d, 0, d).(*mat.Dense)
	a12 := a.Slice(0, d, d, 2*d).(*mat.Dense)
	a21 := a.Slice(d, 2*d, 0, d).(*mat.Dense)
	a22 := a.Slice(d, 2*d, d, 2*d).(*mat.Dense)
	a11.Mul(trainX.T(), trainX)
	a11.Scale(1.0/m, a11)
	a12.Scale(-1, a11)
	a21.Copy(a12)
	a22.Copy(a11)
	for i := 0; i < 2*d; i++ {
		adata[i*4*d+2*d+i] = 1
	}

	c := make([]float64, 4*d)
	for i := 0; i < 2*d; i++ {
		c[i] = 1
	}

	for i := 0; i < d; i++ {
		b[i+d] = -b[i]
	}

	cbar := make([]float64, 4*d)
	bbar := make([]float64, 2*d)
	for i := range bbar {
		bbar[i] = 1
	}
	tol := math.Sqrt(math.Log(d) / m)
	optTr, err := lp.ParametricWithTrace(c, a, b, cbar, bbar, tol, nil, rnd)
	if err != nil {
		panic(err)
	}

	iter := len(optTr.Lambda)
	idx := 0

	// check SSE on validation set
	lo := math.Inf(1)
	yhat_data := make([]float64, n-m)
	yhat := mat.NewVecDense(n-m, yhat_data)
	betahat_data := make([]float64, d)
	betahat := mat.NewVecDense(d, betahat_data)
	for i := 0; i < iter; i++ {
		for j := range betahat_data {
			betahat_data[j] = 0
		}
		basicIdxs := optTr.Idx[2*d*i : 2*d*(i+1)]
		optx := optTr.X[2*d*i : 2*d*(i+1)]
		for i, v := range basicIdxs {
			if v < d {
				betahat_data[v] += optx[i]
			} else if v < 2*d {
				betahat_data[v-d] -= optx[i]
			}
		}
		yhat.MulVec(validateX, betahat)
		yhat.SubVec(yhat, validateY)
		if sse := floats.Dot(yhat_data, yhat_data); sse < lo {
			lo = sse
			idx = i
		}
	}

	finalX := optTr.X[2*d*idx : 2*d*(idx+1)]
	finalBasicIdxs := optTr.Idx[2*d*idx : 2*d*(idx+1)]
	xopt := make([]float64, 4*d)
	for i, v := range finalBasicIdxs {
		xopt[v] = finalX[i]
	}

	// check for correctness of solution
	var bCheck mat.VecDense
	bCheck.MulVec(a, mat.NewVecDense(len(xopt), xopt))
	tmp := make([]float64, len(b))
	floats.AddScaledTo(tmp, b, optTr.Lambda[idx], bbar)
	if !mat.EqualApprox(&bCheck, mat.NewVecDense(len(tmp), tmp), 1e-10) {
		log.Panic("No error in primal but solution infeasible")
	}

	betamap := make(map[int]float64)
	for i, v := range finalBasicIdxs {
		if v < d {
			betamap[v] += finalX[i]
		} else if v < 2*d {
			betamap[v-d] -= finalX[i]
		}
	}

	keys := make([]int, 0, len(betamap))
	values := make([]float64, 0, len(betamap))
	for k, v := range betamap {
		keys = append(keys, k)
		values = append(values, v)
	}

	p := IntFloat64Pair{
		x: keys,
		y: values,
	}
	sort.Sort(p)

	p.x = idxdata
	p.y = nonzero
	sort.Sort(p)

	fmt.Println("in order of decreasing absolute value")
	fmt.Printf("true idx: %v\n", idxdata)
	fmt.Printf("est idx:  %v\n", keys)
	// Output:
	// in order of decreasing absolute value
	// true idx: [186 244 32 227 54 229 33 122]
	// est idx:  [186 244 32 227 54 229 122 33 145 49 165 202 75 194 114 147]
}
