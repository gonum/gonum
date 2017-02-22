// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
 * Cephes Math Library Release 2.1:  January, 1989
 * Copyright 1984, 1987, 1989 by Stephen L. Moshier
 * Direct inquiries to 30 Frost Street, Cambridge, MA 02140
 */

package cephes

import "math"

/*
 * Adapted from scipy's cephes zeta.c
 */

/*
 *     Riemann zeta function of two arguments
 *
 *
 * DESCRIPTION:
 *
 *                 inf.
 *                  -        -x
 *   zeta(x,q)  =   >   (k+q)
 *                  -
 *                 k=0
 *
 * where x > 1 and q is not a negative integer or zero.
 * The Euler-Maclaurin summation formula is used to obtain
 * the expansion
 *
 *                n
 *                -       -x
 * zeta(x,q)  =   >  (k+q)
 *                -
 *               k=1
 *
 *           1-x                 inf.  B   x(x+1)...(x+2j)
 *      (n+q)           1         -     2j
 *  +  ---------  -  -------  +   >    --------------------
 *        x-1              x      -                   x+2j+1
 *                   2(n+q)      j=1       (2j)! (n+q)
 *
 * where the B2j are Bernoulli numbers.  Note that (see zetac.c)
 * zeta(x,1) = zetac(x) + 1.
 *
 *
 * REFERENCE:
 *
 * Gradshteyn, I. S., and I. M. Ryzhik, Tables of Integrals,
 * Series, and Products, p. 1073; Academic Press, 1980.
 *
 */

/* Expansion coefficients
 * for Euler-Maclaurin summation formula
 * (2k)! / B2k
 * where B2k are Bernoulli numbers
 */
var zetaCoefs = []float64{
	12.0,
	-720.0,
	30240.0,
	-1209600.0,
	47900160.0,
	-1.8924375803183791606e9, /*1.307674368e12/691 */
	7.47242496e10,
	-2.950130727918164224e12,  /*1.067062284288e16/3617 */
	1.1646782814350067249e14,  /*5.109094217170944e18/43867 */
	-4.5979787224074726105e15, /*8.028576626982912e20/174611 */
	1.8152105401943546773e17,  /*1.5511210043330985984e23/854513 */
	-7.1661652561756670113e18, /*1.6938241367317436694528e27/236364091 */
}

/* 30 Nov 86 -- error in third coefficient fixed */

// Zeta calculates the Riemann zeta function of two arguments
func Zeta(x, q float64) float64 {
	if x == 1 {
		return math.MaxFloat64
	}

	if x < 1 {
		panic(badParamOutOfBounds)
	}

	if q <= 0 {
		if q == math.Floor(q) {
			panic(badParamFunctionSingularity)
		}
		if x != math.Floor(x) {
			panic(badParamOutOfBounds) // because q^-x not defined
		}
	}

	/* Asymptotic expansion
	 * http://dlmf.nist.gov/25.11#E43
	 */
	if q > 1e8 {
		return (1/(x-1) + 1/(2*q)) * math.Pow(q, 1-x)
	}

	// Euler-Maclaurin summation formula

	/* Permit negative q but continue sum until n+q > +9 .
	 * This case should be handled by a reflection formula.
	 * If q<0 and x is an integer, there is a relation to
	 * the polyGamma function.
	 */
	s := math.Pow(q, -x)
	a := q
	i := 0
	b := 0.0
	for i < 9 || a <= 9 {
		i++
		a += 1.0
		b = math.Pow(a, -x)
		s += b
		if math.Abs(b/s) < machEp {
			return s
		}
	}

	w := a
	s += b * w / (x - 1)
	s -= 0.5 * b
	a = 1.0
	k := 0.0
	for i = 0; i < 12; i++ {
		a *= x + k
		b /= w
		t := a * b / zetaCoefs[i]
		s = s + t
		t = math.Abs(t / s)
		if t < machEp {
			return s
		}
		k += 1.0
		a *= x + k
		b /= w
		k += 1.0
	}
	return s
}
