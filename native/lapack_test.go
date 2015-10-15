// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package native

import (
	"testing"

	"github.com/gonum/lapack/testlapack"
)

var impl = Implementation{}

func TestDgecon(t *testing.T) {
	testlapack.DgeconTest(t, impl)
}

func TestDgelqf(t *testing.T) {
	testlapack.DgelqfTest(t, impl)
}

func TestDgelq2(t *testing.T) {
	testlapack.Dgelq2Test(t, impl)
}

func TestDgels(t *testing.T) {
	testlapack.DgelsTest(t, impl)
}

func TestDgeqr2(t *testing.T) {
	testlapack.Dgeqr2Test(t, impl)
}

func TestDgeqrf(t *testing.T) {
	testlapack.DgeqrfTest(t, impl)
}

func TestDgetri(t *testing.T) {
	testlapack.DgetriTest(t, impl)
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

func TestDlacpy(t *testing.T) {
	testlapack.DlacpyTest(t, impl)
}

func TestDlaev2(t *testing.T) {
	testlapack.Dlaev2Test(t, impl)
}

func TestDlange(t *testing.T) {
	testlapack.DlangeTest(t, impl)
}

func TestDlas2(t *testing.T) {
	testlapack.Dlas2Test(t, impl)
}

func TestDlansy(t *testing.T) {
	testlapack.DlansyTest(t, impl)
}

func TestDlantr(t *testing.T) {
	testlapack.DlantrTest(t, impl)
}

func TestDlarfb(t *testing.T) {
	testlapack.DlarfbTest(t, impl)
}

func TestDlarf(t *testing.T) {
	testlapack.DlarfTest(t, impl)
}

func TestDlarfg(t *testing.T) {
	testlapack.DlarfgTest(t, impl)
}

func TestDlarft(t *testing.T) {
	testlapack.DlarftTest(t, impl)
}

func TestDlartg(t *testing.T) {
	testlapack.DlartgTest(t, impl)
}

func TestDlasr(t *testing.T) {
	testlapack.DlasrTest(t, impl)
}

func TestDorml2(t *testing.T) {
	testlapack.Dorml2Test(t, impl)
}

func TestDormlq(t *testing.T) {
	testlapack.DormlqTest(t, impl)
}

func TestDormqr(t *testing.T) {
	testlapack.DormqrTest(t, impl)
}

func TestDorm2r(t *testing.T) {
	testlapack.Dorm2rTest(t, impl)
}

func TestDpocon(t *testing.T) {
	testlapack.DpoconTest(t, impl)
}

func TestDpotf2(t *testing.T) {
	testlapack.Dpotf2Test(t, impl)
}

func TestDpotrf(t *testing.T) {
	testlapack.DpotrfTest(t, impl)
}

func TestDrscl(t *testing.T) {
	testlapack.DrsclTest(t, impl)
}

func TestDtrcon(t *testing.T) {
	testlapack.DtrconTest(t, impl)
}

func TestDtrti2(t *testing.T) {
	testlapack.Dtrti2Test(t, impl)
}

func TestDtrtri(t *testing.T) {
	testlapack.DtrtriTest(t, impl)
}

func TestIladlc(t *testing.T) {
	testlapack.IladlcTest(t, impl)
}

func TestIladlr(t *testing.T) {
	testlapack.IladlrTest(t, impl)
}
