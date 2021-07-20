package lp

import (
	"golang.org/x/exp/rand"
	"math"
	"testing"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/mat"
)

const lambdaTol = 1e-9

func TestParametric(t *testing.T) {
	// First test specific inputs. Combination of problems from literature and
	// those collected from failures during randomized testing.
	rnd := rand.New(rand.NewSource(1))

	for _, test := range []struct {
		A    mat.Matrix
		b    []float64
		c    []float64
		tol  float64
		want struct {
			optF float64
			optX []float64
			err  error
		}
	}{
		{
			// Feasible LP
			A: mat.NewDense(4, 7, []float64{
				-1, 3, 0, 1, 0, 0, 0,
				3, 3, 0, 0, 1, 0, 0,
				0, 3, 2, 0, 0, 1, 0,
				-3, 0, -5, 0, 0, 0, 1,
			}),
			b: []float64{5, 4, 6, -4},
			c: []float64{3, -11, -2, 0, 0, 0, 0},
			want: struct {
				optF float64
				optX []float64
				err  error
			}{-50.0 / 3, []float64{0, 4.0 / 3, 1, 1, 0, 0, 1}, nil},
		},
		{
			// Unbounded LP
			A: mat.NewDense(3, 6, []float64{
				1, -1, 0, 1, 0, 0,
				-2, 1, 0, 0, 1, 0,
				0, 1, -2, 0, 0, 1,
			}),
			b: []float64{5, 3, 5},
			c: []float64{0, -2, -1, 0, 0, 0},
			want: struct {
				optF float64
				optX []float64
				err  error
			}{math.NaN(), nil, ErrInfeasible},
		},
		{
			// Infeasible LP
			A: mat.NewDense(3, 6, []float64{
				1, 3, 2, 4, 1, 0,
				3, 1, 2, 1, 0, 1,
				5, 3, 3, 3, 0, 0,
			}),
			b: []float64{5, 4, 9},
			c: []float64{-1, -1, -1, -1, 0, 0},
			want: struct {
				optF float64
				optX []float64
				err  error
			}{math.NaN(), nil, ErrInfeasible},
		},
		{
			A: mat.NewDense(4, 5, []float64{
				0, -0.9742386073191733, 0, 0.2816837892460988, 0,
				0, -0.20539646935482825, 0, 0, 0.4487437183866715,
				-1.2144408758278467, 0, 0, 0, 0.510453094527252,
				0, 0, -0.006610967893518338, 0, 0,
			}),
			b: []float64{0, 0, 0, 0},
			c: []float64{-0.8400829172082259, 0, 0.03666949844210843, 0, 0},
			want: struct {
				optF float64
				optX []float64
				err  error
			}{math.NaN(), nil, ErrInfeasible},
		},
		{
			A: mat.NewDense(6, 9, []float64{
				0, 1.6392367266859786, -0.14741113226313596, -0.769911159365356, 0, -0.6382401523215369, 0, 0, 0,
				1.2197160204747788, 0, -0.4861293521246829, -0.31192544198903405, -1.948134108238715, 0, 0, 0, 0,
				0, -1.4605777644731557, 1.1790388360594117, 0, 0.45205146354525966, 0, 0.02909201412376361, 0, 0,
				0, -2.1401618581313238, 0.678116719369447, 0, 1.176117554101876, -0.5464795797539499, 0, 0, 0,
				0, 0, -0.6885246482305794, -1.2044205819013112, 0, 0, 0.9330480368289458, 0, -1.1408920167186356,
				0, 0, 0.7116088199099427, -0.6204671724570668, 0, 0, -0.8651056202786371, -0.9483682123278148, 0,
			}),
			b: []float64{0, -0.3662758531043814, 0, 0, 0, 0},
			c: []float64{0, 0, 0, 0, -0.8893328727966363, 0, 0, 0, 0},
			want: struct {
				optF float64
				optX []float64
				err  error
			}{math.NaN(), nil, ErrInfeasible},
		},
		{
			A: mat.NewDense(8, 11, []float64{
				0, -3.0010362933716372, 0.7695299043896839, 0, 0, 0, 0.7068940135484373, 0.8339224715131716, 0, 0, 0.8762801726937165,
				0, -0.30072592533121284, 1.636167870689285, 0, -1.001900093669902, 0, -0.5969178654006676, 0, 0, 0, 0,
				0, 0, 0, 0.3159774559090017, 0, 0, 0, 0, 0, -2.135325621867796, -0.5472986667952038,
				-1.1118762076948259, 0, 0.5830526440332912, 0, 0, 0, 0, 0, 0, -0.5464623021033211, 0,
				0, 0, 0, 0, 0, 1.6081491415334566, 0, -0.9701015378936042, 0, 0, 0,
				0.4582342400804322, 0, 0.4819613032018917, 0, 0.8855994897783475, -0.7418793528921341, 0, 0.04372539903952011, 0, 0.3360257804029679, 0,
				0, 0, 0, -1.2283105841094097, 0, 0, 0, 0, 0.036260028615893725, 0, 0,
				0, -0.6335131547579882, 0, 0, 0, 0, 0, 0, 1.0954699437501858, 0, 0,
			}),
			b: []float64{0, 0, 0, 0, -0.24203013208240787, -1.308451156438598, -0.08302571655946211, 0},
			c: []float64{0, 0, 0, 0, 0, 0, 0, 0, -1.120375655128679, 0, 0},
			want: struct {
				optF float64
				optX []float64
				err  error
			}{math.NaN(), nil, ErrInfeasible},
		},
		{
			A: mat.NewDense(4, 5, []float64{
				0, 0.268052398505888, 0, 0, 0,
				0, 0, -0.4879133011183309, 0, 0,
				0.8803778393701118, 0.8179955009659348, 0, 0, 0,
				0, 0.6726012375359441, 0, 0.9311880555752978, -1.3901276377747083,
			}),
			b: []float64{-1.2150912536093552, 0, 1.1064714448349866, -0.5312348881892228},
			c: []float64{0, 0, 0, -1.130037953158608, -0.2885385379752342},
			want: struct {
				optF float64
				optX []float64
				err  error
			}{math.NaN(), nil, ErrInfeasible},
		},
		{
			A: mat.NewDense(3, 7, []float64{
				1.4966464507472446, 0.2632657569804575, -0.44911885177422195, 1.1382099517162505, -0.47969637327748016, -0.09294552620408325, 0.3227440593322585,
				1.0966895566442567, -0.003330485058491206, -2.884796181323111, 0.650305178541605, 0.4511143535432509, -0.6255349662329515, -0.10077057548899038,
				3.941329544954186e-06, 1.0436981345081278, 0.3087885408567783, 1.3350887904950457, -0.2438311495245138, 0.16889024616635862, 0.17923883902552618,
			}),
			b: []float64{-0.7429992255960549, -0.3000189383827957, 1.5743763396716914},
			c: []float64{1.2243481103297997, 0.19985352640758625, 1.5910855224147022, 0.6211775827803893, 0.04641879687782291, 0.33622517819534153, -0.3487463126329189},
			want: struct {
				optF float64
				optX []float64
				err  error
			}{-7.7432431619022815, []float64{0, 0, 2.5826160644529463, 0, 24.160446274645135, 0, 37.201554740925083}, nil},
		},
		{
			// Poorly conditioned dual of previous lp, generated random with seed
			// of 1. Must increase lambdaTolerance to solve correctly.
			A: mat.NewDense(7, 13, []float64{
				-1.4966464507472446, -1.0966895566442567, -3.941329544954186e-06, 1.4966464507472446, 1.0966895566442567, 3.941329544954186e-06, 1, 0, 0, 0, 0, 0, 0,
				-0.2632657569804575, 0.003330485058491206, -1.0436981345081278, 0.2632657569804575, -0.003330485058491206, 1.0436981345081278, 0, 1, 0, 0, 0, 0, 0,
				0.44911885177422195, 2.884796181323111, -0.3087885408567783, -0.44911885177422195, -2.884796181323111, 0.3087885408567783, 0, 0, 1, 0, 0, 0, 0,
				-1.1382099517162505, -0.650305178541605, -1.3350887904950457, 1.1382099517162505, 0.650305178541605, 1.3350887904950457, 0, 0, 0, 1, 0, 0, 0,
				0.47969637327748016, -0.4511143535432509, 0.2438311495245138, -0.47969637327748016, 0.4511143535432509, -0.2438311495245138, 0, 0, 0, 0, 1, 0, 0,
				0.09294552620408325, 0.6255349662329515, -0.16889024616635862, -0.09294552620408325, -0.6255349662329515, 0.16889024616635862, 0, 0, 0, 0, 0, 1, 0,
				-0.3227440593322585, 0.10077057548899038, -0.17923883902552618, 0.3227440593322585, -0.10077057548899038, 0.17923883902552618, 0, 0, 0, 0, 0, 0, 1,
			}),
			b:   []float64{1.2243481103297997, 0.19985352640758625, 1.5910855224147022, 0.6211775827803893, 0.04641879687782291, 0.33622517819534153, -0.3487463126329189},
			c:   []float64{-0.7429992255960549, -0.3000189383827957, 1.5743763396716914, 0.7429992255960549, 0.3000189383827957, -1.5743763396716914, 0, 0, 0, 0, 0, 0, 0},
			tol: 1e-9,
			want: struct {
				optF float64
				optX []float64
				err  error
			}{7.7432431619022815, []float64{0, 1.2194607641912651, 4.6274926407688461, 1.1086008326353591, 0, 0, 0.90255273217009591, 4.7336409297806537, 0, 6.3304922853645875, 0, 0.45798768921659117, 0}, nil},
		},
	} {
		// solve and check
		opt, x, _, err := Parametric(test.c, test.A, test.b, test.tol, nil, true, rnd)
		if err != nil {
			if test.want.err == nil {
				t.Errorf("Unexpected error: %s", err)
			} else if err != test.want.err {
				t.Errorf("Error mismatch. Want %v, got %v", test.want.err, err)
			}
			continue
		}
		if err == nil && test.want.err != nil {
			t.Errorf("Did not error during optimization.")
			continue
		}
		if !floats.EqualWithinAbsOrRel(opt, test.want.optF, 1e-10, 1e-10) {
			t.Errorf("Optimum value mismatch. Want %v, got %v", test.want.optF, opt)
		}
		if !floats.EqualApprox(x, test.want.optX, 1e-10) {
			t.Errorf("Optimum solution mismatch. Want %v, got %v", test.want.optX, x)
		}
	}

	// Randomized tests
	testRandomParametric(t, 20000, 0.7, 10, rnd)
	testRandomParametric(t, 20000, 0, 10, rnd)
	testRandomParametric(t, 200, 0, 100, rnd)
	testRandomParametric(t, 2, 0, 400, rnd)
}

func testRandomParametric(t *testing.T, nTest int, pZero float64, maxN int, rnd *rand.Rand) {
	// Try a bunch of random LPs
	for i := 0; i < nTest; i++ {
		n := rnd.Intn(maxN) + 2 // n must be at least two.
		m := rnd.Intn(n-1) + 1  // m must be between 1 and n
		if m == 0 || n == 0 {
			continue
		}
		randValue := func() float64 {
			//var pZero float64
			v := rnd.Float64()
			if v < pZero {
				return 0
			}
			return rnd.NormFloat64()
		}
		a := mat.NewDense(m, n, nil)
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				a.Set(i, j, randValue())
			}
		}
		b := make([]float64, m)
		for i := range b {
			b[i] = randValue()
		}

		c := make([]float64, n)
		for i := range c {
			c[i] = randValue()
		}

		testParametric(t, nil, c, a, b, lambdaTol, rnd)
	}
}

func testParametric(t *testing.T, initialBasic []int, c []float64, a mat.Matrix, b []float64, lambdaTol float64, rnd *rand.Rand) error {
	primalOpt, primalX, _, _, errPrimal := parametric(initialBasic, c, a, b, nil, nil, lambdaTol, false, true, rnd)
	if errPrimal == nil {
		// No error solving the simplex, check that the solution is feasible.
		var bCheck mat.VecDense
		bCheck.MulVec(a, mat.NewVecDense(len(primalX), primalX))
		if !mat.EqualApprox(&bCheck, mat.NewVecDense(len(b), b), 1e-10) {
			t.Errorf("No error in primal but solution infeasible")
		}
	}

	primalInfeasible := errPrimal == ErrInfeasible
	primalUnbounded := errPrimal == ErrUnbounded
	primalBounded := errPrimal == nil
	primalASingular := errPrimal == ErrSingular
	primalZeroRow := errPrimal == ErrZeroRow
	primalZeroCol := errPrimal == ErrZeroColumn
	primalBad := !primalInfeasible && !primalUnbounded && !primalBounded &&
		!primalASingular && !primalZeroRow && !primalZeroCol

	// It's an error if it's not one of the known returned errors. If it's
	// singular the problem is undefined and so the result cannot be compared
	// to the dual.
	if errPrimal == ErrSingular || primalBad {
		if primalBad {
			t.Errorf("non-known error returned: %s", errPrimal)
		}
		return errPrimal
	}

	// Compare the result to the answer found from solving the dual LP.

	// Construct and solve the dual LP.
	// Standard Form:
	//  minimize c^T * x
	//    subject to  A * x = b, x >= 0
	// The dual of this problem is
	//  maximize -b^T * nu
	//   subject to A^T * nu + c >= 0
	// Which is
	//   minimize b^T * nu
	//   subject to -A^T * nu <= c

	negAT := &mat.Dense{}
	negAT.Clone(a.T())
	negAT.Scale(-1, negAT)
	cNew, aNew, bNew := Convert(b, negAT, c, nil, nil)

	dualOpt, dualX, _, _, errDual := parametric(nil, cNew, aNew, bNew, nil, nil, lambdaTol, false, true, rnd)
	if errDual == nil {
		// Check that the dual is feasible
		var bCheck mat.VecDense
		bCheck.MulVec(aNew, mat.NewVecDense(len(dualX), dualX))
		if !mat.EqualApprox(&bCheck, mat.NewVecDense(len(bNew), bNew), 1e-10) {
			t.Errorf("No error in dual but solution infeasible")
		}
	}

	// Check about the zero status.
	if errPrimal == ErrZeroRow || errPrimal == ErrZeroColumn {
		return errPrimal
	}

	// If the primal problem is feasible, then the primal and the dual should
	// be the same answer. We have flopped the sign in the dual (minimizing
	// b^T *nu instead of maximizing -b^T*nu), so flip it back.
	if errPrimal == nil {
		if errDual != nil {
			t.Errorf("Primal feasible but dual errored: %s", errDual)
		}
		dualOpt *= -1
		if !floats.EqualWithinAbsOrRel(dualOpt, primalOpt, lambdaTol, lambdaTol) {
			t.Errorf("Primal and dual value mismatch. Primal %v, dual %v.", primalOpt, dualOpt)
		}
	}
	// If the primal problem is unbounded, then the dual should be infeasible.
	if errPrimal == ErrUnbounded && errDual != ErrInfeasible {
		t.Errorf("Primal unbounded but dual not infeasible. ErrDual = %s", errDual)
	}

	// If the dual is unbounded, then the primal should be infeasible.
	if errDual == ErrUnbounded && errPrimal != ErrInfeasible {
		t.Errorf("Dual unbounded but primal not infeasible. ErrDual = %s", errPrimal)
	}

	// If the primal is infeasible, then the dual should be either infeasible
	// or unbounded.
	if errPrimal == ErrInfeasible {
		if errDual != ErrUnbounded && errDual != ErrInfeasible && errDual != ErrZeroColumn {
			t.Errorf("Primal infeasible but dual not infeasible or unbounded: %s", errDual)
		}
	}

	return errPrimal
}
