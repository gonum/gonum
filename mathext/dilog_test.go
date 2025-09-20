package mathext

import (
	"math"
	"math/cmplx"
	"testing"
)

func TestDiLogValues(t *testing.T) {
	t.Parallel()

	for i, test := range []struct {
		input, want complex128
	}{
		// well known values
		{0 + 0i, 0 + 0i},
		{1 + 0i, math.Pi*math.Pi/6 + 0i},
		{-1 + 0i, -math.Pi*math.Pi/12 + 0i},
		{0.5 + 0i, 0.5822405264650125 + 0i},
		{2 + 0i, 2.467401100272340 - 2.177586090303602i},
		// abs(z) < 0.5
		{0.1 + 0.1i, 0.0997511967987130 + 0.1052203839055999i},
		{-0.3 + 0.39i, -0.3061740467543387 + 0.3375506686212592i},
		{0.001 - 0.49i, -0.0558313932859285 - 0.4781564303723531i},
		// 0.5 < abs(z) < 1
		{0.5 + 0.7i, 0.3593643505226532 + 0.8567677123275905i},
		{complex(math.Pi/4, -math.Pi/7), 0.8280864611553766 - 0.7393453853093411i},
		{-0.8 - 0.0001i, -0.6797815889545422 - 0.0000734733330837i},
		// abs(z) > 1
		{5 + 0i, 1.783719161266631 - 5.056198322111863i},
		{-10 + 0i, -4.198277886858104 + 0i},
		{1000 + 10000i, -42.71073756884990 + 15.39396088869304i},
		{-1791.91931 + 0.5i, -29.70223568904652 + 0.00209038439188i},
	} {

		if got := Li2(test.input); cmplx.Abs(got-test.want) > 1e-12 || math.Abs(real(got)-real(test.want)) > 1e-12 || math.Abs(imag(got)-imag(test.want)) > 1e-12 {
			t.Errorf("test %d Li2(%g) failed: got %g want %g", i, test.input, got, test.want)
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

	// Conjugation symmetry: Li2(conj(z)) == conj(Li2(z)), only valid off the branch cut (1, infinity)
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
