#!/bin/bash

env GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o sshez-osx
env GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o sshez-amd64
env GOOS=linux GOARCH=386 CGO_ENABLED=0 go build -ldflags="-s -w" -o sshez-386
env GOOS=linux GOARCH=arm CGO_ENABLED=0 go build -ldflags="-s -w" -o sshez-arm

echo "done"