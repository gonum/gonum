// Copyright ©2013 The Gonum Authors. All rights reserved.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cmplxs

import (
	"errors"
	"math"
	"math/cmplx"

	"gonum.org/v1/gonum/cmplxs/cscalar"
	"gonum.org/v1/gonum/internal/asm/c128"
)

const (
	zeroLength   = "cmplxs: zero length slice"
	shortSpan    = "cmplxs: slice length less than 2"
	badLength    = "cmplxs: slice lengths do not match"
	badDstLength = "cmplxs: destination slice length does not match input"
)

// Abs calculates the absolute values of the elements of s, and stores them in dst.
// It panics if the argument lengths do not match.
func Abs(dst []float64, s []complex128) {
	if len(dst) != len(s) {
		panic(badDstLength)
	}
	for i, v := range s {
		dst[i] = cmplx.Abs(v)
	}
}

// Add adds, element-wise, the elements of s and dst, and stores the result in dst.
// It panics if the argument lengths do not match.
func Add(dst, s []complex128) {
	if len(dst) != len(s) {
		panic(badLength)
	}
	c128.AxpyUnitaryTo(dst, 1, s, dst)
}

// AddTo adds, element-wise, the elements of s and t and
// stores the result in dst.
// It panics if the argument lengths do not match.
func AddTo(dst, s, t []complex128) []complex128 {
	if len(s) != len(t) {
		panic(badLength)
	}
	if len(dst) != len(s) {
		panic(badDstLength)
	}
	c128.AxpyUnitaryTo(dst, 1, s, t)
	return dst
}

// AddConst adds the scalar c to all of the values in dst.
func AddConst(c complex128, dst []complex128) {
	c128.AddConst(c, dst)
}

// AddScaled performs dst = dst + alpha * s.
// It panics if the slice argument lengths do not match.
func AddScaled(dst []complex128, alpha complex128, s []complex128) {
	if len(dst) != len(s) {
		panic(badLength)
	}
	c128.AxpyUnitaryTo(dst, alpha, s, dst)
}

// AddScaledTo performs dst = y + alpha * s, where alpha is a scalar,
// and dst, y and s are all slices.
// It panics if the slice argument lengths do not match.
//
// At the return of the function, dst[i] = y[i] + alpha * s[i]
func AddScaledTo(dst, y []complex128, alpha complex128, s []complex128) []complex128 {
	if len(s) != len(y) {
		panic(badLength)
	}
	if len(dst) != len(y) {
		panic(badDstLength)
	}
	c128.AxpyUnitaryTo(dst, alpha, s, y)
	return dst
}

// Count applies the function f to every element of s and returns the number
// of times the function returned true.
func Count(f func(complex128) bool, s []complex128) int {
	var n int
	for _, val := range s {
		if f(val) {
			n++
		}
	}
	return n
}

// Complex fills each of the elements of dst with the complex number
// constructed from the corresponding elements of real and imag.
// It panics if the argument lengths do not match.
func Complex(dst []complex128, real, imag []float64) []complex128 {
	if len(real) != len(imag) {
		panic(badLength)
	}
	if len(dst) != len(real) {
		panic(badDstLength)
	}
	if len(dst) == 0 {
		return dst
	}
	for i, r := range real {
		dst[i] = complex(r, imag[i])
	}
	return dst
}

// CumProd finds the cumulative product of elements of s and store it in
// place into dst so that
//
//	dst[i] = s[i] * s[i-1] * s[i-2] * ... * s[0]
//
// It panics if the argument lengths do not match.
func CumProd(dst, s []complex128) []complex128 {
	if len(dst) != len(s) {
		panic(badDstLength)
	}
	if len(dst) == 0 {
		return dst
	}
	return c128.CumProd(dst, s)
}

// CumSum finds the cumulative sum of elements of s and stores it in place
// into dst so that
//
//	dst[i] = s[i] + s[i-1] + s[i-2] + ... + s[0]
//
// It panics if the argument lengths do not match.
func CumSum(dst, s []complex128) []complex128 {
	if len(dst) != len(s) {
		panic(badDstLength)
	}
	if len(dst) == 0 {
		return dst
	}
	return c128.CumSum(dst, s)
}

// Distance computes the L-norm of s - t. See Norm for special cases.
// It panics if the slice argument lengths do not match.
func Distance(s, t []complex128, L float64) float64 {
	if len(s) != len(t) {
		panic(badLength)
	}
	if len(s) == 0 {
		return 0
	}

	var norm float64
	switch {
	case L == 2:
		return c128.L2DistanceUnitary(s, t)
	case L == 1:
		for i, v := range s {
			norm += cmplx.Abs(t[i] - v)
		}
		return norm
	case math.IsInf(L, 1):
		for i, v := range s {
			absDiff := cmplx.Abs(t[i] - v)
			if absDiff > norm {
				norm = absDiff
			}
		}
		return norm
	default:
		for i, v := range s {
			norm += math.Pow(cmplx.Abs(t[i]-v), L)
		}
		return math.Pow(norm, 1/L)
	}
}

// Div performs element-wise division dst / s
// and stores the result in dst.
// It panics if the argument lengths do not match.
func Div(dst, s []complex128) {
	if len(dst) != len(s) {
		panic(badLength)
	}
	c128.Div(dst, s)
}

// DivTo performs element-wise division s / t
// and stores the result in dst.
// It panics if the argument lengths do not match.
func DivTo(dst, s, t []complex128) []complex128 {
	if len(s) != len(t) {
		panic(badLength)
	}
	if len(dst) != len(s) {
		panic(badDstLength)
	}
	return c128.DivTo(dst, s, t)
}

// Dot computes the dot product of s1 and s2, i.e.
// sum_{i = 1}^N conj(s1[i])*s2[i].
// It panics if the argument lengths do not match.
func Dot(s1, s2 []complex128) complex128 {
	if len(s1) != len(s2) {
		panic(badLength)
	}
	return c128.DotUnitary(s1, s2)
}

// Equal returns true when the slices have equal lengths and
// all elements are numerically identical.
func Equal(s1, s2 []complex128) bool {
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

// EqualApprox returns true when the slices have equal lengths and
// all element pairs have an absolute tolerance less than tol or a
// relative tolerance less than tol.
func EqualApprox(s1, s2 []complex128, tol float64) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, a := range s1 {
		if !cscalar.EqualWithinAbsOrRel(a, s2[i], tol, tol) {
			return false
		}
	}
	return true
}

// EqualFunc returns true when the slices have the same lengths
// and the function returns true for all element pairs.
func EqualFunc(s1, s2 []complex128, f func(complex128, complex128) bool) bool {
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

// EqualLengths returns true when all of the slices have equal length,
// and false otherwise. It also returns true when there are no input slices.
func EqualLengths(slices ...[]complex128) bool {
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

// Find applies f to every element of s and returns the indices of the first
// k elements for which the f returns true, or all such elements
// if k < 0.
// Find will reslice inds to have 0 length, and will append
// found indices to inds.
// If k > 0 and there are fewer than k elements in s satisfying f,
// all of the found elements will be returned along with an error.
// At the return of the function, the input inds will be in an undetermined state.
func Find(inds []int, f func(complex128) bool, s []complex128, k int) ([]int, error) {
	// inds is also returned to allow for calling with nil.

	// Reslice inds to have zero length.
	inds = inds[:0]

	// If zero elements requested, can just return.
	if k == 0 {
		return inds, nil
	}

	// If k < 0, return all of the found indices.
	if k < 0 {
		for i, val := range s {
			if f(val) {
				inds = append(inds, i)
			}
		}
		return inds, nil
	}

	// Otherwise, find the first k elements.
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
	// Finished iterating over the loop, which means k elements were not found.
	return inds, errors.New("cmplxs: insufficient elements found")
}

// HasNaN returns true when the slice s has any values that are NaN and false
// otherwise.
func HasNaN(s []complex128) bool {
	for _, v := range s {
		if cmplx.IsNaN(v) {
			return true
		}
	}
	return false
}

// Imag places the imaginary components of src into dst.
// It panics if the argument lengths do not match.
func Imag(dst []float64, src []complex128) []float64 {
	if len(dst) != len(src) {
		panic(badDstLength)
	}
	if len(dst) == 0 {
		return dst
	}
	for i, z := range src {
		dst[i] = imag(z)
	}
	return dst
}

// LogSpan returns a set of n equally spaced points in log space between,
// l and u where N is equal to len(dst). The first element of the
// resulting dst will be l and the final element of dst will be u.
// Panics if len(dst) < 2
// Note that this call will return NaNs if either l or u are negative, and
// will return all zeros if l or u is zero.
// Also returns the mutated slice dst, so that it can be used in range, like:
//
//	for i, x := range LogSpan(dst, l, u) { ... }
func LogSpan(dst []complex128, l, u complex128) []complex128 {
	Span(dst, cmplx.Log(l), cmplx.Log(u))
	for i := range dst {
		dst[i] = cmplx.Exp(dst[i])
	}
	return dst
}

// MaxAbs returns the maximum absolute value in the input slice.
// It panics if s is zero length.
func MaxAbs(s []complex128) complex128 {
	return s[MaxAbsIdx(s)]
}

// MaxAbsIdx returns the index of the maximum absolute value in the input slice.
// If several entries have the maximum absolute value, the first such index is
// returned.
// It panics if s is zero length.
func MaxAbsIdx(s []complex128) int {
	if len(s) == 0 {
		panic(zeroLength)
	}
	max := math.NaN()
	var ind int
	for i, v := range s {
		if cmplx.IsNaN(v) {
			continue
		}
		if a := cmplx.Abs(v); a > max || math.IsNaN(max) {
			max = a
			ind = i
		}
	}
	return ind
}

// MinAbs returns the minimum absolute value in the input slice.
// It panics if s is zero length.
func MinAbs(s []complex128) complex128 {
	return s[MinAbsIdx(s)]
}

// MinAbsIdx returns the index of the minimum absolute value in the input slice. If several
// entries have the minimum absolute value, the first such index is returned.
// It panics if s is zero length.
func MinAbsIdx(s []complex128) int {
	if len(s) == 0 {
		panic(zeroLength)
	}
	min := math.NaN()
	var ind int
	for i, v := range s {
		if cmplx.IsNaN(v) {
			continue
		}
		if a := cmplx.Abs(v); a < min || math.IsNaN(min) {
			min = a
			ind = i
		}
	}
	return ind
}

// Mul performs element-wise multiplication between dst
// and s and stores the result in dst.
// It panics if the argument lengths do not match.
func Mul(dst, s []complex128) {
	if len(dst) != len(s) {
		panic(badLength)
	}
	for i, val := range s {
		dst[i] *= val
	}
}

// MulConj performs element-wise multiplication between dst
// and the conjugate of s and stores the result in dst.
// It panics if the argument lengths do not match.
func MulConj(dst, s []complex128) {
	if len(dst) != len(s) {
		panic(badLength)
	}
	for i, val := range s {
		dst[i] *= cmplx.Conj(val)
	}
}

// MulConjTo performs element-wise multiplication between s
// and the conjugate of t and stores the result in dst.
// It panics if the argument lengths do not match.
func MulConjTo(dst, s, t []complex128) []complex128 {
	if len(s) != len(t) {
		panic(badLength)
	}
	if len(dst) != len(s) {
		panic(badDstLength)
	}
	for i, val := range t {
		dst[i] = cmplx.Conj(val) * s[i]
	}
	return dst
}

// MulTo performs element-wise multiplication between s
// and t and stores the result in dst.
// It panics if the argument lengths do not match.
func MulTo(dst, s, t []complex128) []complex128 {
	if len(s) != len(t) {
		panic(badLength)
	}
	if len(dst) != len(s) {
		panic(badDstLength)
	}
	for i, val := range t {
		dst[i] = val * s[i]
	}
	return dst
}

// NearestIdx returns the index of the element in s
// whose value is nearest to v. If several such
// elements exist, the lowest index is returned.
// It panics if s is zero length.
func NearestIdx(s []complex128, v complex128) int {
	if len(s) == 0 {
		panic(zeroLength)
	}
	switch {
	case cmplx.IsNaN(v):
		return 0
	case cmplx.IsInf(v):
		return MaxAbsIdx(s)
	}
	var ind int
	dist := math.NaN()
	for i, val := range s {
		newDist := cmplx.Abs(v - val)
		// A NaN distance will not be closer.
		if math.IsNaN(newDist) {
			continue
		}
		if newDist < dist || math.IsNaN(dist) {
			dist = newDist
			ind = i
		}
	}
	return ind
}

// Norm returns the L-norm of the slice S, defined as
// (sum_{i=1}^N abs(s[i])^L)^{1/L}
// Special cases:
// L = math.Inf(1) gives the maximum absolute value.
// Does not correctly compute the zero norm (use Count).
func Norm(s []complex128, L float64) float64 {
	// Should this complain if L is not positive?
	// Should this be done in log space for better numerical stability?
	//	would be more cost
	//	maybe only if L is high?
	if len(s) == 0 {
		return 0
	}
	var norm float64
	switch {
	case L == 2:
		return c128.L2NormUnitary(s)
	case L == 1:
		for _, v := range s {
			norm += cmplx.Abs(v)
		}
		return norm
	case math.IsInf(L, 1):
		for _, v := range s {
			norm = math.Max(norm, cmplx.Abs(v))
		}
		return norm
	default:
		for _, v := range s {
			norm += math.Pow(cmplx.Abs(v), L)
		}
		return math.Pow(norm, 1/L)
	}
}

// Prod returns the product of the elements of the slice.
// Returns 1 if len(s) = 0.
func Prod(s []complex128) complex128 {
	prod := 1 + 0i
	for _, val := range s {
		prod *= val
	}
	return prod
}

// Real places the real components of src into dst.
// It panics if the argument lengths do not match.
func Real(dst []float64, src []complex128) []float64 {
	if len(dst) != len(src) {
		panic(badDstLength)
	}
	if len(dst) == 0 {
		return dst
	}
	for i, z := range src {
		dst[i] = real(z)
	}
	return dst
}

// Reverse reverses the order of elements in the slice.
func Reverse(s []complex128) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// Same returns true when the input slices have the same length and all
// elements have the same value with NaN treated as the same.
func Same(s, t []complex128) bool {
	if len(s) != len(t) {
		return false
	}
	for i, v := range s {
		w := t[i]
		if v != w && !(cmplx.IsNaN(v) && cmplx.IsNaN(w)) {
			return false
		}
	}
	return true
}

// Scale multiplies every element in dst by the scalar c.
func Scale(c complex128, dst []complex128) {
	if len(dst) > 0 {
		c128.ScalUnitary(c, dst)
	}
}

// ScaleReal multiplies every element in dst by the real scalar f.
func ScaleReal(f float64, dst []complex128) {
	for i, z := range dst {
		dst[i] = complex(f*real(z), f*imag(z))
	}
}

// ScaleRealTo multiplies the elements in s by the real scalar f and
// stores the result in dst.
// It panics if the slice argument lengths do not match.
func ScaleRealTo(dst []complex128, f float64, s []complex128) []complex128 {
	if len(dst) != len(s) {
		panic(badDstLength)
	}
	for i, z := range s {
		dst[i] = complex(f*real(z), f*imag(z))
	}
	return dst
}

// ScaleTo multiplies the elements in s by c and stores the result in dst.
// It panics if the slice argument lengths do not match.
func ScaleTo(dst []complex128, c complex128, s []complex128) []complex128 {
	if len(dst) != len(s) {
		panic(badDstLength)
	}
	if len(dst) > 0 {
		c128.ScalUnitaryTo(dst, c, s)
	}
	return dst
}

// Span returns a set of N equally spaced points between l and u, where N
// is equal to the length of the destination. The first element of the destination
// is l, the final element of the destination is u.
// It panics if the length of dst is less than 2.
//
// Span also returns the mutated slice dst, so that it can be used in range expressions,
// like:
//
//	for i, x := range Span(dst, l, u) { ... }
func Span(dst []complex128, l, u complex128) []complex128 {
	n := len(dst)
	if n < 2 {
		panic(shortSpan)
	}

	// Special cases for Inf and NaN.
	switch {
	case cmplx.IsNaN(l):
		for i := range dst[:len(dst)-1] {
			dst[i] = cmplx.NaN()
		}
		dst[len(dst)-1] = u
		return dst
	case cmplx.IsNaN(u):
		for i := range dst[1:] {
			dst[i+1] = cmplx.NaN()
		}
		dst[0] = l
		return dst
	case cmplx.IsInf(l) && cmplx.IsInf(u):
		for i := range dst {
			dst[i] = cmplx.Inf()
		}
		return dst
	case cmplx.IsInf(l):
		for i := range dst[:len(dst)-1] {
			dst[i] = l
		}
		dst[len(dst)-1] = u
		return dst
	case cmplx.IsInf(u):
		for i := range dst[1:] {
			dst[i+1] = u
		}
		dst[0] = l
		return dst
	}

	step := (u - l) / complex(float64(n-1), 0)
	for i := range dst {
		dst[i] = l + step*complex(float64(i), 0)
	}
	return dst
}

// Sub subtracts, element-wise, the elements of s from dst.
// It panics if the argument lengths do not match.
func Sub(dst, s []complex128) {
	if len(dst) != len(s) {
		panic(badLength)
	}
	c128.AxpyUnitaryTo(dst, -1, s, dst)
}

// SubTo subtracts, element-wise, the elements of t from s and
// stores the result in dst.
// It panics if the argument lengths do not match.
func SubTo(dst, s, t []complex128) []complex128 {
	if len(s) != len(t) {
		panic(badLength)
	}
	if len(dst) != len(s) {
		panic(badDstLength)
	}
	c128.AxpyUnitaryTo(dst, -1, t, s)
	return dst
}

// Sum returns the sum of the elements of the slice.
func Sum(s []complex128) complex128 {
	return c128.Sum(s)
}
