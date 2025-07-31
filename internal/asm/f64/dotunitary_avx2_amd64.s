// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// AVX2/FMA optimized version of DotUnitary
// sum = x·y = Σ(x[i] * y[i])

//go:build !noasm && !gccgo && !safe
// +build !noasm,!gccgo,!safe

#include "textflag.h"

// func DotUnitaryAVX2(x, y []float64) (sum float64)
TEXT ·DotUnitaryAVX2(SB), NOSPLIT, $0
	MOVQ x+0(FP), SI      // SI = &x[0]
	MOVQ x_len+8(FP), CX  // CX = len(x)
	MOVQ y+24(FP), DI     // DI = &y[0]
	
	// Initialize accumulators
	VXORPD Y0, Y0, Y0     // Y0 = sum[0:3] = 0
	VXORPD Y1, Y1, Y1     // Y1 = sum[4:7] = 0
	VXORPD Y2, Y2, Y2     // Y2 = sum[8:11] = 0
	VXORPD Y3, Y3, Y3     // Y3 = sum[12:15] = 0
	
	// Check if length is 0
	TESTQ CX, CX
	JZ end
	
	// Process 16 elements at a time (4 YMM registers)
	MOVQ CX, BX
	SHRQ $4, BX           // BX = n / 16
	JZ tail_start
	
loop16:
	// Load and multiply-add using FMA
	VMOVUPD (SI), Y4
	VMOVUPD (DI), Y5
	VFMADD231PD Y5, Y4, Y0     // Y0 += x[0:3] * y[0:3]
	
	VMOVUPD 32(SI), Y4
	VMOVUPD 32(DI), Y5
	VFMADD231PD Y5, Y4, Y1     // Y1 += x[4:7] * y[4:7]
	
	VMOVUPD 64(SI), Y4
	VMOVUPD 64(DI), Y5
	VFMADD231PD Y5, Y4, Y2     // Y2 += x[8:11] * y[8:11]
	
	VMOVUPD 96(SI), Y4
	VMOVUPD 96(DI), Y5
	VFMADD231PD Y5, Y4, Y3     // Y3 += x[12:15] * y[12:15]
	
	ADDQ $128, SI         // Advance x by 16 elements
	ADDQ $128, DI         // Advance y by 16 elements
	DECQ BX
	JNZ loop16
	
tail_start:
	// Add all accumulators
	VADDPD Y1, Y0, Y0     // Y0 = Y0 + Y1
	VADDPD Y3, Y2, Y2     // Y2 = Y2 + Y3
	VADDPD Y2, Y0, Y0     // Y0 = Y0 + Y2
	
	// Handle remaining elements
	MOVQ CX, BX
	ANDQ $15, BX          // BX = n % 16
	SHRQ $2, BX           // BX = (n % 16) / 4
	JZ tail_2
	
loop4:
	VMOVUPD (SI), Y4
	VMOVUPD (DI), Y5
	VFMADD231PD Y5, Y4, Y0
	ADDQ $32, SI
	ADDQ $32, DI
	DECQ BX
	JNZ loop4
	
tail_2:
	// Process 2 elements at a time using SSE
	MOVQ CX, BX
	ANDQ $3, BX           // BX = n % 4
	SHRQ $1, BX           // BX = (n % 4) / 2
	JZ tail_1
	
	MOVUPD (SI), X4
	MOVUPD (DI), X5
	MULPD X5, X4
	ADDPD X4, X0          // Use lower half of Y0
	ADDQ $16, SI
	ADDQ $16, DI
	
tail_1:
	// Process last element if present
	MOVQ CX, BX
	ANDQ $1, BX
	JZ reduce
	
	MOVSD (SI), X4
	MOVSD (DI), X5
	MULSD X5, X4
	ADDSD X4, X0
	
reduce:
	// Horizontal sum of Y0
	VEXTRACTF128 $1, Y0, X1    // X1 = upper 128 bits of Y0
	ADDPD X1, X0               // X0 = X0 + X1
	MOVAPD X0, X1
	UNPCKHPD X1, X1            // X1 = X0[1]
	ADDSD X1, X0               // X0 = final sum
	
end:
	MOVSD X0, sum+48(FP)       // Return sum
	VZEROUPPER
	RET
