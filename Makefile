
all: test binaries

first: deps test run

binaries: linux64

linux64:
	GOOS=linux GOARCH=amd64 go build -o bin/contra contra.go

deps:
	dep ensure -v

fmt:
	go fmt $(shell go list ./... | grep -v /vendor/)

vet:
	go vet $(shell go list ./... | grep -v /vendor/)

lint:
	golint -set_exit_status $(shell go list ./... | grep -v /vendor/)

test: fmt vet lint
	go test $(shell go list ./... | grep -v /vendor/)

run: linux64
	./bin/contra

testrun: test run

.PHONY: all deps fmt vet test run testrun
.PHONY: binaries linux64 first
