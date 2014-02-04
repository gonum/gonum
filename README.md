# Gonum BLAS [![Build Status](https://travis-ci.org/gonum/blas.png?branch=master)](https://travis-ci.org/gonum/blas)

A collection of packages to provide BLAS functionality for the Go programming
language (http://golang.org)

## Installation 
```
  go get github.com/gonum/blas
  cd $GOPATH/src/github.com/gonum/blas
  go install ./...
```
## Packages

### blas

Defines BLAS API (www.netlib.org/blas/blast-forum/cinterface.pdf) split in several interfaces

### blas/goblas

Go implementation of the BLAS API (incomplete, implements most of the float64 API)

### blas/cblas

Binding to a C implementation of the cblas interface (e.g. ATLAS, OpenBLAS, intel MKL)

On linux the linker flags (i.e. path to the BLAS library and library name) might have to be adapted.

### blas/dbw

Wrapper for an implementation of the double precision real (i.e. float64) part of the blas API

You have to register an implementation before you can use the BLAS functions:

```
package main

import (
	"fmt"

	"github.com/gonum/blas/cblas"
	"github.com/gonum/blas/dbw"
)

func init() {
	dbw.Register(cblas.Blas{})
}

func main() {
	v := Vector{[]float64{1, 1, 1}, 3, 1}
	fmt.Println("v has length:", dbw.Nrm2(v))
}
```

### blas/zbw

Wrapper for an implementation of the double precision complex (i.e. complex128) part of the blas API
