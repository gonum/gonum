// Copyright ©2024 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
//
// AVX2 optimized version of ScalUnitary
// x[i] *= alpha

//go:build !noasm && !gccgo && !safe
// +build !noasm,!gccgo,!safe

#include "textflag.h"

// func ScalUnitaryAVX2(alpha float64, x []float64)
TEXT ·ScalUnitaryAVX2(SB), NOSPLIT, $0
	MOVSD alpha+0(FP), X0      // X0 = alpha
	VBROADCASTSD X0, Y0        // Y0 = [alpha, alpha, alpha, alpha]
	MOVQ x+8(FP), SI           // SI = &x[0]
	MOVQ x_len+16(FP), CX      // CX = len(x)
	
	// Check if length is 0
	TESTQ CX, CX
	JZ end
	
	// Process 16 elements at a time (4 YMM registers)
	MOVQ CX, BX
	SHRQ $4, BX               // BX = n / 16
	JZ tail_start
	
loop16:
	// Load and multiply using AVX2
	VMOVUPD (SI), Y1
	VMULPD Y0, Y1, Y1
	VMOVUPD Y1, (SI)
	
	VMOVUPD 32(SI), Y2
	VMULPD Y0, Y2, Y2
	VMOVUPD Y2, 32(SI)
	
	VMOVUPD 64(SI), Y3
	VMULPD Y0, Y3, Y3
	VMOVUPD Y3, 64(SI)
	
	VMOVUPD 96(SI), Y4
	VMULPD Y0, Y4, Y4
	VMOVUPD Y4, 96(SI)
	
	ADDQ $128, SI             // Advance by 16 elements
	DECQ BX
	JNZ loop16
	
tail_start:
	// Handle remaining elements
	MOVQ CX, BX
	ANDQ $15, BX              // BX = n % 16
	JZ end
	
	// Process 4 elements at a time using YMM
	MOVQ BX, DX
	SHRQ $2, DX               // DX = (n % 16) / 4
	JZ tail_2
	
loop4:
	VMOVUPD (SI), Y1
	VMULPD Y0, Y1, Y1
	VMOVUPD Y1, (SI)
	ADDQ $32, SI
	DECQ DX
	JNZ loop4
	
tail_2:
	// Process 2 elements at a time using SSE
	MOVQ BX, DX
	ANDQ $3, DX               // DX = n % 4
	SHRQ $1, DX               // DX = (n % 4) / 2
	JZ tail_1
	
	MOVUPD (SI), X1
	MULPD X0, X1
	MOVUPD X1, (SI)
	ADDQ $16, SI
	
tail_1:
	// Process last element if present
	MOVQ BX, DX
	ANDQ $1, DX
	JZ end
	
	MOVSD (SI), X1
	MULSD X0, X1
	MOVSD X1, (SI)
	
end:
	VZEROUPPER
	RET
