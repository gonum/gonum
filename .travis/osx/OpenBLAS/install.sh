set -ex

# fetch and install OpenBLAS using homebrew
brew install homebrew/science/openblas

# fetch and install gonum/blas against OpenBLAS
export CGO_LDFLAGS="-L/usr/local/opt/openblas/lib -lopenblas"
go get github.com/gonum/blas
pushd cgo
go install -v -x
popd

# run the OS common installation script
source ${TRAVIS_BUILD_DIR}/.travis/$TRAVIS_OS_NAME/install.sh
