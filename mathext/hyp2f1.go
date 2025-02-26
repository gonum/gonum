package mathext

import (
	"fmt"
	"math"
)

// Machine epsilon used in Cephes.
var mACHEP float64 = (math.Nextafter(1.0, 2.0) - 1.0) / 2

const eps float64 = 1e-13
const ethresh float64 = 1e-12

// Hyp2f1 returns the value of the Gaussian Hypergeometric function at z.
// For |z| < 1, this implementation follows the Cephes library.
// For |z| > 1, this implementation perform analytic continuation via relevant Hypergeometric identities.
// See https://en.wikipedia.org/wiki/Hypergeometric_function for more details.
func Hyp2f1(a float64, b float64, c float64, z float64) (float64, error) {
	if math.Abs(z) < 1 {
		return hyp2f1(a, b, c, z)
	}

	// Function undefined between the 1 and Inf branch points.
	if z > 0 {
		return math.NaN(), fmt.Errorf("analytic continuation failed")
	}

	// Analytic continuation formula contains infinities from Gamma(a-b) when a == b.
	// Make a and b different using equation 15.3.4 from Abramowitz.
	// M. Abramowitz and I. A. Stegun 1965. Handbook of Mathematical Functions, New York: Dover.
	if a == b {
		y, err := Hyp2f1(a, c-b, c, z/(z-1))
		if err != nil {
			return y, err
		}
		return math.Pow(1-z, -a) * y, nil
	}

	// Analytic continuation based on https://www.johndcook.com/blog/2021/11/03/escaping-the-unit-disk/
	y1, err := hyp2f1(a, 1-c+a, 1-b+a, 1/z)
	if err != nil {
		return y1, err
	}
	y1 *= math.Gamma(c) / math.Gamma(b) * math.Gamma(b-a) / math.Gamma(c-a) * math.Pow(-z, -a)
	y2, err := hyp2f1(b, 1-c+b, 1-a+b, 1/z)
	if err != nil {
		return y2, err
	}
	y2 *= math.Gamma(c) / math.Gamma(a) * math.Gamma(a-b) / math.Gamma(c-b) * math.Pow(-z, -b)
	return y1 + y2, nil
}

/*							hyp2f1.c
 *
 *	Gauss hypergeometric function   F
 *	                               2 1
 *
 *
 * SYNOPSIS:
 *
 * double a, b, c, x, y, hyp2f1();
 *
 * y = hyp2f1( a, b, c, x );
 *
 *
 * DESCRIPTION:
 *
 *
 *  hyp2f1( a, b, c, x )  =   F ( a, b; c; x )
 *                           2 1
 *
 *           inf.
 *            -   a(a+1)...(a+k) b(b+1)...(b+k)   k+1
 *   =  1 +   >   -----------------------------  x   .
 *            -         c(c+1)...(c+k) (k+1)!
 *          k = 0
 *
 *  Cases addressed are
 *	Tests and escapes for negative integer a, b, or c
 *	Linear transformation if c - a or c - b negative integer
 *	Special case c = a or c = b
 *	Linear transformation for  x near +1
 *	Transformation for x < -0.5
 *	Psi function expansion if x > 0.5 and c - a - b integer
 *      Conditionally, a recurrence on c to make c-a-b > 0
 *
 * |x| > 1 is rejected.
 *
 * The parameters a, b, c are considered to be integer
 * valued if they are within 1.0e-14 of the nearest integer
 * (1.0e-13 for IEEE arithmetic).
 *
 * ACCURACY:
 *
 *
 *               Relative error (-1 < x < 1):
 * arithmetic   domain     # trials      peak         rms
 *    IEEE      -1,7        230000      1.2e-11     5.2e-14
 *
 * Several special cases also tested with a, b, c in
 * the range -7 to 7.
 *
 * ERROR MESSAGES:
 *
 * A "partial loss of precision" message is printed if
 * the internally estimated relative error exceeds 1^-12.
 * A "singularity" message is printed on overflow or
 * in cases not addressed (such as x < -1).
 */
/*
Cephes Math Library Release 2.8:  June, 2000
Copyright 1984, 1987, 1992, 2000 by Stephen L. Moshier
*/
func hyp2f1(a float64, b float64, c float64, x float64) (float64, error) {
	var d float64
	var d1 float64
	var d2 float64
	var e float64
	var p float64
	var q float64
	var r float64
	var s float64
	var y float64
	var ax float64
	var ia float64
	var ib float64
	var ic float64
	var id float64
	var err float64
	var flag int
	var i int
	var aid int

	err = 0
	ax = math.Abs(x)
	s = 1 - x
	flag = 0
	ia = math.Round(a) // nearest integer to a
	ib = math.Round(b)

	if a <= 0 {
		if math.Abs(a-ia) < eps { // a is a negative integer
			flag |= 1
		}
	}

	if b <= 0 {
		if math.Abs(b-ib) < eps { // b is a negative integer
			flag |= 2
		}
	}

	if ax < 1 {
		if math.Abs(b-c) < eps { // b == c
			y = math.Pow(s, -a) // s to the -a power
			goto hypdon
		}
		if math.Abs(a-c) < eps { // a = c
			y = math.Pow(s, -b) // s to the -b power
			goto hypdon
		}
	}

	if c <= 0 {
		ic = math.Round(c)        // nearest integer to c
		if math.Abs(c-ic) < eps { // c is a negative integer
			// check if termination before explosion
			if flag&1 != 0 && ia > ic {
				goto hypok
			}
			if flag&2 != 0 && ib > ic {
				goto hypok
			}
			goto hypdiv
		}
	}

	if flag != 0 { // function is a polynomial
		goto hypok
	}

	if ax > 1 { // series diverges
		goto hypdiv
	}

	p = c - a
	ia = math.Round(p)                   // nearest integer to c-a
	if ia <= 0 && math.Abs(p-ia) < eps { // negative int c - a
		flag |= 4
	}

	r = c - b
	ib = math.Round(r)                   // nearest integer to c-b
	if ib <= 0 && math.Abs(r-ib) < eps { // negative int c - b
		flag |= 8
	}

	d = c - a - b
	id = math.Round(d) // nearest integer to d
	q = math.Abs(d - id)

	// Thanks to Christian Burger <BURGER@DMRHRZ11.HRZ.Uni-Marburg.DE> for reporting a bug here.
	if math.Abs(ax-1) < eps { // |x| == 1.0
		if x > 0 {
			if flag&12 != 0 { // negative int c-a or c-b
				if d >= 0 {
					goto hypf
				} else {
					goto hypdiv
				}
			}
			if d <= 0 {
				goto hypdiv
			}
			y = math.Gamma(c) * math.Gamma(d) / (math.Gamma(p) * math.Gamma(r))
			goto hypdon
		}
		if d <= -1 {
			goto hypdiv
		}
	}

	// Conditionally make d > 0 by recurrence on c, AMS55 #15.2.27.
	if d < 0 {
		// Try the power series first
		y = hyt2f1(a, b, c, x, &err)
		if err < ethresh {
			goto hypdon
		}

		// Apply the recurrence if power series fails
		err = 0
		aid = int(2 - id)
		e = c + float64(aid)
		var herr error
		d2, herr = Hyp2f1(a, b, e, x)
		if herr != nil {
			return math.NaN(), herr
		}
		d1, herr = Hyp2f1(a, b, e+1, x)
		if herr != nil {
			return math.NaN(), herr
		}
		q = a + b + 1
		for i = 0; i < aid; i++ {
			r = e - 1
			y = (e*(r-(2*e-q)*x)*d2 + (e-a)*(e-b)*x*d1) / (e * r * s)
			e = r
			d1 = d2
			d2 = y
		}
		goto hypdon
	}

	if flag&12 != 0 { // negative integer c-a or c-b
		goto hypf
	}

hypok:
	y = hyt2f1(a, b, c, x, &err)
hypdon:
	if err > ethresh {
		return y, fmt.Errorf("partial loss of precision")
	}
	return y, nil
hypf:
	y = math.Pow(s, d) * hys2f1(c-a, c-b, c, x, &err)
	goto hypdon
hypdiv:
	return math.MaxFloat64, fmt.Errorf("overflow range error")
}

/* Apply transformations for |x| near 1
 * then call the power series
 */
func hyt2f1(a float64, b float64, c float64, x float64, loss *float64) float64 {
	var p float64
	var q float64
	var r float64
	var s float64
	var t float64
	var y float64
	var d float64
	var err float64
	var err1 float64
	var ax float64
	var id float64
	var d1 float64
	var d2 float64
	var e float64
	var y1 float64
	var i int32
	var aid int32

	err = 0
	s = 1 - x
	if x < -0.5 {
		if b > a {
			y = math.Pow(s, -a) * hys2f1(a, c-b, c, -x/s, &err)
		} else {
			y = math.Pow(s, -b) * hys2f1(c-a, b, c, -x/s, &err)
		}
		goto done
	}

	d = c - a - b
	id = math.Round(d) // nearest integer to d

	if x > 0.9 {
		if math.Abs(d-id) > eps { // test for integer c-a-b
			// Try the power series first
			y = hys2f1(a, b, c, x, &err)
			if err < ethresh {
				goto done
			}

			// If power series fails, then apply AMS55 #15.3.6
			q = hys2f1(a, b, 1-d, s, &err)
			q *= math.Gamma(d) / (math.Gamma(c-a) * math.Gamma(c-b))
			r = math.Pow(s, d) * hys2f1(c-a, c-b, d+1, s, &err1)
			r *= math.Gamma(-d) / (math.Gamma(a) * math.Gamma(b))
			y = q + r

			q = math.Abs(q) // estimate cancellation error
			r = math.Abs(r)
			if q > r {
				r = q
			}
			err += err1 + mACHEP*r/y

			y *= math.Gamma(c)
			goto done
		} else {
			// Psi function expansion, AMS55 #15.3.10, #15.3.11, #15.3.12
			if id >= 0 {
				e = d
				d1 = d
				d2 = 0
				aid = int32(id)
			} else {
				e = -d
				d1 = 0
				d2 = d
				aid = int32(-id)
			}
			ax = math.Log(s)

			// sum for t = 0
			y = Digamma(1) + Digamma(1+e) - Digamma(a+d1) - Digamma(b+d1) - ax
			y /= math.Gamma(e + 1)

			p = (a + d1) * (b + d1) * s / math.Gamma(e+2)
			t = 1
			for {
				r = Digamma(1+t) + Digamma(1+t+e) - Digamma(a+t+d1) - Digamma(b+t+d1) - ax
				q = p * r
				y += q
				p *= s * (a + t + d1) / (t + 1)
				p *= (b + t + d1) / (t + 1 + e)
				t += 1

				if math.Abs(q/y) <= eps {
					break
				}
			}

			if id == 0 {
				y *= math.Gamma(c) / (math.Gamma(a) * math.Gamma(b))
				goto psidon
			}

			y1 = 1

			if aid == int32(1) {
				goto nosum
			}

			t = 0
			p = 1
			for i = int32(1); i < aid; i++ {
				r = 1 - e + t
				p *= s * (a + t + d2) * (b + t + d2) / r
				t += 1
				p /= t
				y1 += p
			}
		nosum:
			p = math.Gamma(c)
			y1 *= math.Gamma(e) * p / (math.Gamma(a+d1) * math.Gamma(b+d1))

			y *= p / (math.Gamma(a+d2) * math.Gamma(b+d2))
			if aid&int32(1) != int32(0) {
				y = -y
			}

			q = math.Pow(s, id) // s to the id power
			if id > 0 {
				y *= q
			} else {
				y1 *= q
			}

			y += y1
		psidon:
			goto done
		}
	}

	// Use defining power series if no special cases
	y = hys2f1(a, b, c, x, &err)

done:
	*loss = err
	return y
}

/* Defining power series expansion of Gauss hypergeometric function */
// loss estimates loss of significance upon return.
func hys2f1(a float64, b float64, c float64, x float64, loss *float64) float64 {
	var f float64
	var g float64
	var h float64
	var k float64
	var m float64
	var s float64
	var u float64
	var umax float64
	var i int32

	i = 0
	umax = 0
	f = a
	g = b
	h = c
	s = 1
	u = 1
	k = 0
	for {
		if math.Abs(h) < 1e-13 {
			*loss = 1
			return math.MaxFloat64
		}
		m = k + 1
		u = u * ((f + k) * (g + k) * x / ((h + k) * m))
		s += u
		k = math.Abs(u) // remember largest term summed
		if k > umax {
			umax = k
		}
		k = m

		i++
		if i > 10000 { // should never happen
			*loss = 1
			return s
		}

		if math.Abs(u/s) <= mACHEP {
			break
		}
	}

	*loss = mACHEP*umax/math.Abs(s) + mACHEP*float64(i)

	return s
}
