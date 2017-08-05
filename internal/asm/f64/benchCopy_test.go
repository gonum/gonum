// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build go1.7

package f64

import (
	"fmt"
	"testing"
)

var incCopy = []struct {
	len int
	inc []int
}{
	{1, []int{1}},
	{3, []int{1, 2, 4, 10}},
	{10, []int{1, 2, 4, 10}},
	{30, []int{1, 2, 4, 10}},
	{1e2, []int{1, 2, 4, 10}},
	{3e2, []int{1, 2, 4, 10}},
	{1e3, []int{1, 2, 4, 10}},
	{3e3, []int{1, 2, 4, 10}},
	{1e4, []int{1, 2, 4, 10}},
}

func BenchmarkCopy(t *testing.B) {
	naivecopy := func(n int, dst []float64, incDst int, src []float64, incSrc int) {
		for i := 0; i < n; i++ {
			dst[i*incDst] = src[i*incSrc]
		}
	}
	tests := []struct {
		name string
		f    func(n int, dst []float64, incDst int, src []float64, incSrc int)
	}{
		{"NaiveCopy", naivecopy},
		{"Copy", Copy},
	}
	for _, tt := range incCopy {
		for _, inc := range tt.inc {
			for _, test := range tests {
				t.Run(fmt.Sprintf("%s-%d-inc(%d)", test.name, tt.len, inc), func(b *testing.B) {
					x := make([]float64, inc*tt.len)
					y := make([]float64, inc*tt.len)
					b.SetBytes(int64(64 * tt.len))
					b.ResetTimer()
					for i := 0; i < b.N; i++ {
						Copy(tt.len, y, inc, x, inc)
					}
				})
			}
		}
	}
}
