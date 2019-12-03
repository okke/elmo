#!/bin/bash
echo generate parser
cd core
go generate

echo run tests
cd ..
go test ./...

echo build binary
go build -ldflags "-X github.com/okke/elmo/core.Build=`git rev-parse HEAD`" -o build/elmo tools/elmo/main.go
