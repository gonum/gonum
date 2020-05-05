// Copyright ©2013 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import (
	"fmt"
	"math"
	"testing"
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
