.DEFAULT_GOAL := build

dependencies:
	go get -d
.PHONY:dependencies

fmt: dependencies
	go fmt ./...
.PHONY:fmt

lint: fmt
	go install golang.org/x/lint/golint@latest
	golint ./...
.PHONY:lint

vet: lint
	go vet ./...
.PHONY:vet

build: vet
	go build
.PHONY:build

