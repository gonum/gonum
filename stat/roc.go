// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat

import (
	"math"
	"sort"
)

// ROC returns paired false positive rate (FPR) and true positive rate
// (TPR) values corresponding to cutoff points on the receiver operator
// characteristic (ROC) curve obtained when y is treated as a binary
// classifier for classes with weights. The cutoff thresholds used to
// calculate the ROC are returned in thresh such that tpr[i] and fpr[i]
// are the true and false positive rates for y >= thresh[i].
//
// The input y and cutoffs must be sorted, and values in y must correspond
// to values in classes and weights. SortWeightedLabeled can be used to
// sort y together with classes and weights.
//
// For a given cutoff value, observations corresponding to entries in y
// greater than the cutoff value are classified as true, while those
// less than or equal to the cutoff value are classified as false. These
// assigned class labels are compared with the true values in the classes
// slice and used to calculate the FPR and TPR.
//
// If weights is nil, all weights are treated as 1. If weights is not nil
// it must have the same length as y and classes, otherwise ROC will panic.
//
// If cutoffs is nil or empty, all possible cutoffs are calculated,
// resulting in fpr and tpr having length one greater than the number of
// unique values in y. Otherwise fpr and tpr will be returned with the
// same length as cutoffs. floats.Span can be used to generate equally
// spaced cutoffs.
//
// More details about ROC curves are available at
// https://en.wikipedia.org/wiki/Receiver_operating_characteristic
func ROC(cutoffs, y []float64, classes []bool, weights []float64) (tpr, fpr, thresh []float64) {
	if len(y) != len(classes) {
		panic("stat: slice length mismatch")
	}
	if weights != nil && len(y) != len(weights) {
		panic("stat: slice length mismatch")
	}
	if !sort.Float64sAreSorted(y) {
		panic("stat: input must be sorted ascending")
	}
	if !sort.Float64sAreSorted(cutoffs) {
		panic("stat: cutoff values must be sorted ascending")
	}
	if len(y) == 0 {
		return nil, nil, nil
	}
	if len(cutoffs) == 0 {
		if cutoffs == nil || cap(cutoffs) < len(y)+1 {
			cutoffs = make([]float64, len(y)+1)
		} else {
			cutoffs = cutoffs[:len(y)+1]
		}
		// Choose all possible cutoffs for unique values in y.
		bin := 0
		cutoffs[bin] = y[0]
		for i, u := range y[1:] {
			if u == y[i] {
				continue
			}
			bin++
			cutoffs[bin] = u
		}
		cutoffs[bin+1] = math.Inf(1)
		cutoffs = cutoffs[:bin+2]
	} else {
		// Don't mutate the provided cutoffs.
		tmp := cutoffs
		cutoffs = make([]float64, len(cutoffs))
		copy(cutoffs, tmp)
	}

	tpr = make([]float64, len(cutoffs))
	fpr = make([]float64, len(cutoffs))
	var bin int
	var nPos, nNeg float64
	for i, u := range classes {
		// Update the bin until it matches the next y value
		// skipping empty bins.
		for bin < len(cutoffs)-1 && y[i] >= cutoffs[bin] {
			bin++
			tpr[bin] = tpr[bin-1]
			fpr[bin] = fpr[bin-1]
		}
		posWeight, negWeight := 1.0, 0.0
		if weights != nil {
			posWeight = weights[i]
		}
		if !u {
			posWeight, negWeight = negWeight, posWeight
		}
		nPos += posWeight
		nNeg += negWeight
		// Count false negatives (in tpr) and true negatives (in fpr).
		if y[i] < cutoffs[bin] {
			tpr[bin] += posWeight
			fpr[bin] += negWeight
		}
	}

	invNeg := 1 / nNeg
	invPos := 1 / nPos
	// Convert negative counts to TPR and FPR.
	// Bins beyond the maximum value in y are skipped
	// leaving these fpr and tpr elements as zero.
	for i := range tpr[:bin+1] {
		// Prevent fused float operations by
		// making explicit float64 conversions.
		tpr[i] = 1 - float64(tpr[i]*invPos)
		fpr[i] = 1 - float64(fpr[i]*invNeg)
	}
	for i, j := 0, len(tpr)-1; i < j; i, j = i+1, j-1 {
		tpr[i], tpr[j] = tpr[j], tpr[i]
		fpr[i], fpr[j] = fpr[j], fpr[i]
	}
	for i, j := 0, len(cutoffs)-1; i < j; i, j = i+1, j-1 {
		cutoffs[i], cutoffs[j] = cutoffs[j], cutoffs[i]
	}

	return tpr, fpr, cutoffs
}
