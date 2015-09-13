// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cgo

import (
	"testing"

	"github.com/gonum/blas"
	"github.com/gonum/lapack/testlapack"
)

var impl = Implementation{}

func TestDlange(t *testing.T) {
	testlapack.DlangeTest(t, impl)
}

func TestDlantr(t *testing.T) {
	testlapack.DlantrTest(t, impl)
}

func TestDpotrf(t *testing.T) {
	testlapack.DpotrfTest(t, impl)
}

func TestDgecon(t *testing.T) {
	testlapack.DgeconTest(t, impl)
}

func TestDgelq2(t *testing.T) {
	testlapack.Dgelq2Test(t, impl)
}

func TestDgels(t *testing.T) {
	testlapack.DgelsTest(t, impl)
}

func TestDgelqf(t *testing.T) {
	testlapack.DgelqfTest(t, impl)
}

func TestDgeqr2(t *testing.T) {
	testlapack.Dgeqr2Test(t, impl)
}

func TestDgeqrf(t *testing.T) {
	testlapack.DgeqrfTest(t, impl)
}

func TestDgetf2(t *testing.T) {
	testlapack.Dgetf2Test(t, impl)
}

func TestDgetrf(t *testing.T) {
	testlapack.DgetrfTest(t, impl)
}

func TestDgetrs(t *testing.T) {
	testlapack.DgetrsTest(t, impl)
}

// blockedTranslate transforms some blocked C calls to be the unblocked algorithms
// for testing, as several of the unblocked algorithms are not defined by the C
// interface.
type blockedTranslate struct {
	Implementation
}

func (d blockedTranslate) Dorm2r(side blas.Side, trans blas.Transpose, m, n, k int, a []float64, lda int, tau, c []float64, ldc int, work []float64) {
	impl.Dormqr(side, trans, m, n, k, a, lda, tau, c, ldc, work, len(work))
}

func (d blockedTranslate) Dorml2(side blas.Side, trans blas.Transpose, m, n, k int, a []float64, lda int, tau, c []float64, ldc int, work []float64) {
	impl.Dormlq(side, trans, m, n, k, a, lda, tau, c, ldc, work, len(work))
}

func TestDormqr(t *testing.T) {
	testlapack.Dorm2rTest(t, blockedTranslate{impl})
}

/*
// Test disabled because of bug in c interface. Leaving stub for easy reproducer.
func TestDormlq(t *testing.T) {
	testlapack.Dorml2Test(t, blockedTranslate{impl})
}
*/

func TestDpocon(t *testing.T) {
	testlapack.DpoconTest(t, impl)
}

func TestDtrcon(t *testing.T) {
	testlapack.DtrconTest(t, impl)
}

func TestDtrtri(t *testing.T) {
	testlapack.DtrtriTest(t, impl)
}
