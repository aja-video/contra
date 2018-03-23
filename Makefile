
all: test binaries

first: deps test run

binaries: linux64

linux64:
	GOOS=linux GOARCH=amd64 go build -o bin/contra contra.go

deps:
#	go get -u -v ./....
	dep ensure -v

fmt:
#	go fmt ./...
	go fmt $(go list ./... | grep -v /vendor/)

vet:
#	go vet ./...
	go vet $(go list ./... | grep -v /vendor/)

lint:
	golint -set_exit_status $(go list ./... | grep -v /vendor/)

test: fmt vet lint
	go test $(go list ./... | grep -v /vendor/)

run: linux64
	./bin/contra

.PHONY: all deps fmt vet test run
.PHONY: binaries linux64 first
