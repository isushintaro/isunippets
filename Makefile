#!/usr/bin/make -f

.PHONY: install
install:
	go install

.PHONY: test
test:
	go test -v ./...

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: check
check: fmt vet test
