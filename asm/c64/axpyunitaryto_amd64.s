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

// func AxpyUnitaryTo(dst []complex64, alpha complex64, x, y []complex64)
TEXT ·AxpyUnitaryTo(SB), NOSPLIT, $0
	MOVQ    dst_base+0(FP), DI // 4
	MOVQ    x_base+32(FP), SI
	MOVQ    y_base+56(FP), DX
	MOVQ    x_len+40(FP), CX
	CMPQ    y_len+64(FP), CX   // 4
	CMOVQLE y_len+64(FP), CX   // 5
	CMPQ    dst_len+8(FP), CX
	CMOVQLE dst_len+8(FP), CX
	CMPQ    CX, $0
	JE      caxy_end
	MOVSD   alpha+24(FP), X0   // (0,0,ar,ai)
	SHUFPD  $0, X0, X0         // (ar,ai,ar,ai)
	MOVAPS  X0, X1
	SHUFPS  $0x11, X1, X1      // (ai,ar,ai,ar)
	XORQ    AX, AX
	MOVQ    DX, BX
	ANDQ    $15, BX            // Align on 16-byte boundary for ADDPS
	JZ      caxy_no_trim

	MOVSD (SI)(AX*8), X3
	MOVSHDUP_X3_X2
	MOVSLDUP_X3_X3
	MULPS X1, X2
	MULPS X0, X3
	ADDSUBPS_X2_X3
	MOVSD (DX)(AX*8), X4
	ADDPS X4, X3
	MOVSD X3, (DI)(AX*8)
	INCQ  AX
	DECQ  CX
	JZ    caxy_end

caxy_no_trim:
	MOVAPS X0, X10
	MOVAPS X1, X11
	MOVQ   CX, BX
	ANDQ   $7, CX
	SHRQ   $3, BX
	JZ     caxy_tail

caxy_loop:
	MOVUPS (SI)(AX*8), X3
	MOVUPS 16(SI)(AX*8), X5
	MOVUPS 32(SI)(AX*8), X7
	MOVUPS 48(SI)(AX*8), X9
	MOVSHDUP_X3_X2          // Load and duplicate real elements (x2r, x2r, x1r, x1r)
	MOVSLDUP_X3_X3          // Load and duplicate imag elements (x2i, x2i, x1i, x1i)
	MOVSHDUP_X5_X4
	MOVSLDUP_X5_X5
	MOVSHDUP_X7_X6
	MOVSLDUP_X7_X7
	MOVSHDUP_X9_X8
	MOVSLDUP_X9_X9
	MULPS  X1, X2           // (ai*x2r, ar*x2r, ai*x1r, ar*x1r)
	MULPS  X0, X3           // (ar*x2i, ai*x2i, ar*x1i, ai*x1i)
	MULPS  X11, X4
	MULPS  X10, X5
	MULPS  X1, X6
	MULPS  X0, X7
	MULPS  X11, X8
	MULPS  X10, X9
	ADDSUBPS_X2_X3          // (ai*x2r+ar*x2i, ar*x2r-ai*x2i, ai*x1r+ar*x1i, ar*x1r-ai*x1i)
	ADDSUBPS_X4_X5
	ADDSUBPS_X6_X7
	ADDSUBPS_X8_X9
	ADDPS  (DX)(AX*8), X3   // Add y2,y1 to a*(x2,x1)
	ADDPS  16(DX)(AX*8), X5
	ADDPS  32(DX)(AX*8), X7
	ADDPS  48(DX)(AX*8), X9
	MOVUPS X3, (DI)(AX*8)   // Write result back to dst2,dst1
	MOVUPS X5, 16(DI)(AX*8)
	MOVUPS X7, 32(DI)(AX*8)
	MOVUPS X9, 48(DI)(AX*8)
	ADDQ   $8, AX
	DECQ   BX
	JNZ    caxy_loop
	CMPQ   CX, $0
	JE     caxy_end

caxy_tail: // Same calculation, but read in values to avoid trampling memory
	MOVSD (SI)(AX*8), X3
	MOVSHDUP_X3_X2       // Load and duplicate real elements (x2r, x2r, x1r, x1r)
	MOVSLDUP_X3_X3       // Load and duplicate imag elements (x2i, x2i, x1i, x1i)
	MULPS X1, X2         // (ai*x2r, ar*x2r, ai*x1r, ar*x1r)
	MULPS X0, X3         // (ar*x2i, ai*x2i, ar*x1i, ai*x1i)
	ADDSUBPS_X2_X3       // Add y2,y1 to a*(x2,x1)
	MOVSD (DX)(AX*8), X4
	ADDPS X4, X3
	MOVSD X3, (DI)(AX*8)
	INCQ  AX
	LOOP  caxy_tail

caxy_end:
	RET
