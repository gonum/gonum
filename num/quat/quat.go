// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package quat

import (
	"bytes"
	"fmt"
	"strconv"
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

// Parse converts the string s to a Quat. The string may be parenthesized and
// has the format [±]N±Ni±Nj±Nk. The order of the components is not strict.
func Parse(s string) (Quat, error) {
	if len(s) == 0 {
		return Quat{}, parseError{state: -1}
	}
	orig := s

	wantClose := s[0] == '('
	if wantClose {
		if s[len(s)-1] != ')' {
			return Quat{}, parseError{string: orig, state: -1}
		}
		s = s[1 : len(s)-1]
	}
	if len(s) == 0 {
		return Quat{}, parseError{string: orig, state: -1}
	}

	var q Quat
	var parts byte
	for i := 0; i < 4; i++ {
		end, p, err := floatPart(s)
		if err != nil {
			return q, parseError{string: orig, state: -1}
		}
		if parts&(1<<p) != 0 {
			return q, parseError{string: orig, state: -1}
		}
		parts |= 1 << p
		var v float64
		switch s[:end] {
		case "-":
			if len(s[end:]) == 0 {
				return q, parseError{string: orig, state: -1}
			}
			v = -1
		case "+":
			if len(s[end:]) == 0 {
				return q, parseError{string: orig, state: -1}
			}
			v = 1
		default:
			v, err = strconv.ParseFloat(s[:end], 64)
			if err != nil {
				return q, err
			}
		}
		s = s[end:]
		switch p {
		case 0:
			q.Real = v
		case 1:
			q.Imag = v
			s = s[1:]
		case 2:
			q.Jmag = v
			s = s[1:]
		case 3:
			q.Kmag = v
			s = s[1:]
		}
		if len(s) == 0 {
			return q, nil
		}
		if !isSign(rune(s[0])) {
			return q, parseError{string: orig, state: -1}
		}
	}

	return q, parseError{string: orig, state: -1}
}

func floatPart(s string) (end int, part uint, err error) {
	const (
		wantMantSign = iota
		wantMantInt
		wantMantFrac
		wantExpSign
		wantExpInt
	)
	var i, state int
	var r rune
	for i, r = range s {
		switch state {
		case wantMantSign:
			switch {
			default:
				return i, 0, parseError{string: s, state: state, rune: r}
			case isSign(r) || isDigit(r):
				state = wantMantInt
			case isDot(r):
				state = wantMantFrac
			}

		case wantMantInt:
			switch {
			default:
				return i, 0, parseError{string: s, state: state, rune: r}
			case isDigit(r):
				// Do nothing
			case isDot(r):
				state = wantMantFrac
			case isExponent(r):
				state = wantExpSign
			case isSign(r):
				return i, 0, nil
			case r == 'i':
				return i, 1, nil
			case r == 'j':
				return i, 2, nil
			case r == 'k':
				return i, 3, nil
			}

		case wantMantFrac:
			switch {
			default:
				return i, 0, parseError{string: s, state: state, rune: r}
			case isDigit(r):
				// Do nothing
			case isExponent(r):
				state = wantExpSign
			case isSign(r):
				return i, 0, nil
			case r == 'i':
				return i, 1, nil
			case r == 'j':
				return i, 2, nil
			case r == 'k':
				return i, 3, nil
			}

		case wantExpSign:
			switch {
			default:
				return i, 0, parseError{string: s, state: state, rune: r}
			case isSign(r) || isDigit(r):
				state = wantExpInt
			}

		case wantExpInt:
			switch {
			default:
				return i, 0, parseError{string: s, state: state, rune: r}
			case isDigit(r):
				// Do nothing
			case isSign(r):
				return i, 0, nil
			case r == 'i':
				return i, 1, nil
			case r == 'j':
				return i, 2, nil
			case r == 'k':
				return i, 3, nil
			}
		}
	}
	switch state {
	case wantMantSign, wantExpSign, wantExpInt:
		if state == wantExpInt && isDigit(r) {
			break
		}
		return i, 0, parseError{string: s, state: state, rune: r}
	}
	return len(s), 0, nil
}

func isSign(r rune) bool {
	return r == '+' || r == '-'
}

func isDigit(r rune) bool {
	return '0' <= r && r <= '9'
}

func isExponent(r rune) bool {
	return r == 'e' || r == 'E'
}

func isDot(r rune) bool {
	return r == '.'
}

type parseError struct {
	string string
	state  int
	rune   rune
}

func (e parseError) Error() string {
	if e.state < 0 {
		return fmt.Sprintf("quat: failed to parse: %q", e.string)
	}
	return fmt.Sprintf("quat: failed to parse in state %d with %q: %q", e.state, e.rune, e.string)
}
