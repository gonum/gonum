brew install homebrew/science/openblas
export CGO_LDFLAGS="-L/usr/local/opt/openblas/lib -lopenblas"
go get github.com/gonum/blas
pushd cgo
go install -v -x
popd
