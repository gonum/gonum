// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build fortran

package amos

import (
	"flag"
	"runtime"
	"testing"

	"golang.org/x/exp/rand"

	"gonum.org/v1/gonum/mathext/internal/amos/amoslib"
)

// BUG(kortschak): Some tests here comparing the direct Go translation
// of the Fortran code fail. Do not delete these tests or this file until
// https://github.com/gonum/gonum/issues/1322 has been satisfactorily
// resolved.
var runFailing = flag.Bool("failing", false, "run known failing cases")

func TestAiryFortran(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < nInputs; i++ {
		in := randInput(rnd)
		zairytestFort(t, in.x, in.kode, in.id)
	}
}

func TestZacaiFortran(t *testing.T) {
	if !*runFailing {
		t.Skip("fails")
	}

	switch runtime.GOARCH {
	case "arm64":
		t.Skipf("skipping on GOARCH=%s", runtime.GOARCH)
	}
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < nInputs; i++ {
		in := randInput(rnd)
		zacaitestFort(t, in.x, in.is, in.tol, in.n, in.yr, in.yi, in.kode)
	}
}

func TestZbknuFortran(t *testing.T) {
	if !*runFailing {
		t.Skip("fails")
	}

	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < nInputs; i++ {
		in := randInput(rnd)
		zbknutestFort(t, in.x, in.is, in.tol, in.n, in.yr, in.yi, in.kode)
	}
}

func TestZasyiFortran(t *testing.T) {
	if !*runFailing {
		t.Skip("fails")
	}

	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < nInputs; i++ {
		in := randInput(rnd)
		zasyitestFort(t, in.x, in.is, in.tol, in.n, in.yr, in.yi, in.kode)
	}
}

func TestZseriFortran(t *testing.T) {
	switch runtime.GOARCH {
	case "arm64":
		t.Skipf("skipping on GOARCH=%s", runtime.GOARCH)
	}
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < nInputs; i++ {
		in := randInput(rnd)
		zseritestFort(t, in.x, in.is, in.tol, in.n, in.yr, in.yi, in.kode)
	}
}

func TestZmlriFortran(t *testing.T) {
	if !*runFailing {
		t.Skip("fails")
	}

	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < nInputs; i++ {
		in := randInput(rnd)
		zmlritestFort(t, in.x, in.is, in.tol, in.n, in.yr, in.yi, in.kode)
	}
}

func TestZksclFortran(t *testing.T) {
	if !*runFailing {
		t.Skip("fails")
	}

	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < nInputs; i++ {
		in := randInput(rnd)
		zkscltestFort(t, in.x, in.is, in.tol, in.n, in.yr, in.yi)
	}
}

func TestZuchkFortran(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < nInputs; i++ {
		in := randInput(rnd)
		zuchktestFort(t, in.x, in.is, in.tol)
	}
}

func TestZs1s2Fortran(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < nInputs; i++ {
		in := randInput(rnd)
		zs1s2testFort(t, in.x, in.is)
	}
}

func zs1s2testFort(t *testing.T, x []float64, is []int) {
	const tol = 1e-11

	type data struct {
		ZRR, ZRI, S1R, S1I, S2R, S2I float64
		NZ                           int
		ASCLE, ALIM                  float64
		IUF                          int
	}

	input := data{
		x[0], x[1], x[2], x[3], x[4], x[5],
		is[0],
		x[6], x[7],
		is[1],
	}

	impl := func(input data) data {
		zrr, zri, s1r, s1i, s2r, s2i, nz, ascle, alim, iuf :=
			amoslib.Zs1s2Fort(input.ZRR, input.ZRI, input.S1R, input.S1I, input.S2R, input.S2I, input.NZ, input.ASCLE, input.ALIM, input.IUF)
		return data{zrr, zri, s1r, s1i, s2r, s2i, nz, ascle, alim, iuf}
	}

	comp := func(input data) data {
		zrr, zri, s1r, s1i, s2r, s2i, nz, ascle, alim, iuf :=
			zs1s2Orig(input.ZRR, input.ZRI, input.S1R, input.S1I, input.S2R, input.S2I, input.NZ, input.ASCLE, input.ALIM, input.IUF)
		return data{zrr, zri, s1r, s1i, s2r, s2i, nz, ascle, alim, iuf}
	}

	oi := impl(input)
	oc := comp(input)

	sameF64Approx(t, "zs1s2 zrr", oc.ZRR, oi.ZRR, tol)
	sameF64Approx(t, "zs1s2 zri", oc.ZRI, oi.ZRI, tol)
	sameF64Approx(t, "zs1s2 s1r", oc.S1R, oi.S1R, tol)
	sameF64Approx(t, "zs1s2 s1i", oc.S1I, oi.S1I, tol)
	sameF64Approx(t, "zs1s2 s2r", oc.S2R, oi.S2R, tol)
	sameF64Approx(t, "zs1s2 s2i", oc.S2I, oi.S2I, tol)
	sameF64Approx(t, "zs1s2 ascle", oc.ASCLE, oi.ASCLE, tol)
	sameF64Approx(t, "zs1s2 alim", oc.ALIM, oi.ALIM, tol)
	sameInt(t, "iuf", oc.IUF, oi.IUF)
	sameInt(t, "nz", oc.NZ, oi.NZ)
}

func zuchktestFort(t *testing.T, x []float64, is []int, tol float64) {
	t.Helper()

	YR := x[0]
	YI := x[1]
	NZ := is[0]
	ASCLE := x[2]
	TOL := tol

	YRfort, YIfort, NZfort, ASCLEfort, TOLfort := zuchkOrig(YR, YI, NZ, ASCLE, TOL)
	YRamoslib, YIamoslib, NZamoslib, ASCLEamoslib, TOLamoslib := amoslib.ZuchkFort(YR, YI, NZ, ASCLE, TOL)

	sameF64(t, "zuchk yr", YRfort, YRamoslib)
	sameF64(t, "zuchk yi", YIfort, YIamoslib)
	sameInt(t, "zuchk nz", NZfort, NZamoslib)
	sameF64(t, "zuchk ascle", ASCLEfort, ASCLEamoslib)
	sameF64(t, "zuchk tol", TOLfort, TOLamoslib)
}

func zkscltestFort(t *testing.T, x []float64, is []int, tol float64, n int, yr, yi []float64) {
	t.Helper()

	ZRR := x[0]
	ZRI := x[1]
	FNU := x[2]
	NZ := is[1]
	ELIM := x[3]
	ASCLE := x[4]
	RZR := x[6]
	RZI := x[7]

	yrfort := make([]float64, len(yr))
	copy(yrfort, yr)
	yifort := make([]float64, len(yi))
	copy(yifort, yi)
	ZRRfort, ZRIfort, FNUfort, Nfort, YRfort, YIfort, NZfort, RZRfort, RZIfort, ASCLEfort, TOLfort, ELIMfort :=
		zksclOrig(ZRR, ZRI, FNU, n, yrfort, yifort, NZ, RZR, RZI, ASCLE, tol, ELIM)

	yramos := make([]float64, len(yr))
	copy(yramos, yr)
	yiamos := make([]float64, len(yi))
	copy(yiamos, yi)
	ZRRamoslib, ZRIamoslib, FNUamoslib, Namoslib, YRamoslib, YIamoslib, NZamoslib, RZRamoslib, RZIamoslib, ASCLEamoslib, TOLamoslib, ELIMamoslib :=
		amoslib.ZksclFort(ZRR, ZRI, FNU, n, yramos, yiamos, NZ, RZR, RZI, ASCLE, tol, ELIM)

	sameF64(t, "zkscl zrr", ZRRfort, ZRRamoslib)
	sameF64(t, "zkscl zri", ZRIfort, ZRIamoslib)
	sameF64(t, "zkscl fnu", FNUfort, FNUamoslib)
	sameInt(t, "zkscl n", Nfort, Namoslib)
	sameInt(t, "zkscl nz", NZfort, NZamoslib)
	sameF64(t, "zkscl rzr", RZRfort, RZRamoslib)
	sameF64(t, "zkscl rzi", RZIfort, RZIamoslib)
	sameF64(t, "zkscl ascle", ASCLEfort, ASCLEamoslib)
	sameF64(t, "zkscl tol", TOLfort, TOLamoslib)
	sameF64(t, "zkscl elim", ELIMfort, ELIMamoslib)

	sameF64SApprox(t, "zkscl yr", YRfort, YRamoslib, 1e-14)
	sameF64SApprox(t, "zkscl yi", YIfort, YIamoslib, 1e-14)
}

func zmlritestFort(t *testing.T, x []float64, is []int, tol float64, n int, yr, yi []float64, kode int) {
	t.Helper()

	ZR := x[0]
	ZI := x[1]
	FNU := x[2]
	KODE := kode
	NZ := is[1]

	yrfort := make([]float64, len(yr))
	copy(yrfort, yr)
	yifort := make([]float64, len(yi))
	copy(yifort, yi)
	ZRfort, ZIfort, FNUfort, KODEfort, Nfort, YRfort, YIfort, NZfort, TOLfort :=
		zmlriOrig(ZR, ZI, FNU, KODE, n, yrfort, yifort, NZ, tol)

	yramos := make([]float64, len(yr))
	copy(yramos, yr)
	yiamos := make([]float64, len(yi))
	copy(yiamos, yi)
	ZRamoslib, ZIamoslib, FNUamoslib, KODEamoslib, Namoslib, YRamoslib, YIamoslib, NZamoslib, TOLamoslib :=
		amoslib.ZmlriFort(ZR, ZI, FNU, KODE, n, yramos, yiamos, NZ, tol)

	sameF64(t, "zmlri zr", ZRfort, ZRamoslib)
	sameF64(t, "zmlri zi", ZIfort, ZIamoslib)
	sameF64(t, "zmlri fnu", FNUfort, FNUamoslib)
	sameInt(t, "zmlri kode", KODEfort, KODEamoslib)
	sameInt(t, "zmlri n", Nfort, Namoslib)
	sameInt(t, "zmlri nz", NZfort, NZamoslib)
	sameF64(t, "zmlri tol", TOLfort, TOLamoslib)

	sameF64S(t, "zmlri yr", YRfort, YRamoslib)
	sameF64S(t, "zmlri yi", YIfort, YIamoslib)
}

func zseritestFort(t *testing.T, x []float64, is []int, tol float64, n int, yr, yi []float64, kode int) {
	t.Helper()

	ZR := x[0]
	ZI := x[1]
	FNU := x[2]
	KODE := kode
	NZ := is[1]
	ELIM := x[3]
	ALIM := x[4]

	yrfort := make([]float64, len(yr))
	copy(yrfort, yr)
	yifort := make([]float64, len(yi))
	copy(yifort, yi)
	ZRfort, ZIfort, FNUfort, KODEfort, Nfort, YRfort, YIfort, NZfort, TOLfort, ELIMfort, ALIMfort :=
		zseriOrig(ZR, ZI, FNU, KODE, n, yrfort, yifort, NZ, tol, ELIM, ALIM)

	yramos := make([]float64, len(yr))
	copy(yramos, yr)
	yiamos := make([]float64, len(yi))
	copy(yiamos, yi)
	y := make([]complex128, len(yramos))
	for i, v := range yramos {
		y[i] = complex(v, yiamos[i])
	}

	ZRamoslib, ZIamoslib, FNUamoslib, KODEamoslib, Namoslib, YRamoslib, YIamoslib, NZamoslib, TOLamoslib, ELIMamoslib, ALIMamoslib :=
		amoslib.ZseriFort(ZR, ZI, FNU, KODE, n, yrfort, yifort, NZ, tol, ELIM, ALIM)

	sameF64(t, "zseri zr", ZRfort, ZRamoslib)
	sameF64(t, "zseri zi", ZIfort, ZIamoslib)
	sameF64(t, "zseri fnu", FNUfort, FNUamoslib)
	sameInt(t, "zseri kode", KODEfort, KODEamoslib)
	sameInt(t, "zseri n", Nfort, Namoslib)
	if *runFailing {
		sameInt(t, "zseri nz", NZfort, NZamoslib)
	}
	sameF64(t, "zseri tol", TOLfort, TOLamoslib)
	sameF64(t, "zseri elim", ELIMfort, ELIMamoslib)
	sameF64(t, "zseri elim", ALIMfort, ALIMamoslib)

	sameF64SApprox(t, "zseri yr", YRfort, YRamoslib, 1e-9)
	sameF64SApprox(t, "zseri yi", YIfort, YIamoslib, 1e-10)
}

func zasyitestFort(t *testing.T, x []float64, is []int, tol float64, n int, yr, yi []float64, kode int) {
	t.Helper()

	ZR := x[0]
	ZI := x[1]
	FNU := x[2]
	KODE := kode
	NZ := is[1]
	ELIM := x[3]
	ALIM := x[4]
	RL := x[5]

	yrfort := make([]float64, len(yr))
	copy(yrfort, yr)
	yifort := make([]float64, len(yi))
	copy(yifort, yi)
	ZRfort, ZIfort, FNUfort, KODEfort, Nfort, YRfort, YIfort, NZfort, RLfort, TOLfort, ELIMfort, ALIMfort :=
		zasyiOrig(ZR, ZI, FNU, KODE, n, yrfort, yifort, NZ, RL, tol, ELIM, ALIM)

	yramos := make([]float64, len(yr))
	copy(yramos, yr)
	yiamos := make([]float64, len(yi))
	copy(yiamos, yi)
	ZRamoslib, ZIamoslib, FNUamoslib, KODEamoslib, Namoslib, YRamoslib, YIamoslib, NZamoslib, RLamoslib, TOLamoslib, ELIMamoslib, ALIMamoslib :=
		amoslib.ZasyiFort(ZR, ZI, FNU, KODE, n, yramos, yiamos, NZ, RL, tol, ELIM, ALIM)

	sameF64(t, "zasyi zr", ZRfort, ZRamoslib)
	sameF64(t, "zasyi zr", ZIfort, ZIamoslib)
	sameF64(t, "zasyi fnu", FNUfort, FNUamoslib)
	sameInt(t, "zasyi kode", KODEfort, KODEamoslib)
	sameInt(t, "zasyi n", Nfort, Namoslib)
	sameInt(t, "zasyi nz", NZfort, NZamoslib)
	sameF64(t, "zasyi rl", RLfort, RLamoslib)
	sameF64(t, "zasyi tol", TOLfort, TOLamoslib)
	sameF64(t, "zasyi elim", ELIMfort, ELIMamoslib)
	sameF64(t, "zasyi alim", ALIMfort, ALIMamoslib)

	sameF64SApprox(t, "zasyi yr", YRfort, YRamoslib, 1e-12)
	sameF64SApprox(t, "zasyi yi", YIfort, YIamoslib, 1e-12)
}

func zbknutestFort(t *testing.T, x []float64, is []int, tol float64, n int, yr, yi []float64, kode int) {
	t.Helper()

	ZR := x[0]
	ZI := x[1]
	FNU := x[2]
	KODE := kode
	NZ := is[1]
	ELIM := x[3]
	ALIM := x[4]

	yrfort := make([]float64, len(yr))
	copy(yrfort, yr)
	yifort := make([]float64, len(yi))
	copy(yifort, yi)
	ZRfort, ZIfort, FNUfort, KODEfort, Nfort, YRfort, YIfort, NZfort, TOLfort, ELIMfort, ALIMfort :=
		zbknuOrig(ZR, ZI, FNU, KODE, n, yrfort, yifort, NZ, tol, ELIM, ALIM)

	yramos := make([]float64, len(yr))
	copy(yramos, yr)
	yiamos := make([]float64, len(yi))
	copy(yiamos, yi)
	ZRamoslib, ZIamoslib, FNUamoslib, KODEamoslib, Namoslib, YRamoslib, YIamoslib, NZamoslib, TOLamoslib, ELIMamoslib, ALIMamoslib :=
		amoslib.ZbknuFort(ZR, ZI, FNU, KODE, n, yramos, yiamos, NZ, tol, ELIM, ALIM)

	sameF64(t, "zbknu zr", ZRfort, ZRamoslib)
	sameF64(t, "zbknu zr", ZIfort, ZIamoslib)
	sameF64(t, "zbknu fnu", FNUfort, FNUamoslib)
	sameInt(t, "zbknu kode", KODEfort, KODEamoslib)
	sameInt(t, "zbknu n", Nfort, Namoslib)
	sameInt(t, "zbknu nz", NZfort, NZamoslib)
	sameF64(t, "zbknu tol", TOLfort, TOLamoslib)
	sameF64(t, "zbknu elim", ELIMfort, ELIMamoslib)
	sameF64(t, "zbknu alim", ALIMfort, ALIMamoslib)

	sameF64SApprox(t, "zbknu yr", YRfort, YRamoslib, 1e-12)
	sameF64SApprox(t, "zbknu yi", YIfort, YIamoslib, 1e-12)
}

func zairytestFort(t *testing.T, x []float64, kode, id int) {
	const tol = 1e-8
	t.Helper()

	ZR := x[0]
	ZI := x[1]
	KODE := kode
	ID := id

	AIRfort, AIIfort, NZfort, IERRfort := zairyOrig(ZR, ZI, ID, KODE)
	AIRamos, AIIamos, NZamos, IERRamos := amoslib.ZairyFort(ZR, ZI, ID, KODE)

	sameF64Approx(t, "zairy air", AIRfort, AIRamos, tol)
	sameF64Approx(t, "zairy aii", AIIfort, AIIamos, tol)
	sameInt(t, "zairy nz", NZfort, NZamos)
	sameInt(t, "zairy ierr", IERRfort, IERRamos)
}

func zacaitestFort(t *testing.T, x []float64, is []int, tol float64, n int, yr, yi []float64, kode int) {
	t.Helper()

	ZR := x[0]
	ZI := x[1]
	FNU := x[2]
	KODE := kode
	NZ := is[1]
	MR := is[2]
	ELIM := x[3]
	ALIM := x[4]
	RL := x[5]

	yrfort := make([]float64, len(yr))
	copy(yrfort, yr)
	yifort := make([]float64, len(yi))
	copy(yifort, yi)
	ZRfort, ZIfort, FNUfort, KODEfort, MRfort, Nfort, YRfort, YIfort, NZfort, RLfort, TOLfort, ELIMfort, ALIMfort :=
		zacaiOrig(ZR, ZI, FNU, KODE, MR, n, yrfort, yifort, NZ, RL, tol, ELIM, ALIM)

	yramos := make([]float64, len(yr))
	copy(yramos, yr)
	yiamos := make([]float64, len(yi))
	copy(yiamos, yi)
	ZRamoslib, ZIamoslib, FNUamoslib, KODEamoslib, MRamoslib, Namoslib, YRamoslib, YIamoslib, NZamoslib, RLamoslib, TOLamoslib, ELIMamoslib, ALIMamoslib :=
		amoslib.ZacaiFort(ZR, ZI, FNU, KODE, MR, n, yramos, yiamos, NZ, RL, tol, ELIM, ALIM)

	sameF64(t, "zacai zr", ZRfort, ZRamoslib)
	sameF64(t, "zacai zi", ZIfort, ZIamoslib)
	sameF64(t, "zacai fnu", FNUfort, FNUamoslib)
	sameInt(t, "zacai kode", KODEfort, KODEamoslib)
	sameInt(t, "zacai mr", MRfort, MRamoslib)
	sameInt(t, "zacai n", Nfort, Namoslib)
	sameInt(t, "zacai nz", NZfort, NZamoslib)
	sameF64(t, "zacai rl", RLfort, RLamoslib)
	sameF64(t, "zacai tol", TOLfort, TOLamoslib)
	sameF64(t, "zacai elim", ELIMfort, ELIMamoslib)
	sameF64(t, "zacai elim", ALIMfort, ALIMamoslib)

	sameF64SApprox(t, "zacai yr", YRfort, YRamoslib, 1e-12)
	sameF64SApprox(t, "zacai yi", YIfort, YIamoslib, 1e-12)
}
