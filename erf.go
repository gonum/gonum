// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mathext

import "math"

// NormalQuantile computes the quantile function (inverse CDF) of the standard
// normal. NormalQuantile panics if the input p is less than 0 or greater than 1.
func NormalQuantile(p float64) float64 {
	switch {
	case p < 0 || 1 < p:
		panic("mathext: quantile out of bounds")
	case p == 1:
		return math.Inf(1)
	case p == 0:
		return math.Inf(-1)
	}
	return zQuantile(p)
}

/*
Copyright (c) 2012 The Probab Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

* Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
* Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
* Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

var (
	zQuantSmallA = []float64{3.387132872796366608, 133.14166789178437745, 1971.5909503065514427, 13731.693765509461125, 45921.953931549871457, 67265.770927008700853, 33430.575583588128105, 2509.0809287301226727}
	zQuantSmallB = []float64{1.0, 42.313330701600911252, 687.1870074920579083, 5394.1960214247511077, 21213.794301586595867, 39307.89580009271061, 28729.085735721942674, 5226.495278852854561}
	zQuantInterA = []float64{1.42343711074968357734, 4.6303378461565452959, 5.7694972214606914055, 3.64784832476320460504, 1.27045825245236838258, 0.24178072517745061177, 0.0227238449892691845833, 7.7454501427834140764e-4}
	zQuantInterB = []float64{1.0, 2.05319162663775882187, 1.6763848301838038494, 0.68976733498510000455, 0.14810397642748007459, 0.0151986665636164571966, 5.475938084995344946e-4, 1.05075007164441684324e-9}
	zQuantTailA  = []float64{6.6579046435011037772, 5.4637849111641143699, 1.7848265399172913358, 0.29656057182850489123, 0.026532189526576123093, 0.0012426609473880784386, 2.71155556874348757815e-5, 2.01033439929228813265e-7}
	zQuantTailB  = []float64{1.0, 0.59983220655588793769, 0.13692988092273580531, 0.0148753612908506148525, 7.868691311456132591e-4, 1.8463183175100546818e-5, 1.4215117583164458887e-7, 2.04426310338993978564e-15}
)

func rateval(a []float64, na int64, b []float64, nb int64, x float64) float64 {
	var (
		u, v, r float64
	)
	u = a[na-1]

	for i := na - 1; i > 0; i-- {
		u = x*u + a[i-1]
	}

	v = b[nb-1]

	for j := nb - 1; j > 0; j-- {
		v = x*v + b[j-1]
	}

	r = u / v

	return r
}

func zQuantSmall(q float64) float64 {
	r := 0.180625 - q*q
	return q * rateval(zQuantSmallA, 8, zQuantSmallB, 8, r)
}

func zQuantIntermediate(r float64) float64 {
	return rateval(zQuantInterA, 8, zQuantInterB, 8, (r - 1.6))
}

func zQuantTail(r float64) float64 {
	return rateval(zQuantTailA, 8, zQuantTailB, 8, (r - 5.0))
}

// Compute the quantile in normalized units
func zQuantile(p float64) float64 {
	switch {
	case p == 1.0:
		return math.Inf(1)
	case p == 0.0:
		return math.Inf(-1)
	}
	var r, x, pp, dp float64
	dp = p - 0.5
	if math.Abs(dp) <= 0.425 {
		return zQuantSmall(dp)
	}
	if p < 0.5 {
		pp = p
	} else {
		pp = 1.0 - p
	}
	r = math.Sqrt(-math.Log(pp))
	if r <= 5.0 {
		x = zQuantIntermediate(r)
	} else {
		x = zQuantTail(r)
	}
	if p < 0.5 {
		return -x
	}
	return x
}
