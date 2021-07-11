package curve

func NewHilbert(order, dimensions int) SpaceFilling {
	if order < 0 {
		// order 0 doesn't really make sense, but it does work
		panic("order must be positive")
	}
	if dimensions < 2 {
		panic("dimensions must be 2 or greater")
	}

	switch dimensions {
	case 2:
		return Hilbert2D{Order: order}
	case 3:
		return Hilbert3D{Order: order}
	default:
		return Hilbert{Order: order, Dimension: dimensions}
	}
}

func mirror(n int, v ...*int) {
	for _, v := range v {
		*v = (1 << n) - 1 - *v
	}
}

type Hilbert2D struct{ Order int }

func (h Hilbert2D) rot(n int, v []int, rx, ry int) {
	if rx == 0 {
		if ry != 0 {
			mirror(n, &v[0], &v[1])
		}
		v[0], v[1] = v[1], v[0]
	}
}

func (h Hilbert2D) Curve(v ...int) Point {
	// Based on https://en.wikipedia.org/w/index.php?title=Hilbert_curve&oldid=1011599190

	var d Point
	for n := h.Order - 1; n >= 0; n-- {
		rx := (v[0] >> n) & 1
		ry := (v[1] >> n) & 1
		rd := (3 * ry) ^ rx
		d += Point(rd) << (2 * n)
		h.rot(h.Order, v, rx, ry)
	}
	return d
}

func (h Hilbert2D) Space2D(d Point) [2]int {
	// Based on https://en.wikipedia.org/w/index.php?title=Hilbert_curve&oldid=1011599190

	var v [2]int
	for n := 0; n < h.Order; n++ {
		e := int(d & 3)
		ry := e >> 1
		rx := (e>>0 ^ e>>1) & 1
		h.rot(n, v[:], rx, ry)
		v[0] += rx << n
		v[1] += ry << n
		d >>= 2
	}
	return v
}

func (h Hilbert2D) Space(v Point) []int {
	xy := h.Space2D(v)
	return xy[:]
}

type Hilbert3D struct{ Order int }

func (h Hilbert3D) rot(n int, v []int, rx, ry, rz int) {
	if rx == 1 {
		if rz == 1 {
			mirror(n, &v[1], &v[2])
		}
		v[1], v[2] = v[2], v[1]
	} else if ry == 1 {
		mirror(n, &v[0], &v[1])
	} else {
		if rz == 1 {
			mirror(n, &v[0], &v[2])
		}
		v[0], v[2] = v[2], v[0]
	}

	// rx == 0 && ry == 0 && rz == 0: *x, *z = *z, *x
	// rx == 1 && ry == 0 && rz == 0: *y, *z = *z, *y
	// rx == 1 && ry == 1 && rz == 0: *y, *z = *z, *y
	// rx == 0 && ry == 1 && rz == 0: mirror(n, x, y)
	// rx == 0 && ry == 1 && rz == 1: mirror(n, x, y)
	// rx == 1 && ry == 1 && rz == 1: *y, *z = *z, *y; mirror(n, y, z)
	// rx == 1 && ry == 0 && rz == 1: *y, *z = *z, *y; mirror(n, y, z)
	// rx == 0 && ry == 0 && rz == 1: *x, *z = *z, *x; mirror(n, x, z)
}

func (h Hilbert3D) Curve(v ...int) Point {
	var d Point
	for n := h.Order - 1; n >= 0; n-- {
		rx := (v[0] >> n) & 1
		ry := (v[1] >> n) & 1
		rz := (v[2] >> n) & 1
		rd := rz<<2 | (rz^ry)<<1 | (rz^ry^rx)<<0
		d += Point(rd) << (3 * n)
		h.rot(h.Order, v, rx, ry, rz)
	}
	return d
}

func (h Hilbert3D) Space3D(d Point) [3]int {
	var v [3]int
	for n := 0; n < h.Order; n++ {
		e := int(d & 7)
		rz := e >> 2
		ry := (e>>1 ^ e>>2) & 1
		rx := (e>>0 ^ e>>1) & 1
		h.rot(n, v[:], rx, ry, rz)
		v[0] += rx << n
		v[1] += ry << n
		v[2] += rz << n
		d >>= 3
	}
	return v
}

func (h Hilbert3D) Space(v Point) []int {
	xy := h.Space3D(v)
	return xy[:]
}

type Hilbert struct{ Order, Dimension int }

func (h Hilbert) rot(n int, v []int, rv []int) {
	switch h.Dimension {
	case 2:
		Hilbert2D{}.rot(n, v, rv[0], rv[1])
	case 3:
		Hilbert3D{}.rot(n, v, rv[0], rv[1], rv[2])
	default:
		panic("TODO")
	}
}

func (h Hilbert) Curve(v ...int) Point {
	var d Point
	rv := make([]int, h.Dimension)
	for n := h.Order - 1; n >= 0; n-- {
		for i := 0; i < h.Dimension; i++ {
			rv[i] = v[i] >> n & 1
		}
		var rd int
		for i := h.Dimension - 1; i >= 0; i-- {
			rd = rd<<1 | (rd^rv[i])&1
		}
		d += Point(rd) << (h.Dimension * n)
		h.rot(h.Order, v, rv)
	}
	return d
}

func (h Hilbert) Space(d Point) []int {
	v := make([]int, h.Dimension)
	rv := make([]int, h.Dimension)
	for n := 0; n < h.Order; n++ {
		e := int(d & (1<<h.Dimension - 1))
		for i := 0; i < h.Dimension; i++ {
			rv[i] = (e ^ e>>1) & 1
			e >>= 1
		}
		h.rot(n, v, rv)
		for i := 0; i < h.Dimension; i++ {
			v[i] += rv[i] << n
		}
		d >>= Point(h.Dimension)
	}
	return v
}
