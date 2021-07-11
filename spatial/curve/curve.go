package curve

// Point is a point on a curve.
type Point uint64

// SpaceFilling is a space filling curve.
type SpaceFilling interface {
	// Curve maps a point in space to a point on the curve.
	Curve(...int) Point

	// Space maps a point on the curve to a point in space.
	Space(Point) []int
}
