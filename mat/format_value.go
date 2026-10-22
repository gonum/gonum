// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"fmt"
	"strconv"
)

const (
	syntaxDefault = iota
	syntaxMATLAB
	syntaxPython

	defaultDotByte = '.'
)

// valueFormatter is a type that can set options for formatting values and that
// satisfies the fmt.Formatter interface.
type valueFormatter interface {
	fmt.Formatter

	setDotByte(b byte)
	setSyntax(s int)
}

// valueFormatOption is a functional option for value formatting.
type valueFormatOption func(valueFormatter)

// valueDotByte sets the dot character to b. The dot character is used to
// replace a value if the result is printed with the fmt ' ' verb flag. Without
// a valueDotByte option, the default dot character is '.'.
func valueDotByte(b byte) valueFormatOption {
	return func(f valueFormatter) { f.setDotByte(b) }
}

// valueFormatMATLAB sets the printing behavior to output MATLAB syntax. If
// MATLAB syntax is specified, the ' ' verb flag is ignored.
func valueFormatMATLAB() valueFormatOption {
	return func(f valueFormatter) { f.setSyntax(syntaxMATLAB) }
}

// valueFormatPython sets the printing behavior to output Python syntax. If
// Python syntax is specified, the ' ' verb flag is ignored.
func valueFormatPython() valueFormatOption {
	return func(f valueFormatter) { f.setSyntax(syntaxPython) }
}

// formattedFloat returns a fmt.Formatter for the floating point value v using
// the given options.
func formattedFloat(v float64, options ...valueFormatOption) fmt.Formatter {
	f := floatFormatter{
		value:  v,
		buf:    make([]byte, 0, 64),
		dot:    defaultDotByte,
		syntax: syntaxDefault,
	}
	for _, o := range options {
		o(&f)
	}
	return f
}

// floatFormatter formats 64-bit a floating point value and satisfies the
// valueFormatter interface. A floatFormatter utilizes an internal buffer while
// generating formatting output and this buffer persists with the formatter.
type floatFormatter struct {
	value  float64
	buf    []byte
	dot    byte
	syntax int

	format func(v float64, buf []byte, dot byte, syntax int, fs fmt.State, c rune)
}

var _ valueFormatter = (*floatFormatter)(nil)

// Format satisfies the fmt.Formatter interface.
func (f floatFormatter) Format(fs fmt.State, c rune) {
	if f.format == nil {
		f.format = formatFloat
	}
	f.format(f.value, f.buf, f.dot, f.syntax, fs, c)
}

func (f *floatFormatter) setDotByte(b byte) { f.dot = b }
func (f *floatFormatter) setSyntax(s int)   { f.syntax = s }

// formatFloat prints a representation of v to the fs io.Writer.  The format
// character c specifies the numerical representation; valid values are those
// for float64 specified in the fmt package, with their associated flags. In
// addition to this, a space preceding a verb indicates that zero values should
// be represented by the dot character.
//
// formatFloat will not provide Go syntax output.
func formatFloat(v float64, buf []byte, dot byte, syntax int, fs fmt.State, c rune) {
	buf = buf[:0]

	if (v == 0) && (syntax == syntaxDefault) && fs.Flag(' ') {
		buf = append(buf, dot)
	} else {
		prec, ok := fs.Precision()
		if !ok {
			prec = -1
		}

		if c == 'v' {
			c = 'g'
		}

		if (v >= 0) && fs.Flag('+') {
			buf = append(buf, '+')
		}

		buf = strconv.AppendFloat(buf, v, byte(c), prec, 64)
	}

	width, ok := fs.Width()
	if ok && (width > len(buf)) {
		l := len(buf) // l is the length in bytes of unpadded content.
		for i := 0; i < width-l; i++ {
			buf = append(buf, ' ')
		}

		if !fs.Flag('-') {
			fs.Write(buf[l:])
			buf = buf[:l]
		}
	}

	fs.Write(buf)
}

// formattedComplex returns a fmt.Formatter for the complex value v using
// the given options.
func formattedComplex(v complex128, options ...valueFormatOption) fmt.Formatter {
	f := complexFormatter{
		value:  v,
		buf:    make([]byte, 0, 128),
		dot:    defaultDotByte,
		syntax: syntaxDefault,
	}
	for _, o := range options {
		o(&f)
	}
	return f
}

// floatComplex formats 128-bit complex value and satisfies the valueFormatter
// interface. A complexFormatter utilizes an internal buffer while generating
// formatting output and this buffer persists with the formatter.
type complexFormatter struct {
	value  complex128
	buf    []byte
	dot    byte
	syntax int

	format func(v complex128, buf []byte, dot byte, syntax int, fs fmt.State, c rune)
}

var _ valueFormatter = (*complexFormatter)(nil)

// Format satisfies the fmt.Formatter interface.
func (f complexFormatter) Format(fs fmt.State, c rune) {
	if f.format == nil {
		f.format = formatComplex
	}
	f.format(f.value, f.buf, f.dot, f.syntax, fs, c)
}

func (f *complexFormatter) setDotByte(b byte) { f.dot = b }
func (f *complexFormatter) setSyntax(s int)   { f.syntax = s }

// formatComplex prints a representation of v to the fs io.Writer.  The format
// character c specifies the numerical representation; valid values are those
// for float64 specified in the fmt package, with their associated flags. In
// addition to this, a space preceding a verb indicates that zero values should
// be represented by the dot character.
//
// formatComplex will not provide Go syntax output.
func formatComplex(v complex128, buf []byte, dot byte, syntax int, fs fmt.State, c rune) {
	buf = buf[:0]

	if (v == 0) && (syntax == syntaxDefault) && fs.Flag(' ') {
		buf = append(buf, dot)
	} else {
		prec, ok := fs.Precision()
		if !ok {
			prec = -1
		}

		if c == 'v' {
			c = 'g'
		}

		if (real(v) >= 0) && fs.Flag('+') {
			buf = append(buf, '+')
		}

		// Append real component.
		if real(v) != 0 {
			buf = strconv.AppendFloat(buf, real(v), byte(c), prec, 64)
		}

		if (imag(v) >= 0) && ((real(v) != 0) || fs.Flag('+')) {
			buf = append(buf, '+')
		}

		// Append imaginary component.
		buf = strconv.AppendFloat(buf, imag(v), byte(c), prec, 64)

		// Append notation for imaginary component.
		if syntax == syntaxPython {
			buf = append(buf, 'j')
		} else {
			buf = append(buf, 'i')
		}
	}

	width, ok := fs.Width()
	if ok && (width > len(buf)) {
		l := len(buf) // l is the length in bytes of unpadded content.
		for i := 0; i < width-l; i++ {
			buf = append(buf, ' ')
		}

		if !fs.Flag('-') {
			fs.Write(buf[l:])
			buf = buf[:l]
		}
	}

	fs.Write(buf)
}
