// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package window_test

import (
	"fmt"
	"log"
	"math"
	"math/cmplx"
	"os"

	"gonum.org/v1/gonum/dsp/fourier"
	"gonum.org/v1/gonum/dsp/window"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgimg"
)

func Example() {
	// The input sequence is a 2.5 period of the Sin function.
	src := make([]float64, 20)
	k := 5 * math.Pi / float64(len(src)-1)
	for i := range src {
		src[i] = math.Sin(k * float64(i))
	}

	// Initialize an FFT and perform the analysis.
	fft := fourier.NewFFT(len(src))
	coeff := fft.Coefficients(nil, src)

	// The result shows that width of the main lobe with center
	// between frequencies 0.1 and 0.15 is small, but that the
	// height of the side lobes is large.
	fmt.Println("Rectangular window (or no window):")
	for i, c := range coeff {
		fmt.Printf("freq=%.4f\tcycles/period, magnitude=%.4f,\tphase=%.4f\n",
			fft.Freq(i), cmplx.Abs(c), cmplx.Phase(c))
	}

	// Initialize an FFT and perform the analysis on a sequence
	// transformed by the Hamming window function.
	fft = fourier.NewFFT(len(src))
	coeff = fft.Coefficients(nil, window.Hamming(src))

	// The result shows that width of the main lobe is wider,
	// but height of the side lobes is lower.
	fmt.Println("Hamming window:")
	// The magnitude of all bins has been decreased by β.
	// Generally in an analysis amplification may be omitted, but to
	// make a comparable data, the result should be amplified by -β
	// of the window function — +5.37 dB for the Hamming window.
	//  -β = 20 log_10(amplifier).
	amplifier := math.Pow(10, 5.37/20.0)
	for i, c := range coeff {
		fmt.Printf("freq=%.4f\tcycles/period, magnitude=%.4f,\tphase=%.4f\n",
			fft.Freq(i), amplifier*cmplx.Abs(c), cmplx.Phase(c))
	}
	// Output:
	//
	// Rectangular window (or no window):
	// freq=0.0000	cycles/period, magnitude=2.2798,	phase=0.0000
	// freq=0.0500	cycles/period, magnitude=2.6542,	phase=0.1571
	// freq=0.1000	cycles/period, magnitude=5.3115,	phase=0.3142
	// freq=0.1500	cycles/period, magnitude=7.3247,	phase=-2.6704
	// freq=0.2000	cycles/period, magnitude=1.6163,	phase=-2.5133
	// freq=0.2500	cycles/period, magnitude=0.7681,	phase=-2.3562
	// freq=0.3000	cycles/period, magnitude=0.4385,	phase=-2.1991
	// freq=0.3500	cycles/period, magnitude=0.2640,	phase=-2.0420
	// freq=0.4000	cycles/period, magnitude=0.1530,	phase=-1.8850
	// freq=0.4500	cycles/period, magnitude=0.0707,	phase=-1.7279
	// freq=0.5000	cycles/period, magnitude=0.0000,	phase=0.0000
	// Hamming window:
	// freq=0.0000	cycles/period, magnitude=0.0506,	phase=0.0000
	// freq=0.0500	cycles/period, magnitude=0.5386,	phase=-2.9845
	// freq=0.1000	cycles/period, magnitude=7.3350,	phase=0.3142
	// freq=0.1500	cycles/period, magnitude=8.9523,	phase=-2.6704
	// freq=0.2000	cycles/period, magnitude=1.7979,	phase=0.6283
	// freq=0.2500	cycles/period, magnitude=0.0957,	phase=0.7854
	// freq=0.3000	cycles/period, magnitude=0.0050,	phase=-2.1991
	// freq=0.3500	cycles/period, magnitude=0.0158,	phase=-2.0420
	// freq=0.4000	cycles/period, magnitude=0.0125,	phase=-1.8850
	// freq=0.4500	cycles/period, magnitude=0.0065,	phase=-1.7279
	// freq=0.5000	cycles/period, magnitude=0.0000,	phase=0.0000
}

func ExampleHamming() {
	src := []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

	// Window functions change data in place. So, if input data
	// needs to stay unchanged, it must be copied.
	srcCpy := append([]float64(nil), src...)
	// Apply window function to srcCpy.
	dst := window.Hamming(srcCpy)

	// src is unchanged.
	fmt.Printf("src:    %f\n", src)
	// srcCpy is altered.
	fmt.Printf("srcCpy: %f\n", srcCpy)
	// dst mirrors the srcCpy slice.
	fmt.Printf("dst:    %f\n", dst)

	// Output:
	//
	// src:    [1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000 1.000000]
	// srcCpy: [0.092577 0.136714 0.220669 0.336222 0.472063 0.614894 0.750735 0.866288 0.950242 0.994379 0.994379 0.950242 0.866288 0.750735 0.614894 0.472063 0.336222 0.220669 0.136714 0.092577]
	// dst:    [0.092577 0.136714 0.220669 0.336222 0.472063 0.614894 0.750735 0.866288 0.950242 0.994379 0.994379 0.950242 0.866288 0.750735 0.614894 0.472063 0.336222 0.220669 0.136714 0.092577]
}

func ExampleValues() {
	src := []float64{1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

	// Create a Sine Window lookup table.
	sine := window.NewValues(window.Sine, len(src))

	// Apply the transformation to the src.
	fmt.Printf("dst: %f\n", sine.Transform(src))

	// Output:
	//
	// dst: [0.078459 0.233445 0.382683 0.522499 0.649448 0.760406 0.852640 0.923880 0.972370 0.996917 0.996917 0.972370 0.923880 0.852640 0.760406 0.649448 0.522499 0.382683 0.233445 0.078459]
}

func Example_plotHamming() {
	const (
		N  = 51
		sz = 7.5 * vg.Centimeter
	)

	ham := window.NewValues(window.Hamming, N)

	fft := fourier.NewFFT(2048)
	tmp := make([]float64, 2048)
	copy(tmp, ham)

	a := fft.Coefficients(nil, tmp)
	mag := -math.MaxFloat64
	for i, v := range a {
		a[i] = 0.5 * v / N
		mag = math.Max(mag, cmplx.Abs(a[i]))
	}

	freq := make([]float64, len(a))
	for i := range freq {
		freq[i] = fft.Freq(i)
	}
	resp := make([]float64, len(freq))
	for i, v := range a {
		resp[i] = 20 * math.Log10(cmplx.Abs(v)/mag)
	}

	plots := make([][]*plot.Plot, 2)

	{
		p, err := plot.New()
		if err != nil {
			log.Fatalf("could not create plot: %+v", err)
		}
		p.Title.Text = "Hamming window"
		p.X.Label.Text = "Sample"
		p.Y.Label.Text = "Amplitude"

		data := make(plotter.XYs, len(ham))
		for i := range ham {
			data[i].X = float64(i)
			data[i].Y = ham[i]
		}

		pts, err := plotter.NewLine(data)
		if err != nil {
			log.Fatalf("could not create line: %+v", err)
		}
		pts.LineStyle.Color = plotutil.SoftColors[2]
		pts.LineStyle.Width = 1.5

		p.Add(plotter.NewGrid(), pts)
		plots[0] = []*plot.Plot{p}
	}
	{
		p, err := plot.New()
		if err != nil {
			log.Fatalf("could not create plot: %+v", err)
		}
		p.Title.Text = "Frequency response of the Hamming window"
		p.X.Label.Text = "Frequency [cycles/sample]"
		p.Y.Label.Text = "Magnitude [dB]"

		data := make(plotter.XYs, len(freq))
		for i := range freq {
			data[i].X = freq[i]
			data[i].Y = resp[i]
		}
		pts, err := plotter.NewLine(data)
		if err != nil {
			log.Fatalf("could not create line: %+v", err)
		}
		pts.LineStyle.Color = plotutil.SoftColors[2]
		pts.LineStyle.Width = 1.5

		p.Add(plotter.NewGrid(), pts)
		plots[1] = []*plot.Plot{p}
	}

	var (
		img = vgimg.New(math.Phi*sz, sz)
		dc  = draw.New(img)
		cs  = plot.Align(plots, draw.Tiles{Rows: 2, Cols: 1}, dc)
	)
	for i := 0; i < 2; i++ {
		plots[i][0].Draw(cs[i][0])
	}

	w, err := os.Create("testdata/hamming.png")
	if err != nil {
		log.Fatalf("could not create output file: %+v", err)
	}
	defer w.Close()

	png := vgimg.PngCanvas{Canvas: img}
	if _, err := png.WriteTo(w); err != nil {
		log.Fatalf("could not save plot: %+v", err)
	}

	// Output:
}
