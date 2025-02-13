#!/bin/bash

test -z "$(go tool golang.org/x/tools/cmd/goimports -d .)"
if [[ -n "$(gofmt -s -l .)" ]]; then
	echo -e '\e[31mCode not gofmt simplified in:\n\n'
	gofmt -s -l .
	echo -e "\e[0"
fi
