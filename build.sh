#!/bin/bash
go build -ldflags "-X github.com/okke/elmo/core.Build=`git rev-parse HEAD`" -o build/elmo tools/elmo/main.go
