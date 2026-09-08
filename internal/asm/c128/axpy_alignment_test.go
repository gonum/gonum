// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package c128_test

import (
	"fmt"
	"slices"
	"testing"
	"unsafe"

	"gonum.org/v1/gonum/internal/asm/c128"
)

func alignedAxpyStorage(t *testing.T, mod uintptr) []complex128 {
	t.Helper()
	storage := new(struct {
		A   [128]complex128
		Pad uint64
		B   [128]complex128
	})
	for _, v := range [][]complex128{storage.A[:], storage.B[:]} {
		for i := 0; i < 4; i++ {
			if uintptr(unsafe.Pointer(&v[i]))%16 == mod {
				return v[i:]
			}
		}
	}
	t.Fatal("alignment fixture")
	return nil
}

func TestAxpyAlignment(t *testing.T) {
	for _, n := range []int{0, 1, 2, 3, 4, 5, 7, 8, 9, 15, 16, 17} {
		for _, mod := range []uintptr{0, 8} {
			for _, op := range []string{"Unitary", "UnitaryTo", "Inc", "IncTo"} {
				for _, alias := range []string{"none", "x", "y"} {
					t.Run(fmt.Sprintf("%s/n=%d/mod=%d/alias=%s", op, n, mod, alias), func(t *testing.T) {
						inc := 1
						if op == "Inc" || op == "IncTo" {
							inc = 2
						}
						span := 1
						if n > 0 {
							span = 1 + (n-1)*inc
						}
						x := make([]complex128, span+1)
						y := alignedAxpyStorage(t, mod)[:span+1]
						dst := make([]complex128, span+1)
						for i := range x {
							x[i] = complex128(complex(float64(i+1), -2))
							y[i] = complex128(complex(3, float64(i)))
							dst[i] = 99
						}
						if alias == "x" {
							dst = x
						}
						if alias == "y" {
							dst = y
						}
						if op == "Unitary" || op == "Inc" {
							dst = y
						}
						beforeX, beforeY := slices.Clone(x), slices.Clone(y)
						want := slices.Clone(dst)
						const alpha = complex128(2 - 3i)
						for i := 0; i < n; i++ {
							j := i * inc
							want[j] = alpha*x[j] + y[j]
						}
						switch op {
						case "Unitary":
							c128.AxpyUnitary(alpha, x[:n], y[:n])
						case "UnitaryTo":
							c128.AxpyUnitaryTo(dst[:n], alpha, x[:n], y[:n])
						case "Inc":
							c128.AxpyInc(alpha, x[:span], y[:span], uintptr(n), uintptr(inc), uintptr(inc), 0, 0)
						case "IncTo":
							c128.AxpyIncTo(dst[:span], uintptr(inc), 0, alpha, x[:span], y[:span], uintptr(n), uintptr(inc), uintptr(inc), 0, 0)
						}
						if !slices.Equal(dst, want) {
							t.Fatalf("got %v, want %v", dst, want)
						}
						if &dst[0] != &x[0] && !slices.Equal(x, beforeX) {
							t.Fatal("x changed")
						}
						if &dst[0] != &y[0] && !slices.Equal(y, beforeY) {
							t.Fatal("y changed")
						}
					})
				}
			}
		}
	}
}
