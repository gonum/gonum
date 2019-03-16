// Copyright ©2013 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package unit

import (
	"fmt"
	"math"
	"testing"
)

var formatTests = []struct {
	unit   Uniter
	format string
	expect string
}{
	{New(9.81, Dimensions{MassDim: 1, TimeDim: -2}), "%f", "9.810000 kg s^-2"},
	{New(9.81, Dimensions{MassDim: 1, TimeDim: -2}), "%1.f", "10 kg s^-2"},
	{New(9.81, Dimensions{MassDim: 1, TimeDim: -2}), "%10.f", "10 kg s^-2"},
	{New(9.81, Dimensions{MassDim: 1, TimeDim: -2}), "%.1f", "9.8 kg s^-2"},
	{New(9.81, Dimensions{MassDim: 1, TimeDim: -2}), "%20f", "    9.810000 kg s^-2"},
	{New(9.81, Dimensions{MassDim: 1, TimeDim: -2}), "%1f", "9.810000 kg s^-2"},
	{New(9.81, Dimensions{MassDim: 1, TimeDim: -2, LengthDim: 0}), "%f", "9.810000 kg s^-2"},
	{New(6.62606957e-34, Dimensions{MassDim: 2, TimeDim: -1}), "%e", "6.626070e-34 kg^2 s^-1"},
	{New(6.62606957e-34, Dimensions{MassDim: 2, TimeDim: -1}), "%20.3e", " 6.626e-34 kg^2 s^-1"},
	{New(6.62606957e-34, Dimensions{MassDim: 2, TimeDim: -1}), "%.3e", "6.626e-34 kg^2 s^-1"},
	{New(6.62606957e-34, Dimensions{MassDim: 2, TimeDim: -1}), "%25e", "   6.626070e-34 kg^2 s^-1"},
	{New(6.62606957e-34, Dimensions{MassDim: 2, TimeDim: -1}), "%1e", "6.626070e-34 kg^2 s^-1"},
	{New(6.62606957e-34, Dimensions{MassDim: 2, TimeDim: -1}), "%v", "6.62606957e-34 kg^2 s^-1"},
	{New(6.62606957e-34, Dimensions{MassDim: 2, TimeDim: -1}), "%s", "%!s(*Unit=6.62606957e-34 kg^2 s^-1)"},
	{Dimless(math.E), "%v", "2.718281828459045"},
	{Dimless(math.E), "%.2v", "2.7"},
	{Dimless(math.E), "%6.2v", "   2.7"},
	{Dimless(math.E), "%20v", "   2.718281828459045"},
	{Dimless(math.E), "%1v", "2.718281828459045"},
	{Dimless(math.E), "%#v", "unit.Dimless(2.718281828459045)"},
	{Dimless(math.E), "%s", "%!s(unit.Dimless=2.718281828459045)"},
	{Angle(1), "%v", "1 rad"},
	{Angle(0.45), "%.1v", "0.5 rad"},
	{Angle(0.45), "%10.1v", "   0.5 rad"},
	{Angle(0.45), "%10v", "  0.45 rad"},
	{Angle(0.45), "%1v", "0.45 rad"},
	{Angle(1), "%#v", "unit.Angle(1)"},
	{Angle(1), "%s", "%!s(unit.Angle=1 rad)"},
	{Current(1), "%v", "1 A"},
	{Current(0.45), "%.1v", "0.5 A"},
	{Current(0.45), "%10.1v", "     0.5 A"},
	{Current(0.45), "%10v", "    0.45 A"},
	{Current(0.45), "%1v", "0.45 A"},
	{Current(1), "%#v", "unit.Current(1)"},
	{Current(1), "%s", "%!s(unit.Current=1 A)"},
	{LuminousIntensity(1), "%v", "1 cd"},
	{LuminousIntensity(0.45), "%.1v", "0.5 cd"},
	{LuminousIntensity(0.45), "%10.1v", "    0.5 cd"},
	{LuminousIntensity(0.45), "%10v", "   0.45 cd"},
	{LuminousIntensity(0.45), "%1v", "0.45 cd"},
	{LuminousIntensity(1), "%#v", "unit.LuminousIntensity(1)"},
	{LuminousIntensity(1), "%s", "%!s(unit.LuminousIntensity=1 cd)"},
	{Mass(1), "%v", "1 kg"},
	{Mass(0.45), "%.1v", "0.5 kg"},
	{Mass(0.45), "%10.1v", "    0.5 kg"},
	{Mass(0.45), "%10v", "   0.45 kg"},
	{Mass(0.45), "%1v", "0.45 kg"},
	{Mass(1), "%#v", "unit.Mass(1)"},
	{Mass(1), "%s", "%!s(unit.Mass=1 kg)"},
	{Mole(1), "%v", "1 mol"},
	{Mole(0.45), "%.1v", "0.5 mol"},
	{Mole(0.45), "%10.1v", "   0.5 mol"},
	{Mole(0.45), "%10v", "  0.45 mol"},
	{Mole(0.45), "%1v", "0.45 mol"},
	{Mole(1), "%#v", "unit.Mole(1)"},
	{Mole(1), "%s", "%!s(unit.Mole=1 mol)"},
	{Length(1.61619926e-35), "%v", "1.61619926e-35 m"},
	{Length(1.61619926e-35), "%.2v", "1.6e-35 m"},
	{Length(1.61619926e-35), "%10.2v", " 1.6e-35 m"},
	{Length(1.61619926e-35), "%20v", "    1.61619926e-35 m"},
	{Length(1.61619926e-35), "%1v", "1.61619926e-35 m"},
	{Length(1.61619926e-35), "%#v", "unit.Length(1.61619926e-35)"},
	{Length(1.61619926e-35), "%s", "%!s(unit.Length=1.61619926e-35 m)"},
	{Temperature(15.2), "%v", "15.2 K"},
	{Temperature(15.2), "%.2v", "15 K"},
	{Temperature(15.2), "%6.2v", "  15 K"},
	{Temperature(15.2), "%10v", "    15.2 K"},
	{Temperature(15.2), "%1v", "15.2 K"},
	{Temperature(15.2), "%#v", "unit.Temperature(15.2)"},
	{Temperature(15.2), "%s", "%!s(unit.Temperature=15.2 K)"},
	{Time(15.2), "%v", "15.2 s"},
	{Time(15.2), "%.2v", "15 s"},
	{Time(15.2), "%6.2v", "  15 s"},
	{Time(15.2), "%10v", "    15.2 s"},
	{Time(15.2), "%1v", "15.2 s"},
	{Time(15.2), "%#v", "unit.Time(15.2)"},
	{Time(15.2), "%s", "%!s(unit.Time=15.2 s)"},
}

func TestFormat(t *testing.T) {
	for _, ts := range formatTests {
		if r := fmt.Sprintf(ts.format, ts.unit); r != ts.expect {
			t.Errorf("Format %q: got: %q expected: %q", ts.format, r, ts.expect)
		}
	}
}

func TestGoStringFormat(t *testing.T) {
	expect1 := `&unit.Unit{dimensions:unit.Dimensions{4:2, 7:-1}, value:6.62606957e-34}`
	expect2 := `&unit.Unit{dimensions:unit.Dimensions{7:-1, 4:2}, value:6.62606957e-34}`
	if r := fmt.Sprintf("%#v", New(6.62606957e-34, Dimensions{MassDim: 2, TimeDim: -1})); r != expect1 && r != expect2 {
		t.Errorf("Format %q: got: %q expected: %q", "%#v", r, expect1)
	}
}

var initializationTests = []struct {
	unit     *Unit
	expValue float64
	expMap   map[Dimension]int
}{
	{New(9.81, Dimensions{MassDim: 1, TimeDim: -2}), 9.81, Dimensions{MassDim: 1, TimeDim: -2}},
	{New(9.81, Dimensions{MassDim: 1, TimeDim: -2, LengthDim: 0, CurrentDim: 0}), 9.81, Dimensions{MassDim: 1, TimeDim: -2}},
}

func TestInitialization(t *testing.T) {
	for _, ts := range initializationTests {
		if ts.expValue != ts.unit.value {
			t.Errorf("Value wrong on initialization: got: %v expected: %v", ts.unit.value, ts.expValue)
		}
		if len(ts.expMap) != len(ts.unit.dimensions) {
			t.Errorf("Map mismatch: got: %#v expected: %#v", ts.unit.dimensions, ts.expMap)
		}
		for key, val := range ts.expMap {
			if ts.unit.dimensions[key] != val {
				t.Errorf("Map mismatch: got: %#v expected: %#v", ts.unit.dimensions, ts.expMap)
			}
		}
	}
}

var dimensionEqualityTests = []struct {
	name        string
	a           Uniter
	b           Uniter
	shouldMatch bool
}{
	{"same_empty", New(1.0, Dimensions{}), New(1.0, Dimensions{}), true},
	{"same_one", New(1.0, Dimensions{TimeDim: 1}), New(1.0, Dimensions{TimeDim: 1}), true},
	{"same_mult", New(1.0, Dimensions{TimeDim: 1, LengthDim: -2}), New(1.0, Dimensions{TimeDim: 1, LengthDim: -2}), true},
	{"diff_one_empty", New(1.0, Dimensions{}), New(1.0, Dimensions{TimeDim: 1, LengthDim: -2}), false},
	{"diff_same_dim", New(1.0, Dimensions{TimeDim: 1}), New(1.0, Dimensions{TimeDim: 2}), false},
	{"diff_same_pow", New(1.0, Dimensions{LengthDim: 1}), New(1.0, Dimensions{TimeDim: 1}), false},
	{"diff_numdim", New(1.0, Dimensions{TimeDim: 1, LengthDim: 2}), New(1.0, Dimensions{TimeDim: 2}), false},
	{"diff_one_same_dim", New(1.0, Dimensions{LengthDim: 1, TimeDim: 1}), New(1.0, Dimensions{LengthDim: 1, TimeDim: 2}), false},
}

func TestDimensionEquality(t *testing.T) {
	for _, ts := range dimensionEqualityTests {
		if DimensionsMatch(ts.a, ts.b) != ts.shouldMatch {
			t.Errorf("Dimension comparison incorrect for case %s. got: %v, expected: %v", ts.name, !ts.shouldMatch, ts.shouldMatch)
		}
	}
}

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

var sink float64

func BenchmarkAddFloat(b *testing.B) {
	sink = 0
	c := 10.0
	for i := 0; i < b.N; i++ {
		sink += c
	}
}

func BenchmarkMulFloat(b *testing.B) {
	sink = 0
	c := 10.0
	for i := 0; i < b.N; i++ {
		sink *= c
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
