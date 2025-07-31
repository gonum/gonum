// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// AVX2/FMA3 optimized 8x8 DGEMM micro-kernel
// Computes C[8x8] += A[8xK] * B[Kx8]

//go:build !noasm && !gccgo && !safe
// +build !noasm,!gccgo,!safe

#include "textflag.h"

// Register usage:
// A matrix pointers: RSI + row offsets
// B matrix: Two YMM registers for current 8 values
// C matrix: 16 YMM registers (Y0-Y15) for 8x8 block
// Each YMM holds 4 doubles, so we need 2 YMM per row

// func gemmKernel8x8AVX2(a *float64, b *float64, c *float64, k int, lda, ldb, ldc int)
TEXT ·gemmKernel8x8AVX2(SB), NOSPLIT, $0
	MOVQ a+0(FP), SI       // SI = &a[0,0]
	MOVQ b+8(FP), DX       // DX = &b[0,0]
	MOVQ c+16(FP), DI      // DI = &c[0,0]
	MOVQ k+24(FP), CX      // CX = k
	MOVQ lda+32(FP), R8    // R8 = lda (in elements)
	MOVQ ldb+40(FP), R9    // R9 = ldb
	MOVQ ldc+48(FP), R10   // R10 = ldc
	
	// Convert strides to bytes
	SHLQ $3, R8            // lda *= 8
	SHLQ $3, R9            // ldb *= 8
	SHLQ $3, R10           // ldc *= 8
	
	// Load C[0:8, 0:8] into Y0-Y15
	VMOVUPD (DI), Y0           // C[0,0:3]
	VMOVUPD 32(DI), Y1         // C[0,4:7]
	ADDQ R10, DI
	VMOVUPD (DI), Y2           // C[1,0:3]
	VMOVUPD 32(DI), Y3         // C[1,4:7]
	ADDQ R10, DI
	VMOVUPD (DI), Y4           // C[2,0:3]
	VMOVUPD 32(DI), Y5         // C[2,4:7]
	ADDQ R10, DI
	VMOVUPD (DI), Y6           // C[3,0:3]
	VMOVUPD 32(DI), Y7         // C[3,4:7]
	ADDQ R10, DI
	VMOVUPD (DI), Y8           // C[4,0:3]
	VMOVUPD 32(DI), Y9         // C[4,4:7]
	ADDQ R10, DI
	VMOVUPD (DI), Y10          // C[5,0:3]
	VMOVUPD 32(DI), Y11        // C[5,4:7]
	ADDQ R10, DI
	VMOVUPD (DI), Y12          // C[6,0:3]
	VMOVUPD 32(DI), Y13        // C[6,4:7]
	ADDQ R10, DI
	VMOVUPD (DI), Y14          // C[7,0:3]
	VMOVUPD 32(DI), Y15        // C[7,4:7]
	
	// Reset DI to start of C
	SUBQ R10, DI
	SUBQ R10, DI
	SUBQ R10, DI
	SUBQ R10, DI
	SUBQ R10, DI
	SUBQ R10, DI
	SUBQ R10, DI
	
	// Check if k == 0
	TESTQ CX, CX
	JZ done
	
	// Calculate row pointers for A
	MOVQ SI, R11              // R11 = &A[0,0]
	LEAQ (SI)(R8*1), R12      // R12 = &A[1,0]
	LEAQ (SI)(R8*2), R13      // R13 = &A[2,0]
	LEAQ (R12)(R8*2), R14     // R14 = &A[3,0]
	LEAQ (SI)(R8*4), R15      // R15 = &A[4,0]
	
	// We'll compute remaining rows using offsets from R15
	// A[5] = R15 + R8
	// A[6] = R15 + 2*R8
	// A[7] = R15 + 3*R8

loop:
	// Load B[k,0:7]
	VMOVUPD (DX), Y28          // B[k,0:3]
	VMOVUPD 32(DX), Y29        // B[k,4:7]
	
	// Row 0
	VBROADCASTSD (R11), Y30
	VFMADD231PD Y30, Y28, Y0
	VFMADD231PD Y30, Y29, Y1
	
	// Row 1
	VBROADCASTSD (R12), Y30
	VFMADD231PD Y30, Y28, Y2
	VFMADD231PD Y30, Y29, Y3
	
	// Row 2
	VBROADCASTSD (R13), Y30
	VFMADD231PD Y30, Y28, Y4
	VFMADD231PD Y30, Y29, Y5
	
	// Row 3
	VBROADCASTSD (R14), Y30
	VFMADD231PD Y30, Y28, Y6
	VFMADD231PD Y30, Y29, Y7
	
	// Row 4
	VBROADCASTSD (R15), Y30
	VFMADD231PD Y30, Y28, Y8
	VFMADD231PD Y30, Y29, Y9
	
	// Row 5
	VBROADCASTSD (R15)(R8*1), Y30
	VFMADD231PD Y30, Y28, Y10
	VFMADD231PD Y30, Y29, Y11
	
	// Row 6
	VBROADCASTSD (R15)(R8*2), Y30
	VFMADD231PD Y30, Y28, Y12
	VFMADD231PD Y30, Y29, Y13
	
	// Row 7
	MOVQ R15, AX
	ADDQ R8, AX
	ADDQ R8, AX
	ADDQ R8, AX
	VBROADCASTSD (AX), Y30
	VFMADD231PD Y30, Y28, Y14
	VFMADD231PD Y30, Y29, Y15
	
	// Advance pointers
	ADDQ $8, R11
	ADDQ $8, R12
	ADDQ $8, R13
	ADDQ $8, R14
	ADDQ $8, R15
	ADDQ R9, DX               // B += ldb
	
	DECQ CX
	JNZ loop

done:
	// Store C back to memory
	VMOVUPD Y0, (DI)
	VMOVUPD Y1, 32(DI)
	ADDQ R10, DI
	VMOVUPD Y2, (DI)
	VMOVUPD Y3, 32(DI)
	ADDQ R10, DI
	VMOVUPD Y4, (DI)
	VMOVUPD Y5, 32(DI)
	ADDQ R10, DI
	VMOVUPD Y6, (DI)
	VMOVUPD Y7, 32(DI)
	ADDQ R10, DI
	VMOVUPD Y8, (DI)
	VMOVUPD Y9, 32(DI)
	ADDQ R10, DI
	VMOVUPD Y10, (DI)
	VMOVUPD Y11, 32(DI)
	ADDQ R10, DI
	VMOVUPD Y12, (DI)
	VMOVUPD Y13, 32(DI)
	ADDQ R10, DI
	VMOVUPD Y14, (DI)
	VMOVUPD Y15, 32(DI)
	
	VZEROUPPER
	RET
