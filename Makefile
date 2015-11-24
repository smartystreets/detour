#!/usr/bin/make -f

generate:
	go generate ./...

test:
	go build ./...
	go generate ./...
	go test -v ./...

docs:
	go get -u github.com/robertkrimen/godocdown/godocdown
	go install github.com/robertkrimen/godocdown/godocdown
	godocdown > README.md
