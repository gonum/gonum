// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat

import "gonum.org/v1/gonum/blas/blas64"

func asBasicMatrix(d *Dense) *basicMatrix              { return (*basicMatrix)(d) }
func asBasicVector(d *VecDense) *basicVector           { return (*basicVector)(d) }
func asBasicSymmetric(s *SymDense) *basicSymmetric     { return (*basicSymmetric)(s) }
func asBasicTriangular(t *TriDense) *basicTriangular   { return (*basicTriangular)(t) }
func asBasicBanded(b *BandDense) *basicBanded          { return (*basicBanded)(b) }
func asBasicSymBanded(s *SymBandDense) *basicSymBanded { return (*basicSymBanded)(s) }
func asBasicTriBanded(t *TriBandDense) *basicTriBanded { return (*basicTriBanded)(t) }
func asBasicDiagonal(d *DiagDense) *basicDiagonal      { return (*basicDiagonal)(d) }

type basicMatrix Dense

var _ Matrix = &basicMatrix{}

func (m *basicMatrix) At(r, c int) float64 { return (*Dense)(m).At(r, c) }
func (m *basicMatrix) Dims() (r, c int)    { return (*Dense)(m).Dims() }
func (m *basicMatrix) T() Matrix           { return Transpose{m} }

type rawMatrix struct {
	*basicMatrix
}

func (a *rawMatrix) RawMatrix() blas64.General {
	return a.mat
}

type basicVector VecDense

var _ Vector = &basicVector{}

func (v *basicVector) At(r, c int) float64 { return (*VecDense)(v).At(r, c) }
func (v *basicVector) Dims() (r, c int)    { return (*VecDense)(v).Dims() }
func (v *basicVector) T() Matrix           { return Transpose{v} }
func (v *basicVector) AtVec(i int) float64 { return (*VecDense)(v).AtVec(i) }
func (v *basicVector) Len() int            { return (*VecDense)(v).Len() }

type rawVector struct {
	*basicVector
}

func (v *rawVector) RawVector() blas64.Vector {
	return v.mat
}

type basicSymmetric SymDense

var _ Symmetric = &basicSymmetric{}

func (m *basicSymmetric) At(r, c int) float64 { return (*SymDense)(m).At(r, c) }
func (m *basicSymmetric) Dims() (r, c int)    { return (*SymDense)(m).Dims() }
func (m *basicSymmetric) T() Matrix           { return m }
func (m *basicSymmetric) SymmetricDim() int   { return (*SymDense)(m).SymmetricDim() }

type basicTriangular TriDense

var _ Triangular = &basicTriangular{}

func (m *basicTriangular) At(r, c int) float64      { return (*TriDense)(m).At(r, c) }
func (m *basicTriangular) Dims() (r, c int)         { return (*TriDense)(m).Dims() }
func (m *basicTriangular) T() Matrix                { return Transpose{m} }
func (m *basicTriangular) Triangle() (int, TriKind) { return (*TriDense)(m).Triangle() }
func (m *basicTriangular) TTri() Triangular         { return TransposeTri{m} }

type basicBanded BandDense

var _ Banded = &basicBanded{}

func (m *basicBanded) At(r, c int) float64     { return (*BandDense)(m).At(r, c) }
func (m *basicBanded) Dims() (r, c int)        { return (*BandDense)(m).Dims() }
func (m *basicBanded) T() Matrix               { return Transpose{m} }
func (m *basicBanded) Bandwidth() (kl, ku int) { return (*BandDense)(m).Bandwidth() }
func (m *basicBanded) TBand() Banded           { return TransposeBand{m} }

type basicSymBanded SymBandDense

var _ SymBanded = &basicSymBanded{}

func (m *basicSymBanded) At(r, c int) float64     { return (*SymBandDense)(m).At(r, c) }
func (m *basicSymBanded) Dims() (r, c int)        { return (*SymBandDense)(m).Dims() }
func (m *basicSymBanded) T() Matrix               { return m }
func (m *basicSymBanded) Bandwidth() (kl, ku int) { return (*SymBandDense)(m).Bandwidth() }
func (m *basicSymBanded) TBand() Banded           { return m }
func (m *basicSymBanded) SymmetricDim() int       { return (*SymBandDense)(m).SymmetricDim() }
func (m *basicSymBanded) SymBand() (n, k int)     { return (*SymBandDense)(m).SymBand() }

type basicTriBanded TriBandDense

var _ TriBanded = &basicTriBanded{}

func (m *basicTriBanded) At(r, c int) float64               { return (*TriBandDense)(m).At(r, c) }
func (m *basicTriBanded) Dims() (r, c int)                  { return (*TriBandDense)(m).Dims() }
func (m *basicTriBanded) T() Matrix                         { return Transpose{m} }
func (m *basicTriBanded) Triangle() (int, TriKind)          { return (*TriBandDense)(m).Triangle() }
func (m *basicTriBanded) TTri() Triangular                  { return TransposeTri{m} }
func (m *basicTriBanded) Bandwidth() (kl, ku int)           { return (*TriBandDense)(m).Bandwidth() }
func (m *basicTriBanded) TBand() Banded                     { return TransposeBand{m} }
func (m *basicTriBanded) TriBand() (n, k int, kind TriKind) { return (*TriBandDense)(m).TriBand() }
func (m *basicTriBanded) TTriBand() TriBanded               { return TransposeTriBand{m} }

type basicDiagonal DiagDense

var _ Diagonal = &basicDiagonal{}

func (m *basicDiagonal) At(r, c int) float64               { return (*DiagDense)(m).At(r, c) }
func (m *basicDiagonal) Dims() (r, c int)                  { return (*DiagDense)(m).Dims() }
func (m *basicDiagonal) T() Matrix                         { return Transpose{m} }
func (m *basicDiagonal) Diag() int                         { return (*DiagDense)(m).Diag() }
func (m *basicDiagonal) SymmetricDim() int                 { return (*DiagDense)(m).SymmetricDim() }
func (m *basicDiagonal) SymBand() (n, k int)               { return (*DiagDense)(m).SymBand() }
func (m *basicDiagonal) Bandwidth() (kl, ku int)           { return (*DiagDense)(m).Bandwidth() }
func (m *basicDiagonal) TBand() Banded                     { return TransposeBand{m} }
func (m *basicDiagonal) Triangle() (int, TriKind)          { return (*DiagDense)(m).Triangle() }
func (m *basicDiagonal) TTri() Triangular                  { return TransposeTri{m} }
func (m *basicDiagonal) TriBand() (n, k int, kind TriKind) { return (*DiagDense)(m).TriBand() }
func (m *basicDiagonal) TTriBand() TriBanded               { return TransposeTriBand{m} }
