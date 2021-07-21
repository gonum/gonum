// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package curve

// Hilbert2D is a 2-dimensional Hilbert curve.
type Hilbert2D struct{ Order int }

// Size returns {2ⁿ, 2ⁿ} where n is the order.
func (h Hilbert2D) Size() []int { return []int{1 << h.Order, 1 << h.Order} }

func (h Hilbert2D) rot(n int, v []int, d int) {
	switch d {
	case 0:
		swap{0, 1}.do(n, v)
	case 3:
		flip{0, 1}.do(n, v)
	}
}

// Curve returns the curve coordinate of V. For order ≥2, Curve modifies V.
func (h Hilbert2D) Curve(v []int) uint64 {
	var d uint64
	for n := h.Order - 1; n >= 0; n-- {
		rx := (v[0] >> n) & 1
		ry := (v[1] >> n) & 1
		rd := ry<<1 | (ry ^ rx)
		d += uint64(rd) << (2 * n)
		h.rot(h.Order, v, rd)
	}
	return d
}

// Space2D returns the spatial coordinates of D.
func (h Hilbert2D) Space2D(d uint64) [2]int {
	var v [2]int
	for n := 0; n < h.Order; n++ {
		e := int(d & 3)
		h.rot(n, v[:], e)

		ry := e >> 1
		rx := (e>>0 ^ e>>1) & 1
		v[0] += rx << n
		v[1] += ry << n
		d >>= 2
	}
	return v
}

// Space returns Space2D as a slice.
func (h Hilbert2D) Space(d uint64) []int {
	xy := h.Space2D(d)
	return xy[:]
}

// Hilbert3D is a 3-dimensional Hilbert curve.
type Hilbert3D struct{ Order int }

// Size returns {2ⁿ, 2ⁿ, 2ⁿ} where n is the order.
func (h Hilbert3D) Size() []int { return []int{1 << h.Order, 1 << h.Order, 1 << h.Order} }

func (h Hilbert3D) rot(reverse bool, n int, v []int, d int) {
	switch d {
	case 0:
		ops{swap{1, 2}, swap{0, 2}}.do(reverse, n, v)
	case 1, 2:
		ops{swap{0, 2}, swap{1, 2}}.do(reverse, n, v)
	case 3, 4:
		invert{0, 1}.do(n, v)
	case 5, 6:
		ops{flip{0, 2}, flip{1, 2}}.do(reverse, n, v)
	case 7:
		ops{flip{1, 2}, flip{0, 2}}.do(reverse, n, v)
	}
}

// Curve returns the curve coordinate of V. For order ≥2, Curve modifies V.
func (h Hilbert3D) Curve(v []int) uint64 {
	var d uint64
	for n := h.Order - 1; n >= 0; n-- {
		rx := (v[0] >> n) & 1
		ry := (v[1] >> n) & 1
		rz := (v[2] >> n) & 1
		rd := rz<<2 | (rz^ry)<<1 | (rz ^ ry ^ rx)
		d += uint64(rd) << (3 * n)
		h.rot(false, h.Order, v, rd)
	}
	return d
}

// Space3D returns the spatial coordinates of D.
func (h Hilbert3D) Space3D(d uint64) [3]int {
	var v [3]int
	for n := 0; n < h.Order; n++ {
		e := int(d & 7)
		h.rot(true, n, v[:], e)

		rz := e >> 2
		ry := (e>>1 ^ e>>2) & 1
		rx := (e>>0 ^ e>>1) & 1
		v[0] += rx << n
		v[1] += ry << n
		v[2] += rz << n
		d >>= 3
	}
	return v
}

// Space returns Space3D as a slice.
func (h Hilbert3D) Space(v uint64) []int {
	xy := h.Space3D(v)
	return xy[:]
}

// Hilbert4D is a 4-dimensional Hilbert curve.
type Hilbert4D struct{ Order int }

// Size returns {2ⁿ, 2ⁿ, 2ⁿ, 2ⁿ} where n is the order.
func (h Hilbert4D) Size() []int { return []int{1 << h.Order, 1 << h.Order, 1 << h.Order, 1 << h.Order} }

func (h Hilbert4D) rot(reverse bool, n int, v []int, d int) {
	switch d {
	case 0:
		ops{swap{1, 3}, swap{0, 3}}.do(reverse, n, v)
	case 1, 2:
		ops{swap{0, 3}, swap{1, 3}}.do(reverse, n, v)
	case 3, 4:
		ops{flip{0, 1}, swap{2, 3}}.do(reverse, n, v)
	case 5, 6:
		ops{flip{1, 2}, swap{2, 3}}.do(reverse, n, v)
	case 7, 8:
		invert{0, 2}.do(n, v)
	case 9, 10:
		ops{flip{1, 2}, flip{2, 3}}.do(reverse, n, v)
	case 11, 12:
		ops{flip{0, 1}, flip{2, 3}}.do(reverse, n, v)
	case 13, 14:
		ops{flip{0, 3}, flip{1, 3}}.do(reverse, n, v)
	case 15:
		ops{flip{1, 3}, flip{0, 3}}.do(reverse, n, v)
	}
}

// Curve returns the curve coordinate of V. For order ≥2, Curve modifies V.
func (h Hilbert4D) Curve(v []int) uint64 {
	var d uint64
	N := 4
	for n := h.Order - 1; n >= 0; n-- {
		var e int
		for i := N - 1; i >= 0; i-- {
			v := v[i] >> n & 1
			e = e<<1 | (e^v)&1
		}

		d += uint64(e) << (N * n)
		h.rot(false, h.Order, v, e)
	}
	return d
}

// Space returns the spatial coordinates of D.
func (h Hilbert4D) Space(d uint64) []int {
	N := 4
	v := make([]int, N)
	for n := 0; n < h.Order; n++ {
		e := int(d & (1<<N - 1))
		h.rot(true, n, v, e)

		for i, e := 0, e; i < N; i++ {
			v[i] += (e ^ e>>1) & 1 << n
			e >>= 1
		}
		d >>= N
	}
	return v
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

// compose multiple coordinate modifications
type ops []op

func (c ops) do(reverse bool, n int, v []int) {
	if reverse {
		for i := len(c) - 1; i >= 0; i-- {
			c[i].do(n, v)
		}
	} else {
		for _, c := range c {
			c.do(n, v)
		}
	}
}
