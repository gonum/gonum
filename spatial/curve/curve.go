// Package curve defines space filling curves. A space filling curve is a curve
// who's range contains the entirety of a finite k-dimensional space. Space
// filling curves can be used to map between a linear space and a 2D, 3D, or 4D
// space.
package curve

// Point is a point on a curve.
type Point uint64

// SpaceFilling is a space filling curve.
type SpaceFilling interface {
	// Size returns the spatial size of the curve. If the X, Y, and Z
	// coordinates of spatial points mapped by a 3 dimensional curve all lie
	// within [0, x), [0, y), and [0, z), respectively, the size of that curve
	// is (x, y, z).
	Size() []int

	// Curve returns the curve coordinate of V.
	Curve(v ...int) Point

	// Space returns the spatial coordinates of D.
	Space(d Point) []int
}
