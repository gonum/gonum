// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unit_test

import (
	"fmt"

	"gonum.org/v1/gonum/unit"
)

func ExampleNew() {
	// Create an acceleration of 3 m/s^2
	accel := unit.New(3.0, unit.Dimensions{unit.LengthDim: 1, unit.TimeDim: -2})
	fmt.Println(accel)

	// Output: 3 m s^-2
}

func ExampleNewDimension() {
	// Create a "trees" dimension
	// Typically, this should be used within an init function
	treeDim := unit.NewDimension("tree")
	countPerArea := unit.New(0.1, unit.Dimensions{treeDim: 1, unit.LengthDim: -2})
	fmt.Println(countPerArea)

	// Output: 0.1 tree m^-2
}

func Example_horsepower() {
	// One mechanical horsepower ≡ 33,000 ft-lbf/min.
	foot := unit.Length(0.3048)
	pound := unit.Mass(0.45359237)

	gravity := unit.New(9.80665, unit.Dimensions{unit.LengthDim: 1, unit.TimeDim: -2})
	poundforce := pound.Unit().Mul(gravity)

	hp := ((33000 * foot).Unit().Mul(poundforce)).Div(unit.Minute)
	fmt.Println("1 hp =", hp)

	watt := unit.New(1, unit.Dimensions{unit.MassDim: 1, unit.LengthDim: 2, unit.TimeDim: -3})
	fmt.Println("W is equivalent to hp?", unit.DimensionsMatch(hp, watt))

	// Output:
	//
	// 1 hp = 745.6998715822701 kg m^2 s^-3
	// W is equivalent to hp? true
}
