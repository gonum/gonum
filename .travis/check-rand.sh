#!/bin/bash

RAND_IMPORTERS=$(find -name "*.go" -printf '%P\n' | xargs grep -l '"math/rand"')
if [[ -n "$RAND_IMPORTERS" ]]; then
	echo -e '\e[31mImports of "math/rand" in:\n'
	for F in $RAND_IMPORTERS; do
		echo $F
	done
	echo -e '\e[0'
	exit 1
fi
