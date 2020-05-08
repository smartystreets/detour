#!/usr/bin/make -f

test: fmt
	go test -timeout=1s -race -covermode=atomic ./...

test.simple:
	go test -timeout=1s -count=1 ./...

fmt:
	go fmt ./...

compile:
	go build ./...

build: test compile

.PHONY: test test.simple fmt compile build
