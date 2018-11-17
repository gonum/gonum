// Copyright ©2018 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package quat

func lift(v float64) Quat {
	return Quat{Real: v}
}

func split(q Quat) (float64, Quat) {
	return q.Real, Quat{Imag: q.Imag, Jmag: q.Jmag, Kmag: q.Kmag}
}

func join(w float64, uv Quat) Quat {
	uv.Real = w
	return uv
}

func unit(q Quat) Quat {
	return Scale(1/Abs(q), q)
}
