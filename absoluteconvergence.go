package optimize

import "math"

// FunctionConvergence implements FunctionConverger
type FunctionConvergence struct {
	Absolute   float64
	Relative   float64
	Iterations int

	absBest float64
	relBest float64
	absIter int
	relIter int
}

func (fc *FunctionConvergence) Init(f float64) {
	fc.absBest = f
	fc.relBest = f
	fc.absIter = 0
	fc.relIter = 0
}

func (fc *FunctionConvergence) FunctionConverged(f float64) Status {
	if f >= fc.absBest-fc.Absolute {
		fc.absIter++
		if fc.Iterations != 0 && fc.absIter >= fc.Iterations {
			return FunctionAbsoluteConvergence
		}
	} else {
		fc.absBest = f
		fc.absIter = 0
	}

	if math.Abs(f-fc.relBest)/math.Max(math.Abs(f), math.Abs(fc.relBest)) <= fc.Relative {
		fc.relIter++
		if fc.Iterations != 0 && fc.relIter >= fc.Iterations {
			return FunctionRelativeConvergence
		}
	} else {
		fc.relBest = f
		fc.relIter = 0
	}

	return NotTerminated
}
