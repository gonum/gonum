Gonum LAPACK  [![Build Status](https://travis-ci.org/gonum/lapack.svg)](https://travis-ci.org/gonum/lapack)  [![Coverage Status](https://img.shields.io/coveralls/gonum/lapack.svg)](https://coveralls.io/r/gonum/lapack)
======

A collection of packages to provide LAPACK functionality for the Go programming
language (http://golang.org)

This is work in progress. Breaking changes are likely to happen.

## Installation

```
  go get github.com/gonum/blas
```


Install OpenBLAS:
```
  git clone https://github.com/xianyi/OpenBLAS
  cd OpenBLAS
  make
```

For Windows you can download binary packages for OpenBLAS at
http://sourceforge.net/projects/openblas/files/

generate lapack bindings
```
  cd $GOPATH/src/github.com/gonum/lapack/clapack
  ./genLapack.pl -L/path/to/OpenBLAS -lopenblas
```

If you want to use the Intel MKL and all of your paths are properly set
```
  cd $GOPATH/src/github.com/gonum/lapack/clapack
  ./genLapack.pl -lmkl_rt
```
should work.

## Packages

### lapack

Defines the LAPACK API based on http://www.netlib.org/lapack/lapacke.html

### lapack/clapack

Binding to a C implementation of the lapacke interface (e.g. OpenBLAS or intel MKL)

The linker flags (i.e. path to the BLAS library and library name) might have to be adapted.

The recommended (free) option for good performance on both linux and darwin is OpenBLAS.

## Issues

If you find any bugs, feel free to file an issue on the github issue tracker. Discussions on API changes, added features, code review, or similar requests are preferred on the gonum-dev Google Group.

https://groups.google.com/forum/#!forum/gonum-dev

## License

Please see github.com/gonum/license for general license information, contributors, authors, etc on the Gonum suite of packages.
