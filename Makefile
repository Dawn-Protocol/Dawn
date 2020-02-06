#!/usr/bin/make -f

all: test clean install lint

# The below include contains the tools and runsim targets.
include contrib/devtools/Makefile

build:
	go build -mod=readonly -o build/dawnd ./cmd/dawnd
	go build -mod=readonly -o build/dawncli ./cmd/dawncli
	go build -mod=readonly -o build/dawnrelayer ./cmd/dawnrelayer

clean:
	rm -rf build/

install:
	go install -mod=readonly ./cmd/dawnd
	go install -mod=readonly ./cmd/dawncli
	go install -mod=readonly ./cmd/dawnrelayer

lint:
	@echo "--> Running linter"
	golangci-lint run
	@find . -name '*.go' -type f -not -path "./vendor*" -not -path "*.git*" | xargs gofmt -d -s
	go mod verify

test:
	go test ./...

.PHONY: all build clean install test lint all