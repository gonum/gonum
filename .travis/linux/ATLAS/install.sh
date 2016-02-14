set -ex

# fetch and install ATLAS
sudo apt-get update -qq
sudo apt-get install -qq libatlas-base-dev


# fetch and install gonum/blas and gonum/matrix
export CGO_LDFLAGS="-L/usr/lib -latlas -llapack_atlas"
go get github.com/gonum/blas
go get github.com/gonum/matrix/mat64

# install lapack against ATLAS
pushd cgo/clapack
go install -v -x
popd
