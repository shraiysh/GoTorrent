#!/bin/bash

# all files should be here
files=( 
	main.go
	parser/*.go
	tracker/*.go
	torrent/*.go
	queue/*.go
	piece/*.go
	args/*.go
)

# script for formatting 
for i in "${files[@]}"; do
	gofmt -w "$i"
done 

# script for running test files
go test ./... -v

#checking goreportcard-cli for issues
set -e
report=$(goreportcard-cli)

startindex=$(($(echo $report | grep -b -o "gofmt" | cut -d: -f1)+7))
endindex=$(($(echo $report|grep -b -o "go_vet" | cut -d: -f1)-1))
gofmt=${report:$startindex:$endindex-$startindex}
echo "gofmt: "$gofmt
if [ $gofmt == "100%" ]; then
        echo "gofmt passed"
fi
if [ $gofmt != "100%" ]; then
        echo $issuecount" issues. Run \`goreportcard-cli -v\` to check. Ignore the issues from \`vendor\` directory"
	exit 1
fi

startindex=$(($(echo $report | grep -b -o "go_vet" | cut -d: -f1)+8))
endindex=$(($(echo $report|grep -b -o "gocyclo" | cut -d: -f1)-1))
go_vet=${report:$startindex:$endindex-$startindex}
echo "go_vet: "$go_vet
if [ $go_vet == "100%" ]; then
        echo "go_vet passed"
fi
if [ $go_vet != "100%" ]; then
        echo $issuecount" issues. Run \`goreportcard-cli -v\` to check. Ignore the issues from \`vendor\` directory"
	exit 1
fi

startindex=$(($(echo $report | grep -b -o "golint" | cut -d: -f1)+8))
endindex=$(($(echo $report|grep -b -o "ineffassign" | cut -d: -f1)-1))
golint=${report:$startindex:$endindex-$startindex}
echo "go_lint: "$golint
if [ $golint == "100%" ]; then
        echo "go_lint passed"
fi
if [ $golint != "100%" ]; then
        echo $issuecount" issues. Run \`goreportcard-cli -v\` to check. Ignore the issues from \`vendor\` directory"
	exit 1
fi

startindex=$(($(echo $report | grep -b -o "ineffassign" | cut -d: -f1)+13))
endindex=$(($(echo $report|grep -b -o "license" | cut -d: -f1)-1))
ineffassign=${report:$startindex:$endindex-$startindex}
echo "ineffassign: "$ineffassign
if [ $ineffassign == "100%" ]; then
        echo "ineffassign passed"
fi
if [ $ineffassign != "100%" ]; then
        echo $issuecount" issues. Run \`goreportcard-cli -v\` to check. Ignore the issues from \`vendor\` directory"
	exit 1
fi

startindex=$(($(echo $report | grep -b -o "license" | cut -d: -f1)+9))
endindex=$(($(echo $report|grep -b -o "misspell" | cut -d: -f1)-1))
license=${report:$startindex:$endindex-$startindex}
echo "license: "$license
if [ $license == "100%" ]; then
        echo "license passed"
fi
if [ $license != "100%" ]; then
        echo $issuecount" issues. Run \`goreportcard-cli -v\` to check. Ignore the issues from \`vendor\` directory"
	exit 1
fi

startindex=$(($(echo $report | grep -b -o "misspell" | cut -d: -f1)+10))
endindex=${#report}
misspell=${report:$startindex:$endindex-$startindex}
echo "misspell: "$misspell
if [ $misspell == "100%" ]; then
        echo "misspell passed"
fi
if [ $misspell != "100%" ]; then
        echo $issuecount" issues. Run \`goreportcard-cli -v\` to check. Ignore the issues from \`vendor\` directory"
	exit 1
fi

echo "Build Successful!"
