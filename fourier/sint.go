// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is a translation of the FFTPACK sint functions by
// Paul N Swarztrauber, placed in the public domain at
// http://www.netlib.org/fftpack/.

package fourier

import "math"

// subroutine sinti initializes the array work which is used in
// subroutine sint. the prime factorization of n together with
// a tabulation of the trigonometric functions are computed and
// stored in work.
//
// input parameter
//
// n       the length of the sequence to be transformed.  the method
//         is most efficient when n+1 is a product of small primes.
//
// output parameter
//
// work    a work array with at least ceil(2.5*n) locations.
//         different work arrays are required for different values
//         of n. the contents of work must not be changed between
//         calls of sint.
//
// ifac    an integer work array of length at least 15.
func sinti(n int, work []float64, ifac []int) {
	if len(work) < 5*(n+1)/2 {
		panic("fourier: short work")
	}
	if len(ifac) < 15 {
		panic("fourier: short ifac")
	}
	if n <= 1 {
		return
	}
	dt := math.Pi / float64(n+1)
	for k := 0; k < n/2; k++ {
		work[k] = 2 * math.Sin(float64(k+1)*dt)
	}
	rffti(n+1, work[n/2:], ifac)
}

// subroutine sint computes the discrete fourier sine transform
// of an odd sequence x(i). the transform is defined below at
// output parameter x.
//
// sint is the unnormalized inverse of itself since a call of sint
// followed by another call of sint will multiply the input sequence
// x by 2*(n+1).
//
// the array work which is used by subroutine sint must be
// initialized by calling subroutine sinti(n,work).
//
// input parameters
//
// n       the length of the sequence to be transformed.  the method
//         is most efficient when n+1 is the product of small primes.
//
// x       an array which contains the sequence to be transformed
//
//
// work    a work array with dimension at least ceil(2.5*n)
//         in the program that calls sint. the work array must be
//         initialized by calling subroutine sinti(n,work) and a
//         different work array must be used for each different
//         value of n. this initialization does not have to be
//         repeated so long as n remains unchanged thus subsequent
//         transforms can be obtained faster than the first.
//
// ifac    an integer work array of length at least 15.
//
// output parameters
//
// x       for i=1,...,n
//           x(i)= the sum from k=1 to k=n
//             2*x(k)*sin(k*i*pi/(n+1))
//
// a call of sint followed by another call of
// sint will multiply the sequence x by 2*(n+1).
// hence sint is the unnormalized inverse
// of itself.
//
// work    contains initialization calculations which must not be
//         destroyed between calls of sint.
// ifac    contains initialization calculations which must not be
//         destroyed between calls of sint.
func sint(n int, x, work []float64, ifac []int) {
	if len(x) < n {
		panic("fourier: short sequence")
	}
	if len(work) < 5*(n+1)/2 {
		panic("fourier: short work")
	}
	if len(ifac) < 15 {
		panic("fourier: short ifac")
	}
	if n == 0 {
		return
	}
	sint1(n, x, work, work[n/2:], work[n/2+n+1:], ifac)
}

func sint1(n int, war, was, xh, x []float64, ifac []int) {
	const sqrt3 = 1.73205080756888

	for i := 0; i < n; i++ {
		xh[i] = war[i]
		war[i] = x[i]
	}

	switch n {
	case 1:
		xh[0] *= 2
	case 2:
		xh[0], xh[1] = sqrt3*(xh[0]+xh[1]), sqrt3*(xh[0]-xh[1])
	default:
		x[0] = 0
		for k := 0; k < n/2; k++ {
			kc := n - k
			t1 := xh[k] - xh[kc-1]
			t2 := was[k] * (xh[k] + xh[kc-1])
			x[k+1] = t1 + t2
			x[kc+1-1] = t2 - t1
		}
		if n%2 != 0 {
			x[n/2+1] = 4 * xh[n/2]
		}
		rfftf1(n+1, x, xh, war, ifac)
		xh[0] = 0.5 * x[0]
		for i := 2; i < n; i += 2 {
			xh[i-1] = -x[i]
			xh[i] = xh[i-2] + x[i-1]
		}
		if n%2 == 0 {
			xh[n-1] = -x[n]
		}
	}

	for i := 0; i < n; i++ {
		x[i] = war[i]
		war[i] = xh[i]
	}
}
