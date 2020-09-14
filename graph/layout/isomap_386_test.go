// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build 386

package layout

// Change the testdata path for calculations done on 386.
func init() {
	arch = "_386"
}
