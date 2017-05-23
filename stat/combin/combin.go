// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package combin implements routines involving combinatorics (permutations,
// combinations, etc.).
package combin

import "math"

const (
	badNegInput = "combin: negative input"
	badSetSize  = "combin: n < k"
	badInput    = "combin: wrong input slice length"
)

// Binomial returns the binomial coefficient of (n,k), also commonly referred to
// as "n choose k".
//
// The binomial coefficient, C(n,k), is the number of unordered combinations of
// k elements in a set that is n elements big, and is defined as
//
//  C(n,k) = n!/((n-k)!k!)
//
// n and k must be non-negative with n >= k, otherwise Binomial will panic.
// No check is made for overflow.
func Binomial(n, k int) int {
	if n < 0 || k < 0 {
		panic(badNegInput)
	}
	if n < k {
		panic(badSetSize)
	}
	// (n,k) = (n, n-k)
	if k > n/2 {
		k = n - k
	}
	b := 1
	for i := 1; i <= k; i++ {
		b = (n - k + i) * b / i
	}
	return b
}

// GeneralizedBinomial returns the generalized binomial coefficient of (n, k),
// defined as
//  Γ(n+1) / (Γ(k+1) Γ(n-k+1))
// where Γ is the Gamma function. GeneralizedBinomial is useful for continuous
// relaxations of the binomial coefficient, or when the binomial coefficient value
// may overflow int. In the latter case, one may use math/big for an exact
// computation.
//
// n and k must be non-negative with n >= k, otherwise GeneralizedBinomial will panic.
func GeneralizedBinomial(n, k float64) float64 {
	return math.Exp(LogGeneralizedBinomial(n, k))
}

// LogGeneralizedBinomial returns the log of the generalized binomial coefficient.
// See GeneralizedBinomial for more information.
func LogGeneralizedBinomial(n, k float64) float64 {
	if n < 0 || k < 0 {
		panic(badNegInput)
	}
	if n < k {
		panic(badSetSize)
	}
	a, _ := math.Lgamma(n + 1)
	b, _ := math.Lgamma(k + 1)
	c, _ := math.Lgamma(n - k + 1)
	return a - b - c
}

// CombinationGenerator generates combinations iteratively. Combinations may be
// called to generate all combinations collectively.
type CombinationGenerator struct {
	n         int
	k         int
	previous  []int
	remaining int
}

// NewCombinationGenerator returns a CombinationGenerator for generating the
// combinations of k elements from a set of size n.
//
// n and k must be non-negative with n >= k, otherwise NewCombinationGenerator
// will panic.
func NewCombinationGenerator(n, k int) *CombinationGenerator {
	return &CombinationGenerator{
		n:         n,
		k:         k,
		remaining: Binomial(n, k),
	}
}

// Next advances the iterator if there are combinations remaining to be generated,
// and returns false if all combinations have been generated. Next must be called
// to initialize the first value before calling Combination or Combination will
// panic. The value returned by Combination is only changed during calls to Next.
func (c *CombinationGenerator) Next() bool {
	if c.remaining <= 0 {
		// Next is called before combination, so c.remaining is set to zero before
		// Combination is called. Thus, Combination cannot panic on zero, and a
		// second sentinel value is needed.
		c.remaining = -1
		return false
	}
	if c.previous == nil {
		c.previous = make([]int, c.k)
		for i := range c.previous {
			c.previous[i] = i
		}
	} else {
		nextCombination(c.previous, c.n, c.k)
	}
	c.remaining--
	return true
}

// Combination generates the next combination. If next is non-nil, it must have
// length k and the result will be stored in-place into combination. If combination
// is nil a new slice will be allocated and returned. If all of the combinations
// have already been constructed (Next() returns false), Combination will panic.
//
// Next must be called to initialize the first value before calling Combination
// or Combination will panic. The value returned by Combination is only changed
// during calls to Next.
func (c *CombinationGenerator) Combination(combination []int) []int {
	if c.remaining == -1 {
		panic("combin: all combinations have been generated")
	}
	if c.previous == nil {
		panic("combin: Combination called before Next")
	}
	if combination == nil {
		combination = make([]int, c.k)
	}
	if len(combination) != c.k {
		panic(badInput)
	}
	copy(combination, c.previous)
	return combination
}

// Combinations generates all of the combinations of k elements from a
// set of size n. The returned slice has length Binomial(n,k) and each inner slice
// has length k.
//
// n and k must be non-negative with n >= k, otherwise Combinations will panic.
//
// CombinationGenerator may alternatively be used to generate the combinations
// iteratively instead of collectively.
func Combinations(n, k int) [][]int {
	combins := Binomial(n, k)
	data := make([][]int, combins)
	if len(data) == 0 {
		return data
	}
	data[0] = make([]int, k)
	for i := range data[0] {
		data[0][i] = i
	}
	for i := 1; i < combins; i++ {
		next := make([]int, k)
		copy(next, data[i-1])
		nextCombination(next, n, k)
		data[i] = next
	}
	return data
}

// nextCombination generates the combination after s, overwriting the input value.
func nextCombination(s []int, n, k int) {
	for j := k - 1; j >= 0; j-- {
		if s[j] == n+j-k {
			continue
		}
		s[j]++
		for l := j + 1; l < k; l++ {
			s[l] = s[j] + l - j
		}
		break
	}
}
