package time

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

func (t Time) In(t2 Time) float64 {
	return float64(t) / float64(t2)
}

// So it can implement a timer interface
func (t Time) Time() Time {
	return l
}
