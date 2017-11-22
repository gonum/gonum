// Copyright 2017 The Gonum Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fit_test

import (
	"bufio"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"testing"

	"gonum.org/v1/gonum/fit"
	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/optimize"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func TestCurve1D(t *testing.T) {
	ExampleCurve1D_gaussian()
	ExampleCurve1D_exponential()
	ExampleCurve1D_poly()
}

func ExampleCurve1D_gaussian() {
	var (
		cst   = 3.0
		mean  = 30.0
		sigma = 20.0
		want  = []float64{cst, mean, sigma}
	)

	xdata, ydata, err := readXY("testdata/gauss-data.txt")

	gauss := func(x, cst, mu, sigma float64) float64 {
		v := (x - mu)
		return cst * math.Exp(-v*v/sigma)
	}

	res, err := fit.Curve1D(
		fit.Func1D{
			F: func(x float64, ps []float64) float64 {
				return gauss(x, ps[0], ps[1], ps[2])
			},
			X:  xdata,
			Y:  ydata,
			Ps: []float64{10, 10, 10},
		},
		nil, &optimize.NelderMead{},
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := res.Status.Err(); err != nil {
		log.Fatal(err)
	}
	if got := res.X; !floats.EqualApprox(got, want, 1e-3) {
		log.Fatalf("got= %v\nwant=%v\n", got, want)
	}

	{
		p, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}
		p.X.Label.Text = "Gauss"
		p.Y.Label.Text = "y-data"

		s, err := plotter.NewScatter(xyFrom(xdata, ydata))
		if err != nil {
			log.Fatal(err)
		}
		s.Color = color.RGBA{0, 0, 255, 255}
		s.Shape = draw.CrossGlyph{}
		p.Add(s)

		f := plotter.NewFunction(func(x float64) float64 {
			return gauss(x, res.X[0], res.X[1], res.X[2])
		})
		f.Color = color.RGBA{255, 0, 0, 255}
		f.Samples = 1000
		p.Add(f)

		p.Add(plotter.NewGrid())

		err = p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/gauss-plot.png")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ExampleCurve1D_exponential() {
	const (
		a   = 0.3
		b   = 0.1
		ndf = 2.0
	)

	xdata, ydata, err := readXY("testdata/exp-data.txt")

	exp := func(x, a, b float64) float64 {
		return math.Exp(a*x + b)
	}

	res, err := fit.Curve1D(
		fit.Func1D{
			F: func(x float64, ps []float64) float64 {
				return exp(x, ps[0], ps[1])
			},
			X: xdata,
			Y: ydata,
			N: 2,
		},
		nil, &optimize.NelderMead{},
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := res.Status.Err(); err != nil {
		log.Fatal(err)
	}
	if got, want := res.X, []float64{a, b}; !floats.EqualApprox(got, want, 0.1) {
		log.Fatalf("got= %v\nwant=%v\n", got, want)
	}

	{
		p, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}
		p.X.Label.Text = "exp(a*x+b)"
		p.Y.Label.Text = "y-data"
		p.Y.Min = 0
		p.Y.Max = 5
		p.X.Min = 0
		p.X.Max = 5

		s, err := plotter.NewScatter(xyFrom(xdata, ydata))
		if err != nil {
			log.Fatal(err)
		}
		s.Color = color.RGBA{0, 0, 255, 255}
		s.Shape = draw.CrossGlyph{}
		p.Add(s)

		f := plotter.NewFunction(func(x float64) float64 {
			return exp(x, res.X[0], res.X[1])
		})
		f.Color = color.RGBA{255, 0, 0, 255}
		f.Samples = 1000
		p.Add(f)

		p.Add(plotter.NewGrid())

		err = p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/exp-plot.png")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ExampleCurve1D_poly() {
	var (
		a    = 1.0
		b    = 2.0
		ps   = []float64{a, b}
		want = []float64{1.38592513, 1.98485122} // from scipy.curve_fit
	)

	poly := func(x float64, ps []float64) float64 {
		return ps[0] + ps[1]*x*x
	}

	xdata, ydata := genXY(100, poly, ps, -10, 10)

	res, err := fit.Curve1D(
		fit.Func1D{
			F:  poly,
			X:  xdata,
			Y:  ydata,
			Ps: []float64{1, 1},
		},
		nil, &optimize.NelderMead{},
	)
	if err != nil {
		log.Fatal(err)
	}

	if err := res.Status.Err(); err != nil {
		log.Fatal(err)
	}

	if got := res.X; !floats.EqualApprox(got, want, 1e-6) {
		log.Fatalf("got= %v\nwant=%v\n", got, want)
	}

	{
		p, err := plot.New()
		if err != nil {
			log.Fatal(err)
		}
		p.X.Label.Text = "f(x) = a + b*x*x"
		p.Y.Label.Text = "y-data"
		p.X.Min = -10
		p.X.Max = +10
		p.Y.Min = 0
		p.Y.Max = 220

		s, err := plotter.NewScatter(xyFrom(xdata, ydata))
		if err != nil {
			log.Fatal(err)
		}
		s.Color = color.RGBA{0, 0, 255, 255}
		s.Shape = draw.CrossGlyph{}
		p.Add(s)

		f := plotter.NewFunction(func(x float64) float64 {
			return poly(x, res.X)
		})
		f.Color = color.RGBA{255, 0, 0, 255}
		f.Samples = 1000
		p.Add(f)

		p.Add(plotter.NewGrid())

		err = p.Save(10*vg.Centimeter, 10*vg.Centimeter, "testdata/poly-plot.png")
		if err != nil {
			log.Fatal(err)
		}
	}
}

func genXY(n int, f func(x float64, ps []float64) float64, ps []float64, xmin, xmax float64) ([]float64, []float64) {
	xdata := make([]float64, n)
	ydata := make([]float64, n)
	rnd := rand.New(rand.NewSource(1234))
	xstep := (xmax - xmin) / float64(n)
	p := make([]float64, len(ps))
	for i := 0; i < n; i++ {
		x := xmin + xstep*float64(i)
		for j := range p {
			v := rnd.NormFloat64()
			p[j] = ps[j] + v*0.2
		}
		xdata[i] = x
		ydata[i] = f(x, p)
	}
	return xdata, ydata
}

func readXY(fname string) (xs, ys []float64, err error) {
	f, err := os.Open(fname)
	if err != nil {
		return xs, ys, err
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		toks := strings.Split(line, " ")
		x, err := strconv.ParseFloat(toks[0], 64)
		if err != nil {
			return xs, ys, err
		}
		xs = append(xs, x)

		y, err := strconv.ParseFloat(toks[1], 64)
		if err != nil {
			return xs, ys, err
		}
		ys = append(ys, y)
	}

	return
}

func readXYerr(fname string) (xs, ys, yerrs []float64, err error) {
	f, err := os.Open(fname)
	if err != nil {
		return xs, ys, yerrs, err
	}
	defer f.Close()

	scan := bufio.NewScanner(f)
	for scan.Scan() {
		line := scan.Text()
		toks := strings.Split(line, " ")
		x, err := strconv.ParseFloat(toks[0], 64)
		if err != nil {
			return xs, ys, yerrs, err
		}
		xs = append(xs, x)

		y, err := strconv.ParseFloat(toks[1], 64)
		if err != nil {
			return xs, ys, yerrs, err
		}
		ys = append(ys, y)

		yerr, err := strconv.ParseFloat(toks[2], 64)
		if err != nil {
			return xs, ys, yerrs, err
		}
		yerrs = append(yerrs, yerr)
	}

	return
}

func xyFrom(xs, ys []float64) plotter.XYer {
	if len(xs) != len(ys) {
		panic("lengths do not match")
	}

	var xys = make(plotter.XYs, len(xs))
	for i := range xs {
		xys[i].X = xs[i]
		xys[i].Y = ys[i]
	}
	return xys
}
