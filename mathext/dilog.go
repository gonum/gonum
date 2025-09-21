package mathext

import (
	"math"
	"math/cmplx"
)

// Li2 returns the dilogarithm Li2(z) on the principal branch.
//
// For |z| < 1, Li2(z) is defined by the power series
//
//     Li2(z) = SUM_{k=1}^{infinity} z^k / k^2
//
// and is analytically continued to the rest of the complex plane.
// The implementation uses reflection and inversion identities to map z into
// a region where the series converges rapidly, then evaluates the series.
//
// Branch cut: Li2 has a logarithmic branch point at z=1 with the standard
// cut on the real axis for z in the interval (1, infinity). The principal value is taken with
// Arg(z) element of (−Pi, Pi].

func Li2(z complex128) complex128 {
	// Special cases
	if z == 0 {
		return 0
	}
	if z == 1 {
		return complex(math.Pi*math.Pi/6, 0)
	}

	// Reflection: map Re(z) > 0.5 into left half-plane for better convergence
	// This formula is applied before inversion on the principal branch and gives very accurate results
	// for real z > 1 because Li2(-1) and similar values are known exactly.
	// Inversion first would also be valid, but is more sensitive to the
	// Arg(-z)= +/-Pi branch choice in log(-z) and can yield wrong imaginary signs
	// if care is not taken.
	if real(z) > 0.5 {
		return complex(math.Pi*math.Pi/6, 0) - cmplx.Log(z)*cmplx.Log(1-z) - Li2(1-z)
	}

	// Inversion: map |z| > 1 into unit disk
	if cmplx.Abs(z) > 1 {
		logmz := cmplx.Log(-z)
		return -complex(math.Pi*math.Pi/6, 0) - 0.5*logmz*logmz - Li2(1/z)
	}

	// Direct series for |z| <= 1 and Re(z) <= 0.5
	return li2Series(z)
}

func li2Series(z complex128) complex128 {
	const tol = 1e-15
	sum := complex(0, 0)
	zk := z // zk = z^k
	for k := 1; cmplx.Abs(zk)/float64(k*k) > tol; k++ {
		sum += zk / complex(float64(k*k), 0)
		zk *= z
	}
	return sum
}
