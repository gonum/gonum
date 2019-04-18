// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bound

import "math"

// Bound represents [Min, Max] bounds.
type Bound struct {
	Min, Max float64
}

// IsValid returns whether the bound is valid
func (b Bound) IsValid() bool {
	return b.Min <= b.Max
}

// Intersection returns the intersection of the input bounds if possible.
// Otherwise a NaN Bound is returned.
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
// Otherwise a NaN Bound is returned.
func Union(bounds ...Bound) Bound {
	if len(bounds) == 0 {
		return Bound{Min: math.NaN(), Max: math.NaN()}
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
