// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.12

package gonum

import (
	"fmt"
	"runtime/debug"
)

const root = "gonum.org/v1/gonum"

// Version returns the version of Gonum and its checksum. The returned
// values are only valid in binaries built with module support.
func Version() (version, sum string) {
	b, ok := debug.ReadBuildInfo()
	if !ok {
		return "", ""
	}
	for _, m := range b.Deps {
		if m.Path == root {
			if m.Replace != nil {
				switch {
				case m.Replace.Version != "" && m.Replace.Path != "":
					return fmt.Sprintf("%s %s", m.Replace.Path, m.Replace.Version), m.Replace.Sum
				case m.Replace.Version != "":
					return m.Replace.Version, m.Replace.Sum
				case m.Replace.Path != "":
					return m.Replace.Path, m.Replace.Sum
				default:
					return m.Version + "*", ""
				}
			}
			return m.Version, m.Sum
		}
	}
	return "", ""
}
