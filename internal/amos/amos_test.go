package amos

import (
	"math"
	"math/rand"
	"strconv"
	"testing"

	"github.com/gonum/mathext/internal/amos/amoslib"
)

type input struct {
	x    []float64
	is   []int
	kode int
	id   int
	yr   []float64
	yi   []float64
	n    int
	tol  float64
}

func randnum(rnd *rand.Rand) float64 {
	r := 2e2 // Fortran has infinite loop if this is set higher than 2e3
	if rnd.Float64() > 0.99 {
		return 0
	}
	return rnd.Float64()*r - r/2
}

func randInput(rnd *rand.Rand) input {
	x := make([]float64, 8)
	for j := range x {
		x[j] = randnum(rnd)
	}
	is := make([]int, 3)
	for j := range is {
		is[j] = rand.Intn(1000)
	}
	kode := rand.Intn(2) + 1
	id := rand.Intn(2)
	n := rand.Intn(5) + 1
	yr := make([]float64, n+1)
	yi := make([]float64, n+1)
	for j := range yr {
		yr[j] = randnum(rnd)
		yi[j] = randnum(rnd)
	}
	tol := 1e-14

	return input{
		x, is, kode, id, yr, yi, n, tol,
	}
}

const nInputs = 100000

func TestAiry(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < nInputs; i++ {
		in := randInput(rnd)
		zairytest(t, in.x, in.kode, in.id)
	}
}

func TestZacai(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < nInputs; i++ {
		in := randInput(rnd)
		zacaitest(t, in.x, in.is, in.tol, in.n, in.yr, in.yi, in.kode)
	}
}

func TestZbknu(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < nInputs; i++ {
		in := randInput(rnd)
		zbknutest(t, in.x, in.is, in.tol, in.n, in.yr, in.yi, in.kode)
	}
}

func TestZasyi(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < nInputs; i++ {
		in := randInput(rnd)
		zasyitest(t, in.x, in.is, in.tol, in.n, in.yr, in.yi, in.kode)
	}
}

func TestZseri(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < nInputs; i++ {
		in := randInput(rnd)
		zseritest(t, in.x, in.is, in.tol, in.n, in.yr, in.yi, in.kode)
	}
}

func TestZmlri(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < nInputs; i++ {
		in := randInput(rnd)
		zmlritest(t, in.x, in.is, in.tol, in.n, in.yr, in.yi, in.kode)
	}
}

func TestZkscl(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < nInputs; i++ {
		in := randInput(rnd)
		zkscltest(t, in.x, in.is, in.tol, in.n, in.yr, in.yi)
	}
}

func TestZuchk(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < nInputs; i++ {
		in := randInput(rnd)
		zuchktest(t, in.x, in.is, in.tol)
	}
}

func TestZs1s2(t *testing.T) {
	rnd := rand.New(rand.NewSource(1))
	for i := 0; i < nInputs; i++ {
		in := randInput(rnd)
		zs1s2test(t, in.x, in.is)
	}
}

func zs1s2test(t *testing.T, x []float64, is []int) {
	ZRR := x[0]
	ZRI := x[1]
	S1R := x[2]
	S1I := x[3]
	S2R := x[4]
	S2I := x[5]
	ASCLE := x[6]
	ALIM := x[7]

	NZ := is[0]
	IUF := is[1]

	ZRRfort, ZRIfort, S1Rfort, S1Ifort, S2Rfort, S2Ifort, NZfort, ASCLEfort, ALIMfort, IUFfort :=
		amoslib.Zs1s2Fort(ZRR, ZRI, S1R, S1I, S2R, S2I, NZ, ASCLE, ALIM, IUF)
	ZRRamos, ZRIamos, S1Ramos, S1Iamos, S2Ramos, S2Iamos, NZamos, ASCLEamos, ALIMamos, IUFamos :=
		Zs1s2(ZRR, ZRI, S1R, S1I, S2R, S2I, NZ, ASCLE, ALIM, IUF)

	SameF64(t, "zs1s2 zrr", ZRRfort, ZRRamos)
	SameF64(t, "zs1s2 zri", ZRIfort, ZRIamos)
	SameF64(t, "zs1s2 s1r", S1Rfort, S1Ramos)
	SameF64(t, "zs1s2 s1i", S1Ifort, S1Iamos)
	SameF64(t, "zs1s2 s2r", S2Rfort, S2Ramos)
	SameF64(t, "zs1s2 s2i", S2Ifort, S2Iamos)
	SameF64(t, "zs1s2 ascle", ASCLEfort, ASCLEamos)
	SameF64(t, "zs1s2 alim", ALIMfort, ALIMamos)
	SameInt(t, "iuf", IUFfort, IUFamos)
	SameInt(t, "nz", NZfort, NZamos)
}

func zuchktest(t *testing.T, x []float64, is []int, tol float64) {
	YR := x[0]
	YI := x[1]
	NZ := is[0]
	ASCLE := x[2]
	TOL := tol

	YRfort, YIfort, NZfort, ASCLEfort, TOLfort := amoslib.ZuchkFort(YR, YI, NZ, ASCLE, TOL)
	YRamos, YIamos, NZamos, ASCLEamos, TOLamos := Zuchk(YR, YI, NZ, ASCLE, TOL)

	SameF64(t, "zuchk yr", YRfort, YRamos)
	SameF64(t, "zuchk yi", YIfort, YIamos)
	SameInt(t, "zuchk nz", NZfort, NZamos)
	SameF64(t, "zuchk ascle", ASCLEfort, ASCLEamos)
	SameF64(t, "zuchk tol", TOLfort, TOLamos)
}

func zkscltest(t *testing.T, x []float64, is []int, tol float64, n int, yr, yi []float64) {
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
		amoslib.ZksclFort(ZRR, ZRI, FNU, n, yrfort[1:], yifort[1:], NZ, RZR, RZI, ASCLE, tol, ELIM)
	YRfort2 := make([]float64, len(yrfort))
	YRfort2[0] = yrfort[0]
	copy(YRfort2[1:], YRfort)
	YIfort2 := make([]float64, len(yifort))
	YIfort2[0] = yifort[0]
	copy(YIfort2[1:], YIfort)

	yramos := make([]float64, len(yr))
	copy(yramos, yr)
	yiamos := make([]float64, len(yi))
	copy(yiamos, yi)
	ZRRamos, ZRIamos, FNUamos, Namos, YRamos, YIamos, NZamos, RZRamos, RZIamos, ASCLEamos, TOLamos, ELIMamos :=
		Zkscl(ZRR, ZRI, FNU, n, yramos, yiamos, NZ, RZR, RZI, ASCLE, tol, ELIM)

	SameF64(t, "zkscl zrr", ZRRfort, ZRRamos)
	SameF64(t, "zkscl zri", ZRIfort, ZRIamos)
	SameF64(t, "zkscl fnu", FNUfort, FNUamos)
	SameInt(t, "zkscl n", Nfort, Namos)
	SameInt(t, "zkscl nz", NZfort, NZamos)
	SameF64(t, "zkscl rzr", RZRfort, RZRamos)
	SameF64(t, "zkscl rzi", RZIfort, RZIamos)
	SameF64(t, "zkscl ascle", ASCLEfort, ASCLEamos)
	SameF64(t, "zkscl tol", TOLfort, TOLamos)
	SameF64(t, "zkscl elim", ELIMfort, ELIMamos)

	SameF64S(t, "zkscl yr", YRfort2, YRamos)
	SameF64S(t, "zkscl yi", YIfort2, YIamos)
}

func zmlritest(t *testing.T, x []float64, is []int, tol float64, n int, yr, yi []float64, kode int) {
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
		amoslib.ZmlriFort(ZR, ZI, FNU, KODE, n, yrfort[1:], yifort[1:], NZ, tol)
	YRfort2 := make([]float64, len(yrfort))
	YRfort2[0] = yrfort[0]
	copy(YRfort2[1:], YRfort)
	YIfort2 := make([]float64, len(yifort))
	YIfort2[0] = yifort[0]
	copy(YIfort2[1:], YIfort)

	yramos := make([]float64, len(yr))
	copy(yramos, yr)
	yiamos := make([]float64, len(yi))
	copy(yiamos, yi)
	ZRamos, ZIamos, FNUamos, KODEamos, Namos, YRamos, YIamos, NZamos, TOLamos :=
		Zmlri(ZR, ZI, FNU, KODE, n, yramos, yiamos, NZ, tol)

	SameF64(t, "zmlri zr", ZRfort, ZRamos)
	SameF64(t, "zmlri zi", ZIfort, ZIamos)
	SameF64(t, "zmlri fnu", FNUfort, FNUamos)
	SameInt(t, "zmlri kode", KODEfort, KODEamos)
	SameInt(t, "zmlri n", Nfort, Namos)
	SameInt(t, "zmlri nz", NZfort, NZamos)
	SameF64(t, "zmlri tol", TOLfort, TOLamos)

	SameF64S(t, "zmlri yr", YRfort2, YRamos)
	SameF64S(t, "zmlri yi", YIfort2, YIamos)
}

func zseritest(t *testing.T, x []float64, is []int, tol float64, n int, yr, yi []float64, kode int) {
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
		amoslib.ZseriFort(ZR, ZI, FNU, KODE, n, yrfort[1:], yifort[1:], NZ, tol, ELIM, ALIM)
	YRfort2 := make([]float64, len(yrfort))
	YRfort2[0] = yrfort[0]
	copy(YRfort2[1:], YRfort)
	YIfort2 := make([]float64, len(yifort))
	YIfort2[0] = yifort[0]
	copy(YIfort2[1:], YIfort)

	yramos := make([]float64, len(yr))
	copy(yramos, yr)
	yiamos := make([]float64, len(yi))
	copy(yiamos, yi)
	ZRamos, ZIamos, FNUamos, KODEamos, Namos, YRamos, YIamos, NZamos, TOLamos, ELIMamos, ALIMamos :=
		Zseri(ZR, ZI, FNU, KODE, n, yramos, yiamos, NZ, tol, ELIM, ALIM)

	SameF64(t, "zseri zr", ZRfort, ZRamos)
	SameF64(t, "zseri zi", ZIfort, ZIamos)
	SameF64(t, "zseri fnu", FNUfort, FNUamos)
	SameInt(t, "zseri kode", KODEfort, KODEamos)
	SameInt(t, "zseri n", Nfort, Namos)
	SameInt(t, "zseri nz", NZfort, NZamos)
	SameF64(t, "zseri tol", TOLfort, TOLamos)
	SameF64(t, "zseri elim", ELIMfort, ELIMamos)
	SameF64(t, "zseri elim", ALIMfort, ALIMamos)

	SameF64S(t, "zseri yr", YRfort2, YRamos)
	SameF64S(t, "zseri yi", YIfort2, YIamos)
}

func zasyitest(t *testing.T, x []float64, is []int, tol float64, n int, yr, yi []float64, kode int) {
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
		amoslib.ZasyiFort(ZR, ZI, FNU, KODE, n, yrfort[1:], yifort[1:], NZ, RL, tol, ELIM, ALIM)
	YRfort2 := make([]float64, len(yrfort))
	YRfort2[0] = yrfort[0]
	copy(YRfort2[1:], YRfort)
	YIfort2 := make([]float64, len(yifort))
	YIfort2[0] = yifort[0]
	copy(YIfort2[1:], YIfort)

	yramos := make([]float64, len(yr))
	copy(yramos, yr)
	yiamos := make([]float64, len(yi))
	copy(yiamos, yi)
	ZRamos, ZIamos, FNUamos, KODEamos, Namos, YRamos, YIamos, NZamos, RLamos, TOLamos, ELIMamos, ALIMamos :=
		Zasyi(ZR, ZI, FNU, KODE, n, yramos, yiamos, NZ, RL, tol, ELIM, ALIM)

	SameF64(t, "zasyi zr", ZRfort, ZRamos)
	SameF64(t, "zasyi zr", ZIfort, ZIamos)
	SameF64(t, "zasyi fnu", FNUfort, FNUamos)
	SameInt(t, "zasyi kode", KODEfort, KODEamos)
	SameInt(t, "zasyi n", Nfort, Namos)
	SameInt(t, "zasyi nz", NZfort, NZamos)
	SameF64(t, "zasyi rl", RLfort, RLamos)
	SameF64(t, "zasyi tol", TOLfort, TOLamos)
	SameF64(t, "zasyi elim", ELIMfort, ELIMamos)
	SameF64(t, "zasyi alim", ALIMfort, ALIMamos)

	SameF64S(t, "zasyi yr", YRfort2, YRamos)
	SameF64S(t, "zasyi yi", YIfort2, YIamos)
}

func zbknutest(t *testing.T, x []float64, is []int, tol float64, n int, yr, yi []float64, kode int) {
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
		amoslib.ZbknuFort(ZR, ZI, FNU, KODE, n, yrfort[1:], yifort[1:], NZ, tol, ELIM, ALIM)
	YRfort2 := make([]float64, len(yrfort))
	YRfort2[0] = yrfort[0]
	copy(YRfort2[1:], YRfort)
	YIfort2 := make([]float64, len(yifort))
	YIfort2[0] = yifort[0]
	copy(YIfort2[1:], YIfort)

	yramos := make([]float64, len(yr))
	copy(yramos, yr)
	yiamos := make([]float64, len(yi))
	copy(yiamos, yi)
	ZRamos, ZIamos, FNUamos, KODEamos, Namos, YRamos, YIamos, NZamos, TOLamos, ELIMamos, ALIMamos :=
		Zbknu(ZR, ZI, FNU, KODE, n, yramos, yiamos, NZ, tol, ELIM, ALIM)

	SameF64(t, "zbknu zr", ZRfort, ZRamos)
	SameF64(t, "zbknu zr", ZIfort, ZIamos)
	SameF64(t, "zbknu fnu", FNUfort, FNUamos)
	SameInt(t, "zbknu kode", KODEfort, KODEamos)
	SameInt(t, "zbknu n", Nfort, Namos)
	SameInt(t, "zbknu nz", NZfort, NZamos)
	SameF64(t, "zbknu tol", TOLfort, TOLamos)
	SameF64(t, "zbknu elim", ELIMfort, ELIMamos)
	SameF64(t, "zbknu alim", ALIMfort, ALIMamos)

	SameF64S(t, "zbknu yr", YRfort2, YRamos)
	SameF64S(t, "zbknu yi", YIfort2, YIamos)
}

func zairytest(t *testing.T, x []float64, kode, id int) {
	ZR := x[0]
	ZI := x[1]
	KODE := kode
	ID := id

	AIRfort, AIIfort, NZfort := amoslib.ZairyFort(ZR, ZI, ID, KODE)
	AIRamos, AIIamos, NZamos := Zairy(ZR, ZI, ID, KODE)

	SameF64(t, "zairy air", AIRfort, AIRamos)
	SameF64(t, "zairy aii", AIIfort, AIIamos)
	SameInt(t, "zairy nz", NZfort, NZamos)
}

func zacaitest(t *testing.T, x []float64, is []int, tol float64, n int, yr, yi []float64, kode int) {
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
		amoslib.ZacaiFort(ZR, ZI, FNU, KODE, MR, n, yrfort[1:], yifort[1:], NZ, RL, tol, ELIM, ALIM)
	YRfort2 := make([]float64, len(yrfort))
	YRfort2[0] = yrfort[0]
	copy(YRfort2[1:], YRfort)
	YIfort2 := make([]float64, len(yifort))
	YIfort2[0] = yifort[0]
	copy(YIfort2[1:], YIfort)

	yramos := make([]float64, len(yr))
	copy(yramos, yr)
	yiamos := make([]float64, len(yi))
	copy(yiamos, yi)
	ZRamos, ZIamos, FNUamos, KODEamos, MRamos, Namos, YRamos, YIamos, NZamos, RLamos, TOLamos, ELIMamos, ALIMamos :=
		Zacai(ZR, ZI, FNU, KODE, MR, n, yramos, yiamos, NZ, RL, tol, ELIM, ALIM)

	SameF64(t, "zacai zr", ZRfort, ZRamos)
	SameF64(t, "zacai zi", ZIfort, ZIamos)
	SameF64(t, "zacai fnu", FNUfort, FNUamos)
	SameInt(t, "zacai kode", KODEfort, KODEamos)
	SameInt(t, "zacai mr", MRfort, MRamos)
	SameInt(t, "zacai n", Nfort, Namos)
	SameInt(t, "zacai nz", NZfort, NZamos)
	SameF64(t, "zacai rl", RLfort, RLamos)
	SameF64(t, "zacai tol", TOLfort, TOLamos)
	SameF64(t, "zacai elim", ELIMfort, ELIMamos)
	SameF64(t, "zacai elim", ALIMfort, ALIMamos)

	SameF64S(t, "zacai yr", YRfort2, YRamos)
	SameF64S(t, "zacai yi", YIfort2, YIamos)
}

func SameF64(t *testing.T, str string, c, native float64) {
	if math.IsNaN(c) && math.IsNaN(native) {
		return
	}
	if c == native {
		return
	}
	cb := math.Float64bits(c)
	nb := math.Float64bits(native)
	t.Errorf("Case %s: Float64 mismatch. c = %v, native = %v\n cb: %v, nb: %v\n", str, c, native, cb, nb)
}

func SameInt(t *testing.T, str string, c, native int) {
	if c != native {
		t.Errorf("Case %s: Int mismatch. c = %v, native = %v.", str, c, native)
	}
}

func SameF64S(t *testing.T, str string, c, native []float64) {
	if len(c) != len(native) {
		panic(str)
	}
	for i, v := range c {
		SameF64(t, str+"_idx_"+strconv.Itoa(i), v, native[i])
	}
}
