// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build !noasm,!appengine

#include "textflag.h"

#define X_PTR R8
#define Y_PTR R9
#define IDX SI
#define LEN DI
#define TAIL DX

// func DdotUnitaryAVX(x, y []float64) (sum float64)
// This function assumes len(y) >= len(x).
TEXT ·DotUnitaryAVX(SB), NOSPLIT, $0
	MOVQ x+0(FP), X_PTR
	MOVQ x_len+8(FP), LEN // n = len(x)
	MOVQ y+24(FP), Y_PTR

	VXORPS Y7, Y7, Y7 // sum = 0
	VXORPS Y8, Y8, Y8 // sum = 0
	XORQ   IDX, IDX   // i = 0

	MOVQ LEN, TAIL
	SHRQ $4, LEN   // LEN = floor( n / 16 )
	JZ   dot_tail8 // if LEN == 0 { goto dot_tail8 }

loop_uni:
	// sum += x[i] * y[i] unrolled 16x.
	VMOVUPS     0(X_PTR)(IDX*8), Y0
	VMOVUPS     32(X_PTR)(IDX*8), Y1
	VMOVUPS     64(X_PTR)(IDX*8), Y2
	VMOVUPS     96(X_PTR)(IDX*8), Y3
	VFMADD231PD 0(Y_PTR)(IDX*8), Y0, Y7
	VFMADD231PD 32(Y_PTR)(IDX*8), Y1, Y8
	VFMADD231PD 64(Y_PTR)(IDX*8), Y2, Y7
	VFMADD231PD 96(Y_PTR)(IDX*8), Y3, Y8

	ADDQ $16, IDX // i += 16
	DECQ LEN
	JNZ  loop_uni // if n > 0 { goto loop_uni }

	ANDQ $15, TAIL // TAIL = n %16
	JZ   end_uni   // if TAIL == 0 { goto end_uni }

dot_tail8:
	TESTQ $8, TAIL
	JZ    dot_tail4

	VMOVUPS     0(X_PTR)(IDX*8), Y0
	VMOVUPS     32(X_PTR)(IDX*8), Y1
	VFMADD231PD 0(Y_PTR)(IDX*8), Y0, Y7
	VFMADD231PD 32(Y_PTR)(IDX*8), Y1, Y8

	ADDQ $8, IDX // i += 8

dot_tail4:
	TESTQ $4, TAIL
	JZ    dot_tail2

	VMOVUPS     0(X_PTR)(IDX*8), Y0
	VFMADD231PD 0(Y_PTR)(IDX*8), Y0, Y7

	ADDQ $4, IDX // i += 4

dot_tail2:
	// Collapse sum to 128-bit register
	// VFMADD will clear bits 128-255 when used with XMM registers
	VADDPD       Y8, Y7, Y7 // Y7 += Y8
	VEXTRACTF128 $1, Y7, X8 // X8 = Y7[2:3] [ X7 = X7[0:1] ]

	TESTQ $2, TAIL
	JZ    dot_tail1

	VMOVUPS     0(X_PTR)(IDX*8), X0
	VFMADD231PD 0(Y_PTR)(IDX*8), X0, X7

	ADDQ $2, IDX // i += 2

dot_tail1:
	TESTQ $1, TAIL
	JZ    end_uni

	VMOVSD      0(X_PTR)(IDX*8), X0
	VFMADD231SD 0(Y_PTR)(IDX*8), X0, X7

	INCQ IDX

end_uni:
	// Add the four sums together.
	VHADDPD X8, X7, X7
	VHADDPD X7, X7, X7
	MOVSD   X7, sum+48(FP) // Return final sum.
	RET
