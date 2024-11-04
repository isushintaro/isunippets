#!/usr/bin/make -f

.PHONY: install
install:
	go install

.PHONY: test
test: install
	go test -v ./...

.PHONY: fmt
fmt: install
	go fmt ./...

.PHONY: vet
vet: install
	go vet ./...

.PHONY: check
check: fmt vet test
