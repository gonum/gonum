package sliceops

import "math"

type InsufficientElements struct{}

func (i InsufficientElements) Error() string {
	return "Insufficient elements found"
}

// Returns the element-wise sum of all the slices with the
// results stored in the first slice.
// Example: Add(a,b) // result will be a[i] = a[i] + b[i]
// a := make([]float64, len(b)); Add(a,b,c,d,e).
// For computational efficiency, it is assumed that all of
// the variadic arguments have the same length. If this is
// in doubt, EqLengths can be called.
func Add(slices ...[]float64) {
	if len(slices) == 0 {
		return
	}
	for i := 1; i < len(slices); i++ {
		for j, val := range slices[i] {
			slices[0][j] += val
		}
	}
}

// Adds a constant to all of the values in s
func AddConst(s []float64, c float64) {
	for i := range s {
		s[i] += c
	}
}

// Applies a function (math.Exp, math.Sin, etc.) to every element
// of the slice
func ApplyFunc(s []float64, f func(float64) float64) {
	for i, val := range s {
		s[i] = f(val)
	}
}

// Finds the cumulative product of the first i elements in
// s and puts them in place into the ith element of the
// destination. Assumes destination is at least as long as s
func Cumprod(dst, s []float64) []float64 {
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

// Finds the cumulative sum of the first i elements in
// s and puts them in place into the ith element of the
// destination. Assumes destination is at least as long as s
func Cumsum(dst, s []float64) {
	dst[0] = s[0]
	for i := 1; i < len(s); i++ {
		dst[i] = dst[i-1] + s[i]
	}
}

// Computes the dot product of s1 and s2, i.e.
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

// Returns false if |s1[i] - s2[i]| > tol for any i.
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

// Returns true if all of the slices have equal length,
// and false otherwise.
// Special case: Returns true if there are no input slices
func Eqlen(slices ...[]float64) bool {
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

// Applies a function returning a boolean to the elements of the slice
// and returns a list of indices for which the value is true
func Find(s []float64, f func(float64) bool) (inds []int) {
	// Not sure what an appropriate capacity is here. Don't want to make
	// it the length of the slice because if the slice is large that is
	// a lot of potentially wasted memory
	inds = make([]int, 0)
	for i, val := range s {
		if f(val) {
			inds = append(inds, i)
		}
	}
	return inds
}

// Applies a function returning a boolean to the elements of the slice
// and returns a list of the first k indices for which the value is true.
// If there are fewer than k indices for which the value is true, it returns
// the found indices and an error.
func FindFirst(s []float64, f func(float64) bool, k int) (inds []int, err error) {
	count := 0
	inds = make([]int, 0, k)
	for i, val := range s {
		if f(val) {
			inds = append(inds, i)
			count++
			if count == k {
				return inds, nil
			}
		}
	}
	return inds, InsufficientElements{}
}

// Returns a set of N equally spaced points between l and u, where N
// is equal to the length of the destination. The first element of the destination
// is l, the final element of the destination is u. Will panic if the destination has
// length < 2
func Linspace(dst []float64, l, u float64) {
	n := len(dst)
	step := (u - l) / float64(n-1)
	for i := range dst {
		dst[i] = l + step*float64(i)
	}
}

// Returns a set of N equally spaced points in log space between l and u, where N
// is equal to the length of the destination. The first element of the destination
// is l, the final element of the destination is u. Will panic if the destination has
// length < 2. Note that this call will return NaNs if l or u are negative, and
// zeros if l or u is zero.
func Logspace(dst []float64, l, u float64) {
	Linspace(dst, math.Log(l), math.Log(u))
	ApplyFunc(dst, math.Exp)
}

// Returns the log of the sum of the exponentials of the values in s
func Logsumexp(s []float64) (logsumexp float64) {
	// Want to do this in a numerically stable way which avoids
	// overflow and underflow
	// First, find the maximum value in the slice.
	maxval, _ := Max(s)
	if math.IsInf(maxval, 0) {
		// If it's infinity either way, the logsumexp will be infinity as well
		// returning now avoids NaNs
		return maxval
	}
	// Subtract off the largest value, so the largest value in
	// the new slice is 0
	AddConst(s, -maxval)
	defer AddConst(s, maxval) // make sure we add it back on at the end

	// compute the sumexp part
	for _, val := range s {
		logsumexp += math.Exp(val)
	}
	// Take the log and add back on the constant taken out
	logsumexp = math.Log(logsumexp) + maxval
	return logsumexp
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

// Returns the L norm of the slice S.
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

// Returns the product of the elements of the slice
// Returns 1 if the input has length zero
func Prod(s []float64) (prod float64) {
	prod = 1
	for _, val := range s {
		prod *= val
	}
	return prod
}

// Multiplies every element in s by a constant in place
func Scale(s []float64, c float64) {
	for i := range s {
		s[i] *= c
	}
}

// Subtract, element-wise, the first argument from the second. Assumes
// the lengths of s and t match (can be tested with EqLen)
func Sub(s, t []float64) {
	for i, val := range t {
		s[i] -= val
	}
}

// Subtract, element-wise, the first argument from the second and
// store the result in destination. Assumes the lengths of s and t match
// (can be tested with EqLen)
func SubDst(dst, s, t []float64) {
	for i, val := range t {
		dst[i] = s[i] - val
	}
}

// Returns the sum of the elements of the slice
func Sum(s []float64) (sum float64) {
	for _, val := range s {
		sum += val
	}
	return
}
