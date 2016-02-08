export CGO_LDFLAGS="-framework Accelerate"
go get github.com/gonum/blas
pushd cgo
go install -v -x
popd
