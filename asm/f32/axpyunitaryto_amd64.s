// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

#include "textflag.h"

// func AxpyUnitaryTo(dst []float32, alpha float32, x, y []float32)
TEXT ·AxpyUnitaryTo(SB), NOSPLIT, $0
	MOVQ    dst_base+0(FP), DI // Load data buffer pointers
	MOVQ    x_base+32(FP), SI
	MOVQ    y_base+56(FP), DX
	MOVQ    x_len+40(FP), CX   // CX = min( len(x), len(y), len(dst) )
	CMPQ    y_len+64(FP), CX
	CMOVLEQ y_len+64(FP), CX
	CMPQ    dst_len+8(FP), CX
	CMOVLEQ dst_len+8(FP), CX
	CMPQ    CX, $0             // Empty return
	JE      caxy_end
	MOVSS   alpha+24(FP), X0
	SHUFPS  $0, X0, X0         // Load alpha 4 times
	MOVUPS  X0, X2             // Copy to X2 for pipelining
	XORQ    AX, AX             // i = 0
	MOVQ    CX, BX
	ANDQ    $7, BX             // BX = len % 8
	SHRQ    $3, CX             // CX = int(len / 8)
	CMPQ    CX, $0
	JE      caxy_tail_start

caxy_loop:
	MOVUPS (SI)(AX*4), X1   // xmm = x[i:i+4]
	MOVUPS 16(SI)(AX*4), X3
	MULPS  X0, X1           // xmm *= a
	MULPS  X2, X3
	ADDPS  (DX)(AX*4), X1   // xmm += y[i:i+4]
	ADDPS  16(DX)(AX*4), X3
	MOVUPS X1, (DI)(AX*4)   // dst[i:i+4] = xmm
	MOVUPS X3, 16(DI)(AX*4)
	ADDQ   $8, AX           // i+=8
	LOOPNE caxy_loop        // while (--CX) > 0
	CMPQ   BX, $0
	JE     caxy_end

caxy_tail_start:
	MOVQ BX, CX

caxy_tail:
	MOVSS  (SI)(AX*4), X1
	MULSS  X0, X1
	ADDSS  (DX)(AX*4), X1
	MOVSS  X1, (DI)(AX*4)
	INCQ   AX
	LOOPNE caxy_tail

caxy_end:
	RET
