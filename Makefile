#!/usr/bin/make -f

test:
	go build ./...
	go generate ./...
	go test -v ./...

docs:
	go install github.com/robertkrimen/godocdown/godocdown
	godocdown > README.md