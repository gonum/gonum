#!/bin/bash

set -ex

# Avoid contaminating the go.mod/go.sum files.
WORK=$(mktemp -d)
pushd $WORK

# Work around go1.11 issue where go get will not work outside a module.
# TODO(kortschak): Remove when go1.11 is no longer supported.
go mod init fake || true

# Required for format check.
go get golang.org/x/tools/cmd/goimports
# Required for imports check.
go get gonum.org/v1/tools/cmd/check-imports
# Required for copyright header check.
go get gonum.org/v1/tools/cmd/check-copyright
# Required for coverage.
go get golang.org/x/tools/cmd/cover
go get github.com/mattn/goveralls
# Required for dot parser checks.
${TRAVIS_BUILD_DIR}/.travis/script.d/install-gocc.sh 66c61e91b3657c517a6f89d2837d370e61fb9430

# Clean up.
# TODO(kortschak): Remove when go1.11 is no longer supported.
popd
rm -rf $WORK
