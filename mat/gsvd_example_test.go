// Copyright ©2017 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat_test

import (
	"fmt"
	"log"
	"math"

	"gonum.org/v1/gonum/mat"
)

func ExampleGSVD() {
	// Perform a GSVD factorization on food production/consumption data for the
	// three years 1990, 2000 and 2014, for Africa and Latin America/Caribbean.
	//
	// See Lee et al. doi:10.1371/journal.pone.0030098 and
	// Alter at al. doi:10.1073/pnas.0530258100 for more details.
	var gsvd mat.GSVD
	ok := gsvd.Factorize(FAO.Africa, FAO.LatinAmericaCaribbean, mat.GSVDU|mat.GSVDV|mat.GSVDQ)
	if !ok {
		log.Fatal("GSVD factorization failed")
	}
	s1 := gsvd.ValuesA(nil)
	s2 := gsvd.ValuesB(nil)

	fmt.Printf("Africa\n\ts1 = %.4f\n\n", s1)
	fmt.Printf("Latin America/Caribbean\n\ts2 = %.4f\n", s2)

	// Calculate the antisymmetric angular distances for each eigenvariable.
	fmt.Println("\nSignificance:")
	for i := 0; i < 3; i++ {
		fmt.Printf("\teigenvar_%d: %+.4f\n", i, math.Atan(s1[i]/s2[i])-math.Pi/4)
	}

	// Output:
	//
	// Africa
	// 	s1 = [1.0000 0.9344 0.5118]
	//
	// Latin America/Caribbean
	// 	s2 = [0.0047 0.3563 0.8591]
	//
	// Significance:
	// 	eigenvar_0: +0.7807
	// 	eigenvar_1: +0.4211
	// 	eigenvar_2: -0.2482
}
