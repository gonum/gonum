// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package curve

import (
	"errors"
	"fmt"
	"math/rand/v2"
	"reflect"
	"slices"
	"strings"
	"testing"
)

func ExampleHilbert2D_Pos() {
	h := Hilbert2D{order: 3}

	for y := 0; y < 1<<h.order; y++ {
		for x := 0; x < 1<<h.order; x++ {
			if x > 0 {
				fmt.Print("  ")
			}
			fmt.Printf("%02x", h.Pos([]int{x, y}))
		}
		fmt.Println()
	}
	// Output:
	// 00  01  0e  0f  10  13  14  15
	// 03  02  0d  0c  11  12  17  16
	// 04  07  08  0b  1e  1d  18  19
	// 05  06  09  0a  1f  1c  1b  1a
	// 3a  39  36  35  20  23  24  25
	// 3b  38  37  34  21  22  27  26
	// 3c  3d  32  33  2e  2d  28  29
	// 3f  3e  31  30  2f  2c  2b  2a
}

func ExampleHilbert3D_Pos() {
	h := Hilbert3D{order: 2}

	for z := 0; z < 1<<h.order; z++ {
		for y := 0; y < 1<<h.order; y++ {
			for x := 0; x < 1<<h.order; x++ {
				if x > 0 {
					fmt.Print("  ")
				}
				fmt.Printf("%02x", h.Pos([]int{x, y, z}))
			}
			fmt.Println()
		}
		fmt.Println()
	}
	// Output:
	// 00  07  08  0b
	// 01  06  0f  0c
	// 1a  1b  10  13
	// 19  18  17  14
	//
	// 03  04  09  0a
	// 02  05  0e  0d
	// 1d  1c  11  12
	// 1e  1f  16  15
	//
	// 3c  3b  36  35
	// 3d  3a  31  32
	// 22  23  2e  2d
	// 21  20  29  2a
	//
	// 3f  38  37  34
	// 3e  39  30  33
	// 25  24  2f  2c
	// 26  27  28  2b
}

func ExampleHilbert4D_Pos() {
	h := Hilbert4D{order: 2}
	N := 1 << h.order
	for z := 0; z < N; z++ {
		if z > 0 {
			s := strings.Repeat("═", N*4-2)
			s = s + strings.Repeat("═╬═"+s, N-1)
			fmt.Println(s)
		}
		for y := 0; y < N; y++ {
			for w := 0; w < N; w++ {
				if w > 0 {
					fmt.Print(" ║ ")
				}
				for x := 0; x < N; x++ {
					if x > 0 {
						fmt.Print("  ")
					}
					fmt.Printf("%02x", h.Pos([]int{x, y, z, w}))
				}
			}
			fmt.Println()
		}
	}
	// Output:
	// 00  0f  10  13 ║ 03  0c  11  12 ║ fc  f3  ee  ed ║ ff  f0  ef  ec
	// 01  0e  1f  1c ║ 02  0d  1e  1d ║ fd  f2  e1  e2 ║ fe  f1  e0  e3
	// 32  31  20  23 ║ 35  36  21  22 ║ ca  c9  de  dd ║ cd  ce  df  dc
	// 33  30  2f  2c ║ 34  37  2e  2d ║ cb  c8  d1  d2 ║ cc  cf  d0  d3
	// ═══════════════╬════════════════╬════════════════╬═══════════════
	// 07  08  17  14 ║ 04  0b  16  15 ║ fb  f4  e9  ea ║ f8  f7  e8  eb
	// 06  09  18  1b ║ 05  0a  19  1a ║ fa  f5  e6  e5 ║ f9  f6  e7  e4
	// 3d  3e  27  24 ║ 3a  39  26  25 ║ c5  c6  d9  da ║ c2  c1  d8  db
	// 3c  3f  28  2b ║ 3b  38  29  2a ║ c4  c7  d6  d5 ║ c3  c0  d7  d4
	// ═══════════════╬════════════════╬════════════════╬═══════════════
	// 76  77  6c  6d ║ 79  78  6b  6a ║ 86  87  94  95 ║ 89  88  93  92
	// 75  74  63  62 ║ 7a  7b  64  65 ║ 85  84  9b  9a ║ 8a  8b  9c  9d
	// 42  41  5c  5d ║ 45  46  5b  5a ║ ba  b9  a4  a5 ║ bd  be  a3  a2
	// 43  40  53  52 ║ 44  47  54  55 ║ bb  b8  ab  aa ║ bc  bf  ac  ad
	// ═══════════════╬════════════════╬════════════════╬═══════════════
	// 71  70  6f  6e ║ 7e  7f  68  69 ║ 81  80  97  96 ║ 8e  8f  90  91
	// 72  73  60  61 ║ 7d  7c  67  66 ║ 82  83  98  99 ║ 8d  8c  9f  9e
	// 4d  4e  5f  5e ║ 4a  49  58  59 ║ b5  b6  a7  a6 ║ b2  b1  a0  a1
	// 4c  4f  50  51 ║ 4b  48  57  56 ║ b4  b7  a8  a9 ║ b3  b0  af  ae
}

func TestConstructors(t *testing.T) {
	const intSize = 32 << (^uint(0) >> 63) // 32 or 64

	t.Run("2D/Ok", func(t *testing.T) {
		_, err := NewHilbert2D(intSize/2 - 1)
		noError(t, err)
	})

	t.Run("3D/Ok", func(t *testing.T) {
		_, err := NewHilbert3D(intSize / 3)
		noError(t, err)
	})

	t.Run("4D/Ok", func(t *testing.T) {
		_, err := NewHilbert4D(intSize/4 - 1)
		noError(t, err)
	})

	t.Run("2D/Underflow", func(t *testing.T) {
		_, err := NewHilbert2D(0)
		errorIs(t, err, ErrUnderflow)
	})

	t.Run("3D/Underflow", func(t *testing.T) {
		_, err := NewHilbert3D(0)
		errorIs(t, err, ErrUnderflow)
	})

	t.Run("4D/Underflow", func(t *testing.T) {
		_, err := NewHilbert4D(0)
		errorIs(t, err, ErrUnderflow)
	})

	t.Run("2D/Overflow", func(t *testing.T) {
		_, err := NewHilbert2D(intSize / 2)
		errorIs(t, err, ErrOverflow)
	})

	t.Run("3D/Overflow", func(t *testing.T) {
		_, err := NewHilbert3D(intSize/3 + 1)
		errorIs(t, err, ErrOverflow)
	})

	t.Run("4D/Overflow", func(t *testing.T) {
		_, err := NewHilbert4D(intSize / 4)
		errorIs(t, err, ErrOverflow)
	})
}

func TestHilbert2D(t *testing.T) {
	for ord := 1; ord <= 4; ord++ {
		t.Run(fmt.Sprintf("Order/%d", ord), func(t *testing.T) {
			testCurve(t, Hilbert2D{order: ord})
		})
	}

	cases := map[int][]int{
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
		t.Run(fmt.Sprintf("Case/%d", order), func(t *testing.T) {
			testCurveCase(t, Hilbert2D{order: order}, order, expected)
		})
	}
}

func TestHilbert3D(t *testing.T) {
	for ord := 1; ord <= 4; ord++ {
		t.Run(fmt.Sprintf("Order/%d", ord), func(t *testing.T) {
			testCurve(t, Hilbert3D{order: ord})
		})
	}

	cases := map[int][]int{
		1: {
			0, 1,
			3, 2,

			7, 6,
			4, 5,
		},
		2: {
			0x00, 0x07, 0x08, 0x0B,
			0x01, 0x06, 0x0F, 0x0C,
			0x1A, 0x1B, 0x10, 0x13,
			0x19, 0x18, 0x17, 0x14,

			0x03, 0x04, 0x09, 0x0A,
			0x02, 0x05, 0x0E, 0x0D,
			0x1D, 0x1C, 0x11, 0x12,
			0x1E, 0x1F, 0x16, 0x15,

			0x3C, 0x3B, 0x36, 0x35,
			0x3D, 0x3A, 0x31, 0x32,
			0x22, 0x23, 0x2E, 0x2D,
			0x21, 0x20, 0x29, 0x2A,

			0x3F, 0x38, 0x37, 0x34,
			0x3E, 0x39, 0x30, 0x33,
			0x25, 0x24, 0x2F, 0x2C,
			0x26, 0x27, 0x28, 0x2B,
		},
	}

	for order, expected := range cases {
		t.Run(fmt.Sprintf("Case/%d", order), func(t *testing.T) {
			testCurveCase(t, Hilbert3D{order: order}, order, expected)
		})
	}
}

func TestHilbert4D(t *testing.T) {
	for ord := 1; ord <= 4; ord++ {
		t.Run(fmt.Sprintf("Order %d", ord), func(t *testing.T) {
			testCurve(t, Hilbert4D{order: ord})
		})
	}
}

func BenchmarkHilbert(b *testing.B) {
	const O = 10
	for N := 2; N <= 4; N++ {
		b.Run(fmt.Sprintf("%dD/Pos", N), func(b *testing.B) {
			for ord := 1; ord <= O; ord++ {
				b.Run(fmt.Sprintf("Order %d", ord), func(b *testing.B) {
					h := newCurve(ord, N)
					v := make([]int, N)
					for i := range v {
						v[i] = rand.IntN(1 << ord)
					}
					u := make([]int, N)
					for n := 0; n < b.N; n++ {
						copy(u, v)
						h.Pos(u)
					}
				})
			}
		})

		b.Run(fmt.Sprintf("%dD/Coord", N), func(b *testing.B) {
			for ord := 1; ord <= O; ord++ {
				b.Run(fmt.Sprintf("Order %d", ord), func(b *testing.B) {
					h := newCurve(ord, N)
					d := rand.IntN(1 << (ord * N))
					v := make([]int, N)
					for n := 0; n < b.N; n++ {
						h.Coord(v, d)
					}
				})
			}
		})
	}
}

type curve interface {
	Dims() []int
	Len() int
	Pos(v []int) int
	Coord(dst []int, pos int) []int
}

func newCurve(order, dim int) curve {
	switch dim {
	case 2:
		return Hilbert2D{order: order}
	case 3:
		return Hilbert3D{order: order}
	case 4:
		return Hilbert4D{order: order}
	}
	panic("invalid dimension")
}

// testCurve verifies that Pos and Coord (of C) are inverses of each other and
// that the spatial coordinates V and U - corresponding to linear coordinates D
// and D+1 - are exactly one unit (euclidean) distant from each other.
func testCurve(t *testing.T, c curve) {
	t.Helper()

	// Stop if the error count reaches 10
	var errc int
	fail := func() {
		if errc < 10 {
			errc++
			t.Fail()
		} else {
			t.FailNow()
		}
	}

	// Map between linear and spatial coordinates, and verify that Pos and Coord
	// are inverses of each other
	m := map[int][]int{}
	curveRange(c, func(v []int) {
		d := c.Pos(slices.Clone(v))
		u := c.Coord(nil, d)
		if !reflect.DeepEqual(v, u) {
			t.Logf("Space is not the inverse of Curve for d=%d %v", d, v)
			fail()
		}

		m[d] = slices.Clone(v)
	})

	D := 1
	for _, v := range c.Dims() {
		D *= v
	}

	// For each possible pairs of linear coordinates D and D+1, verify that the
	// corresponding spatial coordinates V and U are exactly one unit apart
	// (euclidean distance)
	for d := 0; d < D-1; d++ {
		v, u := m[d], m[d+1]
		if !adjacent(v, u) {
			t.Logf("points %x and %x are not adjacent\n    %v -> %v", d, d+1, v, u)
			fail()
		}
	}
}

// curveRange ranges over the n-dimensional coordinate space of the curve,
// calling fn on each element of the space.
func curveRange(c curve, fn func([]int)) {
	size := c.Dims()
	dimRange(len(size), size, make([]int, len(size)), fn)
}

// dimRange ranges over the coordinate space defined by size, calling fn on each
// element of the space. Call dimRange with dim = len(size).
func dimRange(dim int, size []int, v []int, fn func([]int)) {
	if dim == 0 {
		fn(v)
		return
	}

	for i := 0; i < size[dim-1]; i++ {
		v[dim-1] = i
		dimRange(dim-1, size, v, fn)
	}
}

// adjacent returns true if the euclidean distance between v and u is
// exactly one. v and u must be the same length.
//
// In other words, given d(i) = abs(v[i] - u[i]), adjacent returns false if
// d(i) > 1 for any i or if d(i) == 1 for more than a single i, or if d(i)
// == 0 for all i.
func adjacent(v, u []int) bool {
	n := 0
	for i := range v {
		x := v[i] - u[i]
		if x == 0 {
			continue
		}
		if x < -1 || 1 < x {
			return false
		}
		n++
	}

	return n == 1
}

// testCurveCase verifies that the curve produces the expected sequence of
// values.
func testCurveCase(t *testing.T, c curve, order int, expected []int) {
	t.Helper()

	dim := len(c.Dims())
	actual := make([]int, len(expected))
	for i := range expected {
		v := coord(i, order, dim)
		actual[i] = c.Pos(slices.Clone(v))
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("unexpected result:\ngot:  %v\nwant: %v", actual, expected)
	}

	for i, d := range expected {
		v := coord(i, order, dim)
		if !reflect.DeepEqual(v, c.Coord(nil, d)) {
			t.Fatalf("[%d] expected %v for d = %d", i, v, d)
		}
	}
}

// coord returns the nth spatial coordinates for a dim-dimensional space with
// sides 2^ord in row-major order.
func coord(n, ord, dim int) []int {
	v := make([]int, dim)
	for i := 0; i < dim; i++ {
		v[i] = n % (1 << ord)
		n /= (1 << ord)
	}
	return v
}

func noError(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}

func errorIs(t testing.TB, err, target error) {
	t.Helper()
	if err == nil {
		t.Fatal("Expected an error")
	}
	if !errors.Is(err, target) {
		t.Fatalf("Expected %v to be %v", err, target)
	}
}
