// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"testing"

	"gonum.org/v1/gonum/lapack/testlapack"
)

var impl = Implementation{}

func TestDbdsqr(t *testing.T) {
	t.Parallel()
	testlapack.DbdsqrTest(t, impl)
}

func TestDcombssq(t *testing.T) {
	t.Parallel()
	testlapack.DcombssqTest(t, impl)
}

func TestDhseqr(t *testing.T) {
	t.Parallel()
	testlapack.DhseqrTest(t, impl)
}

func TestDgebak(t *testing.T) {
	t.Parallel()
	testlapack.DgebakTest(t, impl)
}

func TestDgebal(t *testing.T) {
	t.Parallel()
	testlapack.DgebalTest(t, impl)
}

func TestDgebd2(t *testing.T) {
	t.Parallel()
	testlapack.Dgebd2Test(t, impl)
}

func TestDgebrd(t *testing.T) {
	t.Parallel()
	testlapack.DgebrdTest(t, impl)
}

func TestDgecon(t *testing.T) {
	t.Parallel()
	testlapack.DgeconTest(t, impl)
}

func TestDgeev(t *testing.T) {
	t.Parallel()
	testlapack.DgeevTest(t, impl)
}

func TestDgehd2(t *testing.T) {
	t.Parallel()
	testlapack.Dgehd2Test(t, impl)
}

func TestDgehrd(t *testing.T) {
	t.Parallel()
	testlapack.DgehrdTest(t, impl)
}

func TestDgelqf(t *testing.T) {
	t.Parallel()
	testlapack.DgelqfTest(t, impl)
}

func TestDgelq2(t *testing.T) {
	t.Parallel()
	testlapack.Dgelq2Test(t, impl)
}

func TestDgeql2(t *testing.T) {
	t.Parallel()
	testlapack.Dgeql2Test(t, impl)
}

func TestDgels(t *testing.T) {
	t.Parallel()
	testlapack.DgelsTest(t, impl)
}

func TestDgerq2(t *testing.T) {
	t.Parallel()
	testlapack.Dgerq2Test(t, impl)
}

func TestDgeqp3(t *testing.T) {
	t.Parallel()
	testlapack.Dgeqp3Test(t, impl)
}

func TestDgeqr2(t *testing.T) {
	t.Parallel()
	testlapack.Dgeqr2Test(t, impl)
}

func TestDgeqrf(t *testing.T) {
	t.Parallel()
	testlapack.DgeqrfTest(t, impl)
}

func TestDgerqf(t *testing.T) {
	t.Parallel()
	testlapack.DgerqfTest(t, impl)
}

func TestDgesvd(t *testing.T) {
	t.Parallel()
	const tol = 1e-13
	testlapack.DgesvdTest(t, impl, tol)
}

func TestDgetri(t *testing.T) {
	t.Parallel()
	testlapack.DgetriTest(t, impl)
}

func TestDgetf2(t *testing.T) {
	t.Parallel()
	testlapack.Dgetf2Test(t, impl)
}

func TestDgetrf(t *testing.T) {
	t.Parallel()
	testlapack.DgetrfTest(t, impl)
}

func TestDgetrs(t *testing.T) {
	t.Parallel()
	testlapack.DgetrsTest(t, impl)
}

func TestDggsvd3(t *testing.T) {
	t.Parallel()
	testlapack.Dggsvd3Test(t, impl)
}

func TestDggsvp3(t *testing.T) {
	t.Parallel()
	testlapack.Dggsvp3Test(t, impl)
}

func TestDlabrd(t *testing.T) {
	t.Parallel()
	testlapack.DlabrdTest(t, impl)
}

func TestDlacn2(t *testing.T) {
	t.Parallel()
	testlapack.Dlacn2Test(t, impl)
}

func TestDlacpy(t *testing.T) {
	t.Parallel()
	testlapack.DlacpyTest(t, impl)
}

func TestDlae2(t *testing.T) {
	t.Parallel()
	testlapack.Dlae2Test(t, impl)
}

func TestDlaev2(t *testing.T) {
	t.Parallel()
	testlapack.Dlaev2Test(t, impl)
}

func TestDlaexc(t *testing.T) {
	t.Parallel()
	testlapack.DlaexcTest(t, impl)
}

func TestDlags2(t *testing.T) {
	t.Parallel()
	testlapack.Dlags2Test(t, impl)
}

func TestDlahqr(t *testing.T) {
	t.Parallel()
	testlapack.DlahqrTest(t, impl)
}

func TestDlahr2(t *testing.T) {
	t.Parallel()
	testlapack.Dlahr2Test(t, impl)
}

func TestDlaln2(t *testing.T) {
	t.Parallel()
	testlapack.Dlaln2Test(t, impl)
}

func TestDlange(t *testing.T) {
	t.Parallel()
	testlapack.DlangeTest(t, impl)
}

func TestDlapy2(t *testing.T) {
	t.Parallel()
	testlapack.Dlapy2Test(t, impl)
}

func TestDlapll(t *testing.T) {
	t.Parallel()
	testlapack.DlapllTest(t, impl)
}

func TestDlapmt(t *testing.T) {
	t.Parallel()
	testlapack.DlapmtTest(t, impl)
}

func TestDlas2(t *testing.T) {
	t.Parallel()
	testlapack.Dlas2Test(t, impl)
}

func TestDlascl(t *testing.T) {
	t.Parallel()
	testlapack.DlasclTest(t, impl)
}

func TestDlaset(t *testing.T) {
	t.Parallel()
	testlapack.DlasetTest(t, impl)
}

func TestDlasrt(t *testing.T) {
	t.Parallel()
	testlapack.DlasrtTest(t, impl)
}

func TestDlassq(t *testing.T) {
	t.Parallel()
	testlapack.DlassqTest(t, impl)
}

func TestDlaswp(t *testing.T) {
	t.Parallel()
	testlapack.DlaswpTest(t, impl)
}

func TestDlasy2(t *testing.T) {
	t.Parallel()
	testlapack.Dlasy2Test(t, impl)
}

func TestDlansb(t *testing.T) {
	t.Parallel()
	testlapack.DlansbTest(t, impl)
}

func TestDlanst(t *testing.T) {
	t.Parallel()
	testlapack.DlanstTest(t, impl)
}

func TestDlansy(t *testing.T) {
	t.Parallel()
	testlapack.DlansyTest(t, impl)
}

func TestDlantb(t *testing.T) {
	t.Parallel()
	testlapack.DlantbTest(t, impl)
}

func TestDlantr(t *testing.T) {
	t.Parallel()
	testlapack.DlantrTest(t, impl)
}

func TestDlanv2(t *testing.T) {
	t.Parallel()
	testlapack.Dlanv2Test(t, impl)
}

func TestDlaqr04(t *testing.T) {
	t.Parallel()
	testlapack.Dlaqr04Test(t, impl)
}

func TestDlaqp2(t *testing.T) {
	t.Parallel()
	testlapack.Dlaqp2Test(t, impl)
}

func TestDlaqps(t *testing.T) {
	t.Parallel()
	testlapack.DlaqpsTest(t, impl)
}

func TestDlaqr1(t *testing.T) {
	t.Parallel()
	testlapack.Dlaqr1Test(t, impl)
}

func TestDlaqr23(t *testing.T) {
	t.Parallel()
	testlapack.Dlaqr23Test(t, impl)
}

func TestDlaqr5(t *testing.T) {
	t.Parallel()
	testlapack.Dlaqr5Test(t, impl)
}

func TestDlarf(t *testing.T) {
	t.Parallel()
	testlapack.DlarfTest(t, impl)
}

func TestDlarfb(t *testing.T) {
	t.Parallel()
	testlapack.DlarfbTest(t, impl)
}

func TestDlarfg(t *testing.T) {
	t.Parallel()
	testlapack.DlarfgTest(t, impl)
}

func TestDlarft(t *testing.T) {
	t.Parallel()
	testlapack.DlarftTest(t, impl)
}

func TestDlarfx(t *testing.T) {
	t.Parallel()
	testlapack.DlarfxTest(t, impl)
}

func TestDlartg(t *testing.T) {
	t.Parallel()
	testlapack.DlartgTest(t, impl)
}

func TestDlasq1(t *testing.T) {
	t.Parallel()
	testlapack.Dlasq1Test(t, impl)
}

func TestDlasq2(t *testing.T) {
	t.Parallel()
	testlapack.Dlasq2Test(t, impl)
}

func TestDlasr(t *testing.T) {
	t.Parallel()
	testlapack.DlasrTest(t, impl)
}

func TestDlasv2(t *testing.T) {
	t.Parallel()
	testlapack.Dlasv2Test(t, impl)
}

func TestDlatbs(t *testing.T) {
	t.Parallel()
	testlapack.DlatbsTest(t, impl)
}

func TestDlatrd(t *testing.T) {
	t.Parallel()
	testlapack.DlatrdTest(t, impl)
}

func TestDlatrs(t *testing.T) {
	t.Parallel()
	testlapack.DlatrsTest(t, impl)
}

func TestDlauu2(t *testing.T) {
	t.Parallel()
	testlapack.Dlauu2Test(t, impl)
}

func TestDlauum(t *testing.T) {
	t.Parallel()
	testlapack.DlauumTest(t, impl)
}

func TestDorg2r(t *testing.T) {
	t.Parallel()
	testlapack.Dorg2rTest(t, impl)
}

func TestDorgbr(t *testing.T) {
	t.Parallel()
	testlapack.DorgbrTest(t, impl)
}

func TestDorghr(t *testing.T) {
	t.Parallel()
	testlapack.DorghrTest(t, impl)
}

func TestDorg2l(t *testing.T) {
	t.Parallel()
	testlapack.Dorg2lTest(t, impl)
}

func TestDorgl2(t *testing.T) {
	t.Parallel()
	testlapack.Dorgl2Test(t, impl)
}

func TestDorglq(t *testing.T) {
	t.Parallel()
	testlapack.DorglqTest(t, impl)
}

func TestDorgql(t *testing.T) {
	t.Parallel()
	testlapack.DorgqlTest(t, impl)
}

func TestDorgqr(t *testing.T) {
	t.Parallel()
	testlapack.DorgqrTest(t, impl)
}

func TestDorgtr(t *testing.T) {
	t.Parallel()
	testlapack.DorgtrTest(t, impl)
}

func TestDormbr(t *testing.T) {
	t.Parallel()
	testlapack.DormbrTest(t, impl)
}

func TestDormhr(t *testing.T) {
	t.Parallel()
	testlapack.DormhrTest(t, impl)
}

func TestDorml2(t *testing.T) {
	t.Parallel()
	testlapack.Dorml2Test(t, impl)
}

func TestDormlq(t *testing.T) {
	t.Parallel()
	testlapack.DormlqTest(t, impl)
}

func TestDormqr(t *testing.T) {
	t.Parallel()
	testlapack.DormqrTest(t, impl)
}

func TestDormr2(t *testing.T) {
	t.Parallel()
	testlapack.Dormr2Test(t, impl)
}

func TestDorm2r(t *testing.T) {
	t.Parallel()
	testlapack.Dorm2rTest(t, impl)
}

func TestDpbcon(t *testing.T) {
	t.Parallel()
	testlapack.DpbconTest(t, impl)
}

func TestDpbtf2(t *testing.T) {
	t.Parallel()
	testlapack.Dpbtf2Test(t, impl)
}

func TestDpbtrf(t *testing.T) {
	t.Parallel()
	testlapack.DpbtrfTest(t, impl)
}

func TestDpbtrs(t *testing.T) {
	t.Parallel()
	testlapack.DpbtrsTest(t, impl)
}

func TestDpocon(t *testing.T) {
	t.Parallel()
	testlapack.DpoconTest(t, impl)
}

func TestDpotf2(t *testing.T) {
	t.Parallel()
	testlapack.Dpotf2Test(t, impl)
}

func TestDpotrf(t *testing.T) {
	t.Parallel()
	testlapack.DpotrfTest(t, impl)
}

func TestDpotri(t *testing.T) {
	t.Parallel()
	testlapack.DpotriTest(t, impl)
}

func TestDpotrs(t *testing.T) {
	t.Parallel()
	testlapack.DpotrsTest(t, impl)
}

func TestDrscl(t *testing.T) {
	t.Parallel()
	testlapack.DrsclTest(t, impl)
}

func TestDsteqr(t *testing.T) {
	t.Parallel()
	testlapack.DsteqrTest(t, impl)
}

func TestDsterf(t *testing.T) {
	t.Parallel()
	testlapack.DsterfTest(t, impl)
}

func TestDsyev(t *testing.T) {
	t.Parallel()
	testlapack.DsyevTest(t, impl)
}

func TestDsytd2(t *testing.T) {
	t.Parallel()
	testlapack.Dsytd2Test(t, impl)
}

func TestDsytrd(t *testing.T) {
	t.Parallel()
	testlapack.DsytrdTest(t, impl)
}

func TestDtgsja(t *testing.T) {
	t.Parallel()
	testlapack.DtgsjaTest(t, impl)
}

func TestDtrcon(t *testing.T) {
	t.Parallel()
	testlapack.DtrconTest(t, impl)
}

func TestDtrevc3(t *testing.T) {
	t.Parallel()
	testlapack.Dtrevc3Test(t, impl)
}

func TestDtrexc(t *testing.T) {
	t.Parallel()
	testlapack.DtrexcTest(t, impl)
}

func TestDtrti2(t *testing.T) {
	t.Parallel()
	testlapack.Dtrti2Test(t, impl)
}

func TestDtrtri(t *testing.T) {
	t.Parallel()
	testlapack.DtrtriTest(t, impl)
}

func TestIladlc(t *testing.T) {
	t.Parallel()
	testlapack.IladlcTest(t, impl)
}

func TestIladlr(t *testing.T) {
	t.Parallel()
	testlapack.IladlrTest(t, impl)
}
