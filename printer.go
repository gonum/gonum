// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package opt

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// Printer writes column-format output to the specified writer as the optimization
// progresses. By default, it writes to Stdout.
type Printer struct {
	Writer          io.Writer
	HeadingInterval int
	ValueInterval   time.Duration

	lastHeading int
	printGrad   bool
	lastValue   time.Time
}

func NewPrinter() *Printer {
	return &Printer{
		Writer:          os.Stdout,
		HeadingInterval: 30,
		ValueInterval:   500 * time.Millisecond,
	}
}

const nPrinterOut = 4

var (
	// TOOD: Add test that these match up
	printerHeadings = [nPrinterOut]string{"Iter", "FunEval", "Obj", "GradNorm"}
)

func (p *Printer) Init(f *FunctionStats) error {
	p.printGrad = f.IsFunGrad || f.IsGradient

	p.lastHeading = p.HeadingInterval + 1          // So the headings are printed the first time
	p.lastValue = time.Now().Add(-p.ValueInterval) // So the values are printed the first time
	return nil
}

func (p *Printer) Record(l Location, eval EvaluationType, iter IterationType, stats *Stats) error {
	// Only print on major iterations
	if iter != Major {
		return nil
	}

	var nPrint int
	if p.printGrad {
		nPrint = nPrinterOut
	} else {
		nPrint = nPrinterOut - 1
	}

	// Make the value strings
	var valueStrings [nPrinterOut]string
	valueStrings[0] = strconv.Itoa(stats.NumMajorIterations)
	valueStrings[1] = strconv.Itoa(stats.NumFunEvals + stats.NumFunGradEvals)
	valueStrings[2] = fmt.Sprintf("%g", l.F)
	if p.printGrad {
		valueStrings[3] = fmt.Sprintf("%g", stats.GradNorm)
	}

	var maxLengths [nPrinterOut]int

	for i := 0; i < nPrint; i++ {
		v := len(printerHeadings[i])
		v2 := len(valueStrings[i])
		if v > v2 {
			maxLengths[i] = v
			continue
		}
		maxLengths[i] = v2
	}

	// First, see if we want to print the headings
	if p.lastHeading >= p.HeadingInterval {
		// Yes we do
		p.lastHeading = 0

		headingString := constructPrinterString(printerHeadings, maxLengths, nPrint)
		// Add an extra newline to heading string
		headingString = "\n" + headingString

		_, err := p.Writer.Write([]byte(headingString))
		if err != nil {
			return err
		}

	}
	// See if we want to print the value
	if time.Since(p.lastValue) > p.ValueInterval {
		// Yes we do
		p.lastHeading++
		p.lastValue = time.Now()

		valueString := constructPrinterString(valueStrings, maxLengths, nPrint)
		_, err := p.Writer.Write([]byte(valueString))
		if err != nil {
			return err
		}
	}
	return nil
}

// pad string adds spaces onto the end of the string to make it the correct length
func padString(s string, l int) string {
	if len(s) > l {
		panic("string too long")
	}
	if len(s) == l {
		return s
	}
	nShort := l - len(s)
	return s + strings.Repeat(" ", nShort)
}

func constructPrinterString(values [nPrinterOut]string, maxLengths [nPrinterOut]int, nPrint int) string {
	var str string
	for i := 0; i < nPrint; i++ {
		s := values[i]
		s = padString(s, maxLengths[i])
		str += s
		if i != nPrint-1 {
			str += "\t"
		}
	}
	str += "\n"
	return str
}
