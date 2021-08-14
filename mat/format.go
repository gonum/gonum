// Copyright ©2013 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"bytes"
	"fmt"
	"unicode/utf8"
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
	return f
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
	return f
}

// formatter is a matrix formatter that satisfies the fmt.Formatter interface.
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
	dimser

	// NewValueFormatter returns a new, zero-valued fmt.Formatter.
	NewValueFormatter(options ...valueFormatOption) fmt.Formatter

	// RevalueFormatter sets the value of a formatter to element i, j.
	// A use is to recycle an existing value formatter for efficiency.
	RevalueFormatter(f fmt.Formatter, i, j int) fmt.Formatter

	// ValueFormatter returns a new fmt.Formatter for element i, j.
	ValueFormatter(i, j int, options ...valueFormatOption) fmt.Formatter

	// Matrix returns the source matrix.
	Mat() dimser
}

// formattableMatrix is a formattable for floating point-valued matrices, those
// that satisfy the Matrix interface.
type formattableMatrix struct{ mat Matrix }

var _ formattable = formattableMatrix{}

func (m formattableMatrix) Dims() (r, c int) { return m.mat.Dims() }
func (m formattableMatrix) Mat() dimser      { return m.mat }
func (m formattableMatrix) NewValueFormatter(options ...valueFormatOption) fmt.Formatter {
	f := formattedFloat(0, options...).(floatFormatter)
	return &f
}
func (m formattableMatrix) RevalueFormatter(f fmt.Formatter, i, j int) fmt.Formatter {
	f.(*floatFormatter).value = m.mat.At(i, j)
	return f
}
func (m formattableMatrix) ValueFormatter(i, j int, options ...valueFormatOption) fmt.Formatter {
	return formattedFloat(m.mat.At(i, j), options...)
}

// formattableCMatrix is a formattable for complex-valued matrices, those that
// satisfy the CMatrix interface.
type formattableCMatrix struct{ mat CMatrix }

var _ formattable = formattableCMatrix{}

func (m formattableCMatrix) Dims() (r, c int) { return m.mat.Dims() }
func (m formattableCMatrix) Mat() dimser      { return m.mat }
func (m formattableCMatrix) NewValueFormatter(options ...valueFormatOption) fmt.Formatter {
	f := formattedComplex(0, options...).(complexFormatter)
	return &f
}
func (m formattableCMatrix) RevalueFormatter(f fmt.Formatter, i, j int) fmt.Formatter {
	f.(*complexFormatter).value = m.mat.At(i, j)
	return f
}
func (m formattableCMatrix) ValueFormatter(i, j int, options ...valueFormatOption) fmt.Formatter {
	return formattedComplex(m.mat.At(i, j), options...)
}

// dimser type is a two-dimensional array having rows and columns with table
// elements accessible by a formattable. Details regarding value type and
// methods for element access are left to the formattable.
type dimser interface {
	Dims() (r, c int)
}

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
		buf        = make([]byte, 0, 3) // Three bytes required because of Unicode decoration.
		f          = m.NewValueFormatter(valueDotByte(dot))
		rows, cols = m.Dims()
		s          = newState(fs)
	)

	// Emit matrix characteristics if given verb is not in printable set.
	switch c {
	case 'v', 'e', 'E', 'f', 'F', 'g', 'G':
	default:
		fmt.Fprintf(fs, "%%!%c(%T=Dims(%d, %d))", c, m.Mat(), rows, cols)
		return
	}

	// Caculate printed margin.
	var printed = max(rows, cols)
	if margin > 0 {
		printed = min(printed, margin)
	}

	// The widthing state must be fit to the formattable in order to
	// horizontally align cells.
	s.fit(m, printed, squeeze, fs, c, f)

	// Prepend dimensions if mat overflows printed region.
	if (rows > 2*printed) || (cols > 2*printed) {
		fmt.Fprintf(fs, "Dims(%d, %d)\n%s", rows, cols, prefix)
	}

	// Write matrix content.
	var (
		lbrack, rbrack rune // lbrack and rbrack runes store Unicode characters.
		n              int  // n is a counter for bytes written by utf.EncodeRune.
	)
	for i := 0; i < rows; i++ {
		// Matrix rows are printed on their own lines.  These lines may
		// be prefixed.
		if i > 0 {
			buf = append(buf[:0], '\n')
			fs.Write(buf)

			if prefix != "" {
				fmt.Fprint(fs, prefix)
			}
		}

		// The first three rows below the upper printable region
		// contain only a dot left indented by one space. Additional
		// rows beyond the upper printable region but above the lower
		// printable region are skipped.
		if (i >= printed) && (i < rows-printed) && (2*printed < rows) {
			buf = append(buf[:0], ' ', '.')
			fs.Write(buf)

			// After three sequences, advance to the next printable
			// row.
			if i == (2 + printed) {
				i = rows - printed - 1
			}
			continue
		}

		// Set left and right pair of brackets (or bracket pieces).
		switch {
		case rows == 1:
			lbrack = '[' // Left square bracket (U+005B).
			rbrack = ']' // Right square bracket (U+005D).
		case i == 0:
			lbrack = '⎡' // Left square bracket upper corner (U+23A1).
			rbrack = '⎤' // Right square bracket upper corner (U+23A4)
		case i < rows-1:
			lbrack = '⎢' // Left square bracket extension (U+23A2).
			rbrack = '⎥' // Right square bracket extension (U+23A5).
		default:
			lbrack = '⎣' // Left square bracket lower corner (U+23A3).
			rbrack = '⎦' // Right square bracket lower corner (U+23A6).
		}

		// Write left bracket (or bracket piece).
		n = utf8.EncodeRune(buf[:3], lbrack)
		fs.Write(buf[:n])

		// Write row content.
		for j := 0; j < cols; j++ {
			// Elements in a matrix row are separated by two spaces.
			if j > 0 {
				buf = append(buf[:0], ' ', ' ')
				fs.Write(buf)
			}

			// The first two columns beyond the left printable
			// region are filled with a three-character sequence,
			// consisting of dots (if on first or last row) or
			// spaces (if on an inner row). Additional columns
			// beyond the left printable region but before the
			// right printable region are skipped.
			if (j >= printed) && (j < cols-printed) {
				if (i == 0) || (i == rows-1) {
					buf = append(buf[:0], '.', '.', '.')
				} else {
					buf = append(buf[:0], ' ', ' ', ' ')
				}
				fs.Write(buf)

				// After two sequences, advance to the next
				// printable column.
				if j == (1 + printed) {
					j = cols - printed - 1
				}
				continue
			}

			// Write cell content.
			f = m.RevalueFormatter(f, i, j)
			f.Format(s.At(j), c)
		}

		// Write right bracket (or bracket piece).
		n = utf8.EncodeRune(buf[:3], rbrack)
		fs.Write(buf[:n])
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
		buf        = make([]byte, 0, 2)
		f          = m.NewValueFormatter(valueFormatMATLAB())
		rows, cols = m.Dims()
		s          = newState(fs)
	)

	// Emit matrix characteristics if given verb is not in printable set.
	switch c {
	case 'v', 'e', 'E', 'f', 'F', 'g', 'G':
	default:
		fmt.Fprintf(fs, "%%!%c(%T=Dims(%d, %d))", c, m.Mat(), rows, cols)
		return
	}

	// Read flag for layout syntax.
	var layout bool = fs.Flag('#')

	// In layout syntax, the widthing state must be fit to the formattable
	// in order to horizontally align cells.
	if layout {
		s.fit(m, max(rows, cols), squeeze, fs, c, f)
	}

	// Write an opening bracket.
	buf = append(buf[:0], '[')
	fs.Write(buf)

	// Write matrix content.
	for i := 0; i < rows; i++ {
		// In standard syntax, matrix rows are separated by a semicolon
		// and one space.
		if !layout && (i > 0) {
			buf = append(buf[:0], ';', ' ')
			fs.Write(buf)
		}

		// In layout syntax, matrix rows are printed on their own lines
		// unless the matrix consists of only a single row. These lines
		// may be prefixed and are always left-indented by one space.
		if layout && (rows > 1) {
			if prefix == "" {
				buf = append(buf[:0], '\n', ' ')
				fs.Write(buf)
			} else {
				fmt.Fprintf(fs, "\n%s ", prefix)
			}
		}

		// Write the row content.
		for j := 0; j < cols; j++ {
			// Elements in a matrix row are separated by one space.
			if j > 0 {
				buf = append(buf[:0], ' ')
				fs.Write(buf)
			}

			// Write cell content.
			f = m.RevalueFormatter(f, i, j)
			f.Format(s.At(j), c)
		}
	}

	// In layout syntax, the closing bracket is printed on its
	// own line unless the matrix consists of only a single row.
	if layout && (rows > 1) {
		buf = append(buf[:0], '\n', ']')
		fs.Write(buf)
		return
	}

	// Write a closing bracket.
	buf = append(buf[:0], ']')
	fs.Write(buf)
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
		buf        = make([]byte, 0, 3)
		f          = m.NewValueFormatter(valueFormatPython())
		rows, cols = m.Dims()
		s          = newState(fs)
	)

	// Emit matrix characteristics if given verb is not in printable set.
	switch c {
	case 'v', 'e', 'E', 'f', 'F', 'g', 'G':
	default:
		fmt.Fprintf(fs, "%%!%c(%T=Dims(%d, %d))", c, m.Mat(), rows, cols)
		return
	}

	// Read flag for layout syntax.
	var layout bool = fs.Flag('#')

	// In layout syntax, the widthing state must be fit to the formattable
	// in order to horizontally align cells.
	if layout {
		s.fit(m, max(rows, cols), squeeze, fs, c, f)
	}

	// Write an opening bracket.
	buf = append(buf[:0], '[')
	fs.Write(buf)

	// Write matrix content.
	for i := 0; i < rows; i++ {
		// In standard syntax, matrix rows are separated by a semicolon
		// and one space.
		if !layout && (i > 0) {
			buf = append(buf[:0], ',', ' ')
			fs.Write(buf)
		}

		// In layout syntax, matrix rows are separated by a comma and
		// a new line.  Matrix rows after the first are printed on
		// their own lines. These lines may be prefixed and are always
		// left-indented by one space.
		if layout && (i > 0) {
			if prefix == "" {
				buf = append(buf[:0], ',', '\n', ' ')
				fs.Write(buf)
			} else {
				fmt.Fprintf(fs, ",\n%s ", prefix)
			}
		}

		// Write left bracket.
		if rows > 1 {
			buf = append(buf[:0], '[')
			fs.Write(buf)
		}

		// Write the row content.
		for j := 0; j < cols; j++ {
			// Elements in a matrix row are separated by a comma
			// and one space.
			if j > 0 {
				buf = append(buf[:0], ',', ' ')
				fs.Write(buf)
			}

			// Write cell content.
			f = m.RevalueFormatter(f, i, j)
			f.Format(s.At(j), c)
		}

		// Write right bracket.
		if rows > 1 {
			buf = append(buf[:0], ']')
			fs.Write(buf)
		}
	}

	// Write a closing bracket.
	buf = append(buf[:0], ']')
	fs.Write(buf)
}

// state satisfies fmt.State and may vary its responses to Width.
type state struct {
	write     func(b []byte) (n int, err error)
	width     func() (wid int, ok bool)
	precision func() (prec int, ok bool)
	flag      func(c int) bool

	w widther
	p int
}

var _ fmt.State = state{}

// newState returns a new state, shadowing given fmt.State fs.
func newState(fs fmt.State) state {
	s := state{
		write:     fs.Write,
		width:     fs.Width,
		precision: fs.Precision,
		flag:      fs.Flag,
	}
	return s
}

// fit fits a state to formattable m, storing a single width for the
// formattable or, if squeeze is true, individual widths are stored for each
// column in the formattable. After fitting, the state responds to Width calls
// with the value in p, which may be updated for a new column by calling At.
//
// The fitted state is not safe for concurrent use.
func (s *state) fit(m formattable, printed int, squeeze bool, fs fmt.State, c rune, f fmt.Formatter) *state {
	rows, cols := m.Dims()

	// Assign appropriate widther.
	if squeeze {
		s.w = make(columnWidth, cols)
	} else {
		s.w = new(uniformWidth)
	}

	// Set state to temporarily write to a buffer instead of fs.Write.
	var buf bytes.Buffer
	s.write = buf.Write

	// Fit state to m.
	for i := 0; i < rows; i++ {
		if i >= printed-1 && i < rows-printed && 2*printed < rows {
			i = rows - printed - 1
			continue
		}
		for j := 0; j < cols; j++ {
			if j >= printed && j < cols-printed {
				continue
			}

			_ = m.RevalueFormatter(f, i, j)
			f.Format(s, c)

			if buf.Len() > s.w.width(j) {
				s.w.setWidth(j, buf.Len())
			}

			buf.Reset()
		}
	}

	// Restore writes by the state to fs.
	s.write = fs.Write

	// State should now respond to width with p, defaulting to fs.Width.
	s.width = func() (wid int, ok bool) { return s.p, true }
	if p, ok := fs.Width(); ok {
		s.p = p
	}

	return s
}

// Write carries the same meaning as in fmt.State.
func (s state) Write(b []byte) (n int, err error) {
	return s.write(b)
}

// Width carries the same meaning as in fmt.State.
func (s state) Width() (wid int, ok bool) {
	return s.width()
}

// Precision carries the same meaning as in fmt.State.
func (s state) Precision() (prec int, ok bool) {
	return s.precision()
}

// Flag carries the same meaning as in fmt.State.
func (s state) Flag(c int) bool {
	return s.flag(c)
}

// At will provide a copy of the state that responds to subsequent Width calls
// with the widther value for column i.
func (s *state) At(i int) *state {
	if s.w != nil {
		s.p = s.w.width(i)
	}
	return s
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

var _ widther = columnWidth{}

func (c columnWidth) width(i int) int   { return c[i] }
func (c columnWidth) setWidth(i, w int) { c[i] = w }
