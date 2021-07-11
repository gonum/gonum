package curve

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type Dim int

func (d Dim) Run(ord int, fn func([]int)) {
	d.run(ord, make([]int, d), fn)
}

func (d Dim) run(ord int, v []int, fn func([]int)) {
	if d == 0 {
		fn(v)
		return
	}

	for i := 0; i < ord; i++ {
		v[d-1] = i
		(d - 1).run(ord, v, fn)
	}
}

func vec(n, ord, dim int) []int {
	v := make([]int, dim)
	for i := 0; i < dim; i++ {
		v[i] = n % (1 << ord)
		n /= (1 << ord)
	}
	return v
}

func BenchmarkHilbert2D(b *testing.B) {
	b.Run("Curve", func(b *testing.B) {
		for ord := 1; ord <= 10; ord++ {
			b.Run(fmt.Sprint(ord), func(b *testing.B) {
				h := Hilbert2D{Order: ord}
				x, y := rand.Intn(ord), rand.Intn(ord)
				for n := 0; n < b.N; n++ {
					h.Curve(x, y)
				}
			})
		}
	})

	b.Run("Space", func(b *testing.B) {
		for ord := 1; ord <= 10; ord++ {
			b.Run(fmt.Sprint(ord), func(b *testing.B) {
				h := Hilbert2D{Order: ord}
				d := Point(rand.Intn(1 << ord))
				for n := 0; n < b.N; n++ {
					h.Space2D(d)
				}
			})
		}
	})
}

func TestHilbert2D(t *testing.T) {
	for ord := 1; ord <= 10; ord++ {
		t.Run(fmt.Sprint(ord), func(t *testing.T) {
			h := Hilbert2D{Order: ord}
			Dim(2).Run(ord, func(v []int) {
				x := make([]int, 2)
				copy(x, v)
				d := h.Curve(v...)
				u := h.Space(d)
				assert.Equal(t, x, u)
			})
		})
	}

	cases := map[int][]Point{
		1: {
			0, 1,
			3, 2,
		},
		2: {
			0x0, 0x3, 0x4, 0x5,
			0x1, 0x2, 0x7, 0x6,
			0xE, 0xD, 0x8, 0x9,
			0xF, 0xC, 0xB, 0xA,
		},
		3: {
			0x00, 0x01, 0x0E, 0x0F, 0x10, 0x13, 0x14, 0x15,
			0x03, 0x02, 0x0D, 0x0C, 0x11, 0x12, 0x17, 0x16,
			0x04, 0x07, 0x08, 0x0B, 0x1E, 0x1D, 0x18, 0x19,
			0x05, 0x06, 0x09, 0x0A, 0x1F, 0x1C, 0x1B, 0x1A,
			0x3A, 0x39, 0x36, 0x35, 0x20, 0x23, 0x24, 0x25,
			0x3B, 0x38, 0x37, 0x34, 0x21, 0x22, 0x27, 0x26,
			0x3C, 0x3D, 0x32, 0x33, 0x2E, 0x2D, 0x28, 0x29,
			0x3F, 0x3E, 0x31, 0x30, 0x2F, 0x2C, 0x2B, 0x2A,
		},
	}

	for order, expected := range cases {
		t.Run(fmt.Sprintf("Order %d", order), func(t *testing.T) {
			h := Hilbert2D{Order: order}

			actual := make([]Point, len(expected))
			for i := range expected {
				v := vec(i, order, 2)
				actual[i] = h.Curve(v...)
			}
			require.Equal(t, expected, actual, "expected curves to equal")

			for i, v := range expected {
				u := vec(i, order, 2)
				require.Equal(t, u, h.Space(v), "[%d] expected (%d, %d) for d = %d", i, u[0], u[1], v)
			}
		})
	}
}

func TestHilbert3D(t *testing.T) {
	for ord := 1; ord <= 10; ord++ {
		t.Run(fmt.Sprint(ord), func(t *testing.T) {
			h := Hilbert3D{Order: ord}
			Dim(3).Run(ord, func(v []int) {
				x := make([]int, 3)
				copy(x, v)
				d := h.Curve(v...)
				u := h.Space(d)
				assert.Equal(t, x, u)
			})
		})
	}

	cases := map[int][]Point{
		1: {
			0, 1,
			3, 2,

			7, 6,
			4, 5,
		},
		2: {
			0x00, 0x07, 0x08, 0x09,
			0x03, 0x04, 0x0F, 0x0E,
			0x1A, 0x1B, 0x10, 0x11,
			0x19, 0x18, 0x17, 0x16,

			0x01, 0x06, 0x0B, 0x0A,
			0x02, 0x05, 0x0C, 0x0D,
			0x1D, 0x1C, 0x13, 0x12,
			0x1E, 0x1F, 0x14, 0x15,

			0x3E, 0x39, 0x34, 0x35,
			0x3D, 0x3A, 0x33, 0x32,
			0x22, 0x23, 0x2C, 0x2D,
			0x21, 0x20, 0x2B, 0x2A,

			0x3F, 0x38, 0x37, 0x36,
			0x3C, 0x3B, 0x30, 0x31,
			0x25, 0x24, 0x2F, 0x2E,
			0x26, 0x27, 0x28, 0x29,
		},
	}

	for order, expected := range cases {
		t.Run(fmt.Sprintf("Order %d", order), func(t *testing.T) {
			h := Hilbert3D{Order: order}

			actual := make([]Point, len(expected))
			for i := range expected {
				v := vec(i, order, 3)
				actual[i] = h.Curve(v...)
			}
			require.Equal(t, expected, actual, "expected curves to equal")

			for i, v := range expected {
				u := vec(i, order, 3)
				require.Equal(t, u, h.Space(v), "[%d] expected (%d, %d, %d) for d = %d", i, u[0], u[1], u[2], v)
			}
		})
	}
}

func TestHilbert(t *testing.T) {
	for dim := 2; dim <= 3; dim++ {
		t.Run(fmt.Sprint(dim), func(t *testing.T) {
			for ord := 1; ord <= 10; ord++ {
				t.Run(fmt.Sprint(ord), func(t *testing.T) {
					h := Hilbert{Order: ord, Dimension: dim}
					Dim(dim).Run(ord, func(v []int) {
						x := make([]int, dim)
						copy(x, v)
						d := h.Curve(v...)
						u := h.Space(d)
						assert.Equal(t, x, u)
					})
				})
			}
		})
	}
}

func ExampleHilbert2D_Curve() {
	h := Hilbert2D{Order: 3}

	for y := 0; y < 1<<h.Order; y++ {
		for x := 0; x < 1<<h.Order; x++ {
			if x > 0 {
				fmt.Print("  ")
			}
			fmt.Printf("%02X", h.Curve(x, y))
		}
		fmt.Println()
	}

	// Output:
	// 00  01  0E  0F  10  13  14  15
	// 03  02  0D  0C  11  12  17  16
	// 04  07  08  0B  1E  1D  18  19
	// 05  06  09  0A  1F  1C  1B  1A
	// 3A  39  36  35  20  23  24  25
	// 3B  38  37  34  21  22  27  26
	// 3C  3D  32  33  2E  2D  28  29
	// 3F  3E  31  30  2F  2C  2B  2A
}

func ExampleHilbert3D_Curve() {
	h := Hilbert3D{Order: 1}

	for z := 0; z < 1<<h.Order; z++ {
		for y := 0; y < 1<<h.Order; y++ {
			for x := 0; x < 1<<h.Order; x++ {
				if x > 0 {
					fmt.Print("  ")
				}
				fmt.Printf("%02X", h.Curve(x, y, z))
			}
			fmt.Println()
		}
		fmt.Println()
	}

	// Output:
	// 00  07  08  09
	// 03  04  0F  0E
	// 1A  1B  10  11
	// 19  18  17  16
	//
	// 01  06  0B  0A
	// 02  05  0C  0D
	// 1D  1C  13  12
	// 1E  1F  14  15
	//
	// 3E  39  34  35
	// 3D  3A  33  32
	// 22  23  2C  2D
	// 21  20  2B  2A
	//
	// 3F  38  37  36
	// 3C  3B  30  31
	// 25  24  2F  2E
	// 26  27  28  29
}
