package mathext

import "github.com/gonum/mathext/internal/amos"

// Airy returns the Airy function at z.
func Airy(z complex128) complex128 {
	id := 0
	kode := 1
	air, aii, _ := amos.Zairy(real(z), imag(z), id, kode)
	return complex(air, aii)
}

// AiryDeriv computes the derivative of the Airy function at z.
func AiryDeriv(z complex128) complex128 {
	id := 1
	kode := 1
	air, aii, _ := amos.Zairy(real(z), imag(z), id, kode)
	return complex(air, aii)
}
