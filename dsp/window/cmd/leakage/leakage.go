// Copyright ©2021 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// The leakage program provides summary characteristics and a plot
// of spectral response for window functions or csv input. It is intended
// to be used to verify window behaviour against foreign implementations.
// For example, the behavior of a NumPy window can be captured using this
// python code:
//
//	import matplotlib.pyplot as plt
//	import numpy as np
//	from numpy.fft import fft
//
//	window = np.blackman(20)
//	print("# beta = %f" % np.mean(window))
//
//	plt.figure()
//
//	A = fft(window, 1000)
//	mag = np.abs(A)
//	with np.errstate(divide='ignore', invalid='ignore'):
//	    mag = 20 * np.log10(mag)
//	mag -= max(mag)
//	freq = np.linspace(0, len(window)/2, len(A)/2)
//
//	for m in mag[:len(mag)/2]:
//		print(m)
//
//	plt.plot(freq, mag[:len(mag)/2])
//	plt.title("Spectral leakage")
//	plt.ylabel("Amplitude (dB)")
//	plt.xlabel("DFT bin")
//
//	plt.show()
//
// and then be exported to leakage and compared with the Gonum
// implementation.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"image/color"
	"io"
	"log"
	"math"
	"math/cmplx"
	"os"
	"strconv"
	"strings"

	"gonum.org/v1/gonum/dsp/fourier"
	"gonum.org/v1/gonum/dsp/window"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/stat"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

var windows = map[string]*builtin{
	"rectangular": {
		name: func(_ float64) string { return "Rectangular" },
		fn:   func(_ float64) func([]float64) []float64 { return window.Rectangular },
		ok:   func(_ float64) bool { return true },
	},
	"sine": {
		name: func(_ float64) string { return "Sine" },
		fn:   func(_ float64) func([]float64) []float64 { return window.Sine },
		ok:   func(_ float64) bool { return true },
	},
	"lanczos": {
		name: func(_ float64) string { return "Lanczos" },
		fn:   func(_ float64) func([]float64) []float64 { return window.Lanczos },
		ok:   func(_ float64) bool { return true },
	},
	"triangular": {
		name: func(_ float64) string { return "Triangular" },
		fn:   func(_ float64) func([]float64) []float64 { return window.Triangular },
		ok:   func(_ float64) bool { return true },
	},
	"hann": {
		name: func(_ float64) string { return "Hann" },
		fn:   func(_ float64) func([]float64) []float64 { return window.Hann },
		ok:   func(_ float64) bool { return true },
	},
	"bartletthann": {
		name: func(_ float64) string { return "BartlettHann" },
		fn:   func(_ float64) func([]float64) []float64 { return window.BartlettHann },
		ok:   func(_ float64) bool { return true },
	},
	"hamming": {
		name: func(_ float64) string { return "Hamming" },
		fn:   func(_ float64) func([]float64) []float64 { return window.Hamming },
		ok:   func(_ float64) bool { return true },
	},
	"blackman": {
		name: func(_ float64) string { return "Blackman" },
		fn:   func(_ float64) func([]float64) []float64 { return window.Blackman },
		ok:   func(_ float64) bool { return true },
	},
	"blackmanharris": {
		name: func(_ float64) string { return "BlackmanHarris" },
		fn:   func(_ float64) func([]float64) []float64 { return window.BlackmanHarris },
		ok:   func(_ float64) bool { return true },
	},
	"nuttall": {
		name: func(_ float64) string { return "Nuttall" },
		fn:   func(_ float64) func([]float64) []float64 { return window.Nuttall },
		ok:   func(_ float64) bool { return true },
	},
	"blackmannuttall": {
		name: func(_ float64) string { return "BlackmanNuttall" },
		fn:   func(_ float64) func([]float64) []float64 { return window.BlackmanNuttall },
		ok:   func(_ float64) bool { return true },
	},
	"flattop": {
		name: func(_ float64) string { return "FlatTop" },
		fn:   func(_ float64) func([]float64) []float64 { return window.FlatTop },
		ok:   func(_ float64) bool { return true },
	},
	"gaussian": {
		name: func(p float64) string { return fmt.Sprintf("Gaussian σ=%v", p) },
		fn:   func(p float64) func([]float64) []float64 { return window.Gaussian{Sigma: p}.Transform },
		ok:   func(p float64) bool { return !math.IsNaN(p) },
	},
	"tukey": {
		name: func(p float64) string { return fmt.Sprintf("Tukey α=%v", p) },
		fn:   func(p float64) func([]float64) []float64 { return window.Tukey{Alpha: p}.Transform },
		ok:   func(p float64) bool { return !math.IsNaN(p) },
	},
}

type builtin struct {
	name func(float64) string
	fn   func(float64) func([]float64) []float64
	ok   func(float64) bool
}

func main() {
	name := flag.String("window", "", "specify a built-in window name (required if csv not given)")
	param := flag.Float64("param", math.NaN(), "specify parameter for parametric windows")
	symm := flag.Bool("symm", true, "specify whether the window is symmetric")
	n := flag.Int("n", 1024, "specify window length")
	m := flag.Int("m", 1<<16, "specify sample length (must be greater than n)")
	csv := flag.String("csv", "", "specify an input file of dB transformed frequency amplitudes (required if window not given)")
	out := flag.String("o", "", "specify output file for plot (required, formats eps, jpg, jpeg, pdf, png, svg, tex or tif)")
	width := flag.Float64("width", 16, "specify plot width (cm)")
	height := flag.Float64("height", 8, "specify plot height (cm)")
	comment := flag.Bool("comment", false, "output a comment line for the window (only for window function)")
	eps := flag.Float64("eps", 1e-3, "warn if parameters have not converged to this relative tolerance at -n (zero disables the check)")
	flag.Parse()

	win := windows[strings.ToLower(*name)]
	if win == nil && *csv == "" {
		fmt.Fprintln(os.Stderr, "missing function name or csv input")
		flag.Usage()
		os.Exit(2)
	}
	if *csv == "" && !win.ok(*param) {
		fmt.Fprintln(os.Stderr, "missing parameter")
		flag.Usage()
		os.Exit(2)
	}
	if *out == "" {
		fmt.Fprintln(os.Stderr, "missing output filename")
		flag.Usage()
		os.Exit(2)
	}

	p := plot.New()
	p.X.Label.Text = "DFT bin"
	p.Y.Label.Text = "Amplitude [dB]"
	p.Add(plotter.NewGrid())

	var (
		c        *characteristics
		spectrum plotter.XYs
		min      float64
		err      error
	)
	if win != nil {
		symmetry := "(symmetric)"
		if !*symm {
			symmetry = "(periodic)"
		}
		p.Title.Text = fmt.Sprintf("Spectral Leakage for %s %s", win.name(*param), symmetry)
		c, spectrum, min, err = funcCharacteristics(win.fn(*param), *n, *m, *symm)
		if err != nil {
			log.Fatal(err)
		}
		if *eps > 0 {
			warnUnconverged(os.Stderr, win.fn(*param), *n, *m, *symm, c, *eps)
		}
		if *comment {
			fmt.Printf("// Spectral leakage parameters: ΔF_0 = %2f, ΔF_0.5 = %.2f, ENBW = %.3f, K = %.2f, ɣ_max = %2f, β = %2f.\n",
				c.deltaF0, c.deltaFhalf, c.enbw, c.k(), c.gammaMax, c.beta)
		}
	} else {
		f, err := os.Open(*csv)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		p.Title.Text = fmt.Sprintf("Spectral Leakage for %s", *csv)
		c, spectrum, min, err = csvCharacteristics(f, *n, *m)
		if err != nil {
			log.Fatal(err)
		}
	}

	spectrumLines, err := plotter.NewLine(spectrum)
	if err != nil {
		log.Fatalf("spectrum: %v", err)
	}

	gammaLine, err := plotter.NewLine(plotter.XYs{
		{X: 0, Y: c.gammaMax}, {X: float64(*n) / 2, Y: c.gammaMax},
	})
	if err != nil {
		log.Fatalf("ɣ_max: %v", err)
	}
	gammaLine.Color = color.RGBA{R: 0x40, G: 0x80, B: 0xff, A: 0xff}

	deltaF0Line, err := plotter.NewLine(plotter.XYs{
		{X: c.deltaF0 / 2, Y: 0}, {X: c.deltaF0 / 2, Y: min},
	})
	if err != nil {
		log.Fatalf("ΔF_0: %v", err)
	}
	deltaF0Line.Color = color.RGBA{R: 0xff, A: 0xff}

	deltaFhalfLine, err := plotter.NewLine(plotter.XYs{
		{X: c.deltaFhalf / 2, Y: 0}, {X: c.deltaFhalf / 2, Y: min},
	})
	if err != nil {
		log.Fatalf("ΔF_0.5: %v", err)
	}
	deltaFhalfLine.Color = color.RGBA{G: 0xff, A: 0xff}

	var blank plotter.Line

	p.Add(
		gammaLine,
		deltaF0Line,
		deltaFhalfLine,
		spectrumLines,
	)
	p.Legend.Add(fmt.Sprintf("ΔF_0 = %.3v", c.deltaF0), deltaF0Line)
	p.Legend.Add(fmt.Sprintf("ΔF_0.5 = %.3v", c.deltaFhalf), deltaFhalfLine)
	p.Legend.Add(fmt.Sprintf("K = %.3v", c.k()), &blank)
	p.Legend.Add(fmt.Sprintf("ɣ_max = %.3v", c.gammaMax), gammaLine)
	p.Legend.Add(fmt.Sprintf("β = %.3v", c.beta), &blank)
	p.Legend.Top = true
	p.Legend.XOffs = -10
	p.Legend.YOffs = -10

	err = p.Save(vg.Length(*width)*vg.Centimeter, vg.Length(*height)*vg.Centimeter, *out)
	if err != nil {
		log.Fatal(err)
	}
}

// characteristics hold DFT window characteristic statistics.
// See http://www.dsplib.ru/content/win/win.html for details.
type characteristics struct {
	deltaF0    float64
	deltaFhalf float64
	gammaMax   float64
	beta       float64

	// enbw is the equivalent noise bandwidth, which is a distinct
	// parameter from deltaFhalf. It is NaN when the window itself is not
	// available, as when reading a spectrum from csv.
	enbw float64
}

// k returns the K window parameter which is the ratio of the window's
// ΔF_0 to the ΔF_0 of the rectangular window.
func (c *characteristics) k() float64 {
	return c.deltaF0 / 2
}

func funcCharacteristics(fn func([]float64) []float64, n, m int, symm bool) (c *characteristics, xy plotter.XYs, min float64, err error) {
	if m < n {
		return nil, nil, 0, fmt.Errorf("window: sequence too short for window: %d < %d", m, n)
	}

	var w []float64
	t := make([]float64, m)
	if symm {
		w = window.NewValues(fn, n)
	} else {
		w = window.NewValues(fn, n+1)[:n]
	}

	copy(t, w)

	var max float64
	xy = make(plotter.XYs, m/2+1)
	fft := fourier.NewFFT(len(t))
	for i, c := range fft.Coefficients(nil, t) {
		a := db(cmplx.Abs(c))
		t[i] = a
		if !math.IsInf(a, -1) && a < min {
			min = a
		}
		if i == 0 {
			max = a
		}
	}
	for i, a := range t[:m/2+1] {
		if math.IsInf(a, -1) {
			a = min
		}
		xy[i] = plotter.XY{X: float64(i) * float64(n) / float64(m), Y: a - max}
	}

	c = &characteristics{beta: db(stat.Mean(w, nil)), enbw: enbw(w)}
	c.deltaF0, c.deltaFhalf, c.gammaMax = parameters(t, n, m)

	return c, xy, min - max, nil
}

func csvCharacteristics(r io.Reader, n, m int) (c *characteristics, xy plotter.XYs, min float64, err error) {
	if m < n {
		return nil, nil, 0, fmt.Errorf("window: sequence too short for window: %d < %d", m, n)
	}
	sc := bufio.NewScanner(r)
	max := math.Inf(-1)
	var t []float64
	for sc.Scan() {
		text := sc.Text()
		if strings.HasPrefix(text, "#") {
			continue
		}
		v, err := strconv.ParseFloat(text, 64)
		if err != nil {
			log.Fatal(err)
		}
		if v > max {
			max = v
		}
		t = append(t, v)
	}

	xy = make(plotter.XYs, len(t))
	for i, a := range t {
		if math.IsInf(a, -1) {
			a = min
		} else if a < min {
			min = a
		}
		if i == 0 {
			max = a
		}
		xy[i] = plotter.XY{X: float64(i) * float64(n) / float64(m), Y: a - max}
	}
	err = sc.Err()
	if err != nil {
		return nil, nil, 0, err
	}

	c = &characteristics{beta: math.NaN()}
	c.deltaF0, c.deltaFhalf, c.gammaMax = parameters(t, n, m)

	return c, xy, min - max, nil
}

func parameters(spectrum []float64, n, m int) (deltaF0, deltaFhalf, gammaMax float64) {
	// Locate the peak rather than assuming it is at DC. A flat top window
	// is designed to have a flat passband and its maximum lies a little
	// away from zero frequency, so taking spectrum[0] as the reference
	// makes its main lobe width and maximum sidelobe meaningless.
	half := len(spectrum) / 2
	pk := 0
	for i := 1; i < half; i++ {
		if spectrum[i] > spectrum[pk] {
			pk = i
		}
	}
	max := spectrum[pk]

	// Frequency of bin i in DFT bins of the window.
	freq := func(i int) float64 { return float64(i) * float64(n) / float64(m) }

	// First null: the first point past the peak where the response stops
	// falling, refined to the local minimum.
	null := half - 1
	for i := pk + 1; i < half; i++ {
		if spectrum[i] > spectrum[i-1] {
			null = i - 1
			break
		}
	}
	deltaF0 = 2 * (freq(null) - freq(pk))

	// Half power width, linearly interpolated between the samples that
	// straddle the -3 dB crossing. Taking the last sample above the
	// threshold instead biases the result low by up to one bin.
	thresh := max - 3
	for i := pk + 1; i <= null; i++ {
		if spectrum[i] <= thresh {
			x0, x1 := freq(i-1), freq(i)
			y0, y1 := spectrum[i-1], spectrum[i]
			cross := x0
			if y0 != y1 {
				cross = x0 + (x1-x0)*(y0-thresh)/(y0-y1)
			}
			deltaFhalf = 2 * (cross - freq(pk))
			break
		}
	}

	// Highest sidelobe, measured beyond the main lobe and relative to the
	// true peak.
	var peaks []float64
	for i := null + 1; i < half-1; i++ {
		if spectrum[i-1] <= spectrum[i] && spectrum[i] > spectrum[i+1] {
			peaks = append(peaks, spectrum[i])
		}
	}
	if len(peaks) == 0 {
		gammaMax = math.NaN()
	} else {
		gammaMax = floats.Max(peaks) - max
	}

	return deltaF0, deltaFhalf, gammaMax
}

// enbw returns the equivalent noise bandwidth of the window w, in DFT bins.
//
// ENBW is the width of the ideal rectangular filter passing the same noise
// power as the window, and is a distinct parameter from the half power width
// ΔF_0.5: for a Hann window ENBW is 1.5 bins while ΔF_0.5 is about 1.44.
func enbw(w []float64) float64 {
	var sum, sumSq float64
	for _, v := range w {
		sum += v
		sumSq += v * v
	}
	if sum == 0 {
		return math.NaN()
	}
	return float64(len(w)) * sumSq / (sum * sum)
}

// warnUnconverged recomputes the characteristics at twice the window length
// and reports any parameter that still moves by more than eps in relative
// terms. The leakage parameters are asymptotic in the window length, so a
// short window reports values that have not settled: at n = 20 a Blackman
// Nuttall window appears to have a maximum sidelobe of -85 dB where the
// converged figure is -98 dB.
func warnUnconverged(w io.Writer, fn func([]float64) []float64, n, m int, symm bool, c *characteristics, eps float64) {
	ref, _, _, err := funcCharacteristics(fn, 2*n, 2*m, symm)
	if err != nil {
		fmt.Fprintf(w, "warning: cannot check convergence: %v\n", err)
		return
	}
	for _, p := range []struct {
		name     string
		got, ref float64
	}{
		{"ΔF_0", c.deltaF0, ref.deltaF0},
		{"ΔF_0.5", c.deltaFhalf, ref.deltaFhalf},
		{"ENBW", c.enbw, ref.enbw},
		{"K", c.k(), ref.k()},
		{"ɣ_max", c.gammaMax, ref.gammaMax},
		{"β", c.beta, ref.beta},
	} {
		if math.IsNaN(p.got) || math.IsNaN(p.ref) || p.ref == 0 {
			continue
		}
		if d := math.Abs(p.got-p.ref) / math.Abs(p.ref); d > eps {
			fmt.Fprintf(w, "warning: %s has not converged at n=%d: %.6g here against %.6g at n=%d (%.2g relative)\n",
				p.name, n, p.got, p.ref, 2*n, d)
		}
	}
}

func db(m float64) float64 {
	return 20 * math.Log10(m)
}
