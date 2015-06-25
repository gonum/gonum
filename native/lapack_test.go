// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package native

import (
	"testing"

	"github.com/gonum/lapack/testlapack"
)

var impl = Implementation{}

func TestDlarf(t *testing.T) {
	testlapack.DlarfTest(t, impl)
}

func TestDpotf2(t *testing.T) {
	testlapack.Dpotf2Test(t, impl)
}

func TestDpotrf(t *testing.T) {
	testlapack.DpotrfTest(t, impl)
}

func TestIladlc(t *testing.T) {
	testlapack.IladlcTest(t, impl)
}

func TestIladlr(t *testing.T) {
	testlapack.IladlrTest(t, impl)
}
