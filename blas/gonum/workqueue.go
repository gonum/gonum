// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gonum

import (
	"sync/atomic"
)

// blockWorkQueue implements parallel iterator over two dimensions with blockSize.
//
// The implementation corresponds to this double loop:
//
//   for i := 0; i < m; i += blockSize {
//      for j := 0; j < n; j += blockSize {
//         next(i, j)
//      }
//   }
type blockWorkQueue struct {
	head int64

	total int
	mod   int
}

// Reset resets the work queue with the parameters.
func (q *blockWorkQueue) Reset(m, n int) {
	q.head = 0

	qm := blocks(m, blockSize)
	qn := blocks(n, blockSize)
	q.total = qm * qn
	q.mod = qn
}

// Next returns work items until everything has been exhausted.
func (q *blockWorkQueue) Next() (i, j int, ok bool) {
	w := int(atomic.AddInt64(&q.head, 1)) - 1
	i = (w / q.mod) * blockSize
	j = (w % q.mod) * blockSize
	return i, j, w < q.total
}
