// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is a translation of the FFTPACK sinq functions by
// Paul N Swarztrauber, placed in the public domain at
// http://www.netlib.org/fftpack/.

package fourier

import "math"

// subroutine sinqi initializes the array work which is used in
// both sinqf and sinqb. the prime factorization of n together with
// a tabulation of the trigonometric functions are computed and
// stored in work.
//
// input parameter
//
// n       the length of the sequence to be transformed. the method
//         is most efficient when n+1 is a product of small primes.
//
// output parameter
//
// work    a work array which must be dimensioned at least 3*n.
//         the same work array can be used for both sinqf and sinqb
//         as long as n remains unchanged. different work arrays
//         are required for different values of n. the contents of
//         work must not be changed between calls of sinqf or sinqb.
//
// ifac    an integer work array of length at least 15.
func sinqi(n int, work []float64, ifac []int) {
	if len(work) < 3*n {
		panic("fourier: short work")
	}
	if len(ifac) < 15 {
		panic("fourier: short ifac")
	}
	dt := 0.5 * math.Pi / float64(n)
	for k := range work[:n] {
		work[k] = math.Cos(float64(k+1) * dt)
	}
	rffti(n, work[n:], ifac)
}

// subroutine sinqf computes the fast fourier transform of quarter
// wave data. that is, sinqf computes the coefficients in a sine
// series representation with only odd wave numbers. the transform
// is defined below at output parameter x.
//
// sinqb is the unnormalized inverse of sinqf since a call of sinqf
// followed by a call of sinqb will multiply the input sequence x
// by 4*n.
//
// the array work which is used by subroutine sinqf must be
// initialized by calling subroutine sinqi(n,work).
//
// input parameters
//
// n       the length of the array x to be transformed.  the method
//         is most efficient when n is a product of small primes.
//
// x       an array which contains the sequence to be transformed
//
// work    a work array which must be dimensioned at least 3*n.
//         in the program that calls sinqf. the work array must be
//         initialized by calling subroutine sinqi(n,work) and a
//         different work array must be used for each different
//         value of n. this initialization does not have to be
//         repeated so long as n remains unchanged thus subsequent
//         transforms can be obtained faster than the first.
//
// ifac    an integer work array of length at least 15.
//
// output parameters
//
// x       for i=0, ..., n-1
//           x[i] = (-1)^(i)*x[n-1]
//             + the sum from k=0 to k=n-2 of
//               2*x[k]*sin((2*i+1)*k*pi/(2*n))
//
//         a call of sinqf followed by a call of
//         sinqb will multiply the sequence x by 4*n.
//         therefore sinqb is the unnormalized inverse
//         of sinqf.
//
// work    contains initialization calculations which must not
//         be destroyed between calls of sinqf or sinqb.
func sinqf(n int, x, work []float64, ifac []int) {
	if len(x) < n {
		panic("fourier: short sequence")
	}
	if len(work) < 3*n {
		panic("fourier: short work")
	}
	if len(ifac) < 15 {
		panic("fourier: short ifac")
	}
	if n == 1 {
		return
	}
	for k := 0; k < n/2; k++ {
		kc := n - k - 1
		x[k], x[kc] = x[kc], x[k]
	}
	cosqf(n, x, work, ifac)
	for k := 1; k < n; k += 2 {
		x[k] = -x[k]
	}
}

// subroutine sinqb computes the fast fourier transform of quarter
// wave data. that is, sinqb computes a sequence from its
// representation in terms of a sine series with odd wave numbers.
// the transform is defined below at output parameter x.
//
// sinqf is the unnormalized inverse of sinqb since a call of sinqb
// followed by a call of sinqf will multiply the input sequence x
// by 4*n.
//
// the array work which is used by subroutine sinqb must be
// initialized by calling subroutine sinqi(n,work).
//
//
// input parameters
//
// n       the length of the array x to be transformed.  the method
//         is most efficient when n is a product of small primes.
//
// x       an array which contains the sequence to be transformed
//
// work    a work array which must be dimensioned at least 3*n.
//         in the program that calls sinqb. the work array must be
//         initialized by calling subroutine sinqi(n,work) and a
//         different work array must be used for each different
//         value of n. this initialization does not have to be
//         repeated so long as n remains unchanged thus subsequent
//         transforms can be obtained faster than the first.
//
// ifac    an integer work array of length at least 15.
//
// output parameters
//
// x       for i=0, ..., n-1
//           x[i]= the sum from k=0 to k=n-1 of
//             4*x[k]*sin((2*k+1)*i*pi/(2*n))
//
//         a call of sinqb followed by a call of
//         sinqf will multiply the sequence x by 4*n.
//         therefore sinqf is the unnormalized inverse
//         of sinqb.
//
// work    contains initialization calculations which must not
//         be destroyed between calls of sinqb or sinqf.
func sinqb(n int, x, work []float64, ifac []int) {
	if len(x) < n {
		panic("fourier: short sequence")
	}
	if len(work) < 3*n {
		panic("fourier: short work")
	}
	if len(ifac) < 15 {
		panic("fourier: short ifac")
	}
	switch n {
	case 1:
		x[0] *= 4
		fallthrough
	case 0:
		return
	default:
		for k := 1; k < n; k += 2 {
			x[k] = -x[k]
		}
		cosqb(n, x, work, ifac)
		for k := 0; k < n/2; k++ {
			kc := n - k - 1
			x[k], x[kc] = x[kc], x[k]
		}
	}
}
