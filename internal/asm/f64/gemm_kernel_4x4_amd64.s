// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// AVX2/FMA3 optimized 4x4 DGEMM micro-kernel
// Computes C[4x4] += A[4xK] * B[Kx4]

//go:build !noasm && !gccgo && !safe
// +build !noasm,!gccgo,!safe

#include "textflag.h"

// Register assignments:
// A matrix pointers: RSI (4 row pointers)
// B matrix pointer: RDX  
// C matrix pointer: RDI
// K (inner dimension): RCX
// lda, ldb, ldc: R8, R9, R10
// C[0,0:3] = Y0
// C[1,0:3] = Y1  
// C[2,0:3] = Y2
// C[3,0:3] = Y3
// B[k,0:3] = Y4
// A broadcasts = Y5-Y8

// func gemmKernel4x4AVX2(a *float64, b *float64, c *float64, k int, lda, ldb, ldc int)
TEXT ·gemmKernel4x4AVX2(SB), NOSPLIT, $0
	MOVQ a+0(FP), SI       // SI = &a[0,0]
	MOVQ b+8(FP), DX       // DX = &b[0,0]
	MOVQ c+16(FP), DI      // DI = &c[0,0]
	MOVQ k+24(FP), CX      // CX = k
	MOVQ lda+32(FP), R8    // R8 = lda (in elements, not bytes)
	MOVQ ldb+40(FP), R9    // R9 = ldb
	MOVQ ldc+48(FP), R10   // R10 = ldc
	
	// Convert strides to bytes
	SHLQ $3, R8            // lda *= 8
	SHLQ $3, R9            // ldb *= 8
	SHLQ $3, R10           // ldc *= 8
	
	// Load C[0:4, 0:4] into registers
	VMOVUPD (DI), Y0          // C[0,0:3]
	ADDQ R10, DI              // DI = &C[1,0]
	VMOVUPD (DI), Y1          // C[1,0:3]
	ADDQ R10, DI              // DI = &C[2,0]
	VMOVUPD (DI), Y2          // C[2,0:3]
	ADDQ R10, DI              // DI = &C[3,0]
	VMOVUPD (DI), Y3          // C[3,0:3]
	SUBQ R10, DI              // Reset DI
	SUBQ R10, DI
	SUBQ R10, DI              // DI = &C[0,0]
	
	// Check if k == 0
	TESTQ CX, CX
	JZ done
	
	// Calculate row pointers for A
	MOVQ SI, R11           // R11 = &A[0,0]
	MOVQ SI, R12           // R12 = &A[1,0]
	ADDQ R8, R12
	MOVQ SI, R13           // R13 = &A[2,0]
	ADDQ R8, R13
	ADDQ R8, R13
	MOVQ SI, R14           // R14 = &A[3,0]
	ADDQ R8, R14
	ADDQ R8, R14
	ADDQ R8, R14

loop:
	// Load B[k,0:3]
	VMOVUPD (DX), Y4
	
	// Broadcast A[0:4,k] and multiply-add
	VBROADCASTSD (R11), Y5     // Y5 = A[0,k]
	VFMADD231PD Y5, Y4, Y0     // C[0,0:3] += A[0,k] * B[k,0:3]
	
	VBROADCASTSD (R12), Y6     // Y6 = A[1,k]
	VFMADD231PD Y6, Y4, Y1     // C[1,0:3] += A[1,k] * B[k,0:3]
	
	VBROADCASTSD (R13), Y7     // Y7 = A[2,k]
	VFMADD231PD Y7, Y4, Y2     // C[2,0:3] += A[2,k] * B[k,0:3]
	
	VBROADCASTSD (R14), Y8     // Y8 = A[3,k]
	VFMADD231PD Y8, Y4, Y3     // C[3,0:3] += A[3,k] * B[k,0:3]
	
	// Advance pointers
	ADDQ $8, R11           // A[0,k] += 1
	ADDQ $8, R12           // A[1,k] += 1
	ADDQ $8, R13           // A[2,k] += 1
	ADDQ $8, R14           // A[3,k] += 1
	ADDQ R9, DX            // B += ldb (next row)
	
	DECQ CX                // k--
	JNZ loop

done:
	// Store C back to memory
	VMOVUPD Y0, (DI)          // C[0,0:3]
	ADDQ R10, DI
	VMOVUPD Y1, (DI)          // C[1,0:3]
	ADDQ R10, DI
	VMOVUPD Y2, (DI)          // C[2,0:3]
	ADDQ R10, DI
	VMOVUPD Y3, (DI)          // C[3,0:3]
	
	VZEROUPPER
	RET
