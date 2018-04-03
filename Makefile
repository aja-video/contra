VERSION=$(shell ./bin/contra -q -version)
all: test binaries

first: deps test run

binaries: linux64

linux64:
	GOOS=linux GOARCH=amd64 go build -o bin/contra contra.go

packages: binaries rpm64 deb64

rpm64:
	rm -rf build/rpm
	mkdir -p build/rpm/contra/usr/local/bin
	mkdir -p build/rpm/contra/etc
	cp bin/contra build/rpm/contra/usr/local/bin/
	cp contra.example.conf build/rpm/contra/etc/contra.conf.dist
	fpm -s dir -t rpm -n contra -a x86_64 --epoch 0 -v $(VERSION) -C build/rpm/contra .
	mv contra-$(VERSION)-1.x86_64.rpm bin/

deb64:
	rm -rf build/deb
	mkdir -p build/deb/contra/usr/local/bin
	cp bin/contra build/deb/contra/usr/local/bin/
	mkdir -p build/deb/contra/etc
	cp contra.example.conf build/deb/contra/etc/contra.conf.dist
	fpm -s dir -t deb -n contra -a amd64 -v $(VERSION) -C build/deb/contra .
	mv contra_$(VERSION)_amd64.deb bin/

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
