// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package f64

import (
	"fmt"
	"testing"
	"unsafe"
	"golang.org/x/sys/cpu"
)

// Direct access to implementations for benchmarking
func BenchmarkAxpyUnitaryAVX2(b *testing.B) {
	sizes := []int{10, 100, 1000, 10000, 100000}
	
	for _, size := range sizes {
		x := make([]float64, size)
		y := make([]float64, size)
		
		// Initialize with some data
		for i := range x {
			x[i] = float64(i)
			y[i] = float64(i * 2)
		}
		
		alpha := 2.5
		
		b.Run("Size"+fmt.Sprintf("%d", size), func(b *testing.B) {
			b.Run("Current", func(b *testing.B) {
				b.SetBytes(int64(size * 8 * 2)) // 2 arrays of float64
				for i := 0; i < b.N; i++ {
					AxpyUnitary(alpha, x, y)
				}
			})
			
			// Benchmark SSE2 directly
			b.Run("SSE2", func(b *testing.B) {
				b.SetBytes(int64(size * 8 * 2))
				for i := 0; i < b.N; i++ {
					axpyUnitarySSE2Asm(alpha, x, y)
				}
			})
			
			// Benchmark AVX2 if available
			if cpu.X86.HasAVX2 {
				b.Run("AVX2", func(b *testing.B) {
					b.SetBytes(int64(size * 8 * 2))
					for i := 0; i < b.N; i++ {
						AxpyUnitaryAVX2(alpha, x, y)
					}
				})
			}
			
			// Benchmark FMA if available
			if cpu.X86.HasAVX2 && cpu.X86.HasFMA {
				b.Run("FMA", func(b *testing.B) {
					b.SetBytes(int64(size * 8 * 2))
					for i := 0; i < b.N; i++ {
						AxpyUnitaryFMA(alpha, x, y)
					}
				})
			}
		})
	}
}

// Benchmark comparing aligned vs unaligned data
func BenchmarkAxpyUnitaryAlignmentAVX2(b *testing.B) {
	size := 10000
	
	// Aligned allocation (32-byte aligned for AVX2)
	xAligned := make([]float64, size+4)
	yAligned := make([]float64, size+4)
	
	// Find 32-byte aligned offset
	xOffset := 0
	for i := 0; i < 4; i++ {
		if uintptr(unsafe.Pointer(&xAligned[i]))%32 == 0 {
			xOffset = i
			break
		}
	}
	yOffset := 0
	for i := 0; i < 4; i++ {
		if uintptr(unsafe.Pointer(&yAligned[i]))%32 == 0 {
			yOffset = i
			break
		}
	}
	
	x := xAligned[xOffset : xOffset+size]
	y := yAligned[yOffset : yOffset+size]
	
	// Unaligned (deliberately misaligned)
	xUnaligned := make([]float64, size+1)[1:]
	yUnaligned := make([]float64, size+1)[1:]
	
	alpha := 2.5
	
	b.Run("Aligned", func(b *testing.B) {
		b.SetBytes(int64(size * 8 * 2))
		for i := 0; i < b.N; i++ {
			AxpyUnitary(alpha, x, y)
		}
	})
	
	b.Run("Unaligned", func(b *testing.B) {
		b.SetBytes(int64(size * 8 * 2))
		for i := 0; i < b.N; i++ {
			AxpyUnitary(alpha, xUnaligned, yUnaligned)
		}
	})
}