package mat

import (
	"fmt"
	"testing"
)

func TestComplexFormat(t *testing.T) {
	t.Parallel()
	type rp struct {
		format string
		output string
	}
	for i, test := range []struct {
		m   fmt.Formatter
		rep []rp
	}{
		{
			m: CFormatted(NewCDense(3, 3, []complex128{0, 0, 0, 0, 0, 0, 0, 0, 0})),
			rep: []rp{
				{"%v", "⎡0i  0i  0i⎤\n⎢0i  0i  0i⎥\n⎣0i  0i  0i⎦"},
				{"% v", "⎡ .   .   .⎤\n⎢ .   .   .⎥\n⎣ .   .   .⎦"},
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
			m: CFormatted(NewCDense(3, 3, []complex128{1, 1, 1, 1, 1, 1, 1, 1, 1}), CPrefix("\t")),
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
			m: CFormatted(NewCDense(3, 3, []complex128{
				1.2345678 + 2.34413124231i, 3.321421 + 5.5231521i, 6.42314231 + 3.321321i,
				1.2345678 + 2.34413124231i, 3.321421 + 5.5231521i, 6.42314231 + 3.321321i,
				1.2345678 + 2.34413124231i, 3.321421 + 5.5231521i, 6.42314231 + 3.321321i,
			})),
			rep: []rp{
				{"%v", "⎡1.2345678+2.34413124231i       3.321421+5.5231521i      6.42314231+3.321321i⎤\n⎢1.2345678+2.34413124231i       3.321421+5.5231521i      6.42314231+3.321321i⎥\n⎣1.2345678+2.34413124231i       3.321421+5.5231521i      6.42314231+3.321321i⎦"},
				{"%.2f", "⎡1.23+2.34i  3.32+5.52i  6.42+3.32i⎤\n⎢1.23+2.34i  3.32+5.52i  6.42+3.32i⎥\n⎣1.23+2.34i  3.32+5.52i  6.42+3.32i⎦"},
				{"% f", "⎡1.2345678+2.34413124231i       3.321421+5.5231521i      6.42314231+3.321321i⎤\n⎢1.2345678+2.34413124231i       3.321421+5.5231521i      6.42314231+3.321321i⎥\n⎣1.2345678+2.34413124231i       3.321421+5.5231521i      6.42314231+3.321321i⎦"},
				{"%#v", fmt.Sprintf("%#v", NewCDense(3, 3, []complex128{
					1.2345678 + 2.34413124231i, 3.321421 + 5.5231521i, 6.42314231 + 3.321321i,
					1.2345678 + 2.34413124231i, 3.321421 + 5.5231521i, 6.42314231 + 3.321321i,
					1.2345678 + 2.34413124231i, 3.321421 + 5.5231521i, 6.42314231 + 3.321321i,
				}))},
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
