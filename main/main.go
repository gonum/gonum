// logdemo.go
package main

import (
	"fmt"

	"gonum.org/v1/gonum/mathext"
)

func main() {
	// Example values to try
	values := []complex128{
		-1 + 0i,
		1 + 1i,
		-2 + 3i,
		complex(0, 1),
		1000 + 0i,
		-1 + 0i,
	}

	for _, z := range values {
		l := mathext.Li2(z) // principal branch log
		fmt.Printf("Log(%v) = %.15f + %.15fi\n", z, real(l), imag(l))
	}

}
