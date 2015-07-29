// Copyright ©2015 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"fmt"
	"math/rand"
	"reflect"

	"github.com/gonum/blas/blas64"
	"github.com/gonum/floats"
	"gopkg.in/check.v1"
)

// legalSizeSameRectangular returns whether the two matrices have the same rectangular shape.
func legalSizeSameRectangular(ar, ac, br, bc int) bool {
	if ar != br {
		return false
	}
	if ac != bc {
		return false
	}
	return true
}

// legalSizeSameSquare returns whether the two matrices have the same square shape.
func legalSizeSameSquare(ar, ac, br, bc int) bool {
	if ar != br {
		return false
	}
	if ac != bc {
		return false
	}
	if ar != ac {
		return false
	}
	return true
}

// legalTypesAll returns true for all Matrix types.
func legalTypesAll(a, b Matrix) bool {
	return true
}

// legalTypesSym returns whether both input arguments are Symmetric.
func legalTypesSym(a, b Matrix) bool {
	if _, ok := a.(Symmetric); !ok {
		return false
	}
	if _, ok := b.(Symmetric); !ok {
		return false
	}
	return true
}

// legalTypesNotVecVec returns whether the first input is an arbitrary Matrix
// and the second input is a *Vector.
func legalTypesNotVecVec(a, b Matrix) bool {
	_, ok := b.(*Vector)
	return ok
}

// legalDims returns whether {m,n} is a valid dimension of the given matrix type.
func legalDims(a Matrix, m, n int) bool {
	switch t := a.(type) {
	default:
		panic("legal dims type not coded")
	case Untransposer:
		return legalDims(t.Untranspose(), n, m)
	case *Dense, *basicMatrix, *basicVectorer:
		if m < 0 || n < 0 {
			return false
		}
		return true
	case *SymDense, *TriDense, *basicSymmetric, *basicTriangular:
		if m < 0 || n < 0 || m != n {
			return false
		}
		return true
	case *Vector:
		if m < 0 || n < 0 {
			return false
		}
		return n == 1
	}
}

// returnAs returns the matrix a with the type of t. Used for making a concrete
// type and changing to the basic form.
func returnAs(a, t Matrix) Matrix {
	switch mat := a.(type) {
	default:
		panic("unknown type for a")
	case *Dense:
		switch t.(type) {
		default:
			panic("bad type")
		case *Dense:
			return mat
		case *basicMatrix:
			return asBasicMatrix(mat)
		case *basicVectorer:
			return asBasicVectorer(mat)
		}
	case *SymDense:
		switch t.(type) {
		default:
			panic("bad type")
		case *SymDense:
			return mat
		case *basicSymmetric:
			return asBasicSymmetric(mat)
		}
	case *TriDense:
		switch t.(type) {
		default:
			panic("bad type")
		case *TriDense:
			return mat
		case *basicTriangular:
			return asBasicTriangular(mat)
		}
	}
}

// retranspose returns the matrix m inside an Untransposer of the type
// of a.
func retranspose(a, m Matrix) Matrix {
	switch a.(type) {
	case Transpose:
		return Transpose{m}
	case TransposeTri:
		return TransposeTri{m.(Triangular)}
	case Untransposer:
		panic("unknown transposer type")
	default:
		panic("a is not an untransposer")
	}
}

// makeRandOf returns a new randomly filled m×n matrix of the underlying matrix type.
func makeRandOf(a Matrix, m, n int) Matrix {
	var matrix Matrix
	switch t := a.(type) {
	default:
		panic("unknown type for make rand of")
	case Untransposer:
		matrix = retranspose(a, makeRandOf(t.Untranspose(), n, m))
	case *Dense, *basicMatrix, *basicVectorer:
		mat := NewDense(m, n, nil)
		for i := 0; i < m; i++ {
			for j := 0; j < n; j++ {
				mat.Set(i, j, rand.Float64())
			}
		}
		matrix = returnAs(mat, t)
	case *Vector:
		if m == 0 && n == 0 {
			return &Vector{}
		}
		if n != 1 {
			panic(fmt.Sprintf("bad vector size: m = %v, n = %v", m, n))
		}
		length := m
		inc := 1
		if t.mat.Inc != 0 {
			inc = t.mat.Inc
		}
		mat := &Vector{
			mat: blas64.Vector{
				Inc:  inc,
				Data: make([]float64, inc*length),
			},
			n: length,
		}
		for i := 0; i < length; i++ {
			mat.SetVec(i, rand.Float64())
		}
		return mat
	case *SymDense, *basicSymmetric:
		if m != n {
			panic("bad size")
		}
		mat := NewSymDense(n, nil)
		for i := 0; i < m; i++ {
			for j := i; j < n; j++ {
				mat.SetSym(i, j, rand.Float64())
			}
		}
		matrix = returnAs(mat, t)
	case *TriDense, *basicTriangular:
		if m != n {
			panic("bad size")
		}
		_, upper := t.(Triangular).Triangle()
		mat := NewTriDense(n, upper, nil)
		if upper {
			for i := 0; i < m; i++ {
				for j := i; j < n; j++ {
					mat.SetTri(i, j, rand.Float64())
				}
			}
		} else {
			for i := 0; i < m; i++ {
				for j := 0; j <= i; j++ {
					mat.SetTri(i, j, rand.Float64())
				}
			}
		}
		matrix = returnAs(mat, t)
	}
	if mr, mc := matrix.Dims(); mr != m || mc != n {
		panic(fmt.Sprintf("makeRandOf for %T returns wrong size: %d×%d != %d×%d", a, m, n, mr, mc))
	}
	return matrix
}

// makeCopyOf returns a copy of the matrix.
func makeCopyOf(a Matrix) Matrix {
	switch t := a.(type) {
	default:
		panic("unknown type in makeCopyOf")
	case Untransposer:
		return retranspose(a, makeCopyOf(t.Untranspose()))
	case *Dense, *basicMatrix, *basicVectorer:
		var m Dense
		m.Clone(a)
		return returnAs(&m, t)
	case *SymDense, *basicSymmetric:
		n := t.(Symmetric).Symmetric()
		m := NewSymDense(n, nil)
		m.CopySym(t.(Symmetric))
		return returnAs(m, t)
	case *TriDense, *basicTriangular:
		n, upper := t.(Triangular).Triangle()
		m := NewTriDense(n, upper, nil)
		if upper {
			for i := 0; i < n; i++ {
				for j := i; j < n; j++ {
					m.SetTri(i, j, t.At(i, j))
				}
			}
		} else {
			for i := 0; i < n; i++ {
				for j := 0; j <= i; j++ {
					m.SetTri(i, j, t.At(i, j))
				}
			}
		}
		return returnAs(m, t)
	case *Vector:
		m := &Vector{
			mat: blas64.Vector{
				Inc:  t.mat.Inc,
				Data: make([]float64, len(t.mat.Data)),
			},
			n: t.n,
		}
		copy(m.mat.Data, t.mat.Data)
		return m
	}
}

// sameType returns true if a and b have the same underlying type.
func sameType(a, b Matrix) bool {
	return reflect.ValueOf(a).Type() == reflect.ValueOf(b).Type()
}

// maybeSame returns true if the two matrices could be represented by the same
// pointer.
func maybeSame(receiver, a Matrix) bool {
	rr, rc := receiver.Dims()
	u, ok := a.(Untransposer)
	if ok {
		a = u.Untranspose()
	}
	ar, ac := a.Dims()
	if !sameType(receiver, a) {
		return false
	}
	return rr == ar && rc == ac
}

// equalApprox returns whether the elements of a and b are the same to within
// the tolerance.
func equalApprox(a, b Matrix, tol float64) bool {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ar != br {
		return false
	}
	if ac != bc {
		return false
	}
	for i := 0; i < ar; i++ {
		for j := 0; j < ac; j++ {
			if !floats.EqualWithinAbsOrRel(a.At(i, j), b.At(i, j), tol, tol) {
				return false
			}
		}
	}
	return true
}

// equal returns true if the matrices have equal entries.
func equal(a, b Matrix) bool {
	ar, ac := a.Dims()
	br, bc := b.Dims()
	if ar != br {
		return false
	}
	if ac != bc {
		return false
	}
	for i := 0; i < ar; i++ {
		for j := 0; j < ac; j++ {
			if a.At(i, j) != b.At(i, j) {
				return false
			}
		}
	}
	return true
}

// underlyingData extracts the underlying data of the matrix a.
func underlyingData(a Matrix) []float64 {
	switch t := a.(type) {
	default:
		panic("matrix type not implemented for extracting underlying data")
	case Untransposer:
		return underlyingData(t.Untranspose())
	case *Dense:
		return t.mat.Data
	case *SymDense:
		return t.mat.Data
	case *Vector:
		return t.mat.Data
	}
}

// testTwoInput tests a method that has two input arguments.
func testTwoInput(c *check.C,
	// name is the name of the method being tested.
	name string,

	// receiver is a value of the receiver type.
	receiver Matrix,

	// method is the generalized receiver.Method(a, b).
	method func(receiver, a, b Matrix),

	// denseComparison performs the same operation as method, but with dense
	// matrices for comparison with the result.
	denseComparison func(receiver, a, b *Dense),

	// legalTypes returns whether the concrete types in Matrix are valid for
	// the method.
	legalTypes func(a, b Matrix) bool,

	// dimsOK returns whether the matrix sizes are valid for the method.
	legalSize func(ar, ac, br, bc int) bool,
) {
	strideVec := &Vector{
		mat: blas64.Vector{
			Inc: 10,
		},
	}
	// It is useful to isolate a single Matrix in the types list during debugging.
	// Ensure that strideVec is always a used variable to avoid compile errors
	// when commenting out types.
	_ = strideVec

	// Loop over all of the matrix types.
	types := []Matrix{
		&Dense{},
		&SymDense{},
		NewTriDense(0, true, nil),
		NewTriDense(0, false, nil),
		NewVector(0, nil),
		Transpose{NewVector(0, nil)},
		strideVec,
		&basicMatrix{},
		&basicVectorer{},
		&basicSymmetric{},
		&basicTriangular{},

		Transpose{&Dense{}},
		Transpose{NewTriDense(0, true, nil)},
		TransposeTri{NewTriDense(0, true, nil)},
		Transpose{NewTriDense(0, false, nil)},
		TransposeTri{NewTriDense(0, false, nil)},
		Transpose{strideVec},
		Transpose{&basicMatrix{}},
		Transpose{&basicVectorer{}},
		Transpose{&basicSymmetric{}},
		Transpose{&basicTriangular{}},
	}

	for _, aMat := range types {
		for _, bMat := range types {
			// Loop over all of the size combinations (bigger, smaller, etc.).
			for _, test := range []struct {
				ar, ac, br, bc int
			}{
				{1, 1, 1, 1},
				{6, 6, 6, 6},
				{7, 7, 7, 7},

				{1, 1, 1, 5},
				{1, 1, 5, 1},
				{1, 5, 1, 1},
				{5, 1, 1, 1},

				{6, 6, 6, 11},
				{6, 6, 11, 6},
				{6, 11, 6, 6},
				{11, 6, 6, 6},
				{11, 11, 11, 6},
				{11, 11, 6, 11},
				{11, 6, 11, 11},
				{6, 11, 11, 11},

				{1, 1, 5, 5},
				{1, 5, 1, 5},
				{1, 5, 5, 1},
				{5, 1, 1, 5},
				{5, 1, 5, 1},
				{5, 5, 1, 1},
				{6, 6, 11, 11},
				{6, 11, 6, 11},
				{6, 11, 11, 6},
				{11, 6, 6, 11},
				{11, 6, 11, 6},
				{11, 11, 6, 6},

				{1, 1, 17, 11},
				{1, 1, 11, 17},
				{1, 11, 1, 17},
				{1, 17, 1, 11},
				{1, 11, 17, 1},
				{1, 17, 11, 1},
				{11, 1, 1, 17},
				{17, 1, 1, 11},
				{11, 1, 17, 1},
				{17, 1, 11, 1},
				{11, 17, 1, 1},
				{17, 11, 1, 1},

				{6, 6, 1, 11},
				{6, 6, 11, 1},
				{6, 11, 6, 1},
				{6, 1, 6, 11},
				{6, 11, 1, 6},
				{6, 1, 11, 6},
				{11, 6, 6, 1},
				{1, 6, 6, 11},
				{11, 6, 1, 6},
				{1, 6, 11, 6},
				{11, 1, 6, 6},
				{1, 11, 6, 6},

				{6, 6, 17, 1},
				{6, 6, 1, 17},
				{6, 1, 6, 17},
				{6, 17, 6, 1},
				{6, 1, 17, 6},
				{6, 17, 1, 6},
				{1, 6, 6, 17},
				{17, 6, 6, 1},
				{1, 6, 17, 6},
				{17, 6, 1, 6},
				{1, 17, 6, 6},
				{17, 1, 6, 6},

				{6, 6, 17, 11},
				{6, 6, 11, 17},
				{6, 11, 6, 17},
				{6, 17, 6, 11},
				{6, 11, 17, 6},
				{6, 17, 11, 6},
				{11, 6, 6, 17},
				{17, 6, 6, 11},
				{11, 6, 17, 6},
				{17, 6, 11, 6},
				{11, 17, 6, 6},
				{17, 11, 6, 6},
			} {
				// Skip the test if any argument would not be assignable to the
				// method's corresponding input parameter or it is not possible
				// to construct an argument of the requested size.
				if !legalTypes(aMat, bMat) {
					continue
				}
				if !legalDims(aMat, test.ar, test.ac) {
					continue
				}
				if !legalDims(bMat, test.br, test.bc) {
					continue
				}
				a := makeRandOf(aMat, test.ar, test.ac)
				b := makeRandOf(bMat, test.br, test.bc)

				// Compute the true answer if the sizes are legal.
				dimsOK := legalSize(test.ar, test.ac, test.br, test.bc)
				var want Dense
				if dimsOK {
					var aDense, bDense Dense
					aDense.Clone(a)
					bDense.Clone(b)
					denseComparison(&want, &aDense, &bDense)
				}
				aCopy := makeCopyOf(a)
				bCopy := makeCopyOf(b)

				// Test the method for a zero-value of the receiver.
				aType, aTrans := untranspose(a)
				bType, bTrans := untranspose(b)
				errStr := fmt.Sprintf("%T.%s(%T, %T), sizes: %#v, atrans %v, btrans %v", receiver, name, aType, bType, test, aTrans, bTrans)
				zero := makeRandOf(receiver, 0, 0)
				panicked, err := panics(func() { method(zero, a, b) })
				if !dimsOK && !panicked {
					c.Errorf("Did not panic with illegal size: %s", errStr)
					continue
				}
				if dimsOK && panicked {
					c.Errorf("Panicked with legal size: %s %s", errStr, err)
					continue
				}
				if !equal(a, aCopy) {
					c.Errorf("First input argument changed in call: %s", errStr)
				}
				if !equal(b, bCopy) {
					c.Errorf("Second input argument changed in call: %s", errStr)
				}
				if !dimsOK {
					continue
				}
				if !equalApprox(zero, &want, 1e-14) {
					c.Errorf("Answer mismatch with zero receiver: %s", errStr)
					continue
				}

				// Test the method with a non-zero-value of the receiver.
				// The receiver has been overwritten in place so use its size
				// to construct a new random matrix.
				rr, rc := zero.Dims()
				nonZero := makeRandOf(receiver, rr, rc)
				panicked, _ = panics(func() { method(nonZero, a, b) })
				if panicked {
					c.Errorf("Panicked with non-zero receiver: %s", errStr)
				}
				if !equalApprox(nonZero, &want, 1e-14) {
					c.Errorf("Answer mismatch non-zero receiver: %s", errStr)
				}

				// Test with an incorrectly sized matrix.
				switch receiver.(type) {
				default:
					panic("matrix type not coded for incorrect receiver size")
				case *Dense:
					wrongSize := makeRandOf(receiver, rr+1, rc)
					panicked, _ = panics(func() { method(wrongSize, a, b) })
					if !panicked {
						c.Errorf("Did not panic with wrong number of rows: %s", errStr)
					}
					wrongSize = makeRandOf(receiver, rr, rc+1)
					panicked, _ = panics(func() { method(wrongSize, a, b) })
					if !panicked {
						c.Errorf("Did not panic with wrong number of columns: %s", errStr)
					}
				case *TriDense, *SymDense:
					// Add to the square size.
					wrongSize := makeRandOf(receiver, rr+1, rc+1)
					panicked, _ = panics(func() { method(wrongSize, a, b) })
					if !panicked {
						c.Errorf("Did not panic with wrong size: %s", errStr)
					}
				case *Vector:
					// Add to the column length.
					wrongSize := makeRandOf(receiver, rr+1, rc)
					panicked, _ = panics(func() { method(wrongSize, a, b) })
					if !panicked {
						c.Errorf("Did not panic with wrong number of rows: %s", errStr)
					}
				}

				// The receiver and an input may share a matrix pointer
				// if the type and size of the receiver and one of the
				// arguments match. Test the method works properly
				// when this is the case.
				aMaybeSame := maybeSame(nonZero, a)
				bMaybeSame := maybeSame(nonZero, b)
				if aMaybeSame {
					aSame := makeCopyOf(a)
					receiver = aSame
					u, ok := aSame.(Untransposer)
					if ok {
						receiver = u.Untranspose()
					}
					preData := underlyingData(receiver)
					panicked, err = panics(func() { method(receiver, aSame, b) })
					if panicked {
						c.Errorf("Panics when a maybeSame: %s: %s", errStr, err)
					} else {
						if !equalApprox(receiver, &want, 1e-14) {
							c.Errorf("Wrong answer when a maybeSame: %s", errStr)
						}
						postData := underlyingData(receiver)
						if !floats.Equal(preData, postData) {
							c.Errorf("Original data slice not modified when a maybeSame: %s", errStr)
						}
					}
				}
				if bMaybeSame {
					bSame := makeCopyOf(b)
					receiver = bSame
					u, ok := bSame.(Untransposer)
					if ok {
						receiver = u.Untranspose()
					}
					preData := underlyingData(receiver)
					panicked, err = panics(func() { method(receiver, a, bSame) })
					if panicked {
						c.Errorf("Panics when b maybeSame: %s", errStr)
					} else {
						if !equalApprox(receiver, &want, 1e-14) {
							c.Errorf("Wrong answer when b maybeSame: %s: %s", errStr, err)
						}
						postData := underlyingData(receiver)
						if !floats.Equal(preData, postData) {
							c.Errorf("Original data slice not modified when b maybeSame: %s", errStr)
						}
					}
				}
				if aMaybeSame && bMaybeSame {
					aSame := makeCopyOf(a)
					receiver = aSame
					u, ok := aSame.(Untransposer)
					if ok {
						receiver = u.Untranspose()
					}
					// Ensure that b is the correct transpose type if applicable.
					// The receiver is always a concrete type so use it.
					bSame := receiver
					u, ok = bSame.(Untransposer)
					if ok {
						bSame = retranspose(bSame, receiver)
					}
					// Compute the real answer for this case. It is different
					// from the inital answer since now a and b have the
					// same data.
					zero = makeRandOf(zero, 0, 0)
					method(zero, aSame, bSame)
					preData := underlyingData(receiver)
					panicked, err = panics(func() { method(receiver, aSame, bSame) })
					if panicked {
						c.Errorf("Panics when both maybeSame: %s: %s", errStr, err)
					} else {
						if !equalApprox(receiver, zero, 1e-14) {
							c.Errorf("Wrong answer when both maybeSame: %s", errStr)
						}
						postData := underlyingData(receiver)
						if !floats.Equal(preData, postData) {
							c.Errorf("Original data slice not modified when both maybeSame: %s", errStr)
						}
					}
				}
			}
		}
	}
}
