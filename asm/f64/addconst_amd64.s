// Copyright ©2016 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

#include "textflag.h"

// func Daddconst(alpha float64, x []float64)
TEXT ·AddConst(SB), NOSPLIT, $0
	MOVQ    x_base+8(FP), SI
	MOVQ    x_len+16(FP), CX
	CMPQ    CX, $0
	JE      ac_end
	MOVSD   alpha+0(FP), X0
	SHUFPD  $0, X0, X0
	SUBQ    $4, CX
	JL      ac_tail
ac_loop:  
	MOVUPS  (SI), X1
	MOVUPS  16(SI), X2
	ADDPD   X0, X1
	ADDPD   X0, X2
	MOVUPS  X1, (SI)
	MOVUPS  X2, 16(SI)
	ADDQ    $32, SI
	SUBQ    $4, CX
	JGE     ac_loop
	ADDQ    $4, CX
	JE      ac_end
ac_tail:  
	MOVSD   (SI), X1
	ADDSD   X0, X1
	MOVSD   X1, (SI)
	ADDQ    $8, SI
	SUBQ    $1, CX
	JG      ac_tail
ac_end:
	RET
