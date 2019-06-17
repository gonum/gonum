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

// Intersection returns a Bound that is the intersection of the input bounds.
// If the intersection is empty or if the input length is zero,
// then the NaN Bound will be returned.
func Intersection(bounds ...Bound) Bound {
	if len(bounds) == 0 {
		return Bound{Min: math.NaN(), Max: math.NaN()}
	}

	ret := Bound{Min: bounds[0].Min, Max: bounds[0].Max}
	for _, b := range bounds[1:] {
		ret.Min = math.Max(ret.Min, b.Min)
		ret.Max = math.Min(ret.Max, b.Max)
	}

	if !ret.IsValid() {
		return Bound{Min: math.NaN(), Max: math.NaN()}
	}

	return ret
}
