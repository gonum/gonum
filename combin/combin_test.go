// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package combin

import (
	"math/big"
	"testing"

	"github.com/gonum/floats"
)

// intSosMatch returns true if the two slices of slices are equal.
func intSosMatch(a, b [][]int) bool {
	if len(a) != len(b) {
		return false
	}
	for i, s := range a {
		if len(s) != len(b[i]) {
			return false
		}
		for j, v := range s {
			if v != b[i][j] {
				return false
			}
		}
	}
	return true
}

var binomialTests = []struct {
	n, k, ans int
}{
	{0, 0, 1},
	{5, 0, 1},
	{5, 1, 5},
	{5, 2, 10},
	{5, 3, 10},
	{5, 4, 5},
	{5, 5, 1},

	{6, 0, 1},
	{6, 1, 6},
	{6, 2, 15},
	{6, 3, 20},
	{6, 4, 15},
	{6, 5, 6},
	{6, 6, 1},

	{20, 0, 1},
	{20, 1, 20},
	{20, 2, 190},
	{20, 3, 1140},
	{20, 4, 4845},
	{20, 5, 15504},
	{20, 6, 38760},
	{20, 7, 77520},
	{20, 8, 125970},
	{20, 9, 167960},
	{20, 10, 184756},
	{20, 11, 167960},
	{20, 12, 125970},
	{20, 13, 77520},
	{20, 14, 38760},
	{20, 15, 15504},
	{20, 16, 4845},
	{20, 17, 1140},
	{20, 18, 190},
	{20, 19, 20},
	{20, 20, 1},
}

func TestBinomial(t *testing.T) {
	for cas, test := range binomialTests {
		ans := Binomial(test.n, test.k)
		if ans != test.ans {
			t.Errorf("Case %v: Binomial mismatch. Got %v, want %v.", cas, ans, test.ans)
		}
	}
	var (
		n    = 61
		want big.Int
		got  big.Int
	)
	for k := 0; k <= n; k++ {
		want.Binomial(int64(n), int64(k))
		got.SetInt64(int64(Binomial(n, k)))
		if want.Cmp(&got) != 0 {
			t.Errorf("Case n=%v,k=%v: Binomial mismatch for large n. Got %v, want %v.", n, k, got, want)
		}
	}
}

func TestGeneralizedBinomial(t *testing.T) {
	for cas, test := range binomialTests {
		ans := GeneralizedBinomial(float64(test.n), float64(test.k))
		if !floats.EqualWithinAbsOrRel(ans, float64(test.ans), 1e-14, 1e-14) {
			t.Errorf("Case %v: Binomial mismatch. Got %v, want %v.", cas, ans, test.ans)
		}
	}
}

func TestCombinations(t *testing.T) {
	for cas, test := range []struct {
		n, k int
		data [][]int
	}{
		{
			n:    1,
			k:    1,
			data: [][]int{{0}},
		},
		{
			n:    2,
			k:    1,
			data: [][]int{{0}, {1}},
		},
		{
			n:    2,
			k:    2,
			data: [][]int{{0, 1}},
		},
		{
			n:    3,
			k:    1,
			data: [][]int{{0}, {1}, {2}},
		},
		{
			n:    3,
			k:    2,
			data: [][]int{{0, 1}, {0, 2}, {1, 2}},
		},
		{
			n:    3,
			k:    3,
			data: [][]int{{0, 1, 2}},
		},
		{
			n:    4,
			k:    1,
			data: [][]int{{0}, {1}, {2}, {3}},
		},
		{
			n:    4,
			k:    2,
			data: [][]int{{0, 1}, {0, 2}, {0, 3}, {1, 2}, {1, 3}, {2, 3}},
		},
		{
			n:    4,
			k:    3,
			data: [][]int{{0, 1, 2}, {0, 1, 3}, {0, 2, 3}, {1, 2, 3}},
		},
		{
			n:    4,
			k:    4,
			data: [][]int{{0, 1, 2, 3}},
		},
	} {
		data := Combinations(test.n, test.k)
		if !intSosMatch(data, test.data) {
			t.Errorf("Cas %v: Generated combinations mismatch. Got %v, want %v.", cas, data, test.data)
		}
	}
}

func TestCombinationGenerator(t *testing.T) {
	for n := 0; n <= 10; n++ {
		for k := 1; k <= n; k++ {
			combinations := Combinations(n, k)
			cg := NewCombinationGenerator(n, k)
			genCombs := make([][]int, 0, len(combinations))
			for cg.Next() {
				genCombs = append(genCombs, cg.Combination(nil))
			}
			if !intSosMatch(combinations, genCombs) {
				t.Errorf("Combinations and generated combinations do not match. n = %v, k = %v", n, k)
			}
		}
	}
}
