package main

import (
	"go/format"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

type Unit struct {
	Name          string
	Offset        int    // From normal (for example, mass base unit is kg, not kg)
	PrintString   string // print string for the unit (kg for mass)
	ExtraConstant []Constant
	Suffix        string
	TypeComment   string // Text to comment the type
	Dimensions    []Dimension
	ErForm        string //For Xxxer interface
}

type Dimension struct {
	Name  string
	Power int
}

const (
	TimeName   string = "TimeDim"
	LengthName string = "LengthDim"
	MassName   string = "MassDim"
)

type Constant struct {
	Name  string
	Value string
}

type Prefix struct {
	Name  string
	Power int
}

var Prefixes = []Prefix{
	{
		Name:  "Yotta",
		Power: 24,
	},
	{
		Name:  "Zetta",
		Power: 21,
	},
	{
		Name:  "Exa",
		Power: 18,
	},
	{
		Name:  "Peta",
		Power: 15,
	},
	{
		Name:  "Tera",
		Power: 12,
	},
	{
		Name:  "Giga",
		Power: 9,
	},
	{
		Name:  "Mega",
		Power: 6,
	},
	{
		Name:  "Kilo",
		Power: 3,
	},
	{
		Name:  "Hecto",
		Power: 2,
	},
	{
		Name:  "Deca",
		Power: 1,
	},
	{
		Name:  "",
		Power: 0,
	},
	{
		Name:  "Deci",
		Power: -1,
	},
	{
		Name:  "Centi",
		Power: -2,
	},
	{
		Name:  "Milli",
		Power: -3,
	},
	{
		Name:  "Micro",
		Power: -6,
	},
	{
		Name:  "Nano",
		Power: -9,
	},
	{
		Name:  "Pico",
		Power: -12,
	},
	{
		Name:  "Femto",
		Power: -15,
	},
	{
		Name:  "Atto",
		Power: -18,
	},
	{
		Name:  "Zepto",
		Power: -21,
	},
	{
		Name:  "Yocto",
		Power: -24,
	},
}

var Units = []Unit{
	{
		Name:        "Mass",
		Offset:      -3,
		PrintString: "kg",
		Suffix:      "gram",
		TypeComment: "Mass represents a mass in kilograms",
		Dimensions: []Dimension{
			{
				Name:  MassName,
				Power: 1,
			},
		},
	},
	{
		Name:        "Length",
		PrintString: "m",
		Suffix:      "meter",
		TypeComment: "Length represents a length in meters",
		Dimensions: []Dimension{
			{
				Name:  LengthName,
				Power: 1,
			},
		},
	},
	{
		Name:        "Time",
		PrintString: "s",
		Suffix:      "second",
		TypeComment: "Time represents a time in seconds",
		ExtraConstant: []Constant{
			{
				Name:  "Hour",
				Value: "3600",
			},
			{
				Name:  "Minute",
				Value: "60",
			},
		},
		Dimensions: []Dimension{
			{
				Name:  TimeName,
				Power: 1,
			},
		},
		ErForm: "Timer",
	},
}

var gopath string
var unitPkgPath string

func init() {
	gopath = os.Getenv("GOPATH")
	if gopath == "" {
		log.Fatal("no gopath")
	}

	unitPkgPath = filepath.Join(gopath, "github.com", "gonum", "unit")
}

var licenceHeader = `// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.`

var imports = `import (
	"errors"
	"fmt"
	"math"
)`

// Generate generates a file for each of the units
func main() {
	for _, unit := range Units {
		generate(unit)
	}
}

func generate(unit Unit) {
	lowerName := strings.ToLower(unit.Name)
	filename := filepath.Join(unitPkgPath, lowerName+".go")
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Add the headers and type declarations
	strs := make([]string, 0, 100)
	strs = append(strs, licenceHeader)
	strs = append(strs, "")
	strs = append(strs, "package unit")
	strs = append(strs, "")
	strs = append(strs, imports)
	strs = append(strs, "// "+unit.TypeComment)
	strs = append(strs, "type "+unit.Name+" float64")
	strs = append(strs, "")

	// Add all of the constants
	strs = append(strs, "const (")

	for _, con := range unit.ExtraConstant {
		strs = append(strs, con.Name+" "+unit.Name+" = "+con.Value)
	}

	a := []rune(unit.Suffix)
	a[0] = unicode.ToUpper(a[0])
	upperSuffix := string(a)
	for _, prefix := range Prefixes {
		var str string
		if prefix.Name != "" {
			str = prefix.Name + unit.Suffix
		} else {
			str = upperSuffix
		}
		str += " " + unit.Name
		pow := prefix.Power + unit.Offset
		if pow == 0 {
			str += " = 1.0"
		} else {
			str += " = 1e" + strconv.Itoa(pow)
		}
		strs = append(strs, str)
	}
	strs = append(strs, ")")

	// Add the "Unit" method
	dim := make([]string, 0, 8)
	dim = append(dim, "Dimensions {")
	for _, dimension := range unit.Dimensions {
		dim = append(dim, dimension.Name+": "+strconv.Itoa(dimension.Power)+",")
	}
	dim = append(dim, "}")

	dimStr := strings.Join(dim, "\n")

	strs = append(strs, "// Unit converts the "+unit.Name+" to a *Unit")
	strs = append(strs, "func ("+lowerName+" "+unit.Name+") Unit() *Unit{")
	strs = append(strs, "return New(float64("+lowerName+"),"+dimStr+")")
	strs = append(strs, "}")

	// Add the Xxxer method
	erForm := unit.ErForm
	if erForm == "" {
		erForm = unit.Name + "er"
	}

	strs = append(strs, "// "+unit.Name+" allows "+unit.Name+" to implement a "+erForm+" interface")
	strs = append(strs, "func ("+lowerName+" "+unit.Name+")"+unit.Name+"() "+unit.Name+"{")
	strs = append(strs, "return "+lowerName)
	strs = append(strs, "}")

	// Add the "From" method
	strs = append(strs, "// From converts a Uniter to a "+unit.Name+". Returns an error if")
	strs = append(strs, "// there is a mismatch in dimension")
	strs = append(strs, "func ("+lowerName+" *"+unit.Name+") From(u Uniter) error {")
	strs = append(strs, "if !DimensionsMatch(u, "+upperSuffix+"){")
	strs = append(strs, "(*"+lowerName+") = "+unit.Name+"(math.NaN())")
	strs = append(strs, "return errors.New(\"Dimension mismatch\")")
	strs = append(strs, "}")
	strs = append(strs, "(*"+lowerName+") = "+unit.Name+"(u.Unit().Value())")
	strs = append(strs, "return nil")
	strs = append(strs, "}")

	// Add the "Format" method
	// case 'v'
	strs = append(strs, "")
	strs = append(strs, "func ("+lowerName+" "+unit.Name+") Format(fs fmt.State, c rune ) {")
	strs = append(strs, "switch c {")
	strs = append(strs, "case 'v':")
	strs = append(strs, "if fs.Flag('#') {")
	strs = append(strs, "fmt.Fprintf(fs, \"%T(%v)\", "+lowerName+", float64("+lowerName+"))")
	strs = append(strs, "return")
	strs = append(strs, "}")
	strs = append(strs, "fallthrough")

	// case others
	str := `	case 'e', 'E', 'f', 'F', 'g', 'G':
p, pOk := fs.Precision()
if !pOk {
p = -1
}
w, wOk := fs.Width()
if !wOk {
w = -1
}
fmt.Fprintf(fs, "%*.*"+string(c), w, p, float64(`
	str += lowerName + "))"

	strs = append(strs, str)
	strs = append(strs, "fmt.Fprint(fs, \" "+unit.PrintString+"\")")

	//case default
	strs = append(strs, "default:")
	strs = append(strs, "fmt.Fprintf(fs, \"%%!%c(%T=%g "+unit.PrintString+")\", c, "+lowerName+", float64("+lowerName+"))")
	strs = append(strs, "return")
	strs = append(strs, "}")
	strs = append(strs, "}")
	// Run go fmt
	s := strings.Join(strs, "\n")
	b, err := format.Source([]byte(s))
	if err != nil {
		f.WriteString(s) // This is here to debug bad format
		log.Fatal(err)
	}

	f.Write(b)
}
