// Copyright ©2025 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mathext

import (
	"math"
	"math/cmplx"
	"testing"
)

func TestDiLogValues(t *testing.T) {
	t.Parallel()

	for _, test := range []struct {
		input, want complex128
	}{
		// Reference values were generated using Python's scipy.special.spence
		// (where Li2(z) = spence(1 - z)) and verified with a computer algebra system.

		// well known values
		{0 + 0i, 0 + 0i},
		{1 + 0i, math.Pi*math.Pi/6 + 0i},
		{-1 + 0i, -math.Pi*math.Pi/12 + 0i},
		{0.5 + 0i, 0.582240526465012 + 0i},
		{2 + 0i, 2.467401100272340 - 2.177586090303602i},
		// abs(z) < 0.5
		{0.1 + 0.1i, 0.099751196798713 + 0.105220383905600i},
		{-0.3 + 0.39i, -0.306174046754339 + 0.337550668621259i},
		{0.001 - 0.49i, -0.055831393285929 - 0.478156430372353i},
		// 0.5 < abs(z) < 1
		{0.5 + 0.7i, 0.359364350522653 + 0.856767712327590i},
		{complex(math.Pi/4, -math.Pi/7), 0.828086461155377 - 0.739345385309341i},
		{-0.8 - 0.0001i, -0.679781588954542 - 0.000073473333084i},
		// abs(z) > 1
		{5 + 0i, 1.783719161266631 - 5.056198322111862i},
		{-10 + 0i, -4.198277886858104 + 0i},
		{1000 + 10000i, -42.71073756884990 + 15.39396088869304i},
		{-1791.91931 + 0.5i, -29.70223568904652 + 0.00209038439188i},
	} {
		got := Li2(test.input)
		const tol = 1e-12
		diff := cmplx.Abs(got - test.want)
		if cmplx.Abs(test.want) > 0 {
			if diff/cmplx.Abs(test.want) > tol {
				t.Errorf("Li2(%g) relative error %g exceeds tol %g", test.input, diff/cmplx.Abs(test.want), tol)
			}
		} else if diff > tol {
			t.Errorf("Li2(%g) abs error %g exceeds tol %g", test.input, diff, tol)
		}
	}
}

func TestDiLogProperties(t *testing.T) {
	t.Parallel()
	const tol = 1e-12

	// Duplication formula: Li2(z^2) = 2 (Li2(z) + Li2(-z))
	for i, z := range []complex128{
		0.1 + 0.1i,
		-0.3 + 0.4i,
		0.5,
		200,
	} {
		lhs := Li2(z * z)
		rhs := 2 * (Li2(z) + Li2(-z))
		if math.Abs(real(lhs)-real(rhs)) > tol || math.Abs(imag(lhs)-imag(rhs)) > tol {
			t.Errorf("duplication formula failed for case %d, z=%v: got %v want %v", i, z, lhs, rhs)
		}
	}

	// Conjugation symmetry: Li2(conj(z)) == conj(Li2(z))
	// Valid for all z not on the branch cut (real z > 1).
	for i, z := range []complex128{
		0.3 + 0.5i,
		-0.8 + 0.1i,
		2 + 0.7i,
	} {
		lhs := Li2(cmplx.Conj(z))
		rhs := cmplx.Conj(Li2(z))
		if math.Abs(real(lhs)-real(rhs)) > tol || math.Abs(imag(lhs)-imag(rhs)) > tol {
			t.Errorf("conjugation symmetry failed for case %d, z=%v: got %v want %v", i, z, lhs, rhs)
		}
	}
}
