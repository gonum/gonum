package lp

import (
	"errors"
	"golang.org/x/exp/rand"
	"math"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

var (
	ErrDegenerate = errors.New("lp: current perturbation is degenerate")
)

const (
	// absZeroTol is the absolute tolerance on testing if the elements of an
	// update vector are zero.
	absZeroTol = 1e-12
	// relZeroTol is the tolerance on testing if the elements of an update
	// vector are zero relative to the current values.
	relZeroTol = 1e-12
	// xZeroTol is the tolerance on testing if the elements of the
	// optimal solution vector are zero (degeneracy).
	xZeroTol = 1e-14
	// swapCap is the default max capacity for the Swap object
	swapCap = 25
	// swapCondTol is the tolerence on the condition bound of the swap matrix;
	// this is lower than mat.ConditionTolerance as to avoid the accumalation
	// of "decimal dust" through repeated use of an ill-conditioned operator
	swapCondTol = 1e8
)

// Trace stores the per-iteration output of the parametric simplex method.
type Trace struct {
	Lambda, X []float64
	Idx       []int
}

// Parametric solves a linear program in standard form using the parametric
// simplex method, a primal-dual algorithm. A description of the algorithm can
// be found in Pang, H. et al. (2017). "Parametric Simplex Method for Sparse Learning".
// Advances in Neural Information Processing Systems(30): p. 187-196.
func Parametric(c []float64, A mat.Matrix, b []float64, tol float64, initialBasic []int, dense bool, rnd *rand.Rand) (optF float64, optX []float64, basicIdxs []int, err error) {
	ans, x, idxs, _, err := parametric(initialBasic, c, A, b, nil, nil, tol, false, dense, rnd)
	if dense {
		return ans, x, nil, err
	}
	return ans, x, idxs, err
}

// ParametricWithTrace solves a linear program in standard form using the parametric
// simplex method and returns the output of each iteration in a Trace object. This
// function is useful for homotopy in L1 regularized minimization problems. Note that
//  minimize	L(θ) + λ‖θ‖_1
// can be written equivalently as
//  minimize ‖θ‖_1
//  s.t. 		‖∇L(θ)‖_∞ ≤ λ.
// Either fomulation can both be expressed as a linear program when L(θ)
// is piecewise linear or quadratic, respectively. Therefore the parametric simplex
// method allows for recovery of the entire regularization path of at the cost of
// solving a single linear program.
func ParametricWithTrace(c []float64, A mat.Matrix, b, cbar, bbar []float64, tol float64, initialBasic []int, rnd *rand.Rand) (optTr *Trace, err error) {
	_, _, _, tr, err := parametric(initialBasic, c, A, b, cbar, bbar, tol, true, false, rnd)
	return tr, err
}

func parametric(initialBasic []int, c []float64, A mat.Matrix, b,
	cbar, bbar []float64, tol float64, trace, dense bool, rnd *rand.Rand) (float64, []float64, []int, *Trace, error) {
	err := verifyInputs(initialBasic, c, A, b)
	if err != nil {
		if err == ErrUnbounded {
			return math.Inf(-1), nil, nil, nil, ErrUnbounded
		}
		return math.NaN(), nil, nil, nil, err
	}
	verifyPerturbations(c, b, cbar, bbar)
	m, n := A.Dims()

	var basicIdxs, nonBasicIdxs []int // The indices of the non-zero and zero x values, respectively.
	var ab, an *mat.Dense             // Basic and non-basic subsets of A
	var xb, xbbar, xtot, dx []float64 // The non-zero elements of primal dictionary
	var zn, znbar, ztot, dz []float64 // The non-zero elements of dual dictionary
	var cb, cn []float64
	var lu mat.LU
	var swap *Swap
	var tr *Trace

	if trace {
		tr = &Trace{
			Lambda: make([]float64, 0, 4),
			X:      make([]float64, 0, 4*m),
			Idx:    make([]int, 0, 4*m),
		}
	}

	if initialBasic != nil {
		// InitialBasic supplied. Panic if incorrect length or infeasible.
		if len(initialBasic) != m {
			panic("lp: incorrect number of initial vectors")
		}
		ab = mat.NewDense(m, m, nil)
		extractColumns(ab, A, initialBasic)
		xb = make([]float64, m)
		lu.Factorize(ab)
		if lu.Cond() > mat.ConditionTolerance {
			panic("lp: initalBasic produces singular matrix")
		}
		basicIdxs = make([]int, m)
		copy(basicIdxs, initialBasic)
	} else {
		// No initial basis supplied. Find non-singular submatrix
		ab = mat.NewDense(m, m, nil)
		basicIdxs = findLinearlyIndependent(A)
		if len(basicIdxs) != m {
			return math.NaN(), nil, nil, nil, ErrSingular
		}
		extractColumns(ab, A, basicIdxs)
		lu.Factorize(ab)
		if lu.Cond() > mat.ConditionTolerance {
			panic("lp: initalBasic produces singular matrix")
		}
	}

	// nonBasicIdxs is the set of nonbasic variables.
	nonBasicIdxs = make([]int, 0, n-m)
	inBasic := make(map[int]struct{})
	for _, v := range basicIdxs {
		inBasic[v] = struct{}{}
	}
	for i := 0; i < n; i++ {
		_, ok := inBasic[i]
		if !ok {
			nonBasicIdxs = append(nonBasicIdxs, i)
		}
	}

	// cb is the subset of c for the basic variables. an and cn
	// are the equivalents to ab and cb but for the nonbasic variables.
	cb = make([]float64, m)
	for i, idx := range basicIdxs {
		cb[i] = c[idx]
	}
	cn = make([]float64, n-m)
	for i, idx := range nonBasicIdxs {
		cn[i] = c[idx]
	}

	an = mat.NewDense(m, n-m, nil)
	extractColumns(an, A, nonBasicIdxs)

	xb = make([]float64, m)
	xbbar = make([]float64, len(xb))
	xtot = make([]float64, len(xb))
	dx = make([]float64, len(xb))
	zn = make([]float64, n-m)
	znbar = make([]float64, n-m)
	ztot = make([]float64, len(zn))
	dz = make([]float64, len(zn))
	swap = &Swap{Dim: m}

	xbVec := mat.NewVecDense(m, xb)
	xbbarVec := mat.NewVecDense(m, xbbar)
	dxVec := mat.NewVecDense(len(dx), dx)
	znVec := mat.NewVecDense(n-m, zn)
	znbarVec := mat.NewVecDense(n-m, znbar)
	dzVec := mat.NewVecDense(len(dz), dz)

	// to produce perturbations of appropriate scale
	bnorm := floats.Norm(b, math.Inf(1))
	if bnorm == 0 {
		bnorm = 1
	}
	cnorm := floats.Norm(c, math.Inf(1))
	if cnorm == 0 {
		cnorm = 1
	}

	// initialize xb and xbbar
	lu.SolveVec(xbVec, false, mat.NewVecDense(m, b))
	if bbar != nil {
		if len(bbar) != 0 {
			lu.SolveVec(xbbarVec, false, mat.NewVecDense(m, bbar))
		}
	} else {
		// initialize xbbar directly to guarantee primal feasibility
		for i := 0; i < m; i++ {
			xbbar[i] = rnd.Float64() * bnorm
		}
	}

	tmp := make([]float64, m)
	tmpVec := mat.NewVecDense(m, tmp)

	// initialize zn and znbar
	err = lu.SolveVec(tmpVec, true, mat.NewVecDense(m, cb))
	if err != nil {
		return math.NaN(), nil, nil, nil, err
	}
	znVec.MulVec(an.T(), tmpVec)
	floats.SubTo(zn, cn, zn)
	if cbar != nil {
		if len(cbar) != 0 {
			cbbar := make([]float64, m)
			for i, idx := range basicIdxs {
				cbbar[i] = cbar[idx]
			}
			cnbar := make([]float64, n-m)
			for i, idx := range nonBasicIdxs {
				cnbar[i] = cbar[idx]
			}
			err = lu.SolveVec(tmpVec, true, mat.NewVecDense(m, cbbar))
			if err != nil {
				return math.NaN(), nil, nil, nil, err
			}
			znbarVec.MulVec(an.T(), tmpVec)
			floats.SubTo(znbar, cnbar, znbar)
		}
	} else {
		// initialize znbar directly to guarantee dual feasibility
		for i := 0; i < n-m; i++ {
			znbar[i] = rnd.Float64() * cnorm
		}
	}

	var enIdx, lvIdx int
	oldLambdaMin := math.Inf(1)
	boolean, index, lambdaMin, lambdaMax := lambdaFn(xb, xbbar, zn, znbar)

	for {
		floats.AddScaledTo(xtot, xb, lambdaMin, xbbar)
		if trace {
			// save lambdaMin, xtot, and basicIdxs to caches
			tr.Lambda = append(tr.Lambda, lambdaMin)
			tr.X = append(tr.X, xtot...)
			tr.Idx = append(tr.Idx, basicIdxs...)
		}

		if lambdaMax < lambdaMin || lambdaMin >= oldLambdaMin {
			// current perturbation creates a degenerate lp
			if trace {
				return math.NaN(), nil, nil, tr, ErrDegenerate
			} else {
				// if no trace required, choose new perterbations;
				// probabality that new system is also degenerate is
				// nearly zero
				for i := 0; i < m; i++ {
					xbbar[i] = rnd.Float64() * bnorm
				}
				for i := 0; i < n-m; i++ {
					znbar[i] = rnd.Float64() * cnorm
				}
				boolean, index, lambdaMin, lambdaMax = lambdaFn(xb, xbbar, zn, znbar)
			}
		}

		if lambdaMin <= tol {
			// bailout
			break
		}

		if boolean {
			enIdx = index
			err = computePrimal(dxVec, an, enIdx, &lu, swap)
			if err != nil {
				break
			}
			lvIdx, err = selectIdx(xtot, dx)
			if err != nil {
				break
			}
			err = computeDual(dzVec, an, lvIdx, &lu, swap)
			if err != nil {
				break
			}
		} else {
			lvIdx = index
			err = computeDual(dzVec, an, lvIdx, &lu, swap)
			if err != nil {
				break
			}
			floats.AddScaledTo(ztot, zn, lambdaMin, znbar)
			enIdx, err = selectIdx(ztot, dz)
			if err != nil {
				break
			}
			err = computePrimal(dxVec, an, enIdx, &lu, swap)
			if err != nil {
				break
			}
		}
		// update primal and dual dictionaries
		t := xb[lvIdx] / dx[lvIdx]
		tbar := xbbar[lvIdx] / dx[lvIdx]
		s := zn[enIdx] / dz[enIdx]
		sbar := znbar[enIdx] / dz[enIdx]
		floats.AddScaled(xb, -t, dx)
		floats.AddScaled(xbbar, -tbar, dx)
		floats.AddScaled(zn, -s, dz)
		floats.AddScaled(znbar, -sbar, dz)
		xb[lvIdx] = t
		xbbar[lvIdx] = tbar
		zn[enIdx] = s
		znbar[enIdx] = sbar
		basicIdxs[lvIdx], nonBasicIdxs[enIdx] = nonBasicIdxs[enIdx], basicIdxs[lvIdx]

		// perform implicit swap
		lu.SolveVec(tmpVec, false, an.ColView(enIdx).(*mat.VecDense))
		swap.SolveVec(tmpVec, false, tmpVec)
		swap.Append(tmp, lvIdx)
		tmpCol1 := mat.Col(nil, lvIdx, ab)
		tmpCol2 := mat.Col(nil, enIdx, an)
		ab.SetCol(lvIdx, tmpCol2)
		an.SetCol(enIdx, tmpCol1)

		// if condition bound is large, refactorize matrix
		if swap.Cond() > swapCondTol || swap.Len() >= swapCap {
			lu.Factorize(ab)
			if lu.Cond() > mat.ConditionTolerance {
				panic("lp: basicIdxs produces singular matrix")
			}
			swap.Reset()
		}

		// update lambdaMin
		oldLambdaMin = lambdaMin
		boolean, index, lambdaMin, lambdaMax = lambdaFn(xb, xbbar, zn, znbar)
	}

	var xopt []float64

	for i, v := range basicIdxs {
		cb[i] = c[v]
		if xb[i] < xZeroTol {
			xb[i] = 0
		}
	}
	opt := floats.Dot(cb, xb)

	if dense {
		xopt = make([]float64, n)
		for i, v := range basicIdxs {
			xopt[v] = xb[i]
		}
	} else {
		xopt = xb
	}

	return opt, xopt, basicIdxs, tr, err
}

func computePrimal(dxVec *mat.VecDense, an *mat.Dense, enIdx int, lu *mat.LU, swap *Swap) error {
	err := lu.SolveVec(dxVec, false, an.ColView(enIdx).(*mat.VecDense))
	if err != nil {
		return err
	}
	err = swap.SolveVec(dxVec, false, dxVec)
	if err != nil {
		return err
	}
	return nil
}

func computeDual(dzVec *mat.VecDense, an *mat.Dense, lvIdx int, lu *mat.LU, swap *Swap) error {
	tmp := indicator(swap.Dim, lvIdx)
	err := swap.SolveVec(tmp, true, tmp)
	if err != nil {
		return err
	}
	err = lu.SolveVec(tmp, true, tmp)
	if err != nil {
		return err
	}
	dzVec.MulVec(an.T(), tmp)
	dzVec.ScaleVec(-1, dzVec)
	return nil
}

func verifyPerturbations(c, b, cbar, bbar []float64) {
	if cbar != nil && len(cbar) != 0 && len(cbar) != len(c) {
		panic("lp: cbar vector incorrect length")
	}
	if bbar != nil && len(bbar) != 0 && len(bbar) != len(b) {
		panic("lp: bbar vector incorrect length")
	}
}

func indicator(n, i int) *mat.VecDense {
	e := make([]float64, n)
	e[i] = 1
	return mat.NewVecDense(n, e)
}

func selectIdx(x, dx []float64) (int, error) {
	n := len(x)
	m := math.Inf(1)
	xnorm := floats.Norm(x, math.Inf(1))
	var idx int

	for i := 0; i < n; i++ {
		if v := dx[i]; v > absZeroTol && v > relZeroTol*xnorm {
			if w := x[i] / v; w < m {
				m = w
				idx = i
			}
		}
	}
	if math.IsInf(m, 1) {
		return -1, ErrInfeasible
	}
	return idx, nil
}

func lambdaFn(xb, xbbar, zn, znbar []float64) (bool, int, float64, float64) {
	nx := len(xb)
	nz := len(zn)

	var boolean bool
	var index int
	var lambdaMin float64 = math.Inf(-1)
	var lambdaMax float64 = math.Inf(1)

	for i := 0; i < nx; i++ {
		if v := xbbar[i]; v != 0 {
			d := -xb[i] / v
			if v > 0 && d > lambdaMin {
				boolean = false
				lambdaMin = d
				index = i
			} else if v < 0 && d < lambdaMax {
				lambdaMax = d
			}
		}
	}
	for i := 0; i < nz; i++ {
		if v := znbar[i]; v != 0 {
			d := -zn[i] / v
			if v > 0 && d > lambdaMin {
				boolean = true
				lambdaMin = d
				index = i
			} else if v < 0 && d < lambdaMax {
				lambdaMax = d
			}
		}
	}
	return boolean, index, lambdaMin, lambdaMax
}
