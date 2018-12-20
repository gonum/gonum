// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dualcmplx_test

import (
	"fmt"

	"gonum.org/v1/gonum/num/dualcmplx"
)

// Example point, displacement and rotation from Euclidean Space Dual Complex Number page:
// http://www.euclideanspace.com/maths/algebra/realNormedAlgebra/other/dualComplex/index.htm

func Example_displace() {
	// Displace a point [3, 4] by [4, 3].

	// Point to be transformed in the dual imaginary vector.
	p := dualcmplx.Number{Real: 1, Dual: 3 + 4i}

	// Displacement vector, [4, 3], in the dual imaginary vector.
	d := dualcmplx.Number{Real: 1, Dual: 4 + 3i}

	fmt.Println(dualcmplx.Mul(d, p).Dual)
	// Output:
	//
	// (7+7i)
}

func Example_rotate() {
	// Rotate a point [3, 4] by 90° around the origin.

	// Point to be transformed in the dual imaginary vector.
	p := dualcmplx.Number{Real: 1, Dual: 3 + 4i}

	// Rotation in the real quaternion.
	r := dualcmplx.Number{Real: 0 + 1i}

	fmt.Println(dualcmplx.Mul(r, p).Dual)
	// Output:
	//
	// (-4+3i)
}

func Example_displaceAndRotate() {
	// Displace a point [3, 4] by [4, 3] and then rotate
	// by 90° around the origin.

	// Point to be transformed in the dual imaginary vector.
	p := dualcmplx.Number{Real: 1, Dual: 3 + 4i}

	// Displacement vector, [4, 3], in the dual imaginary vector.
	d := dualcmplx.Number{Real: 1, Dual: 4 + 3i}

	// Rotation in the real quaternion.
	r := dualcmplx.Number{Real: 0 + 1i}

	// Combine the rotation and displacement so
	// the displacement is performed first.
	q := dualcmplx.Mul(r, d)

	fmt.Println(dualcmplx.Mul(q, p).Dual)
	// Output:
	//
	// (-7+7i)
}

func Example_rotateAndDisplace() {
	// Rotate a point [3, 4] by 90° around the origin and then
	// displace by [4, 3].

	// Point to be transformed in the dual imaginary vector.
	p := dualcmplx.Number{Real: 1, Dual: 3 + 4i}

	// Displacement vector, [4, 3], in the dual imaginary vector.
	d := dualcmplx.Number{Real: 1, Dual: 4 + 3i}

	// Rotation in the real quaternion.
	r := dualcmplx.Number{Real: 0 + 1i}

	// Combine the rotation and displacement so
	// the displacement is performed first.
	q := dualcmplx.Mul(d, r)

	fmt.Println(dualcmplx.Mul(q, p).Dual)
	// Output:
	//
	// (-7+7i)
}
