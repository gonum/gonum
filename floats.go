package floats

import "math"

// InsufficientElements is an error type used by FindFirst
type InsufficientElements struct{}

// Error makes InsufficientElements satisfy the error interface
func (i *InsufficientElements) Error() string {
	return "Insufficient elements found"
}

// Add returns the element-wise sum of all the slices with the
// results stored in the first slice.
// Example: Add(a,b) // result will be a[i] = a[i] + b[i]
// a := make([]float64, len(b)); Add(a,b,c,d,e).
// For computational efficiency, it is assumed that all of
// the variadic arguments have the same length. If this is
// in doubt, EqLengths can be called.
func Add(dst []float64, slices ...[]float64) {
	if len(slices) == 0 {
		return
	}
	for i := 0; i < len(slices); i++ {
		for j, val := range slices[i] {
			dst[j] += val
		}
	}
}

// AddConst adds a constant to all of the values in s
func AddConst(s []float64, c float64) {
	for i := range s {
		s[i] += c
	}
}

// ApplyFunc applies a function (math.Exp, math.Sin, etc.) to every element
// of the slice
func Apply(s []float64, f func(float64) float64) {
	for i, val := range s {
		s[i] = f(val)
	}
}

// Cumprod finds the cumulative product of the first i elements in
// s and puts them in place into the ith element of the
// destination. Assumes destination is at least as long as s
func CumProd(dst, s []float64) []float64 {
	if dst == nil {
		dst = make([]float64, len(s))
	}
	if len(s) == 0 {
		return dst[:0]
	}
	dst[0] = s[0]
	for i := 1; i < len(s); i++ {
		dst[i] = dst[i-1] * s[i]
	}
	return dst
}

// Cumsum finds the cumulative sum of the first i elements in
// s and puts them in place into the ith element of the
// destination. Assumes destination is at least as long as s
func CumSum(dst, s []float64) {
	dst[0] = s[0]
	for i := 1; i < len(s); i++ {
		dst[i] = dst[i-1] + s[i]
	}
}

// Dot computes the dot product of s1 and s2, i.e.
// sum_{i = 1}^N s1[n]*s2[n]
// Assumes the slices are of equal length. If this is
// in doubt it should be checked with Eqlen
func Dot(s1, s2 []float64) float64 {
	var sum float64
	for i, val := range s1 {
		sum += val * s2[i]
	}
	return sum
}

// Eq returns false if |s1[i] - s2[i]| > tol for any i.
// Assumes that the slices are of equal length. If this
// is in doubt it should be checked with Eqlen
func Eq(s1, s2 []float64, tol float64) bool {
	for i, val := range s1 {
		if math.Abs(s2[i]-val) > tol {
			return false
		}
	}
	return true
}

// Eqlen returns true if all of the slices have equal length,
// and false otherwise.
// Special case: Returns true if there are no input slices
func EqLen(slices ...[]float64) bool {
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

// Find finds the first k indices of the slice s for which
// the function f returns true and stores them in the slice
// inds. If k < 0, all such elements are found.
// Find will reslice inds to have 0 length, and will append
// found indices to inds.
// If there are fewer than k elements in s satisfying f,
// all of the found elements will be returned along with an
// InsufficientElements error
// TODO: Add example for nil slice, appending to the end of a larger slice
// reusing memory, insufficient elements
func Find(inds []int, k int, s []float64, f func(float64) bool) ([]int, error) {

	// inds is also returned to allow for calling with nil

	// Reslice inds to have zero length
	inds = inds[:0]

	// If zero elements requested, can just return
	if k == 0 {
		return inds, nil
	}

	// If k < 0, return all of the found indices
	if k < 0 {
		for i, val := range s {
			if f(val) {
				inds = append(inds, i)
			}
		}
		return inds, nil
	}

	// Otherwise, find the first k elements
	nFound := 0
	for i, val := range s {
		if f(val) {
			inds = append(inds, i)
			nFound++
			if nFound == k {
				return inds, nil
			}
		}
	}
	// Finished iterating over the loop, which means k elements were not found
	return inds, &InsufficientElements{}
}

// LogSpan returns a set of N equally spaced points in log space between l and u, where N
// is equal to the length of the destination. The first element of the destination
// is l, the final element of the destination is u. Will panic if the destination has
// length < 2. Note that this call will return NaNs if l or u are negative, and
// zeros if l or u is zero.
func LogSpan(dst []float64, l, u float64) {
	Span(dst, math.Log(l), math.Log(u))
	Apply(dst, math.Exp)
}

// Logsumexp returns the log of the sum of the exponentials of the values in s
func LogSumExp(s []float64) (lse float64) {
	// Want to do this in a numerically stable way which avoids
	// overflow and underflow
	// First, find the maximum value in the slice.
	maxval, _ := Max(s)
	if math.IsInf(maxval, 0) {
		// If it's infinity either way, the logsumexp will be infinity as well
		// returning now avoids NaNs
		return maxval
	}
	// Compute the sumexp part
	for _, val := range s {
		lse += math.Exp(val - maxval)
	}
	// Take the log and add back on the constant taken out
	lse = math.Log(lse) + maxval
	return lse
}

// Max returns the maximum value in the slice and the location of
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

// Min returns the minimum value in the slice and the index of
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

// Norm returns the L norm of the slice S.
// Special cases:
// L = math.Inf(1) gives the maximum value
// Does not correctly compute the zero norm, as the zero norm is a count.
func Norm(s []float64, L float64) (norm float64) {
	// Should this complain if L is not positive?
	// Should this be done in log space for better numerical stability?
	//	would be more cost
	//	maybe only if L is high?
	if L == 2 {
		for _, val := range s {
			norm += val * val
		}
		return math.Pow(norm, 0.5)
	}
	if L == 1 {
		for _, val := range s {
			norm += math.Abs(val)
		}
		return norm
	}
	if math.IsInf(L, 1) {
		norm, _ = Max(s)
		return norm
	}
	for _, val := range s {
		norm += math.Pow(math.Abs(val), L)
	}
	return math.Pow(norm, 1/L)
}

// Prod returns the product of the elements of the slice
// Returns 1 if the input has length zero
func Prod(s []float64) (prod float64) {
	prod = 1
	for _, val := range s {
		prod *= val
	}
	return prod
}

// Scale multiplies every element in s by a constant in place
func Scale(s []float64, c float64) {
	for i := range s {
		s[i] *= c
	}
}

// Span returns a set of N equally spaced points between l and u, where N
// is equal to the length of the destination. The first element of the destination
// is l, the final element of the destination is u. Will panic if the destination has
// length < 2
func Span(dst []float64, l, u float64) {
	n := len(dst)
	step := (u - l) / float64(n-1)
	for i := range dst {
		dst[i] = l + step*float64(i)
	}
}

// Sub subtracts, element-wise, the first argument from the second. Assumes
// the lengths of s and t match (can be tested with EqLen)
func Sub(s, t []float64) {
	for i, val := range t {
		s[i] -= val
	}
}

// SubDst subtracts, element-wise, the first argument from the second and
// store the result in destination. Assumes the lengths of s and t match
// (can be tested with EqLen)
func SubDst(dst, s, t []float64) {
	for i, val := range t {
		dst[i] = s[i] - val
	}
}

// Sum returns the sum of the elements of the slice
func Sum(s []float64) (sum float64) {
	for _, val := range s {
		sum += val
	}
	return
}
