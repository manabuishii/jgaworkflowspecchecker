#!/bin/bash
VERSION=0.8.0


env GOOS=darwin GOARCH=amd64 go build -o jgaworkflowmanager-mac-${VERSION} -ldflags "-X main.version=${VERSION} -X main.revision=$(git rev-parse --short HEAD)" main.go

env GOOS=linux GOARCH=amd64 go build -o jgaworkflowmanager-linux-${VERSION} -ldflags "-X main.version=${VERSION} -X main.revision=$(git rev-parse --short HEAD)" main.go

