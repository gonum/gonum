// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat_test

import (
	"fmt"
	"log"

	"gonum.org/v1/gonum/mat"
)

func ExampleSVD_SolveTo() {
	// The system described by A is rank deficient.
	a := mat.NewDense(5, 3, []float64{
		-1.7854591879711257, -0.42687285925779594, -0.12730256811265162,
		-0.5728984211439724, -0.10093393134001777, -0.1181901192353067,
		1.2484316018707418, 0.5646683943038734, -0.48229492403243485,
		0.10174927665169475, -0.5805410929482445, 1.3054473231942054,
		-1.134174808195733, -0.4732430202414438, 0.3528489486370508,
	})

	// Perform an SVD retaining all singular vectors.
	var svd mat.SVD
	ok := svd.Factorize(a, mat.SVDFull)
	if !ok {
		log.Fatal("failed to factorize A")
	}

	// Determine the rank of the A matrix with a near zero condition threshold.
	const rcond = 1e-15
	rank := svd.Rank(rcond)
	if rank == 0 {
		log.Fatal("zero rank system")
	}

	b := mat.NewDense(5, 2, []float64{
		-2.318, -4.35,
		-0.715, 1.451,
		1.836, -0.119,
		-0.357, 3.094,
		-1.636, 0.021,
	})

	// Find a least-squares solution using the determined parts of the system.
	var x mat.Dense
	svd.SolveTo(&x, b, rank)

	fmt.Printf("singular values = %v\nrank = %d\nx = %.15f",
		format(svd.Values(nil), 4, rcond), rank, mat.Formatted(&x, mat.Prefix("    ")))

	// Output:
	// singular values = [2.685 1.526 <1e-15]
	// rank = 2
	// x = ⎡ 1.212064313552347   1.507467451093930⎤
	//     ⎢ 0.415400738264774  -0.624498607705372⎥
	//     ⎣-0.183184442255280   2.221334193689124⎦
}

func ExampleSVD_SolveVecTo() {
	// The system described by A is rank deficient.
	a := mat.NewDense(5, 3, []float64{
		-1.7854591879711257, -0.42687285925779594, -0.12730256811265162,
		-0.5728984211439724, -0.10093393134001777, -0.1181901192353067,
		1.2484316018707418, 0.5646683943038734, -0.48229492403243485,
		0.10174927665169475, -0.5805410929482445, 1.3054473231942054,
		-1.134174808195733, -0.4732430202414438, 0.3528489486370508,
	})

	// Perform an SVD retaining all singular vectors.
	var svd mat.SVD
	ok := svd.Factorize(a, mat.SVDFull)
	if !ok {
		log.Fatal("failed to factorize A")
	}

	// Determine the rank of the A matrix with a near zero condition threshold.
	const rcond = 1e-15
	rank := svd.Rank(rcond)
	if rank == 0 {
		log.Fatal("zero rank system")
	}

	b := mat.NewVecDense(5, []float64{-2.318, -0.715, 1.836, -0.357, -1.636})

	// Find a least-squares solution using the determined parts of the system.
	var x mat.VecDense
	svd.SolveVecTo(&x, b, rank)

	fmt.Printf("singular values = %v\nrank = %d\nx = %.15f",
		format(svd.Values(nil), 4, rcond), rank, mat.Formatted(&x, mat.Prefix("    ")))

	// Output:
	// singular values = [2.685 1.526 <1e-15]
	// rank = 2
	// x = ⎡ 1.212064313552347⎤
	//     ⎢ 0.415400738264774⎥
	//     ⎣-0.183184442255280⎦
}

func format(vals []float64, prec int, eps float64) []string {
	s := make([]string, len(vals))
	for i, v := range vals {
		if v < eps {
			s[i] = fmt.Sprintf("<%.*g", prec, eps)
			continue
		}
		s[i] = fmt.Sprintf("%.*g", prec, v)
	}
	return s
}
