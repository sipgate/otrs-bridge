#!/usr/bin/env bash

dep ensure
go build -ldflags "-s -w"