// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is a translation of the FFTPACK cosq functions by
// Paul N Swarztrauber, placed in the public domain at
// http://www.netlib.org/fftpack/.

package fourier

import "math"

// subroutine cosqi initializes the array work which is used in
// both cosqf and cosqb. the prime factorization of n together with
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
//         the same work array can be used for both cosqf and cosqb
//         as long as n remains unchanged. different work arrays
//         are required for different values of n. the contents of
//         work must not be changed between calls of cosqf or cosqb.
//
// ifac    an integer work array of length at least 15.
func cosqi(n int, work []float64, ifac []int) {
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

// subroutine cosqf computes the fast fourier transform of quarter
// wave data. that is, cosqf computes the coefficients in a cosine
// series representation with only odd wave numbers. the transform
// is defined below at output parameter x.
//
// cosqb is the unnormalized inverse of cosqf since a call of cosqf
// followed by a call of cosqb will multiply the input sequence x
// by 4*n.
//
// the array work which is used by subroutine cosqf must be
// initialized by calling subroutine cosqi(n,work).
//
// input parameters
//
// n       the length of the array x to be transformed.  the method
//         is most efficient when n is a product of small primes.
//
// x       an array which contains the sequence to be transformed
//
// work    a work array which must be dimensioned at least 3*n.
//         in the program that calls cosqf. the work array must be
//         initialized by calling subroutine cosqi(n,work) and a
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
//           x[i] = x[i] + the sum from k=0 to k=n-2 of
//               2*x[k]*cos((2*i+1)*k*pi/(2*n))
//
//         a call of cosqf followed by a call of
//         cosqb will multiply the sequence x by 4*n.
//         therefore cosqb is the unnormalized inverse
//         of cosqf.
//
// work    contains initialization calculations which must not
//         be destroyed between calls of cosqf or cosqb.
func cosqf(n int, x, work []float64, ifac []int) {
	if len(x) < n {
		panic("fourier: short sequence")
	}
	if len(work) < 3*n {
		panic("fourier: short work")
	}
	if len(ifac) < 15 {
		panic("fourier: short ifac")
	}
	if n < 2 {
		return
	}
	if n == 2 {
		tsqx := math.Sqrt2 * x[1]
		x[1] = x[0] - tsqx
		x[0] += tsqx
		return
	}
	cosqf1(n, x, work, work[n:], ifac)
}

func cosqf1(n int, x, w, xh []float64, ifac []int) {
	for k := 1; k < (n+1)/2; k++ {
		kc := n - k
		xh[k] = x[k] + x[kc]
		xh[kc] = x[k] - x[kc]
	}
	if n%2 == 0 {
		xh[(n+1)/2] = 2 * x[(n+1)/2]
	}
	for k := 1; k < (n+1)/2; k++ {
		kc := n - k
		x[k] = w[k-1]*xh[kc] + w[kc-1]*xh[k]
		x[kc] = w[k-1]*xh[k] - w[kc-1]*xh[kc]
	}
	if n%2 == 0 {
		x[(n+1)/2] = w[(n-1)/2] * xh[(n+1)/2]
	}
	rfftf(n, x, xh, ifac)
	for i := 2; i < n; i += 2 {
		x[i-1], x[i] = x[i-1]-x[i], x[i-1]+x[i]
	}
}

// subroutine cosqb computes the fast fourier transform of quarter
// wave data. that is, cosqb computes a sequence from its
// representation in terms of a cosine series with odd wave numbers.
// the transform is defined below at output parameter x.
//
// cosqf is the unnormalized inverse of cosqb since a call of cosqb
// followed by a call of cosqf will multiply the input sequence x
// by 4*n.
//
// the array work which is used by subroutine cosqb must be
// initialized by calling subroutine cosqi(n,work).
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
//         in the program that calls cosqb. the work array must be
//         initialized by calling subroutine cosqi(n,work) and a
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
//             4*x[k]*cos((2*k+1)*i*pi/(2*n))
//
//         a call of cosqb followed by a call of
//         cosqf will multiply the sequence x by 4*n.
//         therefore cosqf is the unnormalized inverse
//         of cosqb.
//
// work    contains initialization calculations which must not
//         be destroyed between calls of cosqb or cosqf.
func cosqb(n int, x, work []float64, ifac []int) {
	if len(x) < n {
		panic("fourier: short sequence")
	}
	if len(work) < 3*n {
		panic("fourier: short work")
	}
	if len(ifac) < 15 {
		panic("fourier: short ifac")
	}

	if n < 2 {
		x[0] *= 4
		return
	}
	if n == 2 {
		x[0], x[1] = 4*(x[0]+x[1]), 2*math.Sqrt2*(x[0]-x[1])
		return
	}
	cosqb1(n, x, work, work[n:], ifac)
}

func cosqb1(n int, x, w, xh []float64, ifac []int) {
	for i := 2; i < n; i += 2 {
		x[i-1], x[i] = x[i-1]+x[i], x[i]-x[i-1]
	}
	x[0] *= 2
	if n%2 == 0 {
		x[n-1] *= 2
	}
	rfftb(n, x, xh, ifac)
	for k := 1; k < (n+1)/2; k++ {
		kc := n - k
		xh[k] = w[k-1]*x[kc] + w[kc-1]*x[k]
		xh[kc] = w[k-1]*x[k] - w[kc-1]*x[kc]
	}
	if n%2 == 0 {
		x[(n+1)/2] *= 2 * w[(n-1)/2]
	}
	for k := 1; k < (n+1)/2; k++ {
		x[k] = xh[k] + xh[n-k]
		x[n-k] = xh[k] - xh[n-k]
	}
	x[0] *= 2
}
