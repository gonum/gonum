sudo apt-get update -qq
sudo apt-get install -qq libatlas-base-dev

export CGO_LDFLAGS="-L/usr/lib -latlas -llapack_atlas"

go get github.com/gonum/blas
go get github.com/gonum/matrix/mat64

pushd cgo
sudo perl genLapack.pl -L/usr/lib -latlas -llapack_atlas
popd