package sliceops

import (
	"fmt"
)

// Set of examples for all the functions

func ExampleAdd() {
	s1 := []float64{1, 2, 3, 4}
	s2 := []float64{5, 6, 7, 8}
	Add(s1, s2)
	fmt.Println("s1 = ", s1)
	fmt.Println("s2 = ", s2)
	// Output:
	// s1 =  [6 8 10 12]
	// s2 =  [5 6 7 8]
}
