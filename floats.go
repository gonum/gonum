// Copyright 2013 The Gonum Authors. All rights reserved.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file

// package floats provides a set of helper routines for dealing with slices
// of float64. The functions avoid allocations to allow for use within tight
// loops without garbage collection overhead.

package floats

import (
	"errors"
	"math"
	"sort"
)

// Add returns the element-wise sum of all the slices with the
// results stored in the first slice.
// For computational efficiency, it is assumed that all of
// the variadic arguments have the same length. If this is
// in doubt, EqLen can be used.
func Add(dst []float64, slices ...[]float64) []float64 {
	if len(slices) == 0 {
		return nil
	}
	if len(dst) != len(slices[0]) {
		panic("floats: length of destination does not match length of the slices")
	}
	for _, slice := range slices {
		for j, val := range slice {
			dst[j] += val
		}
	}
	return dst
}

// AddConst adds the value c to all of the values in s.
func AddConst(c float64, s []float64) {
	for i := range s {
		s[i] += c
	}
}

// AddScaled performs dst = dst + alpha * s.
// It panics if the lengths of dst and s are not equal.
func AddScaled(dst []float64, alpha float64, s []float64) {
	if len(dst) != len(s) {
		panic("floats: length of destination and source to not match")
	}
	for i, val := range s {
		dst[i] += alpha * val
	}
}

// AddScaledTo performs dst = y + alpha * s.
// It panics if the lengths of dst, y, and s are not equal.
func AddScaledTo(dst []float64, y []float64, alpha float64, s []float64) []float64 {
	if len(dst) != len(s) || len(dst) != len(y) {
		panic("floats: lengths of slices do not match")
	}
	for i, val := range s {
		dst[i] = y[i] + alpha*val
	}
	return dst
}

type argsort struct {
	s    []float64
	inds []int
}

func (a argsort) Len() int {
	return len(a.s)
}

func (a argsort) Less(i, j int) bool {
	return a.s[i] < a.s[j]
}

func (a argsort) Swap(i, j int) {
	a.s[i], a.s[j] = a.s[j], a.s[i]
	a.inds[i], a.inds[j] = a.inds[j], a.inds[i]
}

// Apply applies a function f (math.Exp, math.Sin, etc.) to every element
// of the slice s.
func Apply(f func(float64) float64, s []float64) {
	for i, val := range s {
		s[i] = f(val)
	}
}

// Argsort sorts the elements of s while tracking their original order.
// At the conclusion of Argsort, s will contain the original elements of s
// but sorted in increasing order, and inds will contain the original position
// of the elements in the slice such that s[i] = sOrig[inds[i]].
func Argsort(s []float64, inds []int) {
	if len(s) != len(inds) {
		panic("floats: length of inds does not match length of slice")
	}
	for i := range s {
		inds[i] = i
	}

	a := argsort{s: s, inds: inds}
	sort.Sort(a)
}

// Count applies the function f to every element of s and returns the number
// of times the function returned true.
func Count(f func(float64) bool, s []float64) int {
	var n int
	for _, val := range s {
		if f(val) {
			n++
		}
	}
	return n
}

// CumProd finds the cumulative product of the first i elements in
// s and puts them in place into the ith element of the
// destination. A panic will occur if lengths of do not match.
func CumProd(dst, s []float64) []float64 {
	if len(dst) != len(s) {
		panic("floats: length of destination does not match length of the source")
	}
	dst[0] = s[0]
	for i := 1; i < len(s); i++ {
		dst[i] = dst[i-1] * s[i]
	}
	return dst
}

// CumSum finds the cumulative sum of the first i elements in
// s and puts them in place into the ith element of the
// destination. A panic will occur if lengths of arguments do not match.
func CumSum(dst, s []float64) []float64 {
	if len(dst) != len(s) {
		panic("floats: length of destination does not match length of the source")
	}
	dst[0] = s[0]
	for i := 1; i < len(s); i++ {
		dst[i] = dst[i-1] + s[i]
	}
	return dst
}

// Dist computes the L-norm of s - t. See Norm for special cases.
func Distance(s []float64, t []float64, L float64) float64 {
	if len(s) != len(t) {
		panic("floats: slice lengths do not match")
	}
	if len(s) == 0 {
		return 0
	}
	var norm float64
	if L == 2 {
		for i, v := range s {
			diff := t[i] - v
			norm = math.Hypot(norm, diff)
		}
		return norm
	}
	if L == 1 {
		for i, v := range s {
			norm += math.Abs(t[i] - v)
		}
		return norm
	}
	if math.IsInf(L, 1) {
		for i, v := range s {
			absDiff := math.Abs(t[i] - v)
			if absDiff > norm {
				norm = absDiff
			}
		}
		return norm
	}
	for i, v := range s {
		norm += math.Pow(math.Abs(t[i]-v), L)
	}
	return math.Pow(norm, 1/L)
}

// Div performs element-wise division between s
// and t and stores the value in s. It panics if the
// lengths of s and t are not equal.
func Div(s []float64, t []float64) {
	if len(s) != len(t) {
		panic("floats: slice lengths do not match")
	}
	for i, val := range t {
		s[i] /= val
	}
}

// DivTo performs element-wise division between s
// and t and stores the value in dst. It panics if the
// lengths of s, t, and dst are not equal.
func DivTo(dst []float64, s []float64, t []float64) []float64 {
	if len(s) != len(t) || len(dst) != len(t) {
		panic("floats: slice lengths do not match")
	}
	for i, val := range t {
		dst[i] = s[i] / val
	}
	return dst
}

// Dot computes the dot product of s1 and s2, i.e.
// sum_{i = 1}^N s1[i]*s2[i].
// A panic will occur if lengths of arguments do not match.
func Dot(s1, s2 []float64) float64 {
	if len(s1) != len(s2) {
		panic("floats: lengths of the slices do not match")
	}
	var sum float64
	for i, val := range s1 {
		sum += val * s2[i]
	}
	return sum
}

// Equal returns true if the slices have equal lengths and
// all elements are numerically identical.
func Equal(s1, s2 []float64) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, val := range s1 {
		if s2[i] != val {
			return false
		}
	}
	return true
}

// EqualApprox returns true if the slices have equal lengths and
// all element pairs have an absolute tolerance less than tol or a
// relative tolerance less than tol.
func EqualApprox(s1, s2 []float64, tol float64) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, a := range s1 {
		if !EqualWithinAbsOrRel(a, s2[i], tol, tol) {
			return false
		}
	}
	return true
}

// EqualFunc returns true if the slices have the same lengths
// and the function returns true for all element pairs.
func EqualFunc(s1, s2 []float64, f func(float64, float64) bool) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, val := range s1 {
		if !f(val, s2[i]) {
			return false
		}
	}
	return true

}

// EqualWithinAbs returns true if a and b have an absolute
// difference of less than tol.
func EqualWithinAbs(a, b, tol float64) bool {
	return a == b || math.Abs(a-b) <= tol
}

const minNormalFloat64 = 2.2250738585072014e-308

// EqualWithinRel returns true if the difference between a and b
// is not greater than tol times the greater value.
func EqualWithinRel(a, b, tol float64) bool {
	if a == b {
		return true
	}
	delta := math.Abs(a - b)
	if delta <= minNormalFloat64 {
		return delta <= tol*minNormalFloat64
	}
	// We depend on the division in this relationship to identify
	// infinities (we rely on the NaN to fail the test) otherwise
	// we compare Infs of the same sign and evaluate Infs as equal
	// independent of sign.
	return delta/math.Max(math.Abs(a), math.Abs(b)) <= tol
}

// EqualWithinAbsOrRel returns true if a and b are equal to within
// the absolute tolerance.
func EqualWithinAbsOrRel(a, b, absTol, relTol float64) bool {
	if EqualWithinAbs(a, b, absTol) {
		return true
	}
	return EqualWithinRel(a, b, relTol)
}

// EqualWithinULP returns true if a and b are equal to within
// the specified number of floating point units in the last place.
func EqualWithinULP(a, b float64, ulp uint) bool {
	if a == b {
		return true
	}
	if math.IsNaN(a) || math.IsNaN(b) {
		return false
	}
	if math.Signbit(a) != math.Signbit(b) {
		return math.Float64bits(math.Abs(a))+math.Float64bits(math.Abs(b)) <= uint64(ulp)
	}
	return ulpDiff(math.Float64bits(a), math.Float64bits(b)) <= uint64(ulp)
}

func ulpDiff(a, b uint64) uint64 {
	if a > b {
		a, b = b, a
	}
	return b - a
}

// EqualLengths returns true if all of the slices have equal length,
// and false otherwise. Returns true if there are no input slices.
func EqualLengths(slices ...[]float64) bool {
	// This length check is needed: http://play.golang.org/p/sdty6YiLhM
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

// Fill loops over the elements of s and stores a value generated from f.
// f is called n times, where n = len(s)
func Fill(f func() float64, s []float64) {
	for i := range s {
		s[i] = f()
	}
}

// Find applies f to every element of s and returns the indices of the first
// k elements for which the f returns true, or all such elements
// if k < 0.
// Find will reslice inds to have 0 length, and will append
// found indices to inds.
// If k > 0 and there are fewer than k elements in s satisfying f,
// all of the found elements will be returned along with an error.
func Find(inds []int, f func(float64) bool, s []float64, k int) ([]int, error) {

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
	return inds, errors.New("floats: insufficient elements found")
}

// HasNaN returns true if the slice s has any values that are NaN and false
// otherwise.
func HasNaN(s []float64) bool {
	for _, v := range s {
		if math.IsNaN(v) {
			return true
		}
	}
	return false
}

// LogSpan returns a set of n equally spaced points in log space between,
// l and u where N is equal to len(dst). The first element of the
// resulting dst will be l and the final element of dst will be u.
// Panics if len(dst) < 2
// Note that this call will return NaNs if either l or u are negative, and
// will return all zeros if l or u is zero.
func LogSpan(dst []float64, l, u float64) []float64 {
	Span(dst, math.Log(l), math.Log(u))
	Apply(math.Exp, dst)
	return dst
}

// LogSumExp returns the log of the sum of the exponentials of the values in s
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
	return math.Log(lse) + maxval
}

// Max returns the maximum value in the slice and the location of
// the maximum value. If the input slice is empty, Max will panic.
func Max(s []float64) (max float64, ind int) {
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
// the minimum value. If the input slice is empty, Min will panic.
func Min(s []float64) (min float64, ind int) {
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

// Mul performs element-wise multiplication between s
// and t and stores the value in s. Panics if the
// lengths of s and t are not equal.
func Mul(s []float64, t []float64) {
	if len(s) != len(t) {
		panic("floats: slice lengths do not match")
	}
	for i, val := range t {
		s[i] *= val
	}
}

// MulTo performs element-wise multiplication between s
// and t and stores the value in dst. Panics if the
// lengths of s, t, and dst are not equal.
func MulTo(dst []float64, s []float64, t []float64) []float64 {
	if len(s) != len(t) || len(dst) != len(t) {
		panic("floats: slice lengths do not match")
	}
	for i, val := range t {
		dst[i] = val * s[i]
	}
	return dst
}

// Nearest returns the index of the element in s
// whose value is nearest to v.  If several such
// elements exist, the lowest index is returned.
func Nearest(s []float64, v float64) (ind int) {
	dist := math.Abs(v - s[0])
	ind = 0
	for i, val := range s {
		newDist := math.Abs(v - val)
		if newDist < dist {
			dist = newDist
			ind = i
		}
	}
	return
}

// NearestWithinSpan return the index of a hypothetical vector created
// by Span with length n and bounds l and u whose value is closest
// to v. Assumes u > l. If the value is greater than u or less than
// l, the function will panic.
func NearestWithinSpan(n int, l, u float64, v float64) int {
	if v < l || v > u {
		panic("floats: value outside span bounds")
	}

	// Can't guarantee anything about exactly halfway between
	// because of floating point weirdness
	return int((float64(n)-1)/(u-l)*(v-l) + 0.5)
}

// Norm returns the L norm of the slice S, defined as
// (sum_{i=1}^N s[i]^N)^{1/N}
// Special cases:
// L = math.Inf(1) gives the maximum value
// Does not correctly compute the zero norm (use Count).
func Norm(s []float64, L float64) (norm float64) {
	// Should this complain if L is not positive?
	// Should this be done in log space for better numerical stability?
	//	would be more cost
	//	maybe only if L is high?
	if len(s) == 0 {
		return 0
	}
	if L == 2 {
		twoNorm := math.Abs(s[0])
		for i := 1; i < len(s); i++ {
			twoNorm = math.Hypot(twoNorm, s[i])
		}
		return twoNorm
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
// Returns 1 if len(s) = 0.
func Prod(s []float64) (prod float64) {
	prod = 1
	for _, val := range s {
		prod *= val
	}
	return prod
}

// Scale multiplies every element in s by c.
func Scale(c float64, s []float64) {
	for i := range s {
		s[i] *= c
	}
}

// Span returns a set of N equally spaced points between l and u, where N
// is equal to the length of the destination. The first element of the destination
// is l, the final element of the destination is u.
// Panics if len(dst) < 2.
func Span(dst []float64, l, u float64) []float64 {
	n := len(dst)
	if n < 2 {
		panic("floats: destination must have length >1")
	}
	step := (u - l) / float64(n-1)
	for i := range dst {
		dst[i] = l + step*float64(i)
	}
	return dst
}

// Sub subtracts, element-wise, the second argument from the first. Panics if
// the lengths of s and t do not match.
func Sub(s, t []float64) {
	if len(s) != len(t) {
		panic("floats: length of the slices do not match")
	}
	for i, val := range t {
		s[i] -= val
	}
}

// SubTo subtracts, element-wise, the third argument from the second and
// stores the result in dst. Panics if the lengths of s and t do not match.
func SubTo(dst, s, t []float64) []float64 {
	if len(s) != len(t) {
		panic("floats: length of subtractor and subtractee do not match")
	}
	if len(dst) != len(s) {
		panic("floats: length of destination does not match length of subtractor")
	}
	for i, val := range t {
		dst[i] = s[i] - val
	}
	return dst
}

// Sum returns the sum of the elements of the slice.
func Sum(s []float64) (sum float64) {
	for _, val := range s {
		sum += val
	}
	return
}
