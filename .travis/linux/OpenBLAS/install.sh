set -ex

# fetch fortran to build OpenBLAS
sudo apt-get update -qq && sudo apt-get install -qq gfortran

# fetch OpenBLAS
pushd ~
sudo git clone --depth=1 git://github.com/xianyi/OpenBLAS

# make OpenBLAS
pushd OpenBLAS
echo OpenBLAS $(git rev-parse HEAD)
sudo make FC=gfortran &> /dev/null && sudo make PREFIX=/usr install
popd

# fetch cblas reference lib
curl http://www.netlib.org/blas/blast-forum/cblas.tgz | tar -zx

# make cblas and install
pushd CBLAS
sudo mv Makefile.LINUX Makefile.in
sudo BLLIB=/usr/lib/libopenblas.a make alllib
sudo mv lib/cblas_LINUX.a /usr/lib/libcblas.a
popd
popd

# install gonum/blas against OpenBLAS
export CGO_LDFLAGS="-L/usr/lib -lopenblas"
go get github.com/gonum/blas
pushd cgo
go install -v -x
popd

# run the OS common installation script
source ${TRAVIS_BUILD_DIR}/.travis/$TRAVIS_OS_NAME/install.sh
