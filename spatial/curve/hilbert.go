// Hilbert Curves
//
// The Hilbert curve is a continuous, space filling curve first described by
// David Hilbert. The implementation of Hilbert2D is based on example code from
// the Wikipedia article
// (https://en.wikipedia.org/w/index.php?title=Hilbert_curve&oldid=1011599190).
// The implementation of Hilbert3D and Hilbert4D are extrapolated from
// Hilbert2D.
//
// For the first-order k-dimensional Hilbert curve, a spatial point V is mapped
// to a point on the curve D by XOR - each dimension of V is expected to be 0 or
// 1:
//
//     func map1stOrder(k int, v []int) (d int) {
//         for i := k - 1; i >= 0; i-- {
//             d = d<<1 | (d^v[i])&1
//         }
//         return d
//     }
//
// In a 2-space with the origin at the bottom left, this results in a ⊐ shape,
// wound counter clockwise.
package curve

// NewHilbert returns a 2-, 3-, or 4-dimensional Hilbert curve. NewHilbert will
// panic if dimension is not in the range [2, 4] or if order is not in the range
// [1, ∞).
//
// The runtime of Space and Curve scales as O(k∙n) where k is the dimension and
// n is the order. The length of the curve is 2^(k∙n).
func NewHilbert(order, dimension int) SpaceFilling {
	if order < 1 {
		panic("order must be positive")
	}

	switch dimension {
	case 2:
		return Hilbert2D{Order: order}
	case 3:
		return Hilbert3D{Order: order}
	case 4:
		return Hilbert4D{Order: order}
	default:
		if dimension < 2 {
			panic("dimensions must be 2 or greater")
		}
		panic("more than 4 dimensions are not supported")
	}
}

// V may get mangled by Curve, so copy it
func dup(v []int) []int {
	u := make([]int, len(v))
	copy(u, v)
	return u
}

type coordInvert struct{ I, J int } // invert I and J
type coordSwap struct{ I, J int }   // swap I and J
type coordFlip struct{ I, J int }   // swap and invert I and J

func (c coordInvert) do(n int, v []int) { v[c.I], v[c.J] = v[c.I]^(1<<n-1), v[c.J]^(1<<n-1) }
func (c coordSwap) do(n int, v []int)   { v[c.I], v[c.J] = v[c.J], v[c.I] }
func (c coordFlip) do(n int, v []int)   { v[c.I], v[c.J] = v[c.J]^(1<<n-1), v[c.I]^(1<<n-1) }

// compose multiple coordinate modifications
type multiCoord []interface{ do(int, []int) }

func (c multiCoord) do(reverse bool, n int, v []int) {
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

// Hilbert2D is a 2-dimensional Hilbert curve.
type Hilbert2D struct{ Order int }

// Size returns {2ⁿ, 2ⁿ} where n is the order.
func (h Hilbert2D) Size() []int { return []int{1 << h.Order, 1 << h.Order} }

func (h Hilbert2D) rot(n int, v []int, d int) {
	switch d {
	case 0:
		coordSwap{0, 1}.do(n, v)
	case 3:
		coordFlip{0, 1}.do(n, v)
	}
}

// Curve returns the curve coordinate of V.
func (h Hilbert2D) Curve(v ...int) Point {
	v = dup(v)
	var d Point
	for n := h.Order - 1; n >= 0; n-- {
		rx := (v[0] >> n) & 1
		ry := (v[1] >> n) & 1
		rd := ry<<1 | (ry ^ rx)
		d += Point(rd) << (2 * n)
		h.rot(h.Order, v, rd)
	}
	return d
}

// Space2D returns the spatial coordinates of D.
func (h Hilbert2D) Space2D(d Point) [2]int {
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
func (h Hilbert2D) Space(d Point) []int {
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
		multiCoord{coordSwap{1, 2}, coordSwap{0, 2}}.do(reverse, n, v)
	case 1, 2:
		multiCoord{coordSwap{0, 2}, coordSwap{1, 2}}.do(reverse, n, v)
	case 3, 4:
		coordInvert{0, 1}.do(n, v)
	case 5, 6:
		multiCoord{coordFlip{0, 2}, coordFlip{1, 2}}.do(reverse, n, v)
	case 7:
		multiCoord{coordFlip{1, 2}, coordFlip{0, 2}}.do(reverse, n, v)
	}
}

// Curve returns the curve coordinate of V.
func (h Hilbert3D) Curve(v ...int) Point {
	v = dup(v)
	var d Point
	for n := h.Order - 1; n >= 0; n-- {
		rx := (v[0] >> n) & 1
		ry := (v[1] >> n) & 1
		rz := (v[2] >> n) & 1
		rd := rz<<2 | (rz^ry)<<1 | (rz ^ ry ^ rx)
		d += Point(rd) << (3 * n)
		h.rot(false, h.Order, v, rd)
	}
	return d
}

// Space3D returns the spatial coordinates of D.
func (h Hilbert3D) Space3D(d Point) [3]int {
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
func (h Hilbert3D) Space(v Point) []int {
	xy := h.Space3D(v)
	return xy[:]
}

// Hilbert4D is a 4-dimensional Hilbert curve.
type Hilbert4D struct{ Order int }

// Size returns {2ⁿ, 2ⁿ, 2ⁿ, 2ⁿ} where n is the order.
func (h Hilbert4D) Size() []int { return []int{1 << h.Order, 1 << h.Order, 1 << h.Order, 1 << h.Order} }

func (h Hilbert4D) rot(reverse bool, n int, v []int, d int) {
	switch d {
	case 0x0:
		multiCoord{coordSwap{1, 3}, coordSwap{0, 3}}.do(reverse, n, v)
	case 0x1, 0x2:
		multiCoord{coordSwap{0, 3}, coordSwap{1, 3}}.do(reverse, n, v)
	case 0x3, 0x4:
		multiCoord{coordFlip{0, 1}, coordSwap{2, 3}}.do(reverse, n, v)
	case 0x5, 0x6:
		multiCoord{coordFlip{1, 2}, coordSwap{2, 3}}.do(reverse, n, v)
	case 0x7, 0x8:
		coordInvert{0, 2}.do(n, v)
	case 0x9, 0xA:
		multiCoord{coordFlip{1, 2}, coordFlip{2, 3}}.do(reverse, n, v)
	case 0xB, 0xC:
		multiCoord{coordFlip{0, 1}, coordFlip{2, 3}}.do(reverse, n, v)
	case 0xD, 0xE:
		multiCoord{coordFlip{0, 3}, coordFlip{1, 3}}.do(reverse, n, v)
	case 0xF:
		multiCoord{coordFlip{1, 3}, coordFlip{0, 3}}.do(reverse, n, v)
	}
}

// Curve returns the curve coordinate of V.
func (h Hilbert4D) Curve(v ...int) Point {
	v = dup(v)
	var d Point
	N := 4
	for n := h.Order - 1; n >= 0; n-- {
		var e int
		for i := N - 1; i >= 0; i-- {
			v := v[i] >> n & 1
			e = e<<1 | (e^v)&1
		}

		d += Point(e) << (N * n)
		h.rot(false, h.Order, v, e)
	}
	return d
}

// Space returns the spatial coordinates of D.
func (h Hilbert4D) Space(d Point) []int {
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
