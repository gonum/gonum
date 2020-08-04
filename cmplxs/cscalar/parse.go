// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file

package cscalar

import (
	"fmt"
	"math/cmplx"
	"strconv"
	"strings"
)

// parse converts the string s to a complex128. The string may be parenthesized and
// has the format [±]N±Ni. The order of the components is not strict.
func parse(s string) (complex128, error) {
	if len(s) == 0 {
		return 0, parseError{state: -1}
	}
	orig := s

	wantClose := s[0] == '('
	if wantClose {
		if s[len(s)-1] != ')' {
			return 0, parseError{string: orig, state: -1}
		}
		s = s[1 : len(s)-1]
	}
	if len(s) == 0 {
		return 0, parseError{string: orig, state: -1}
	}
	switch s[0] {
	case 'n', 'N':
		if strings.ToLower(s) == "nan" {
			return cmplx.NaN(), nil
		}
	case 'i', 'I':
		if strings.ToLower(s) == "inf" {
			return cmplx.Inf(), nil
		}
	}

	var q complex128
	var parts byte
	for i := 0; i < 4; i++ {
		beg, end, p, err := floatPart(s)
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
			v, err = strconv.ParseFloat(s[beg:end], 64)
			if err != nil {
				return q, err
			}
		}
		s = s[end:]
		switch p {
		case 0:
			q += complex(v, 0)
		case 1:
			q += complex(0, v)
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

func floatPart(s string) (beg, end int, part uint, err error) {
	const (
		wantMantSign = iota
		wantMantIntInit
		wantMantInt
		wantMantFrac
		wantExpSign
		wantExpInt

		wantInfN
		wantInfF
		wantCloseInf

		wantNaNA
		wantNaNN
		wantCloseNaN
	)
	var i, state int
	var r rune
	for i, r = range s {
		switch state {
		case wantMantSign:
			switch {
			default:
				return 0, i, 0, parseError{string: s, state: state, rune: r}
			case isSign(r):
				state = wantMantIntInit
			case isDigit(r):
				state = wantMantInt
			case isDot(r):
				state = wantMantFrac
			case r == 'i', r == 'I':
				state = wantInfN
			case r == 'n', r == 'N':
				state = wantNaNA
			}

		case wantMantIntInit:
			switch {
			default:
				return 0, i, 0, parseError{string: s, state: state, rune: r}
			case isDigit(r):
				state = wantMantInt
			case isDot(r):
				state = wantMantFrac
			case r == 'i':
				// We need to sneak a look-ahead here.
				if i == len(s)-1 || s[i+1] == '-' || s[i+1] == '+' {
					return 0, i, 1, nil
				}
				fallthrough
			case r == 'I':
				state = wantInfN
			case r == 'n', r == 'N':
				state = wantNaNA
			}

		case wantMantInt:
			switch {
			default:
				return 0, i, 0, parseError{string: s, state: state, rune: r}
			case isDigit(r):
				// Do nothing
			case isDot(r):
				state = wantMantFrac
			case isExponent(r):
				state = wantExpSign
			case isSign(r):
				return 0, i, 0, nil
			case r == 'i':
				return 0, i, 1, nil
			}

		case wantMantFrac:
			switch {
			default:
				return 0, i, 0, parseError{string: s, state: state, rune: r}
			case isDigit(r):
				// Do nothing
			case isExponent(r):
				state = wantExpSign
			case isSign(r):
				return 0, i, 0, nil
			case r == 'i':
				return 0, i, 1, nil
			}

		case wantExpSign:
			switch {
			default:
				return 0, i, 0, parseError{string: s, state: state, rune: r}
			case isSign(r) || isDigit(r):
				state = wantExpInt
			}

		case wantExpInt:
			switch {
			default:
				return 0, i, 0, parseError{string: s, state: state, rune: r}
			case isDigit(r):
				// Do nothing
			case isSign(r):
				return 0, i, 0, nil
			case r == 'i':
				return 0, i, 1, nil
			}

		case wantInfN:
			if r != 'n' && r != 'N' {
				return 0, i, 0, parseError{string: s, state: state, rune: r}
			}
			state = wantInfF
		case wantInfF:
			if r != 'f' && r != 'F' {
				return 0, i, 0, parseError{string: s, state: state, rune: r}
			}
			state = wantCloseInf
		case wantCloseInf:
			switch {
			default:
				return 0, i, 0, parseError{string: s, state: state, rune: r}
			case isSign(r):
				return 0, i, 0, nil
			case r == 'i':
				return 0, i, 1, nil
			}

		case wantNaNA:
			if r != 'a' && r != 'A' {
				return 0, i, 0, parseError{string: s, state: state, rune: r}
			}
			state = wantNaNN
		case wantNaNN:
			if r != 'n' && r != 'N' {
				return 0, i, 0, parseError{string: s, state: state, rune: r}
			}
			state = wantCloseNaN
		case wantCloseNaN:
			if isSign(rune(s[0])) {
				beg = 1
			}
			switch {
			default:
				return beg, i, 0, parseError{string: s, state: state, rune: r}
			case isSign(r):
				return beg, i, 0, nil
			case r == 'i':
				return beg, i, 1, nil
			}
		}
	}
	switch state {
	case wantMantSign, wantExpSign, wantExpInt:
		if state == wantExpInt && isDigit(r) {
			break
		}
		return 0, i, 0, parseError{string: s, state: state, rune: r}
	}
	return 0, len(s), 0, nil
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
