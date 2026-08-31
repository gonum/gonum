// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"math"

	"gonum.org/v1/gonum/blas/blas64"
)

// Dlaic1 applies one step of incremental condition estimation in its simplest
// version. It is called as a part of the condition estimator in LAPACK routines
// such as Dtrcon and Dtrsna.
//
// Let x, op(A)*x = sest*w, where op(A) is the triangular matrix A or its
// transpose. Dlaic1 computes sestpr, s, c such that the vector
//
//	[s*x]
//	[ c ]
//
// is an approximate singular vector of
//
//	[op(A)  w]
//	[  0    γ]
//
// in the sense that
//
//	diag(sest, sestpr) ≈ sigma(op(A_hat))
//
// where A_hat is the (j+1)×(j+1) matrix
//
//	A_hat = [op(A)  w]
//	        [  0    γ]
//
// and sigma(A_hat) are its singular values.
//
// If job == 1, Dlaic1 computes an updated estimate of the largest singular
// value. If job == 2, it computes an updated estimate of the smallest singular
// value.
//
// j is the length of vectors x and w. x and w must have length at least j when
// j > 0, otherwise Dlaic1 will panic. sest must be non-negative.
//
// On return, sestpr is the updated singular value estimate, and s, c are the
// sine and cosine of the rotation used in the update.
//
// Dlaic1 is an internal routine. It is exported for testing purposes.
func (impl Implementation) Dlaic1(job int, j int, x []float64, sest float64, w []float64, gamma float64) (sestpr, s, c float64) {
	switch {
	case job != 1 && job != 2:
		panic("lapack: bad job")
	case j < 0:
		panic("lapack: j < 0")
	}
	if j > 0 {
		if len(x) < j {
			panic(shortX)
		}
		if len(w) < j {
			panic(shortW)
		}
	}

	eps := dlamchE

	bi := blas64.Implementation()
	alpha := bi.Ddot(j, x, 1, w, 1)

	absalp := math.Abs(alpha)
	absgam := math.Abs(gamma)
	absest := math.Abs(sest)

	if job == 1 {
		if sest == 0 {
			s1 := math.Max(absgam, absalp)
			if s1 == 0 {
				s = 0
				c = 1
				sestpr = 0
			} else {
				s = alpha / s1
				c = gamma / s1
				tmp := math.Sqrt(s*s + c*c)
				s /= tmp
				c /= tmp
				sestpr = s1 * tmp
			}
			return
		} else if absgam <= eps*absest {
			s = 1
			c = 0
			tmp := math.Max(absest, absalp)
			s1 := absest / tmp
			s2 := absalp / tmp
			sestpr = tmp * math.Sqrt(s1*s1+s2*s2)
			return
		} else if absalp <= eps*absest {
			s1 := absgam
			s2 := absest
			if s1 <= s2 {
				s = 1
				c = 0
				sestpr = s2
			} else {
				s = 0
				c = 1
				sestpr = s1
			}
			return
		} else if absest <= eps*absalp || absest <= eps*absgam {
			s1 := absgam
			s2 := absalp
			if s1 <= s2 {
				tmp := s1 / s2
				s = math.Sqrt(1 + tmp*tmp)
				sestpr = s2 * s
				c = (gamma / s2) / s
				s = math.Copysign(1, alpha) / s
			} else {
				tmp := s2 / s1
				c = math.Sqrt(1 + tmp*tmp)
				sestpr = s1 * c
				s = (alpha / s1) / c
				c = math.Copysign(1, gamma) / c
			}
			return
		}

		zeta1 := alpha / absest
		zeta2 := gamma / absest

		b := (1 - zeta1*zeta1 - zeta2*zeta2) * 0.5
		c = zeta1 * zeta1
		var t float64
		if b > 0 {
			t = c / (b + math.Sqrt(b*b+c))
		} else {
			t = math.Sqrt(b*b+c) - b
		}

		sine := -zeta1 / t
		cosine := -zeta2 / (1 + t)
		tmp := math.Sqrt(sine*sine + cosine*cosine)
		s = sine / tmp
		c = cosine / tmp
		sestpr = math.Sqrt(t+1) * absest
		return
	}

	// job == 2
	if sest == 0 {
		sestpr = 0
		if math.Max(absgam, absalp) == 0 {
			sine := 1.0
			cosine := 0.0
			s1 := math.Max(math.Abs(sine), math.Abs(cosine))
			s = sine / s1
			c = cosine / s1
			tmp := math.Sqrt(s*s + c*c)
			s /= tmp
			c /= tmp
		} else {
			sine := -gamma
			cosine := alpha
			s1 := math.Max(math.Abs(sine), math.Abs(cosine))
			s = sine / s1
			c = cosine / s1
			tmp := math.Sqrt(s*s + c*c)
			s /= tmp
			c /= tmp
		}
		return
	} else if absgam <= eps*absest {
		s = 0
		c = 1
		sestpr = absgam
		return
	} else if absalp <= eps*absest {
		s1 := absgam
		s2 := absest
		if s1 <= s2 {
			s = 0
			c = 1
			sestpr = s1
		} else {
			s = 1
			c = 0
			sestpr = s2
		}
		return
	} else if absest <= eps*absalp || absest <= eps*absgam {
		s1 := absgam
		s2 := absalp
		if s1 <= s2 {
			tmp := s1 / s2
			c = math.Sqrt(1 + tmp*tmp)
			sestpr = absest * (tmp / c)
			s = -(gamma / s2) / c
			c = math.Copysign(1, alpha) / c
		} else {
			tmp := s2 / s1
			s = math.Sqrt(1 + tmp*tmp)
			sestpr = absest / s
			c = (alpha / s1) / s
			s = -math.Copysign(1, gamma) / s
		}
		return
	}

	zeta1 := alpha / absest
	zeta2 := gamma / absest

	norma := math.Max(
		1+zeta1*zeta1+math.Abs(zeta1*zeta2),
		math.Abs(zeta1*zeta2)+zeta2*zeta2,
	)

	test := 1 + 2*(zeta1-zeta2)*(zeta1+zeta2)
	var sine, cosine float64
	if test >= 0 {
		b := (zeta1*zeta1 + zeta2*zeta2 + 1) * 0.5
		c = zeta2 * zeta2
		t := c / (b + math.Sqrt(math.Abs(b*b-c)))
		sine = zeta1 / (1 - t)
		cosine = -zeta2 / t
		sestpr = math.Sqrt(t+4*eps*eps*norma) * absest
	} else {
		b := (zeta2*zeta2 + zeta1*zeta1 - 1) * 0.5
		c = zeta1 * zeta1
		var t float64
		if b >= 0 {
			t = -c / (b + math.Sqrt(b*b+c))
		} else {
			t = b - math.Sqrt(b*b+c)
		}
		sine = -zeta1 / t
		cosine = -zeta2 / (1 + t)
		sestpr = math.Sqrt(1+t+4*eps*eps*norma) * absest
	}
	tmp := math.Sqrt(sine*sine + cosine*cosine)
	s = sine / tmp
	c = cosine / tmp
	return
}
