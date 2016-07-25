set -ex

CACHE_DIR=${TRAVIS_BUILD_DIR}/.travis/${BLAS_LIB}.cache

# fetch fortran to build OpenBLAS
sudo apt-get update -qq && sudo apt-get install -qq gfortran

# check if cache exists
if [ -e ${CACHE_DIR}/last_commit_id ]; then
    echo "Cache $CACHE_DIR hit"
    LAST_COMMIT="$(git ls-remote git://github.com/xianyi/OpenBLAS HEAD | grep -o '^\S*')"
    CACHED_COMMIT="$(cat ${CACHE_DIR}/last_commit_id)"
    # determine current OpenBLAS master commit id and compare
    # with commit id in cache directory
    if [ "$LAST_COMMIT" != "$CACHED_COMMIT" ]; then
        echo "Cache Directory $CACHE_DIR has stale commit"
        # if commit is different, delete the cache
        rm -rf ${CACHE_DIR}
    fi
fi

# check if cache directory exists
if [ ! -e ${CACHE_DIR}/last_commit_id ]; then
    if [ -d ${CACHE_DIR} ]; then
        # Travis automatically creates the cache directory if it does not exist,
        # so this is needed for initialization.
        rm -rf ${CACHE_DIR}
    fi

    # cache generation
    echo "Building cache at $CACHE_DIR"
    mkdir ${CACHE_DIR}
    sudo git clone --depth=1 git://github.com/xianyi/OpenBLAS

    pushd OpenBLAS
    echo OpenBLAS $(git rev-parse HEAD)
    # write commit id to cache location
    echo $(git rev-parse HEAD) > ${CACHE_DIR}/last_commit_id
    sudo make FC=gfortran &> /dev/null && sudo make PREFIX=${CACHE_DIR} install
    popd

    curl http://www.netlib.org/blas/blast-forum/cblas.tgz | tar -zx

    pushd CBLAS
    sudo mv Makefile.LINUX Makefile.in
    sudo BLLIB=${CACHE_DIR}/lib/libopenblas.a make alllib
    sudo mv lib/cblas_LINUX.a ${CACHE_DIR}/lib/libcblas.a
    popd

fi

# copy the cache files into /usr
sudo cp -r ${CACHE_DIR}/* /usr/

# install gonum/blas against OpenBLAS
# fetch and install gonum/blas and gonum/matrix
export CGO_LDFLAGS="-L/usr/lib -lopenblas"
go get github.com/gonum/blas
go get github.com/gonum/matrix/mat64

# install clapack against OpenBLAS
pushd cgo/clapack
go install -v -x
popd

# travis compiles commands in script and then executes in bash.  By adding
# set -e we are changing the travis build script's behavior, and the set
# -e lives on past the commands we are providing it.  Some of the travis
# commands are supposed to exit with non zero status, but then continue
# executing.  set -x makes the travis log files extremely verbose and
# difficult to understand.
# 
# see travis-ci/travis-ci#5120
set +ex
