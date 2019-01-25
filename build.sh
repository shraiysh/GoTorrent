#!/bin/bash

# all files should be here
files=( 
	main.go
	parser/*.go
	tracker/*.go
)

test_files=(

)

# script for formatting 
for i in "${files[@]}"; do
	gofmt -w "$i"
done 

goreportcard-cli -v


# script for running test files
for i in "${test_files[@]}"; do
	go test -coverprofile "$i"
done

echo "Build Successful!"
