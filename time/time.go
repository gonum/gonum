package unit

type Time float64

const (
	JulianYear Time = 365.24
	Second     Time = 1.0
	Picosecond Time = 1E-12
	Minute     Time = 60.0
	Hour       Time = 24.0 * Minute
	//Year       Time = 365.24 * Hour
)

func (t Time) Unit() *Unit {
	return CreateUnit(float64(t), &Dimensions{Time: 1})
}

func (t Time) Picoseconds() float64 {
	return float64(t) / float64(Picosecond)
}
