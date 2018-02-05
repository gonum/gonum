// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This is a translation of the FFTPACK rfft functions by
// Paul N Swarztrauber, placed in the public domain at
// http://www.netlib.org/fftpack/.

package fourier

import "math"

// rffti initializes the array work which is used in both rfftf
// and rfftb. the prime factorization of n together with a
// tabulation of the trigonometric functions are computed and
// stored in work.
//
//  input parameter
//
//  n      the length of the sequence to be transformed.
//
//  output parameters
//
//  work   a work array which must be dimensioned at least 2*n.
//         the same work array can be used for both rfftf and rfftb
//         as long as n remains unchanged. different work arrays
//         are required for different values of n. the contents of
//         work must not be changed between calls of rfftf or rfftb.
//
//  ifac   a work array containing the factors of n. ifac must have
//         length 15.
func rffti(n int, work []float64, ifac []int) {
	if n == 1 {
		return
	}
	rffti1(n, work[n:], ifac)
}

func rffti1(n int, wa []float64, ifac []int) {
	ntryh := [4]int{4, 2, 3, 5}

	nl := n
	nf := 0

outer:
	for j, ntry := 0, 0; ; j++ {
		if j < 4 {
			ntry = ntryh[j]
		} else {
			ntry += 2
		}
		for {
			if nl%ntry != 0 {
				continue outer
			}

			ifac[nf+2] = ntry
			nl /= ntry
			nf++

			if ntry == 2 && nf != 1 {
				for i := 1; i < nf; i++ {
					ib := nf - i + 1
					ifac[ib+1] = ifac[ib]
				}
				ifac[2] = 2
			}

			if nl == 1 {
				break outer
			}
		}
	}

	ifac[0] = n
	ifac[1] = nf
	if nf == 1 {
		return
	}
	argh := 2 * math.Pi / float64(n)

	is := 0
	l1 := 1
	for k1 := 0; k1 < nf-1; k1++ {
		ip := ifac[k1+2]
		ld := 0
		l2 := l1 * ip
		ido := n / l2
		for j := 0; j < ip-1; j++ {
			ld += l1
			i := is
			fi := 0.
			argld := float64(ld) * argh
			for ii := 2; ii < ido; ii += 2 {
				fi++
				arg := fi * argld
				wa[i] = math.Cos(arg)
				wa[i+1] = math.Sin(arg)
				i += 2
			}
			is += ido
		}
		l1 = l2
	}
}

// rfftf computes the fourier coefficients of a real
// perodic sequence (fourier analysis). the transform is defined
// below at output parameter r.
//
//  input parameters
//
//  n      the length of the array r to be transformed.  the method
//         is most efficient when n is a product of small primes.
//         n may change so long as different work arrays are provided
//
//  r      a real array of length n which contains the sequence
//         to be transformed
//
//  work   a work array which must be dimensioned at least 2*n.
//         in the program that calls rfftf. the work array must be
//         initialized by calling subroutine rffti(n,work) and a
//         different work array must be used for each different
//         value of n. this initialization does not have to be
//         repeated so long as n remains unchanged thus subsequent
//         transforms can be obtained faster than the first.
//         the same work array can be used by rfftf and rfftb.
//
//  ifac   a work array containing the factors of n. ifac must have
//         length 15.
//
//  output parameters
//
//  r      r(1) = the sum from i=1 to i=n of r(i)
//
//         if n is even set l=n/2, if n is odd set l = (n+1)/2
//           then for k = 2,...,l
//             r(2*k-2) = the sum from i = 1 to i = n of
//               r(i)*cos((k-1)*(i-1)*2*pi/n)
//             r(2*k-1) = the sum from i = 1 to i = n of
//               -r(i)*sin((k-1)*(i-1)*2*pi/n)
//
//         if n is even
//           r(n) = the sum from i = 1 to i = n of
//             (-1)**(i-1)*r(i)
//
//  This transform is unnormalized since a call of rfftf
//  followed by a call of rfftb will multiply the input
//  sequence by n.
//
//  work   contains results which must not be destroyed between
//         calls of rfftf or rfftb.
//  ifac   contains results which must not be destroyed between
//         calls of rfftf or rfftb.
func rfftf(n int, r, work []float64, ifac []int) {
	if n == 1 {
		return
	}
	rfftf1(n, r, work, work[n:], ifac)
}

func rfftf1(n int, c, ch []float64, wa oneArray, ifac []int) {
	nf := ifac[1]
	na := 1
	l2 := n
	iw := n

	for k1 := 1; k1 <= nf; k1++ {
		kh := nf - k1
		ip := ifac[kh+2]
		l1 := l2 / ip
		ido := n / l2
		idl1 := ido * l1
		iw -= (ip - 1) * ido
		na = 1 - na

		switch ip {
		case 4:
			ix2 := iw + ido
			ix3 := ix2 + ido
			if na == 0 {
				radf4(ido, l1, c, ch, wa.sliceFrom(iw), wa.sliceFrom(ix2), wa.sliceFrom(ix3))
			} else {
				radf4(ido, l1, ch, c, wa.sliceFrom(iw), wa.sliceFrom(ix2), wa.sliceFrom(ix3))
			}
		case 2:
			if na == 0 {
				radf2(ido, l1, c, ch, wa.sliceFrom(iw))
			} else {
				radf2(ido, l1, ch, c, wa.sliceFrom(iw))
			}
		case 3:
			ix2 := iw + ido
			if na == 0 {
				radf3(ido, l1, c, ch, wa.sliceFrom(iw), wa.sliceFrom(ix2))
			} else {
				radf3(ido, l1, ch, c, wa.sliceFrom(iw), wa.sliceFrom(ix2))
			}
		case 5:
			ix2 := iw + ido
			ix3 := ix2 + ido
			ix4 := ix3 + ido
			if na == 0 {
				radf5(ido, l1, c, ch, wa.sliceFrom(iw), wa.sliceFrom(ix2), wa.sliceFrom(ix3), wa.sliceFrom(ix4))
			} else {
				radf5(ido, l1, ch, c, wa.sliceFrom(iw), wa.sliceFrom(ix2), wa.sliceFrom(ix3), wa.sliceFrom(ix4))
			}
		default:
			if ido == 1 {
				na = 1 - na
			}
			if na == 0 {
				radfg(ido, ip, l1, idl1, c, c, c, ch, ch, wa.sliceFrom(iw))
				na = 1
			} else {
				radfg(ido, ip, l1, idl1, ch, ch, ch, c, c, wa.sliceFrom(iw))
				na = 0
			}
		}

		l2 = l1
	}

	if na == 1 {
		return
	}
	for i := 0; i < n; i++ {
		c[i] = ch[i]
	}
}

func radf2(ido, l1 int, cc, ch []float64, wa1 oneArray) {
	cc3 := newThreeArray(ido, l1, 2, cc)
	ch3 := newThreeArray(ido, 2, l1, ch)

	for k := 1; k <= l1; k++ {
		ch3.set(1, 1, k, cc3.at(1, k, 1)+cc3.at(1, k, 2))
		ch3.set(ido, 2, k, cc3.at(1, k, 1)-cc3.at(1, k, 2))
	}
	if ido < 2 {
		return
	}
	if ido > 2 {
		idp2 := ido + 2
		for k := 1; k <= l1; k++ {
			for i := 3; i <= ido; i += 2 {
				ic := idp2 - i
				tr2 := wa1.at(i-2)*cc3.at(i-1, k, 2) + wa1.at(i-1)*cc3.at(i, k, 2)
				ti2 := wa1.at(i-2)*cc3.at(i, k, 2) - wa1.at(i-1)*cc3.at(i-1, k, 2)
				ch3.set(i, 1, k, cc3.at(i, k, 1)+ti2)
				ch3.set(ic, 2, k, ti2-cc3.at(i, k, 1))
				ch3.set(i-1, 1, k, cc3.at(i-1, k, 1)+tr2)
				ch3.set(ic-1, 2, k, cc3.at(i-1, k, 1)-tr2)
			}
		}
		if ido%2 == 1 {
			return
		}
	}
	for k := 1; k <= l1; k++ {
		ch3.set(1, 2, k, -cc3.at(ido, k, 2))
		ch3.set(ido, 1, k, cc3.at(ido, k, 1))
	}
}

func radf3(ido, l1 int, cc, ch []float64, wa1, wa2 oneArray) {
	const (
		taur = -0.5
		taui = 0.866025403784439 // sqrt(3)/2
	)

	cc3 := newThreeArray(ido, l1, 3, cc)
	ch3 := newThreeArray(ido, 3, l1, ch)

	for k := 1; k <= l1; k++ {
		cr2 := cc3.at(1, k, 2) + cc3.at(1, k, 3)
		ch3.set(1, 1, k, cc3.at(1, k, 1)+cr2)
		ch3.set(1, 3, k, taui*(cc3.at(1, k, 3)-cc3.at(1, k, 2)))
		ch3.set(ido, 2, k, cc3.at(1, k, 1)+taur*cr2)
	}
	if ido < 2 {
		return
	}
	idp2 := ido + 2
	for k := 1; k <= l1; k++ {
		for i := 3; i <= ido; i += 2 {
			ic := idp2 - i
			dr2 := wa1.at(i-2)*cc3.at(i-1, k, 2) + wa1.at(i-1)*cc3.at(i, k, 2)
			di2 := wa1.at(i-2)*cc3.at(i, k, 2) - wa1.at(i-1)*cc3.at(i-1, k, 2)
			dr3 := wa2.at(i-2)*cc3.at(i-1, k, 3) + wa2.at(i-1)*cc3.at(i, k, 3)
			di3 := wa2.at(i-2)*cc3.at(i, k, 3) - wa2.at(i-1)*cc3.at(i-1, k, 3)
			cr2 := dr2 + dr3
			ci2 := di2 + di3
			ch3.set(i-1, 1, k, cc3.at(i-1, k, 1)+cr2)
			ch3.set(i, 1, k, cc3.at(i, k, 1)+ci2)
			tr2 := cc3.at(i-1, k, 1) + taur*cr2
			ti2 := cc3.at(i, k, 1) + taur*ci2
			tr3 := taui * (di2 - di3)
			ti3 := taui * (dr3 - dr2)
			ch3.set(i-1, 3, k, tr2+tr3)
			ch3.set(ic-1, 2, k, tr2-tr3)
			ch3.set(i, 3, k, ti2+ti3)
			ch3.set(ic, 2, k, ti3-ti2)
		}
	}
}

func radf4(ido, l1 int, cc, ch []float64, wa1, wa2, wa3 oneArray) {
	const hsqt2 = math.Sqrt2 / 2

	cc3 := newThreeArray(ido, l1, 4, cc)
	ch3 := newThreeArray(ido, 4, l1, ch)

	for k := 1; k <= l1; k++ {
		tr1 := cc3.at(1, k, 2) + cc3.at(1, k, 4)
		tr2 := cc3.at(1, k, 1) + cc3.at(1, k, 3)
		ch3.set(1, 1, k, tr1+tr2)
		ch3.set(ido, 4, k, tr2-tr1)
		ch3.set(ido, 2, k, cc3.at(1, k, 1)-cc3.at(1, k, 3))
		ch3.set(1, 3, k, cc3.at(1, k, 4)-cc3.at(1, k, 2))
	}
	if ido < 2 {
		return
	}
	if ido > 2 {
		idp2 := ido + 2
		for k := 1; k <= l1; k++ {
			for i := 3; i <= ido; i += 2 {
				ic := idp2 - i
				cr2 := wa1.at(i-2)*cc3.at(i-1, k, 2) + wa1.at(i-1)*cc3.at(i, k, 2)
				ci2 := wa1.at(i-2)*cc3.at(i, k, 2) - wa1.at(i-1)*cc3.at(i-1, k, 2)
				cr3 := wa2.at(i-2)*cc3.at(i-1, k, 3) + wa2.at(i-1)*cc3.at(i, k, 3)
				ci3 := wa2.at(i-2)*cc3.at(i, k, 3) - wa2.at(i-1)*cc3.at(i-1, k, 3)
				cr4 := wa3.at(i-2)*cc3.at(i-1, k, 4) + wa3.at(i-1)*cc3.at(i, k, 4)
				ci4 := wa3.at(i-2)*cc3.at(i, k, 4) - wa3.at(i-1)*cc3.at(i-1, k, 4)
				tr1 := cr2 + cr4
				tr4 := cr4 - cr2
				ti1 := ci2 + ci4
				ti4 := ci2 - ci4
				ti2 := cc3.at(i, k, 1) + ci3
				ti3 := cc3.at(i, k, 1) - ci3
				tr2 := cc3.at(i-1, k, 1) + cr3
				tr3 := cc3.at(i-1, k, 1) - cr3
				ch3.set(i-1, 1, k, tr1+tr2)
				ch3.set(ic-1, 4, k, tr2-tr1)
				ch3.set(i, 1, k, ti1+ti2)
				ch3.set(ic, 4, k, ti1-ti2)
				ch3.set(i-1, 3, k, ti4+tr3)
				ch3.set(ic-1, 2, k, tr3-ti4)
				ch3.set(i, 3, k, tr4+ti3)
				ch3.set(ic, 2, k, tr4-ti3)
			}
		}

		if ido%2 == 1 {
			return
		}
	}
	for k := 1; k <= l1; k++ {
		ti1 := -hsqt2 * (cc3.at(ido, k, 2) + cc3.at(ido, k, 4))
		tr1 := hsqt2 * (cc3.at(ido, k, 2) - cc3.at(ido, k, 4))
		ch3.set(ido, 1, k, tr1+cc3.at(ido, k, 1))
		ch3.set(ido, 3, k, cc3.at(ido, k, 1)-tr1)
		ch3.set(1, 2, k, ti1-cc3.at(ido, k, 3))
		ch3.set(1, 4, k, ti1+cc3.at(ido, k, 3))
	}
}

func radf5(ido, l1 int, cc, ch []float64, wa1, wa2, wa3, wa4 oneArray) {
	const (
		tr11 = 0.309016994374947
		ti11 = 0.951056516295154
		tr12 = -0.809016994374947
		ti12 = 0.587785252292473
	)

	cc3 := newThreeArray(ido, l1, 5, cc)
	ch3 := newThreeArray(ido, 5, l1, ch)

	for k := 1; k <= l1; k++ {
		cr2 := cc3.at(1, k, 5) + cc3.at(1, k, 2)
		ci5 := cc3.at(1, k, 5) - cc3.at(1, k, 2)
		cr3 := cc3.at(1, k, 4) + cc3.at(1, k, 3)
		ci4 := cc3.at(1, k, 4) - cc3.at(1, k, 3)
		ch3.set(1, 1, k, cc3.at(1, k, 1)+cr2+cr3)
		ch3.set(ido, 2, k, cc3.at(1, k, 1)+tr11*cr2+tr12*cr3)
		ch3.set(1, 3, k, ti11*ci5+ti12*ci4)
		ch3.set(ido, 4, k, cc3.at(1, k, 1)+tr12*cr2+tr11*cr3)
		ch3.set(1, 5, k, ti12*ci5-ti11*ci4)
	}

	if ido < 2 {
		return
	}
	idp2 := ido + 2
	for k := 1; k <= l1; k++ {
		for i := 3; i <= ido; i += 2 {
			ic := idp2 - i
			dr2 := wa1.at(i-2)*cc3.at(i-1, k, 2) + wa1.at(i-1)*cc3.at(i, k, 2)
			di2 := wa1.at(i-2)*cc3.at(i, k, 2) - wa1.at(i-1)*cc3.at(i-1, k, 2)
			dr3 := wa2.at(i-2)*cc3.at(i-1, k, 3) + wa2.at(i-1)*cc3.at(i, k, 3)
			di3 := wa2.at(i-2)*cc3.at(i, k, 3) - wa2.at(i-1)*cc3.at(i-1, k, 3)
			dr4 := wa3.at(i-2)*cc3.at(i-1, k, 4) + wa3.at(i-1)*cc3.at(i, k, 4)
			di4 := wa3.at(i-2)*cc3.at(i, k, 4) - wa3.at(i-1)*cc3.at(i-1, k, 4)
			dr5 := wa4.at(i-2)*cc3.at(i-1, k, 5) + wa4.at(i-1)*cc3.at(i, k, 5)
			di5 := wa4.at(i-2)*cc3.at(i, k, 5) - wa4.at(i-1)*cc3.at(i-1, k, 5)
			cr2 := dr2 + dr5
			ci5 := dr5 - dr2
			cr5 := di2 - di5
			ci2 := di2 + di5
			cr3 := dr3 + dr4
			ci4 := dr4 - dr3
			cr4 := di3 - di4
			ci3 := di3 + di4
			ch3.set(i-1, 1, k, cc3.at(i-1, k, 1)+cr2+cr3)
			ch3.set(i, 1, k, cc3.at(i, k, 1)+ci2+ci3)
			tr2 := cc3.at(i-1, k, 1) + tr11*cr2 + tr12*cr3
			ti2 := cc3.at(i, k, 1) + tr11*ci2 + tr12*ci3
			tr3 := cc3.at(i-1, k, 1) + tr12*cr2 + tr11*cr3
			ti3 := cc3.at(i, k, 1) + tr12*ci2 + tr11*ci3
			tr5 := ti11*cr5 + ti12*cr4
			ti5 := ti11*ci5 + ti12*ci4
			tr4 := ti12*cr5 - ti11*cr4
			ti4 := ti12*ci5 - ti11*ci4
			ch3.set(i-1, 3, k, tr2+tr5)
			ch3.set(ic-1, 2, k, tr2-tr5)
			ch3.set(i, 3, k, ti2+ti5)
			ch3.set(ic, 2, k, ti5-ti2)
			ch3.set(i-1, 5, k, tr3+tr4)
			ch3.set(ic-1, 4, k, tr3-tr4)
			ch3.set(i, 5, k, ti3+ti4)
			ch3.set(ic, 4, k, ti4-ti3)
		}
	}
}

func radfg(ido, ip, l1, idl1 int, cc, c1, c2, ch, ch2 []float64, wa oneArray) {
	cc3 := newThreeArray(ido, ip, l1, cc)
	c13 := newThreeArray(ido, l1, ip, c1)
	ch3 := newThreeArray(ido, l1, ip, ch)
	c2m := newTwoArray(idl1, ip, c2)
	ch2m := newTwoArray(idl1, ip, ch2)

	arg := 2 * math.Pi / float64(ip)
	dcp := math.Cos(arg)
	dsp := math.Sin(arg)
	ipph := (ip + 1) / 2
	ipp2 := ip + 2
	idp2 := ido + 2
	nbd := (ido - 1) / 2

	if ido == 1 {
		for ik := 1; ik <= idl1; ik++ {
			c2m.set(ik, 1, ch2m.at(ik, 1))
		}
	} else {
		for ik := 1; ik <= idl1; ik++ {
			ch2m.set(ik, 1, c2m.at(ik, 1))
		}
		for j := 2; j <= ip; j++ {
			for k := 1; k <= l1; k++ {
				ch3.set(1, k, j, c13.at(1, k, j))
			}
		}

		is := -ido
		if nbd > l1 {
			for j := 2; j <= ip; j++ {
				is += ido
				for k := 1; k <= l1; k++ {
					idij := is
					for i := 3; i <= ido; i += 2 {
						idij += 2
						ch3.set(i-1, k, j, wa.at(idij-1)*c13.at(i-1, k, j)+wa.at(idij)*c13.at(i, k, j))
						ch3.set(i, k, j, wa.at(idij-1)*c13.at(i, k, j)-wa.at(idij)*c13.at(i-1, k, j))
					}
				}
			}
		} else {
			for j := 2; j <= ip; j++ {
				is += ido
				idij := is
				for i := 3; i <= ido; i += 2 {
					idij += 2
					for k := 1; k <= l1; k++ {
						ch3.set(i-1, k, j, wa.at(idij-1)*c13.at(i-1, k, j)+wa.at(idij)*c13.at(i, k, j))
						ch3.set(i, k, j, wa.at(idij-1)*c13.at(i, k, j)-wa.at(idij)*c13.at(i-1, k, j))
					}
				}
			}
		}
		if nbd < l1 {
			for j := 2; j <= ipph; j++ {
				jc := ipp2 - j
				for i := 3; i <= ido; i += 2 {
					for k := 1; k <= l1; k++ {
						c13.set(i-1, k, j, ch3.at(i-1, k, j)+ch3.at(i-1, k, jc))
						c13.set(i-1, k, jc, ch3.at(i, k, j)-ch3.at(i, k, jc))
						c13.set(i, k, j, ch3.at(i, k, j)+ch3.at(i, k, jc))
						c13.set(i, k, jc, ch3.at(i-1, k, jc)-ch3.at(i-1, k, j))
					}
				}
			}
		} else {
			for j := 2; j <= ipph; j++ {
				jc := ipp2 - j
				for k := 1; k <= l1; k++ {
					for i := 3; i <= ido; i += 2 {
						c13.set(i-1, k, j, ch3.at(i-1, k, j)+ch3.at(i-1, k, jc))
						c13.set(i-1, k, jc, ch3.at(i, k, j)-ch3.at(i, k, jc))
						c13.set(i, k, j, ch3.at(i, k, j)+ch3.at(i, k, jc))
						c13.set(i, k, jc, ch3.at(i-1, k, jc)-ch3.at(i-1, k, j))
					}
				}
			}
		}
	}

	for j := 2; j <= ipph; j++ {
		jc := ipp2 - j
		for k := 1; k <= l1; k++ {
			c13.set(1, k, j, ch3.at(1, k, j)+ch3.at(1, k, jc))
			c13.set(1, k, jc, ch3.at(1, k, jc)-ch3.at(1, k, j))
		}
	}
	ar1 := 1.0
	ai1 := 0.0
	for l := 2; l <= ipph; l++ {
		lc := ipp2 - l
		ar1h := dcp*ar1 - dsp*ai1
		ai1 = dcp*ai1 + dsp*ar1
		ar1 = ar1h
		for ik := 1; ik <= idl1; ik++ {
			ch2m.set(ik, l, c2m.at(ik, 1)+ar1*c2m.at(ik, 2))
			ch2m.set(ik, lc, ai1*c2m.at(ik, ip))
		}
		dc2 := ar1
		ds2 := ai1
		ar2 := ar1
		ai2 := ai1
		for j := 3; j <= ipph; j++ {
			jc := ipp2 - j
			ar2h := dc2*ar2 - ds2*ai2
			ai2 = dc2*ai2 + ds2*ar2
			ar2 = ar2h
			for ik := 1; ik <= idl1; ik++ {
				ch2m.set(ik, l, ch2m.at(ik, l)+ar2*c2m.at(ik, j))
				ch2m.set(ik, lc, ch2m.at(ik, lc)+ai2*c2m.at(ik, jc))
			}
		}
	}
	for j := 2; j <= ipph; j++ {
		for ik := 1; ik <= idl1; ik++ {
			ch2m.set(ik, 1, ch2m.at(ik, 1)+c2m.at(ik, j))
		}
	}

	if ido < l1 {
		for i := 1; i <= ido; i++ {
			for k := 1; k <= l1; k++ {
				cc3.set(i, 1, k, ch3.at(i, k, 1))
			}
		}
	} else {
		for k := 1; k <= l1; k++ {
			for i := 1; i <= ido; i++ {
				cc3.set(i, 1, k, ch3.at(i, k, 1))
			}
		}
	}
	for j := 2; j <= ipph; j++ {
		jc := ipp2 - j
		j2 := 2 * j
		for k := 1; k <= l1; k++ {
			cc3.set(ido, j2-2, k, ch3.at(1, k, j))
			cc3.set(1, j2-1, k, ch3.at(1, k, jc))
		}
	}

	if ido == 1 {
		return
	}
	if nbd < l1 {
		for j := 2; j <= ipph; j++ {
			jc := ipp2 - j
			j2 := 2 * j
			for i := 3; i <= ido; i += 2 {
				ic := idp2 - i
				for k := 1; k <= l1; k++ {
					cc3.set(i-1, j2-1, k, ch3.at(i-1, k, j)+ch3.at(i-1, k, jc))
					cc3.set(ic-1, j2-2, k, ch3.at(i-1, k, j)-ch3.at(i-1, k, jc))
					cc3.set(i, j2-1, k, ch3.at(i, k, j)+ch3.at(i, k, jc))
					cc3.set(ic, j2-2, k, ch3.at(i, k, jc)-ch3.at(i, k, j))
				}
			}
		}
		return
	}
	for j := 2; j <= ipph; j++ {
		jc := ipp2 - j
		j2 := 2 * j
		for k := 1; k <= l1; k++ {
			for i := 3; i <= ido; i += 2 {
				ic := idp2 - i
				cc3.set(i-1, j2-1, k, ch3.at(i-1, k, j)+ch3.at(i-1, k, jc))
				cc3.set(ic-1, j2-2, k, ch3.at(i-1, k, j)-ch3.at(i-1, k, jc))
				cc3.set(i, j2-1, k, ch3.at(i, k, j)+ch3.at(i, k, jc))
				cc3.set(ic, j2-2, k, ch3.at(i, k, jc)-ch3.at(i, k, j))
			}
		}
	}
}

// rfftb computes the real perodic sequence from its fourier
// coefficients (fourier synthesis). the transform is defined
// below at output parameter r.
//
//  input parameters
//
//  n      the length of the array r to be transformed.  the method
//         is most efficient when n is a product of small primes.
//         n may change so long as different work arrays are provided
//
//  r      a real array of length n which contains the sequence
//         to be transformed
//
//  work   a work array which must be dimensioned at least 2*n.
//         in the program that calls rfftb. the work array must be
//         initialized by calling subroutine rffti(n,work) and a
//         different work array must be used for each different
//         value of n. this initialization does not have to be
//         repeated so long as n remains unchanged thus subsequent
//         transforms can be obtained faster than the first.
//         the same work array can be used by rfftf and rfftb.
//
//  ifac   a work array containing the factors of n. ifac must have
//         length 15.
//
//  output parameters
//
//  r      for n even and for i = 1,...,n
//           r(i) = r(1)+(-1)**(i-1)*r(n)
//             plus the sum from k=2 to k=n/2 of
//               2.*r(2*k-2)*cos((k-1)*(i-1)*2*pi/n)
//               -2.*r(2*k-1)*sin((k-1)*(i-1)*2*pi/n)
//
//         for n odd and for i = 1,...,n
//           r(i) = r(1) plus the sum from k=2 to k=(n+1)/2 of
//             2.*r(2*k-2)*cos((k-1)*(i-1)*2*pi/n)
//             -2.*r(2*k-1)*sin((k-1)*(i-1)*2*pi/n)
//
//  This transform is unnormalized since a call of rfftf
//  followed by a call of rfftb will multiply the input
//  sequence by n.
//
//  work   contains results which must not be destroyed between
//         calls of rfftf or rfftb.
//  ifac   contains results which must not be destroyed between
//         calls of rfftf or rfftb.
func rfftb(n int, r, work []float64, ifac []int) {
	if n == 1 {
		return
	}
	rfftb1(n, r, work, work[n:], ifac)
}

func rfftb1(n int, c, ch []float64, wa oneArray, ifac []int) {
	nf := ifac[1]
	na := 0
	l1 := 1
	iw := 1

	for k1 := 1; k1 <= nf; k1++ {
		ip := ifac[k1+1]
		l2 := ip * l1
		ido := n / l2
		idl1 := ido * l1

		switch ip {
		case 4:
			ix2 := iw + ido
			ix3 := ix2 + ido
			if na == 0 {
				radb4(ido, l1, c, ch, wa.sliceFrom(iw), wa.sliceFrom(ix2), wa.sliceFrom(ix3))
			} else {
				radb4(ido, l1, ch, c, wa.sliceFrom(iw), wa.sliceFrom(ix2), wa.sliceFrom(ix3))
			}
			na = 1 - na
		case 2:
			if na == 0 {
				radb2(ido, l1, c, ch, wa.sliceFrom(iw))
			} else {
				radb2(ido, l1, ch, c, wa.sliceFrom(iw))
			}
			na = 1 - na
		case 3:
			ix2 := iw + ido
			if na == 0 {
				radb3(ido, l1, c, ch, wa.sliceFrom(iw), wa.sliceFrom(ix2))
			} else {
				radb3(ido, l1, ch, c, wa.sliceFrom(iw), wa.sliceFrom(ix2))
			}
			na = 1 - na
		case 5:
			ix2 := iw + ido
			ix3 := ix2 + ido
			ix4 := ix3 + ido
			if na == 0 {
				radb5(ido, l1, c, ch, wa.sliceFrom(iw), wa.sliceFrom(ix2), wa.sliceFrom(ix3), wa.sliceFrom(ix4))
			} else {
				radb5(ido, l1, ch, c, wa.sliceFrom(iw), wa.sliceFrom(ix2), wa.sliceFrom(ix3), wa.sliceFrom(ix4))
			}
			na = 1 - na
		default:
			if na == 0 {
				radbg(ido, ip, l1, idl1, c, c, c, ch, ch, wa.sliceFrom(iw))
			} else {
				radbg(ido, ip, l1, idl1, ch, ch, ch, c, c, wa.sliceFrom(iw))
			}
			if ido == 1 {
				na = 1 - na
			}
		}

		l1 = l2
		iw += (ip - 1) * ido
	}

	if na == 0 {
		return
	}
	for i := 0; i < n; i++ {
		c[i] = ch[i]
	}
}

func radb2(ido, l1 int, cc, ch []float64, wa1 oneArray) {
	cc3 := newThreeArray(ido, 2, l1, cc)
	ch3 := newThreeArray(ido, l1, 2, ch)

	for k := 1; k <= l1; k++ {
		ch3.set(1, k, 1, cc3.at(1, 1, k)+cc3.at(ido, 2, k))
		ch3.set(1, k, 2, cc3.at(1, 1, k)-cc3.at(ido, 2, k))
	}

	if ido < 2 {
		return
	}
	if ido > 2 {
		idp2 := ido + 2
		for k := 1; k <= l1; k++ {
			for i := 3; i <= ido; i += 2 {
				ic := idp2 - i
				ch3.set(i-1, k, 1, cc3.at(i-1, 1, k)+cc3.at(ic-1, 2, k))
				tr2 := cc3.at(i-1, 1, k) - cc3.at(ic-1, 2, k)
				ch3.set(i, k, 1, cc3.at(i, 1, k)-cc3.at(ic, 2, k))
				ti2 := cc3.at(i, 1, k) + cc3.at(ic, 2, k)
				ch3.set(i-1, k, 2, wa1.at(i-2)*tr2-wa1.at(i-1)*ti2)
				ch3.set(i, k, 2, wa1.at(i-2)*ti2+wa1.at(i-1)*tr2)
			}
		}

		if ido%2 == 1 {
			return
		}
	}
	for k := 1; k <= l1; k++ {
		ch3.set(ido, k, 1, cc3.at(ido, 1, k)+cc3.at(ido, 1, k))
		ch3.set(ido, k, 2, -(cc3.at(1, 2, k) + cc3.at(1, 2, k)))
	}
}

func radb3(ido, l1 int, cc, ch []float64, wa1, wa2 oneArray) {
	const (
		taur = -0.5
		taui = 0.866025403784439 // sqrt(3)/2
	)

	cc3 := newThreeArray(ido, 3, l1, cc)
	ch3 := newThreeArray(ido, l1, 3, ch)

	for k := 1; k <= l1; k++ {
		tr2 := cc3.at(ido, 2, k) + cc3.at(ido, 2, k)
		cr2 := cc3.at(1, 1, k) + taur*tr2
		ch3.set(1, k, 1, cc3.at(1, 1, k)+tr2)
		ci3 := taui * (cc3.at(1, 3, k) + cc3.at(1, 3, k))
		ch3.set(1, k, 2, cr2-ci3)
		ch3.set(1, k, 3, cr2+ci3)
	}

	if ido == 1 {
		return
	}

	idp2 := ido + 2
	for k := 1; k <= l1; k++ {
		for i := 3; i <= ido; i += 2 {
			ic := idp2 - i
			tr2 := cc3.at(i-1, 3, k) + cc3.at(ic-1, 2, k)
			cr2 := cc3.at(i-1, 1, k) + taur*tr2
			ch3.set(i-1, k, 1, cc3.at(i-1, 1, k)+tr2)
			ti2 := cc3.at(i, 3, k) - cc3.at(ic, 2, k)
			ci2 := cc3.at(i, 1, k) + taur*ti2
			ch3.set(i, k, 1, cc3.at(i, 1, k)+ti2)
			cr3 := taui * (cc3.at(i-1, 3, k) - cc3.at(ic-1, 2, k))
			ci3 := taui * (cc3.at(i, 3, k) + cc3.at(ic, 2, k))
			dr2 := cr2 - ci3
			dr3 := cr2 + ci3
			di2 := ci2 + cr3
			di3 := ci2 - cr3
			ch3.set(i-1, k, 2, wa1.at(i-2)*dr2-wa1.at(i-1)*di2)
			ch3.set(i, k, 2, wa1.at(i-2)*di2+wa1.at(i-1)*dr2)
			ch3.set(i-1, k, 3, wa2.at(i-2)*dr3-wa2.at(i-1)*di3)
			ch3.set(i, k, 3, wa2.at(i-2)*di3+wa2.at(i-1)*dr3)
		}
	}
}

func radb4(ido, l1 int, cc, ch []float64, wa1, wa2, wa3 oneArray) {
	cc3 := newThreeArray(ido, 4, l1, cc)
	ch3 := newThreeArray(ido, l1, 4, ch)

	for k := 1; k <= l1; k++ {
		tr1 := cc3.at(1, 1, k) - cc3.at(ido, 4, k)
		tr2 := cc3.at(1, 1, k) + cc3.at(ido, 4, k)
		tr3 := cc3.at(ido, 2, k) + cc3.at(ido, 2, k)
		tr4 := cc3.at(1, 3, k) + cc3.at(1, 3, k)
		ch3.set(1, k, 1, tr2+tr3)
		ch3.set(1, k, 2, tr1-tr4)
		ch3.set(1, k, 3, tr2-tr3)
		ch3.set(1, k, 4, tr1+tr4)
	}

	if ido < 2 {
		return
	}
	if ido > 2 {
		idp2 := ido + 2
		for k := 1; k <= l1; k++ {
			for i := 3; i <= ido; i += 2 {
				ic := idp2 - i
				ti1 := cc3.at(i, 1, k) + cc3.at(ic, 4, k)
				ti2 := cc3.at(i, 1, k) - cc3.at(ic, 4, k)
				ti3 := cc3.at(i, 3, k) - cc3.at(ic, 2, k)
				tr4 := cc3.at(i, 3, k) + cc3.at(ic, 2, k)
				tr1 := cc3.at(i-1, 1, k) - cc3.at(ic-1, 4, k)
				tr2 := cc3.at(i-1, 1, k) + cc3.at(ic-1, 4, k)
				ti4 := cc3.at(i-1, 3, k) - cc3.at(ic-1, 2, k)
				tr3 := cc3.at(i-1, 3, k) + cc3.at(ic-1, 2, k)
				ch3.set(i-1, k, 1, tr2+tr3)
				cr3 := tr2 - tr3
				ch3.set(i, k, 1, ti2+ti3)
				ci3 := ti2 - ti3
				cr2 := tr1 - tr4
				cr4 := tr1 + tr4
				ci2 := ti1 + ti4
				ci4 := ti1 - ti4
				ch3.set(i-1, k, 2, wa1.at(i-2)*cr2-wa1.at(i-1)*ci2)
				ch3.set(i, k, 2, wa1.at(i-2)*ci2+wa1.at(i-1)*cr2)
				ch3.set(i-1, k, 3, wa2.at(i-2)*cr3-wa2.at(i-1)*ci3)
				ch3.set(i, k, 3, wa2.at(i-2)*ci3+wa2.at(i-1)*cr3)
				ch3.set(i-1, k, 4, wa3.at(i-2)*cr4-wa3.at(i-1)*ci4)
				ch3.set(i, k, 4, wa3.at(i-2)*ci4+wa3.at(i-1)*cr4)
			}
		}

		if ido%2 == 1 {
			return
		}
	}
	for k := 1; k <= l1; k++ {
		ti1 := cc3.at(1, 2, k) + cc3.at(1, 4, k)
		ti2 := cc3.at(1, 4, k) - cc3.at(1, 2, k)
		tr1 := cc3.at(ido, 1, k) - cc3.at(ido, 3, k)
		tr2 := cc3.at(ido, 1, k) + cc3.at(ido, 3, k)
		ch3.set(ido, k, 1, tr2+tr2)
		ch3.set(ido, k, 2, math.Sqrt2*(tr1-ti1))
		ch3.set(ido, k, 3, ti2+ti2)
		ch3.set(ido, k, 4, -math.Sqrt2*(tr1+ti1))
	}
}

func radb5(ido, l1 int, cc, ch []float64, wa1, wa2, wa3, wa4 oneArray) {
	const (
		tr11 = 0.309016994374947
		ti11 = 0.951056516295154
		tr12 = -0.809016994374947
		ti12 = 0.587785252292473
	)

	cc3 := newThreeArray(ido, 5, l1, cc)
	ch3 := newThreeArray(ido, l1, 5, ch)

	for k := 1; k <= l1; k++ {
		ti5 := cc3.at(1, 3, k) + cc3.at(1, 3, k)
		ti4 := cc3.at(1, 5, k) + cc3.at(1, 5, k)
		tr2 := cc3.at(ido, 2, k) + cc3.at(ido, 2, k)
		tr3 := cc3.at(ido, 4, k) + cc3.at(ido, 4, k)
		ch3.set(1, k, 1, cc3.at(1, 1, k)+tr2+tr3)
		cr2 := cc3.at(1, 1, k) + tr11*tr2 + tr12*tr3
		cr3 := cc3.at(1, 1, k) + tr12*tr2 + tr11*tr3
		ci5 := ti11*ti5 + ti12*ti4
		ci4 := ti12*ti5 - ti11*ti4
		ch3.set(1, k, 2, cr2-ci5)
		ch3.set(1, k, 3, cr3-ci4)
		ch3.set(1, k, 4, cr3+ci4)
		ch3.set(1, k, 5, cr2+ci5)
	}

	if ido == 1 {
		return
	}

	idp2 := ido + 2
	for k := 1; k <= l1; k++ {
		for i := 3; i <= ido; i += 2 {
			ic := idp2 - i
			ti5 := cc3.at(i, 3, k) + cc3.at(ic, 2, k)
			ti2 := cc3.at(i, 3, k) - cc3.at(ic, 2, k)
			ti4 := cc3.at(i, 5, k) + cc3.at(ic, 4, k)
			ti3 := cc3.at(i, 5, k) - cc3.at(ic, 4, k)
			tr5 := cc3.at(i-1, 3, k) - cc3.at(ic-1, 2, k)
			tr2 := cc3.at(i-1, 3, k) + cc3.at(ic-1, 2, k)
			tr4 := cc3.at(i-1, 5, k) - cc3.at(ic-1, 4, k)
			tr3 := cc3.at(i-1, 5, k) + cc3.at(ic-1, 4, k)
			ch3.set(i-1, k, 1, cc3.at(i-1, 1, k)+tr2+tr3)
			ch3.set(i, k, 1, cc3.at(i, 1, k)+ti2+ti3)
			cr2 := cc3.at(i-1, 1, k) + tr11*tr2 + tr12*tr3
			ci2 := cc3.at(i, 1, k) + tr11*ti2 + tr12*ti3
			cr3 := cc3.at(i-1, 1, k) + tr12*tr2 + tr11*tr3
			ci3 := cc3.at(i, 1, k) + tr12*ti2 + tr11*ti3
			cr5 := ti11*tr5 + ti12*tr4
			ci5 := ti11*ti5 + ti12*ti4
			cr4 := ti12*tr5 - ti11*tr4
			ci4 := ti12*ti5 - ti11*ti4
			dr3 := cr3 - ci4
			dr4 := cr3 + ci4
			di3 := ci3 + cr4
			di4 := ci3 - cr4
			dr5 := cr2 + ci5
			dr2 := cr2 - ci5
			di5 := ci2 - cr5
			di2 := ci2 + cr5
			ch3.set(i-1, k, 2, wa1.at(i-2)*dr2-wa1.at(i-1)*di2)
			ch3.set(i, k, 2, wa1.at(i-2)*di2+wa1.at(i-1)*dr2)
			ch3.set(i-1, k, 3, wa2.at(i-2)*dr3-wa2.at(i-1)*di3)
			ch3.set(i, k, 3, wa2.at(i-2)*di3+wa2.at(i-1)*dr3)
			ch3.set(i-1, k, 4, wa3.at(i-2)*dr4-wa3.at(i-1)*di4)
			ch3.set(i, k, 4, wa3.at(i-2)*di4+wa3.at(i-1)*dr4)
			ch3.set(i-1, k, 5, wa4.at(i-2)*dr5-wa4.at(i-1)*di5)
			ch3.set(i, k, 5, wa4.at(i-2)*di5+wa4.at(i-1)*dr5)
		}
	}
}

func radbg(ido, ip, l1, idl1 int, cc, c1, c2, ch, ch2 []float64, wa oneArray) {
	cc3 := newThreeArray(ido, ip, l1, cc)
	c13 := newThreeArray(ido, l1, ip, c1)
	ch3 := newThreeArray(ido, l1, ip, ch)
	c2m := newTwoArray(idl1, ip, c2)
	ch2m := newTwoArray(idl1, ip, ch2)

	arg := 2 * math.Pi / float64(ip)
	dcp := math.Cos(arg)
	dsp := math.Sin(arg)
	ipph := (ip + 1) / 2
	ipp2 := ip + 2
	idp2 := ido + 2
	nbd := (ido - 1) / 2

	if ido < l1 {
		for i := 1; i <= ido; i++ {
			for k := 1; k <= l1; k++ {
				ch3.set(i, k, 1, cc3.at(i, 1, k))
			}
		}
	} else {
		for k := 1; k <= l1; k++ {
			for i := 1; i <= ido; i++ {
				ch3.set(i, k, 1, cc3.at(i, 1, k))
			}
		}
	}

	for j := 2; j <= ipph; j++ {
		jc := ipp2 - j
		j2 := 2 * j
		for k := 1; k <= l1; k++ {
			ch3.set(1, k, j, cc3.at(ido, j2-2, k)+cc3.at(ido, j2-2, k))
			ch3.set(1, k, jc, cc3.at(1, j2-1, k)+cc3.at(1, j2-1, k))
		}
	}

	if ido != 1 {
		if nbd < l1 {
			for j := 2; j <= ipph; j++ {
				jc := ipp2 - j
				for i := 3; i <= ido; i += 2 {
					ic := idp2 - i
					for k := 1; k <= l1; k++ {
						ch3.set(i-1, k, j, cc3.at(i-1, 2*j-1, k)+cc3.at(ic-1, 2*j-2, k))
						ch3.set(i-1, k, jc, cc3.at(i-1, 2*j-1, k)-cc3.at(ic-1, 2*j-2, k))
						ch3.set(i, k, j, cc3.at(i, 2*j-1, k)-cc3.at(ic, 2*j-2, k))
						ch3.set(i, k, jc, cc3.at(i, 2*j-1, k)+cc3.at(ic, 2*j-2, k))
					}
				}
			}
		} else {
			for j := 2; j <= ipph; j++ {
				jc := ipp2 - j
				for k := 1; k <= l1; k++ {
					for i := 3; i <= ido; i += 2 {
						ic := idp2 - i
						ch3.set(i-1, k, j, cc3.at(i-1, 2*j-1, k)+cc3.at(ic-1, 2*j-2, k))
						ch3.set(i-1, k, jc, cc3.at(i-1, 2*j-1, k)-cc3.at(ic-1, 2*j-2, k))
						ch3.set(i, k, j, cc3.at(i, 2*j-1, k)-cc3.at(ic, 2*j-2, k))
						ch3.set(i, k, jc, cc3.at(i, 2*j-1, k)+cc3.at(ic, 2*j-2, k))
					}
				}
			}
		}
	}

	ar1 := 1.0
	ai1 := 0.0
	for l := 2; l <= ipph; l++ {
		lc := ipp2 - l
		ar1h := dcp*ar1 - dsp*ai1
		ai1 = dcp*ai1 + dsp*ar1
		ar1 = ar1h
		for ik := 1; ik <= idl1; ik++ {
			c2m.set(ik, l, ch2m.at(ik, 1)+ar1*ch2m.at(ik, 2))
			c2m.set(ik, lc, ai1*ch2m.at(ik, ip))
		}
		dc2 := ar1
		ds2 := ai1
		ar2 := ar1
		ai2 := ai1
		for j := 3; j <= ipph; j++ {
			jc := ipp2 - j
			ar2h := dc2*ar2 - ds2*ai2
			ai2 = dc2*ai2 + ds2*ar2
			ar2 = ar2h
			for ik := 1; ik <= idl1; ik++ {
				c2m.set(ik, l, c2m.at(ik, l)+ar2*ch2m.at(ik, j))
				c2m.set(ik, lc, c2m.at(ik, lc)+ai2*ch2m.at(ik, jc))
			}
		}
	}

	for j := 2; j <= ipph; j++ {
		for ik := 1; ik <= idl1; ik++ {
			ch2m.set(ik, 1, ch2m.at(ik, 1)+ch2m.at(ik, j))
		}
	}
	for j := 2; j <= ipph; j++ {
		jc := ipp2 - j
		for k := 1; k <= l1; k++ {
			ch3.set(1, k, j, c13.at(1, k, j)-c13.at(1, k, jc))
			ch3.set(1, k, jc, c13.at(1, k, j)+c13.at(1, k, jc))
		}
	}

	if ido != 1 {
		if nbd < l1 {
			for j := 2; j <= ipph; j++ {
				jc := ipp2 - j
				for i := 3; i <= ido; i += 2 {
					for k := 1; k <= l1; k++ {
						ch3.set(i-1, k, j, c13.at(i-1, k, j)-c13.at(i, k, jc))
						ch3.set(i-1, k, jc, c13.at(i-1, k, j)+c13.at(i, k, jc))
						ch3.set(i, k, j, c13.at(i, k, j)+c13.at(i-1, k, jc))
						ch3.set(i, k, jc, c13.at(i, k, j)-c13.at(i-1, k, jc))
					}
				}
			}
		} else {
			for j := 2; j <= ipph; j++ {
				jc := ipp2 - j
				for k := 1; k <= l1; k++ {
					for i := 3; i <= ido; i += 2 {
						ch3.set(i-1, k, j, c13.at(i-1, k, j)-c13.at(i, k, jc))
						ch3.set(i-1, k, jc, c13.at(i-1, k, j)+c13.at(i, k, jc))
						ch3.set(i, k, j, c13.at(i, k, j)+c13.at(i-1, k, jc))
						ch3.set(i, k, jc, c13.at(i, k, j)-c13.at(i-1, k, jc))
					}
				}
			}
		}
	}

	if ido == 1 {
		return
	}
	for ik := 1; ik <= idl1; ik++ {
		c2m.set(ik, 1, ch2m.at(ik, 1))
	}
	for j := 2; j <= ip; j++ {
		for k := 1; k <= l1; k++ {
			c13.set(1, k, j, ch3.at(1, k, j))
		}
	}

	is := -ido
	if nbd > l1 {
		for j := 2; j <= ip; j++ {
			is += ido
			for k := 1; k <= l1; k++ {
				idij := is
				for i := 3; i <= ido; i += 2 {
					idij = idij + 2
					c13.set(i-1, k, j, wa.at(idij-1)*ch3.at(i-1, k, j)-wa.at(idij)*ch3.at(i, k, j))
					c13.set(i, k, j, wa.at(idij-1)*ch3.at(i, k, j)+wa.at(idij)*ch3.at(i-1, k, j))
				}
			}
		}
		return
	}
	for j := 2; j <= ip; j++ {
		is += ido
		idij := is
		for i := 3; i <= ido; i += 2 {
			idij += 2
			for k := 1; k <= l1; k++ {
				c13.set(i-1, k, j, wa.at(idij-1)*ch3.at(i-1, k, j)-wa.at(idij)*ch3.at(i, k, j))
				c13.set(i, k, j, wa.at(idij-1)*ch3.at(i, k, j)+wa.at(idij)*ch3.at(i-1, k, j))
			}
		}
	}
}
