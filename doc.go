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
// This is not perfect -- something about newton-meter vs. Joule
// If one wants to add a new dimensionality NewUnit should be called
// to guarantee non-overlap
// Should not use integer literals, only the values
// Other packages can create the new units in their init functions.
// This is so two packages will not accidtally use the same unit
package unit
