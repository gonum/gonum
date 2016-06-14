#!/bin/bash

go get -u github.com/Kunde21/sift

go test ./...  -bench Axpy | tee >( sift markL | sed s/markL/mark/ > old.tst ) |  sift mark[^L] >new.tst
benchcmp old.tst new.tst
