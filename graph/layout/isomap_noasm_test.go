// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build noasm || gccgo || safe
// +build noasm gccgo safe

package layout_test

// Change the testdata path for calculations done without assembly kernels.
func init() {
	tag = "_noasm"
}
