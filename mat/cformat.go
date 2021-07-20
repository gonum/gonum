package mat

import (
	"fmt"
	"strconv"
)

type ColumnWidths []int

func (c ColumnWidths) width(i int) int   { return c[i] }
func (c ColumnWidths) setWidth(i, w int) { c[i] = w }

type UniformWidth int

func (u *UniformWidth) width(_ int) int   { return int(*u) }
func (u *UniformWidth) setWidth(_, w int) { *u = UniformWidth(w) }

type CFormatOption func(*cformatter)

type cformatter struct {
	matrix  CMatrix
	prefix  string
	margin  int
	dot     byte
	squeeze bool

	format func(m CMatrix, prefix string, margin int, dot byte, squueze bool, fs fmt.State, c rune)
}

// Formatted returns a fmt.Formatter for the matrix m using the given options.
func CFormatted(m CMatrix, options ...CFormatOption) fmt.Formatter {
	f := cformatter{
		matrix: m,
		dot:    '.',
	}

	for _, o := range options {
		o(&f)
	}

	return f
}

func (f cformatter) Format(fs fmt.State, c rune) {
	if c == 'v' && fs.Flag('#') && f.format == nil {
		fmt.Fprintf(fs, "%#v", f.matrix)
		return
	}
	if f.format == nil {
		f.format = cformat
	}
	f.format(f.matrix, f.prefix, f.margin, f.dot, f.squeeze, fs, c)
}

func CPrefix(p string) CFormatOption {
	return func(f *cformatter) { f.prefix = p }
}

func CExcerpt(m int) CFormatOption {
	return func(f *cformatter) { f.margin = m }
}

func CDotByte(b byte) CFormatOption {
	return func(f *cformatter) { f.dot = b }
}

func CSqueeze() CFormatOption {
	return func(f *cformatter) { f.squeeze = true }
}

func cformat(m CMatrix, prefix string, margin int, dot byte, squueze bool, fs fmt.State, c rune) {
	var (
		maxWidth int
		widths   widther
		buf, pad []byte
	)
	rows, cols := m.Dims()

	var printed int
	if margin <= 0 {
		printed = rows
		if cols > printed {
			printed = cols
		}
	} else {
		printed = margin
	}

	prec, p0k := fs.Precision()
	if !p0k {
		prec = -1
	}

	if squueze {
		widths = make(ColumnWidths, cols)
	} else {
		widths = new(UniformWidth)
	}
	switch c {
	case 'v', 'e', 'E', 'f', 'F', 'g', 'G':
		if c == 'v' {
			buf, maxWidth = cmaxCellWidth(m, 'g', printed, prec, widths)
		} else {
			buf, maxWidth = cmaxCellWidth(m, c, printed, prec, widths)
		}
	default:
		fmt.Fprintf(fs, "%%!%c(%T=Dims(%d, %d))", c, m, rows, cols)
		return
	}

	width, _ := fs.Width()
	width = max(width, maxWidth)
	pad = make([]byte, max(width, 1))
	for i := range pad {
		pad[i] = ' '
	}

	first := true
	if rows > 2*printed || cols > 2*printed {
		first = false
		fmt.Fprintf(fs, "Dims(%d %d)\n", rows, cols)
	}

	skipZero := fs.Flag(' ')
	for i := 0; i < rows; i++ {
		var el string

		if !first {
			fmt.Fprint(fs, prefix)
		}
		first = false

		switch {
		case rows == 1:
			fmt.Fprint(fs, "[")
			el = "]"
		case i == 0:
			fmt.Fprint(fs, "⎡")
			el = "⎤\n"
		case i < rows-1:
			fmt.Fprint(fs, "⎢")
			el = "⎥\n"
		default:
			fmt.Fprint(fs, "⎣")
			el = "⎦"
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

			v := m.At(i, j)
			if v == 0 && skipZero {
				buf = buf[:1]
				buf[0] = dot
			} else {
				if c == 'v' {
					if real(v) != 0 {
						buf = strconv.AppendFloat(buf[:0], real(v), 'g', prec, 64)
						buf = append(buf, "+"...)
						buf = strconv.AppendFloat(buf, imag(v), 'g', prec, 64)
						buf = append(buf, "i"...)
					} else {
						buf = strconv.AppendFloat(buf[:0], imag(v), 'g', prec, 64)
						buf = append(buf, "i"...)
					}
				} else {
					if real(v) != 0 {
						buf = strconv.AppendFloat(buf[:0], real(v), byte(c), prec, 64)
						buf = append(buf, "+"...)
						buf = strconv.AppendFloat(buf, imag(v), byte(c), prec, 64)
						buf = append(buf, "i"...)
					} else {
						buf = strconv.AppendFloat(buf[:0], imag(v), byte(c), prec, 64)
						buf = append(buf, "i"...)
					}
				}
			}
			if fs.Flag('-') {
				fs.Write(buf)
				fs.Write(pad[:widths.width(j)-len(buf)])
			} else {
				fs.Write(pad[:widths.width(j)-len(buf)])
				fs.Write(buf)
			}

			if j < cols-1 {
				fs.Write(pad[:2])
			}
		}

		fmt.Fprint(fs, el)

		if i >= printed-1 && i < rows-printed && 2*printed < rows {
			i = rows - printed - 1
			fmt.Fprintf(fs, "%s .\n%[1]s .\n%[1]s .\n", prefix)
			continue
		}
	}
}

func cmaxCellWidth(m CMatrix, c rune, printed, prec int, w widther) ([]byte, int) {
	var (
		buf        = make([]byte, 0, 64)
		rows, cols = m.Dims()
		max        int
	)
	for i := 0; i < rows; i++ {
		if i >= printed-1 && i < rows-printed && 2*printed < rows {
			i = rows - printed - 1
			continue
		}
		for j := 0; j < cols; j++ {
			if j >= printed && j < cols-printed {
				continue
			}

			value := m.At(i, j)
			if real(value) != 0 {
				buf = strconv.AppendFloat(buf, real(value), byte(c), prec, 64)
				buf = append(buf, "+"...)
			}
			buf = strconv.AppendFloat(buf, imag(value), byte(c), prec, 64)
			buf = append(buf, "i"...)
			if len(buf) > max {
				max = len(buf)
			}
			if len(buf) > w.width(j) {
				w.setWidth(j, len(buf))
			}
			buf = buf[:0]
		}
	}
	return buf, max
}
