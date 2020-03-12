#!/usr/bin/make -f

simple-test:
	go test -timeout=1s -count=1 ./...

test:
	go test -timeout=1s -race -covermode=atomic ./...

compile:
	go build ./...

build: test compile

.PHONY: simple-test test compile build
