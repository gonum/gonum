// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package goblas

import (
	"fmt"
	"runtime"

	"github.com/gonum/blas"
)

const (
	blockSize   = 50 // b x b matrix
	minParBlock = 4  // minimum number of blocks needed to go parallel
)

// Dgemm computes c := beta * C + alpha * A * B. If tA or tB is blas.Trans,
// A or B is transposed.
// m is the number of rows in A or A transpose
// n is the number of columns in B or B transpose
// k is the columns of A and rows of B
func (Blas) Dgemm(tA, tB blas.Transpose, m, n, k int, alpha float64, a []float64, lda int, b []float64, ldb int, beta float64, c []float64, ldc int) {
	var amat, bmat, cmat general
	if tA == blas.Trans {
		amat = general{
			data:   a,
			rows:   k,
			cols:   m,
			stride: lda,
		}
	} else {
		amat = general{
			data:   a,
			rows:   m,
			cols:   k,
			stride: lda,
		}
	}
	err := amat.check()
	if err != nil {
		panic(err)
	}
	if tB == blas.Trans {
		bmat = general{
			data:   b,
			rows:   n,
			cols:   k,
			stride: ldb,
		}
	} else {
		bmat = general{
			data:   b,
			rows:   k,
			cols:   n,
			stride: ldb,
		}
	}

	err = bmat.check()
	if err != nil {
		panic(err)
	}
	cmat = general{
		data:   c,
		rows:   m,
		cols:   n,
		stride: ldc,
	}
	err = cmat.check()
	if err != nil {
		panic(err)
	}
	if tA != blas.Trans && tA != blas.NoTrans {
		panic(badTranspose)
	}
	if tB != blas.Trans && tB != blas.NoTrans {
		panic(badTranspose)
	}

	// scale c
	if beta != 1 {
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				cmat.data[i*cmat.stride+j] *= beta
			}
		}
	}

	dgemmParallel(tA, tB, amat, bmat, cmat, alpha)
}

func dgemmParallel(tA, tB blas.Transpose, a, b, c general, alpha float64) {
	// dgemmParallel computes a parallel matrix multiplication by partitioning
	// a and b into sub-blocks, and updating c with the multiplication of the sub-block
	// In all cases,
	// A = [ 	A_11	A_12 ... 	A_1j
	//			A_21	A_22 ...	A_2j
	//				...
	//			A_i1	A_i2 ...	A_ij]
	//
	// and same for B. All of the submatrix sizes are blockSize*blockSize except
	// at the edges.
	// In all cases, there is one dimension for each matrix along which
	// C must be updated sequentially.
	// Cij = \sum_k Aik Bki,	(A * B)
	// Cij = \sum_k Aki Bkj,	(A^T * B)
	// Cij = \sum_k Aik Bjk,	(A * B^T)
	// Cij = \sum_k Aki Bjk,	(A^T * B^T)
	//
	// For each of these cases, this code computes all that it can in parallel,
	// and when one of the cases finishes, it sends the next possible element away.
	// As an example, for A*B, the code sends out A_i1 * B_1j for all i, j. When
	// this solution is computed, A_i2 * B_2j is sent to be computed. This
	// partitioning allows Cij to be updated in-place without race-conditions.
	// Instead of launching a goroutine for each possible concurrent computation,
	// a number of worker goroutines are created and channels are used to pass
	// available and completed cases.
	//
	// http://alexkr.com/docs/matrixmult.pdf is a good reference on matrix-matrix
	// multiplies, though this code does not copy matrices to attempt to eliminate
	// cache misses.

	aTrans := tA == blas.Trans
	bTrans := tB == blas.Trans

	maxKLen, parBlocks, totalBlocks := computeNumBlocks(a, b, aTrans, bTrans)

	if parBlocks < minParBlock {
		// The matrix multiplication is small in the dimensions where it can be
		// computed concurrently. Just do it in serial.
		dgemmSerial(tA, tB, a, b, c, alpha)
		return
	}

	// Number of workers is the minimum of the number of processors and the number
	// of possible parallel computations
	nWorkers := runtime.GOMAXPROCS(0)
	if parBlocks < nWorkers {
		nWorkers = parBlocks
	}

	// Create channels. The structure of the code avoids race conditions in updating C,
	// but the channel sends don't synchronize between the original block and the
	// updated blocks. The buffering avoids deadlock, and as a bonus allows for
	// for minimal communication blocking
	sendChan := make(chan *subMul, totalBlocks)
	ansChan := make(chan *subMul, totalBlocks)
	quit := make(chan struct{})

	// launch workers. A worker receives a submatrix, computes it, and sends
	// a message back to announce the completion.
	for i := 0; i < nWorkers; i++ {
		go func() {
			for {
				select {
				case <-quit:
					return
				case sub := <-sendChan:
					dgemmSerial(tA, tB, sub.a, sub.b, sub.c, sub.alpha)
					ansChan <- sub
				}
			}
		}()
	}

	// Send out the submatrices of the first block along k. For example, in
	// Cij = \sum Aik * Bki, send out Ai1 * B1j for all i and j.
	// The lenX variables are the number of elements along that dimension. Normally
	// that will be blockSize, but at the edges may need to be shorter
	go func() {
		lenk := blockSize
		if lenk > maxKLen {
			lenk = maxKLen
		}
		for i := 0; i < c.rows; i += blockSize {
			leni := blockSize
			if i+leni > c.rows {
				leni = c.rows - i
			}
			for j := 0; j < c.cols; j += blockSize {
				lenj := blockSize
				if j+lenj > c.cols {
					lenj = c.cols - j
				}
				// Take views of the larger a and b matrices to get the appropriate
				// subblock.
				var aSub, bSub general
				if aTrans {
					aSub = a.view(0, i, lenk, leni)
				} else {
					aSub = a.view(i, 0, leni, lenk)
				}
				if bTrans {
					bSub = b.view(j, 0, lenj, lenk)
				} else {
					bSub = b.view(0, j, lenk, lenj)
				}
				sendChan <- &subMul{
					i:     i,
					j:     j,
					k:     0,
					a:     aSub,
					b:     bSub,
					c:     c.view(i, j, leni, lenj),
					alpha: alpha,
				}
			}
		}
	}()

	// Read in the cases as they come in. Send the next k block along that same
	// {i,j} pair to be computed if there is a block left to compute.
	var totalReceived int
	for sub := range ansChan {
		totalReceived++
		if totalReceived == totalBlocks {
			close(quit)
			return
		}
		// If we already computed the final k, don't need to send anything else
		// along this pair.
		if sub.k+blockSize >= maxKLen {
			continue
		}
		// Otherwise, get the block of k along the {i,j} pair. Since i and j are
		// constant, the view of C stays the same.
		k := sub.k + blockSize
		lenk := blockSize
		if lenk > maxKLen {
			lenk = maxKLen - k
		}
		var aSub, bSub general
		if aTrans {
			aSub = a.view(k, sub.i, lenk, sub.a.cols)
		} else {
			aSub = a.view(sub.i, k, sub.a.rows, lenk)
		}
		if bTrans {
			bSub = b.view(sub.j, k, sub.b.rows, lenk)
		} else {
			bSub = b.view(k, sub.j, lenk, sub.b.cols)
		}
		sub.a = aSub
		sub.b = bSub
		sub.k = k
		sendChan <- sub
	}
	close(quit)
	return
}

type subMul struct {
	i, j, k int     // index of block
	a, b, c general // submatrices of the global a,b,c
	alpha   float64
}

// computeNumBlocks says how many blocks there are to compute. maxKLen says the length of the
// k dimension, parBlocks is the number of blocks that could be computed in parallel
// (the submatrices in i and j). totalBlocks is the full number of blocks.
func computeNumBlocks(a, b general, aTrans, bTrans bool) (maxKLen, parBlocks, totalBlocks int) {
	aRowBlocks := a.rows / blockSize
	if a.rows%blockSize != 0 {
		aRowBlocks++
	}
	aColBlocks := a.cols / blockSize
	if a.cols%blockSize != 0 {
		aColBlocks++
	}
	bRowBlocks := b.rows / blockSize
	if b.rows%blockSize != 0 {
		bRowBlocks++
	}
	bColBlocks := b.cols / blockSize
	if b.cols%blockSize != 0 {
		bColBlocks++
	}

	switch {
	case !aTrans && !bTrans:
		// Cij = \sum_k Aik Bki
		maxKLen = a.cols
		parBlocks = aRowBlocks * bColBlocks
		totalBlocks = parBlocks * aColBlocks
	case aTrans && !bTrans:
		// Cij = \sum_k Aki Bkj
		maxKLen = a.rows
		parBlocks = aColBlocks * bColBlocks
		totalBlocks = parBlocks * aRowBlocks
	case !aTrans && bTrans:
		// Cij = \sum_k Aik Bjk
		maxKLen = a.cols
		parBlocks = aRowBlocks * bRowBlocks
		totalBlocks = parBlocks * aColBlocks
	case aTrans && bTrans:
		// Cij = \sum_k Aki Bjk
		maxKLen = a.rows
		parBlocks = aColBlocks * bRowBlocks
		totalBlocks = parBlocks * aRowBlocks
	}
	return
}

// dgemmSerial is serial matrix multiply
func dgemmSerial(tA, tB blas.Transpose, a, b, c general, alpha float64) {
	switch {
	case tA == blas.NoTrans && tB == blas.NoTrans:
		dgemmSerialNotNot(a, b, c, alpha)
		return
	case tA == blas.Trans && tB == blas.NoTrans:
		dgemmSerialTransNot(a, b, c, alpha)
		return
	case tA == blas.NoTrans && tB == blas.Trans:
		dgemmSerialNotTrans(a, b, c, alpha)
		return
	case tA == blas.Trans && tB == blas.Trans:
		dgemmSerialTransTrans(a, b, c, alpha)
		return
	default:
		panic("unreachable")
	}
}

// dgemmSerial where neither a nor b are transposed
func dgemmSerialNotNot(a, b, c general, alpha float64) {
	if debug {
		if a.cols != b.rows {
			panic("inner dimension mismatch")
		}
		if a.rows != c.rows {
			panic("outer dimension mismatch")
		}
		if b.cols != c.cols {
			panic("outer dimension mismatch")
		}
	}
	for i := 0; i < a.rows; i++ {
		for l := 0; l < a.cols; l++ {
			tmp := alpha * a.at(i, l)
			if tmp != 0 {
				for j := 0; j < b.cols; j++ {
					c.data[i*c.stride+j] += tmp * b.at(l, j)
				}
			}
		}
	}
}

// dgemmSerial where neither a is transposed and b is not
func dgemmSerialTransNot(a, b, c general, alpha float64) {
	if debug {
		if a.rows != b.rows {
			fmt.Println(a.rows, b.rows)
			panic("inner dimension mismatch")
		}
		if a.cols != c.rows {
			panic("outer dimension mismatch")
		}
		if b.cols != c.cols {
			panic("outer dimension mismatch")
		}
	}
	for i := 0; i < a.cols; i++ {
		for l := 0; l < a.rows; l++ {
			tmp := alpha * a.at(l, i)
			if tmp != 0 {
				for j := 0; j < b.cols; j++ {
					c.data[i*c.stride+j] += tmp * b.at(l, j)
				}
			}
		}
	}
}

// dgemmSerial where neither a is not transposed and b is
func dgemmSerialNotTrans(a, b, c general, alpha float64) {
	if debug {
		if a.cols != b.cols {
			panic("inner dimension mismatch")
		}
		if a.rows != c.rows {
			panic("outer dimension mismatch")
		}
		if b.rows != c.cols {
			panic("outer dimension mismatch")
		}
	}
	for i := 0; i < a.rows; i++ {
		for j := 0; j < b.rows; j++ {
			var tmp float64
			for l := 0; l < a.cols; l++ {
				tmp += a.at(i, l) * b.at(j, l)
			}
			c.data[i*c.stride+j] += alpha * tmp
		}
	}
}

// dgemmSerial where both are transposed
func dgemmSerialTransTrans(a, b, c general, alpha float64) {
	if debug {
		if a.rows != b.cols {
			panic("inner dimension mismatch")
		}
		if a.cols != c.rows {
			panic("outer dimension mismatch")
		}
		if b.rows != c.cols {
			panic("outer dimension mismatch")
		}
	}
	for i := 0; i < a.cols; i++ {
		for l := 0; l < a.rows; l++ {
			v := a.at(l, i)
			if v != 0 {
				tmp := alpha * v
				for j := 0; j < b.rows; j++ {
					c.data[i*c.stride+j] += tmp * b.at(j, l)
				}
			}
		}
	}
}
