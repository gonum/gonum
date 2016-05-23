// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

//+build !noasm,!appengine

#include "textflag.h"

// func AxpyIncTo(dst []complex64, incDst, idst uintptr, alpha complex64, x, y []complex64, n, incX, incY, ix, iy uintptr)
TEXT ·AxpyIncTo(SB), NOSPLIT, $0
	MOVQ   dst_base+0(FP), R8
	MOVQ   x_base+48(FP), SI
	MOVQ   y_base+72(FP), DI
	MOVQ   n+96(FP), DX
	XORQ   AX, AX             // i = 0
	CMPQ   DX, $0             // if n==0, return
	JE     axpyi_end
	MOVQ   ix+120(FP), CX     // Load the first indicies
	MOVQ   iy+128(FP), BX
	MOVQ   idst+32(FP), R9
	LEAQ   (SI)(CX*8), SI     // Calculate addrress of first indicies
	LEAQ   (DI)(BX*8), DI
	LEAQ   (R8)(R9*8), R8
	MOVQ   incX+104(FP), CX   // Incrementors*8 for easy iteration (ADDQ)
	SHLQ   $3, CX
	MOVQ   incY+112(FP), BX
	SHLQ   $3, BX
	MOVQ   incDst+24(FP), R9
	SHLQ   $3, R9
	MOVSD  alpha+40(FP), X0   // XO:(0,0,ar,ai)
	SHUFPD $0, X0, X0         // X0:(ar,ai,ar,ai)
	MOVAPS X0, X1
	SHUFPS $0x11, X1, X1      // X1:(ai,ar,ai,ar)

axpyi_loop:
	// MOVSHDUP (SI), X2	// Load and duplicate real elements (x1r, x1r)
	// MOVSLDUP (SI), X3	// Load and duplicate imag elements (x1i, x1i)
	BYTE $0xF3; BYTE $0x0F; BYTE $0x16; BYTE $0x16
	BYTE $0xF3; BYTE $0x0F; BYTE $0x12; BYTE $0x1E

	MULPS X1, X2 // (ai*x1r, ar*x1r)
	MULPS X0, X3 // (ar*x1i, ai*x1i)

	// ADDSUBPS X2, X3  	// (ai*x1r+ar*x1i, ar*x1r-ai*x1i)
	BYTE $0xF2; BYTE $0x0F; BYTE $0xD0; BYTE $0xDA

	ADDPS (DI), X3   // Add y to a*x
	MOVSD X3, (R8)   // Write result back to dst
	ADDQ  CX, SI     // Increment addresses
	ADDQ  BX, DI
	ADDQ  R9, R8
	INCQ  AX         // i++
	CMPQ  AX, DX     // while i < n
	JL    axpyi_loop

axpyi_end:
	RET
