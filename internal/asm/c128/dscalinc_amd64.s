// Copyright ©2017 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

#include "textflag.h"

#define SRC SI
#define DST SI
#define LEN CX
#define TAIL BX
#define INC R9
#define INC3 R10
#define ALPHA X0
#define ALPHA_2 X1

// func DscalInc(alpha float64, x []complex128, n, inc int)
TEXT ·DscalInc(SB), NOSPLIT, $0
	MOVQ x_base+8(FP), SRC // SRC = &x
	MOVQ n+32(FP), LEN     // LEN = n
	CMPQ LEN, $0           // if LEN == 0 { return }
	JE   dscal_end

	MOVSD alpha+0(FP), ALPHA // ALPHA = alpha
	MOVQ  inc+40(FP), INC    // INC = inc
	SHLQ  $4, INC            // INC = INC * sizeof(complex128)
	LEAQ  (INC)(INC*2), INC3 // INC3 = 3 * INC
	MOVSD ALPHA, ALPHA_2     // Copy ALPHA and ALPHA_2 for pipelining
	MOVQ  LEN, TAIL          // TAIL = LEN
	SHRQ  $2, LEN            // LEN = floor( n / 4 )
	JZ    dscal_tail         // if LEN == 0 { goto dscal_tail }

dscal_loop: // do {
	MOVSD (SRC), X2         // X_i = real(x[i])
	MOVSD (SRC)(INC*1), X3
	MOVSD (SRC)(INC*2), X4
	MOVSD (SRC)(INC3*1), X5

	MULSD ALPHA, X2   // X_i *= alpha
	MULSD ALPHA_2, X3
	MULSD ALPHA, X4
	MULSD ALPHA_2, X5

	MOVSD X2, (DST)         // real(x[i]) = X_i
	MOVSD X3, (DST)(INC*1)
	MOVSD X4, (DST)(INC*2)
	MOVSD X5, (DST)(INC3*1)

	LEAQ (SRC)(INC*4), SRC // SRC += INC*4
	DECQ LEN
	JNZ  dscal_loop        // } while --LEN > 0

dscal_tail:
	ANDQ $3, TAIL  // TAIL = TAIL % 4
	JE   dscal_end // if TAIL == 0 { return }

dscal_tail_loop: // do {
	MOVSD (SRC), X2       // X_i = real(x[i])
	MULSD ALPHA, X2       // X_i *= alpha
	MOVSD X2, (DST)       // real(x[i]) = X_i
	ADDQ  INC, SRC        // SRC += INC
	DECQ  TAIL
	JNZ   dscal_tail_loop // } while --TAIL > 0

dscal_end:
	RET
