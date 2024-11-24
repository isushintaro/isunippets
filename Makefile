#!/usr/bin/make -f

.PHONY: test
test:
	go test -v ./...

.PHONY: bench
bench:
	go test -v -bench . -cpu 1 -benchmem ./... -run Benchmark

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: vet
vet:
	go vet ./...

.PHONY: check
check: fmt vet test bench

.PHONY: pprof
pprof:
	go tool pprof -http=:3104 cpu.pprof
