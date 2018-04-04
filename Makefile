VERSION=$(shell ./bin/contra -q -version)
all: test binaries

first: deps test run

binaries: linux64

linux64:
	GOOS=linux GOARCH=amd64 go build -o bin/contra contra.go

packages: binaries rpm64 deb64

rpm64: binaries
	rm -rf build/rpm
	mkdir -p build/rpm/contra/usr/local/bin
	mkdir -p build/rpm/contra/etc
	cp bin/contra build/rpm/contra/usr/local/bin/
	cp contra.example.conf build/rpm/contra/etc/contra.conf
	fpm --description "Configuration Tracking for Network Devices" --url "https://gitlab.aja.com/go/contra" \
		--license "mit" -m "it@aja.com" -p bin/ -s dir -t rpm -n contra -a x86_64 --rpm-user "contra" --epoch 0 -v $(VERSION) -C build/rpm/contra .

deb64: binaries
	rm -rf build/deb
	mkdir -p build/deb/contra/usr/local/bin
	cp bin/contra build/deb/contra/usr/local/bin/
	mkdir -p build/deb/contra/etc
	cp contra.example.conf build/deb/contra/etc/contra.conf
	fpm --description "Configuration Tracking for Network Devices" --url "https://gitlab.aja.com/go/contra" \
		--license "mit" -m "it@aja.com" -p bin/ -s dir -t deb -n contra -a amd64 -v $(VERSION) -C build/deb/contra .

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
