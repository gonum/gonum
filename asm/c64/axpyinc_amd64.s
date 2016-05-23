// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

#include "textflag.h"

// func AxpyInc(alpha complex64, x, y []complex64, n, incX, incY, ix, iy uintptr)
TEXT ·AxpyInc(SB), NOSPLIT, $0
	MOVQ   x_base+8(FP), SI
	MOVQ   y_base+32(FP), DI
	MOVQ   n+56(FP), DX
	XORQ   AX, AX
	CMPQ   DX, $0
	JE     axpyi_end
	MOVQ   ix+80(FP), CX
	MOVQ   iy+88(FP), BX
	LEAQ   (SI)(CX*8), SI    // Calculate location of first indicies
	LEAQ   (DI)(BX*8), DI
	MOVQ   incX+64(FP), CX
	SHLQ   $3, CX
	MOVQ   incY+72(FP), BX
	SHLQ   $3, BX
	MOVSD  alpha+0(FP), X0   // (0,0,ar,ai)
	SHUFPD $0, X0, X0        // (ar,ai,ar,ai)
	MOVAPS X0, X1
	SHUFPS $0x11, X1, X1     // (ai,ar,ai,ar)

axpyi_loop:
	// MOVSHDUP (SI), X2
	// MOVSLDUP (SI), X3
	BYTE $0xF3; BYTE $0x0F; BYTE $0x16; BYTE $0x16
	BYTE $0xF3; BYTE $0x0F; BYTE $0x12; BYTE $0x1E

	MULPS X1, X2 // (ai*x2r, ar*x2r, ai*x1r, ar*x1r)
	MULPS X0, X3 // (ar*x2i, ai*x2i, ar*x1i, ai*x1i)

	// ADDSUBPS X2, X3  	// (ai*x2r+ar*x2i, ar*x2r-ai*x2i, ai*x1r+ar*x1i, ar*x1r-ai*x1i)
	BYTE $0xF2; BYTE $0x0F; BYTE $0xD0; BYTE $0xDA

	ADDPS (DI), X3
	MOVSD X3, (DI)
	ADDQ  CX, SI
	ADDQ  BX, DI
	INCQ  AX
	CMPQ  AX, DX
	JL    axpyi_loop

axpyi_end:
	RET
