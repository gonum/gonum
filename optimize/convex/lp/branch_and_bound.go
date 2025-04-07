package lp

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

// wholeNumberTol is the tolerance for a number differing from a whole value
const wholeNumberTol = 1e-14

// Branch and Bound uses simplex to resolve an integer programming problem where some of the constraints may be on the variables
// being whole numbers.
//
// When BnB performs simplex and finds a solution that doesn't hold the whole number constraints, it creates two new simplex problems
// with constraints on the problem variable to be bellow the floor and above the ceil of the existing solution, repeating this process
// until it reaches a solution that holds the integer constraints.
//
// the `whole` parameter is a list of booleans the same length as the variable list that should indicate which variables should have a whole number constraint.
func BNB(c []float64, A mat.Matrix, b []float64, G mat.Matrix, h []float64, whole []bool, tol float64) (float64, []float64, error) {
	q := []problem{}

	q = append(q, problem{
		g: G,
		h: h,
	})

	best_attempt := attempt{
		fitness: math.Inf(1),
	}

	for len(q) > 0 {
		p := q[len(q)-1]
		q = q[:len(q)-1]

		cNew, ANew, bNew := Convert(c, p.g, p.h, A, b)

		fit, x, _, err := simplex(nil, cNew, ANew, bNew, tol)

		if err != nil {
			switch err {
			// Problems are programatically created,
			// if one of them is infeasible, it could be that the branch we're testing
			// now is infeasible but some other branch isn't, just skip this problem
			// and don't add any new ones, if no branches have any feasible solutions
			// then the whole problem is infeasible.
			case ErrInfeasible:
				continue
			default:
				return math.NaN(), nil, err
			}
		}

		//check if the integer variable constraints hold
		broken_whole := 0
		is_whole := true
		for i, b := range whole {
			if b && !checkWholeNumber(x[i]) {
				is_whole = false
				broken_whole = i
				break
			}
		}

		if is_whole {
			if fit < best_attempt.fitness {
				best_attempt = attempt{
					x:       x,
					fitness: fit,
				}
			}
		} else {
			if fit > best_attempt.fitness {
				continue
			}
			lowX := math.Floor(x[broken_whole])
			highX := math.Ceil(x[broken_whole])

			row, col := p.g.Dims()
			// new problem lower bounded
			lowG := mat.NewDense(row+1, col, nil)
			lowG.Copy(p.g)
			lowG.Set(row, broken_whole, 1)
			lowh := make([]float64, row+1)
			copy(lowh, p.h)
			lowh[row] = lowX
			q = append(q, problem{
				g: lowG,
				h: lowh,
			})

			// new problem higher bounded
			highG := mat.NewDense(row+1, col, nil)
			highG.Copy(p.g)
			highG.Set(row, broken_whole, -1)
			highh := make([]float64, row+1)
			copy(highh, p.h)
			highh[row] = -highX
			q = append(q, problem{
				g: highG,
				h: highh,
			})
		}
	}

	if math.IsInf(best_attempt.fitness, 0) {
		return math.NaN(), nil, ErrInfeasible
	}

	return best_attempt.fitness, best_attempt.x, nil
}

type problem struct {
	g mat.Matrix
	h []float64
}

type attempt struct {
	fitness float64
	x       []float64
}

func checkWholeNumber(x float64) bool {
	diff := math.Abs(x - math.Round(x))
	return diff < wholeNumberTol
}
