#!/usr/bin/make -f

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

.PHONY: pprof
pprof:
	go tool pprof -http=:3104 cpu.pprof
