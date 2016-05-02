set -ex

# This script contains common installation commands for linux.  It should be run
# prior to more specific installation commands for a particular blas library.
go get golang.org/x/tools/cmd/cover
go get github.com/mattn/goveralls
go get github.com/gonum/floats
