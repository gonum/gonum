export CGO_LDFLAGS="-framework Accelerate"
go get github.com/jonlawlor/blas
pushd cgo
go install -v -x
popd
