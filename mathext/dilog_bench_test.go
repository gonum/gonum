// Copyright ©2025 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mathext

import (
	"fmt"
	"testing"
)

var li2Sink complex128

func BenchmarkLi2(b *testing.B) {
	cases := []complex128{
		0.1 + 0.1i, -0.3 + 0.39i, 0.001 - 0.49i, // |z| < 0.5
		-0.9999 + 0.001i, 0.5 + 0.7i, -0.8 - 0.0001i, // 0.5 < |z| < 1
		-1.1 + 0.1i, 5 + 0i, -10 + 0i, 1000 + 1e4i, -1791.91931 + 0.5i, // |z| > 1
	}
	for _, z := range cases {
		b.Run(fmt.Sprintf("z=%.6g", z), func(b *testing.B) {
			b.ReportAllocs()
			var r complex128
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				r = Li2(z)
			}
			li2Sink = r
		})
	}
}
