// Copyright ©2016 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stat_test

import (
	"fmt"
	"math"

	"gonum.org/v1/gonum/floats"
	"gonum.org/v1/gonum/integrate"
	"gonum.org/v1/gonum/stat"
)

func ExampleROC_weighted() {
	y := []float64{0, 3, 5, 6, 7.5, 8}
	classes := []bool{false, true, false, true, true, true}
	weights := []float64{4, 1, 6, 3, 2, 2}

	tpr, fpr, _ := stat.ROC(nil, y, classes, weights)
	fmt.Printf("true  positive rate: %v\n", tpr)
	fmt.Printf("false positive rate: %v\n", fpr)

	// Output:
	// true  positive rate: [0 0.25 0.5 0.875 0.875 1 1]
	// false positive rate: [0 0 0 0 0.6 0.6 1]
}

func ExampleROC_unweighted() {
	y := []float64{0, 3, 5, 6, 7.5, 8}
	classes := []bool{false, true, false, true, true, true}

	tpr, fpr, _ := stat.ROC(nil, y, classes, nil)
	fmt.Printf("true  positive rate: %v\n", tpr)
	fmt.Printf("false positive rate: %v\n", fpr)

	// Output:
	// true  positive rate: [0 0.25 0.5 0.75 0.75 1 1]
	// false positive rate: [0 0 0 0 0.5 0.5 1]
}

func ExampleROC_threshold() {
	y := []float64{0.1, 0.4, 0.35, 0.8}
	classes := []bool{false, false, true, true}
	stat.SortWeightedLabeled(y, classes, nil)

	tpr, fpr, thresh := stat.ROC(nil, y, classes, nil)
	fmt.Printf("true  positive rate: %v\n", tpr)
	fmt.Printf("false positive rate: %v\n", fpr)
	fmt.Printf("cutoff thresholds: %v\n", thresh)

	// Output:
	// true  positive rate: [0 0.5 0.5 1 1]
	// false positive rate: [0 0 0.5 0.5 1]
	// cutoff thresholds: [+Inf 0.8 0.4 0.35 0.1]
}

func ExampleROC_unsorted() {
	y := []float64{8, 7.5, 6, 5, 3, 0}
	classes := []bool{true, true, true, false, true, false}
	weights := []float64{2, 2, 3, 6, 1, 4}

	stat.SortWeightedLabeled(y, classes, weights)

	tpr, fpr, _ := stat.ROC(nil, y, classes, weights)
	fmt.Printf("true  positive rate: %v\n", tpr)
	fmt.Printf("false positive rate: %v\n", fpr)

	// Output:
	// true  positive rate: [0 0.25 0.5 0.875 0.875 1 1]
	// false positive rate: [0 0 0 0 0.6 0.6 1]
}

func ExampleROC_knownCutoffs() {
	y := []float64{8, 7.5, 6, 5, 3, 0}
	classes := []bool{true, true, true, false, true, false}
	weights := []float64{2, 2, 3, 6, 1, 4}
	cutoffs := []float64{-1, 3, 4}

	stat.SortWeightedLabeled(y, classes, weights)

	tpr, fpr, _ := stat.ROC(cutoffs, y, classes, weights)
	fmt.Printf("true  positive rate: %v\n", tpr)
	fmt.Printf("false positive rate: %v\n", fpr)

	// Output:
	// true  positive rate: [0.875 1 1]
	// false positive rate: [0.6 0.6 1]
}

func ExampleROC_equallySpacedCutoffs() {
	y := []float64{8, 7.5, 6, 5, 3, 0}
	classes := []bool{true, true, true, false, true, true}
	weights := []float64{2, 2, 3, 6, 1, 4}
	n := 9

	stat.SortWeightedLabeled(y, classes, weights)
	cutoffs := make([]float64, n)
	floats.Span(cutoffs, math.Nextafter(y[0], y[0]-1), y[len(y)-1])

	tpr, fpr, _ := stat.ROC(cutoffs, y, classes, weights)
	fmt.Printf("true  positive rate: %.3v\n", tpr)
	fmt.Printf("false positive rate: %.3v\n", fpr)

	// Output:
	// true  positive rate: [0.167 0.333 0.583 0.583 0.583 0.667 0.667 0.667 1]
	// false positive rate: [0 0 0 1 1 1 1 1 1]
}

func ExampleROC_aUC_unweighted() {
	y := []float64{0.1, 0.35, 0.4, 0.8}
	classes := []bool{true, false, true, false}

	tpr, fpr, _ := stat.ROC(nil, y, classes, nil)

	// Compute Area Under Curve.
	auc := integrate.Trapezoidal(fpr, tpr)
	fmt.Printf("true  positive rate: %v\n", tpr)
	fmt.Printf("false positive rate: %v\n", fpr)
	fmt.Printf("auc: %v\n", auc)

	// Output:
	// true  positive rate: [0 0 0.5 0.5 1]
	// false positive rate: [0 0.5 0.5 1 1]
	// auc: 0.25
}

func ExampleROC_aUC_weighted() {
	y := []float64{0.1, 0.35, 0.4, 0.8}
	classes := []bool{true, false, true, false}
	weights := []float64{1, 2, 2, 1}

	tpr, fpr, _ := stat.ROC(nil, y, classes, weights)

	// Compute Area Under Curve.
	auc := integrate.Trapezoidal(fpr, tpr)
	fmt.Printf("auc: %f\n", auc)

	// Output:
	// auc: 0.444444

}

func ExampleTOC() {
	classes := []bool{
		false, false, false, false, false, false,
		false, false, false, false, false, false,
		false, false, true, true, true, true,
		true, true, true, false, false, true,
		false, true, false, false, true, false,
	}

	min, ntp, max := stat.TOC(classes, nil)
	fmt.Printf("minimum bound: %v\n", min)
	fmt.Printf("TOC:           %v\n", ntp)
	fmt.Printf("maximum bound: %v\n", max)

	// Output:
	// minimum bound: [0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 1 2 3 4 5 6 7 8 9 10]
	// TOC:           [0 0 1 1 1 2 2 3 3 3 4 5 6 7 8 9 10 10 10 10 10 10 10 10 10 10 10 10 10 10 10]
	// maximum bound: [0 1 2 3 4 5 6 7 8 9 10 10 10 10 10 10 10 10 10 10 10 10 10 10 10 10 10 10 10 10 10]
}

func ExampleTOC_unsorted() {
	y := []float64{8, 7.5, 6, 5, 3, 0}
	classes := []bool{true, false, true, false, false, false}
	weights := []float64{4, 1, 6, 3, 2, 2}

	stat.SortWeightedLabeled(y, classes, weights)

	min, ntp, max := stat.TOC(classes, weights)
	fmt.Printf("minimum bound: %v\n", min)
	fmt.Printf("TOC:           %v\n", ntp)
	fmt.Printf("maximum bound: %v\n", max)

	// Output:
	// minimum bound: [0 0 0 3 6 8 10]
	// TOC:           [0 4 4 10 10 10 10]
	// maximum bound: [0 4 5 10 10 10 10]
}

func ExampleTOC_aUC_unweighted() {
	classes := []bool{true, false, true, false}

	_, ntp, _ := stat.TOC(classes, nil)
	pos := ntp[len(ntp)-1]
	base := float64(len(classes)) - pos

	// Compute the area under ntp and under the
	// minimum bound.
	x := floats.Span(make([]float64, len(classes)+1), 0, float64(len(classes)))
	aucNTP := integrate.Trapezoidal(x, ntp)
	aucMin := pos * pos / 2

	// Calculate the area under the curve
	// within the bounding parallelogram.
	auc := aucNTP - aucMin

	// Calculate the area within the bounding
	// parallelogram.
	par := pos * base

	// The AUC is the ratio of the area under
	// the curve within the bounding parallelogram
	// and the total parallelogram bound.
	auc /= par

	fmt.Printf("number of true positives: %v\n", ntp)
	fmt.Printf("auc: %v\n", auc)

	// Output:
	// number of true positives: [0 0 1 1 2]
	// auc: 0.25
}

func ExampleTOC_aUC_weighted() {
	classes := []bool{true, false, true, false}
	weights := []float64{1, 2, 2, 1}

	min, ntp, max := stat.TOC(classes, weights)

	// Compute the area under ntp and under the
	// minimum and maximum bounds.
	x := make([]float64, len(classes)+1)
	floats.CumSum(x[1:], weights)
	aucNTP := integrate.Trapezoidal(x, ntp)
	aucMin := integrate.Trapezoidal(x, min)
	aucMax := integrate.Trapezoidal(x, max)

	// Calculate the area under the curve
	// within the bounding parallelogram.
	auc := aucNTP - aucMin

	// Calculate the area within the bounding
	// parallelogram.
	par := aucMax - aucMin

	// The AUC is the ratio of the area under
	// the curve within the bounding parallelogram
	// and the total parallelogram bound.
	auc /= par

	fmt.Printf("number of true positives: %v\n", ntp)
	fmt.Printf("auc: %f\n", auc)

	// Output:
	// number of true positives: [0 0 2 2 3]
	// auc: 0.444444
}
