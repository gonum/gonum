// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package quat

import (
	"math"
	"math/cmplx"
	"testing"
)

var sinTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{}},
	{q: Quat{Real: math.Pi / 2}, want: Quat{Real: 1}},
	{q: Quat{Imag: math.Pi / 2}, want: Quat{Imag: imag(cmplx.Sin(complex(0, math.Pi/2)))}},
	{q: Quat{Jmag: math.Pi / 2}, want: Quat{Jmag: imag(cmplx.Sin(complex(0, math.Pi/2)))}},
	{q: Quat{Kmag: math.Pi / 2}, want: Quat{Kmag: imag(cmplx.Sin(complex(0, math.Pi/2)))}},

	// Exercises from Real Quaternionic Calculus Handbook doi:10.1007/978-3-0348-0622-0
	// Ex 6.159 (a) and (b).
	{q: Quat{Real: 1, Imag: 1, Jmag: 1, Kmag: 1}, want: func() Quat {
		p := math.Cos(1) * math.Sinh(math.Sqrt(3)) / math.Sqrt(3)
		// An error exists in the book's given solution for the real part.
		return Quat{Real: math.Sin(1) * math.Cosh(math.Sqrt(3)), Imag: p, Jmag: p, Kmag: p}
	}()},
	{q: Quat{Imag: -2, Jmag: 1}, want: func() Quat {
		s := math.Sinh(math.Sqrt(5)) / math.Sqrt(5)
		return Quat{Imag: -2 * s, Jmag: s}
	}()},
}

func TestSin(t *testing.T) {
	const tol = 1e-14
	for _, test := range sinTests {
		got := Sin(test.q)
		if !equalApprox(got, test.want, tol) {
			t.Errorf("unexpected result for Sin(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var sinhTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{}},
	{q: Quat{Real: math.Pi / 2}, want: Quat{Real: math.Sinh(math.Pi / 2)}},
	{q: Quat{Imag: math.Pi / 2}, want: Quat{Imag: imag(cmplx.Sinh(complex(0, math.Pi/2)))}},
	{q: Quat{Jmag: math.Pi / 2}, want: Quat{Jmag: imag(cmplx.Sinh(complex(0, math.Pi/2)))}},
	{q: Quat{Kmag: math.Pi / 2}, want: Quat{Kmag: imag(cmplx.Sinh(complex(0, math.Pi/2)))}},
	{q: Quat{Real: 1, Imag: -1, Jmag: -1}, want: func() Quat {
		// This was based on the example on p118, but it too has an error.
		q := Quat{Real: 1, Imag: -1, Jmag: -1}
		return Scale(0.5, Sub(Exp(q), Exp(Scale(-1, q))))
	}()},
	{q: Quat{1, 1, 1, 1}, want: func() Quat {
		q := Quat{1, 1, 1, 1}
		return Scale(0.5, Sub(Exp(q), Exp(Scale(-1, q))))
	}()},
	{q: Asinh(Quat{1, 1, 1, 1}), want: Quat{1, 1, 1, 1}},
	{q: Asinh(Quat{1, 1, 1, 1}), want: func() Quat {
		q := Asinh(Quat{1, 1, 1, 1})
		return Scale(0.5, Sub(Exp(q), Exp(Scale(-1, q))))
	}()},
}

func TestSinh(t *testing.T) {
	const tol = 1e-14
	for _, test := range sinhTests {
		got := Sinh(test.q)
		if !equalApprox(got, test.want, tol) {
			t.Errorf("unexpected result for Sinh(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var cosTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{Real: 1}},
	{q: Quat{Real: math.Pi / 2}, want: Quat{Real: 0}},
	{q: Quat{Imag: math.Pi / 2}, want: Quat{Real: real(cmplx.Cos(complex(0, math.Pi/2)))}},
	{q: Quat{Jmag: math.Pi / 2}, want: Quat{Real: real(cmplx.Cos(complex(0, math.Pi/2)))}},
	{q: Quat{Kmag: math.Pi / 2}, want: Quat{Real: real(cmplx.Cos(complex(0, math.Pi/2)))}},

	// Example from Real Quaternionic Calculus Handbook doi:10.1007/978-3-0348-0622-0
	// p108.
	{q: Quat{Real: 1, Imag: 1, Jmag: 1, Kmag: 1}, want: func() Quat {
		p := math.Sin(1) * math.Sinh(math.Sqrt(3)) / math.Sqrt(3)
		return Quat{Real: math.Cos(1) * math.Cosh(math.Sqrt(3)), Imag: -p, Jmag: -p, Kmag: -p}
	}()},
}

func TestCos(t *testing.T) {
	const tol = 1e-14
	for _, test := range cosTests {
		got := Cos(test.q)
		if !equalApprox(got, test.want, tol) {
			t.Errorf("unexpected result for Cos(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var coshTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{Real: 1}},
	{q: Quat{Real: math.Pi / 2}, want: Quat{Real: math.Cosh(math.Pi / 2)}},
	{q: Quat{Imag: math.Pi / 2}, want: Quat{Imag: imag(cmplx.Cosh(complex(0, math.Pi/2)))}},
	{q: Quat{Jmag: math.Pi / 2}, want: Quat{Jmag: imag(cmplx.Cosh(complex(0, math.Pi/2)))}},
	{q: Quat{Kmag: math.Pi / 2}, want: Quat{Kmag: imag(cmplx.Cosh(complex(0, math.Pi/2)))}},
	{q: Quat{Real: 1, Imag: -1, Jmag: -1}, want: func() Quat {
		q := Quat{Real: 1, Imag: -1, Jmag: -1}
		return Scale(0.5, Add(Exp(q), Exp(Scale(-1, q))))
	}()},
	{q: Quat{1, 1, 1, 1}, want: func() Quat {
		q := Quat{1, 1, 1, 1}
		return Scale(0.5, Add(Exp(q), Exp(Scale(-1, q))))
	}()},
}

func TestCosh(t *testing.T) {
	const tol = 1e-14
	for _, test := range coshTests {
		got := Cosh(test.q)
		if !equalApprox(got, test.want, tol) {
			t.Errorf("unexpected result for Cosh(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var tanTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{}},
	{q: Quat{Real: math.Pi / 4}, want: Quat{Real: math.Tan(math.Pi / 4)}},
	{q: Quat{Imag: math.Pi / 4}, want: Quat{Imag: imag(cmplx.Tan(complex(0, math.Pi/4)))}},
	{q: Quat{Jmag: math.Pi / 4}, want: Quat{Jmag: imag(cmplx.Tan(complex(0, math.Pi/4)))}},
	{q: Quat{Kmag: math.Pi / 4}, want: Quat{Kmag: imag(cmplx.Tan(complex(0, math.Pi/4)))}},

	// From exercise from Real Quaternionic Calculus Handbook doi:10.1007/978-3-0348-0622-0
	{q: Quat{Imag: 1}, want: Mul(Sin(Quat{Imag: 1}), Inv(Cos(Quat{Imag: 1})))},
	{q: Quat{1, 1, 1, 1}, want: Mul(Sin(Quat{1, 1, 1, 1}), Inv(Cos(Quat{1, 1, 1, 1})))},
}

func TestTan(t *testing.T) {
	const tol = 1e-14
	for _, test := range tanTests {
		got := Tan(test.q)
		if !equalApprox(got, test.want, tol) {
			t.Errorf("unexpected result for Tan(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var tanhTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{}},
	{q: Quat{Real: math.Pi / 4}, want: Quat{Real: math.Tanh(math.Pi / 4)}},
	{q: Quat{Imag: math.Pi / 4}, want: Quat{Imag: imag(cmplx.Tanh(complex(0, math.Pi/4)))}},
	{q: Quat{Jmag: math.Pi / 4}, want: Quat{Jmag: imag(cmplx.Tanh(complex(0, math.Pi/4)))}},
	{q: Quat{Kmag: math.Pi / 4}, want: Quat{Kmag: imag(cmplx.Tanh(complex(0, math.Pi/4)))}},
	{q: Quat{Imag: 1}, want: Mul(Sinh(Quat{Imag: 1}), Inv(Cosh(Quat{Imag: 1})))},
	{q: Quat{1, 1, 1, 1}, want: Mul(Sinh(Quat{1, 1, 1, 1}), Inv(Cosh(Quat{1, 1, 1, 1})))},
}

func TestTanh(t *testing.T) {
	const tol = 1e-14
	for _, test := range tanhTests {
		got := Tanh(test.q)
		if !equalApprox(got, test.want, tol) {
			t.Errorf("unexpected result for Tanh(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var asinTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{}},
	{q: Quat{Real: 1}, want: Quat{Real: math.Pi / 2}},
	{q: Quat{Imag: 1}, want: Quat{Imag: real(cmplx.Asinh(1))}},
	{q: Quat{Jmag: 1}, want: Quat{Jmag: real(cmplx.Asinh(1))}},
	{q: Quat{Kmag: 1}, want: Quat{Kmag: real(cmplx.Asinh(1))}},
	{q: Sin(Quat{1, 1, 1, 1}), want: Quat{1, 1, 1, 1}},
}

func TestAsin(t *testing.T) {
	const tol = 1e-14
	for _, test := range asinTests {
		got := Asin(test.q)
		if !equalApprox(got, test.want, tol) {
			t.Errorf("unexpected result for Asin(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var asinhTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{}},
	{q: Quat{Real: 1}, want: Quat{Real: math.Asinh(1)}},
	{q: Quat{Imag: 1}, want: Quat{Imag: math.Pi / 2}},
	{q: Quat{Jmag: 1}, want: Quat{Jmag: math.Pi / 2}},
	{q: Quat{Kmag: 1}, want: Quat{Kmag: math.Pi / 2}},
	{q: Quat{1, 1, 1, 1}, want: func() Quat {
		q := Quat{1, 1, 1, 1}
		return Log(Add(q, Sqrt(Add(Mul(q, q), Quat{Real: 1}))))
	}()},
	{q: Sinh(Quat{Real: 1}), want: Quat{Real: 1}},
	{q: Sinh(Quat{Imag: 1}), want: Quat{Imag: 1}},
	{q: Sinh(Quat{Imag: 1, Jmag: 1}), want: Quat{Imag: 1, Jmag: 1}},
	{q: Sinh(Quat{Real: 1, Imag: 1, Jmag: 1}), want: Quat{Real: 1, Imag: 1, Jmag: 1}},
	// The following fails:
	// {q: Sinh(Quat{1, 1, 1, 1}), want: Quat{1, 1, 1, 1}},
	// but this passes...
	{q: Sinh(Quat{1, 1, 1, 1}), want: func() Quat {
		q := Sinh(Quat{1, 1, 1, 1})
		return Log(Add(q, Sqrt(Add(Mul(q, q), Quat{Real: 1}))))
	}()},
	// And see the Sinh tests that do the reciprocal operation.
}

func TestAsinh(t *testing.T) {
	const tol = 1e-14
	for _, test := range asinhTests {
		got := Asinh(test.q)
		if !equalApprox(got, test.want, tol) {
			t.Errorf("unexpected result for Asinh(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var acosTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{Real: math.Pi / 2}},
	{q: Quat{Real: 1}, want: Quat{Real: 0}},
	{q: Quat{Imag: 1}, want: Quat{Real: real(cmplx.Acos(1i)), Imag: imag(cmplx.Acos(1i))}},
	{q: Quat{Jmag: 1}, want: Quat{Real: real(cmplx.Acos(1i)), Jmag: imag(cmplx.Acos(1i))}},
	{q: Quat{Kmag: 1}, want: Quat{Real: real(cmplx.Acos(1i)), Kmag: imag(cmplx.Acos(1i))}},
	{q: Cos(Quat{1, 1, 1, 1}), want: Quat{1, 1, 1, 1}},
}

func TestAcos(t *testing.T) {
	const tol = 1e-14
	for _, test := range acosTests {
		got := Acos(test.q)
		if !equalApprox(got, test.want, tol) {
			t.Errorf("unexpected result for Acos(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var acoshTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{Real: math.Pi / 2}},
	{q: Quat{Real: 1}, want: Quat{Real: math.Acosh(1)}},
	{q: Quat{Imag: 1}, want: Quat{Real: real(cmplx.Acosh(1i)), Imag: imag(cmplx.Acosh(1i))}},
	{q: Quat{Jmag: 1}, want: Quat{Real: real(cmplx.Acosh(1i)), Jmag: imag(cmplx.Acosh(1i))}},
	{q: Quat{Kmag: 1}, want: Quat{Real: real(cmplx.Acosh(1i)), Kmag: imag(cmplx.Acosh(1i))}},
	{q: Cosh(Quat{1, 1, 1, 1}), want: Quat{1, 1, 1, 1}},
	{q: Quat{1, 1, 1, 1}, want: func() Quat {
		q := Quat{1, 1, 1, 1}
		return Log(Add(q, Sqrt(Sub(Mul(q, q), Quat{Real: 1}))))
	}()},
	// The following fails by a factor of -1.
	// {q: Cosh(Quat{1, 1, 1, 1}), want: func() Quat {
	// 	q := Cosh(Quat{1, 1, 1, 1})
	// 	return Log(Add(q, Sqrt(Sub(Mul(q, q), Quat{Real: 1}))))
	// }()},
}

func TestAcosh(t *testing.T) {
	const tol = 1e-14
	for _, test := range acoshTests {
		got := Acosh(test.q)
		if !equalApprox(got, test.want, tol) {
			t.Errorf("unexpected result for Acosh(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var atanTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{}},
	{q: Quat{Real: 1}, want: Quat{Real: math.Pi / 4}},
	{q: Quat{Imag: 0.5}, want: Quat{Real: real(cmplx.Atan(0.5i)), Imag: imag(cmplx.Atan(0.5i))}},
	{q: Quat{Jmag: 0.5}, want: Quat{Real: real(cmplx.Atan(0.5i)), Jmag: imag(cmplx.Atan(0.5i))}},
	{q: Quat{Kmag: 0.5}, want: Quat{Real: real(cmplx.Atan(0.5i)), Kmag: imag(cmplx.Atan(0.5i))}},
	{q: Tan(Quat{1, 1, 1, 1}), want: Quat{1, 1, 1, 1}},
}

func TestAtan(t *testing.T) {
	const tol = 1e-14
	for _, test := range atanTests {
		got := Atan(test.q)
		if !equalApprox(got, test.want, tol) {
			t.Errorf("unexpected result for Atan(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}

var atanhTests = []struct {
	q    Quat
	want Quat
}{
	{q: Quat{}, want: Quat{}},
	{q: Quat{Real: 1}, want: Quat{Real: math.Atanh(1)}},
	{q: Quat{Imag: 0.5}, want: Quat{Real: real(cmplx.Atanh(0.5i)), Imag: imag(cmplx.Atanh(0.5i))}},
	{q: Quat{Jmag: 0.5}, want: Quat{Real: real(cmplx.Atanh(0.5i)), Jmag: imag(cmplx.Atanh(0.5i))}},
	{q: Quat{Kmag: 0.5}, want: Quat{Real: real(cmplx.Atanh(0.5i)), Kmag: imag(cmplx.Atanh(0.5i))}},
	{q: Quat{1, 1, 1, 1}, want: func() Quat {
		q := Quat{1, 1, 1, 1}
		return Scale(0.5, Sub(Log(Add(Quat{Real: 1}, q)), Log(Sub(Quat{Real: 1}, q))))
	}()},
	{q: Tanh(Quat{Real: 1}), want: Quat{Real: 1}},
	{q: Tanh(Quat{Imag: 1}), want: Quat{Imag: 1}},
	{q: Tanh(Quat{Imag: 1, Jmag: 1}), want: Quat{Imag: 1, Jmag: 1}},
	{q: Tanh(Quat{Real: 1, Imag: 1, Jmag: 1}), want: Quat{Real: 1, Imag: 1, Jmag: 1}},
	// The following fails
	// {q: Tanh(Quat{1, 1, 1, 1}), want: Quat{1, 1, 1, 1}},
	// but...
	{q: Tanh(Quat{1, 1, 1, 1}), want: func() Quat {
		q := Tanh(Quat{1, 1, 1, 1})
		return Scale(0.5, Sub(Log(Add(Quat{Real: 1}, q)), Log(Sub(Quat{Real: 1}, q))))
	}()},
}

func TestAtanh(t *testing.T) {
	const tol = 1e-14
	for _, test := range atanhTests {
		got := Atanh(test.q)
		if !equalApprox(got, test.want, tol) {
			t.Errorf("unexpected result for Atanh(%v): got:%v want:%v", test.q, got, test.want)
		}
	}
}
