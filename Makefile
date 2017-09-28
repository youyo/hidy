Name := hidy
Version := $(shell git describe --tags --abbrev=0)
OWNER := youyo
.DEFAULT_GOAL := help

## Setup
setup:
	go get github.com/golang/dep/cmd/dep
	go get github.com/Songmu/make2help/cmd/make2help

## Install dependencies
deps: setup
	dep ensure

## Build
build: deps
	go build -o artifact/$(Name)

## Release
release: build
	ghr -t ${GITHUB_TOKEN} -u $(OWNER) -r $(Name) --replace $(Version) artifact/

## Show help
help:
	@make2help $(MAKEFILE_LIST)

.PHONY: setup deps build release help
