package unit

import "testing"

type UnitStructer interface {
	UnitStruct() *UnitStruct
}

type UnitStruct struct {
	current     int
	length      int
	luminosity  int
	mass        int
	temperature int
	time        int
	chemamt     int // For mol
	value       float64
}

// Check if the dimensions of two units are the same
func DimensionsMatchStruct(aU, bU UnitStructer) bool {
	a := aU.UnitStruct()
	b := bU.UnitStruct()
	if a.length != b.length {
		return false
	}
	if a.time != b.time {
		return false
	}
	if a.mass != b.mass {
		return false
	}
	if a.current != b.current {
		return false
	}
	if a.temperature != b.temperature {
		return false
	}
	if a.luminosity != b.luminosity {
		return false
	}
	if a.chemamt != b.chemamt {
		return false
	}
	return true
}

func (u *UnitStruct) UnitStruct() *UnitStruct {
	return u
}

func (u *UnitStruct) Add(aU UnitStructer) *UnitStruct {
	a := aU.UnitStruct()
	if !DimensionsMatchStruct(a, u) {
		panic("dimension mismatch")
	}
	u.value += a.value
	return u
}

func (u *UnitStruct) Mul(aU UnitStructer) *UnitStruct {
	a := aU.UnitStruct()
	u.length += a.length
	u.time += a.time
	u.mass += a.mass
	u.current += a.current
	u.temperature += a.temperature
	u.luminosity += a.luminosity
	u.chemamt += a.chemamt
	u.value *= a.value
	return u
}

var u3 *UnitStruct

func BenchmarkAddStruct(b *testing.B) {
	u1 := &UnitStruct{current: 1, chemamt: 5, value: 10}
	u2 := &UnitStruct{current: 1, chemamt: 5, value: 100}
	for i := 0; i < b.N; i++ {
		u2.Add(u1)
	}
}

func BenchmarkMulStruct(b *testing.B) {
	u1 := &UnitStruct{current: 1, chemamt: 5, value: 10}
	u2 := &UnitStruct{mass: 1, time: 1, value: 100}
	for i := 0; i < b.N; i++ {
		u2.Mul(u1)
	}
}

type UnitMapper interface {
	UnitMap() *UnitMap
}

type dimensionMap int

const (
	LengthB dimensionMap = iota
	TimeB
	MassB
	CurrentB
	TemperatureB
	LuminosityB
	ChemAmtB
)

type UnitMap struct {
	dimension map[dimensionMap]int
	value     float64
}

// Check if the dimensions of two units are the same
func DimensionsMatchMap(aU, bU UnitMapper) bool {
	a := aU.UnitMap()
	b := bU.UnitMap()
	if len(a.dimension) != len(b.dimension) {
		panic("Unequal dimension")
	}
	for key, dimA := range a.dimension {
		dimB, ok := b.dimension[key]
		if !ok || dimA != dimB {
			panic("Unequal dimension")
		}
	}
	return true
}

func (u *UnitMap) UnitMap() *UnitMap {
	return u
}

func (u *UnitMap) Add(aU UnitMapper) *UnitMap {
	a := aU.UnitMap()
	if !DimensionsMatchMap(a, u) {
		panic("dimension mismatch")
	}
	u.value += a.value
	return u
}

func (u *UnitMap) Mul(aU UnitMapper) *UnitMap {
	a := aU.UnitMap()
	for key, val := range a.dimension {
		u.dimension[key] += val
	}
	u.value *= a.value
	return u
}

func BenchmarkAddFloat(b *testing.B) {
	a := 0.0
	c := 10.0
	for i := 0; i < b.N; i++ {
		a += c
	}
}

func BenchmarkMulFloat(b *testing.B) {
	a := 0.0
	c := 10.0
	for i := 0; i < b.N; i++ {
		a *= c
	}
}

func BenchmarkAddMapSmall(b *testing.B) {
	u1 := &UnitMap{value: 10}
	u1.dimension = make(map[dimensionMap]int)
	u1.dimension[CurrentB] = 1
	u1.dimension[ChemAmtB] = 5

	u2 := &UnitMap{value: 10}
	u2.dimension = make(map[dimensionMap]int)
	u2.dimension[CurrentB] = 1
	u2.dimension[ChemAmtB] = 5
	for i := 0; i < b.N; i++ {
		u2.Add(u1)
	}
}

func BenchmarkMulMapSmallDiff(b *testing.B) {
	u1 := &UnitMap{value: 10}
	u1.dimension = make(map[dimensionMap]int)
	u1.dimension[LengthB] = 1

	u2 := &UnitMap{value: 10}
	u2.dimension = make(map[dimensionMap]int)
	u2.dimension[MassB] = 1
	for i := 0; i < b.N; i++ {
		u2.Mul(u1)
	}
}

func BenchmarkMulMapSmallSame(b *testing.B) {
	u1 := &UnitMap{value: 10}
	u1.dimension = make(map[dimensionMap]int)
	u1.dimension[LengthB] = 1

	u2 := &UnitMap{value: 10}
	u2.dimension = make(map[dimensionMap]int)
	u2.dimension[LengthB] = 2
	for i := 0; i < b.N; i++ {
		u2.Mul(u1)
	}
}

func BenchmarkMulMapLargeDiff(b *testing.B) {
	u1 := &UnitMap{value: 10}
	u1.dimension = make(map[dimensionMap]int)
	u1.dimension[LengthB] = 1
	u1.dimension[MassB] = 1
	u1.dimension[ChemAmtB] = 1
	u1.dimension[TemperatureB] = 1
	u1.dimension[LuminosityB] = 1
	u1.dimension[TimeB] = 1
	u1.dimension[CurrentB] = 1

	u2 := &UnitMap{value: 10}
	u2.dimension = make(map[dimensionMap]int)
	u2.dimension[MassB] = 1
	for i := 0; i < b.N; i++ {
		u2.Mul(u1)
	}
}

func BenchmarkMulMapLargeSame(b *testing.B) {
	u1 := &UnitMap{value: 10}
	u1.dimension = make(map[dimensionMap]int)
	u1.dimension[LengthB] = 2
	u1.dimension[MassB] = 2
	u1.dimension[ChemAmtB] = 2
	u1.dimension[TemperatureB] = 2
	u1.dimension[LuminosityB] = 2
	u1.dimension[TimeB] = 2
	u1.dimension[CurrentB] = 2

	u2 := &UnitMap{value: 10}
	u2.dimension = make(map[dimensionMap]int)
	u2.dimension[LengthB] = 3
	u2.dimension[MassB] = 3
	u2.dimension[ChemAmtB] = 3
	u2.dimension[TemperatureB] = 3
	u2.dimension[LuminosityB] = 3
	u2.dimension[TimeB] = 3
	u2.dimension[CurrentB] = 3
	for i := 0; i < b.N; i++ {
		u2.Mul(u1)
	}
}
