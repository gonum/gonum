// Copyright ©2013 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"fmt"
	"math"
	"math/cmplx"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/blas/cblas128"
)

func TestFormat(t *testing.T) {
	t.Parallel()
	type rp struct {
		format string
		output string
	}
	sqrt := func(_, _ int, v float64) float64 { return math.Sqrt(v) }
	for i, test := range []struct {
		m   fmt.Formatter
		rep []rp
	}{
		// Dense matrix representation
		{
			m: Formatted(NewDense(3, 3, []float64{0, 0, 0, 0, 0, 0, 0, 0, 0})),
			rep: []rp{
				{"%v", "⎡0  0  0⎤\n⎢0  0  0⎥\n⎣0  0  0⎦"},
				{"% f", "⎡.  .  .⎤\n⎢.  .  .⎥\n⎣.  .  .⎦"},
				{"%#v", fmt.Sprintf("%#v", NewDense(3, 3, nil))},
				{"%s", "%!s(*mat.Dense=Dims(3, 3))"},
			},
		},
		{
			m: Formatted(NewDense(3, 3, []float64{1, 1, 1, 1, 1, 1, 1, 1, 1})),
			rep: []rp{
				{"%v", "⎡1  1  1⎤\n⎢1  1  1⎥\n⎣1  1  1⎦"},
				{"% f", "⎡1  1  1⎤\n⎢1  1  1⎥\n⎣1  1  1⎦"},
				{"%#v", fmt.Sprintf("%#v", NewDense(3, 3, []float64{1, 1, 1, 1, 1, 1, 1, 1, 1}))},
			},
		},
		{
			m: Formatted(NewDense(3, 3, []float64{1, 1, 1, 1, 1, 1, 1, 1, 1}), Prefix("\t")),
			rep: []rp{
				{"%v", "⎡1  1  1⎤\n\t⎢1  1  1⎥\n\t⎣1  1  1⎦"},
				{"% f", "⎡1  1  1⎤\n\t⎢1  1  1⎥\n\t⎣1  1  1⎦"},
				{"%#v", fmt.Sprintf("%#v", NewDense(3, 3, []float64{1, 1, 1, 1, 1, 1, 1, 1, 1}))},
			},
		},
		{
			m: Formatted(NewDense(3, 3, []float64{1, 0, 0, 0, 1, 0, 0, 0, 1})),
			rep: []rp{
				{"%v", "⎡1  0  0⎤\n⎢0  1  0⎥\n⎣0  0  1⎦"},
				{"% f", "⎡1  .  .⎤\n⎢.  1  .⎥\n⎣.  .  1⎦"},
				{"%#v", fmt.Sprintf("%#v", NewDense(3, 3, []float64{1, 0, 0, 0, 1, 0, 0, 0, 1}))},
			},
		},
		{
			m: Formatted(NewDense(2, 3, []float64{1, 2, 3, 4, 5, 6})),
			rep: []rp{
				{"%v", "⎡1  2  3⎤\n⎣4  5  6⎦"},
				{"% f", "⎡1  2  3⎤\n⎣4  5  6⎦"},
				{"%#v", fmt.Sprintf("%#v", NewDense(2, 3, []float64{1, 2, 3, 4, 5, 6}))},
			},
		},
		{
			m: Formatted(NewDense(3, 2, []float64{1, 2, 3, 4, 5, 6})),
			rep: []rp{
				{"%v", "⎡1  2⎤\n⎢3  4⎥\n⎣5  6⎦"},
				{"% f", "⎡1  2⎤\n⎢3  4⎥\n⎣5  6⎦"},
				{"%#v", fmt.Sprintf("%#v", NewDense(3, 2, []float64{1, 2, 3, 4, 5, 6}))},
			},
		},
		{
			m: func() fmt.Formatter {
				m := NewDense(2, 3, []float64{0, 1, 2, 3, 4, 5})
				m.Apply(sqrt, m)
				return Formatted(m)
			}(),
			rep: []rp{
				{"%v", "⎡                 0                   1  1.4142135623730951⎤\n⎣1.7320508075688772                   2    2.23606797749979⎦"},
				{"%.2f", "⎡0.00  1.00  1.41⎤\n⎣1.73  2.00  2.24⎦"},
				{"% f", "⎡                 .                   1  1.4142135623730951⎤\n⎣1.7320508075688772                   2    2.23606797749979⎦"},
				{"%#v", fmt.Sprintf("%#v", NewDense(2, 3, []float64{0, 1, 1.4142135623730951, 1.7320508075688772, 2, 2.23606797749979}))},
			},
		},
		{
			m: func() fmt.Formatter {
				m := NewDense(3, 2, []float64{0, 1, 2, 3, 4, 5})
				m.Apply(sqrt, m)
				return Formatted(m)
			}(),
			rep: []rp{
				{"%v", "⎡                 0                   1⎤\n⎢1.4142135623730951  1.7320508075688772⎥\n⎣                 2    2.23606797749979⎦"},
				{"%.2f", "⎡0.00  1.00⎤\n⎢1.41  1.73⎥\n⎣2.00  2.24⎦"},
				{"% f", "⎡                 .                   1⎤\n⎢1.4142135623730951  1.7320508075688772⎥\n⎣                 2    2.23606797749979⎦"},
				{"%#v", fmt.Sprintf("%#v", NewDense(3, 2, []float64{0, 1, 1.4142135623730951, 1.7320508075688772, 2, 2.23606797749979}))},
			},
		},
		{
			m: func() fmt.Formatter {
				m := NewDense(2, 3, []float64{0, 1, 2, 3, 4, 5})
				m.Apply(sqrt, m)
				return Formatted(m, Squeeze())
			}(),
			rep: []rp{
				{"%v", "⎡                 0  1  1.4142135623730951⎤\n⎣1.7320508075688772  2    2.23606797749979⎦"},
				{"%.2f", "⎡0.00  1.00  1.41⎤\n⎣1.73  2.00  2.24⎦"},
				{"% f", "⎡                 .  1  1.4142135623730951⎤\n⎣1.7320508075688772  2    2.23606797749979⎦"},
				{"%#v", fmt.Sprintf("%#v", NewDense(2, 3, []float64{0, 1, 1.4142135623730951, 1.7320508075688772, 2, 2.23606797749979}))},
			},
		},
		{
			m: func() fmt.Formatter {
				m := NewDense(1, 10, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
				return Formatted(m, Excerpt(3))
			}(),
			rep: []rp{
				{"%v", "Dims(1, 10)\n[ 1   2   3  ...  ...   8   9  10]"},
			},
		},
		{
			m: func() fmt.Formatter {
				m := NewDense(10, 1, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10})
				return Formatted(m, Excerpt(3))
			}(),
			rep: []rp{
				{"%v", "Dims(10, 1)\n⎡ 1⎤\n⎢ 2⎥\n⎢ 3⎥\n .\n .\n .\n⎢ 8⎥\n⎢ 9⎥\n⎣10⎦"},
			},
		},
		{
			m: func() fmt.Formatter {
				m := NewDense(10, 10, nil)
				for i := 0; i < 10; i++ {
					m.Set(i, i, 1)
				}
				return Formatted(m, Excerpt(3))
			}(),
			rep: []rp{
				{"%v", "Dims(10, 10)\n⎡1  0  0  ...  ...  0  0  0⎤\n⎢0  1  0            0  0  0⎥\n⎢0  0  1            0  0  0⎥\n .\n .\n .\n⎢0  0  0            1  0  0⎥\n⎢0  0  0            0  1  0⎥\n⎣0  0  0  ...  ...  0  0  1⎦"},
			},
		},
		{
			m: Formatted(NewDense(9, 1, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}), FormatMATLAB()),
			rep: []rp{
				{"%v", "[1; 2; 3; 4; 5; 6; 7; 8; 9]"},
				{"%#v", "[\n 1\n 2\n 3\n 4\n 5\n 6\n 7\n 8\n 9\n]"},
				{"%s", "%!s(*mat.Dense=Dims(9, 1))"},
				{"%#s", "%!s(*mat.Dense=Dims(9, 1))"},
			},
		},
		{
			m: Formatted(NewDense(1, 9, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}), FormatMATLAB()),
			rep: []rp{
				{"%v", "[1 2 3 4 5 6 7 8 9]"},
				{"%#v", "[1 2 3 4 5 6 7 8 9]"},
			},
		},
		{
			m: Formatted(NewDense(3, 3, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}), FormatMATLAB()),
			rep: []rp{
				{"%v", "[1 2 3; 4 5 6; 7 8 9]"},
				{"%#v", "[\n 1 2 3\n 4 5 6\n 7 8 9\n]"},
			},
		},
		{
			m: Formatted(NewDense(9, 1, []float64{1, -2, 3, 4, 5, 6, 7, 8, 9}), FormatMATLAB()),
			rep: []rp{
				{"%v", "[1; -2; 3; 4; 5; 6; 7; 8; 9]"},
				{"%#v", "[\n  1\n -2\n  3\n  4\n  5\n  6\n  7\n  8\n  9\n]"},
			},
		},
		{
			m: Formatted(NewDense(1, 9, []float64{1, -2, 3, 4, 5, 6, 7, 8, 9}), FormatMATLAB()),
			rep: []rp{
				{"%v", "[1 -2 3 4 5 6 7 8 9]"},
				{"%#v", "[ 1 -2  3  4  5  6  7  8  9]"},
			},
		},
		{
			m: Formatted(NewDense(3, 3, []float64{1, -2, 3, 4, 5, 6, 7, 8, 9}), FormatMATLAB()),
			rep: []rp{
				{"%v", "[1 -2 3; 4 5 6; 7 8 9]"},
				{"%#v", "[\n  1 -2  3\n  4  5  6\n  7  8  9\n]"},
			},
		},

		{
			m: Formatted(NewDense(9, 1, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}), FormatPython()),
			rep: []rp{
				{"%v", "[[1], [2], [3], [4], [5], [6], [7], [8], [9]]"},
				{"%#v", "[[1],\n [2],\n [3],\n [4],\n [5],\n [6],\n [7],\n [8],\n [9]]"},
				{"%s", "%!s(*mat.Dense=Dims(9, 1))"},
				{"%#s", "%!s(*mat.Dense=Dims(9, 1))"},
			},
		},
		{
			m: Formatted(NewDense(1, 9, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}), FormatPython()),
			rep: []rp{
				{"%v", "[1, 2, 3, 4, 5, 6, 7, 8, 9]"},
				{"%#v", "[1, 2, 3, 4, 5, 6, 7, 8, 9]"},
			},
		},
		{
			m: Formatted(NewDense(3, 3, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9}), FormatPython()),
			rep: []rp{
				{"%v", "[[1, 2, 3], [4, 5, 6], [7, 8, 9]]"},
				{"%#v", "[[1, 2, 3],\n [4, 5, 6],\n [7, 8, 9]]"},
			},
		},
		{
			m: Formatted(NewDense(9, 1, []float64{1, -2, 3, 4, 5, 6, 7, 8, 9}), FormatPython()),
			rep: []rp{
				{"%v", "[[1], [-2], [3], [4], [5], [6], [7], [8], [9]]"},
				{"%#v", "[[ 1],\n [-2],\n [ 3],\n [ 4],\n [ 5],\n [ 6],\n [ 7],\n [ 8],\n [ 9]]"},
			},
		},
		{
			m: Formatted(NewDense(1, 9, []float64{1, -2, 3, 4, 5, 6, 7, 8, 9}), FormatPython()),
			rep: []rp{
				{"%v", "[1, -2, 3, 4, 5, 6, 7, 8, 9]"},
				{"%#v", "[ 1, -2,  3,  4,  5,  6,  7,  8,  9]"},
			},
		},
		{
			m: Formatted(NewDense(3, 3, []float64{1, -2, 3, 4, 5, 6, 7, 8, 9}), FormatPython()),
			rep: []rp{
				{"%v", "[[1, -2, 3], [4, 5, 6], [7, 8, 9]]"},
				{"%#v", "[[ 1, -2,  3],\n [ 4,  5,  6],\n [ 7,  8,  9]]"},
			},
		},
	} {
		for j, rp := range test.rep {
			got := fmt.Sprintf(rp.format, test.m)
			if got != rp.output {
				t.Errorf("unexpected format result test %d part %d:\ngot:\n%s\nwant:\n%s", i, j, got, rp.output)
			}
		}
	}
}

func TestCFormat(t *testing.T) {
	t.Parallel()
	type rp struct {
		format string
		output string
	}
	for i, test := range []struct {
		m   fmt.Formatter
		rep []rp
	}{
		// Dense matrix representation with complex data
		{
			m: CFormatted(NewCDense(3, 3, []complex128{0, 0, 0, 0, 0, 0, 0, 0, 0})),
			rep: []rp{
				{"%v", "⎡0i  0i  0i⎤\n⎢0i  0i  0i⎥\n⎣0i  0i  0i⎦"},
				{"% f", "⎡.  .  .⎤\n⎢.  .  .⎥\n⎣.  .  .⎦"},
				{"%#v", fmt.Sprintf("%#v", NewCDense(3, 3, nil))},
				{"%s", "%!s(*mat.CDense=Dims(3, 3))"},
			},
		},
		{
			m: CFormatted(NewCDense(3, 3, []complex128{1, 1, 1, 1, 1, 1, 1, 1, 1})),
			rep: []rp{
				{"%v", "⎡1+0i  1+0i  1+0i⎤\n⎢1+0i  1+0i  1+0i⎥\n⎣1+0i  1+0i  1+0i⎦"},
				{"% f", "⎡1+0i  1+0i  1+0i⎤\n⎢1+0i  1+0i  1+0i⎥\n⎣1+0i  1+0i  1+0i⎦"},
				{"%#v", fmt.Sprintf("%#v", NewCDense(3, 3, []complex128{1, 1, 1, 1, 1, 1, 1, 1, 1}))},
			},
		},
		{
			m: CFormatted(NewCDense(3, 3, []complex128{1, 1, 1, 1, 1, 1, 1, 1, 1}), Prefix("\t")),
			rep: []rp{
				{"%v", "⎡1+0i  1+0i  1+0i⎤\n\t⎢1+0i  1+0i  1+0i⎥\n\t⎣1+0i  1+0i  1+0i⎦"},
				{"% f", "⎡1+0i  1+0i  1+0i⎤\n\t⎢1+0i  1+0i  1+0i⎥\n\t⎣1+0i  1+0i  1+0i⎦"},
				{"%#v", fmt.Sprintf("%#v", NewCDense(3, 3, []complex128{1, 1, 1, 1, 1, 1, 1, 1, 1}))},
			},
		},
		{
			m: CFormatted(NewCDense(3, 3, []complex128{1, 0, 0, 0, 1, 0, 0, 0, 1})),
			rep: []rp{
				{"%v", "⎡1+0i    0i    0i⎤\n⎢  0i  1+0i    0i⎥\n⎣  0i    0i  1+0i⎦"},
				{"% f", "⎡1+0i     .     .⎤\n⎢   .  1+0i     .⎥\n⎣   .     .  1+0i⎦"},
				{"%#v", fmt.Sprintf("%#v", NewCDense(3, 3, []complex128{1, 0, 0, 0, 1, 0, 0, 0, 1}))},
			},
		},
		{
			m: CFormatted(NewCDense(2, 3, []complex128{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i, 9 + 10i, 11 + 12i})),
			rep: []rp{
				{"%v", "⎡  1+2i    3+4i    5+6i⎤\n⎣  7+8i   9+10i  11+12i⎦"},
				{"% f", "⎡  1+2i    3+4i    5+6i⎤\n⎣  7+8i   9+10i  11+12i⎦"},
				{"%#v", fmt.Sprintf("%#v", NewCDense(2, 3, []complex128{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i, 9 + 10i, 11 + 12i}))},
			},
		},
		{
			m: CFormatted(NewCDense(3, 2, []complex128{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i, 9 + 10i, 11 + 12i})),
			rep: []rp{
				{"%v", "⎡  1+2i    3+4i⎤\n⎢  5+6i    7+8i⎥\n⎣ 9+10i  11+12i⎦"},
				{"% f", "⎡  1+2i    3+4i⎤\n⎢  5+6i    7+8i⎥\n⎣ 9+10i  11+12i⎦"},
				{"%#v", fmt.Sprintf("%#v", NewCDense(3, 2, []complex128{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i, 9 + 10i, 11 + 12i}))},
			},
		},
		{
			m: func() fmt.Formatter {
				M, N := 2, 3
				m := NewCDense(M, N, []complex128{0 + 0i, 1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i, 9 + 10i})
				for i := 0; i < M; i++ {
					for j := 0; j < N; j++ {
						m.Set(i, j, cmplx.Sqrt(m.At(i, j)))
					}
				}
				return CFormatted(m)
			}(),
			rep: []rp{
				{"%v", "⎡                                   0i  1.272019649514069+0.7861513777574233i                                   2+1i⎤\n⎣2.5308348104831593+1.185379617655596i   2.9690188457413544+1.34724641634978i  3.350643523793132+1.4922506570736884i⎦"},
				{"%.2f", "⎡     0.00i  1.27+0.79i  2.00+1.00i⎤\n⎣2.53+1.19i  2.97+1.35i  3.35+1.49i⎦"},
				{"% f", "⎡                                    .  1.272019649514069+0.7861513777574233i                                   2+1i⎤\n⎣2.5308348104831593+1.185379617655596i   2.9690188457413544+1.34724641634978i  3.350643523793132+1.4922506570736884i⎦"},
				{"%#v", fmt.Sprintf("%#v", NewCDense(2, 3, []complex128{
					(0 + 0i), (1.272019649514069 + 0.7861513777574233i),
					(2 + 1i), (2.5308348104831593 + 1.185379617655596i),
					(2.9690188457413544 + 1.34724641634978i), (3.350643523793132 + 1.4922506570736884i),
				}))},
			},
		},
		{
			m: func() fmt.Formatter {
				M, N := 3, 2
				m := NewCDense(M, N, []complex128{0 + 0i, 1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i, 9 + 10i})
				for i := 0; i < M; i++ {
					for j := 0; j < N; j++ {
						m.Set(i, j, cmplx.Sqrt(m.At(i, j)))
					}
				}
				return CFormatted(m)
			}(),
			rep: []rp{
				{"%v", "⎡                                   0i  1.272019649514069+0.7861513777574233i⎤\n⎢                                 2+1i  2.5308348104831593+1.185379617655596i⎥\n⎣ 2.9690188457413544+1.34724641634978i  3.350643523793132+1.4922506570736884i⎦"},
				{"%.2f", "⎡     0.00i  1.27+0.79i⎤\n⎢2.00+1.00i  2.53+1.19i⎥\n⎣2.97+1.35i  3.35+1.49i⎦"},
				{"% f", "⎡                                    .  1.272019649514069+0.7861513777574233i⎤\n⎢                                 2+1i  2.5308348104831593+1.185379617655596i⎥\n⎣ 2.9690188457413544+1.34724641634978i  3.350643523793132+1.4922506570736884i⎦"},
				{"%#v", fmt.Sprintf("%#v", NewCDense(3, 2, []complex128{
					(0 + 0i), (1.272019649514069 + 0.7861513777574233i),
					(2 + 1i), (2.5308348104831593 + 1.185379617655596i),
					(2.9690188457413544 + 1.34724641634978i), (3.350643523793132 + 1.4922506570736884i),
				}))},
			},
		},
		{
			m: func() fmt.Formatter {
				M, N := 2, 3
				m := NewCDense(M, N, []complex128{0 + 0i, 1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i, 9 + 10i})
				for i := 0; i < M; i++ {
					for j := 0; j < N; j++ {
						m.Set(i, j, cmplx.Sqrt(m.At(i, j)))
					}
				}
				return CFormatted(m, Squeeze())
			}(),
			rep: []rp{
				{"%v", "⎡                                   0i  1.272019649514069+0.7861513777574233i                                   2+1i⎤\n⎣2.5308348104831593+1.185379617655596i   2.9690188457413544+1.34724641634978i  3.350643523793132+1.4922506570736884i⎦"},
				{"%.2f", "⎡     0.00i  1.27+0.79i  2.00+1.00i⎤\n⎣2.53+1.19i  2.97+1.35i  3.35+1.49i⎦"},
				{"% f", "⎡                                    .  1.272019649514069+0.7861513777574233i                                   2+1i⎤\n⎣2.5308348104831593+1.185379617655596i   2.9690188457413544+1.34724641634978i  3.350643523793132+1.4922506570736884i⎦"},
				{"%#v", fmt.Sprintf("%#v", NewCDense(2, 3, []complex128{
					(0 + 0i), (1.272019649514069 + 0.7861513777574233i),
					(2 + 1i), (2.5308348104831593 + 1.185379617655596i),
					(2.9690188457413544 + 1.34724641634978i), (3.350643523793132 + 1.4922506570736884i),
				}))},
			},
		},
		{
			m: func() fmt.Formatter {
				M, N := 1, 10
				m := NewCDense(M, N, nil)
				for i := 0; i < M*N; i++ {
					m.Set(i%M, int(i/M), complex(float64(1+i*2), float64(2+i*2)))
				}
				return CFormatted(m, Excerpt(3))
			}(),
			rep: []rp{
				{"%v", "Dims(1, 10)\n[  1+2i    3+4i    5+6i  ...  ...  15+16i  17+18i  19+20i]"},
			},
		},
		{
			m: func() fmt.Formatter {
				M, N := 10, 1
				m := NewCDense(M, N, nil)
				for i := 0; i < M*N; i++ {
					m.Set(i%M, int(i/M), complex(float64(1+i*2), float64(2+i*2)))
				}
				return CFormatted(m, Excerpt(3))
			}(),
			rep: []rp{
				{"%v", "Dims(10, 1)\n⎡  1+2i⎤\n⎢  3+4i⎥\n⎢  5+6i⎥\n .\n .\n .\n⎢15+16i⎥\n⎢17+18i⎥\n⎣19+20i⎦"},
			},
		},
		{
			m: func() fmt.Formatter {
				M, N := 10, 10
				m := NewCDense(M, N, nil)
				for i := 0; i < M*N; i++ {
					m.Set(i%M, int(i/M), complex(float64(i%10), float64(i%10)))
				}
				return CFormatted(m, Excerpt(1))
			}(),
			rep: []rp{
				{"%v", "Dims(10, 10)\n⎡  0i  ...  ...    0i⎤\n .\n .\n .\n⎣9+9i  ...  ...  9+9i⎦"},
			},
		},
		{
			m: CFormatted(NewCDense(4, 1, []complex128{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i}), FormatMATLAB()),
			rep: []rp{
				{"%v", "[1+2i; 3+4i; 5+6i; 7+8i]"},
				{"%#v", "[\n 1+2i\n 3+4i\n 5+6i\n 7+8i\n]"},
				{"%s", "%!s(*mat.CDense=Dims(4, 1))"},
				{"%#s", "%!s(*mat.CDense=Dims(4, 1))"},
			},
		},
		{
			m: CFormatted(NewCDense(1, 4, []complex128{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i}), FormatMATLAB()),
			rep: []rp{
				{"%v", "[1+2i 3+4i 5+6i 7+8i]"},
				{"%#v", "[1+2i 3+4i 5+6i 7+8i]"},
			},
		},
		{
			m: CFormatted(NewCDense(2, 2, []complex128{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i}), FormatMATLAB()),
			rep: []rp{
				{"%v", "[1+2i 3+4i; 5+6i 7+8i]"},
				{"%#v", "[\n 1+2i 3+4i\n 5+6i 7+8i\n]"},
			},
		},
		{
			m: CFormatted(NewCDense(4, 1, []complex128{1 + 2i, -3 + 4i, 5 + 6i, 7 + 8i}), FormatMATLAB()),
			rep: []rp{
				{"%v", "[1+2i; -3+4i; 5+6i; 7+8i]"},
				{"%#v", "[\n  1+2i\n -3+4i\n  5+6i\n  7+8i\n]"},
			},
		},
		{
			m: CFormatted(NewCDense(1, 4, []complex128{1 + 2i, -3 + 4i, 5 + 6i, 7 + 8i}), FormatMATLAB()),
			rep: []rp{
				{"%v", "[1+2i -3+4i 5+6i 7+8i]"},
				{"%#v", "[ 1+2i -3+4i  5+6i  7+8i]"},
			},
		},
		{
			m: CFormatted(NewCDense(2, 2, []complex128{1 + 2i, -3 + 4i, 5 + 6i, 7 + 8i}), FormatMATLAB()),
			rep: []rp{
				{"%v", "[1+2i -3+4i; 5+6i 7+8i]"},
				{"%#v", "[\n  1+2i -3+4i\n  5+6i  7+8i\n]"},
			},
		},
		{
			m: CFormatted(NewCDense(4, 1, []complex128{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i}), FormatPython()),
			rep: []rp{
				{"%v", "[[1+2j], [3+4j], [5+6j], [7+8j]]"},
				{"%#v", "[[1+2j],\n [3+4j],\n [5+6j],\n [7+8j]]"},
				{"%s", "%!s(*mat.CDense=Dims(4, 1))"},
				{"%#s", "%!s(*mat.CDense=Dims(4, 1))"},
			},
		},
		{
			m: CFormatted(NewCDense(1, 4, []complex128{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i}), FormatPython()),
			rep: []rp{
				{"%v", "[1+2j, 3+4j, 5+6j, 7+8j]"},
				{"%#v", "[1+2j, 3+4j, 5+6j, 7+8j]"},
			},
		},
		{
			m: CFormatted(NewCDense(2, 2, []complex128{1 + 2i, 3 + 4i, 5 + 6i, 7 + 8i}), FormatPython()),
			rep: []rp{
				{"%v", "[[1+2j, 3+4j], [5+6j, 7+8j]]"},
				{"%#v", "[[1+2j, 3+4j],\n [5+6j, 7+8j]]"},
			},
		},
		{
			m: CFormatted(NewCDense(4, 1, []complex128{1 + 2i, -3 + 4i, 5 + 6i, 7 + 8i}), FormatPython()),
			rep: []rp{
				{"%v", "[[1+2j], [-3+4j], [5+6j], [7+8j]]"},
				{"%#v", "[[ 1+2j],\n [-3+4j],\n [ 5+6j],\n [ 7+8j]]"},
			},
		},
		{
			m: CFormatted(NewCDense(1, 4, []complex128{1 + 2i, -3 + 4i, 5 + 6i, 7 + 8i}), FormatPython()),
			rep: []rp{
				{"%v", "[1+2j, -3+4j, 5+6j, 7+8j]"},
				{"%#v", "[ 1+2j, -3+4j,  5+6j,  7+8j]"},
			},
		},
		{
			m: CFormatted(NewCDense(2, 2, []complex128{1 + 2i, -3 + 4i, 5 + 6i, 7 + 8i}), FormatPython()),
			rep: []rp{
				{"%v", "[[1+2j, -3+4j], [5+6j, 7+8j]]"},
				{"%#v", "[[ 1+2j, -3+4j],\n [ 5+6j,  7+8j]]"},
			},
		},
	} {
		for j, rp := range test.rep {
			got := fmt.Sprintf(rp.format, test.m)
			if got != rp.output {
				t.Errorf("unexpected format result test %d part %d:\ngot:\n%s\nwant:\n%s", i, j, got, rp.output)
			}
		}
	}
}

func BenchmarkFormat(b *testing.B) {
	formats := []struct {
		name string
		form string
		fopt FormatOption
	}{
		{"General", "%v", nil},
		{"DotByte", "% f", DotByte('*')},
		{"Excerpt", "%v", Excerpt(3)},
		{"Prefix", "%v", Prefix("\t")},
		{"Squeeze", "%v", Squeeze()},
		{"MATLAB", "%v", FormatMATLAB()},
		{"MATLAB#", "%#v", FormatMATLAB()},
		{"Python", "%v", FormatPython()},
		{"Python#", "%#v", FormatPython()},
	}
	for i := 10; i <= 1000; i *= 10 {
		src := rand.NewSource(1)
		a, _ := randDense(i, 0.95, src)
		for _, j := range formats {
			b.Run(fmt.Sprintf("%d/%s", i, j.name), func(b *testing.B) {
				for k := 0; k < b.N; k++ {
					_ = fmt.Sprintf(j.form, Formatted(a), j.fopt)
				}
			})
		}
	}
}

func randCDense(size int, rho float64, src rand.Source) (*CDense, error) {
	if size == 0 {
		return nil, ErrZeroLength
	}
	a := &CDense{
		mat: cblas128.General{
			Rows: size, Cols: size, Stride: size,
			Data: make([]complex128, size*size),
		},
		capRows: size, capCols: size,
	}
	rnd := rand.New(src)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if rnd.Float64() < rho {
				a.set(i, j, complex(rnd.NormFloat64(), rnd.NormFloat64()))
			}
		}
	}
	return a, nil
}

func BenchmarkCFormat(b *testing.B) {
	formats := []struct {
		name string
		form string
		fopt FormatOption
	}{
		{"General", "%v", nil},
		{"DotByte", "% f", DotByte('*')},
		{"Excerpt", "%v", Excerpt(3)},
		{"Prefix", "%v", Prefix("\t")},
		{"Squeeze", "%v", Squeeze()},
		{"MATLAB", "%v", FormatMATLAB()},
		{"MATLAB#", "%#v", FormatMATLAB()},
		{"Python", "%v", FormatPython()},
		{"Python#", "%#v", FormatPython()},
	}
	for i := 10; i <= 1000; i *= 10 {
		src := rand.NewSource(1)
		a, _ := randCDense(i, 0.95, src)
		for _, j := range formats {
			b.Run(fmt.Sprintf("%d/%s", i, j.name), func(b *testing.B) {
				for k := 0; k < b.N; k++ {
					_ = fmt.Sprintf(j.form, CFormatted(a), j.fopt)
				}
			})
		}
	}
}
