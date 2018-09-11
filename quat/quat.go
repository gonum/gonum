// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package quat

import (
	"fmt"
	"strings"
)

var zero Quat

// Quat is a float64 precision quaternion.
type Quat struct {
	Real, Imag, Jmag, Kmag float64
}

// Format implements fmt.Formatter.
func (q Quat) Format(fs fmt.State, c rune) {
	prec, pOk := fs.Precision()
	if !pOk {
		prec = -1
	}
	width, wOk := fs.Precision()
	if !wOk {
		width = -1
	}
	switch c {
	case 'v':
		if fs.Flag('#') {
			fmt.Fprintf(fs, "%T{%v, %v, %v, %v}", q, q.Real, q.Imag, q.Jmag, q.Kmag)
			return
		}
		c = 'g'
		prec = -1
		fallthrough
	case 'e', 'E', 'f', 'F', 'g', 'G':
		fre := fmtString(fs, c, prec, width, false)
		fim := fmtString(fs, c, prec, width, true)
		fmt.Fprintf(fs, fmt.Sprintf("(%s%[2]si%[2]sj%[2]sk)", fre, fim), q.Real, q.Imag, q.Jmag, q.Kmag)
	default:
		fmt.Fprintf(fs, "%%!%c(%T=%[2]v)", c, q)
		return
	}
}

// This is horrible, but it's what we have.
func fmtString(fs fmt.State, c rune, prec, width int, wantPlus bool) string {
	var b strings.Builder
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
func Add(x, y Quat) Quat {
	return Quat{
		Real: x.Real + y.Real,
		Imag: x.Imag + y.Imag,
		Jmag: x.Jmag + y.Jmag,
		Kmag: x.Kmag + y.Kmag,
	}
}

// Sub returns the difference of x and y, x-y.
func Sub(x, y Quat) Quat {
	return Quat{
		Real: x.Real - y.Real,
		Imag: x.Imag - y.Imag,
		Jmag: x.Jmag - y.Jmag,
		Kmag: x.Kmag - y.Kmag,
	}
}

// Mul returns the Hamiltonian product of x and y.
func Mul(x, y Quat) Quat {
	return Quat{
		Real: x.Real*y.Real - x.Imag*y.Imag - x.Jmag*y.Jmag - x.Kmag*y.Kmag,
		Imag: x.Real*y.Imag + x.Imag*y.Real + x.Jmag*y.Kmag - x.Kmag*y.Jmag,
		Jmag: x.Real*y.Jmag - x.Imag*y.Kmag + x.Jmag*y.Real + x.Kmag*y.Imag,
		Kmag: x.Real*y.Kmag + x.Imag*y.Jmag - x.Jmag*y.Imag + x.Kmag*y.Real,
	}
}

// Scale returns q scaled by f.
func Scale(f float64, q Quat) Quat {
	return Quat{Real: f * q.Real, Imag: f * q.Imag, Jmag: f * q.Jmag, Kmag: f * q.Kmag}
}
