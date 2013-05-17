package sliceops

import "math"

// Returns the minimum value in the slice and the index of
// the minimum value. If the input slice is empty, zero is returned
// as the minimum value and -1 is returned as the index.
// Use: val,ind := sliceops.Min(slice)
func Min(s []float64) (min float64, ind int) {
	if len(s) == 0 {
		return min, -1 // Ind is -1 to make clear it's not the zeroth index.
	}
	min = s[0]
	ind = 0
	for i, val := range s {
		if val < min {
			min = val
			ind = i
		}
	}
	return min, ind
}

// Returns the maximum value in the slice and the location of
// the maximum value. If the input slice is empty, zero is returned
// as the minimum value and -1 is returned as the index.
// Use: val,ind := sliceops.Max(slice)
func Max(s []float64) (max float64, ind int) {
	if len(s) == 0 {
		return max, -1 // Ind is -1 to make clear it's not the zeroth index.
	}
	max = s[0]
	ind = 0
	for i, val := range s {
		if val > max {
			max = val
			ind = i
		}
	}
	return max, ind
}

// Returns the sum of the elements of the slice
func Sum(s []float64) (sum float64) {
	for _, val := range s {
		sum += val
	}
	return
}

// Returns the product of the elements of the slice
// Returns 1 if the input has length zero
func Prod(s []float64) (prod float64) {
	prod = 1
	for _, val := range s {
		prod *= val
	}
	return prod
}

// Returns true if all of the slices have equal length,
// and false otherwise. 
// Special case: Returns true if there are no input slices
func HasEqLen(slices ...[]float64) bool {
	if len(slices) == 0 {
		return true
	}
	l := len(slices[0])
	for i := 1; i < len(slices); i++ {
		if len(slices[i]) != l {
			return false
		}
	}
	return true
}

// Returns false if |s1[i] - s2[i]| > tol for any i.
// Assumes that the slices are of equal length. If this
// is in doubt it should be checked with HasEqLen
func Eq(s1, s2 []float64, tol float64) bool {
	for i, val := range s1 {
		if math.Abs(s2[i]-val) > tol {
			return false
		}
	}
	return true
}

// Finds the cumulative sum of the first i elements in 
// s and puts them in place into the ith element of the
// receiver. If the receiver is nil a new slice is created
func CumSum(receiver, s []float64) []float64 {
	if receiver == nil {
		receiver = make([]float64, len(s))
	}
	if len(s) == 0 {
		return receiver[:0]
	}
	receiver[0] = s[0]
	for i := 1; i < len(s); i++ {
		receiver[i] = receiver[i-1] + s[i]
	}
	return receiver
}

// Finds the cumulative product of the first i elements in 
// s and puts them in place into the ith element of the
// receiver. If the receiver is nil a new slice is created
func CumProd(receiver, s []float64) []float64 {
	if receiver == nil {
		receiver = make([]float64, len(s))
	}
	if len(s) == 0 {
		return receiver[:0]
	}
	receiver[0] = s[0]
	for i := 1; i < len(s); i++ {
		receiver[i] = receiver[i-1] * s[i]
	}
	return receiver
}

// Returns the element-wise sum of the last n slices
// and puts them in place into the first argument. If the 
// receiver is nil, a new slice of floats is created.
// For computational efficiency, it is assumed that all of
// the variadic arguments have the same length. If this is
// in doubt, EqLengths can be called. If no slices are input,
// the receiver is unchanged.
func ElemSum(receiver []float64, slices ...[]float64) []float64 {
	if len(slices) == 0 {
		return receiver
	}
	if receiver == nil {
		receiver = make([]float64, len(slices[0]))
	}
	for i, val := range slices[0] {
		receiver[i] = val
	}
	for i := 1; i < len(slices); i++ {
		for j, val := range slices[i] {
			receiver[j] += val
		}
	}
	return receiver
}

// Returns the L norm of the slice S.
// Special cases: 
// L = math.Inf(1) gives the maximum value
// Does not deal with the Zero norm. See Zero norm instead
func Norm(s []float64, L float64) (norm float64) {
	// Should this complain if L is not positive?
	// Should this be done in log space for better numerical stability?
	//	would be more cost
	//	maybe only if L is high?
	if math.IsInf(L, 1) {
		norm, _ = Max(s)
		return norm
	}
	for _, val := range s {
		norm += L * math.Pow(math.Abs(val), L)
	}
	return math.Pow(norm, 1/L)
}
