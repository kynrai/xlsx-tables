all: test build

test:
	@go test -v ./...
.PHONY: test

bench:
	@go test -v -bench=. ./...
.PHONY: bench

deps:
	@go get -u && go mod tidy && go mod vendor
.PHONY: deps

build:
	@go install -v ./...
.PHONY: build

