// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

#include "textflag.h"

// MOVSHDUP X3, X2
#define MOVSHDUP_X3_X2 BYTE $0xF3; BYTE $0x0F; BYTE $0x16; BYTE $0xD3
// MOVSLDUP X3, X3
#define MOVSLDUP_X3_X3 BYTE $0xF3; BYTE $0x0F; BYTE $0x12; BYTE $0xDB
// ADDSUBPS X2, X3
#define ADDSUBPS_X2_X3 BYTE $0xF2; BYTE $0x0F; BYTE $0xD0; BYTE $0xDA

// func AxpyUnitary(alpha complex64, x, y []complex64)
TEXT ·AxpyUnitary(SB), NOSPLIT, $0
	MOVQ    x_base+8(FP), SI
	MOVQ    y_base+32(FP), DI
	MOVQ    x_len+16(FP), CX
	CMPQ    y_len+40(FP), CX
	CMOVLEQ y_len+40(FP), CX
	CMPQ    CX, $0
	JE      caxy_end
	PXOR    X0, X0            // Clear work registers and cache-align loop
	PXOR    X1, X1
	MOVSD   alpha+0(FP), X0   // (0,0,ar,ai)
	SHUFPD  $0, X0, X0        // (ar,ai,ar,ai)
	MOVAPS  X0, X1
	SHUFPS  $0x11, X1, X1     // (ai,ar,ai,ar)
	XORQ    AX, AX
	MOVQ    DI, BX
	ANDQ    $15, BX           // Align on 16-byte boundary for ADDPS
	JZ      caxy_no_trim

	// Trim first value in unaligned buffer
	XORPS X2, X2         // Clear work registers and cache-align loop
	XORPS X3, X3
	XORPS X4, X4
	MOVSD (SI)(AX*8), X3
	MOVSHDUP_X3_X2
	MOVSLDUP_X3_X3
	MULPS X1, X2
	MULPS X0, X3
	ADDSUBPS_X2_X3
	MOVSD (DI)(AX*8), X4
	ADDPS X4, X3
	MOVSD X3, (DI)(AX*8)
	INCQ  AX
	DECQ  CX
	JZ    caxy_end

caxy_no_trim:
	MOVQ CX, BX
	ANDQ $1, BX
	SHRQ $1, CX
	JZ   caxy_tail

caxy_loop:
	MOVUPS (SI)(AX*8), X3
	MOVSHDUP_X3_X2        // Load and duplicate real elements (x2r, x2r, x1r, x1r)
	MOVSLDUP_X3_X3        // Load and duplicate imag elements (x2i, x2i, x1i, x1i)
	MULPS  X1, X2         // (ai*x2r, ar*x2r, ai*x1r, ar*x1r)
	MULPS  X0, X3         // (ar*x2i, ai*x2i, ar*x1i, ai*x1i)
	ADDSUBPS_X2_X3        // (ai*x2r+ar*x2i, ar*x2r-ai*x2i, ai*x1r+ar*x1i, ar*x1r-ai*x1i)
	ADDPS  (DI)(AX*8), X3 // Add y2,y1 to a*(x2,x1)
	MOVUPS X3, (DI)(AX*8) // Write result back to y2,y1
	ADDQ   $2, AX
	LOOPNE caxy_loop      // While (--CX) > 0
	CMPQ   BX, $0
	JE     caxy_end

caxy_tail:
	MOVSD (SI)(AX*8), X3
	MOVSHDUP_X3_X2
	MOVSLDUP_X3_X3
	MULPS X1, X2         // (ai*x2r, ar*x2r, ai*x1r, ar*x1r)
	MULPS X0, X3         // (ar*x2i, ai*x2i, ar*x1i, ai*x1i)
	ADDSUBPS_X2_X3       // (ai*x2r+ar*x2i, ar*x2r-ai*x2i, ai*x1r+ar*x1i, ar*x1r-ai*x1i)
	MOVSD (DI)(AX*8), X4
	ADDPS X4, X3
	MOVSD X3, (DI)(AX*8)

caxy_end:
	RET
