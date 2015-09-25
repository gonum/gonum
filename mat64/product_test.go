// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"math/rand"

	"gopkg.in/check.v1"
)

type dims struct{ r, c int }

var productTests = []struct {
	n       int
	factors []dims
	product dims
	panics  bool
}{
	{
		n:       1,
		factors: []dims{{3, 4}},
		product: dims{3, 4},
		panics:  false,
	},
	{
		n:       1,
		factors: []dims{{2, 4}},
		product: dims{3, 4},
		panics:  true,
	},
	{
		n:       3,
		factors: []dims{{10, 30}, {30, 5}, {5, 60}},
		product: dims{10, 60},
		panics:  false,
	},
	{
		n:       3,
		factors: []dims{{100, 30}, {30, 5}, {5, 60}},
		product: dims{10, 60},
		panics:  true,
	},
	{
		n:       7,
		factors: []dims{{60, 5}, {5, 5}, {5, 4}, {4, 10}, {10, 22}, {22, 45}, {45, 10}},
		product: dims{60, 10},
		panics:  false,
	},
	{
		n:       7,
		factors: []dims{{60, 5}, {5, 5}, {5, 400}, {4, 10}, {10, 22}, {22, 45}, {45, 10}},
		product: dims{60, 10},
		panics:  true,
	},
	{
		n:       3,
		factors: []dims{{1, 1000}, {1000, 2}, {2, 2}},
		product: dims{1, 2},
		panics:  false,
	},

	// Random chains.
	{
		n:       0,
		product: dims{0, 0},
		panics:  false,
	},
	{
		n:       2,
		product: dims{60, 10},
		panics:  false,
	},
	{
		n:       3,
		product: dims{60, 10},
		panics:  false,
	},
	{
		n:       4,
		product: dims{60, 10},
		panics:  false,
	},
	{
		n:       10,
		product: dims{60, 10},
		panics:  false,
	},
}

func (s *S) TestProduct(c *check.C) {
	for _, test := range productTests {
		dimensions := test.factors
		if dimensions == nil && test.n > 0 {
			dimensions = make([]dims, test.n)
			for i := range dimensions {
				if i != 0 {
					dimensions[i].r = dimensions[i-1].c
				}
				dimensions[i].c = rand.Intn(50) + 1
			}
			dimensions[0].r = test.product.r
			dimensions[test.n-1].c = test.product.c
		}
		factors := make([]Matrix, test.n)
		for i, d := range dimensions {
			data := make([]float64, d.r*d.c)
			for i := range data {
				data[i] = rand.Float64()
			}
			factors[i] = NewDense(d.r, d.c, data)
		}

		want := &Dense{}
		if !test.panics {
			a := &Dense{}
			for i, b := range factors {
				if i == 0 {
					want.Clone(b)
					continue
				}
				a, want = want, &Dense{}
				want.Mul(a, b)
			}
		}

		got := NewDense(test.product.r, test.product.c, nil)
		panicked, message := panics(func() {
			got.Product(factors...)
		})
		if test.panics {
			if !panicked {
				c.Errorf("fail to panic with product chain dimentions: %+v result dimension: %+v",
					dimensions, test.product)
			}
			continue
		} else if panicked {
			c.Errorf("unexpected panic %q with product chain dimentions: %+v result dimension: %+v",
				message, dimensions, test.product)
			continue
		}

		if !EqualApprox(got, want, 1e-14) {
			c.Errorf("unexpected result from product chain dimensions: %+v", dimensions)
		}
	}
}
