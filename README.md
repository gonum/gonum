LAPACK package for the Go language
======

A collection of packages to provide LAPACK functionality for the Go programming
language (http://golang.org)

## Installation 

```
  go get github.com/dane-unltd/blas
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
  cd $GOPATH/src/github.com/dane-unltd/lapack/clapack
  ./genLapack /path/to/OpenBLAS openblas
```

If you want to use the Intel MKL and all of your paths are properly set
```
  cd $GOPATH/src/github.com/dane-unltd/lapack/clapack
  ./genLapack mkl_rt
```
should work.
