// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

#include "textflag.h"

// func AxpyUnitary(alpha complex64, x, y []complex64)
TEXT ·AxpyUnitary(SB), NOSPLIT, $0
	MOVQ    x_base+8(FP), SI
	MOVQ    y_base+32(FP), DI
	MOVQ    x_len+16(FP), DX
	CMPQ    y_len+40(FP), DX
	CMOVLEQ y_len+40(FP), DX
	CMPQ    DX, $0
	JE      caxy_end
	MOVSD   alpha+0(FP), X0   // (0,0,ar,ai)
	SHUFPD  $0, X0, X0        // (ar,ai,ar,ai)
	MOVAPS  X0, X1
	SHUFPS  $0x11, X1, X1     // (ai,ar,ai,ar)
	XORQ    AX, AX
	CMPQ    DX, $2
	JL      caxy_tail
	DECQ    DX

caxy_loop:
	// MOVSHDUP (SI)(AX*8), X2	// Load and duplicate real elements (x2r, x2r, x1r, x1r)
	// MOVSLDUP (SI)(AX*8), X3	// Load and duplicate imag elements (x2i, x2i, x1i, x1i)
	BYTE $0xF3; BYTE $0x0F; BYTE $0x16; BYTE $0x14; BYTE $0xC6
	BYTE $0xF3; BYTE $0x0F; BYTE $0x12; BYTE $0x1C; BYTE $0xC6

	MULPS X1, X2 // (ai*x2r, ar*x2r, ai*x1r, ar*x1r)
	MULPS X0, X3 // (ar*x2i, ai*x2i, ar*x1i, ai*x1i)

	// ADDSUBPS X3, X2  	// (ai*x2r+ar*x2i, ar*x2r-ai*x2i, ai*x1r+ar*x1i, ar*x1r-ai*x1i)
	BYTE $0xF2; BYTE $0x0F; BYTE $0xD0; BYTE $0xDA

	ADDPS  (DI)(AX*8), X3 // Add y2,y1 to a*(x2,x1)
	MOVUPS X3, (DI)(AX*8) // Write result back to y2,y1
	ADDQ   $2, AX
	CMPQ   AX, DX
	JL     caxy_loop
	JG     caxy_end

caxy_tail: // Same calculation, but read in values to avoid trampling memory
	MOVSD (SI)(AX*8), X2
	MOVSD (SI)(AX*8), X3

	// MOVSHDUP X2, X2
	// MOVSLDUP X3, X3
	BYTE $0xF3; BYTE $0x0F; BYTE $0x16; BYTE $0xD2
	BYTE $0xF3; BYTE $0x0F; BYTE $0x12; BYTE $0xDB

	MULPS X1, X2 // (ai*x2r, ar*x2r, ai*x1r, ar*x1r)
	MULPS X0, X3 // (ar*x2i, ai*x2i, ar*x1i, ai*x1i)

	// ADDSUBPS X2, X3  	// (ai*x2r+ar*x2i, ar*x2r-ai*x2i, ai*x1r+ar*x1i, ar*x1r-ai*x1i)
	BYTE $0xF2; BYTE $0x0F; BYTE $0xD0; BYTE $0xDA

	MOVSD (DI)(AX*8), X4
	ADDPS X4, X3
	MOVSD X3, (DI)(AX*8)

caxy_end:
	RET
