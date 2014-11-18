// Copyright ©2014 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"sync"
//	"runtime"
	
	"github.com/gonum/matrix/mat64"
)


type covMatSlice struct {
	i, j int
	x, y []float64
}

// CovarianceMatrix calculates a covariance matrix (also known as a
// variance-covariance matrix) from a matrix of data, using a two-pass
// algorithm.  It will have better performance if a BLAS engine is
// registered in gonum/matrix/mat64.
//
// The matrix returned will be symmetric, square, and positive-semidefinite.  
func CovarianceMatrix(x mat64.Matrix) *mat64.Dense {

	// matrix version of the two pass algorithm.  This doesn't use
	// the correction found in the Covariance and Variance functions.
	r, c := x.Dims()
	
	if x, ok := x.(mat64.Vectorer); ok {
		cols := make([][]float64, c)
		// perform the covariance or variance as required
		blockSize := 1024
		if blockSize > c {
			blockSize = c
		}
		var wg sync.WaitGroup
		wg.Add(c)
		for j := 0; j < c; j++ {
			go func(j int) {
				// pull the columns out and subtract the means
				cols[j] = make([]float64, r)
				x.Col(cols[j], j)
				mean := Mean(cols[j], nil)
				for i := range cols[j] {
					cols[j][i] -= mean
				}
				wg.Done()
			}(j)
		}
		wg.Wait()
		
		colCh := make(chan covMatSlice, blockSize)

		wg.Add(blockSize)
		m := mat64.NewDense(c, c, nil)
		for i := 0; i < blockSize; i++ {
			go func(in <-chan covMatSlice) {
				for {
					xy, more := <-in
					if !more {
						wg.Done()
						return
					}

					if xy.i == xy.j {
						m.Set(xy.i, xy.j, centeredVariance(xy.x))
						continue
					}
					v := centeredCovariance(xy.x, xy.y)
					m.Set(xy.i, xy.j, v)
					m.Set(xy.j, xy.i, v)
					}
			}(colCh)
		}
		go func(out chan<- covMatSlice) {
			for i := 0; i < c; i++ {
				for j := 0; j <= i; j++ {
					out <- covMatSlice{
						i: i,
						j: j,
						x: cols[i],
						y: cols[j],
					}
				}
			}
			close(out)
		}(colCh)
		// create the output matrix
		wg.Wait()
		return m
	}
	// determine the mean of each of the columns
	b := ones(1, r)
	b.Mul(b, x)
	b.Scale(1/float64(r), b)
	mu := b.RowView(0)

	// subtract the mean from the data
	xc := mat64.DenseCopyOf(x)
	for i := 0; i < r; i++ {
		rv := xc.RowView(i)
		for j, mean := range mu {
			rv[j] -= mean
		}
	}

	// todo: avoid matrix copy?
	var xt mat64.Dense
	xt.TCopy(xc)

	// It would be nice if we could indicate that this was a symmetric
	// matrix.
	var ss mat64.Dense
	ss.Mul(&xt, xc)
	ss.Scale(1/float64(r-1), &ss)
	return &ss
}

// ones is a matrix of all ones.
func ones(r, c int) *mat64.Dense {
	x := make([]float64, r*c)
	for i := range x {
		x[i] = 1
	}
	return mat64.NewDense(r, c, x)
}

func centeredVariance(x []float64) float64 {
	var ss float64
	for _, xv := range x {
		ss += xv * xv
	}
	return ss / float64(len(x)-1)
}

func centeredCovariance(x, y []float64) float64 {
	var ss float64
	for i, xv := range x {
		ss += xv * y[i]
	}
	return ss / float64(len(x)-1)
}
