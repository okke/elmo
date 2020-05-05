#!/bin/bash
echo generate parser
cd core
go generate

echo run tests
cd ..
go test ./...

echo build binary
go build -ldflags "-X github.com/okke/elmo/core.CommitHash=`git rev-parse HEAD` -X github.com/okke/elmo/core.BranchName=`git rev-parse --abbrev-ref HEAD`" -o build/elmo tools/elmo/main.go
