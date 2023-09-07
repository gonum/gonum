// Copyright ©2023 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math"

	"gonum.org/v1/gonum/blas"
	"gonum.org/v1/gonum/blas/blas64"
	"gonum.org/v1/gonum/lapack"
)

// Dhgeqz computes the eigenvalues of a real matrix pair (H,T),
// where H is an upper Hessenberg matrix and T is upper triangular,
// using the double-shift QZ method.
// Matrix pairs of this type are produced by the reduction to
// generalized upper Hessenberg form of a real matrix pair (A,B):
//
//	A = Q1*H*Z1ᵀ,  B = Q1*T*Z1ᵀ,
//
// as computed by DGGHRD.
// If JOB='S', then the Hessenberg-triangular pair (H,T) is
// also reduced to generalized Schur form,
//
//	H = Q*S*Zᵀ,  T = Q*P*Zᵀ,
//
// where Q and Z are orthogonal matrices, P is an upper triangular
// matrix, and S is a quasi-triangular matrix with 1-by-1 and 2-by-2
// diagonal blocks.
// The 1-by-1 blocks correspond to real eigenvalues of the matrix pair
// (H,T) and the 2-by-2 blocks correspond to complex conjugate pairs of
// eigenvalues.
// Additionally, the 2-by-2 upper triangular diagonal blocks of P
// corresponding to 2-by-2 blocks of S are reduced to positive diagonal
// form, i.e., if S(j+1,j) is non-zero, then P(j+1,j) = P(j,j+1) = 0,
// P(j,j) > 0, and P(j+1,j+1) > 0.
//
// Optionally, the orthogonal matrix Q from the generalized Schur
// factorization may be postmultiplied into an input matrix Q1, and the
// orthogonal matrix Z may be postmultiplied into an input matrix Z1.
// If Q1 and Z1 are the orthogonal matrices from DGGHRD that reduced
// the matrix pair (A,B) to generalized upper Hessenberg form, then the
// output matrices Q1*Q and Z1*Z are the orthogonal factors from the
// generalized Schur factorization of (A,B):
//
//	A = (Q1*Q)*S*(Z1*Z)ᵀ,  B = (Q1*Q)*P*(Z1*Z)ᵀ.
//
// To avoid overflow, eigenvalues of the matrix pair (H,T) (equivalently,
// of (A,B)) are computed as a pair of values (alpha,beta), where alpha is
// complex and beta real.
// If beta is nonzero, lambda = alpha / beta is an eigenvalue of the
// generalized nonsymmetric eigenvalue problem (GNEP)
//
//	A*x = lambda*B*x
//
// and if alpha is nonzero, mu = beta / alpha is an eigenvalue of the
// alternate form of the GNEP
//
//	mu*A*y = B*y.
//
// Real eigenvalues can be read directly from the generalized Schur
// form:
//
//	alpha = S(i,i), beta = P(i,i).
//
// Argument info:
//   - info=-1: successful exit
//   - info>=0:
//   - info<=n: The QZ iteration did not converge. (H,T) is not in Schur form
//     but alphar(i), alphai(i), and beta(i), i=info+1,...,n should be correct.
//   - info=n+1...2*n: The shift calculation failed. (H,T) is not in Schur form
//     but alphar(i), alphai(i), and beta(i), i=info-n+1,...,n should be correct.
//   - alphar, alphai, beta are vectors of length n.
//   - work is a vector of length max(1, n)
//
// Ref: C.B. Moler & G.W. Stewart, "An Algorithm for Generalized Matrix Eigenvalue Problems", SIAM J. Numer. Anal., 10(1973), pp. 241--256.
// https://doi.org/10.1137/0710024
func (impl Implementation) Dhgeqz(job lapack.SchurJob, compq, compz lapack.SchurComp, n, ilo, ihi int, h []float64, ldh int, t []float64, ldt int, alphar, alphai, beta, q []float64, ldq int, z []float64, ldz int, work []float64, workspaceQuery bool) (info int) {
	_ = _columns(h, 1, 0, 0)
	var (
		jiter int // counts QZ iterations in main loop.
		// counts iterations run since ILAST was last changed.
		//This is therefore reset only when a 1-by-1 or  2-by-2 block deflates off the bottom.
		iiter                                                                int
		ilschr, ilq, ilz, ilazro, ilazr2, ilpivt                             bool
		icompq, icompz, ifirst, istart, j, maxiter, ilast, ilastm, ifrstm    int
		c, s, s1, s2, wr, wr2, wi, scale, temp, tempr, temp2, tempi, t2, t3  float64 // Trigonometric temporary variables.
		b22, b11, sr, cr, sl, cl, cz, t1, szr, szi, wabs, an, bn             float64
		a11, a21, a12, a22, c11r, c11i, c12, c21, c22r, c22i                 float64
		cq, sqr, sqi, a1r, a1i, a2r, a2i, ad11, ad21, ad12, ad22             float64
		u12, ad11l, ad21l, ad12l, ad22l, ad32l, u12l, vs, eshift             float64
		b1r, b1i, b1a, b2r, b2i, b2a, s1inv, tau, u1, u2, w11, w22, w12, w21 float64
		v                                                                    [3]float64
	)

	switch job {
	case lapack.EigenvaluesOnly:
		ilschr = false
	case lapack.EigenvaluesAndSchur:
		ilschr = true
	default:
		panic(badSchurJob)
	}

	switch compq {
	case lapack.SchurNone:
		ilq = false
		icompq = 1
	case lapack.SchurOrig:
		ilq = true
		icompq = 2
	case lapack.SchurHess:
		ilq = true
		icompq = 3
	default:
		panic(badSchurComp)
	}
	switch compz {
	case lapack.SchurNone:
		ilz = false
		icompz = 1
	case lapack.SchurOrig:
		ilz = true
		icompz = 2
	case lapack.SchurHess:
		ilz = true
		icompz = 3
	default:
		panic(badSchurComp)
	}
	lwork := len(work)
	switch {
	case n < 0:
		panic(nLT0)
	case ilo < 0 || ilo >= n:
		panic(badIlo)
	case ihi < ilo-1 || ihi >= n:
		panic(badIhi)
	case ldh < n:
		panic(badLdH)
	case ldt < n:
		panic(badLdT)
	case ldq < 1 || (ilq && ldq < n):
		panic(badLdQ)
	case ldz < 1 || (ilz && ldz < n):
		panic(badLdZ)
	case lwork < max(1, n) && !workspaceQuery:
		panic(badLWork)
	case n == 0 || workspaceQuery:
		// Quick return or is workspace query.
		work[0] = float64(max(1, n))
		return -1
	}

	// Initialize Q and Z.
	if icompq == 3 {
		impl.Dlaset(blas.All, n, n, 0, 1, q, ldq)
	}
	if icompz == 3 {
		impl.Dlaset(blas.All, n, n, 0, 1, z, ldz)
	}

	// Machine constants.
	const (
		safmin = dlamchS
		safmax = 1. / safmin
		ulp    = dlamchE * dlamchB
	)

	in := ihi + 1 - ilo
	bi := blas64.Implementation()
	anorm := impl.Dlange('F', in, in, h[ilo*ldh+ilo:], ldh, work)
	bnorm := impl.Dlange('F', in, in, t[ilo*ldt+ilo:], ldt, work)

	atol := math.Max(safmin, ulp*anorm)
	btol := math.Max(safmin, ulp*bnorm)
	ascale := 1. / math.Max(safmin, anorm)
	bscale := 1. / math.Max(safmin, bnorm)

	// Set eigenvalues ihi+1 to n.
	for j = ihi + 1; j < n; j++ {
		if t[j*ldt+j] < 0 {
			if job == lapack.EigenvaluesAndSchur {
				for jr := 0; jr <= j; jr++ {
					h[jr*ldh+j] *= -1
					t[jr*ldt+j] *= -1
				}
			} else {
				h[j*ldh+j] *= -1
				t[j*ldt+j] *= -1
			}
			if ilz {
				for jr := 0; jr < n; jr++ {
					z[jr*ldz+j] *= -1
				}
			}
		}
		alphar[j] = h[j*ldh+j]
		alphai[j] = 0
		beta[j] = t[j*ldt+j]
	}

	// If ihi<ilo, skip QZ steps.
	if ihi < ilo {
		goto ThreeEighty
	}

	// MAIN QZ ITERATION LOOP
	// Initialize dynamic indices.
	// Eigenvalues ILAST+1:N have been found.
	// Column operations modify rows IFRSTM:whatever.
	// Row operations modify columns whatever:ILASTM.
	//
	// iiter counts iterations since last eigenvalue was found
	// to tell when to use an extraordinary shift.
	// Maxit is the maximum number of QZ sweeps allowed.
	ilast = ihi
	ifrstm = ilo
	ilastm = ihi
	if ilschr {
		ifrstm = 0
		ilastm = n - 1
	}
	maxiter = 30 * (ihi - ilo + 1)
	for jiter = 1; jiter <= maxiter; jiter++ {
		// Split the matrix if possible.
		// Two tests:
		//  1: H(j,j-1)=0  or  j=ILO
		//  2: T(j,j)=0
		if ilast == ilo {
			// special case: j=ilast
			goto Eighty
		}
		if math.Abs(h[ilast*ldh+ilast-1]) <= math.Max(safmin, ulp*(math.Abs(h[ilast*ldh+ilast])+math.Abs(h[(ilast-1)*ldh+ilast-1]))) {
			h[ilast*ldh+ilast-1] = 0
			goto Eighty
		}
		if math.Abs(t[ilast*ldt+ilast]) <= btol {
			t[ilast*ldt+ilast] = 0
			goto Seventy
		}

		// General case: j<ilast.

		for j = ilast - 1; j >= ilo; j-- {
			// Test 1: for H(j,j-1)=0 or j=ILO
			if j == ilo {
				ilazro = true
			} else {
				if math.Abs(h[j*ldh+j-1]) <= math.Max(safmin, ulp*(math.Abs(h[j*ldh+j])+math.Abs(h[(j-1)*ldh+j-1]))) {
					h[j*ldh+j-1] = 0
					ilazro = true
				} else {
					ilazro = false
				}
			}

			// Test 2: for T(j,j)=0

			if math.Abs(t[j*ldt+j]) < btol {
				t[j*ldt+j] = 0
				//Test 1a: Check for 2 consecutive small subdiagonals in A.
				ilazr2 = false
				if !ilazro {
					temp = math.Abs(h[j*ldh+j-1])
					temp2 := math.Abs(h[j*ldh+j])
					tempr = math.Max(temp, temp2)
					if tempr < 1 && tempr != 0 {
						temp /= tempr
						temp2 /= tempr
					}
					if temp*(ascale*math.Abs(h[(j+1)*ldh+j])) <= temp2*(ascale*atol) {
						ilazr2 = true
					}
				}
				// If both tests pass (1 & 2), i.e., the leading diagonal
				// element of B in the block is zero, split a 1x1 block off
				// at the top. (I.e., at the J-th row/column) The leading
				// diagonal element of the remainder can also be zero, so
				// this may have to be done repeatedly.
				if ilazro || ilazr2 {
					for jch := j; jch <= ilast-1; jch++ {
						temp := h[jch*ldh+jch]
						c, s, h[jch*ldh+jch] = impl.Dlartg(temp, h[(jch+1)*ldh+jch])
						h[(jch+1)*ldh+jch] = 0
						bi.Drot(ilastm-jch, h[jch*ldh+jch+1:], 1, h[(jch+1)*ldh+jch+1:], 1, c, s)
						bi.Drot(ilastm-jch, t[jch*ldt+jch+1:], 1, t[(jch+1)*ldt+jch+1:], 1, c, s)
						if ilq {
							bi.Drot(n, q[jch:], ldq, q[jch+1:], ldq, c, s)
						}
						if ilazr2 {
							h[jch*ldh+jch-1] *= c
						}
						ilazr2 = false
						if math.Abs(t[(jch+1)*ldt+jch+1]) >= btol {
							if jch+1 >= ilast {
								goto Eighty
							} else {
								ifirst = jch + 1
								goto OneTen
							}
						}
						t[(jch+1)*ldt+jch+1] = 0
					}
					goto Seventy
				} else {
					// Only test 2 passed -- chase the zero to T(ILAST,ILAST).
					// Then process as is in the case t(ilast, ilast)=0
					for jch := j; jch <= ilast-1; jch++ {
						temp = t[jch*ldt+jch+1]
						c, s, t[jch*ldt+jch+1] = impl.Dlartg(temp, t[(jch+1)*ldt+jch+1])
						t[(jch+1)*ldt+jch+1] = 0

						if jch < ilastm-1 {
							bi.Drot(ilastm-jch-1, t[jch*ldt+jch+2:], 1, t[(jch+1)*ldt+jch+2:], 1, c, s)
						}
						bi.Drot(ilastm-jch+2, h[jch*ldh+jch-1:], 1, h[(jch+1)*ldh+jch-1:], 1, c, s)

						if ilq {
							bi.Drot(n, q[jch:], ldq, q[jch+1:], ldq, c, s)
						}
						temp = h[(jch+1)*ldh+jch]
						c, s, h[(jch+1)*ldh+jch] = impl.Dlartg(temp, h[(jch+1)*ldh+jch-1])
						h[(jch+1)*ldh+jch-1] = 0
						bi.Drot(jch+1-ifrstm, h[ifrstm*ldh+jch:], ldh, h[ifrstm*ldh+jch-1:], ldh, c, s)
						bi.Drot(jch-ifrstm, t[ifrstm*ldt+jch:], ldt, t[ifrstm*ldt+jch-1:], ldt, c, s)
						if ilz {
							bi.Drot(n, z[jch:], ldz, z[jch-1:], ldz, c, s)
						}
					}
					goto Seventy
				}
			} else if ilazro {
				// Only test 1 passed -- work on j:ilast.
				ifirst = j
				goto OneTen
			}
			// Neither test passed -- try next j.
		}
		panic("unreachable")

	Seventy:
		// T(ILAST,ILAST)=0 -- clear H(ILAST,ILAST-1) to split off a
		// 1x1 block.
		temp = h[ilast*ldh+ilast-1]
		c, s, h[ilast*ldh+ilast] = impl.Dlartg(temp, h[ilast*ldh+ilast-1])
		h[ilast*ldh+ilast-1] = 0
		bi.Drot(ilast-ifrstm, h[ifrstm*ldh+ilast:], ldh, h[ifrstm*ldh+ilast-1:], ldh, c, s)
		bi.Drot(ilast-ifrstm, t[ifrstm*ldt+ilast:], ldt, t[ifrstm*ldt+ilast-1:], ldt, c, s)
		if ilz {
			bi.Drot(n, z[ilast:], ldz, z[ilast-1:], ldz, c, s)
		}

		// H(ILAST,ILAST-1)=0 -- Standardize B, set ALPHAR, ALPHAI, and BETA.

	Eighty:
		if t[ilast*ldt+ilast] < 0 {
			if ilschr {
				for j := ifrstm; j <= ilast; j++ {
					h[j*ldh+ilast] *= -1
					t[j*ldt+ilast] *= -1
				}
			} else {
				h[ilast*ldh+ilast] *= -1
				t[ilast*ldt+ilast] *= -1
			}
			if ilz {
				for j := 0; j < n; j++ {
					z[j*ldz+ilast] *= -1
				}
			}
		}
		alphar[ilast] = h[ilast*ldh+ilast]
		alphai[ilast] = 0
		beta[ilast] = t[ilast*ldt+ilast]

		// Go to next block -- exit if finished.
		ilast--
		if ilast < ilo {
			goto ThreeEighty
		}
		// Reset counters.

		iiter = 0
		eshift = 0.0
		if !ilschr {
			ilastm = ilast
			if ifrstm > ilast {
				ifrstm = ilo
			}
		}
		goto ThreeFifty

		// QZ Step
		// This iteration only involves rows/columns IFIRST:ILAST. We
		// assume IFIRST < ILAST, and that the diagonal of B is non-zero.

	OneTen:
		iiter++
		if !ilschr {
			ifrstm = ifirst
		}

		// Compute single shifts.
		// At this point, IFIRST < ILAST, and the diagonal elements of
		// T(IFIRST:ILAST,IFIRST,ILAST) are larger than BTOL (in magnitude).

		if (iiter/10)*10 == iiter {
			// Exceptional shift. Chosen for no particularly good reason. (single shift only)
			if float64(maxiter)*safmin*math.Abs(h[ilast*ldh+ilast-1]) < math.Abs(t[(ilast-1)*ldt+ilast-1]) {
				eshift = h[ilast*ldh+ilast-1] / t[(ilast-1)*ldt+ilast-1]
			} else {
				eshift += 1 / (safmin * float64(maxiter))
			}
			s1 = 1
			wr = eshift
		} else {
			// Shifts based on the generalized eigenvalues of the
			// bottom-right 2x2 block of A and B. The first eignevalue
			// returned by Dlag2 is the wilkinson shift (AEP p.512).
			s1, s2, wr, wr2, wi = impl.Dlag2(h[(ilast-1)*ldh+ilast-1:], ldh, t[(ilast-1)*ldt+ilast-1:], ldt)
			hlast := h[ilast*ldh+ilast]
			tlast := t[ilast*ldt+ilast]
			if math.Abs((wr/s1)*tlast-hlast) > math.Abs((wr2/s2)*tlast-hlast) {
				wr, wr2 = wr2, wr
				s1, s2 = s2, s1
			}
			temp = math.Max(s1, safmin*math.Max(1, math.Max(math.Abs(wr), math.Abs(wi))))
			if wi != 0 {
				goto TwoHundred
			}
		}

		// Fiddle with shift to avoid overflow.
		temp = math.Min(ascale, 1) * (safmax / 2)
		if s1 > temp {
			scale = temp / s1
		} else {
			scale = 1
		}

		temp = math.Min(bscale, 1) * (safmax / 2)
		if math.Abs(wr) > temp {
			scale = math.Min(scale, temp/math.Abs(wr))
		}
		s1 *= scale
		wr *= scale

		// Now check for two consecutive small subdiagonals.
		for j = ilast - 1; j >= ifirst+1; j-- {
			istart = j
			temp = math.Abs(s1 * h[j*ldh+j-1])
			temp2 := math.Abs(s1*h[j*ldh+j] - wr*t[j*ldt+j])
			tempr = math.Max(temp, temp2)
			if tempr < 1 && tempr != 0 {
				temp /= tempr
				temp2 /= tempr
			}
			if math.Abs(ascale*h[(j+1)*ldh+j]*temp) <= ascale*atol*temp2 {
				goto OneThirty
			}
		}
		istart = ifirst

	OneThirty:

		// Do an implicit-shift QZ sweep.
		// Initial Q.

		temp = s1*h[istart*ldh+istart] - wr*t[istart*ldt+istart]
		c, s, tempr = impl.Dlartg(temp, s1*h[(istart+1)*ldh+istart])

		// Sweep.
		for j = istart; j <= ilast-1; j++ {
			if j > istart {
				temp = h[j*ldh+j-1]
				c, s, h[j*ldh+j-1] = impl.Dlartg(temp, h[(j+1)*ldh+j-1])
				h[(j+1)*ldh+j-1] = 0
			}
			for jc := j; jc <= ilastm; jc++ {
				temp = c*h[j*ldh+jc] + s*h[(j+1)*ldh+jc]
				h[(j+1)*ldh+jc] = -s*h[j*ldh+jc] + c*h[(j+1)*ldh+jc]
				h[j*ldh+jc] = temp
				temp2 := c*t[j*ldt+jc] + s*t[(j+1)*ldt+jc]
				t[(j+1)*ldt+jc] = -s*t[j*ldt+jc] + c*t[(j+1)*ldt+jc]
				t[j*ldt+jc] = temp2
			}
			if ilq {
				for jr := 0; jr < n; jr++ {
					temp = c*q[jr*ldq+j] + s*q[jr*ldq+j+1]
					q[jr*ldq+j+1] = -s*q[jr*ldq+j] + c*q[jr*ldq+j+1]
					q[jr*ldq+j] = temp
				}
			}

			temp = t[(j+1)*ldt+j+1]
			c, s, t[(j+1)*ldt+j+1] = impl.Dlartg(temp, t[(j+1)*ldt+j])
			t[(j+1)*ldt+j] = 0

			maxjr := min(j+2, ilast)
			for jr := ifrstm; jr <= maxjr; jr++ {
				temp = c*h[jr*ldh+j+1] + s*h[jr*ldh+j]
				h[jr*ldh+j] = -s*h[jr*ldh+j+1] + c*h[jr*ldh+j]
				h[jr*ldh+j+1] = temp
			}
			for jr := ifrstm; jr <= j; jr++ {
				temp = c*t[jr*ldt+j+1] + s*t[jr*ldt+j]
				t[jr*ldt+j] = -s*t[jr*ldt+j+1] + c*t[jr*ldt+j]
				t[jr*ldt+j+1] = temp
			}
			if ilz {
				for jr := 0; jr < n; jr++ {
					temp = c*z[jr*ldz+j+1] + s*z[jr*ldz+j]
					z[jr*ldz+j] = -s*z[jr*ldz+j+1] + c*z[jr*ldz+j]
					z[jr*ldz+j+1] = temp
				}
			}
		}
		goto ThreeFifty

		// Use Francis double-shift.

	TwoHundred:
		if ifirst+1 == ilast {
			// Special case -- 2x2 block with complex eigenvectors.
			// Step 1: Standardize, that is, rotate so that
			// B =  (B11  0 )
			//      ( 0  B22)   With B11 non-negative.
			b22, b11, sr, cr, sl, cl = impl.Dlasv2(t[(ilast-1)*ldt+ilast-1], t[(ilast-1)*ldt+ilast], t[ilast*ldt+ilast])
			if b11 < 0 {
				cr = -cr
				sr = -sr
				b11 = -b11
				b22 = -b22
			}
			bi.Drot(ilastm+1-ifirst, h[(ilast-1)*ldh+ilast-1:], 1, h[ilast*ldh+ilast-1:], 1, cl, sl)
			bi.Drot(ilast+1-ifrstm, h[ifrstm*ldh+ilast-1:], ldh, h[ifrstm*ldh+ilast:], ldh, cr, sr)

			if ilast < ilastm {
				bi.Drot(ilastm-ilast, t[(ilast-1)*ldt+ilast+1:], 1, t[ilast*ldt+ilast+1:], 1, cl, sl)
			}
			if ifrstm < ilast-1 {
				bi.Drot(ifirst-ifrstm, t[ifrstm*ldt+ilast-1:], ldt, t[ifrstm*ldt+ilast:], ldt, cr, sr)
			}

			if ilq {
				bi.Drot(n, q[ilast-1:], ldq, q[ilast:], ldq, cl, sl)
			}
			if ilz {
				bi.Drot(n, z[ilast-1:], ldz, z[ilast:], ldz, cr, sr)
			}

			t[(ilast-1)*ldt+ilast-1] = b11
			t[(ilast-1)*ldt+ilast] = 0
			t[ilast*ldt+ilast-1] = 0
			t[ilast*ldt+ilast] = b22

			// If B22 is negative, negate column ilast.
			if b22 < 0 {
				b22 = -b22
				for j := ifrstm; j <= ilast; j++ {
					h[j*ldh+ilast] *= -1
					t[j*ldt+ilast] *= -1
				}
				if ilz {
					for j := 0; j < n; j++ {
						z[j*ldz+ilast] *= -1
					}
				}
			}

			// Step 2: compute alphar, alphai, and beta.
			// Recompute shift.
			s1, _, wr, _, wi = impl.Dlag2(h[(ilast-1)*ldh+ilast-1:], ldh, t[(ilast-1)*ldt+ilast-1:], ldt)
			if wi == 0 {
				// If standardization has perturbed the shift onto real line, do another QR step.
				goto ThreeFifty
			}
			s1inv = 1 / s1

			// Do EISPACK (QZVAL) computation of alpha and beta.
			a11 = h[(ilast-1)*ldh+ilast-1]
			a21 = h[ilast*ldh+ilast-1]
			a12 = h[(ilast-1)*ldh+ilast]
			a22 = h[ilast*ldh+ilast]

			// Compute complex Givens rotation on right assuming some element of C = (sA -wB)>unfl:
			// ( sA - wB ) (  CZ    -^SZ  )
			//             (  SZ     CZ  )

			c11r = s1*a11 - wr*b11
			c11i = -wi * b11
			c12 = s1 * a12
			c21 = s1 * a21
			c22r = s1*a22 - wr*b22
			c22i = -wi * b22

			if math.Abs(c11r)+math.Abs(c11i)+math.Abs(c12) > math.Abs(c21)+math.Abs(c22r)+math.Abs(c22i) {
				t1 = impl.Dlapy3(c12, c11r, c11i)
				cz = c12 / t1
				szr = -c11r / t1
				szi = -c11i / t1
			} else {
				cz = impl.Dlapy2(c22r, c22i)
				if cz <= safmin {
					cz = 0
					szr = 1
					szi = 0
				} else {
					tempr = c22r / cz
					tempi = c22i / cz
					t1 = impl.Dlapy2(cz, c21)
					cz = cz / t1
					szr = -c21 * tempr / t1
					szi = c21 * tempi / t1
				}
			}

			// Compute Givens rotation on left
			// ( CQ   SQ )
			// (-^SQ   CQ )   A or B.

			an = math.Abs(a11) + math.Abs(a12) + math.Abs(a21) + math.Abs(a22)
			bn = math.Abs(b11) + math.Abs(b22)
			wabs = math.Abs(wr) + math.Abs(wi)
			if s1*an > wabs*bn {
				cq = cz * b11
				sqr = szr * b22
				sqi = -szi * b22
			} else {
				a1r = cz*a11 + szr*a12
				a1i = szi * a12
				a2r = cz*a21 + szr*a22
				a2i = szi * a22
				cq = impl.Dlapy2(a1r, a1i)
				if cq <= safmin {
					cq = 0
					sqr = 1
					sqi = 0
				} else {
					tempr = a1r / cq
					tempi = a1i / cq
					sqr = tempr*a2r + tempi*a2i
					sqi = tempi*a2r - tempr*a2i
				}
			}
			t1 = impl.Dlapy3(cq, sqr, sqi)
			cq /= t1
			sqr /= t1
			sqi /= t1

			// Compute diagonal elements of QBZ.
			tempr = sqr*szr - sqi*szi
			tempi = sqr*szi + sqi*szr
			b1r = cq*cz*b11 + tempr*b22
			b1i = tempi * b22
			b1a = impl.Dlapy2(b1r, b1i)
			b2r = cq*cz*b22 + tempr*b11
			b2i = -tempi * b11
			b2a = impl.Dlapy2(b2r, b2i)

			// Normalize so beta>0 and imag(alpha1) > 0.
			beta[ilast-1] = b1a
			beta[ilast] = b2a
			alphar[ilast-1] = (wr * b1a) * s1inv
			alphai[ilast-1] = (wi * b1a) * s1inv
			alphar[ilast] = (wr * b2a) * s1inv
			alphai[ilast] = -(wi * b2a) * s1inv

			// Step 3: Go to next block -- exit if finished.
			ilast = ifirst - 1
			if ilast < ilo {
				goto ThreeEighty
			}

			// Reset counters.
			iiter = 0
			eshift = 0.0
			if !ilschr {
				ilastm = ilast
				if ifrstm > ilast {
					ifrstm = ilo
				}
			}
			goto ThreeFifty
		} else {
			// Usual case: 3x3 or larger block, using Francis implicit double shift.
			//
			// Eigenvalue equation is w² - c w + d = 0,
			// so compute 1st column of (A B⁻¹)² - c A B⁻¹ + d
			// using the formula in QZIT (from EISPACK).

			ad11 = (ascale * h[(ilast-1)*ldh+ilast-1]) /
				(bscale * t[(ilast-1)*ldt+ilast-1])
			ad21 = (ascale * h[(ilast)*ldh+ilast-1]) /
				(bscale * t[(ilast-1)*ldt+ilast-1])
			ad12 = (ascale * h[(ilast-1)*ldh+ilast]) /
				(bscale * t[(ilast)*ldt+ilast])
			ad22 = (ascale * h[(ilast)*ldh+ilast]) /
				(bscale * t[(ilast)*ldt+ilast])
			u12 = t[(ilast-1)*ldt+ilast] / t[(ilast)*ldt+ilast]
			ad11l = (ascale * h[(ifirst)*ldh+ifirst]) /
				(bscale * t[(ifirst)*ldt+ifirst])
			ad21l = (ascale * h[(ifirst+1)*ldh+ifirst]) /
				(bscale * t[(ifirst)*ldt+ifirst])
			ad12l = (ascale * h[(ifirst)*ldh+ifirst+1]) /
				(bscale * t[(ifirst+1)*ldt+ifirst+1])
			ad22l = (ascale * h[(ifirst+1)*ldh+ifirst+1]) /
				(bscale * t[(ifirst+1)*ldt+ifirst+1])
			ad32l = (ascale * h[(ifirst+2)*ldh+ifirst+1]) /
				(bscale * t[(ifirst+1)*ldt+ifirst+1])
			u12l = t[(ifirst)*ldt+ifirst+1] / t[(ifirst+1)*ldt+ifirst+1]

			v[0] = (ad11-ad11l)*(ad22-ad11l) - ad12*ad21 +
				ad21*u12*ad11l + (ad12l-ad11l*u12l)*ad21l
			v[1] = ((ad22l - ad11l) - ad21l*u12l - (ad11 - ad11l) -
				(ad22 - ad11l) + ad21*u12) * ad21l
			v[2] = ad32l * ad21l
			istart = ifirst

			_, tau = impl.Dlarfg(3, v[0], v[1:], 1)
			v[0] = 1

			// Sweep.
			for j = istart; j <= ilast-2; j++ { // Loop 290.
				// All but last elements: use 3x3 Householder transforms.
				if j > istart {
					v[0] = h[j*ldh+j-1]
					v[1] = h[(j+1)*ldh+j-1]
					v[2] = h[(j+2)*ldh+j-1]

					h[j*ldh+j-1], tau = impl.Dlarfg(3, h[j*ldh+j-1], v[1:], 1)
					v[0] = 1
					h[(j+1)*ldh+j-1] = 0
					h[(j+2)*ldh+j-1] = 0
				}
				t2 = tau * v[1]
				t3 = tau * v[2]
				for jc := j; jc <= ilastm; jc++ {
					temp = h[j*ldh+jc] + v[1]*h[(j+1)*ldh+jc] + v[2]*h[(j+2)*ldh+jc]
					h[j*ldh+jc] -= temp * tau
					h[(j+1)*ldh+jc] -= temp * t2
					h[(j+2)*ldh+jc] -= temp * t3
					temp2 = t[j*ldt+jc] + v[1]*t[(j+1)*ldt+jc] + v[2]*t[(j+2)*ldt+jc]
					t[j*ldt+jc] -= temp2 * tau
					t[(j+1)*ldt+jc] -= temp2 * t2
					t[(j+2)*ldt+jc] -= temp2 * t3
				}
				if ilq {
					for jr := 0; jr < n; jr++ {
						temp = q[jr*ldq+j] + v[1]*q[jr*ldq+j+1] + v[2]*q[jr*ldq+j+2]
						q[jr*ldq+j] -= temp * tau
						q[jr*ldq+j+1] -= temp * t2
						q[jr*ldq+j+2] -= temp * t3
					}
				}

				// Zero j-th column of B (see Dlagbc for details).
				// Swap rows to pivot.

				ilpivt = false
				temp = math.Max(math.Abs(t[(j+1)*ldt+j+1]), math.Abs(t[(j+1)*ldt+j+2]))
				temp2 = math.Max(math.Abs(t[(j+2)*ldt+j+1]), math.Abs(t[(j+2)*ldt+j+2]))
				if math.Max(temp, temp2) < safmin {
					scale = 0
					u1 = 1
					u2 = 0
					goto TwoFifty
				} else if temp >= temp2 {
					w11 = t[(j+1)*ldt+j+1]
					w21 = t[(j+2)*ldt+j+1]
					w12 = t[(j+1)*ldt+j+2]
					w22 = t[(j+2)*ldt+j+2]
					u1 = t[(j+1)*ldt+j]
					u2 = t[(j+2)*ldt+j]
				} else {
					w21 = t[(j+1)*ldt+j+1]
					w11 = t[(j+2)*ldt+j+1]
					w22 = t[(j+1)*ldt+j+2]
					w12 = t[(j+2)*ldt+j+2]
					u2 = t[(j+1)*ldt+j]
					u1 = t[(j+2)*ldt+j]
				}

				// Swap columns if necessary.
				if math.Abs(w12) > math.Abs(w11) {
					ilpivt = true
					w12, w11 = w11, w12
					w22, w21 = w21, w22
				}

				// LU Factor.
				temp = w21 / w11
				u2 -= temp * u1
				w22 -= temp * w12
				w21 = 0

				// Compute scale.
				scale = 1
				if math.Abs(w22) < safmin {
					scale = 0
					u2 = 1
					u1 = -w12 / w11
					goto TwoFifty
				}
				if math.Abs(w22) < math.Abs(u2) {
					scale = math.Abs(w22 / u2)
				}
				if math.Abs(w11) < math.Abs(u1) {
					scale = math.Min(scale, math.Abs(w11/u1))
				}

				// Solve.
				u2 = (scale * u2) / w22
				u1 = (scale*u1 - w12*u2) / w11

			TwoFifty: // Continue 250.

				if ilpivt {
					u1, u2 = u2, u1
				}

				// Compute Householder Vector.
				t1 = math.Sqrt(scale*scale + u1*u1 + u2*u2)
				tau = 1 + scale/t1
				vs = -1 / (scale + t1)
				v[0] = 1
				v[1] = vs * u1
				v[2] = vs * u2

				// Apply transformations from the right.
				t2 = tau * v[1]
				t3 = tau * v[2]
				jrmax := min(j+3, ilast)
				for jr := ifrstm; jr <= jrmax; jr++ {
					temp = h[jr*ldh+j] + v[1]*h[jr*ldh+j+1] + v[2]*h[jr*ldh+j+2]
					h[jr*ldh+j] -= temp * tau
					h[jr*ldh+j+1] -= temp * t2
					h[jr*ldh+j+2] -= temp * t3
				}
				for jr := ifrstm; jr <= j+2; jr++ {
					temp = t[jr*ldt+j] + v[1]*t[jr*ldt+j+1] + v[2]*t[jr*ldt+j+2]
					t[jr*ldt+j] -= temp * tau
					t[jr*ldt+j+1] -= temp * t2
					t[jr*ldt+j+2] -= temp * t3
				}
				if ilz {
					for jr := 0; jr < n; jr++ {
						temp = tau * (z[jr*ldz+j] + v[1]*z[jr*ldz+j+1] + v[2]*z[jr*ldz+j+2])
						z[jr*ldz+j] -= temp * tau
						z[jr*ldz+j+1] -= temp * t2
						z[jr*ldz+j+2] -= temp * t3
					}
				}
				t[(j+1)*ldt+j] = 0
				t[(j+2)*ldt+j] = 0
			} // Continue 290.

			// Last elements: Use Givens rotations.
			// Rotations from the left.

			j = ilast - 1
			temp = h[j*ldh+j-1]
			c, s, h[j*ldh+j-1] = impl.Dlartg(temp, h[(j+1)*ldh+j-1])
			h[(j+1)*ldh+j-1] = 0

			for jc := j; jc <= ilastm; jc++ {
				temp = c*h[j*ldh+jc] + s*h[(j+1)*ldh+jc]
				h[(j+1)*ldh+jc] = -s*h[j*ldh+jc] + c*h[(j+1)*ldh+jc]
				h[j*ldh+jc] = temp
				temp2 = c*t[j*ldt+jc] + s*t[(j+1)*ldt+jc]
				t[(j+1)*ldt+jc] = -s*t[j*ldt+jc] + c*t[(j+1)*ldt+jc]
				t[j*ldt+jc] = temp2
			}
			if ilq {
				for jr := 0; jr < n; jr++ {
					temp = c*q[jr*ldq+j] + s*q[jr*ldq+j+1]
					q[jr*ldq+j+1] = -s*q[jr*ldq+j] + c*q[jr*ldq+j+1]
					q[jr*ldq+j] = temp
				}
			}

			// Rotations from the right.
			temp = t[(j+1)*ldt+j+1]
			c, s, t[(j+1)*ldt+j+1] = impl.Dlartg(temp, t[(j+1)*ldt+j])
			t[(j+1)*ldt+j] = 0

			for jr := ifrstm; jr <= ilast; jr++ {
				temp = c*h[jr*ldh+j+1] + s*h[jr*ldh+j]
				h[jr*ldh+j] = -s*h[jr*ldh+j+1] + c*h[jr*ldh+j]
				h[jr*ldh+j+1] = temp
			}
			for jr := ifrstm; jr <= ilast-1; jr++ {
				temp = c*t[jr*ldt+j+1] + s*t[jr*ldt+j]
				t[jr*ldt+j] = -s*t[jr*ldt+j+1] + c*t[jr*ldt+j]
				t[jr*ldt+j+1] = temp
			}
			if ilz {
				for jr := 0; jr < n; jr++ {
					temp = c*z[jr*ldz+j+1] + s*z[jr*ldz+j]
					z[jr*ldz+j] = -s*z[jr*ldz+j+1] + c*z[jr*ldz+j]
					z[jr*ldz+j+1] = temp
				}
			}
			// End of double-shift code.
		} // End of big if.

	ThreeFifty: // Continue with next QZ iteration (jiter++).

	}
	// Drop through, non convergence.
	work[0] = float64(n)
	return ilast

ThreeEighty:
	// Succesful completion of all QZ steps.
	// Set Eigenvalues 1:ILO-1
	for j = 0; j <= ilo-1; j++ {
		if t[j*ldt+j] < 0 {
			if ilschr {
				for jr := 0; jr <= j; jr++ {
					h[jr*ldh+j] *= -1
					t[jr*ldt+j] *= -1
				}
			} else {
				h[j*ldh+j] *= -1
				t[j*ldt+j] *= -1
			}
			if ilz {
				for jr := 0; jr < n; jr++ {
					z[jr*ldz+j] *= -1
				}
			}
		}
		alphar[j] = h[j*ldh+j]
		alphai[j] = 0
		beta[j] = t[j*ldt+j]
	}

	return -1 // Normal termination.
}

// _column copies the j-th column of z into a new slice of length m.
func _column(z []float64, ldz, j, m int) []float64 {
	v := make([]float64, m)
	for i := range v {
		v[i] = z[i*ldz+j]
	}
	return v
}

func _columns(z []float64, ldz, n, m int) [][]float64 {
	v := make([][]float64, n)
	for i := range v {
		v[i] = _column(z, ldz, i, m)
	}
	return v
}
