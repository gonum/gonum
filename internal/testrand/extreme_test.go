// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testrand

import (
	"math"
	"testing"

	"golang.org/x/exp/rand"
)

func TestExtreme_NaN(t *testing.T) {
	src := rand.NewSource(1)
	rnd := rand.New(src)
	ext := newExtreme(0, 1, rnd)

	check64 := func(v float64) {
		if !math.IsNaN(v) {
			t.Errorf("expected NaN, got %v", v)
		}
	}
	check32 := func(v float32) {
		if !math.IsNaN(float64(v)) {
			t.Errorf("expected NaN, got %v", v)
		}
	}

	check64(ext.ExpFloat64())
	check64(ext.Float64())
	check64(ext.NormFloat64())
	check32(ext.Float32())
}
