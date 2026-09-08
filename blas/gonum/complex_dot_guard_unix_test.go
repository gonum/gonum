// Copyright ©2026 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//go:build linux || darwin

package gonum

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"
	"unsafe"

	"gonum.org/v1/gonum/internal/asm/c64"
)

// Each backend runs in a child so an assembly overread fails this test without
// terminating unrelated package tests. The children exercise only valid spans.
func TestComplexDotIncGuardPages(t *testing.T) {
	type dotFunc func(int, []complex64, int, []complex64, int) complex64
	calls := []struct {
		name      string
		conjugate bool
		call      dotFunc
	}{
		{"Cdotc", true, Implementation{}.Cdotc},
		{"Cdotu", false, Implementation{}.Cdotu},
		{"DotcInc", true, func(n int, x []complex64, incX int, y []complex64, incY int) complex64 {
			return c64.DotcInc(x, y, uintptr(n), uintptr(incX), uintptr(incY), uintptr(complexDotGuardStart(n, incX)), uintptr(complexDotGuardStart(n, incY)))
		}},
		{"DotuInc", false, func(n int, x []complex64, incX int, y []complex64, incY int) complex64 {
			return c64.DotuInc(x, y, uintptr(n), uintptr(incX), uintptr(incY), uintptr(complexDotGuardStart(n, incX)), uintptr(complexDotGuardStart(n, incY)))
		}},
	}
	const childKey = "GONUM_COMPLEX_DOT_GUARD_CHILD"
	child := os.Getenv(childKey)
	for _, test := range calls {
		if child == "" {
			t.Run(test.name, func(t *testing.T) {
				binary, err := os.Executable()
				if err != nil {
					t.Fatal(err)
				}
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				defer cancel()
				cmd := exec.CommandContext(ctx, binary, "-test.run=^TestComplexDotIncGuardPages$", "-test.v", "-test.count=1")
				cmd.Env = append(os.Environ(), childKey+"="+test.name)
				output, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("guard child: %v\n%s", err, output)
				}
			})
			continue
		}
		if child != test.name {
			continue
		}
		// n=1/inc=2 is the original public valid-input regression. The y
		// allocation ends at the protected page, with no spare element.
		strides := [][2]int{{2, 2}, {1, 1}, {3, 5}, {-2, -2}, {2, -3}, {-3, 2}}
		for _, n := range []int{1, 0, 2, 3, 4, 5, 6, 7, 8, 9, 15, 16, 17, 31, 32, 33, 63, 64, 65} {
			for _, inc := range strides {
				for _, xEnd := range []bool{false, true} {
					for _, yEnd := range []bool{true, false} {
						name := fmt.Sprintf("%s/n=%d/inc=%d,%d/end=%t,%t", test.name, n, inc[0], inc[1], xEnd, yEnd)
						t.Run(name, func(t *testing.T) {
							x, protectX := complexDotGuardSlice(t, complexDotGuardSpan(n, inc[0]), xEnd)
							y, protectY := complexDotGuardSlice(t, complexDotGuardSpan(n, inc[1]), yEnd)
							ix, iy := complexDotGuardStart(n, inc[0]), complexDotGuardStart(n, inc[1])
							var wr, wi int64
							for k := 0; k < n; k++ {
								xr, xi := int64(k%11-5), int64(k%7-3)
								yr, yi := int64(k%13-6), int64(k%5-2)
								x[ix], y[iy] = complex(float32(xr)/4, float32(xi)/4), complex(float32(yr)/8, float32(yi)/8)
								if test.conjugate {
									wr += xr*yr + xi*yi
									wi += xr*yi - xi*yr
								} else {
									wr += xr*yr - xi*yi
									wi += xr*yi + xi*yr
								}
								ix += inc[0]
								iy += inc[1]
							}
							protectX()
							protectY()
							want := complex(float32(wr)/32, float32(wi)/32)
							if got := test.call(n, x, inc[0], y, inc[1]); got != want {
								t.Fatalf("got %v want %v", got, want)
							}
						})
					}
				}
			}
		}
		return
	}
	if child != "" {
		t.Fatalf("unknown child %q", child)
	}
}

func complexDotGuardStart(n, inc int) int {
	if n > 0 && inc < 0 {
		return (1 - n) * inc
	}
	return 0
}

func complexDotGuardSpan(n, inc int) int {
	if n <= 1 {
		return n
	}
	if inc < 0 {
		inc = -inc
	}
	return 1 + (n-1)*inc
}

func complexDotGuardSlice(t *testing.T, n int, atEnd bool) ([]complex64, func()) {
	t.Helper()
	if n == 0 {
		return nil, func() {}
	}
	page := syscall.Getpagesize()
	if n < 0 || n > page/8 {
		t.Fatal("guard fixture exceeds one page")
	}
	data, err := syscall.Mmap(-1, 0, 3*page, syscall.PROT_NONE, syscall.MAP_PRIVATE|syscall.MAP_ANON)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := syscall.Munmap(data); err != nil {
			t.Error(err)
		}
	})
	middle := data[page : 2*page]
	if err := syscall.Mprotect(middle, syscall.PROT_READ|syscall.PROT_WRITE); err != nil {
		t.Fatal(err)
	}
	start := page
	if atEnd {
		start = 2*page - 8*n
	}
	x := unsafe.Slice((*complex64)(unsafe.Pointer(&data[start])), n)
	for i := range x {
		x[i] = complex(1000, -1000)
	}
	return x, func() {
		if err := syscall.Mprotect(middle, syscall.PROT_READ); err != nil {
			t.Fatal(err)
		}
	}
}
