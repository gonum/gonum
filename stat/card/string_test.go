// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !safe,!appengine

package card

import "unsafe"

func tmpString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
