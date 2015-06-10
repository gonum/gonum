sudo apt-get update -qq
sudo apt-get install -qq libatlas-base-dev

export CGO_LDFLAGS="-L/usr/lib -lblas"
go get github.com/gonum/blas
