#!/bin/bash
go build -ldflags "-X github.com/okke/elmo/core.Build=`git rev-parse HEAD`" -o elmo tools/elmo/main.go
