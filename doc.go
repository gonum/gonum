// Package unit provides a set of types to facilitate the use of
// the International System of Units (SI) in numerical codes.
//
// Unit provides functionality in two main ways.
//
// 1)
// Unit provides a number of types representing SI units. Each type is named for
// the unit it represents (Length, Mass, Pressure, etc.) and is
// based on a float64 with the value representing an amount in the
// corresponding SI base unit, derived unit, or other use accepted
// for use in SI (Kilogram, Meter, Pascal, etc.). A number of other
// SI constants are also defined as multiples of the base unit (Bar,
// electron volt, etc.). The use of these unit types will declare
// a representation of the unit in SI.
//
// 2)
// Unit provides the type Unit to help guarantee consistency
// of dimensionality when doing unit multiplication or division.
// Unit represents a dimensional value where the dimensions are
// not fixed.
package unit
