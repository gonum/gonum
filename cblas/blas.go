// Do not manually edit this file. It was created by the genBlas.pl script from cblas.h.

// Copyright ©2014 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package cblas implements the blas interfaces.
package cblas

/*
#cgo CFLAGS: -g -O2
#include "cblas.h"
*/
import "C"

import (
	"unsafe"

	"github.com/gonum/blas"
)

// Type check assertions:
var (
	_ blas.Float32    = Blas{}
	_ blas.Float64    = Blas{}
	_ blas.Complex64  = Blas{}
	_ blas.Complex128 = Blas{}
)

// Type order is used to specify the matrix storage format. We still interact with
// an API that allows client calls to specify order, so this is here to document that fact.
type order int

const (
	rowMajor order = 101 + iota
)

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type Blas struct{}

// Special cases...

type srotmParams struct {
	flag float32
	h    [4]float32
}

type drotmParams struct {
	flag float64
	h    [4]float64
}

func (Blas) Srotg(a float32, b float32) (c float32, s float32, r float32, z float32) {
	C.cblas_srotg((*C.float)(&a), (*C.float)(&b), (*C.float)(&c), (*C.float)(&s))
	return c, s, a, b
}
func (Blas) Srotmg(d1 float32, d2 float32, b1 float32, b2 float32) (p blas.SrotmParams, rd1 float32, rd2 float32, rb1 float32) {
	var pi srotmParams
	C.cblas_srotmg((*C.float)(&d1), (*C.float)(&d2), (*C.float)(&b1), C.float(b2), (*C.float)(unsafe.Pointer(&pi)))
	return blas.SrotmParams{Flag: blas.Flag(pi.flag), H: pi.h}, d1, d2, b1
}
func (Blas) Srotm(n int, x []float32, incX int, y []float32, incY int, p blas.SrotmParams) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (n-1)*incX >= len(x) {
		panic("blas: index out of range")
	}
	if (n-1)*incY >= len(y) {
		panic("blas: index out of range")
	}
	if p.Flag < blas.Identity || p.Flag > blas.Diagonal {
		panic("blas: illegal blas.Flag value")
	}
	pi := srotmParams{
		flag: float32(p.Flag),
		h:    p.H,
	}
	C.cblas_srotm(C.int(n), (*C.float)(&x[0]), C.int(incX), (*C.float)(&y[0]), C.int(incY), (*C.float)(unsafe.Pointer(&pi)))
}
func (Blas) Drotg(a float64, b float64) (c float64, s float64, r float64, z float64) {
	C.cblas_drotg((*C.double)(&a), (*C.double)(&b), (*C.double)(&c), (*C.double)(&s))
	return c, s, a, b
}
func (Blas) Drotmg(d1 float64, d2 float64, b1 float64, b2 float64) (p blas.DrotmParams, rd1 float64, rd2 float64, rb1 float64) {
	var pi drotmParams
	C.cblas_drotmg((*C.double)(&d1), (*C.double)(&d2), (*C.double)(&b1), C.double(b2), (*C.double)(unsafe.Pointer(&pi)))
	return blas.DrotmParams{Flag: blas.Flag(pi.flag), H: pi.h}, d1, d2, b1
}
func (Blas) Drotm(n int, x []float64, incX int, y []float64, incY int, p blas.DrotmParams) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (n-1)*incX >= len(x) {
		panic("blas: index out of range")
	}
	if (n-1)*incY >= len(y) {
		panic("blas: index out of range")
	}
	if p.Flag < blas.Identity || p.Flag > blas.Diagonal {
		panic("blas: illegal blas.Flag value")
	}
	pi := drotmParams{
		flag: float64(p.Flag),
		h:    p.H,
	}
	C.cblas_drotm(C.int(n), (*C.double)(&x[0]), C.int(incX), (*C.double)(&y[0]), C.int(incY), (*C.double)(unsafe.Pointer(&pi)))
}
func (Blas) Cdotu(n int, x []complex64, incX int, y []complex64, incY int) (dotu complex64) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (n-1)*incX >= len(x) {
		panic("blas: index out of range")
	}
	if (n-1)*incY >= len(y) {
		panic("blas: index out of range")
	}
	C.cblas_cdotu_sub(C.int(n), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&y[0]), C.int(incY), unsafe.Pointer(&dotu))
	return dotu
}
func (Blas) Cdotc(n int, x []complex64, incX int, y []complex64, incY int) (dotc complex64) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (n-1)*incX >= len(x) {
		panic("blas: index out of range")
	}
	if (n-1)*incY >= len(y) {
		panic("blas: index out of range")
	}
	C.cblas_cdotc_sub(C.int(n), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&y[0]), C.int(incY), unsafe.Pointer(&dotc))
	return dotc
}
func (Blas) Zdotu(n int, x []complex128, incX int, y []complex128, incY int) (dotu complex128) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (n-1)*incX >= len(x) {
		panic("blas: index out of range")
	}
	if (n-1)*incY >= len(y) {
		panic("blas: index out of range")
	}
	C.cblas_zdotu_sub(C.int(n), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&y[0]), C.int(incY), unsafe.Pointer(&dotu))
	return dotu
}
func (Blas) Zdotc(n int, x []complex128, incX int, y []complex128, incY int) (dotc complex128) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (n-1)*incX >= len(x) {
		panic("blas: index out of range")
	}
	if (n-1)*incY >= len(y) {
		panic("blas: index out of range")
	}
	C.cblas_zdotc_sub(C.int(n), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&y[0]), C.int(incY), unsafe.Pointer(&dotc))
	return dotc
}

func (Blas) Sdsdot(n int, alpha float32, x []float32, incX int, y []float32, incY int) float32 {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	return float32(C.cblas_sdsdot(C.int(n), C.float(alpha), (*C.float)(&x[0]), C.int(incX), (*C.float)(&y[0]), C.int(incY)))
}
func (Blas) Dsdot(n int, x []float32, incX int, y []float32, incY int) float64 {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	return float64(C.cblas_dsdot(C.int(n), (*C.float)(&x[0]), C.int(incX), (*C.float)(&y[0]), C.int(incY)))
}
func (Blas) Sdot(n int, x []float32, incX int, y []float32, incY int) float32 {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	return float32(C.cblas_sdot(C.int(n), (*C.float)(&x[0]), C.int(incX), (*C.float)(&y[0]), C.int(incY)))
}
func (Blas) Ddot(n int, x []float64, incX int, y []float64, incY int) float64 {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	return float64(C.cblas_ddot(C.int(n), (*C.double)(&x[0]), C.int(incX), (*C.double)(&y[0]), C.int(incY)))
}
func (Blas) Snrm2(n int, x []float32, incX int) float32 {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incX < 0 {
		return 0
	}
	if incX > 0 && (n-1)*incX >= len(x) {
		panic("blas: x index out of range")
	}
	return float32(C.cblas_snrm2(C.int(n), (*C.float)(&x[0]), C.int(incX)))
}
func (Blas) Sasum(n int, x []float32, incX int) float32 {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incX < 0 {
		return 0
	}
	if incX > 0 && (n-1)*incX >= len(x) {
		panic("blas: x index out of range")
	}
	return float32(C.cblas_sasum(C.int(n), (*C.float)(&x[0]), C.int(incX)))
}
func (Blas) Dnrm2(n int, x []float64, incX int) float64 {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incX < 0 {
		return 0
	}
	if incX > 0 && (n-1)*incX >= len(x) {
		panic("blas: x index out of range")
	}
	return float64(C.cblas_dnrm2(C.int(n), (*C.double)(&x[0]), C.int(incX)))
}
func (Blas) Dasum(n int, x []float64, incX int) float64 {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incX < 0 {
		return 0
	}
	if incX > 0 && (n-1)*incX >= len(x) {
		panic("blas: x index out of range")
	}
	return float64(C.cblas_dasum(C.int(n), (*C.double)(&x[0]), C.int(incX)))
}
func (Blas) Scnrm2(n int, x []complex64, incX int) float32 {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incX < 0 {
		return 0
	}
	if incX > 0 && (n-1)*incX >= len(x) {
		panic("blas: x index out of range")
	}
	return float32(C.cblas_scnrm2(C.int(n), unsafe.Pointer(&x[0]), C.int(incX)))
}
func (Blas) Scasum(n int, x []complex64, incX int) float32 {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incX < 0 {
		return 0
	}
	if incX > 0 && (n-1)*incX >= len(x) {
		panic("blas: x index out of range")
	}
	return float32(C.cblas_scasum(C.int(n), unsafe.Pointer(&x[0]), C.int(incX)))
}
func (Blas) Dznrm2(n int, x []complex128, incX int) float64 {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incX < 0 {
		return 0
	}
	if incX > 0 && (n-1)*incX >= len(x) {
		panic("blas: x index out of range")
	}
	return float64(C.cblas_dznrm2(C.int(n), unsafe.Pointer(&x[0]), C.int(incX)))
}
func (Blas) Dzasum(n int, x []complex128, incX int) float64 {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incX < 0 {
		return 0
	}
	if incX > 0 && (n-1)*incX >= len(x) {
		panic("blas: x index out of range")
	}
	return float64(C.cblas_dzasum(C.int(n), unsafe.Pointer(&x[0]), C.int(incX)))
}
func (Blas) Isamax(n int, x []float32, incX int) int {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if n == 0 || incX < 0 {
		return -1
	}
	if incX > 0 && (n-1)*incX >= len(x) {
		panic("blas: x index out of range")
	}
	return int(C.cblas_isamax(C.int(n), (*C.float)(&x[0]), C.int(incX)))
}
func (Blas) Idamax(n int, x []float64, incX int) int {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if n == 0 || incX < 0 {
		return -1
	}
	if incX > 0 && (n-1)*incX >= len(x) {
		panic("blas: x index out of range")
	}
	return int(C.cblas_idamax(C.int(n), (*C.double)(&x[0]), C.int(incX)))
}
func (Blas) Icamax(n int, x []complex64, incX int) int {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if n == 0 || incX < 0 {
		return -1
	}
	if incX > 0 && (n-1)*incX >= len(x) {
		panic("blas: x index out of range")
	}
	return int(C.cblas_icamax(C.int(n), unsafe.Pointer(&x[0]), C.int(incX)))
}
func (Blas) Izamax(n int, x []complex128, incX int) int {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if n == 0 || incX < 0 {
		return -1
	}
	if incX > 0 && (n-1)*incX >= len(x) {
		panic("blas: x index out of range")
	}
	return int(C.cblas_izamax(C.int(n), unsafe.Pointer(&x[0]), C.int(incX)))
}
func (Blas) Sswap(n int, x []float32, incX int, y []float32, incY int) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	C.cblas_sswap(C.int(n), (*C.float)(&x[0]), C.int(incX), (*C.float)(&y[0]), C.int(incY))
}
func (Blas) Scopy(n int, x []float32, incX int, y []float32, incY int) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	C.cblas_scopy(C.int(n), (*C.float)(&x[0]), C.int(incX), (*C.float)(&y[0]), C.int(incY))
}
func (Blas) Saxpy(n int, alpha float32, x []float32, incX int, y []float32, incY int) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	C.cblas_saxpy(C.int(n), C.float(alpha), (*C.float)(&x[0]), C.int(incX), (*C.float)(&y[0]), C.int(incY))
}
func (Blas) Dswap(n int, x []float64, incX int, y []float64, incY int) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	C.cblas_dswap(C.int(n), (*C.double)(&x[0]), C.int(incX), (*C.double)(&y[0]), C.int(incY))
}
func (Blas) Dcopy(n int, x []float64, incX int, y []float64, incY int) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	C.cblas_dcopy(C.int(n), (*C.double)(&x[0]), C.int(incX), (*C.double)(&y[0]), C.int(incY))
}
func (Blas) Daxpy(n int, alpha float64, x []float64, incX int, y []float64, incY int) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	C.cblas_daxpy(C.int(n), C.double(alpha), (*C.double)(&x[0]), C.int(incX), (*C.double)(&y[0]), C.int(incY))
}
func (Blas) Cswap(n int, x []complex64, incX int, y []complex64, incY int) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	C.cblas_cswap(C.int(n), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&y[0]), C.int(incY))
}
func (Blas) Ccopy(n int, x []complex64, incX int, y []complex64, incY int) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	C.cblas_ccopy(C.int(n), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&y[0]), C.int(incY))
}
func (Blas) Caxpy(n int, alpha complex64, x []complex64, incX int, y []complex64, incY int) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	C.cblas_caxpy(C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&y[0]), C.int(incY))
}
func (Blas) Zswap(n int, x []complex128, incX int, y []complex128, incY int) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	C.cblas_zswap(C.int(n), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&y[0]), C.int(incY))
}
func (Blas) Zcopy(n int, x []complex128, incX int, y []complex128, incY int) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	C.cblas_zcopy(C.int(n), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&y[0]), C.int(incY))
}
func (Blas) Zaxpy(n int, alpha complex128, x []complex128, incX int, y []complex128, incY int) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	C.cblas_zaxpy(C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&y[0]), C.int(incY))
}
func (Blas) Srot(n int, x []float32, incX int, y []float32, incY int, c float32, s float32) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	C.cblas_srot(C.int(n), (*C.float)(&x[0]), C.int(incX), (*C.float)(&y[0]), C.int(incY), C.float(c), C.float(s))
}
func (Blas) Drot(n int, x []float64, incX int, y []float64, incY int, c float64, s float64) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	C.cblas_drot(C.int(n), (*C.double)(&x[0]), C.int(incX), (*C.double)(&y[0]), C.int(incY), C.double(c), C.double(s))
}
func (Blas) Sscal(n int, alpha float32, x []float32, incX int) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incX < 0 {
		return
	}
	if incX > 0 && (n-1)*incX >= len(x) {
		panic("blas: x index out of range")
	}
	C.cblas_sscal(C.int(n), C.float(alpha), (*C.float)(&x[0]), C.int(incX))
}
func (Blas) Dscal(n int, alpha float64, x []float64, incX int) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incX < 0 {
		return
	}
	if incX > 0 && (n-1)*incX >= len(x) {
		panic("blas: x index out of range")
	}
	C.cblas_dscal(C.int(n), C.double(alpha), (*C.double)(&x[0]), C.int(incX))
}
func (Blas) Cscal(n int, alpha complex64, x []complex64, incX int) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incX < 0 {
		return
	}
	if incX > 0 && (n-1)*incX >= len(x) {
		panic("blas: x index out of range")
	}
	C.cblas_cscal(C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&x[0]), C.int(incX))
}
func (Blas) Zscal(n int, alpha complex128, x []complex128, incX int) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incX < 0 {
		return
	}
	if incX > 0 && (n-1)*incX >= len(x) {
		panic("blas: x index out of range")
	}
	C.cblas_zscal(C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&x[0]), C.int(incX))
}
func (Blas) Csscal(n int, alpha float32, x []complex64, incX int) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incX < 0 {
		return
	}
	if incX > 0 && (n-1)*incX >= len(x) {
		panic("blas: x index out of range")
	}
	C.cblas_csscal(C.int(n), C.float(alpha), unsafe.Pointer(&x[0]), C.int(incX))
}
func (Blas) Zdscal(n int, alpha float64, x []complex128, incX int) {
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	C.cblas_zdscal(C.int(n), C.double(alpha), unsafe.Pointer(&x[0]), C.int(incX))
}
func (Blas) Sgemv(tA blas.Transpose, m int, n int, alpha float32, a []float32, lda int, x []float32, incX int, beta float32, y []float32, incY int) {
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	var lenX, lenY int
	if tA == blas.NoTrans {
		lenX, lenY = n, m
	} else {
		lenX, lenY = m, n
	}
	if (incX > 0 && (lenX-1)*incX >= len(x)) || (incX < 0 && (1-lenX)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (lenY-1)*incY >= len(y)) || (incY < 0 && (1-lenY)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(m-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_sgemv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_TRANSPOSE(tA), C.int(m), C.int(n), C.float(alpha), (*C.float)(&a[0]), C.int(lda), (*C.float)(&x[0]), C.int(incX), C.float(beta), (*C.float)(&y[0]), C.int(incY))
}
func (Blas) Sgbmv(tA blas.Transpose, m int, n int, kL int, kU int, alpha float32, a []float32, lda int, x []float32, incX int, beta float32, y []float32, incY int) {
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if kL < 0 {
		panic("blas: kL < 0")
	}
	if kU < 0 {
		panic("blas: kU < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	var lenX, lenY int
	if tA == blas.NoTrans {
		lenX, lenY = n, m
	} else {
		lenX, lenY = m, n
	}
	if (incX > 0 && (lenX-1)*incX >= len(x)) || (incX < 0 && (1-lenX)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (lenY-1)*incY >= len(y)) || (incY < 0 && (1-lenY)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(m-1)+kL+kU+1 > len(a) || lda < kL+kU+1 {
		panic("blas: index of a out of range")
	}
	C.cblas_sgbmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_TRANSPOSE(tA), C.int(m), C.int(n), C.int(kL), C.int(kU), C.float(alpha), (*C.float)(&a[0]), C.int(lda), (*C.float)(&x[0]), C.int(incX), C.float(beta), (*C.float)(&y[0]), C.int(incY))
}
func (Blas) Strmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, a []float32, lda int, x []float32, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if lda*(n-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_strmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), (*C.float)(&a[0]), C.int(lda), (*C.float)(&x[0]), C.int(incX))
}
func (Blas) Stbmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, k int, a []float32, lda int, x []float32, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if lda*(n-1)+k+1 > len(a) || lda < k+1 {
		panic("blas: index of a out of range")
	}
	C.cblas_stbmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), C.int(k), (*C.float)(&a[0]), C.int(lda), (*C.float)(&x[0]), C.int(incX))
}
func (Blas) Stpmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, ap []float32, x []float32, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if n*(n+1)/2 > len(ap) {
		panic("blas: index of ap out of range")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	C.cblas_stpmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), (*C.float)(&ap[0]), (*C.float)(&x[0]), C.int(incX))
}
func (Blas) Strsv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, a []float32, lda int, x []float32, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if lda*(n-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_strsv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), (*C.float)(&a[0]), C.int(lda), (*C.float)(&x[0]), C.int(incX))
}
func (Blas) Stbsv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, k int, a []float32, lda int, x []float32, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if lda*(n-1)+k+1 > len(a) || lda < k+1 {
		panic("blas: index of a out of range")
	}
	C.cblas_stbsv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), C.int(k), (*C.float)(&a[0]), C.int(lda), (*C.float)(&x[0]), C.int(incX))
}
func (Blas) Stpsv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, ap []float32, x []float32, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if n*(n+1)/2 > len(ap) {
		panic("blas: index of ap out of range")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	C.cblas_stpsv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), (*C.float)(&ap[0]), (*C.float)(&x[0]), C.int(incX))
}
func (Blas) Dgemv(tA blas.Transpose, m int, n int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int) {
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	var lenX, lenY int
	if tA == blas.NoTrans {
		lenX, lenY = n, m
	} else {
		lenX, lenY = m, n
	}
	if (incX > 0 && (lenX-1)*incX >= len(x)) || (incX < 0 && (1-lenX)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (lenY-1)*incY >= len(y)) || (incY < 0 && (1-lenY)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(m-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_dgemv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_TRANSPOSE(tA), C.int(m), C.int(n), C.double(alpha), (*C.double)(&a[0]), C.int(lda), (*C.double)(&x[0]), C.int(incX), C.double(beta), (*C.double)(&y[0]), C.int(incY))
}
func (Blas) Dgbmv(tA blas.Transpose, m int, n int, kL int, kU int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int) {
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if kL < 0 {
		panic("blas: kL < 0")
	}
	if kU < 0 {
		panic("blas: kU < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	var lenX, lenY int
	if tA == blas.NoTrans {
		lenX, lenY = n, m
	} else {
		lenX, lenY = m, n
	}
	if (incX > 0 && (lenX-1)*incX >= len(x)) || (incX < 0 && (1-lenX)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (lenY-1)*incY >= len(y)) || (incY < 0 && (1-lenY)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(m-1)+kL+kU+1 > len(a) || lda < kL+kU+1 {
		panic("blas: index of a out of range")
	}
	C.cblas_dgbmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_TRANSPOSE(tA), C.int(m), C.int(n), C.int(kL), C.int(kU), C.double(alpha), (*C.double)(&a[0]), C.int(lda), (*C.double)(&x[0]), C.int(incX), C.double(beta), (*C.double)(&y[0]), C.int(incY))
}
func (Blas) Dtrmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, a []float64, lda int, x []float64, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if lda*(n-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_dtrmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), (*C.double)(&a[0]), C.int(lda), (*C.double)(&x[0]), C.int(incX))
}
func (Blas) Dtbmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, k int, a []float64, lda int, x []float64, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if lda*(n-1)+k+1 > len(a) || lda < k+1 {
		panic("blas: index of a out of range")
	}
	C.cblas_dtbmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), C.int(k), (*C.double)(&a[0]), C.int(lda), (*C.double)(&x[0]), C.int(incX))
}
func (Blas) Dtpmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, ap []float64, x []float64, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if n*(n+1)/2 > len(ap) {
		panic("blas: index of ap out of range")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	C.cblas_dtpmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), (*C.double)(&ap[0]), (*C.double)(&x[0]), C.int(incX))
}
func (Blas) Dtrsv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, a []float64, lda int, x []float64, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if lda*(n-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_dtrsv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), (*C.double)(&a[0]), C.int(lda), (*C.double)(&x[0]), C.int(incX))
}
func (Blas) Dtbsv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, k int, a []float64, lda int, x []float64, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if lda*(n-1)+k+1 > len(a) || lda < k+1 {
		panic("blas: index of a out of range")
	}
	C.cblas_dtbsv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), C.int(k), (*C.double)(&a[0]), C.int(lda), (*C.double)(&x[0]), C.int(incX))
}
func (Blas) Dtpsv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, ap []float64, x []float64, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if n*(n+1)/2 > len(ap) {
		panic("blas: index of ap out of range")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	C.cblas_dtpsv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), (*C.double)(&ap[0]), (*C.double)(&x[0]), C.int(incX))
}
func (Blas) Cgemv(tA blas.Transpose, m int, n int, alpha complex64, a []complex64, lda int, x []complex64, incX int, beta complex64, y []complex64, incY int) {
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	var lenX, lenY int
	if tA == blas.NoTrans {
		lenX, lenY = n, m
	} else {
		lenX, lenY = m, n
	}
	if (incX > 0 && (lenX-1)*incX >= len(x)) || (incX < 0 && (1-lenX)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (lenY-1)*incY >= len(y)) || (incY < 0 && (1-lenY)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(m-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_cgemv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_TRANSPOSE(tA), C.int(m), C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&beta), unsafe.Pointer(&y[0]), C.int(incY))
}
func (Blas) Cgbmv(tA blas.Transpose, m int, n int, kL int, kU int, alpha complex64, a []complex64, lda int, x []complex64, incX int, beta complex64, y []complex64, incY int) {
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if kL < 0 {
		panic("blas: kL < 0")
	}
	if kU < 0 {
		panic("blas: kU < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	var lenX, lenY int
	if tA == blas.NoTrans {
		lenX, lenY = n, m
	} else {
		lenX, lenY = m, n
	}
	if (incX > 0 && (lenX-1)*incX >= len(x)) || (incX < 0 && (1-lenX)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (lenY-1)*incY >= len(y)) || (incY < 0 && (1-lenY)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(m-1)+kL+kU+1 > len(a) || lda < kL+kU+1 {
		panic("blas: index of a out of range")
	}
	C.cblas_cgbmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_TRANSPOSE(tA), C.int(m), C.int(n), C.int(kL), C.int(kU), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&beta), unsafe.Pointer(&y[0]), C.int(incY))
}
func (Blas) Ctrmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, a []complex64, lda int, x []complex64, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if lda*(n-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_ctrmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&x[0]), C.int(incX))
}
func (Blas) Ctbmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, k int, a []complex64, lda int, x []complex64, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if lda*(n-1)+k+1 > len(a) || lda < k+1 {
		panic("blas: index of a out of range")
	}
	C.cblas_ctbmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), C.int(k), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&x[0]), C.int(incX))
}
func (Blas) Ctpmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, ap []complex64, x []complex64, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if n*(n+1)/2 > len(ap) {
		panic("blas: index of ap out of range")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	C.cblas_ctpmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), unsafe.Pointer(&ap[0]), unsafe.Pointer(&x[0]), C.int(incX))
}
func (Blas) Ctrsv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, a []complex64, lda int, x []complex64, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if lda*(n-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_ctrsv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&x[0]), C.int(incX))
}
func (Blas) Ctbsv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, k int, a []complex64, lda int, x []complex64, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if lda*(n-1)+k+1 > len(a) || lda < k+1 {
		panic("blas: index of a out of range")
	}
	C.cblas_ctbsv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), C.int(k), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&x[0]), C.int(incX))
}
func (Blas) Ctpsv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, ap []complex64, x []complex64, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if n*(n+1)/2 > len(ap) {
		panic("blas: index of ap out of range")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	C.cblas_ctpsv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), unsafe.Pointer(&ap[0]), unsafe.Pointer(&x[0]), C.int(incX))
}
func (Blas) Zgemv(tA blas.Transpose, m int, n int, alpha complex128, a []complex128, lda int, x []complex128, incX int, beta complex128, y []complex128, incY int) {
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	var lenX, lenY int
	if tA == blas.NoTrans {
		lenX, lenY = n, m
	} else {
		lenX, lenY = m, n
	}
	if (incX > 0 && (lenX-1)*incX >= len(x)) || (incX < 0 && (1-lenX)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (lenY-1)*incY >= len(y)) || (incY < 0 && (1-lenY)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(m-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_zgemv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_TRANSPOSE(tA), C.int(m), C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&beta), unsafe.Pointer(&y[0]), C.int(incY))
}
func (Blas) Zgbmv(tA blas.Transpose, m int, n int, kL int, kU int, alpha complex128, a []complex128, lda int, x []complex128, incX int, beta complex128, y []complex128, incY int) {
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if kL < 0 {
		panic("blas: kL < 0")
	}
	if kU < 0 {
		panic("blas: kU < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	var lenX, lenY int
	if tA == blas.NoTrans {
		lenX, lenY = n, m
	} else {
		lenX, lenY = m, n
	}
	if (incX > 0 && (lenX-1)*incX >= len(x)) || (incX < 0 && (1-lenX)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (lenY-1)*incY >= len(y)) || (incY < 0 && (1-lenY)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(m-1)+kL+kU+1 > len(a) || lda < kL+kU+1 {
		panic("blas: index of a out of range")
	}
	C.cblas_zgbmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_TRANSPOSE(tA), C.int(m), C.int(n), C.int(kL), C.int(kU), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&beta), unsafe.Pointer(&y[0]), C.int(incY))
}
func (Blas) Ztrmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, a []complex128, lda int, x []complex128, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if lda*(n-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_ztrmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&x[0]), C.int(incX))
}
func (Blas) Ztbmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, k int, a []complex128, lda int, x []complex128, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if lda*(n-1)+k+1 > len(a) || lda < k+1 {
		panic("blas: index of a out of range")
	}
	C.cblas_ztbmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), C.int(k), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&x[0]), C.int(incX))
}
func (Blas) Ztpmv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, ap []complex128, x []complex128, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if n*(n+1)/2 > len(ap) {
		panic("blas: index of ap out of range")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	C.cblas_ztpmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), unsafe.Pointer(&ap[0]), unsafe.Pointer(&x[0]), C.int(incX))
}
func (Blas) Ztrsv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, a []complex128, lda int, x []complex128, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if lda*(n-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_ztrsv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&x[0]), C.int(incX))
}
func (Blas) Ztbsv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, k int, a []complex128, lda int, x []complex128, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if lda*(n-1)+k+1 > len(a) || lda < k+1 {
		panic("blas: index of a out of range")
	}
	C.cblas_ztbsv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), C.int(k), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&x[0]), C.int(incX))
}
func (Blas) Ztpsv(ul blas.Uplo, tA blas.Transpose, d blas.Diag, n int, ap []complex128, x []complex128, incX int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if n*(n+1)/2 > len(ap) {
		panic("blas: index of ap out of range")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	C.cblas_ztpsv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(n), unsafe.Pointer(&ap[0]), unsafe.Pointer(&x[0]), C.int(incX))
}
func (Blas) Ssymv(ul blas.Uplo, n int, alpha float32, a []float32, lda int, x []float32, incX int, beta float32, y []float32, incY int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(n-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_ssymv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), C.float(alpha), (*C.float)(&a[0]), C.int(lda), (*C.float)(&x[0]), C.int(incX), C.float(beta), (*C.float)(&y[0]), C.int(incY))
}
func (Blas) Ssbmv(ul blas.Uplo, n int, k int, alpha float32, a []float32, lda int, x []float32, incX int, beta float32, y []float32, incY int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(n-1)+k+1 > len(a) || lda < k+1 {
		panic("blas: index of a out of range")
	}
	C.cblas_ssbmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), C.int(k), C.float(alpha), (*C.float)(&a[0]), C.int(lda), (*C.float)(&x[0]), C.int(incX), C.float(beta), (*C.float)(&y[0]), C.int(incY))
}
func (Blas) Sspmv(ul blas.Uplo, n int, alpha float32, ap []float32, x []float32, incX int, beta float32, y []float32, incY int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if n*(n+1)/2 > len(ap) {
		panic("blas: index of ap out of range")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	C.cblas_sspmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), C.float(alpha), (*C.float)(&ap[0]), (*C.float)(&x[0]), C.int(incX), C.float(beta), (*C.float)(&y[0]), C.int(incY))
}
func (Blas) Sger(m int, n int, alpha float32, x []float32, incX int, y []float32, incY int, a []float32, lda int) {
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (m-1)*incX >= len(x)) || (incX < 0 && (1-m)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(m-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_sger(C.enum_CBLAS_ORDER(rowMajor), C.int(m), C.int(n), C.float(alpha), (*C.float)(&x[0]), C.int(incX), (*C.float)(&y[0]), C.int(incY), (*C.float)(&a[0]), C.int(lda))
}
func (Blas) Ssyr(ul blas.Uplo, n int, alpha float32, x []float32, incX int, a []float32, lda int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if lda*(n-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_ssyr(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), C.float(alpha), (*C.float)(&x[0]), C.int(incX), (*C.float)(&a[0]), C.int(lda))
}
func (Blas) Sspr(ul blas.Uplo, n int, alpha float32, x []float32, incX int, ap []float32) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if n*(n+1)/2 > len(ap) {
		panic("blas: index of ap out of range")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	C.cblas_sspr(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), C.float(alpha), (*C.float)(&x[0]), C.int(incX), (*C.float)(&ap[0]))
}
func (Blas) Ssyr2(ul blas.Uplo, n int, alpha float32, x []float32, incX int, y []float32, incY int, a []float32, lda int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(n-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_ssyr2(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), C.float(alpha), (*C.float)(&x[0]), C.int(incX), (*C.float)(&y[0]), C.int(incY), (*C.float)(&a[0]), C.int(lda))
}
func (Blas) Sspr2(ul blas.Uplo, n int, alpha float32, x []float32, incX int, y []float32, incY int, ap []float32) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if n*(n+1)/2 > len(ap) {
		panic("blas: index of ap out of range")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	C.cblas_sspr2(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), C.float(alpha), (*C.float)(&x[0]), C.int(incX), (*C.float)(&y[0]), C.int(incY), (*C.float)(&ap[0]))
}
func (Blas) Dsymv(ul blas.Uplo, n int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(n-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_dsymv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), C.double(alpha), (*C.double)(&a[0]), C.int(lda), (*C.double)(&x[0]), C.int(incX), C.double(beta), (*C.double)(&y[0]), C.int(incY))
}
func (Blas) Dsbmv(ul blas.Uplo, n int, k int, alpha float64, a []float64, lda int, x []float64, incX int, beta float64, y []float64, incY int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(n-1)+k+1 > len(a) || lda < k+1 {
		panic("blas: index of a out of range")
	}
	C.cblas_dsbmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), C.int(k), C.double(alpha), (*C.double)(&a[0]), C.int(lda), (*C.double)(&x[0]), C.int(incX), C.double(beta), (*C.double)(&y[0]), C.int(incY))
}
func (Blas) Dspmv(ul blas.Uplo, n int, alpha float64, ap []float64, x []float64, incX int, beta float64, y []float64, incY int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if n*(n+1)/2 > len(ap) {
		panic("blas: index of ap out of range")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	C.cblas_dspmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), C.double(alpha), (*C.double)(&ap[0]), (*C.double)(&x[0]), C.int(incX), C.double(beta), (*C.double)(&y[0]), C.int(incY))
}
func (Blas) Dger(m int, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64, lda int) {
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (m-1)*incX >= len(x)) || (incX < 0 && (1-m)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(m-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_dger(C.enum_CBLAS_ORDER(rowMajor), C.int(m), C.int(n), C.double(alpha), (*C.double)(&x[0]), C.int(incX), (*C.double)(&y[0]), C.int(incY), (*C.double)(&a[0]), C.int(lda))
}
func (Blas) Dsyr(ul blas.Uplo, n int, alpha float64, x []float64, incX int, a []float64, lda int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if lda*(n-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_dsyr(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), C.double(alpha), (*C.double)(&x[0]), C.int(incX), (*C.double)(&a[0]), C.int(lda))
}
func (Blas) Dspr(ul blas.Uplo, n int, alpha float64, x []float64, incX int, ap []float64) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if n*(n+1)/2 > len(ap) {
		panic("blas: index of ap out of range")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	C.cblas_dspr(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), C.double(alpha), (*C.double)(&x[0]), C.int(incX), (*C.double)(&ap[0]))
}
func (Blas) Dsyr2(ul blas.Uplo, n int, alpha float64, x []float64, incX int, y []float64, incY int, a []float64, lda int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(n-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_dsyr2(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), C.double(alpha), (*C.double)(&x[0]), C.int(incX), (*C.double)(&y[0]), C.int(incY), (*C.double)(&a[0]), C.int(lda))
}
func (Blas) Dspr2(ul blas.Uplo, n int, alpha float64, x []float64, incX int, y []float64, incY int, ap []float64) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if n*(n+1)/2 > len(ap) {
		panic("blas: index of ap out of range")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	C.cblas_dspr2(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), C.double(alpha), (*C.double)(&x[0]), C.int(incX), (*C.double)(&y[0]), C.int(incY), (*C.double)(&ap[0]))
}
func (Blas) Chemv(ul blas.Uplo, n int, alpha complex64, a []complex64, lda int, x []complex64, incX int, beta complex64, y []complex64, incY int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(n-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_chemv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&beta), unsafe.Pointer(&y[0]), C.int(incY))
}
func (Blas) Chbmv(ul blas.Uplo, n int, k int, alpha complex64, a []complex64, lda int, x []complex64, incX int, beta complex64, y []complex64, incY int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(n-1)+k+1 > len(a) || lda < k+1 {
		panic("blas: index of a out of range")
	}
	C.cblas_chbmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), C.int(k), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&beta), unsafe.Pointer(&y[0]), C.int(incY))
}
func (Blas) Chpmv(ul blas.Uplo, n int, alpha complex64, ap []complex64, x []complex64, incX int, beta complex64, y []complex64, incY int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if n*(n+1)/2 > len(ap) {
		panic("blas: index of ap out of range")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	C.cblas_chpmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&ap[0]), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&beta), unsafe.Pointer(&y[0]), C.int(incY))
}
func (Blas) Cgeru(m int, n int, alpha complex64, x []complex64, incX int, y []complex64, incY int, a []complex64, lda int) {
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (m-1)*incX >= len(x)) || (incX < 0 && (1-m)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(m-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_cgeru(C.enum_CBLAS_ORDER(rowMajor), C.int(m), C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&y[0]), C.int(incY), unsafe.Pointer(&a[0]), C.int(lda))
}
func (Blas) Cgerc(m int, n int, alpha complex64, x []complex64, incX int, y []complex64, incY int, a []complex64, lda int) {
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (m-1)*incX >= len(x)) || (incX < 0 && (1-m)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(m-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_cgerc(C.enum_CBLAS_ORDER(rowMajor), C.int(m), C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&y[0]), C.int(incY), unsafe.Pointer(&a[0]), C.int(lda))
}
func (Blas) Cher(ul blas.Uplo, n int, alpha float32, x []complex64, incX int, a []complex64, lda int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if lda*(n-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_cher(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), C.float(alpha), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&a[0]), C.int(lda))
}
func (Blas) Chpr(ul blas.Uplo, n int, alpha float32, x []complex64, incX int, ap []complex64) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if n*(n+1)/2 > len(ap) {
		panic("blas: index of ap out of range")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	C.cblas_chpr(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), C.float(alpha), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&ap[0]))
}
func (Blas) Cher2(ul blas.Uplo, n int, alpha complex64, x []complex64, incX int, y []complex64, incY int, a []complex64, lda int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(n-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_cher2(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&y[0]), C.int(incY), unsafe.Pointer(&a[0]), C.int(lda))
}
func (Blas) Chpr2(ul blas.Uplo, n int, alpha complex64, x []complex64, incX int, y []complex64, incY int, ap []complex64) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if n*(n+1)/2 > len(ap) {
		panic("blas: index of ap out of range")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	C.cblas_chpr2(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&y[0]), C.int(incY), unsafe.Pointer(&ap[0]))
}
func (Blas) Zhemv(ul blas.Uplo, n int, alpha complex128, a []complex128, lda int, x []complex128, incX int, beta complex128, y []complex128, incY int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(n-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_zhemv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&beta), unsafe.Pointer(&y[0]), C.int(incY))
}
func (Blas) Zhbmv(ul blas.Uplo, n int, k int, alpha complex128, a []complex128, lda int, x []complex128, incX int, beta complex128, y []complex128, incY int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(n-1)+k+1 > len(a) || lda < k+1 {
		panic("blas: index of a out of range")
	}
	C.cblas_zhbmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), C.int(k), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&beta), unsafe.Pointer(&y[0]), C.int(incY))
}
func (Blas) Zhpmv(ul blas.Uplo, n int, alpha complex128, ap []complex128, x []complex128, incX int, beta complex128, y []complex128, incY int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if n*(n+1)/2 > len(ap) {
		panic("blas: index of ap out of range")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	C.cblas_zhpmv(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&ap[0]), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&beta), unsafe.Pointer(&y[0]), C.int(incY))
}
func (Blas) Zgeru(m int, n int, alpha complex128, x []complex128, incX int, y []complex128, incY int, a []complex128, lda int) {
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (m-1)*incX >= len(x)) || (incX < 0 && (1-m)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(m-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_zgeru(C.enum_CBLAS_ORDER(rowMajor), C.int(m), C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&y[0]), C.int(incY), unsafe.Pointer(&a[0]), C.int(lda))
}
func (Blas) Zgerc(m int, n int, alpha complex128, x []complex128, incX int, y []complex128, incY int, a []complex128, lda int) {
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (m-1)*incX >= len(x)) || (incX < 0 && (1-m)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(m-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_zgerc(C.enum_CBLAS_ORDER(rowMajor), C.int(m), C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&y[0]), C.int(incY), unsafe.Pointer(&a[0]), C.int(lda))
}
func (Blas) Zher(ul blas.Uplo, n int, alpha float64, x []complex128, incX int, a []complex128, lda int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if lda*(n-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_zher(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), C.double(alpha), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&a[0]), C.int(lda))
}
func (Blas) Zhpr(ul blas.Uplo, n int, alpha float64, x []complex128, incX int, ap []complex128) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if n*(n+1)/2 > len(ap) {
		panic("blas: index of ap out of range")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	C.cblas_zhpr(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), C.double(alpha), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&ap[0]))
}
func (Blas) Zher2(ul blas.Uplo, n int, alpha complex128, x []complex128, incX int, y []complex128, incY int, a []complex128, lda int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	if lda*(n-1)+n > len(a) || lda < max(1, n) {
		panic("blas: index of a out of range")
	}
	C.cblas_zher2(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&y[0]), C.int(incY), unsafe.Pointer(&a[0]), C.int(lda))
}
func (Blas) Zhpr2(ul blas.Uplo, n int, alpha complex128, x []complex128, incX int, y []complex128, incY int, ap []complex128) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if n*(n+1)/2 > len(ap) {
		panic("blas: index of ap out of range")
	}
	if incX == 0 {
		panic("blas: zero x index increment")
	}
	if incY == 0 {
		panic("blas: zero y index increment")
	}
	if (incX > 0 && (n-1)*incX >= len(x)) || (incX < 0 && (1-n)*incX >= len(x)) {
		panic("blas: x index out of range")
	}
	if (incY > 0 && (n-1)*incY >= len(y)) || (incY < 0 && (1-n)*incY >= len(y)) {
		panic("blas: y index out of range")
	}
	C.cblas_zhpr2(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&x[0]), C.int(incX), unsafe.Pointer(&y[0]), C.int(incY), unsafe.Pointer(&ap[0]))
}
func (Blas) Sgemm(tA blas.Transpose, tB blas.Transpose, m int, n int, k int, alpha float32, a []float32, lda int, b []float32, ldb int, beta float32, c []float32, ldc int) {
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if tB != blas.NoTrans && tB != blas.Trans && tB != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	var rowA, colA, rowB, colB int
	if tA == blas.NoTrans {
		rowA, colA = m, k
	} else {
		rowA, colA = k, m
	}
	if tB == blas.NoTrans {
		rowB, colB = k, n
	} else {
		rowB, colB = n, k
	}
	if lda*(rowA-1)+colA > len(a) || lda < max(1, colA) {
		panic("blas: index of a out of range")
	}
	if ldb*(rowB-1)+colB > len(b) || ldb < max(1, colB) {
		panic("blas: index of b out of range")
	}
	if ldc*(m-1)+n > len(c) || ldc < max(1, n) {
		panic("blas: index of c out of range")
	}
	C.cblas_sgemm(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_TRANSPOSE(tB), C.int(m), C.int(n), C.int(k), C.float(alpha), (*C.float)(&a[0]), C.int(lda), (*C.float)(&b[0]), C.int(ldb), C.float(beta), (*C.float)(&c[0]), C.int(ldc))
}
func (Blas) Ssymm(s blas.Side, ul blas.Uplo, m int, n int, alpha float32, a []float32, lda int, b []float32, ldb int, beta float32, c []float32, ldc int) {
	if s != blas.Left && s != blas.Right {
		panic("blas: illegal side")
	}
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	var k int
	if s == blas.Left {
		k = m
	} else {
		k = n
	}
	if lda*(k-1)+k > len(a) || lda < max(1, k) {
		panic("blas: index of a out of range")
	}
	if ldb*(m-1)+n > len(b) || ldb < max(1, n) {
		panic("blas: index of b out of range")
	}
	if ldc*(m-1)+n > len(c) || ldc < max(1, n) {
		panic("blas: index of c out of range")
	}
	C.cblas_ssymm(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_SIDE(s), C.enum_CBLAS_UPLO(ul), C.int(m), C.int(n), C.float(alpha), (*C.float)(&a[0]), C.int(lda), (*C.float)(&b[0]), C.int(ldb), C.float(beta), (*C.float)(&c[0]), C.int(ldc))
}
func (Blas) Ssyrk(ul blas.Uplo, t blas.Transpose, n int, k int, alpha float32, a []float32, lda int, beta float32, c []float32, ldc int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if t != blas.NoTrans && t != blas.Trans && t != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	var row, col int
	if t == blas.NoTrans {
		row, col = n, k
	} else {
		row, col = k, n
	}
	if lda*(row-1)+col > len(a) || lda < max(1, col) {
		panic("blas: index of a out of range")
	}
	if ldc*(n-1)+n > len(c) || ldc < max(1, n) {
		panic("blas: index of c out of range")
	}
	C.cblas_ssyrk(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(t), C.int(n), C.int(k), C.float(alpha), (*C.float)(&a[0]), C.int(lda), C.float(beta), (*C.float)(&c[0]), C.int(ldc))
}
func (Blas) Ssyr2k(ul blas.Uplo, t blas.Transpose, n int, k int, alpha float32, a []float32, lda int, b []float32, ldb int, beta float32, c []float32, ldc int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if t != blas.NoTrans && t != blas.Trans && t != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	var row, col int
	if t == blas.NoTrans {
		row, col = n, k
	} else {
		row, col = k, n
	}
	if lda*(row-1)+col > len(a) || lda < max(1, col) {
		panic("blas: index of a out of range")
	}
	if ldb*(row-1)+col > len(b) || ldb < max(1, col) {
		panic("blas: index of b out of range")
	}
	if ldc*(n-1)+n > len(c) || ldc < max(1, n) {
		panic("blas: index of c out of range")
	}
	C.cblas_ssyr2k(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(t), C.int(n), C.int(k), C.float(alpha), (*C.float)(&a[0]), C.int(lda), (*C.float)(&b[0]), C.int(ldb), C.float(beta), (*C.float)(&c[0]), C.int(ldc))
}
func (Blas) Strmm(s blas.Side, ul blas.Uplo, tA blas.Transpose, d blas.Diag, m int, n int, alpha float32, a []float32, lda int, b []float32, ldb int) {
	if s != blas.Left && s != blas.Right {
		panic("blas: illegal side")
	}
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	var k int
	if s == blas.Left {
		k = m
	} else {
		k = n
	}
	if lda*(k-1)+k > len(a) || lda < max(1, k) {
		panic("blas: index of a out of range")
	}
	if ldb*(m-1)+n > len(b) || ldb < max(1, n) {
		panic("blas: index of b out of range")
	}
	C.cblas_strmm(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_SIDE(s), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(m), C.int(n), C.float(alpha), (*C.float)(&a[0]), C.int(lda), (*C.float)(&b[0]), C.int(ldb))
}
func (Blas) Strsm(s blas.Side, ul blas.Uplo, tA blas.Transpose, d blas.Diag, m int, n int, alpha float32, a []float32, lda int, b []float32, ldb int) {
	if s != blas.Left && s != blas.Right {
		panic("blas: illegal side")
	}
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	var k int
	if s == blas.Left {
		k = m
	} else {
		k = n
	}
	if lda*(k-1)+k > len(a) || lda < max(1, k) {
		panic("blas: index of a out of range")
	}
	if ldb*(m-1)+n > len(b) || ldb < max(1, n) {
		panic("blas: index of b out of range")
	}
	C.cblas_strsm(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_SIDE(s), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(m), C.int(n), C.float(alpha), (*C.float)(&a[0]), C.int(lda), (*C.float)(&b[0]), C.int(ldb))
}
func (Blas) Dgemm(tA blas.Transpose, tB blas.Transpose, m int, n int, k int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int) {
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if tB != blas.NoTrans && tB != blas.Trans && tB != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	var rowA, colA, rowB, colB int
	if tA == blas.NoTrans {
		rowA, colA = m, k
	} else {
		rowA, colA = k, m
	}
	if tB == blas.NoTrans {
		rowB, colB = k, n
	} else {
		rowB, colB = n, k
	}
	if lda*(rowA-1)+colA > len(a) || lda < max(1, colA) {
		panic("blas: index of a out of range")
	}
	if ldb*(rowB-1)+colB > len(b) || ldb < max(1, colB) {
		panic("blas: index of b out of range")
	}
	if ldc*(m-1)+n > len(c) || ldc < max(1, n) {
		panic("blas: index of c out of range")
	}
	C.cblas_dgemm(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_TRANSPOSE(tB), C.int(m), C.int(n), C.int(k), C.double(alpha), (*C.double)(&a[0]), C.int(lda), (*C.double)(&b[0]), C.int(ldb), C.double(beta), (*C.double)(&c[0]), C.int(ldc))
}
func (Blas) Dsymm(s blas.Side, ul blas.Uplo, m int, n int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int) {
	if s != blas.Left && s != blas.Right {
		panic("blas: illegal side")
	}
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	var k int
	if s == blas.Left {
		k = m
	} else {
		k = n
	}
	if lda*(k-1)+k > len(a) || lda < max(1, k) {
		panic("blas: index of a out of range")
	}
	if ldb*(m-1)+n > len(b) || ldb < max(1, n) {
		panic("blas: index of b out of range")
	}
	if ldc*(m-1)+n > len(c) || ldc < max(1, n) {
		panic("blas: index of c out of range")
	}
	C.cblas_dsymm(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_SIDE(s), C.enum_CBLAS_UPLO(ul), C.int(m), C.int(n), C.double(alpha), (*C.double)(&a[0]), C.int(lda), (*C.double)(&b[0]), C.int(ldb), C.double(beta), (*C.double)(&c[0]), C.int(ldc))
}
func (Blas) Dsyrk(ul blas.Uplo, t blas.Transpose, n int, k int, alpha float64, a []float64, lda int, beta float64, c []float64, ldc int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if t != blas.NoTrans && t != blas.Trans && t != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	var row, col int
	if t == blas.NoTrans {
		row, col = n, k
	} else {
		row, col = k, n
	}
	if lda*(row-1)+col > len(a) || lda < max(1, col) {
		panic("blas: index of a out of range")
	}
	if ldc*(n-1)+n > len(c) || ldc < max(1, n) {
		panic("blas: index of c out of range")
	}
	C.cblas_dsyrk(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(t), C.int(n), C.int(k), C.double(alpha), (*C.double)(&a[0]), C.int(lda), C.double(beta), (*C.double)(&c[0]), C.int(ldc))
}
func (Blas) Dsyr2k(ul blas.Uplo, t blas.Transpose, n int, k int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if t != blas.NoTrans && t != blas.Trans && t != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	var row, col int
	if t == blas.NoTrans {
		row, col = n, k
	} else {
		row, col = k, n
	}
	if lda*(row-1)+col > len(a) || lda < max(1, col) {
		panic("blas: index of a out of range")
	}
	if ldb*(row-1)+col > len(b) || ldb < max(1, col) {
		panic("blas: index of b out of range")
	}
	if ldc*(n-1)+n > len(c) || ldc < max(1, n) {
		panic("blas: index of c out of range")
	}
	C.cblas_dsyr2k(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(t), C.int(n), C.int(k), C.double(alpha), (*C.double)(&a[0]), C.int(lda), (*C.double)(&b[0]), C.int(ldb), C.double(beta), (*C.double)(&c[0]), C.int(ldc))
}
func (Blas) Dtrmm(s blas.Side, ul blas.Uplo, tA blas.Transpose, d blas.Diag, m int, n int, alpha float64, a []float64, lda int, b []float64, ldb int) {
	if s != blas.Left && s != blas.Right {
		panic("blas: illegal side")
	}
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	var k int
	if s == blas.Left {
		k = m
	} else {
		k = n
	}
	if lda*(k-1)+k > len(a) || lda < max(1, k) {
		panic("blas: index of a out of range")
	}
	if ldb*(m-1)+n > len(b) || ldb < max(1, n) {
		panic("blas: index of b out of range")
	}
	C.cblas_dtrmm(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_SIDE(s), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(m), C.int(n), C.double(alpha), (*C.double)(&a[0]), C.int(lda), (*C.double)(&b[0]), C.int(ldb))
}
func (Blas) Dtrsm(s blas.Side, ul blas.Uplo, tA blas.Transpose, d blas.Diag, m int, n int, alpha float64, a []float64, lda int, b []float64, ldb int) {
	if s != blas.Left && s != blas.Right {
		panic("blas: illegal side")
	}
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	var k int
	if s == blas.Left {
		k = m
	} else {
		k = n
	}
	if lda*(k-1)+k > len(a) || lda < max(1, k) {
		panic("blas: index of a out of range")
	}
	if ldb*(m-1)+n > len(b) || ldb < max(1, n) {
		panic("blas: index of b out of range")
	}
	C.cblas_dtrsm(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_SIDE(s), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(m), C.int(n), C.double(alpha), (*C.double)(&a[0]), C.int(lda), (*C.double)(&b[0]), C.int(ldb))
}
func (Blas) Cgemm(tA blas.Transpose, tB blas.Transpose, m int, n int, k int, alpha complex64, a []complex64, lda int, b []complex64, ldb int, beta complex64, c []complex64, ldc int) {
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if tB != blas.NoTrans && tB != blas.Trans && tB != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	var rowA, colA, rowB, colB int
	if tA == blas.NoTrans {
		rowA, colA = m, k
	} else {
		rowA, colA = k, m
	}
	if tB == blas.NoTrans {
		rowB, colB = k, n
	} else {
		rowB, colB = n, k
	}
	if lda*(rowA-1)+colA > len(a) || lda < max(1, colA) {
		panic("blas: index of a out of range")
	}
	if ldb*(rowB-1)+colB > len(b) || ldb < max(1, colB) {
		panic("blas: index of b out of range")
	}
	if ldc*(m-1)+n > len(c) || ldc < max(1, n) {
		panic("blas: index of c out of range")
	}
	C.cblas_cgemm(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_TRANSPOSE(tB), C.int(m), C.int(n), C.int(k), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&b[0]), C.int(ldb), unsafe.Pointer(&beta), unsafe.Pointer(&c[0]), C.int(ldc))
}
func (Blas) Csymm(s blas.Side, ul blas.Uplo, m int, n int, alpha complex64, a []complex64, lda int, b []complex64, ldb int, beta complex64, c []complex64, ldc int) {
	if s != blas.Left && s != blas.Right {
		panic("blas: illegal side")
	}
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	var k int
	if s == blas.Left {
		k = m
	} else {
		k = n
	}
	if lda*(k-1)+k > len(a) || lda < max(1, k) {
		panic("blas: index of a out of range")
	}
	if ldb*(m-1)+n > len(b) || ldb < max(1, n) {
		panic("blas: index of b out of range")
	}
	if ldc*(m-1)+n > len(c) || ldc < max(1, n) {
		panic("blas: index of c out of range")
	}
	C.cblas_csymm(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_SIDE(s), C.enum_CBLAS_UPLO(ul), C.int(m), C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&b[0]), C.int(ldb), unsafe.Pointer(&beta), unsafe.Pointer(&c[0]), C.int(ldc))
}
func (Blas) Csyrk(ul blas.Uplo, t blas.Transpose, n int, k int, alpha complex64, a []complex64, lda int, beta complex64, c []complex64, ldc int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if t != blas.NoTrans && t != blas.Trans {
		panic("blas: illegal transpose")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	var row, col int
	if t == blas.NoTrans {
		row, col = n, k
	} else {
		row, col = k, n
	}
	if lda*(row-1)+col > len(a) || lda < max(1, col) {
		panic("blas: index of a out of range")
	}
	if ldc*(n-1)+n > len(c) || ldc < max(1, n) {
		panic("blas: index of c out of range")
	}
	C.cblas_csyrk(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(t), C.int(n), C.int(k), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&beta), unsafe.Pointer(&c[0]), C.int(ldc))
}
func (Blas) Csyr2k(ul blas.Uplo, t blas.Transpose, n int, k int, alpha complex64, a []complex64, lda int, b []complex64, ldb int, beta complex64, c []complex64, ldc int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if t != blas.NoTrans && t != blas.Trans {
		panic("blas: illegal transpose")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	var row, col int
	if t == blas.NoTrans {
		row, col = n, k
	} else {
		row, col = k, n
	}
	if lda*(row-1)+col > len(a) || lda < max(1, col) {
		panic("blas: index of a out of range")
	}
	if ldb*(row-1)+col > len(b) || ldb < max(1, col) {
		panic("blas: index of b out of range")
	}
	if ldc*(n-1)+n > len(c) || ldc < max(1, n) {
		panic("blas: index of c out of range")
	}
	C.cblas_csyr2k(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(t), C.int(n), C.int(k), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&b[0]), C.int(ldb), unsafe.Pointer(&beta), unsafe.Pointer(&c[0]), C.int(ldc))
}
func (Blas) Ctrmm(s blas.Side, ul blas.Uplo, tA blas.Transpose, d blas.Diag, m int, n int, alpha complex64, a []complex64, lda int, b []complex64, ldb int) {
	if s != blas.Left && s != blas.Right {
		panic("blas: illegal side")
	}
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	var k int
	if s == blas.Left {
		k = m
	} else {
		k = n
	}
	if lda*(k-1)+k > len(a) || lda < max(1, k) {
		panic("blas: index of a out of range")
	}
	if ldb*(m-1)+n > len(b) || ldb < max(1, n) {
		panic("blas: index of b out of range")
	}
	C.cblas_ctrmm(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_SIDE(s), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(m), C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&b[0]), C.int(ldb))
}
func (Blas) Ctrsm(s blas.Side, ul blas.Uplo, tA blas.Transpose, d blas.Diag, m int, n int, alpha complex64, a []complex64, lda int, b []complex64, ldb int) {
	if s != blas.Left && s != blas.Right {
		panic("blas: illegal side")
	}
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	var k int
	if s == blas.Left {
		k = m
	} else {
		k = n
	}
	if lda*(k-1)+k > len(a) || lda < max(1, k) {
		panic("blas: index of a out of range")
	}
	if ldb*(m-1)+n > len(b) || ldb < max(1, n) {
		panic("blas: index of b out of range")
	}
	C.cblas_ctrsm(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_SIDE(s), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(m), C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&b[0]), C.int(ldb))
}
func (Blas) Zgemm(tA blas.Transpose, tB blas.Transpose, m int, n int, k int, alpha complex128, a []complex128, lda int, b []complex128, ldb int, beta complex128, c []complex128, ldc int) {
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if tB != blas.NoTrans && tB != blas.Trans && tB != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	var rowA, colA, rowB, colB int
	if tA == blas.NoTrans {
		rowA, colA = m, k
	} else {
		rowA, colA = k, m
	}
	if tB == blas.NoTrans {
		rowB, colB = k, n
	} else {
		rowB, colB = n, k
	}
	if lda*(rowA-1)+colA > len(a) || lda < max(1, colA) {
		panic("blas: index of a out of range")
	}
	if ldb*(rowB-1)+colB > len(b) || ldb < max(1, colB) {
		panic("blas: index of b out of range")
	}
	if ldc*(m-1)+n > len(c) || ldc < max(1, n) {
		panic("blas: index of c out of range")
	}
	C.cblas_zgemm(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_TRANSPOSE(tB), C.int(m), C.int(n), C.int(k), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&b[0]), C.int(ldb), unsafe.Pointer(&beta), unsafe.Pointer(&c[0]), C.int(ldc))
}
func (Blas) Zsymm(s blas.Side, ul blas.Uplo, m int, n int, alpha complex128, a []complex128, lda int, b []complex128, ldb int, beta complex128, c []complex128, ldc int) {
	if s != blas.Left && s != blas.Right {
		panic("blas: illegal side")
	}
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	var k int
	if s == blas.Left {
		k = m
	} else {
		k = n
	}
	if lda*(k-1)+k > len(a) || lda < max(1, k) {
		panic("blas: index of a out of range")
	}
	if ldb*(m-1)+n > len(b) || ldb < max(1, n) {
		panic("blas: index of b out of range")
	}
	if ldc*(m-1)+n > len(c) || ldc < max(1, n) {
		panic("blas: index of c out of range")
	}
	C.cblas_zsymm(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_SIDE(s), C.enum_CBLAS_UPLO(ul), C.int(m), C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&b[0]), C.int(ldb), unsafe.Pointer(&beta), unsafe.Pointer(&c[0]), C.int(ldc))
}
func (Blas) Zsyrk(ul blas.Uplo, t blas.Transpose, n int, k int, alpha complex128, a []complex128, lda int, beta complex128, c []complex128, ldc int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if t != blas.NoTrans && t != blas.Trans {
		panic("blas: illegal transpose")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	var row, col int
	if t == blas.NoTrans {
		row, col = n, k
	} else {
		row, col = k, n
	}
	if lda*(row-1)+col > len(a) || lda < max(1, col) {
		panic("blas: index of a out of range")
	}
	if ldc*(n-1)+n > len(c) || ldc < max(1, n) {
		panic("blas: index of c out of range")
	}
	C.cblas_zsyrk(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(t), C.int(n), C.int(k), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&beta), unsafe.Pointer(&c[0]), C.int(ldc))
}
func (Blas) Zsyr2k(ul blas.Uplo, t blas.Transpose, n int, k int, alpha complex128, a []complex128, lda int, b []complex128, ldb int, beta complex128, c []complex128, ldc int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if t != blas.NoTrans && t != blas.Trans {
		panic("blas: illegal transpose")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	var row, col int
	if t == blas.NoTrans {
		row, col = n, k
	} else {
		row, col = k, n
	}
	if lda*(row-1)+col > len(a) || lda < max(1, col) {
		panic("blas: index of a out of range")
	}
	if ldb*(row-1)+col > len(b) || ldb < max(1, col) {
		panic("blas: index of b out of range")
	}
	if ldc*(n-1)+n > len(c) || ldc < max(1, n) {
		panic("blas: index of c out of range")
	}
	C.cblas_zsyr2k(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(t), C.int(n), C.int(k), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&b[0]), C.int(ldb), unsafe.Pointer(&beta), unsafe.Pointer(&c[0]), C.int(ldc))
}
func (Blas) Ztrmm(s blas.Side, ul blas.Uplo, tA blas.Transpose, d blas.Diag, m int, n int, alpha complex128, a []complex128, lda int, b []complex128, ldb int) {
	if s != blas.Left && s != blas.Right {
		panic("blas: illegal side")
	}
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	var k int
	if s == blas.Left {
		k = m
	} else {
		k = n
	}
	if lda*(k-1)+k > len(a) || lda < max(1, k) {
		panic("blas: index of a out of range")
	}
	if ldb*(m-1)+n > len(b) || ldb < max(1, n) {
		panic("blas: index of b out of range")
	}
	C.cblas_ztrmm(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_SIDE(s), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(m), C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&b[0]), C.int(ldb))
}
func (Blas) Ztrsm(s blas.Side, ul blas.Uplo, tA blas.Transpose, d blas.Diag, m int, n int, alpha complex128, a []complex128, lda int, b []complex128, ldb int) {
	if s != blas.Left && s != blas.Right {
		panic("blas: illegal side")
	}
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if tA != blas.NoTrans && tA != blas.Trans && tA != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if d != blas.NonUnit && d != blas.Unit {
		panic("blas: illegal diagonal")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	var k int
	if s == blas.Left {
		k = m
	} else {
		k = n
	}
	if lda*(k-1)+k > len(a) || lda < max(1, k) {
		panic("blas: index of a out of range")
	}
	if ldb*(m-1)+n > len(b) || ldb < max(1, n) {
		panic("blas: index of b out of range")
	}
	C.cblas_ztrsm(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_SIDE(s), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(tA), C.enum_CBLAS_DIAG(d), C.int(m), C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&b[0]), C.int(ldb))
}
func (Blas) Chemm(s blas.Side, ul blas.Uplo, m int, n int, alpha complex64, a []complex64, lda int, b []complex64, ldb int, beta complex64, c []complex64, ldc int) {
	if s != blas.Left && s != blas.Right {
		panic("blas: illegal side")
	}
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	var k int
	if s == blas.Left {
		k = m
	} else {
		k = n
	}
	if lda*(k-1)+k > len(a) || lda < max(1, k) {
		panic("blas: index of a out of range")
	}
	if ldb*(m-1)+n > len(b) || ldb < max(1, n) {
		panic("blas: index of b out of range")
	}
	if ldc*(m-1)+n > len(c) || ldc < max(1, n) {
		panic("blas: index of c out of range")
	}
	C.cblas_chemm(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_SIDE(s), C.enum_CBLAS_UPLO(ul), C.int(m), C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&b[0]), C.int(ldb), unsafe.Pointer(&beta), unsafe.Pointer(&c[0]), C.int(ldc))
}
func (Blas) Cherk(ul blas.Uplo, t blas.Transpose, n int, k int, alpha float32, a []complex64, lda int, beta float32, c []complex64, ldc int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if t != blas.NoTrans && t != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	var row, col int
	if t == blas.NoTrans {
		row, col = n, k
	} else {
		row, col = k, n
	}
	if lda*(row-1)+col > len(a) || lda < max(1, col) {
		panic("blas: index of a out of range")
	}
	if ldc*(n-1)+n > len(c) || ldc < max(1, n) {
		panic("blas: index of c out of range")
	}
	C.cblas_cherk(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(t), C.int(n), C.int(k), C.float(alpha), unsafe.Pointer(&a[0]), C.int(lda), C.float(beta), unsafe.Pointer(&c[0]), C.int(ldc))
}
func (Blas) Cher2k(ul blas.Uplo, t blas.Transpose, n int, k int, alpha complex64, a []complex64, lda int, b []complex64, ldb int, beta float32, c []complex64, ldc int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if t != blas.NoTrans && t != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	var row, col int
	if t == blas.NoTrans {
		row, col = n, k
	} else {
		row, col = k, n
	}
	if lda*(row-1)+col > len(a) || lda < max(1, col) {
		panic("blas: index of a out of range")
	}
	if ldb*(row-1)+col > len(b) || ldb < max(1, col) {
		panic("blas: index of b out of range")
	}
	if ldc*(n-1)+n > len(c) || ldc < max(1, n) {
		panic("blas: index of c out of range")
	}
	C.cblas_cher2k(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(t), C.int(n), C.int(k), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&b[0]), C.int(ldb), C.float(beta), unsafe.Pointer(&c[0]), C.int(ldc))
}
func (Blas) Zhemm(s blas.Side, ul blas.Uplo, m int, n int, alpha complex128, a []complex128, lda int, b []complex128, ldb int, beta complex128, c []complex128, ldc int) {
	if s != blas.Left && s != blas.Right {
		panic("blas: illegal side")
	}
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if m < 0 {
		panic("blas: m < 0")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	var k int
	if s == blas.Left {
		k = m
	} else {
		k = n
	}
	if lda*(k-1)+k > len(a) || lda < max(1, k) {
		panic("blas: index of a out of range")
	}
	if ldb*(m-1)+n > len(b) || ldb < max(1, n) {
		panic("blas: index of b out of range")
	}
	if ldc*(m-1)+n > len(c) || ldc < max(1, n) {
		panic("blas: index of c out of range")
	}
	C.cblas_zhemm(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_SIDE(s), C.enum_CBLAS_UPLO(ul), C.int(m), C.int(n), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&b[0]), C.int(ldb), unsafe.Pointer(&beta), unsafe.Pointer(&c[0]), C.int(ldc))
}
func (Blas) Zherk(ul blas.Uplo, t blas.Transpose, n int, k int, alpha float64, a []complex128, lda int, beta float64, c []complex128, ldc int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if t != blas.NoTrans && t != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	var row, col int
	if t == blas.NoTrans {
		row, col = n, k
	} else {
		row, col = k, n
	}
	if lda*(row-1)+col > len(a) || lda < max(1, col) {
		panic("blas: index of a out of range")
	}
	if ldc*(n-1)+n > len(c) || ldc < max(1, n) {
		panic("blas: index of c out of range")
	}
	C.cblas_zherk(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(t), C.int(n), C.int(k), C.double(alpha), unsafe.Pointer(&a[0]), C.int(lda), C.double(beta), unsafe.Pointer(&c[0]), C.int(ldc))
}
func (Blas) Zher2k(ul blas.Uplo, t blas.Transpose, n int, k int, alpha complex128, a []complex128, lda int, b []complex128, ldb int, beta float64, c []complex128, ldc int) {
	if ul != blas.Upper && ul != blas.Lower {
		panic("blas: illegal triangle")
	}
	if t != blas.NoTrans && t != blas.ConjTrans {
		panic("blas: illegal transpose")
	}
	if n < 0 {
		panic("blas: n < 0")
	}
	if k < 0 {
		panic("blas: k < 0")
	}
	var row, col int
	if t == blas.NoTrans {
		row, col = n, k
	} else {
		row, col = k, n
	}
	if lda*(row-1)+col > len(a) || lda < max(1, col) {
		panic("blas: index of a out of range")
	}
	if ldb*(row-1)+col > len(b) || ldb < max(1, col) {
		panic("blas: index of b out of range")
	}
	if ldc*(n-1)+n > len(c) || ldc < max(1, n) {
		panic("blas: index of c out of range")
	}
	C.cblas_zher2k(C.enum_CBLAS_ORDER(rowMajor), C.enum_CBLAS_UPLO(ul), C.enum_CBLAS_TRANSPOSE(t), C.int(n), C.int(k), unsafe.Pointer(&alpha), unsafe.Pointer(&a[0]), C.int(lda), unsafe.Pointer(&b[0]), C.int(ldb), C.double(beta), unsafe.Pointer(&c[0]), C.int(ldc))
}
