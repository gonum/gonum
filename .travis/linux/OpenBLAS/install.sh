sudo apt-get update -qq && sudo apt-get install -qq gfortran
if [ $? != 0 ]; then exit 1; fi
pushd ~
sudo git clone --depth=1 git://github.com/xianyi/OpenBLAS
if [ $? != 0 ]; then popd; exit 1; fi
pushd OpenBLAS
echo OpenBLAS $(git rev-parse HEAD)
sudo make FC=gfortran &> /dev/null && sudo make PREFIX=/usr install
if [ $? != 0 ]; then popd; popd; exit 1; fi
popd
curl http://www.netlib.org/blas/blast-forum/cblas.tgz | tar -zx
if [ $? != 0 ]; then popd; exit 1; fi
pushd CBLAS
sudo mv Makefile.LINUX Makefile.in
sudo BLLIB=/usr/lib/libopenblas.a make alllib
sudo mv lib/cblas_LINUX.a /usr/lib/libcblas.a
popd
popd
export CGO_LDFLAGS="-L/usr/lib -lopenblas"
go get github.com/gonum/blas
if [ $? != 0 ]; then exit 1; fi
pushd cgo
go install -v -x
if [ $? != 0 ]; then popd; exit 1; fi
popd
