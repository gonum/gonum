// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testblas

import "math"

// You Can Write FORTRAN in any Language.
//
// drotmg_f77 is a direct translation of the FORTRAN 77 drotmg reference implementation.
//
// SHOUTING comments are from the original code.
//
//     CONSTRUCT THE MODIFIED GIVENS TRANSFORMATION MATRIX H WHICH ZEROS
//     THE SECOND COMPONENT OF THE 2-VECTOR  (DSQRT(DD1)*DX1,DSQRT(DD2)*
//     DY2)**T.
//     WITH DPARAM(1)=DFLAG, H HAS ONE OF THE FOLLOWING FORMS..
//
//     DFLAG=-1.D0     DFLAG=0.D0        DFLAG=1.D0     DFLAG=-2.D0
//
//       (DH11  DH12)    (1.D0  DH12)    (DH11  1.D0)    (1.D0  0.D0)
//     H=(          )    (          )    (          )    (          )
//       (DH21  DH22),   (DH21  1.D0),   (-1.D0 DH22),   (0.D0  1.D0).
//     LOCATIONS 2-4 OF DPARAM CONTAIN DH11, DH21, DH12, AND DH22
//     RESPECTIVELY. (VALUES OF 1.D0, -1.D0, OR 0.D0 IMPLIED BY THE
//     VALUE OF DPARAM(1) ARE NOT STORED IN DPARAM.)
//
//     THE VALUES OF GAMSQ AND RGAMSQ SET IN THE DATA STATEMENT MAY BE
//     INEXACT.  THIS IS OK AS THEY ARE ONLY USED FOR TESTING THE SIZE
//     OF DD1 AND DD2.  ALL ACTUAL SCALING OF DATA IS DONE USING GAM.
//
// Note that the documentation above refers to a DY2 which does not exist.
func drotmg_f77(d1, d2, x1, y1 float64) (rd1, rd2, rx1 float64, p [5]float64) {
	var p1, p2, q1, q2, u, tmp float64
	var flag, h11, h12, h21, h22 float64

	const (
		// Constants are type float64 to force similar
		// behaviour to the FORTRAN 77 as much as possible.
		gam    float64 = 4096
		gamsq  float64 = 16777216
		rgamsq float64 = 5.9604645e-8
	)

	// Simulate assigned goto with a switch below.
	var igo int

	if !(d1 < 0) {
		goto L10
	}
	// GO ZERO-H-D-AND-DX1..
	goto L60
L10:
	// CASE-DD1-NONNEGATIVE
	p2 = d2 * y1
	if !(p2 == 0) {
		goto L20
	}
	flag = -2
	goto L260
	// REGULAR-CASE..
L20:
	p1 = d1 * x1
	q2 = p2 * y1
	q1 = p1 * x1
	if !(math.Abs(q1) > math.Abs(q2)) {
		goto L40
	}
	h21 = -y1 / x1
	h12 = p2 / p1
	u = 1 - h12*h21
	if !(u <= 0) {
		goto L30
	}
	// GO ZERO-H-D-AND-DX1..
	goto L60
L30:
	flag = 0
	d1 = d1 / u
	d2 = d2 / u
	x1 = x1 * u
	// GO SCALE-CHECK..
	goto L100
L40:
	if !(q2 < 0) {
		goto L50
	}
	// GO ZERO-H-D-AND-DX1..
	goto L60
L50:
	flag = 1
	h11 = p1 / p2
	h22 = x1 / y1
	u = 1 + h11*h22
	tmp = d2 / u
	d2 = d1 / u
	d1 = tmp
	x1 = y1 * u
	// GO SCALE-CHECK
	goto L100
	// PROCEDURE..ZERO-H-D-AND-DX1..
L60:
	flag = -1
	h11 = 0
	h12 = 0
	h21 = 0
	h22 = 0
	d1 = 0
	d2 = 0
	x1 = 0
	// RETURN..
	goto L220
	// PROCEDURE..FIX-H..
L70:
	if !(flag >= 0) {
		goto L90
	}
	if !(flag == 0) {
		goto L80
	}
	h11 = 1
	h22 = 1
	flag = -1
	goto L90
L80:
	h21 = -1
	h12 = 1
	flag = -1
L90:
	switch igo {
	case 120:
		goto L120
	case 150:
		goto L150
	case 180:
		goto L180
	case 210:
		goto L210
	default:
		panic("igo not assigned")
	}
	// PROCEDURE..SCALE-CHECK
L100:
L110:
	if !(d1 <= rgamsq) {
		goto L130
	}
	if d1 == 0 {
		goto L160
	}
	igo = 120
	// FIX-H..
	goto L70
L120:
	d1 = d1 * (gam * gam)
	x1 = x1 / gam
	h11 = h11 / gam
	h12 = h12 / gam
	goto L110
L130:
L140:
	if !(d1 >= gamsq) {
		goto L160
	}
	igo = 150
	// FIX-H..
	goto L70
L150:
	d1 = d1 / (gam * gam)
	x1 = x1 * gam
	h11 = h11 * gam
	h12 = h12 * gam
	goto L140
L160:
L170:
	if !(math.Abs(d2) <= rgamsq) {
		goto L190
	}
	if d2 == 0 {
		goto L220
	}
	igo = 180
	// FIX-H..
	goto L70
L180:
	d2 = d2 * (gam * gam)
	h21 = h21 / gam
	h22 = h22 / gam
	goto L170
L190:
L200:
	if !(math.Abs(d2) >= gamsq) {
		goto L220
	}
	igo = 210
	// FIX-H..
	goto L70
L210:
	d2 = d2 / (gam * gam)
	h21 = h21 * gam
	h22 = h22 * gam
	goto L200
L220:
	switch {
	case flag < 0:
		goto L250
	default:
		goto L230
	case flag > 0:
		goto L240
	}
L230:
	p[2] = h21
	p[3] = h12
	goto L260
L240:
	p[1] = h11
	p[4] = h22
	goto L260
L250:
	p[1] = h11
	p[2] = h21
	p[3] = h12
	p[4] = h22
L260:
	p[0] = flag
	return d1, d2, x1, p
}
