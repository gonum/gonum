// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build arm64

package layout

// Change the testdata path for calculations done on arm64.
func init() {
	arch = "_arm64"
}
