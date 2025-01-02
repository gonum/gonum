// Copyright ©2020 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package uid implements unique ID provision for graphs.
package uid

import (
	"math"
	"math/rand/v2"
	"testing"
)

func TestSetChurn(t *testing.T) {
	rnd := rand.New(rand.NewPCG(1, 1))

	set := NewSet()

	// Iterate over a number of ID allocations,
	// occasionally deleting IDs from the store.
	seen := make(map[int64]bool)
	for k := 0; k < 2; k++ {
		for i := 0; i < 1e4; i++ {
			id := set.NewID()
			if seen[id] {
				t.Fatalf("NewID returned already used ID")
			}
			set.Use(id)
			seen[id] = true
			if rnd.Float64() < 0.01 {
				j := rnd.IntN(10)
				for id := range seen {
					set.Release(id)
					delete(seen, id)
					j--
					if j <= 0 {
						break
					}
				}
			}
		}

		// Kick the set into scavenging mode.
		set.Use(math.MaxInt64)
	}
}
