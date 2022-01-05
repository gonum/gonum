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
	var (
		x = []float64{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
		y = []float64{1, 6, 17, 34, 57, 86, 121, 162, 209, 262, 321}

		degree = 2

		a = Vandermonde(x, degree+1)
		b = mat.NewDense(len(y), 1, y)
		c = mat.NewDense(degree+1, 1, nil)
	)

	var qr mat.QR
	qr.Factorize(a)

	const trans = false
	err := qr.SolveTo(c, trans, b)
	if err != nil {
		log.Fatalf("could not solve QR: %+v", err)
	}
	fmt.Printf("%.3f\n", mat.Formatted(c))

	// Output:
	// ⎡1.000⎤
	// ⎢2.000⎥
	// ⎣3.000⎦
}

func Vandermonde(a []float64, d int) *mat.Dense {
	x := mat.NewDense(len(a), d, nil)
	for i := range a {
		for j, p := 0, 1.0; j < d; j, p = j+1, p*a[i] {
			x.Set(i, j, p)
		}
	}
	return x
}
