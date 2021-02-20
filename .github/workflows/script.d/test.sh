#!/bin/bash

# Travis clobbers GOARCH env vars so this is necessary here.
if [[ -n $FORCE_GOARCH ]]; then
	export GOARCH=$FORCE_GOARCH
fi

go get -d -v ./...
go build -v ./...
go test $TAGS -v ./...
