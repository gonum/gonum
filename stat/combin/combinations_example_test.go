package combin_test

import (
	"fmt"

	"gonum.org/v1/gonum/stat/combin"
)

func Example() {
	// combin provides several ways to work with the combinations of
	// different objects.
	// The first way is to generate them all directly:
	fmt.Println("Generate list:")
	n := 5
	k := 3
	list := combin.Combinations(n, k)
	for i, v := range list {
		fmt.Println(i, v)
	}
	// This is easy, especially for computing functions of different
	// combinations concurrently. However, the number of combinations
	// can be very large, and generating all at once can use a lot
	// of memory.

	// The second way is we can use a generator to step through each
	// combination.
	fmt.Println("\nUse generator:")
	gen := combin.NewCombinationGenerator(n, k)
	idx := 0
	for gen.Next() {
		fmt.Println(idx, gen.Combination(nil)) // can also store in-place.
		idx++
	}

	// The third way is using IndexToCombination to gain random access.
	// This provides two-way access between integers and combinations.
	fmt.Println("\nUse indexing:")
	comb := make([]int, k)
	for i := 0; i < combin.Binomial(n, k); i++ {
		combin.IndexToCombination(comb, i, n, k) // can also use nil.
		idx := combin.CombinationIndex(comb, n, k)
		fmt.Println(i, comb, idx)
	}
	// Note that the iteration order is the same for all methods.

	// Output:
	// Generate list:
	// 0 [0 1 2]
	// 1 [0 1 3]
	// 2 [0 1 4]
	// 3 [0 2 3]
	// 4 [0 2 4]
	// 5 [0 3 4]
	// 6 [1 2 3]
	// 7 [1 2 4]
	// 8 [1 3 4]
	// 9 [2 3 4]
	//
	// Use generator:
	// 0 [0 1 2]
	// 1 [0 1 3]
	// 2 [0 1 4]
	// 3 [0 2 3]
	// 4 [0 2 4]
	// 5 [0 3 4]
	// 6 [1 2 3]
	// 7 [1 2 4]
	// 8 [1 3 4]
	// 9 [2 3 4]
	//
	// Use indexing:
	// 0 [0 1 2] 0
	// 1 [0 1 3] 1
	// 2 [0 1 4] 2
	// 3 [0 2 3] 3
	// 4 [0 2 4] 4
	// 5 [0 3 4] 5
	// 6 [1 2 3] 6
	// 7 [1 2 4] 7
	// 8 [1 3 4] 8
	// 9 [2 3 4] 9
}
