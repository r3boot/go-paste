#!/bin/sh

# Install required packages
apk add git go ca-certificates

# Configure go
export GOPATH=/workspace
go get -v gopkg.in/redis.v3

# Build software
cd ${GOPATH}/src/github.com/r3boot/go-paste \
&& go build -v -o build/go-paste commands/go-paste/go-paste.go \
&& install -o root -g root -m 0755 build/go-paste /opt/go-paste
