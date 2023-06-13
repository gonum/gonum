// Copyright ©2015 The Gonum Authors. All rights reserved.
// Use of this code is governed by a BSD-style
// license that can be found in the LICENSE file

package sampleuv

import (
	"golang.org/x/exp/rand"
)

// Weighted provides sampling without replacement from a collection of items with
// non-uniform probability.
type Weighted struct {
	weights []float64
	// heap is a weight heap.
	//
	// It keeps a heap-organised sum of remaining
	// index weights that are available to be taken
	// from.
	//
	// Each element holds the sum of weights for
	// the corresponding index, plus the sum of
	// its children's weights; the children of
	// an element i can be found at positions
	// 2*(i+1)-1 and 2*(i+1). The root of the
	// weight heap is at element 0.
	//
	// See comments in container/heap for an
	// explanation of the layout of a heap.
	heap []float64
	rnd  *rand.Rand
}

// NewWeighted returns a Weighted for the weights w. If src is nil, rand.Rand is
// used as the random number generator.
//
// Note that sampling from weights with a high variance or overall low absolute
// value sum may result in problems with numerical stability.
func NewWeighted(w []float64, src rand.Source) Weighted {
	s := Weighted{
		weights: make([]float64, len(w)),
		heap:    make([]float64, len(w)),
	}
	if src != nil {
		s.rnd = rand.New(src)
	}
	s.ReweightAll(w)
	return s
}

// Len returns the number of items held by the Weighted, including items
// already taken.
func (s Weighted) Len() int { return len(s.weights) }

func (s Weighted) findIndex(r float64) int {
	i := 0

	for {
		r -= s.weights[i]
		if r < 0 {
			break // Fall within item i.
		}

		ni := i<<1 + 1 // Move to left child.
		// Left node should exist, because r is non-negative, but there are could be floating point errors,
		// so we check index explicitly.
		if ni >= len(s.heap) {
			break
		}

		i = ni

		d := s.heap[i]
		if r >= d {
			// If enough r to pass left child try to move to the right child.
			r -= d
			ni := i + 1

			if ni >= len(s.heap) {
				break
			}

			i = ni
		}
	}

	return i
}

// Take returns an index from the Weighted with probability proportional
// to the weight of the item. The weight of the item is then set to zero.
// Take returns false if there are no items remaining.
func (s Weighted) Take() (idx int, ok bool) {
	if s.heap[0] == 0 {
		return -1, false
	}

	var r float64
	if s.rnd == nil {
		r = rand.Float64()
	} else {
		r = s.rnd.Float64()
	}

	i := s.findIndex(s.heap[0] * r)

	s.Reweight(i, 0)

	return i, true
}

// Reweight sets the weight of item idx to w.
func (s Weighted) Reweight(idx int, w float64) {
	s.weights[idx] = w

	for {
		w = s.weights[idx]
		// We want to keep heap state here like after reset call.
		// Therefore, we sum weights in such order, because floating point sum not associative.

		ri := idx*2 + 2
		if ri < len(s.heap) {
			w += s.heap[ri]
		}

		li := ri - 1
		if li < len(s.heap) {
			w += s.heap[li]
		}

		s.heap[idx] = w

		if idx == 0 {
			break
		}

		idx = (idx - 1) / 2
	}
}

// ReweightAll sets the weight of all items in the Weighted. ReweightAll
// panics if len(w) != s.Len.
func (s Weighted) ReweightAll(w []float64) {
	if len(w) != s.Len() {
		panic("floats: length of the slices do not match")
	}
	copy(s.weights, w)
	s.reset()
}

func (s Weighted) reset() {
	copy(s.heap, s.weights)
	for i := len(s.heap) - 1; i > 0; i-- {
		// Sometimes 1-based counting makes sense.
		s.heap[((i+1)>>1)-1] += s.heap[i]
	}
}
