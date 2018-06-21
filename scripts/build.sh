#!/usr/bin/env bash

dep ensure

echo "Running gofmt"
gofmt -s -w .

go build -ldflags "-s -w"