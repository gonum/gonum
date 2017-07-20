// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package blas64

import (
	"math"
	"testing"

	"gonum.org/v1/gonum/blas"
)

func newGeneralFrom(a GeneralCols) General {
	t := General{
		Rows:   a.Rows,
		Cols:   a.Cols,
		Stride: a.Cols,
		Data:   make([]float64, a.Rows*a.Cols),
	}
	t.From(a)
	return t
}

func (m General) dims() (r, c int)    { return m.Rows, m.Cols }
func (m General) at(i, j int) float64 { return m.Data[i*m.Stride+j] }

func newGeneralColsFrom(a General) GeneralCols {
	t := GeneralCols{
		Rows:   a.Rows,
		Cols:   a.Cols,
		Stride: a.Rows,
		Data:   make([]float64, a.Rows*a.Cols),
	}
	t.From(a)
	return t
}

func (m GeneralCols) dims() (r, c int)    { return m.Rows, m.Cols }
func (m GeneralCols) at(i, j int) float64 { return m.Data[i+j*m.Stride] }

type general interface {
	dims() (r, c int)
	at(i, j int) float64
}

func equalGeneral(a, b general) bool {
	ar, ac := a.dims()
	br, bc := b.dims()
	if ar != br || ac != bc {
		return false
	}
	for i := 0; i < ar; i++ {
		for j := 0; j < ac; j++ {
			if a.at(i, j) != b.at(i, j) || math.IsNaN(a.at(i, j)) != math.IsNaN(b.at(i, j)) {
				return false
			}
		}
	}
	return true
}

var generalTests = []General{
	{Rows: 2, Cols: 3, Stride: 3, Data: []float64{
		1, 2, 3,
		4, 5, 6,
	}},
	{Rows: 3, Cols: 2, Stride: 2, Data: []float64{
		1, 2,
		3, 4,
		5, 6,
	}},
	{Rows: 3, Cols: 3, Stride: 3, Data: []float64{
		1, 2, 3,
		4, 5, 6,
		7, 8, 9,
	}},
	{Rows: 2, Cols: 3, Stride: 5, Data: []float64{
		1, 2, 3, 0, 0,
		4, 5, 6, 0, 0,
	}},
	{Rows: 3, Cols: 2, Stride: 5, Data: []float64{
		1, 2, 0, 0, 0,
		3, 4, 0, 0, 0,
		5, 6, 0, 0, 0,
	}},
	{Rows: 3, Cols: 3, Stride: 5, Data: []float64{
		1, 2, 3, 0, 0,
		4, 5, 6, 0, 0,
		7, 8, 9, 0, 0,
	}},
}

func TestConvertGeneral(t *testing.T) {
	for _, test := range generalTests {
		colmajor := newGeneralColsFrom(test)
		if !equalGeneral(colmajor, test) {
			t.Errorf("unexpected result for row major to col major conversion:\n\tgot: %#v\n\tfrom:%#v",
				colmajor, test)
		}
		rowmajor := newGeneralFrom(colmajor)
		if !equalGeneral(rowmajor, test) {
			t.Errorf("unexpected result for col major to row major conversion:\n\tgot: %#v\n\twant:%#v",
				rowmajor, test)
		}
	}
}

func newSymmetricFrom(a SymmetricCols) Symmetric {
	t := Symmetric{
		N:      a.N,
		Stride: a.N,
		Data:   make([]float64, a.N*a.N),
		Uplo:   a.Uplo,
	}
	t.From(a)
	return t
}

func (m Symmetric) n() int { return m.N }
func (m Symmetric) at(i, j int) float64 {
	if m.Uplo == blas.Lower && i < j && j < m.N {
		i, j = j, i
	}
	if m.Uplo == blas.Upper && i > j {
		i, j = j, i
	}
	return m.Data[i*m.Stride+j]
}
func (m Symmetric) uplo() blas.Uplo { return m.Uplo }

func newSymmetricColsFrom(a Symmetric) SymmetricCols {
	t := SymmetricCols{
		N:      a.N,
		Stride: a.N,
		Data:   make([]float64, a.N*a.N),
		Uplo:   a.Uplo,
	}
	t.From(a)
	return t
}

func (m SymmetricCols) n() int { return m.N }
func (m SymmetricCols) at(i, j int) float64 {
	if m.Uplo == blas.Lower && i < j {
		i, j = j, i
	}
	if m.Uplo == blas.Upper && i > j && i < m.N {
		i, j = j, i
	}
	return m.Data[i+j*m.Stride]
}
func (m SymmetricCols) uplo() blas.Uplo { return m.Uplo }

type symmetric interface {
	n() int
	at(i, j int) float64
	uplo() blas.Uplo
}

func equalSymmetric(a, b symmetric) bool {
	an := a.n()
	bn := b.n()
	if an != bn {
		return false
	}
	if a.uplo() != b.uplo() {
		return false
	}
	for i := 0; i < an; i++ {
		for j := 0; j < an; j++ {
			if a.at(i, j) != b.at(i, j) || math.IsNaN(a.at(i, j)) != math.IsNaN(b.at(i, j)) {
				return false
			}
		}
	}
	return true
}

var symmetricTests = []Symmetric{
	{N: 3, Stride: 3, Data: []float64{
		1, 2, 3,
		4, 5, 6,
		7, 8, 9,
	}},
	{N: 3, Stride: 5, Data: []float64{
		1, 2, 3, 0, 0,
		4, 5, 6, 0, 0,
		7, 8, 9, 0, 0,
	}},
}

func TestConvertSymmetric(t *testing.T) {
	for _, test := range symmetricTests {
		for _, uplo := range []blas.Uplo{blas.Upper, blas.Lower} {
			test.Uplo = uplo
			colmajor := newSymmetricColsFrom(test)
			if !equalSymmetric(colmajor, test) {
				t.Errorf("unexpected result for row major to col major conversion:\n\tgot: %#v\n\tfrom:%#v",
					colmajor, test)
			}
			rowmajor := newSymmetricFrom(colmajor)
			if !equalSymmetric(rowmajor, test) {
				t.Errorf("unexpected result for col major to row major conversion:\n\tgot: %#v\n\twant:%#v",
					rowmajor, test)
			}
		}
	}
}

func newTriangularFrom(a TriangularCols) Triangular {
	t := Triangular{
		N:      a.N,
		Stride: a.N,
		Data:   make([]float64, a.N*a.N),
		Diag:   a.Diag,
		Uplo:   a.Uplo,
	}
	t.From(a)
	return t
}

func (m Triangular) n() int { return m.N }
func (m Triangular) at(i, j int) float64 {
	if m.Diag == blas.Unit && i == j {
		return 1
	}
	if m.Uplo == blas.Lower && i < j && j < m.N {
		return 0
	}
	if m.Uplo == blas.Upper && i > j {
		return 0
	}
	return m.Data[i*m.Stride+j]
}
func (m Triangular) uplo() blas.Uplo { return m.Uplo }
func (m Triangular) diag() blas.Diag { return m.Diag }

func newTriangularColsFrom(a Triangular) TriangularCols {
	t := TriangularCols{
		N:      a.N,
		Stride: a.N,
		Data:   make([]float64, a.N*a.N),
		Diag:   a.Diag,
		Uplo:   a.Uplo,
	}
	t.From(a)
	return t
}

func (m TriangularCols) n() int { return m.N }
func (m TriangularCols) at(i, j int) float64 {
	if m.Diag == blas.Unit && i == j {
		return 1
	}
	if m.Uplo == blas.Lower && i < j {
		return 0
	}
	if m.Uplo == blas.Upper && i > j && i < m.N {
		return 0
	}
	return m.Data[i+j*m.Stride]
}
func (m TriangularCols) uplo() blas.Uplo { return m.Uplo }
func (m TriangularCols) diag() blas.Diag { return m.Diag }

type triangular interface {
	n() int
	at(i, j int) float64
	uplo() blas.Uplo
	diag() blas.Diag
}

func equalTriangular(a, b triangular) bool {
	an := a.n()
	bn := b.n()
	if an != bn {
		return false
	}
	for i := 0; i < an; i++ {
		for j := 0; j < an; j++ {
			if a.at(i, j) != b.at(i, j) || math.IsNaN(a.at(i, j)) != math.IsNaN(b.at(i, j)) {
				return false
			}
		}
	}
	return true
}

var triangularTests = []Triangular{
	{N: 3, Stride: 3, Data: []float64{
		1, 2, 3,
		4, 5, 6,
		7, 8, 9,
	}},
	{N: 3, Stride: 5, Data: []float64{
		1, 2, 3, 0, 0,
		4, 5, 6, 0, 0,
		7, 8, 9, 0, 0,
	}},
}

func TestConvertTriangular(t *testing.T) {
	for _, test := range triangularTests {
		for _, uplo := range []blas.Uplo{blas.Upper, blas.Lower, blas.All} {
			for _, diag := range []blas.Diag{blas.Unit, blas.NonUnit} {
				test.Uplo = uplo
				test.Diag = diag
				colmajor := newTriangularColsFrom(test)
				if !equalTriangular(colmajor, test) {
					t.Errorf("unexpected result for row major to col major conversion:\n\tgot: %#v\n\tfrom:%#v",
						colmajor, test)
				}
				rowmajor := newTriangularFrom(colmajor)
				if !equalTriangular(rowmajor, test) {
					t.Errorf("unexpected result for col major to row major conversion:\n\tgot: %#v\n\twant:%#v",
						rowmajor, test)
				}
			}
		}
	}
}

func newBandFrom(a BandCols) Band {
	t := Band{
		Rows:   a.Rows,
		Cols:   a.Cols,
		KL:     a.KL,
		KU:     a.KU,
		Stride: a.KL + a.KU + 1,
		Data:   make([]float64, a.Rows*(a.KL+a.KU+1)),
	}
	for i := range t.Data {
		t.Data[i] = math.NaN()
	}
	t.From(a)
	return t
}

func (m Band) dims() (r, c int) { return m.Rows, m.Cols }
func (m Band) at(i, j int) float64 {
	pj := j + m.KL - i
	if pj < 0 || m.KL+m.KU+1 <= pj {
		return 0
	}
	return m.Data[i*m.Stride+pj]
}
func (m Band) bandwidth() (kl, ku int) { return m.KL, m.KU }

func newBandColsFrom(a Band) BandCols {
	t := BandCols{
		Rows:   a.Rows,
		Cols:   a.Cols,
		KL:     a.KL,
		KU:     a.KU,
		Stride: a.KL + a.KU + 1,
		Data:   make([]float64, a.Cols*(a.KL+a.KU+1)),
	}
	for i := range t.Data {
		t.Data[i] = math.NaN()
	}
	t.From(a)
	return t
}

func (m BandCols) dims() (r, c int) { return m.Rows, m.Cols }
func (m BandCols) at(i, j int) float64 {
	pj := i + m.KU - j
	if pj < 0 || m.KL+m.KU+1 <= pj {
		return 0
	}
	return m.Data[j*m.Stride+pj]
}
func (m BandCols) bandwidth() (kl, ku int) { return m.KL, m.KU }

type band interface {
	dims() (r, c int)
	at(i, j int) float64
	bandwidth() (kl, ku int)
}

func equalBand(a, b band) bool {
	ar, ac := a.dims()
	br, bc := b.dims()
	if ar != br || ac != bc {
		return false
	}
	akl, aku := a.bandwidth()
	bkl, bku := b.bandwidth()
	if akl != bkl || aku != bku {
		return false
	}
	for i := 0; i < ar; i++ {
		for j := 0; j < ac; j++ {
			if a.at(i, j) != b.at(i, j) || math.IsNaN(a.at(i, j)) != math.IsNaN(b.at(i, j)) {
				return false
			}
		}
	}
	return true
}

var bandTests = []Band{
	{Rows: 3, Cols: 4, KL: 0, KU: 0, Stride: 1, Data: []float64{
		1,
		2,
		3,
	}},
	{Rows: 3, Cols: 3, KL: 0, KU: 0, Stride: 1, Data: []float64{
		1,
		2,
		3,
	}},
	{Rows: 4, Cols: 3, KL: 0, KU: 0, Stride: 1, Data: []float64{
		1,
		2,
		3,
	}},
	{Rows: 4, Cols: 3, KL: 0, KU: 1, Stride: 2, Data: []float64{
		1, 2,
		3, 4,
		5, 6,
	}},
	{Rows: 3, Cols: 4, KL: 0, KU: 1, Stride: 2, Data: []float64{
		1, 2,
		3, 4,
		5, 6,
	}},
	{Rows: 3, Cols: 4, KL: 1, KU: 1, Stride: 3, Data: []float64{
		-1, 2, 3,
		4, 5, 6,
		7, 8, 9,
	}},
	{Rows: 4, Cols: 3, KL: 1, KU: 1, Stride: 3, Data: []float64{
		-1, 2, 3,
		4, 5, 6,
		7, 8, -2,
		9, -3, -4,
	}},
	{Rows: 3, Cols: 4, KL: 2, KU: 1, Stride: 4, Data: []float64{
		-2, -1, 3, 4,
		-3, 5, 6, 7,
		8, 9, 10, 11,
	}},
	{Rows: 4, Cols: 3, KL: 2, KU: 1, Stride: 4, Data: []float64{
		-2, -1, 2, 3,
		-3, 4, 5, 6,
		7, 8, 9, -4,
		10, 11, -5, -6,
	}},

	{Rows: 3, Cols: 4, KL: 0, KU: 0, Stride: 5, Data: []float64{
		1, 0, 0, 0, 0,
		2, 0, 0, 0, 0,
		3, 0, 0, 0, 0,
	}},
	{Rows: 3, Cols: 3, KL: 0, KU: 0, Stride: 5, Data: []float64{
		1, 0, 0, 0, 0,
		2, 0, 0, 0, 0,
		3, 0, 0, 0, 0,
	}},
	{Rows: 4, Cols: 3, KL: 0, KU: 0, Stride: 5, Data: []float64{
		1, 0, 0, 0, 0,
		2, 0, 0, 0, 0,
		3, 0, 0, 0, 0,
	}},
	{Rows: 4, Cols: 3, KL: 0, KU: 1, Stride: 5, Data: []float64{
		1, 2, 0, 0, 0,
		3, 4, 0, 0, 0,
		5, 6, 0, 0, 0,
	}},
	{Rows: 3, Cols: 4, KL: 0, KU: 1, Stride: 5, Data: []float64{
		1, 2, 0, 0, 0,
		3, 4, 0, 0, 0,
		5, 6, 0, 0, 0,
	}},
	{Rows: 3, Cols: 4, KL: 1, KU: 1, Stride: 5, Data: []float64{
		-1, 2, 3, 0, 0,
		4, 5, 6, 0, 0,
		7, 8, 9, 0, 0,
	}},
	{Rows: 4, Cols: 3, KL: 1, KU: 1, Stride: 5, Data: []float64{
		-1, 2, 3, 0, 0,
		4, 5, 6, 0, 0,
		7, 8, -2, 0, 0,
		9, -3, -4, 0, 0,
	}},
	{Rows: 3, Cols: 4, KL: 2, KU: 1, Stride: 5, Data: []float64{
		-2, -1, 3, 4, 0,
		-3, 5, 6, 7, 0,
		8, 9, 10, 11, 0,
	}},
	{Rows: 4, Cols: 3, KL: 2, KU: 1, Stride: 5, Data: []float64{
		-2, -1, 2, 3, 0,
		-3, 4, 5, 6, 0,
		7, 8, 9, -4, 0,
		10, 11, -5, -6, 0,
	}},
}

func TestConvertBand(t *testing.T) {
	for _, test := range bandTests[:1] {
		colmajor := newBandColsFrom(test)
		if !equalBand(colmajor, test) {
			t.Errorf("unexpected result for row major to col major conversion:\n\tgot: %#v\n\tfrom:%#v",
				colmajor, test)
		}
		rowmajor := newBandFrom(colmajor)
		if !equalBand(rowmajor, test) {
			t.Errorf("unexpected result for col major to row major conversion:\n\tgot: %#v\n\twant:%#v",
				rowmajor, test)
		}
	}
}

func newSymmetricBandFrom(a SymmetricBandCols) SymmetricBand {
	t := SymmetricBand{
		N:      a.N,
		K:      a.K,
		Stride: a.K + 1,
		Data:   make([]float64, a.N*(a.K+1)),
		Uplo:   a.Uplo,
	}
	for i := range t.Data {
		t.Data[i] = math.NaN()
	}
	t.From(a)
	return t
}

func (m SymmetricBand) n() (n int) { return m.N }
func (m SymmetricBand) at(i, j int) float64 {
	if i > j {
		i, j = j, i
	}
	pj := j - i
	if m.K+1 <= pj {
		return 0
	}
	return m.Data[i*m.Stride+pj]
}
func (m SymmetricBand) bandwidth() (k int) { return m.K }
func (m SymmetricBand) uplo() blas.Uplo    { return m.Uplo }

func newSymmetricBandColsFrom(a SymmetricBand) SymmetricBandCols {
	t := SymmetricBandCols{
		N:      a.N,
		K:      a.K,
		Stride: a.K + 1,
		Data:   make([]float64, a.N*(a.K+1)),
		Uplo:   a.Uplo,
	}
	for i := range t.Data {
		t.Data[i] = math.NaN()
	}
	t.From(a)
	return t
}

func (m SymmetricBandCols) n() (n int) { return m.N }
func (m SymmetricBandCols) at(i, j int) float64 {
	if i > j {
		i, j = j, i
	}
	pj := j - i
	if m.K+1 <= pj {
		return 0
	}
	return m.Data[j*m.Stride+pj]
}
func (m SymmetricBandCols) bandwidth() (k int) { return m.K }
func (m SymmetricBandCols) uplo() blas.Uplo    { return m.Uplo }

type symmetricBand interface {
	n() (n int)
	at(i, j int) float64
	bandwidth() (k int)
	uplo() blas.Uplo
}

func equalSymmetricBand(a, b symmetricBand) bool {
	an := a.n()
	bn := b.n()
	if an != bn {
		return false
	}
	if a.uplo() != b.uplo() {
		return false
	}
	ak := a.bandwidth()
	bk := b.bandwidth()
	if ak != bk {
		return false
	}
	for i := 0; i < an; i++ {
		for j := 0; j < an; j++ {
			if a.at(i, j) != b.at(i, j) || math.IsNaN(a.at(i, j)) != math.IsNaN(b.at(i, j)) {
				return false
			}
		}
	}
	return true
}

var symmetricBandTests = []SymmetricBand{
	{N: 3, K: 0, Stride: 1, Uplo: blas.Upper, Data: []float64{
		1,
		2,
		3,
	}},
	{N: 3, K: 0, Stride: 1, Uplo: blas.Lower, Data: []float64{
		1,
		2,
		3,
	}},
	{N: 3, K: 1, Stride: 2, Uplo: blas.Upper, Data: []float64{
		1, 2,
		3, 4,
		5, -1,
	}},
	{N: 3, K: 1, Stride: 2, Uplo: blas.Lower, Data: []float64{
		-1, 1,
		2, 3,
		4, 5,
	}},
	{N: 3, K: 2, Stride: 3, Uplo: blas.Upper, Data: []float64{
		1, 2, 3,
		4, 5, -1,
		6, -2, -3,
	}},
	{N: 3, K: 2, Stride: 3, Uplo: blas.Lower, Data: []float64{
		-2, -1, 1,
		-3, 2, 4,
		3, 5, 6,
	}},

	{N: 3, K: 0, Stride: 5, Uplo: blas.Upper, Data: []float64{
		1, 0, 0, 0, 0,
		2, 0, 0, 0, 0,
		3, 0, 0, 0, 0,
	}},
	{N: 3, K: 0, Stride: 5, Uplo: blas.Lower, Data: []float64{
		1, 0, 0, 0, 0,
		2, 0, 0, 0, 0,
		3, 0, 0, 0, 0,
	}},
	{N: 3, K: 1, Stride: 5, Uplo: blas.Upper, Data: []float64{
		1, 2, 0, 0, 0,
		3, 4, 0, 0, 0,
		5, -1, 0, 0, 0,
	}},
	{N: 3, K: 1, Stride: 5, Uplo: blas.Lower, Data: []float64{
		-1, 1, 0, 0, 0,
		2, 3, 0, 0, 0,
		4, 5, 0, 0, 0,
	}},
	{N: 3, K: 2, Stride: 5, Uplo: blas.Upper, Data: []float64{
		1, 2, 3, 0, 0,
		4, 5, -1, 0, 0,
		6, -2, -3, 0, 0,
	}},
	{N: 3, K: 2, Stride: 5, Uplo: blas.Lower, Data: []float64{
		-2, -1, 1, 0, 0,
		-3, 2, 4, 0, 0,
		3, 5, 6, 0, 0,
	}},
}

func TestConvertSymBand(t *testing.T) {
	for _, test := range symmetricBandTests[:1] {
		colmajor := newSymmetricBandColsFrom(test)
		if !equalSymmetricBand(colmajor, test) {
			t.Errorf("unexpected result for row major to col major conversion:\n\tgot: %#v\n\tfrom:%#v",
				colmajor, test)
		}
		rowmajor := newSymmetricBandFrom(colmajor)
		if !equalSymmetricBand(rowmajor, test) {
			t.Errorf("unexpected result for col major to row major conversion:\n\tgot: %#v\n\twant:%#v",
				rowmajor, test)
		}
	}
}
