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

// MOVSHDUP X5, X4
#define MOVSHDUP_X5_X4 BYTE $0xF3; BYTE $0x0F; BYTE $0x16; BYTE $0xE5
// MOVSLDUP X5, X5
#define MOVSLDUP_X5_X5 BYTE $0xF3; BYTE $0x0F; BYTE $0x12; BYTE $0xED
// ADDSUBPS X4, X5
#define ADDSUBPS_X4_X5 BYTE $0xF2; BYTE $0x0F; BYTE $0xD0; BYTE $0xEC

// MOVSHDUP X7, X6
#define MOVSHDUP_X7_X6 BYTE $0xF3; BYTE $0x0F; BYTE $0x16; BYTE $0xF7
// MOVSLDUP X7, X7
#define MOVSLDUP_X7_X7 BYTE $0xF3; BYTE $0x0F; BYTE $0x12; BYTE $0xFF
// ADDSUBPS X6, X7
#define ADDSUBPS_X6_X7 BYTE $0xF2; BYTE $0x0F; BYTE $0xD0; BYTE $0xFE

// MOVSHDUP X9, X8
#define MOVSHDUP_X9_X8 BYTE $0xF3; BYTE $0x45; BYTE $0x0F; BYTE $0x16; BYTE $0xC1
// MOVSLDUP X9, X9
#define MOVSLDUP_X9_X9 BYTE $0xF3; BYTE $0x45; BYTE $0x0F; BYTE $0x12; BYTE $0xC9
// ADDSUBPS X8, X9
#define ADDSUBPS_X8_X9 BYTE $0xF2; BYTE $0x45; BYTE $0x0F; BYTE $0xD0; BYTE $0xC8

// func AxpyInc(alpha complex64, x, y []complex64, n, incX, incY, ix, iy uintptr)
TEXT ·AxpyInc(SB), NOSPLIT, $0
	MOVQ   x_base+8(FP), SI
	MOVQ   y_base+32(FP), DI
	MOVQ   n+56(FP), CX
	CMPQ   CX, $0            // if n==0, return
	JE     axpyi_end
	MOVQ   ix+80(FP), R8     // Load the first indicies
	MOVQ   iy+88(FP), R9
	LEAQ   (SI)(R8*8), SI    // Calculate addrress of first indicies
	LEAQ   (DI)(R9*8), DI
	MOVQ   incX+64(FP), R8   // Incrementors*8 for easy iteration (ADDQ)
	SHLQ   $3, R8
	MOVQ   incY+72(FP), R9
	SHLQ   $3, R9
	MOVSD  alpha+0(FP), X0   // XO:(0,0,ar,ai)
	MOVAPS X0, X1
	SHUFPS $0x11, X1, X1     // X1:(0,0,ai,ar)
	MOVAPS X0, X10
	MOVAPS X1, X11
	MOVQ   CX, BX
	ANDQ   $3, CX
	SHRQ   $2, BX
	JZ     axpyi_tail

axpyi_loop:
	MOVSD (SI), X3
	MOVSD (SI)(R8*1), X5
	LEAQ  (SI)(R8*2), SI
	MOVSD (SI), X7
	MOVSD (SI)(R8*1), X9
	MOVSHDUP_X3_X2       // Load and duplicate real elements (x2r, x2r, x1r, x1r)
	MOVSLDUP_X3_X3       // Load and duplicate imag elements (x2i, x2i, x1i, x1i)
	MOVSHDUP_X5_X4
	MOVSLDUP_X5_X5
	MOVSHDUP_X7_X6
	MOVSLDUP_X7_X7
	MOVSHDUP_X9_X8
	MOVSLDUP_X9_X9
	MULPS X1, X2         // (ai*xr, ar*xr)
	MULPS X0, X3         // (ar*xi, ai*xi)
	MULPS X11, X4
	MULPS X10, X5
	MULPS X1, X6
	MULPS X0, X7
	MULPS X11, X8
	MULPS X10, X9
	ADDSUBPS_X2_X3       // (ai*x1r+ar*x1i, ar*x1r-ai*x1i)
	ADDSUBPS_X4_X5
	ADDSUBPS_X6_X7
	ADDSUBPS_X8_X9
	MOVSD (DI), X2
	MOVSD (DI)(R9*1), X4
	ADDPS X2, X3         // Add y to a*x
	ADDPS X4, X5
	MOVSD X3, (DI)       // Write result back to dst
	MOVSD X5, (DI)(R9*1)
	LEAQ  (DI)(R9*2), DI
	MOVSD (DI), X6
	MOVSD (DI)(R9*1), X8
	ADDPS X6, X7
	ADDPS X8, X9
	MOVSD X7, (DI)       // Write result back to dst
	MOVSD X9, (DI)(R9*1)
	LEAQ  (SI)(R8*2), SI // Increment addresses
	LEAQ  (DI)(R9*2), DI
	DECQ  BX
	JNZ   axpyi_loop
	CMPQ  CX, $0
	JE    axpyi_end

axpyi_tail:
	MOVSD (SI), X3
	MOVSHDUP_X3_X2   // Load and duplicate real elements (x2r, x2r, x1r, x1r)
	MOVSLDUP_X3_X3   // Load and duplicate imag elements (x2i, x2i, x1i, x1i)
	MULPS X1, X2     // (ai*x2r, ar*x2r, ai*x1r, ar*x1r)
	MULPS X0, X3     // (ar*x2i, ai*x2i, ar*x1i, ai*x1i)
	ADDSUBPS_X2_X3   // (ai*x1r+ar*x1i, ar*x1r-ai*x1i)
	MOVSD (DI), X4
	ADDPS X4, X3     // Add y2,y1 to a*(x2,x1)
	MOVSD X3, (DI)
	ADDQ  R8, SI
	ADDQ  R9, DI
	LOOP  axpyi_tail

axpyi_end:
	RET
