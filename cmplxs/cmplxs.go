// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file

package cmplxs

import (
	"errors"
	"math"
	"math/cmplx"

	"gonum.org/v1/gonum/internal/asm/c128"
)

// Abs generates cmplx.Abs, element-wise, for the elements of s and stores in dst.
// dst is also returned.
// It panics if the lengths of dst and s are not equal.
func Abs(dst []float64, s []complex128) []float64 {
	if len(dst) != len(s) {
		panic("cmplxs.Abs: length of the slices do not match")
	}
	for i, v := range s {
		dst[i] = cmplx.Abs(v)
	}
	return dst
}

// Add adds, element-wise, the elements of s and dst, and stores in dst.
// It panics if the lengths of dst and s are not equal.
func Add(dst, s []complex128) {
	if len(dst) != len(s) {
		panic("cmplxs.Add: length of the slices do not match")
	}
	c128.AxpyUnitaryTo(dst, 1, s, dst)
}

// AddConst adds the scalar c to all of the values in dst.
func AddConst(c complex128, dst []complex128) {
	for i := range dst {
		dst[i] += c
	}
}

// AddTo adds, element-wise, the elements of s and t and
// stores the result in dst. Panics if the lengths of s, t and dst do not match.
func AddTo(dst, s, t []complex128) []complex128 {
	if len(s) != len(t) {
		panic("cmplxs.AddTo: length of adders do not match")
	}
	if len(dst) != len(s) {
		panic("cmplxs.AddTo: length of destination does not match length of adder")
	}
	c128.AxpyUnitaryTo(dst, 1, s, t)
	return dst
}

// AddScaled performs dst = dst + alpha * s.
// It panics if the lengths of dst and s are not equal.
func AddScaled(dst []complex128, alpha complex128, s []complex128) {
	if len(dst) != len(s) {
		panic("cmplxs: length of destination and source to not match")
	}
	c128.AxpyUnitaryTo(dst, alpha, s, dst)
}

// AddScaledTo performs dst = y + alpha * s, where alpha is a scalar,
// and dst, y and s are all slices.
// It panics if the lengths of dst, y, and s are not equal.
//
// At the return of the function, dst[i] = y[i] + alpha * s[i]
func AddScaledTo(dst, y []complex128, alpha complex128, s []complex128) []complex128 {
	if len(dst) != len(s) || len(dst) != len(y) {
		panic("cmplxs.AddScaledTo: lengths of slices do not match")
	}
	c128.AxpyUnitaryTo(dst, alpha, s, y)
	return dst
}

// Arg generates cmplx.Phase, element-wise, for the elements of s, and stores in dst.
// Phase is in radians.
// It panics if the lengths of dst and s are not equal.
func Arg(dst []float64, s []complex128) []float64 {
	if len(dst) != len(s) {
		panic("cmplxs.Arg: length of the slices do not match")
	}
	for i, v := range s {
		dst[i] = cmplx.Phase(v)
	}
	return dst
}

// Conj generates cmplx.Conj, element-wise, for the elements of s, and stores in s.
func Conj(s []complex128) {
	for i, v := range s {
		s[i] = cmplx.Conj(v)
	}
	return
}

// ConjTo generates cmplx.Conj, element-wise, for the elements of s, and stores in dst.
// It panics if the lengths of dst and s are not equal.
func ConjTo(dst, s []complex128) []complex128 {
	if len(dst) != len(s) {
		panic("cmplxs.ConjTo: length of the slices do not match")
	}
	for i, v := range s {
		dst[i] = cmplx.Conj(v)
	}
	return dst
}

// CumProd finds the cumulative product of the first i elements in
// s and puts them in place into the ith element of the
// destination dst. A panic will occur if the lengths of arguments
// do not match.
//
// At the return of the function, dst[i] = s[i] * s[i-1] * s[i-2] * ...
func CumProd(dst, s []complex128) []complex128 {
	if len(dst) != len(s) {
		panic("cmplxs.CumProd: length of destination does not match length of the source")
	}
	if len(dst) == 0 {
		return dst
	}
	dst[0] = s[0]
	for i, v := range s[1:] {
		dst[i+1] = dst[i] * v
	}
	return dst
}

// CumSum finds the cumulative sum of the first i elements in
// s and puts them in place into the ith element of the
// destination dst. A panic will occur if the lengths of arguments
// do not match.
//
// At the return of the function, dst[i] = s[i] + s[i-1] + s[i-2] + ...
func CumSum(dst, s []complex128) []complex128 {
	if len(dst) != len(s) {
		panic("cmplxs.CumSum: length of destination does not match length of the source")
	}
	if len(dst) == 0 {
		return dst
	}
	dst[0] = s[0]
	for i, v := range s[1:] {
		dst[i+1] = dst[i] + v
	}
	return dst
}

// Deg converts radians to degrees, element-wise, for the elements of s, and stores in s.
func Deg(s []float64) {
	for i, v := range s {
		s[i] = v * 180 / math.Pi
	}
	return
}

// DegTo converts radians to degrees, element-wise, for the elements of s, and stores in dst.
// It panics if the lengths of dst and s are not equal.
func DegTo(dst, s []float64) []float64 {
	if len(dst) != len(s) {
		panic("cmplxs.DegTo: length of destination does not match length of the source")
	}
	for i, v := range s {
		dst[i] = v * 180 / math.Pi
	}
	return dst
}

// Distance computes the L-norm of s - t. See Norm for special cases.
// A panic will occur if the lengths of s and t do not match.
func Distance(s, t []complex128, L float64) float64 {
	if len(s) != len(t) {
		panic("cmplxs.Distance: slice lengths do not match")
	}
	if len(s) == 0 {
		return 0
	}
	if L == 2 {
		return L2DistanceUnitary(s, t)
	}
	var norm float64
	if L == 1 {
		for i, v := range s {
			norm += cmplx.Abs(t[i] - v)
		}
		return norm
	}
	if math.IsInf(L, 1) {
		for i, v := range s {
			absDiff := cmplx.Abs(t[i] - v)
			if absDiff > norm {
				norm = absDiff
			}
		}
		return norm
	}
	for i, v := range s {
		norm += math.Pow(cmplx.Abs(t[i]-v), L)
	}
	return math.Pow(norm, 1/L)
}

// Div performs element-wise division dst / s
// and stores the value in dst. It panics if the
// lengths of s and t are not equal.
func Div(dst, s []complex128) {
	if len(dst) != len(s) {
		panic("cmplxs.Div: slice lengths do not match")
	}
	for i, v := range s {
		dst[i] /= v
	}
}

// DivTo performs element-wise division s / t
// and stores the value in dst. It panics if the
// lengths of s, t, and dst are not equal.
func DivTo(dst, s, t []complex128) []complex128 {
	if len(s) != len(t) || len(dst) != len(t) {
		panic("cmplxs.DivTo: slice lengths do not match")
	}
	for i, v := range s {
		dst[i] = v / t[i]
	}
	return dst
}

// Dot computes the dot product of s1 and s2, i.e.
// sum_{i = 1}^N s1[i]*s2[i].
// A panic will occur if lengths of arguments do not match.
func Dot(s1, s2 []complex128, conj bool) complex128 {
	if len(s1) != len(s2) {
		panic("cmplxs.Dot: lengths of the slices do not match")
	}
	switch conj {
	case true:
		return c128.DotcUnitary(s1, s2)
	default:
		return c128.DotuUnitary(s1, s2)
	}
}

// Equal returns true if the slices have equal lengths and
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

// EqualApprox returns true if the slices have equal lengths and
// all element pairs have an absolute tolerance less than tol or a
// relative tolerance less than tol.
func EqualApprox(s1, s2 []complex128, tol float64) bool {
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

const minNormalFloat64 = 2.2250738585072014e-308

// EqualWithinAbs returns true if a and b have an absolute
// difference of less than tol.
func EqualWithinAbs(a, b complex128, tol float64) bool {
	return a == b || (math.Abs(real(a-b)) <= tol && math.Abs(imag(a-b)) <= tol)
}

// EqualWithinAbsOrRel returns true if a and b are equal to within
// the absolute tolerance.
func EqualWithinAbsOrRel(a, b complex128, absTol, relTol float64) bool {
	if EqualWithinAbs(a, b, absTol) {
		return true
	}
	return EqualWithinRel(a, b, relTol)
}

// EqualWithinRel returns true if the difference between a and b
// is not greater than tol times the greater value.
func EqualWithinRel(a, b complex128, tol float64) bool {
	// if it's equal, return immediately
	if a == b {
		return true
	}

	// check if real or imaginary numbers are infinity
	switch {
	case math.IsInf(real(a), 0) && !math.IsInf(real(b), 0):
		return false
	case !math.IsInf(real(a), 0) && math.IsInf(real(b), 0):
		return false
	case math.IsInf(imag(a), 0) && !math.IsInf(imag(b), 0):
		return false
	case !math.IsInf(imag(a), 0) && math.IsInf(imag(b), 0):
		return false
	}

	delta := a - b
	if math.Abs(real(delta)) <= minNormalFloat64 && math.Abs(imag(delta)) <= minNormalFloat64 {
		return math.Abs(real(delta)) <= tol*minNormalFloat64 && math.Abs(imag(delta)) <= tol*minNormalFloat64
	}

	// We depend on the division in this relationship to identify
	// infinities (we rely on the NaN to fail the test) otherwise
	// we compare Infs of the same sign and evaluate Infs as equal
	// independent of sign.
	// math.Abs(delta)/math.Max(math.Abs(a), math.Abs(b)) <= tol
	switch {
	case math.Max(math.Abs(real(a)), math.Abs(real(b))) == 0:
		return math.Abs(imag(delta))/math.Max(math.Abs(imag(a)), math.Abs(imag(b))) <= tol
	case math.Max(math.Abs(imag(a)), math.Abs(imag(b))) == 0:
		return math.Abs(real(delta))/math.Max(math.Abs(real(a)), math.Abs(real(b))) <= tol
	default:
		return math.Abs(real(delta))/math.Max(math.Abs(real(a)), math.Abs(real(b))) <= tol && math.Abs(imag(delta))/math.Max(math.Abs(imag(a)), math.Abs(imag(b))) <= tol
	}
}

// EqualWithinULP returns true if a and b are equal to within
// the specified number of floating point units in the last place.
func EqualWithinULP(a, b complex128, ulp uint) bool {
	if a == b {
		return true
	}
	if cmplx.IsNaN(a) || cmplx.IsNaN(b) {
		return false
	}
	if math.Signbit(real(a)) != math.Signbit(real(b)) && math.Signbit(imag(a)) != math.Signbit(imag(b)) {
		return math.Float64bits(math.Abs(real(a)))+math.Float64bits(math.Abs(real(b))) <= uint64(ulp) && math.Float64bits(math.Abs(imag(a)))+math.Float64bits(math.Abs(imag(b))) <= uint64(ulp)
	}
	return ulpDiff(math.Float64bits(real(a)), math.Float64bits(real(b))) <= uint64(ulp) && ulpDiff(math.Float64bits(imag(a)), math.Float64bits(imag(b))) <= uint64(ulp)
}

func ulpDiff(a, b uint64) uint64 {
	if a > b {
		return a - b
	}
	return b - a
}

// EqualLengths returns true if all of the slices have equal length,
// and false otherwise. Returns true if there are no input slices.
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
	return inds, errors.New("cmplxs.Find: insufficient elements found")
}

// HasNaN returns true if the slice s has any values that are NaN and false
// otherwise.
func HasNaN(s []complex128) bool {
	for _, v := range s {
		if cmplx.IsNaN(v) {
			return true
		}
	}
	return false
}

// Hermitian generates the hermitian (conjugate transpose) of s with r rows and c cols,
// and stores in s.
// It panics if the length of s is not equal to r * c.
func Hermitian(r, c int, s []complex128) {
	if (r * c) != len(s) {
		panic("cmplxs.Hermitian: rows & cols do not match with length of s")
	}
	dst := make([]complex128, len(s))
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			dst[i+r*j] = cmplx.Conj(s[i*c+j])
		}
	}
	for i, v := range dst {
		s[i] = v
	}
	return
}

// HermitianTo generates the hermitian (conjugate transpose) of s with r rows and c cols,
// and stores in dst.
// It panics if the length of s is not equal to r * c and if the length of s and dst are
// not equal.
func HermitianTo(r, c int, dst, s []complex128) []complex128 {
	if (r * c) != len(s) {
		panic("cmplxs.HermitianTo: rows & cols do not match with length of s")
	}
	if len(s) != len(dst) {
		panic("cmplxs.HermitianTo: length of s and dst do not match")
	}
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			dst[i+r*j] = cmplx.Conj(s[i*c+j])
		}
	}
	return dst
}

// Imag calculates imag(s) and stores in dst. dst is returned.
// It panics if the lengths of dst and s are not equal.
func Imag(dst []float64, s []complex128) []float64 {
	if len(s) != len(dst) {
		panic("cmplxs.Imag: length of s and dst do not match")
	}
	if len(s) == 0 {
		return dst
	}
	for i, v := range s {
		dst[i] = imag(v)
	}
	return dst
}

// L1Dist returns the 1-Norm of t - s
func L1Dist(s, t []complex128) float64 {
	var norm float64
	switch {
	case len(s) > len(t):
		for i, v := range t {
			norm += cmplx.Abs(v - s[i])
		}
	default:
		for i, v := range s {
			norm += cmplx.Abs(t[i] - v)
		}
	}
	return norm
}

// LinfDist returns the inf-Norm of t - s
func LinfDist(s, t []complex128) float64 {
	var norm float64
	if len(s) == 0 {
		return 0
	}
	norm = cmplx.Abs(t[0] - s[0])
	switch {
	case len(s) > len(t):
		for i, v := range t[1:] {
			absDiff := cmplx.Abs(v - s[i+1])
			if absDiff > norm || math.IsNaN(norm) {
				norm = absDiff
			}
		}
	default:
		for i, v := range s[1:] {
			absDiff := cmplx.Abs(t[i+1] - v)
			if absDiff > norm || math.IsNaN(norm) {
				norm = absDiff
			}
		}
	}
	return norm
}

// L1Norm returns the 1-Norm of x
func L1Norm(x []complex128) (sum float64) {
	for _, v := range x {
		sum += cmplx.Abs(v)
	}
	return sum
}

// L1NormInc returns the 1-Norm of incX incremented values of x
func L1NormInc(x []complex128, n, incX int) (sum float64) {
	for i := 0; i < n*incX; i += incX {
		sum += cmplx.Abs(x[i])
	}
	return sum
}

// L2NormUnitary returns the L2-norm of x.
func L2NormUnitary(x []complex128) (norm float64) {
	var scale float64
	sumSquares := 1.0
	for _, v := range x {
		if v == 0 {
			continue
		}
		absxi := cmplx.Abs(v)
		if math.IsNaN(absxi) {
			return math.NaN()
		}
		if scale < absxi {
			s := scale / absxi
			sumSquares = 1 + sumSquares*s*s
			scale = absxi
		} else {
			s := absxi / scale
			sumSquares += s * s
		}
	}
	if math.IsInf(scale, 1) {
		return math.Inf(1)
	}
	return scale * math.Sqrt(sumSquares)
}

// L2NormInc returns the L2-norm of x.
func L2NormInc(x []complex128, n, incX int) (norm float64) {
	var scale float64
	sumSquares := 1.0
	for ix := 0; ix < n*incX; ix += incX {
		val := x[ix]
		if val == 0 {
			continue
		}
		absxi := cmplx.Abs(val)
		if math.IsNaN(absxi) {
			return math.NaN()
		}
		if scale < absxi {
			s := scale / absxi
			sumSquares = 1 + sumSquares*s*s
			scale = absxi
		} else {
			s := absxi / scale
			sumSquares += s * s
		}
	}
	if math.IsInf(scale, 1) {
		return math.Inf(1)
	}
	return scale * math.Sqrt(sumSquares)
}

// L2DistanceUnitary returns the L2-norm of x-y.
func L2DistanceUnitary(x, y []complex128) (norm float64) {
	var scale float64
	sumSquares := 1.0
	for i, v := range x {
		v -= y[i]
		if v == 0 {
			continue
		}
		absxi := cmplx.Abs(v)
		if math.IsNaN(absxi) {
			return math.NaN()
		}
		if scale < absxi {
			s := scale / absxi
			sumSquares = 1 + sumSquares*s*s
			scale = absxi
		} else {
			s := absxi / scale
			sumSquares += s * s
		}
	}
	if math.IsInf(scale, 1) {
		return math.Inf(1)
	}
	return scale * math.Sqrt(sumSquares)
}

// Mul performs element-wise multiplication between dst
// and s and stores the value in dst. Panics if the
// lengths of s and t are not equal.
func Mul(dst, s []complex128) {
	if len(dst) != len(s) {
		panic("cmplxs: slice lengths do not match")
	}
	for i, val := range s {
		dst[i] *= val
	}
}

// MulTo performs element-wise multiplication between s
// and t and stores the value in dst. Panics if the
// lengths of s, t, and dst are not equal.
func MulTo(dst, s, t []complex128) []complex128 {
	if len(s) != len(t) || len(dst) != len(t) {
		panic("cmplxs: slice lengths do not match")
	}
	for i, val := range t {
		dst[i] = val * s[i]
	}
	return dst
}

// Norm returns the L norm of the slice S, defined as
// (sum_{i=1}^N s[i]^L)^{1/L}
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
	if L == 2 {
		return L2NormUnitary(s)
	}
	var norm float64
	if L == 1 {
		for _, val := range s {
			norm += cmplx.Abs(val)
		}
		return norm
	}
	if math.IsInf(L, 1) {
		for _, val := range s {
			norm = math.Max(norm, cmplx.Abs(val))
		}
		return norm
	}
	for _, val := range s {
		norm += math.Pow(cmplx.Abs(val), L)
	}
	return math.Pow(norm, 1/L)
}

// Polar generates cmplx.Polar, element-wise, for the elements of s, and returns slices.
// phase is in radians.
func Polar(r, p []float64, s []complex128) ([]float64, []float64) {
	if len(r) != len(p) || len(r) != len(s) {
		panic("cmplxs.Arg: length of the slices do not match")
	}
	for i, v := range s {
		r[i], p[i] = cmplx.Polar(v)
	}
	return r, p
}

// Prod returns the product of the elements of the slice.
// Returns 1 if len(s) = 0.
func Prod(s []complex128) complex128 {
	prod := complex(1.0, 0)
	for _, val := range s {
		prod *= val
	}
	return prod
}

// Rad converts degrees to radians for the elements of s and stores in s.
func Rad(s []float64) {
	for i, v := range s {
		s[i] = v * math.Pi / 180
	}
	return
}

// RadTo converts degrees to radians for the elements of s and stores in dst.
// dst is returned.
// It panics if the lengths of dst and s are not equal.
func RadTo(dst, s []float64) []float64 {
	if len(dst) != len(s) {
		panic("cmplxs.Real: length of s and dst do not match")
	}
	if len(s) == 0 {
		return dst
	}
	for i, v := range s {
		dst[i] = v * math.Pi / 180
	}
	return dst
}

// Real calculates real(s) and stores in dst. dst is returned.
// It panics if the lengths of dst and s are not equal.
func Real(dst []float64, s []complex128) []float64 {
	if len(s) != len(dst) {
		panic("cmplxs.Real: length of s and dst do not match")
	}
	if len(s) == 0 {
		return dst
	}
	for i, v := range s {
		dst[i] = real(v)
	}
	return dst
}

// Rect calculates cmplx.Rect(r, p) for the elements of r & p, and returns slice.
// phase is in radians.
func Rect(dst []complex128, r, p []float64) []complex128 {
	if len(r) != len(p) || len(r) != len(dst) {
		panic("cmplxs: Rect: dst, r & p slice lengths do not match")
	}
	for i, v := range r {
		dst[i] = cmplx.Rect(v, p[i])
	}
	return dst
}

// Reverse reverses the order of elements in the slice.
func Reverse(s []complex128) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

// ReverseTo reverses the order of elements in s and stores in dst.
// dst is returned.
// It panics if the lengths of dst and s are not equal.
func ReverseTo(dst, s []complex128) []complex128 {
	if len(s) != len(dst) {
		panic("cmplxs.ReverseTo: length of s and dst do not match")
	}
	for i := 0; i < len(s); i++ {
		dst[i] = s[len(s)-1-i]
	}
	return dst
}

// Round returns the half away from zero rounded value of x with prec precision.
//
// Special cases are:
// 	Round(±0) = +0
// 	Round(±Inf) = ±Inf
// 	Round(NaN) = NaN
func Round(x complex128, prec int) complex128 {
	if x == 0 {
		// Make sure zero is returned
		// without the negative bit set.
		return 0
	}
	// Fast path for positive precision on integers.
	if prec >= 0 && real(x) == math.Trunc(real(x)) && imag(x) == math.Trunc(imag(x)) {
		return x
	}
	pow := math.Pow10(prec)
	intermed := x * complex(pow, 0)
	if cmplx.IsInf(intermed) {
		return x
	}
	xre := real(intermed)
	xim := imag(intermed)
	if xre < 0 {
		xre = math.Ceil(xre - 0.5)
	} else {
		xre = math.Floor(xre + 0.5)
	}
	if xim < 0 {
		xim = math.Ceil(xim - 0.5)
	} else {
		xim = math.Floor(xim + 0.5)
	}

	x = complex(xre, xim)

	if x == 0 {
		return 0
	}

	return x / complex(pow, 0)
}

// RoundEven returns the half even rounded value of x with prec precision.
//
// Special cases are:
// 	RoundEven(±0) = +0
// 	RoundEven(±Inf) = ±Inf
// 	RoundEven(NaN) = NaN
func RoundEven(x complex128, prec int) complex128 {
	if x == 0 {
		// Make sure zero is returned
		// without the negative bit set.
		return 0
	}
	// Fast path for positive precision on integers.
	if prec >= 0 && real(x) == math.Trunc(real(x)) && imag(x) == math.Trunc(imag(x)) {
		return x
	}
	pow := math.Pow10(prec)
	intermed := x * complex(pow, 0)
	if cmplx.IsInf(intermed) {
		return x
	}
	xre := real(intermed)
	if isHalfway(xre) {
		correction, _ := math.Modf(math.Mod(xre, 2))
		xre += correction
		if xre > 0 {
			xre = math.Floor(xre)
		} else {
			xre = math.Ceil(xre)
		}
	} else {
		if xre < 0 {
			xre = math.Ceil(xre - 0.5)
		} else {
			xre = math.Floor(xre + 0.5)
		}
	}
	xim := imag(intermed)
	if isHalfway(xim) {
		correction, _ := math.Modf(math.Mod(xim, 2))
		xim += correction
		if xim > 0 {
			xim = math.Floor(xim)
		} else {
			xim = math.Ceil(xim)
		}
	} else {
		if xim < 0 {
			xim = math.Ceil(xim - 0.5)
		} else {
			xim = math.Floor(xim + 0.5)
		}
	}

	x = complex(xre, xim)
	if x == 0 {
		return 0
	}

	return x / complex(pow, 0)
}

// Copied from floats
func isHalfway(x float64) bool {
	_, frac := math.Modf(x)
	frac = math.Abs(frac)
	return frac == 0.5 || (math.Nextafter(frac, math.Inf(-1)) < 0.5 && math.Nextafter(frac, math.Inf(1)) > 0.5)
}

// Same returns true if the input slices have the same length and the all elements
// have the same value with NaN treated as the same.
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

// ScaleTo multiplies the elements in s by c and stores the result in dst.
func ScaleTo(dst []complex128, c complex128, s []complex128) []complex128 {
	if len(dst) != len(s) {
		panic("cmplxs: lengths of slices do not match")
	}
	if len(dst) > 0 {
		c128.ScalUnitaryTo(dst, c, s)
	}
	return dst
}

// Sub subtracts, element-wise, the elements of s from dst. Panics if
// the lengths of dst and s do not match.
func Sub(dst, s []complex128) {
	if len(dst) != len(s) {
		panic("cmplxs: length of the slices do not match")
	}
	c128.AxpyUnitaryTo(dst, -1, s, dst)
}

// SubTo subtracts, element-wise, the elements of t from s and
// stores the result in dst. Panics if the lengths of s, t and dst do not match.
func SubTo(dst, s, t []complex128) []complex128 {
	if len(s) != len(t) {
		panic("cmplxs: length of subtractor and subtractee do not match")
	}
	if len(dst) != len(s) {
		panic("cmplxs: length of destination does not match length of subtractor")
	}
	c128.AxpyUnitaryTo(dst, -1, t, s)
	return dst
}

// Sum returns the sum of the elements of the slice.
func Sum(s []complex128) complex128 {
	var sum complex128
	for _, v := range s {
		sum += v
	}
	return sum
}
