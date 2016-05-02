set -ex

# run the OS common installation script
source ${TRAVIS_BUILD_DIR}/.travis/$TRAVIS_OS_NAME/install.sh

# change to native directory so we don't test code that depends on an external
# blas library
cd native
