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
//		func UnitAtmosphericConditions(t unit.Temperature, pressure unit.Pressure){
//			...
//		}
//		func main(){
//			t := 300 * unit.Kelvin
//			p := 5 * unit.Bar
//			UnitAtmosphericConditions(p, t) // gives compile-time error
//		}
// gives a compile-time error (temperature type does not match pressure type)
// while the corresponding code using float64 runs without error.
//		func Float64AtmosphericConditions(temperature, pressure float64){
//			...
//		}
//		func main(){
//			t := 300
//			p := 50000 / Pa
//			Float64AtmosphericConditions(p, t) // no error
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
// Unit provides the type Unit to help guarantee consistency
// of dimensionality when doing unit multiplication or division.
// Unit represents a dimensional value where the dimensions are
// not fixed.
// This is not perfect -- something about newton-meter vs. Joule
// If one wants to add a new dimensionality NewUnit should be called
// to guarantee non-overlap
// Should not use integer literals, only the values
// Other packages can create the new units in their init functions.
// This is so two packages will not accidtally use the same unit
package unit
