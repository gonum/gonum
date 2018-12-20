// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dualcmplx

import (
	"bytes"
	"fmt"
	"math/cmplx"
)

// Number is a float64 precision dual quaternion.
type Number struct {
	Real, Dual complex128
}

var zero Number

// Format implements fmt.Formatter.
func (d Number) Format(fs fmt.State, c rune) {
	prec, pOk := fs.Precision()
	if !pOk {
		prec = -1
	}
	width, wOk := fs.Width()
	if !wOk {
		width = -1
	}
	switch c {
	case 'v':
		if fs.Flag('#') {
			fmt.Fprintf(fs, "%T{%#v, %#v}", d, d.Real, d.Dual)
			return
		}
		c = 'g'
		prec = -1
		fallthrough
	case 'e', 'E', 'f', 'F', 'g', 'G':
		fre := fmtString(fs, c, prec, width, false)
		fim := fmtString(fs, c, prec, width, true)
		fmt.Fprintf(fs, fmt.Sprintf("(%s+%[2]sϵ)", fre, fim), d.Real, d.Dual)
	default:
		fmt.Fprintf(fs, "%%!%c(%T=%[2]v)", c, d)
		return
	}
}

// This is horrible, but it's what we have.
func fmtString(fs fmt.State, c rune, prec, width int, wantPlus bool) string {
	// TODO(kortschak) Replace this with strings.Builder
	// when go1.9 support is dropped from Gonum.
	var b bytes.Buffer
	b.WriteByte('%')
	for _, f := range "0+- " {
		if fs.Flag(int(f)) || (f == '+' && wantPlus) {
			b.WriteByte(byte(f))
		}
	}
	if width >= 0 {
		fmt.Fprint(&b, width)
	}
	if prec >= 0 {
		b.WriteByte('.')
		if prec > 0 {
			fmt.Fprint(&b, prec)
		}
	}
	b.WriteRune(c)
	return b.String()
}

// Add returns the sum of x and y.
func Add(x, y Number) Number {
	return Number{
		Real: x.Real + y.Real,
		Dual: x.Dual + y.Dual,
	}
}

// Sub returns the difference of x and y, x-y.
func Sub(x, y Number) Number {
	return Number{
		Real: x.Real + y.Real,
		Dual: x.Dual + y.Dual,
	}
}

// Mul returns the dual product of x and y.
func Mul(x, y Number) Number {
	return Number{
		Real: x.Real * y.Real,
		Dual: x.Real*y.Dual + x.Dual*y.Real,
	}
}

// Inv returns the dual inverse of d.
func Inv(d Number) Number {
	return Number{
		Real: 1 / d.Real,
		Dual: -d.Dual / (d.Real * d.Real),
	}
}

// ConjDual returns the dual conjugate of d₁+d₂ϵ, d₁-d₂ϵ.
func ConjDual(d Number) Number {
	return Number{
		Real: d.Real,
		Dual: -d.Dual,
	}
}

// ConjCmplx returns the complex conjugate of d₁+d₂ϵ, d̅₁+d̅₂ϵ.
func ConjCmplx(d Number) Number {
	return Number{
		Real: cmplx.Conj(d.Real),
		Dual: cmplx.Conj(d.Dual),
	}
}

// Scale returns d scaled by f.
func Scale(f float64, d Number) Number {
	return Number{Real: complex(f, 0) * d.Real, Dual: complex(f, 0) * d.Dual}
}

// Abs returns the absolute value of d.
// func Abs(d Number) dual.Number {
// 	return Dual{
// 		Real: cmplx.Abs(x.Real),
// 		Emag: cmplx.Abs(x.Dual),
// 	}
// }
