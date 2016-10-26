#!/usr/bin/make -f

test:
	go build ./...
	go generate ./...
	go test -short ./...

docs:
	go get -u github.com/robertkrimen/godocdown/godocdown
	go install github.com/robertkrimen/godocdown/godocdown
	godocdown > README.md

