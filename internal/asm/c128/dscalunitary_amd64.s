// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

#include "textflag.h"

#define SRC SI
#define DST SI
#define LEN CX
#define IDX AX
#define TAIL BX
#define ALPHA X0
#define ALPHA_2 X1

// func DscalUnitary(alpha float64, x []complex128)
TEXT ·DscalUnitary(SB), NOSPLIT, $0
	MOVQ x_base+8(FP), SRC // SRC = &x
	MOVQ x_len+16(FP), LEN // LEN = len(x)
	CMPQ LEN, $0           // if LEN == 0 { return }
	JE   dscal_end

	MOVSD alpha+0(FP), ALPHA // ALPHA = alpha
	XORQ  IDX, IDX           // IDX = 0
	MOVSD ALPHA, ALPHA_2     // Copy ALPHA to ALPHA_2 for pipelining
	MOVQ  LEN, TAIL          // TAIL = LEN
	SHRQ  $2, LEN            // LEN = floor( n / 4 )
	JZ    dscal_tail         // if LEN == 0 { goto caxy_tail }

dscal_loop: // do {
	MOVSD (SRC)(IDX*8), X2   // X_i = real(x[i])
	MOVSD 16(SRC)(IDX*8), X3
	MOVSD 32(SRC)(IDX*8), X4
	MOVSD 48(SRC)(IDX*8), X5

	MULSD ALPHA, X2   // X_i *= alpha
	MULSD ALPHA_2, X3
	MULSD ALPHA, X4
	MULSD ALPHA_2, X5

	MOVSD X2, (DST)(IDX*8)   // real(x[i]) = X_i
	MOVSD X3, 16(DST)(IDX*8)
	MOVSD X4, 32(DST)(IDX*8)
	MOVSD X5, 48(DST)(IDX*8)

	ADDQ $8, IDX    // IDX += 8
	DECQ LEN
	JNZ  dscal_loop // } while --LEN > 0

dscal_tail:
	ANDQ $3, TAIL  // TAIL = TAIL % 4
	JE   dscal_end // if TAIL == 0 { return }

dscal_tail_loop: // do {
	MOVSD (SRC)(IDX*8), X2 // X_i = real(x[i])
	MULSD ALPHA, X2        // X_i *= alpha
	MOVSD X2, (DST)(IDX*8) // real(x[i]) = X_i
	ADDQ  $2, IDX          // IDX += 2
	DECQ  TAIL
	JNZ   dscal_tail_loop  // } while --TAIL > 0

dscal_end:
	RET
