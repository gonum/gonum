#!/bin/bash

set -ex

# Required for format check.
go install golang.org/x/tools/cmd/goimports@latest
# Required for imports check.
go install gonum.org/v1/tools/cmd/check-imports@latest
# Required for copyright header check.
go install gonum.org/v1/tools/cmd/check-copyright@latest
# Required for coverage.
go install golang.org/x/tools/cmd/cover@latest
# Required for dot parser checks.
go install github.com/goccmack/gocc@66c61e9
# Required for rdf parser checks.
go install golang.org/x/tools/cmd/stringer@latest
