REVISION = $(shell git rev-parse --short HEAD)

all: pre test

pre: fmt lint vet

fmt:
	go fmt ./...

lint:
	golint $$(go list ./... | grep -v /vendor/)

vet:
	go vet $$(go list ./... | grep -v /vendor/)

test:
	go test -cover -v -race -timeout 10s ./...

bench:
	go test -bench . -benchmem -gcflags="-m -l" ./...

dep-init:
	dep ensure

dep-update:
	dep ensure -update

.PHONY: all pre fmt lint vet test bench dep-init dep-update
