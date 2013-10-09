// Package unit provides a set of types and constants that facilitate
// the use of the International System of Units (SI).
//
// Unit provides two main functionalities.
//
// 1)
// It provides a number of types representing either an SI base unit
// or a common combination of base units, named for the unit it
// represents (Length, Mass, Pressure, etc.). Each type has
// a float64 as the underlying unit, and its value represents the
// number of that underlying unit (Kilogram, Meter, Pascal, etc.).
// For example,
//		height := 1.6 * unit.Meter
//		acc := unit.Acceleration(9.8)
// creates a variable named 'height' with a value of 1.6 meters, and
// a variable named 'acc' with a value of 9.8 meters per second squared.
// These types can be used to add compile-time safety to code. For
// example,
//		func UnitDensity(t unit.Temperature, pressure unit.Pressure) (unit.Density){
//			...
//		}
//		func main(){
//			t := 300 * unit.Kelvin
//			p := 5 * unit.Bar
//			rho := UnitDensity(p, t) // gives compile-time error
//		}
// gives a compile-time error (temperature type does not match pressure type)
// while the corresponding code using float64 runs without error.
//		func Float64Density(temperature, pressure float64) (float64){
//			...
//		}
//		func main(){
//			t := 300.0 // degrees kelvin
//			p := 50000.0 // Pascals
//			rho := Float64Density(p, t) // no error
//		}
// Many types have constants defined representing named SI units (Meter,
// Kilogram, etc. ) or SI derived units (Bar, Milligram, etc.). These are
// all defined as multiples of the base unit, so, for example, the following
// are euqivalent
//		l := 0.001 * unit.Meter
//		k := 1 * unit.Millimeter
//		j := unit.Length(0.001)
//
// 2)
// Unit provides the type "Unit", meant to represent a general dimensional
// value. unit.Unit can be used to help prevent errors of dimensionality
// when multiplying or dividing dimensional numbers. This package also
// provides the "Uniter" interface which is satisfied by any type which can
// be converted to a unit. New varibles of type Unit can be created with
// the NewUnit function, and custom dimensions can be created with the use of
// NewDimension. Please see the rest of the package docs for more
// details on usage.
//
// Please note that Unit cannot catch all errors related to dimensionality.
// Different physical ideas are sometimes expressed with the same dimensions
// and Unit is incapable of catcing these mismatches. For example, energy and
// torque are both expressed as force times distance (Newton-meters in SI),
// but it is wrong to say that a torque of 10 N-m == 10 J.
package unit
