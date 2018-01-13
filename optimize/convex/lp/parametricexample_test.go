// Copyright ©2016 The Gonum Authors. All rights reserved.
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
	if p.x[i] < p.x[j] {
		return true
	} else if p.x[i] == p.x[j] && p.y[i] < p.y[j] {
		return true
	}
	return false
}

func ExampleParametric() {
	rnd := rand.New(rand.NewSource(0))
	c := []float64{-1, -2, 0, 0}
	A := mat.NewDense(2, 4, []float64{-1, 2, 1, 0, 3, 1, 0, 1})
	b := []float64{4, 9}

	opt, x, _, err := lp.Parametric(c, A, b, 0, nil, true, rnd)
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
	// The Annals of Statistics, 2313-2351.
	const n, d, s = 100, 250, 8
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
	for _, v := range idxdata {
		if a, b := rnd.Float64(), 1+rnd.NormFloat64(); a < 0.5 {
			betadata[v] = b
		} else {
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

	epsdata := make([]float64, n)
	for i := range epsdata {
		epsdata[i] = rnd.NormFloat64()
	}
	eps := mat.NewVecDense(n, epsdata)

	ydata := make([]float64, n)
	y := mat.NewVecDense(n, ydata)
	y.MulVec(X, beta)
	y.AddVec(y, eps)

	b := make([]float64, 2*d)
	XTy := mat.NewVecDense(d, b[0:d])
	XTy.MulVec(X.T(), y)
	XTy.ScaleVec(1.0/n, XTy)

	adata := make([]float64, 8*d*d)
	A := mat.NewDense(2*d, 4*d, adata)
	A11 := A.Slice(0, d, 0, d).(*mat.Dense)
	A12 := A.Slice(0, d, d, 2*d).(*mat.Dense)
	A21 := A.Slice(d, 2*d, 0, d).(*mat.Dense)
	A22 := A.Slice(d, 2*d, d, 2*d).(*mat.Dense)
	A11.Mul(X.T(), X)
	A11.Scale(1.0/n, A11)
	A12.Scale(-1, A11)
	A21.Copy(A12)
	A22.Copy(A11)
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
	tol := math.Sqrt(math.Log(d) / n)
	secs := time.Now().UnixNano()
	optTr, err := lp.ParametricWithTrace(c, A, b, cbar, bbar, tol, nil, rnd)
	if err != nil {
		panic(err)
	}

	iter := len(optTr.Lambda)
	idx := 0
	for i := 0; i < iter; i++ {
		tmp := 0
		basicIdxs := optTr.Idx[2*d*i : 2*d*(i+1)]
		for _, v := range basicIdxs {
			if v < 2*d {
				tmp++
			}
		}
		if tmp > s {
			idx = i - 1
			break
		}
	}

	finalX := optTr.X[2*d*idx : 2*d*(idx+1)]
	finalBasicIdxs := optTr.Idx[2*d*idx : 2*d*(idx+1)]
	xopt := make([]float64, 4*d)
	for i, v := range finalBasicIdxs {
		xopt[v] = finalX[i]
	}

	var bCheck mat.VecDense
	bCheck.MulVec(A, mat.NewVecDense(len(xopt), xopt))
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
	fmt.Printf("num selected: %v\n", len(betamap))
	fmt.Printf("idx: %v\n", keys)
	fmt.Printf("est: %v\n", values)
	// Output:
	// num selected: 8
	// idx: [32 33 54 122 186 227 229 244]
	// est: [-1.3767259886061176 -0.20328979569981956 0.7595477723815178 -0.20771185072465692 2.139833618293029 -0.9784601263043735 0.7263435538015337 -2.0692216936117083]
}
