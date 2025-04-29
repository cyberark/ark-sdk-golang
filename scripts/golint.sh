#!/bin/bash

if [ -z "$GOPATH" ]
then
	export GOPATH=$HOME/go
fi

golint_output=$($GOPATH/bin/golint ./... | grep -v "should have comment" | grep -v "don't use an underscore in package name")

if [[ $golint_output ]]; then
    echo "$golint_output"
    exit 1
else
    echo "Golint executed successfully. No problems found."
fi
exit 0
