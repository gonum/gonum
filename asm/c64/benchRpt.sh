#!/bin/bash

go test -bench . | tee >( sift markF | sed s/markF/mark/ > old.tst ) |  sift mark[^F] >new.tst
benchcmp old.tst new.tst
