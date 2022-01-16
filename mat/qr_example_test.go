// Copyright ©2022 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat_test

import (
	"fmt"
	"log"

	"gonum.org/v1/gonum/mat"
)

func ExampleQR_solveTo() {
	// QR factorization can be used for solving linear inverse problems,
	// as this is a more numerically stable technique than direct
	// matrix inversion.
	//
	// Here, we want to solve:
	//   Ax = b

	var (
		a = mat.NewDense(4, 2, []float64{0, 1, 1, 1, 1, 1, 2, 1})
		b = mat.NewDense(4, 1, []float64{1, 0, 2, 1})
		x = mat.NewDense(2, 1, nil)
	)

	var qr mat.QR
	qr.Factorize(a)

	err := qr.SolveTo(x, false, b)
	if err != nil {
		log.Fatalf("could not solve QR: %+v", err)
	}
	fmt.Printf("%.3f\n", mat.Formatted(x))

	// Output:
	// ⎡0.000⎤
	// ⎣1.000⎦
}
