// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"fmt"
	"testing"
)

func TestGerNegativeStrides(t *testing.T) {
	t.Run("Sger", func(t *testing.T) { testGerNegativeStrides(t, Implementation{}.Sger, nil) })
	t.Run("Dger", func(t *testing.T) { testGerNegativeStrides(t, Implementation{}.Dger, nil) })
}
func testGerNegativeStrides[T ~float32 | ~float64](t *testing.T, ger func(int, int, T, []T, int, []T, int, []T, int), alloc func(*testing.T, int) []T) {
	if alloc == nil {
		alloc = func(_ *testing.T, n int) []T { return make([]T, n) }
	}
	for _, shape := range [][2]int{{1, 1}, {1, 8}, {8, 1}, {2, 3}, {3, 2}, {4, 4}, {5, 9}, {7, 15}, {8, 8}, {64, 8}} {
		m, n := shape[0], shape[1]
		for _, incs := range [][2]int{{-3, -3}, {-3, 1}, {1, -3}, {-1, 3}, {3, -1}, {1, 1}, {3, 7}} {
			incX, incY := incs[0], incs[1]
			t.Run(fmt.Sprintf("m=%d/n=%d/incX=%d/incY=%d", m, n, incX, incY), func(t *testing.T) {
				defer func() {
					if p := recover(); p != nil {
						t.Fatalf("unexpected memory access or panic: %v", p)
					}
				}()
				abs := func(v int) int {
					if v < 0 {
						return -v
					}
					return v
				}
				x, y := alloc(t, 1+(m-1)*abs(incX)), alloc(t, 1+(n-1)*abs(incY))
				for i := range x {
					x[i] = T(i%19+1) / 8
				}
				for i := range y {
					y[i] = T(i%13+1) / 4
				}
				lda := n + 3
				a := alloc(t, (m-1)*lda+n)
				for i := range a {
					a[i] = T(i%17-8) / 16
				}
				originalX, originalY := append([]T(nil), x...), append([]T(nil), y...)
				want := append([]T(nil), a...)
				ix, iy := 0, 0
				if incX < 0 {
					ix = (1 - m) * incX
				}
				if incY < 0 {
					iy = (1 - n) * incY
				}
				for i := 0; i < m; i++ {
					for j := 0; j < n; j++ {
						want[i*lda+j] += (T(.5) * x[ix+i*incX]) * y[iy+j*incY]
					}
				}
				ger(m, n, T(.5), x, incX, y, incY, a, lda)
				for i, v := range a {
					if v != want[i] {
						t.Fatalf("matrix/padding index%d got%g want%g", i, v, want[i])
					}
				}
				for i, v := range x {
					if v != originalX[i] {
						t.Fatalf("x index%d changed", i)
					}
				}
				for i, v := range y {
					if v != originalY[i] {
						t.Fatalf("y index%d changed", i)
					}
				}
			})
		}
	}
}
