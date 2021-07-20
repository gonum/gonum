// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package lp

import (
	"fmt"
	"math"
	"time"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
)

// AffineScaling solves using the affine scaling method
// This is the method that was famously re-discovered in 1986
// Dikin 1967
// Vanderbei, Meketon, and Freedman 1986
func AffineScaling(c []float64, A mat.Matrix, b []float64, tol float64, initialPoint []float64, maxIterations int, baseBeta float64, fractionRecenter float64, fractionExtendSearch float64) (optF float64, optX []float64, err error) {
	// set up sizes and vectors
	var newPoint []float64
	nVars := len(c)
	rA, cA := A.Dims()
	bVec := mat.NewVecDense(rA, b)
	cVec := mat.NewVecDense(cA, c)
	currentPoint := initialPoint

	// if we didnt get a starting point find one
	if currentPoint == nil {
		basicIdxs, _, xb, err := findInitialBasic(A, b)
		if err != nil {
			return math.NaN(), nil, err
		}
		currentPoint = make([]float64, nVars)
		for i, v := range basicIdxs {
			currentPoint[v] = xb[i]
		}

		// feasibility check
		fCheck := mat.NewVecDense(rA, nil)
		curVec := mat.NewVecDense(cA, currentPoint)
		fCheck.MulVec(A, curVec)
		checkDiff := mat.NewVecDense(rA, nil)
		checkDiff.SubVec(bVec, fCheck)
		fmt.Println("init check: ", checkDiff)
		checkSize := mat.Norm(checkDiff, 2.0)
		if checkSize > tol {
			panic("initial point seems infeasible")
			//return math.NaN(), nil, ErrInfeasible
		}
	}

	// check that initial point is in the feasible region
	initCheck := mat.NewVecDense(rA, nil)
	initCheck.MulVec(A, mat.NewVecDense(nVars, initialPoint))
	initDistance := mat.NewVecDense(rA, nil)
	initDistance.SubVec(bVec, initCheck)
	initCheckSize := mat.Norm(initDistance, 2.0)
	fmt.Println("init check size: ", initCheckSize)

	// compute where we are
	oldFValue := mat.Dot(cVec, mat.NewVecDense(nVars, currentPoint))
	curVec := mat.NewVecDense(cA, currentPoint)
	newFValue := oldFValue

	betaScale := baseBeta
	newBeta := 0.0
	stepSizes := make([]float64, 0)
	thisSize := 0.0
	minEle := 0.0
	infeasCount := 0.0
	recenter := false
	optCondMin := 0.0
	dSize := 0.0
	for i := 0; i < maxIterations; i++ {
		// and if its feasible
		// feasibility check
		fCheck := mat.NewVecDense(rA, nil)
		// for old
		fCheck.MulVec(A, curVec)
		checkDiff := mat.NewVecDense(rA, nil)
		checkDiff.SubVec(bVec, fCheck)
		//fmt.Println("check: ", checkDiff)
		//oldCheckSize := mat.Norm(checkDiff, 2.0)
		//fmt.Println("old check size: ", oldCheckSize)

		// occasionally recenter
		if rand.Float64() < fractionRecenter {
			recenter = true
		}

		// step
		newPoint, newFValue, newBeta, thisSize, minEle, optCondMin, dSize = affineScalingStep(currentPoint, c, A, betaScale, recenter)
		_ = newBeta
		//fmt.Println("thisSize: ", thisSize)
		//fmt.Println("minEle: ", minEle)

		// compute where we are now
		fDelta := oldFValue - newFValue
		newVec := mat.NewVecDense(nVars, newPoint)
		fCheck.MulVec(A, newVec)
		checkDiff.SubVec(bVec, fCheck)
		newCheckSize := mat.Norm(checkDiff, 2.0)

		// how far did we step?
		diffVec := mat.NewVecDense(nVars, nil)
		diffVec.SubVec(curVec, newVec)
		stepSize := mat.Norm(diffVec, 2.0)
		stepSizes = append(stepSizes, stepSize)

		// and check the stopping conditions now
		fmt.Println(time.Now(), " , old/new fValue: ", oldFValue, ",", newFValue, "  ", i)
		//fmt.Println("thisSize / minEle: ", thisSize, ",", minEle)
		//fmt.Println("optCondMin: ", optCondMin)

		// track how many consecutive infeasible entries
		// so we can cut beta more quickly
		if newCheckSize < tol {
			infeasCount = 0
		} else {
			infeasCount += 1
		}

		if fDelta > 0 && newCheckSize < tol {
			// f still dropping and still feasible
			// just keep chugging
			// this is the base iterating case
			//fmt.Println("step ", time.Now())
		} else if optCondMin > -tol && newCheckSize < tol {
			// optimal condition: c - AtV >= 0
			// this looks like the good one
			fmt.Println("van optimal 1")
			return newFValue, newPoint, nil
		} else if thisSize < tol && minEle >= -tol && newCheckSize < tol {
			// optimal condition from vanderbrei 1986
			fmt.Println("van optimal 2")
			return newFValue, newPoint, nil
		} else if dSize < tol && newCheckSize < tol {
			// optimal condition dk small
			fmt.Println("van optimal 3")
			return newFValue, newPoint, nil
		} else if stepSize < tol && minEle >= -tol && newCheckSize < tol {
			// my condition
			fmt.Println("my optimal")
			return newFValue, newPoint, nil
		} else if newCheckSize > tol {
			//fmt.Println("new point infeasible")
			// infeasible point, reduce step size
			betaScale = betaScale * math.Pow(10, -(1+infeasCount))
			continue
		} else if math.IsNaN(stepSize) {
			fmt.Println("NaN step, returning error")
			return oldFValue, currentPoint, ErrUnbounded
		}
		infeasCount = 0
		betaScale = baseBeta
		currentPoint = newPoint
		//fmt.Println("updating point")
		curVec = mat.NewVecDense(cA, currentPoint)
		oldFValue = newFValue
		recenter = false
		// if we are making progress do not stop
		if fDelta > 0 && i > 0 {
			i -= int(float64(maxIterations) / fractionExtendSearch)
			if i < 0 {
				i = 0
			}
		}
	}
	// FIXME - can be unbounded here too?
	// it is for sure feasible
	fmt.Println("done all iters")
	return oldFValue, currentPoint, nil
}

// the solver step through the feasible region
// variable names are intended to match standard notation
func affineScalingStep(startPoint []float64, cIn []float64, A mat.Matrix, betaScale float64, recenter bool) ([]float64, float64, float64, float64, float64, float64, float64) {
	nVars := len(startPoint)
	rA, cA := A.Dims()

	// ones
	ones := mat.NewVecDense(nVars, nil)
	for i := 0; i < nVars; i++ {
		ones.SetVec(i, 1.0)
	}
	// set up variables
	c := mat.NewVecDense(nVars, cIn)
	// diagonal matrix for current point
	D := mat.NewDiagDense(nVars, startPoint)

	// scaled condition matrix
	Ak := mat.NewDense(rA, cA, nil)
	Ak.Mul(A, D)

	// scaled objective function
	ck := mat.NewVecDense(nVars, nil)
	ck.MulVec(D, c)

	// scaled equality targets
	wk := mat.NewVecDense(rA, nil)
	wk.MulVec(Ak, ck)

	// compute Ak*Ak-transpose
	AkT := Ak.T()
	systemMat := mat.NewDense(rA, rA, nil)
	systemMat.Mul(Ak, AkT)

	// FIXME - this is the slow part
	// factor out / plug-in specialized versions
	// symMat is symmetric positive definite
	var chol mat.Cholesky
	symMat := mat.NewSymDense(rA, systemMat.RawMatrix().Data)
	cholOk := chol.Factorize(symMat)
	// solve into vk
	vk := mat.NewVecDense(rA, nil)
	if cholOk {
		// FIXME - why does this sometimes fail? likely a numeric thing
		// probably also related to occasional jumps outside the feasible region
		chol.SolveVecTo(vk, wk)
	} else {
		vk.SolveVec(systemMat, wk)
		// the dumb version is:
		//sysInv := mat.NewDense(rA, rA, nil)
		//sysInv.Inverse(systemMat)
		//vk.MulVec(sysInv, wk)
	}

	// compute step direction
	gk := mat.NewVecDense(nVars, nil)
	gk.MulVec(AkT, vk)
	dk := mat.NewVecDense(nVars, nil)
	if recenter {
		dk.SubVec(ones, gk)
	} else {
		dk.SubVec(gk, ck)
	}

	rk := mat.NewVecDense(nVars, nil)
	rk.SubVec(ck, gk)
	sizeVec := mat.NewVecDense(nVars, nil)
	sizeVec.MulVec(D, rk)
	thisSum := mat.Sum(sizeVec)
	minEle := mat.Min(sizeVec)

	// and scale it to (near) the boundary
	maxDi := 0.0
	for i := 0; i < nVars; i++ {
		thisDi := dk.At(i, 0)
		if thisDi < 0.0 {
			neg := -thisDi
			if neg > maxDi {
				maxDi = math.Pow(neg, 1.0)
			}
		}
	}
	// compute scale
	beta := 0.0
	if maxDi != 0.0 {
		beta = betaScale / maxDi
	}

	// xform step is z = 1 + beta*dk
	betaDk := mat.NewVecDense(nVars, nil)
	betaDk.ScaleVec(beta, dk)
	z := mat.NewVecDense(nVars, nil)
	z.AddVec(ones, betaDk)

	// step to new point
	newPoint := mat.NewVecDense(nVars, nil)
	newPoint.MulVec(D, z)

	// reevaluate function
	funcValue := mat.Dot(c, newPoint)

	// and optimality condition
	optCondVec := mat.NewVecDense(nVars, nil)
	optCondVec.SubVec(newPoint, gk)
	optCondMin := mat.Min(optCondVec)

	// size of direction vec
	dSize := mat.Norm(dk, 2.0)

	// return
	return newPoint.RawVector().Data, funcValue, beta, thisSum, minEle, optCondMin, dSize
}

// AffineScalingFromInitialBasic solves from an initial basic simplex solution
// this exists to merge testing with simplex
func AffineScalingFromInitialBasic(c []float64, A mat.Matrix, b []float64, tol float64, initialBasic []int, maxIterations int, baseBeta float64, fractionRecenter float64, fractionExtendSearch float64) (float64, []float64, error) {
	_, _, basicIdxs, _, ab, xb, err := simplexPresolve(initialBasic, c, A, b, tol)

	if initialBasic != nil {
		extractColumns(ab, A, initialBasic)
		err = initializeFromBasic(xb, ab, b)
		if err != nil {
			panic(err)
		}
		copy(basicIdxs, initialBasic)

		initialPoint := make([]float64, len(c))
		for i, v := range initialBasic {
			initialPoint[v] = xb[i]
		}
		return AffineScaling(c, A, b, tol, initialPoint, maxIterations, baseBeta, fractionRecenter, fractionExtendSearch)
	}
	return AffineScaling(c, A, b, tol, nil, maxIterations, baseBeta, fractionRecenter, fractionExtendSearch)
}
