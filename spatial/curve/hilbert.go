// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package curve

import (
	"errors"
	"fmt"
)

// ErrUnderflow is returned by Hilbert curve constructors when the power is less
// than 1.
var ErrUnderflow = errors.New("order is less than 1")

// ErrOverflow is returned (wrapped) by Hilbert curve constructors when the
// power would cause Len and Pos to overflow.
var ErrOverflow = errors.New("overflow int")

// The size of an int. Taken from src/math/const.go.
const intSize = 32 << (^uint(0) >> 63) // 32 or 64

// Hilbert2D is a 2-dimensional Hilbert curve.
type Hilbert2D struct{ order int }

// NewHilbert2D constructs a [Hilbert2D] of the given order. NewHilbert2D
// returns [ErrOverflow] (wrapped) if the order would cause Len and Pos to
// overflow.
func NewHilbert2D(order int) (Hilbert2D, error) {
	v := Hilbert2D{order: order}

	// The order must be greater than zero.
	if order < 1 {
		return v, ErrUnderflow
	}

	// The product of the order and number of dimensions must not exceed or
	// equal the size of an int.
	if order*2 >= intSize {
		return v, fmt.Errorf("a 2-dimensional, %d-order Hilbert curve will %w", order, ErrOverflow)
	}

	return v, nil
}

// Dims returns the spatial dimensions of the curve, which is {2ᵏ, 2ᵏ}, where k
// is the order.
func (h Hilbert2D) Dims() []int { return []int{1 << h.order, 1 << h.order} }

// Len returns the length of the curve, which is 2ⁿᵏ, where n is the dimension
// (2) and k is the order.
//
// Len will overflow if order is ≥ 16 on architectures where [int] is 32-bits,
// or ≥ 32 on architectures where [int] is 64-bits.
func (h Hilbert2D) Len() int { return 1 << (2 * h.order) }

func (h Hilbert2D) rot(n int, v []int, d int) {
	switch d {
	case 0:
		swap{0, 1}.do(n, v)
	case 3:
		flip{0, 1}.do(n, v)
	}
}

// Pos returns the linear position of the 3-spatial coordinate v along the
// curve. Pos modifies v.
//
// Pos will overflow if order is ≥ 16 on architectures where [int] is 32-bits,
// or ≥ 32 on architectures where [int] is 64-bits.
func (h Hilbert2D) Pos(v []int) int {
	var d int
	const dims = 2
	for n := h.order - 1; n >= 0; n-- {
		rx := (v[0] >> n) & 1
		ry := (v[1] >> n) & 1
		rd := ry<<1 | (ry ^ rx)
		d += rd << (dims * n)
		h.rot(h.order, v, rd)
	}
	return d
}

// Coord writes the spatial coordinates of pos to dst and returns it. If dst is
// nil, Coord allocates a new slice of length 2; otherwise Coord is
// allocation-free.
//
// Coord panics if dst is not nil and len(dst) is not equal to 2.
func (h Hilbert2D) Coord(dst []int, pos int) []int {
	if dst == nil {
		dst = make([]int, 2)
	} else if len(dst) != 2 {
		panic("len(dst) must equal 2")
	}
	for n := 0; n < h.order; n++ {
		e := pos & 3
		h.rot(n, dst[:], e)

		ry := e >> 1
		rx := (e>>0 ^ e>>1) & 1
		dst[0] += rx << n
		dst[1] += ry << n
		pos >>= 2
	}
	return dst
}

// Hilbert3D is a 3-dimensional Hilbert curve.
type Hilbert3D struct{ order int }

// NewHilbert3D constructs a [Hilbert3D] of the given order. NewHilbert3D
// returns [ErrOverflow] (wrapped) if the order would cause Len and Pos to
// overflow.
func NewHilbert3D(order int) (Hilbert3D, error) {
	v := Hilbert3D{order: order}

	// The order must be greater than zero.
	if order < 1 {
		return v, ErrUnderflow
	}

	// The product of the order and number of dimensions must not exceed or
	// equal the size of an int.
	if order*3 >= intSize {
		return v, fmt.Errorf("a 3-dimensional, %d-order Hilbert curve will %w", order, ErrOverflow)
	}

	return v, nil
}

// Dims returns the spatial dimensions of the curve, which is {2ᵏ, 2ᵏ, 2ᵏ}, where
// k is the order.
func (h Hilbert3D) Dims() []int { return []int{1 << h.order, 1 << h.order, 1 << h.order} }

// Len returns the length of the curve, which is 2ⁿᵏ, where n is the dimension
// (3) and k is the order.
//
// Len will overflow if order is ≥ 11 on architectures where [int] is 32-bits,
// or ≥ 21 on architectures where [int] is 64-bits.
func (h Hilbert3D) Len() int { return 1 << (3 * h.order) }

func (h Hilbert3D) rot(reverse bool, n int, v []int, d int) {
	switch d {
	case 0:
		do2(reverse, n, v, swap{1, 2}, swap{0, 2})
	case 1, 2:
		do2(reverse, n, v, swap{0, 2}, swap{1, 2})
	case 3, 4:
		invert{0, 1}.do(n, v)
	case 5, 6:
		do2(reverse, n, v, flip{0, 2}, flip{1, 2})
	case 7:
		do2(reverse, n, v, flip{1, 2}, flip{0, 2})
	}
}

// Pos returns the linear position of the 4-spatial coordinate v along the
// curve. Pos modifies v.
//
// Pos will overflow if order is ≥ 11 on architectures where [int] is 32-bits,
// or ≥ 21 on architectures where [int] is 64-bits.
func (h Hilbert3D) Pos(v []int) int {
	var d int
	const dims = 3
	for n := h.order - 1; n >= 0; n-- {
		rx := (v[0] >> n) & 1
		ry := (v[1] >> n) & 1
		rz := (v[2] >> n) & 1
		rd := rz<<2 | (rz^ry)<<1 | (rz ^ ry ^ rx)
		d += rd << (dims * n)
		h.rot(false, h.order, v, rd)
	}
	return d
}

// Coord writes the spatial coordinates of pos to dst and returns it. If dst is
// nil, Coord allocates a new slice of length 3; otherwise Coord is
// allocation-free.
//
// Coord panics if dst is not nil and len(dst) is not equal to 3.
func (h Hilbert3D) Coord(dst []int, pos int) []int {
	if dst == nil {
		dst = make([]int, 3)
	} else if len(dst) != 3 {
		panic("len(dst) must equal 3")
	}
	for n := 0; n < h.order; n++ {
		e := pos & 7
		h.rot(true, n, dst[:], e)

		rz := e >> 2
		ry := (e>>1 ^ e>>2) & 1
		rx := (e>>0 ^ e>>1) & 1
		dst[0] += rx << n
		dst[1] += ry << n
		dst[2] += rz << n
		pos >>= 3
	}
	return dst
}

// Hilbert4D is a 4-dimensional Hilbert curve.
type Hilbert4D struct{ order int }

// NewHilbert4D constructs a [Hilbert4D] of the given order. NewHilbert4D
// returns [ErrOverflow] (wrapped) if the order would cause Len and Pos to
// overflow.
func NewHilbert4D(order int) (Hilbert4D, error) {
	v := Hilbert4D{order: order}

	// The order must be greater than zero.
	if order < 1 {
		return v, ErrUnderflow
	}

	// The product of the order and number of dimensions must not exceed or
	// equal the size of an int.
	if order*4 >= intSize {
		return v, fmt.Errorf("a 4-dimensional, %d-order Hilbert curve will %w", order, ErrOverflow)
	}

	return v, nil
}

// Dims returns the spatial dimensions of the curve, which is {2ᵏ, 2ᵏ, 2ᵏ, 2ᵏ},
// where k is the order.
func (h Hilbert4D) Dims() []int { return []int{1 << h.order, 1 << h.order, 1 << h.order, 1 << h.order} }

// Len returns the length of the curve, which is 2ⁿᵏ, where n is the dimension
// (4) and k is the order.
//
// Len will overflow if order is ≥ 8 on architectures where [int] is 32-bits, or
// ≥ 16 on architectures where [int] is 64-bits.
func (h Hilbert4D) Len() int { return 1 << (4 * h.order) }

func (h Hilbert4D) rot(reverse bool, n int, v []int, d int) {
	switch d {
	case 0:
		do2(reverse, n, v, swap{1, 3}, swap{0, 3})
	case 1, 2:
		do2(reverse, n, v, swap{0, 3}, swap{1, 3})
	case 3, 4:
		do2(reverse, n, v, flip{0, 1}, swap{2, 3})
	case 5, 6:
		do2(reverse, n, v, flip{1, 2}, swap{2, 3})
	case 7, 8:
		invert{0, 2}.do(n, v)
	case 9, 10:
		do2(reverse, n, v, flip{1, 2}, flip{2, 3})
	case 11, 12:
		do2(reverse, n, v, flip{0, 1}, flip{2, 3})
	case 13, 14:
		do2(reverse, n, v, flip{0, 3}, flip{1, 3})
	case 15:
		do2(reverse, n, v, flip{1, 3}, flip{0, 3})
	}
}

// Pos returns the linear position of the 2-spatial coordinate v along the
// curve. Pos modifies v.
//
// Pos will overflow if order is ≥ 8 on architectures where [int] is 32-bits, or
// ≥ 16 on architectures where [int] is 64-bits.
func (h Hilbert4D) Pos(v []int) int {
	var d int
	const dims = 4
	for n := h.order - 1; n >= 0; n-- {
		var e int
		for i := dims - 1; i >= 0; i-- {
			v := v[i] >> n & 1
			e = e<<1 | (e^v)&1
		}

		d += e << (dims * n)
		h.rot(false, h.order, v, e)
	}
	return d
}

// Coord writes the spatial coordinates of pos to dst and returns it. If dst is
// nil, Coord allocates a new slice of length 4; otherwise Coord is
// allocation-free.
//
// Coord panics if dst is not nil and len(dst) is not equal to 4.
func (h Hilbert4D) Coord(dst []int, pos int) []int {
	if dst == nil {
		dst = make([]int, 4)
	} else if len(dst) != 4 {
		panic("len(dst) must equal 4")
	}
	N := 4
	for n := 0; n < h.order; n++ {
		e := pos & (1<<N - 1)
		h.rot(true, n, dst[:], e)

		for i, e := 0, e; i < N; i++ {
			dst[i] += (e ^ e>>1) & 1 << n
			e >>= 1
		}
		pos >>= N
	}
	return dst
}

type op interface{ do(int, []int) }

// invert I and J
type invert struct{ i, j int }

func (c invert) do(n int, v []int) { v[c.i], v[c.j] = v[c.i]^(1<<n-1), v[c.j]^(1<<n-1) }

// swap I and J
type swap struct{ i, j int }

func (c swap) do(n int, v []int) { v[c.i], v[c.j] = v[c.j], v[c.i] }

// swap and invert I and J
type flip struct{ i, j int }

func (c flip) do(n int, v []int) { v[c.i], v[c.j] = v[c.j]^(1<<n-1), v[c.i]^(1<<n-1) }

// do2 executes the given operations, optionally in reverse.
//
// Generic specialization reduces allocation (because it can eliminate interface
// value boxing) and improves performance
func do2[A, B op](reverse bool, n int, v []int, a A, b B) {
	if reverse {
		b.do(n, v)
		a.do(n, v)
	} else {
		a.do(n, v)
		b.do(n, v)
	}
}
