// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmv

const (
	badQuantile      = "distmv: quantile not between 0 and 1"
	badOutputLen     = "distmv: output slice is not nil or the correct length"
	badInputLength   = "distmv: input slice length mismatch"
	badSizeMismatch  = "distmv: size mismatch"
	badZeroDimension = "distmv: zero dimensional input"
	nonPosDimension  = "distmv: non-positive dimension input"
)

const logTwoPi = 1.8378770664093454835606594728112352797227949472755668

// reuseAs returns a slice of length n. If len(dst) is n, dst is returned,
// otherwise dst must be nil or reuseAs will panic.
func reuseAs(dst []float64, n int) []float64 {
	if dst == nil {
		dst = make([]float64, n)
	}
	if len(dst) != n {
		panic(badOutputLen)
	}
	return dst
}
