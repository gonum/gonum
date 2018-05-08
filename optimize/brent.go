package optimize

import "math"

const (
	phi           = 1.618043 // phi is the golden ratio
	tiny          = 1e-21
	limit float64 = 110
)

type brentIterType int

var (
	_ Method = &Brent{}
)

const (
	brentGo brentIterType = iota
	brentBrakA
	brentBrakB
	brentBrakC
	brentBrak1
	brentBrak2
	brentBrak3
	brentBrak4
	brentBrak5 // noop
)

// Brent is an optimization method from Richard Brent's "Algorithms for Minimization without Derivatives" (page 79)
type Brent struct {
	Min, Max, Limit float64 // brackets of absiccas

	x              float64
	a, b, c        float64
	fa, fb, fc, fw float64
	denom, w, wlim float64
	iter           brentIterType
}

func (b *Brent) Init(loc *Location) (Operation, error) {
	if len(loc.X) != 1 {
		panic("Expect only 1 parameter")
	}
	if b.Max <= b.Min && b.Max == 0 {
		b.Max = 1
	}

	b.a = b.Min
	b.b = b.Max
	b.x = loc.X[0]

	loc.X[0] = b.Min
	b.iter = brentBrakA
	return FuncEvaluation, nil
}

func (b *Brent) Iterate(loc *Location) (Operation, error) {
	switch b.iter {
	case brentBrakA, brentBrakB, brentBrakC:
		return b.bracket(b.iter, loc)
	case brentBrak1, brentBrak2, brentBrak3, brentBrak4, brentBrak5:
		return b.bracketLoop(loc)
	default:
	}
	return FuncEvaluation, nil
}

func (*Brent) Needs() struct {
	Gradient bool
	Hessian  bool
} {
	return struct {
		Gradient bool
		Hessian  bool
	}{false, false}
}

func (b *Brent) bracket(iter brentIterType, loc *Location) (Operation, error) {
	switch iter {
	case brentBrakA:
		b.fa = loc.F
		loc.X[0] = b.b
		b.iter = brentBrakB
		return FuncEvaluation, nil
	case brentBrakB:
		b.fb = loc.F
		if b.fa < b.fb {
			b.a, b.b = b.b, b.a
			b.fa, b.fb = b.fb, b.fa
		}
		loc.X[0] = b.b + phi*(b.b-b.a)
		b.iter = brentBrakC
		return FuncEvaluation, nil
	case brentBrakC:
		b.fc = loc.F
		return b.bracketLoop(loc)
	}
	panic("HALP")
}

func (b *Brent) bracketLoop(loc *Location) (Operation, error) {
	if b.fc > b.fb {
		loc.X[0] = b.x
		b.iter = brentGo
		return FuncEvaluation, nil
	}
	tmp1 := (b.b - b.a) * (b.fb - b.fc)
	tmp2 := (b.b - b.c) * (b.fb - b.fa)
	tmp3 := tmp2 - tmp1

	b.denom = 2 * tmp3
	if math.Abs(tmp3) < tiny {
		b.denom = 2 * tiny
	}

	b.w = b.b - ((b.b-b.c)*tmp2-(b.b-b.a)*tmp1)/b.denom
	b.wlim = b.b + b.Limit*(b.c-b.b)

	switch {
	case (b.w-b.c)*(b.b-b.w) > 0:
		loc.X[0] = b.w
		b.iter = brentBrak1
		return FuncEvaluation, nil
	case (b.w-b.wlim)*(b.wlim-b.c) >= 0:
		b.w = b.wlim
		loc.X[0] = b.w
		b.iter = brentBrak2
		return FuncEvaluation, nil
	case (b.w-b.wlim)*(b.c-b.w) > 0:
		loc.X[0] = b.w
		b.iter = brentBrak3
		return FuncEvaluation, nil
	default:
		b.w = b.c + phi*(b.c-b.b)
		loc.X[0] = b.w
		b.iter = brentBrak4
		return FuncEvaluation, nil
	}

}

func (b *Brent) doBracketLoop(iter brentIterType, loc *Location) (Operation, error) {
	switch iter {
	case brentBrak1:
		b.fw = loc.F
		switch {
		case b.fw < b.fc:
			b.a = b.b
			b.b = b.w
			b.fa = b.fb
			b.fb = b.fw

			loc.X[0] = b.x
			b.iter = brentGo
			return FuncEvaluation, nil
		case b.fw > b.fb:
			b.c = b.w
			b.fc = b.fw

			loc.X[0] = b.x
			b.iter = brentGo
			return FuncEvaluation, nil
		default:
			b.w = b.c + phi*(b.c-b.b)
			loc.X[0] = b.w
			b.iter = brentBrak5
			return FuncEvaluation, nil
		}
	case brentBrak2:
		b.fw = loc.F
		loc.X[0] = b.w
		b.iter = brentBrak5
		return FuncEvaluation, nil
	case brentBrak3:
		b.fw = loc.F
		if b.fw < b.c {
			b.b = b.c
			b.c = b.w
			b.w = b.c + phi*(b.c-b.b)
			b.fb = b.fc
			b.fc = b.fw
			loc.X[0] = b.w
			b.iter = brentBrak5
			return FuncEvaluation, nil
		}
	case brentBrak4:
		b.w = b.c + phi*(b.c-b.b)
		loc.X[0] = b.w
		b.iter = brentBrak5
		return FuncEvaluation, nil
	case brentBrak5:
		b.a, b.b, b.c = b.b, b.c, b.w
		b.fa, b.fb, b.fc = b.fb, b.fc, b.fw

		loc.X[0] = b.x
		b.iter = brentGo
		return FuncEvaluation, nil
	}
	panic("Unreachable")
}

// bracket finds a triple of a,b,c such that a < b < c and f(a) > f(b) < f(c).
func bracket(f func(float64) float64, min, max, limit float64) (a, b, c, fa, fb, fc float64) {
	fa = f(min)
	fb = f(max)

	a, b = min, max
	if fa < fb {
		a, b = b, a
		fa, fb = fb, fa
	}

	c = b + phi*(b-a)
	fc = f(c)

	for fc < fb {
		var fw float64
		tmp1 := (b - a) * (fb - fc)
		tmp2 := (b - c) * (fb - fa)
		tmp3 := tmp2 - tmp1

		denom := 2 * tmp3
		if math.Abs(tmp3) < tiny {
			denom = 2 * tiny
		}

		w := b - ((b-c)*tmp2-(b-a)*tmp1)/denom
		wlim := b + limit*(c-b)

		switch {
		case (w-c)*(b-w) > 0:
			fw = f(w)

			switch {
			case fw < fc:
				a, b = b, w
				fa, fb = fb, fw
				return
			case fw > fb:
				c, fc = w, fw
				return
			default:
				w = c + phi*(c-b)
				fw = f(w)
			}
		case (w-wlim)*(wlim-c) >= 0:
			w = wlim
			fw = f(w)

		case (w-wlim)*(c-w) > 0:
			fw = f(w)
			if fw < fc {
				b, c = c, w
				w = c + phi*(c-b)
				fb, fc = fc, fw
				fw = f(w)
			}
		default:
			w = c + phi*(c-b)
			fw = f(w)
		}
		a, b, c = b, c, w
		fa, fb, fc = fb, fc, fw
	}
	return
}
