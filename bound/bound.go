// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bound

import (
	"math"
	"sort"
)

// Bound represents [Min, Max] bounds.
type Bound struct {
	Min, Max float64
}

// Len returns the length of the bound
func (b Bound) Len() float64 {
	if math.IsNaN(b.Min) && math.IsNaN(b.Max) {
		return 0
	}
	return b.Max - b.Min
}

// IsValid returns whether the bound is valid. A valid bound will have
// the minimum less than or equal to the maximum.
func (b Bound) IsValid() bool {
	return b.Min <= b.Max
}

// JaccardIndex returns a jaccard index between bounds. The result is only valid for
// sets of bounds that have a contiguous union.
// If all sets of inputs bounds are empty, the result is 1.
func JaccardIndex(bounds ...Bound) float64 {
	if len(bounds) == 0 {
		return math.NaN()
	}
	u := Union(bounds...)
	if !u.IsValid() {
		return math.NaN()
	}
	i := Intersection(bounds...)
	di := i.Len()
	du := u.Len()

	if di == 0 {
		// If the union is empty, the intersection must
		// also be empty. J(Ø, Ø) is 1 by definition.
		if du == 0 {
			return 1
		}

		if !math.IsNaN(du) {
			return 0
		}
		return math.NaN()
	}

	return di / du
}

// Intersection returns the intersection of the input bounds. If the
// intersection is empty an invalid Bound is returned.
func Intersection(bounds ...Bound) Bound {
	if len(bounds) == 0 {
		return Bound{Min: math.NaN(), Max: math.NaN()}
	}

	intersection := Bound{Min: bounds[0].Min, Max: bounds[0].Max}
	for _, b := range bounds[1:] {
		intersection.Min = math.Max(intersection.Min, b.Min)
		intersection.Max = math.Min(intersection.Max, b.Max)
	}

	if !intersection.IsValid() {
		return Bound{Min: math.NaN(), Max: math.NaN()}
	}

	return intersection
}

// Union returns the contiguous union of the input bounds if possible.
// Otherwise an invalid Bound is returned. If bounds is a slice of Bound
// and is not sorted, the order of elements will be changed so that
// they are ordered ascending by Min.
func Union(bounds ...Bound) Bound {
	if len(bounds) == 0 {
		return Bound{Min: math.NaN(), Max: math.NaN()}
	}
	if len(bounds) > 1 && !sort.IsSorted(byMin(bounds)) {
		sort.Sort(byMin(bounds))
	}

	union := Bound{Min: bounds[0].Min, Max: bounds[0].Max}
	for _, b := range bounds[1:] {
		if b.Max < union.Min || union.Max < b.Min {
			return Bound{Min: math.NaN(), Max: math.NaN()}
		}
		union.Min = math.Min(union.Min, b.Min)
		union.Max = math.Max(union.Max, b.Max)
	}

	if !union.IsValid() {
		return Bound{Min: math.NaN(), Max: math.NaN()}
	}

	return union
}

type byMin []Bound

func (b byMin) Len() int           { return len(b) }
func (b byMin) Less(i, j int) bool { return b[i].Min < b[j].Max }
func (b byMin) Swap(i, j int)      { b[i], b[j] = b[j], b[i] }
