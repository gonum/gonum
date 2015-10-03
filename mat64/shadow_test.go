// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"math/rand"
	"testing"
)

func TestDenseOverlaps(t *testing.T) {
	type view struct {
		i, j, r, c int
		*Dense
	}

	for r := 1; r < 20; r++ {
		for c := 1; c < 20; c++ {
			m := NewDense(r, c, nil)
			panicked, message := panics(func() { m.checkOverlap(m.RawMatrix()) })
			if !panicked {
				t.Error("expected matrix overlap with self")
			}
			if message != regionIdentity {
				t.Errorf("unexpected panic message for self overlap: got: %q want: %q", message, regionIdentity)
			}

			for i := 0; i < 1000; i++ {
				var views [2]view
				for k := range views {
					if r > 1 {
						views[k].i = rand.Intn(r - 1)
					}
					if c > 1 {
						views[k].j = rand.Intn(c - 1)
					}
					if r > 1 {
						views[k].r = rand.Intn(r-views[k].i-1) + 1
					} else {
						views[k].r = 1
					}
					if c > 1 {
						views[k].c = rand.Intn(c-views[k].j-1) + 1
					} else {
						views[k].c = 1
					}
					views[k].Dense = m.View(views[k].i, views[k].j, views[k].r, views[k].c).(*Dense)

					panicked, _ = panics(func() { m.checkOverlap(views[k].RawMatrix()) })
					if !panicked {
						t.Errorf("expected matrix (%d×%d) overlap with view {rows=%d:%d, cols=%d:%d}",
							r, c, views[k].i, views[k].i+views[k].r, views[k].j, views[k].j+views[k].c)
					}
					panicked, _ = panics(func() { views[k].checkOverlap(m.RawMatrix()) })
					if !panicked {
						t.Errorf("expected view {rows=%d:%d, cols=%d:%d} overlap with parent (%d×%d)",
							views[k].i, views[k].i+views[k].r, views[k].j, views[k].j+views[k].c, r, c)
					}
				}

				overlapRows := intervalsOverlap(
					interval{views[0].i, views[0].i + views[0].r},
					interval{views[1].i, views[1].i + views[1].r},
				)
				overlapCols := intervalsOverlap(
					interval{views[0].j, views[0].j + views[0].c},
					interval{views[1].j, views[1].j + views[1].c},
				)
				want := overlapRows && overlapCols

				for k, v := range views {
					w := views[1-k]
					got, _ := panics(func() { v.checkOverlap(w.RawMatrix()) })
					if got != want {
						t.Errorf("unexpected result for overlap test for {rows=%d:%d, cols=%d:%d} with {rows=%d:%d, cols=%d:%d}: got: %t want: %t",
							v.i, v.i+v.r, v.j, v.j+v.c,
							w.i, w.i+w.r, w.j, w.j+w.c,
							got, want)
					}
				}
			}
		}
	}
}

type interval struct{ from, to int }

func intervalsOverlap(a, b interval) bool {
	return a.to > b.from && b.to > a.from
}
