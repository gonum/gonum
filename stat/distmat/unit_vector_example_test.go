// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package distmat_test

import (
	"fmt"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
	"gonum.org/v1/gonum/stat/distmat"
)

// ExampleUnitVector uses the UnitVector distribution to take
// a random walk in n-space.  At the end it computes how far
// from the origin the walk finished.
func ExampleUnitVector() {
	src := rand.NewSource(1)
	rnd := rand.New(src)
	dist := distmat.NewUnitVector(src)

	// dimension of the space to walk in
	nDim := 1 + rnd.Intn(100)
	// current position
	position := mat.NewVecDense(nDim, nil)
	// each step goes here
	step := mat.NewVecDense(nDim, nil)

	// how many steps to take
	nSteps := 1 + rnd.Intn(100)
	for i := 0; i < nSteps; i++ {
		// draw and then take a unit step
		dist.UnitVecTo(step)
		position.AddVec(position, step)
	}

	// compute distance from the origin
	distance := mat.Norm(position, 2)
	fmt.Printf("took %d steps in %d-space, walked %1.1f in total", nSteps, nDim, distance)
	// Output:
	// took 22 steps in 52-space, walked 5.3 in total
}
