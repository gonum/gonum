// Copyright ©2019 The Gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package combin_test

import (
	"fmt"

	"gonum.org/v1/gonum/stat/combin"
)

func ExampleCombinations() {
	data := []string{"a", "b", "c", "d", "e"}
	cs := combin.Combinations(len(data), 2)
	for _, c := range cs {
		fmt.Printf("%s%s\n", data[c[0]], data[c[1]])
	}

	// Output:
	// ab
	// ac
	// ad
	// ae
	// bc
	// bd
	// be
	// cd
	// ce
	// de
}
