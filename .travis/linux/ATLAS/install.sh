sudo apt-get update -qq && sudo apt-get install -qq libatlas-base-dev
if [ $? != 0 ]; then exit 1; fi

export CGO_LDFLAGS="-L/usr/lib -lblas"
go get github.com/gonum/blas
if [ $? != 0 ]; then exit 1; fi
