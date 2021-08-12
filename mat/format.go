// Copyright ©2013 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"bytes"
	"fmt"
)

// Formatted returns a fmt.Formatter for the matrix m using the given options.
func Formatted(m Matrix, options ...FormatOption) fmt.Formatter {
	f := formatter{
		matrix: formattableMatrix{mat: m},
		dot:    '.',
	}
	for _, o := range options {
		o(&f)
	}
	return &f
}

// CFormatted returns a fmt.Formatter for the complex-valued matrix m using the
// given options.
func CFormatted(m CMatrix, options ...FormatOption) fmt.Formatter {
	f := formatter{
		matrix: formattableCMatrix{mat: m},
		dot:    '.',
	}
	for _, o := range options {
		o(&f)
	}
	return &f
}

// formatter is a matrix formatter that satisfies the fmt.Formatter interface
type formatter struct {
	matrix  formattable
	prefix  string
	margin  int
	dot     byte
	squeeze bool

	format func(m formattable, prefix string, margin int, dot byte, squeeze bool, fs fmt.State, c rune)
}

var _ fmt.Formatter = (*formatter)(nil)

// FormatOption is a functional option for matrix formatting.
type FormatOption func(*formatter)

// Prefix sets the formatted prefix to the string p. Prefix is a string that is prepended to
// each line of output after the first line.
func Prefix(p string) FormatOption {
	return func(f *formatter) { f.prefix = p }
}

// Excerpt sets the maximum number of rows and columns to print at the margins of the matrix
// to m. If m is zero or less all elements are printed.
func Excerpt(m int) FormatOption {
	return func(f *formatter) { f.margin = m }
}

// DotByte sets the dot character to b. The dot character is used to replace zero elements
// if the result is printed with the fmt ' ' verb flag. Without a DotByte option, the default
// dot character is '.'.
func DotByte(b byte) FormatOption {
	return func(f *formatter) { f.dot = b }
}

// Squeeze sets the printing behavior to minimise column width for each individual column.
func Squeeze() FormatOption {
	return func(f *formatter) { f.squeeze = true }
}

// FormatMATLAB sets the printing behavior to output MATLAB syntax. If MATLAB syntax is
// specified, the ' ' verb flag and Excerpt option are ignored. If the alternative syntax
// verb flag, '#' is used the matrix is formatted in rows and columns.
func FormatMATLAB() FormatOption {
	return func(f *formatter) { f.format = formatMATLAB }
}

// FormatPython sets the printing behavior to output Python syntax. If Python syntax is
// specified, the ' ' verb flag and Excerpt option are ignored. If the alternative syntax
// verb flag, '#' is used the matrix is formatted in rows and columns.
func FormatPython() FormatOption {
	return func(f *formatter) { f.format = formatPython }
}

// Format satisfies the fmt.Formatter interface.
func (f formatter) Format(fs fmt.State, c rune) {
	if c == 'v' && fs.Flag('#') && f.format == nil {
		fmt.Fprintf(fs, "%#v", f.matrix.Mat())
		return
	}
	if f.format == nil {
		f.format = format
	}
	f.format(f.matrix, f.prefix, f.margin, f.dot, f.squeeze, fs, c)
}

// formattable holds a matrix and provides methods for formatting values.
type formattable interface {
	tabular

	// NewValueFormatter returns a new, zero-valued fmt.Formatter
	NewValueFormatter(options ...vFormatOption) fmt.Formatter

	// RevalueFormatter sets the value of a formatter to element i, j.
	// A use is to recycle an existing value formatter for efficiency.
	RevalueFormatter(f fmt.Formatter, i, j int) fmt.Formatter

	// ValueFormatter returns a new fmt.Formatter for element i, j
	ValueFormatter(i, j int, options ...vFormatOption) fmt.Formatter

	// Matrix returns the source matrix
	Mat() tabular
}

// formattableMatrix is a formattable for floating point-valued matrices, those
// that satisfy the Matrix interface.
type formattableMatrix struct{ mat Matrix }

var _ formattable = (*formattableMatrix)(nil)

func (m formattableMatrix) Dims() (r, c int) { return m.mat.Dims() }
func (m formattableMatrix) Mat() tabular     { return m.mat }
func (m formattableMatrix) NewValueFormatter(options ...vFormatOption) fmt.Formatter {
	return formattedFloat(0, options...)
}
func (m formattableMatrix) RevalueFormatter(f fmt.Formatter, i, j int) fmt.Formatter {
	f.(*floatFormatter).value = m.mat.At(i, j)
	return f
}
func (m formattableMatrix) ValueFormatter(i, j int, options ...vFormatOption) fmt.Formatter {
	return formattedFloat(m.mat.At(i, j), options...)
}

// formattableCMatrix is a formattable for complex-valued matrices, those that
// satisfy the CMatrix interface
type formattableCMatrix struct{ mat CMatrix }

var _ formattable = (*formattableCMatrix)(nil)

func (m formattableCMatrix) Dims() (r, c int) { return m.mat.Dims() }
func (m formattableCMatrix) Mat() tabular     { return m.mat }
func (m formattableCMatrix) NewValueFormatter(options ...vFormatOption) fmt.Formatter {
	return formattedComplex(0, options...)
}
func (m formattableCMatrix) RevalueFormatter(f fmt.Formatter, i, j int) fmt.Formatter {
	f.(*complexFormatter).value = m.mat.At(i, j)
	return f
}
func (m formattableCMatrix) ValueFormatter(i, j int, options ...vFormatOption) fmt.Formatter {
	return formattedComplex(m.mat.At(i, j), options...)
}

// tabular type is a two-dimensional array having rows and columns with table
// elements accessible by a formattable. Details regarding value type and
// methods for element access are left to the formattable.
type tabular interface {
	Dims() (r, c int)
}

var _ tabular = (Matrix)(nil)

// format prints a pretty representation of m to the fs io.Writer. The format character c
// specifies the numerical representation of elements; valid values are those for float64
// specified in the fmt package, with their associated flags. In addition to this, a space
// preceding a verb indicates that zero values should be represented by the dot character.
// The printed range of the matrix can be limited by specifying a positive value for margin;
// If margin is greater than zero, only the first and last margin rows/columns of the matrix
// are output. If squeeze is true, column widths are determined on a per-column basis.
//
// format will not provide Go syntax output.
func format(m formattable, prefix string, margin int, dot byte, squeeze bool, fs fmt.State, c rune) {
	var (
		f          = m.NewValueFormatter(vDotByte(dot))
		rows, cols = m.Dims()
	)

	var printed int
	if margin <= 0 {
		printed = rows
		if cols > printed {
			printed = cols
		}
	} else {
		printed = margin
	}

	var ws = newState(fs)
	switch c {
	case 'v', 'e', 'E', 'f', 'F', 'g', 'G':
		ws.fit(m, printed, squeeze, fs, c, f)
	default:
		fmt.Fprintf(fs, "%%!%c(%T=Dims(%d, %d))", c, m.Mat(), rows, cols)
		return
	}

	first := true
	if rows > 2*printed || cols > 2*printed {
		first = false
		fmt.Fprintf(fs, "Dims(%d, %d)\n", rows, cols)
	}

	var buf = make([]byte, 0, 4)
	for i := 0; i < rows; i++ {
		if !first {
			fmt.Fprint(fs, prefix)
		}
		first = false
		switch {
		case rows == 1:
			fmt.Fprint(fs, "[")
			buf = append(buf[:0], []byte("]")...)
		case i == 0:
			fmt.Fprint(fs, "⎡")
			buf = append(buf[:0], []byte("⎤\n")...)
		case i < rows-1:
			fmt.Fprint(fs, "⎢")
			buf = append(buf[:0], []byte("⎥\n")...)
		default:
			fmt.Fprint(fs, "⎣")
			buf = append(buf[:0], []byte("⎦")...)
		}

		for j := 0; j < cols; j++ {
			if j >= printed && j < cols-printed {
				j = cols - printed - 1
				if i == 0 || i == rows-1 {
					fmt.Fprint(fs, "...  ...  ")
				} else {
					fmt.Fprint(fs, "          ")
				}
				continue
			}

			f = m.RevalueFormatter(f, i, j)
			f.Format(ws.At(j), c)

			if j < cols-1 {
				fmt.Fprintf(fs, "  ")
			}
		}

		fs.Write(buf)

		if i >= printed-1 && i < rows-printed && 2*printed < rows {
			i = rows - printed - 1
			fmt.Fprintf(fs, "%s .\n%[1]s .\n%[1]s .\n", prefix)
			continue
		}
	}
}

// formatMATLAB prints a MATLAB representation of m to the fs io.Writer. The format character c
// specifies the numerical representation of elements; valid values are those for float64
// specified in the fmt package, with their associated flags.
// The printed range of the matrix can be limited by specifying a positive value for margin;
// If squeeze is true, column widths are determined on a per-column basis.
//
// formatMATLAB will not provide Go syntax output.
func formatMATLAB(m formattable, prefix string, _ int, _ byte, squeeze bool, fs fmt.State, c rune) {
	var (
		f          = m.NewValueFormatter(vFormatMATLAB())
		rows, cols = m.Dims()
	)

	if !fs.Flag('#') {
		switch c {
		case 'v', 'e', 'E', 'f', 'F', 'g', 'G':
		default:
			fmt.Fprintf(fs, "%%!%c(%T=Dims(%d, %d))", c, m.Mat(), rows, cols)
			return
		}
		fmt.Fprint(fs, "[")
		for i := 0; i < rows; i++ {
			if i != 0 {
				fmt.Fprint(fs, "; ")
			}
			for j := 0; j < cols; j++ {
				if j != 0 {
					fmt.Fprint(fs, " ")
				}
				f = m.RevalueFormatter(f, i, j)
				f.Format(fs, c)
			}
		}
		fmt.Fprint(fs, "]")
		return
	}

	printed := rows
	if cols > printed {
		printed = cols
	}

	var ws = newState(fs)
	switch c {
	case 'v', 'e', 'E', 'f', 'F', 'g', 'G':
		ws.fit(m, printed, squeeze, fs, c, f)
	default:
		fmt.Fprintf(fs, "%%!%c(%T=Dims(%d, %d))", c, m.Mat(), rows, cols)
		return
	}

	var buf = make([]byte, 0, 4)
	for i := 0; i < rows; i++ {
		switch {
		case rows == 1:
			fmt.Fprint(fs, "[")
			buf = append(buf[:0], []byte("]")...)
		case i == 0:
			fmt.Fprint(fs, "[\n"+prefix+" ")
			buf = append(buf[:0], []byte("\n")...)
		case i < rows-1:
			fmt.Fprint(fs, prefix+" ")
			buf = append(buf[:0], []byte("\n")...)
		default:
			fmt.Fprint(fs, prefix+" ")
			buf = append(buf[:0], []byte("\n"+prefix+"]")...)
		}

		for j := 0; j < cols; j++ {
			f = m.RevalueFormatter(f, i, j)
			f.Format(ws.At(j), c)

			if j < cols-1 {
				fmt.Fprint(fs, " ")
			}
		}

		fs.Write(buf)
	}
}

// formatPython prints a Python representation of m to the fs io.Writer. The format character c
// specifies the numerical representation of elements; valid values are those for float64
// specified in the fmt package, with their associated flags.
// The printed range of the matrix can be limited by specifying a positive value for margin;
// If squeeze is true, column widths are determined on a per-column basis.
//
// formatPython will not provide Go syntax output.
func formatPython(m formattable, prefix string, _ int, _ byte, squeeze bool, fs fmt.State, c rune) {
	var (
		f          = m.NewValueFormatter(vFormatPython())
		rows, cols = m.Dims()
	)

	if !fs.Flag('#') {
		switch c {
		case 'v', 'e', 'E', 'f', 'F', 'g', 'G':
		default:
			fmt.Fprintf(fs, "%%!%c(%T=Dims(%d, %d))", c, m.Mat(), rows, cols)
			return
		}
		fmt.Fprint(fs, "[")
		if rows > 1 {
			fmt.Fprint(fs, "[")
		}
		for i := 0; i < rows; i++ {
			if i != 0 {
				fmt.Fprint(fs, "], [")
			}
			for j := 0; j < cols; j++ {
				if j != 0 {
					fmt.Fprint(fs, ", ")
				}
				f = m.RevalueFormatter(f, i, j)
				f.Format(fs, c)
			}
		}
		if rows > 1 {
			fmt.Fprint(fs, "]")
		}
		fmt.Fprint(fs, "]")
		return
	}

	printed := rows
	if cols > printed {
		printed = cols
	}

	var ws = newState(fs)
	switch c {
	case 'v', 'e', 'E', 'f', 'F', 'g', 'G':
		ws.fit(m, printed, squeeze, fs, c, f)
	default:
		fmt.Fprintf(fs, "%%!%c(%T=Dims(%d, %d))", c, m.Mat(), rows, cols)
		return
	}

	var buf = make([]byte, 0, 4)
	for i := 0; i < rows; i++ {
		if i != 0 {
			fmt.Fprint(fs, prefix)
		}
		switch {
		case rows == 1:
			fmt.Fprint(fs, "[")
			buf = append(buf[:0], []byte("]")...)
		case i == 0:
			fmt.Fprint(fs, "[[")
			buf = append(buf[:0], []byte("],\n")...)
		case i < rows-1:
			fmt.Fprint(fs, " [")
			buf = append(buf[:0], []byte("],\n")...)
		default:
			fmt.Fprint(fs, " [")
			buf = append(buf[:0], []byte("]]")...)
		}

		for j := 0; j < cols; j++ {
			f = m.RevalueFormatter(f, i, j)
			f.Format(ws.At(j), c)

			if j < cols-1 {
				fmt.Fprint(fs, ", ")
			}
		}

		fs.Write(buf)
	}
}

// state satisfies fmt.State and may varies its responses to Width.
type state struct {
	write     func(b []byte) (n int, err error)
	width     func() (wid int, ok bool)
	precision func() (prec int, ok bool)
	flag      func(c int) bool

	w widther
	p int
}

var _ fmt.State = (*state)(nil)

// newState return a new state, shadowing given fmt.State fs.
func newState(fs fmt.State) *state {
	ws := state{
		write:     fs.Write,
		width:     fs.Width,
		precision: fs.Precision,
		flag:      fs.Flag,
	}
	return &ws
}

// fit fits a state to formattable m, storing a single width for the
// formattable or, if squeeze is true, individual widths are stored for each
// column in the formattable. After fitment, the state responds to Width calls
// with the value in p, which may be updated for a new column by calling At.
//
// The fitted state is not safe for concurrent use.
func (ws *state) fit(m formattable, printed int, squeeze bool, fs fmt.State, c rune, f fmt.Formatter) *state {
	rows, cols := m.Dims()

	// assign appropriate widther
	if squeeze {
		ws.w = make(columnWidth, cols)
	} else {
		ws.w = new(uniformWidth)
	}

	// set state to temporarily write to a buffer instead of fs.Write
	var buf bytes.Buffer
	ws.write = buf.Write

	// fit state to m
	for i := 0; i < rows; i++ {
		if i >= printed-1 && i < rows-printed && 2*printed < rows {
			i = rows - printed - 1
			continue
		}
		for j := 0; j < cols; j++ {
			if j >= printed && j < cols-printed {
				continue
			}

			f = m.RevalueFormatter(f, i, j)
			f.Format(ws, c)

			if buf.Len() > ws.w.width(j) {
				ws.w.setWidth(j, buf.Len())
			}

			buf.Reset()
		}
	}

	// restore writes by the state to fs
	ws.write = fs.Write

	// state should now response to width with p, defaulting to fs.Width
	ws.width = func() (wid int, ok bool) { return ws.p, true }
	if p, ok := fs.Width(); ok {
		ws.p = p
	}

	return ws
}

// Write carries the same meaning as in fmt.State
func (ws state) Write(b []byte) (n int, err error) {
	return ws.write(b)
}

// Width carries the same meaning as in fmt.State
func (ws state) Width() (wid int, ok bool) {
	return ws.width()
}

// Precision carries the same meaning as in fmt.State
func (ws state) Precision() (prec int, ok bool) {
	return ws.precision()
}

// Flag carries the same meaning as in fmt.State
func (ws state) Flag(c int) bool {
	return ws.flag(c)
}

// At will provide a copy of the state that responds to subsequent Width calls
// with the widther value for column i.
func (ws *state) At(i int) *state {
	ws.p = ws.w.width(i)
	return ws
}

type widther interface {
	width(i int) int
	setWidth(i, w int)
}

type uniformWidth int

var _ widther = (*uniformWidth)(nil)

func (u *uniformWidth) width(_ int) int   { return int(*u) }
func (u *uniformWidth) setWidth(_, w int) { *u = uniformWidth(w) }

type columnWidth []int

var _ widther = (*columnWidth)(nil)

func (c columnWidth) width(i int) int   { return c[i] }
func (c columnWidth) setWidth(i, w int) { c[i] = w }
