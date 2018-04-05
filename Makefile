VERSION=$(shell ./bin/contra -q -version)
BINARY=contra
RPMDIR=rpm
DEBDIR=deb
all: test binaries

first: deps test run

binaries: staging clean linux64

linux64:
	@echo -----linux64-----
	GOOS=linux GOARCH=amd64 go build -o bin/$(BINARY) contra.go

packages: binaries rpm64 deb64

rpm64: binaries
	@echo -----rpm64-----
	mkdir -p build/$(RPMDIR)/contra/usr/local/bin
	mkdir -p build/$(RPMDIR)/contra/etc/systemd/system
	cp bin/$(BINARY) build/$(RPMDIR)/contra/usr/local/bin/
	cp contra.example.conf build/$(RPMDIR)/contra/etc/contra.conf
	cp files/contra.service build/$(RPMDIR)/contra/etc/systemd/system/contra.service
	fpm --description "Configuration Tracking for Network Devices" --url "https://gitlab.aja.com/go/contra" \
		--license "mit" -m "it@aja.com" -p bin/ -s dir -t rpm -n contra -a x86_64 --epoch 0 -v $(VERSION) -C build/$(RPMDIR)/$(BINARY) .

deb64: binaries
	@echo -----deb64-----
	mkdir -p build/$(DEBDIR)/contra/usr/local/bin
	cp bin/$(BINARY) build/$(DEBDIR)/contra/usr/local/bin/
	mkdir -p build/$(DEBDIR)/contra/etc
	cp contra.example.conf build/$(DEBDIR)/contra/etc/contra.conf
	fpm --description "Configuration Tracking for Network Devices" --url "https://gitlab.aja.com/go/contra" \
		--license "mit" -m "it@aja.com" -p bin/ -s dir -t deb -n contra -a amd64 -v $(VERSION) -C build/$(DEBDIR)/$(BINARY) .

clean:
	@echo -----clean-----
	find bin -name $(BINARY) -type f -exec rm {} \;
	find bin -name '*.rpm' -type f -exec rm {} \;
	find bin -name '*.deb' -type f -exec rm {} \;
	find build -name "$(DEBDIR)" -type d -prune -exec rm -rf "{}" \;
	find build -name "$(RPMDIR)" -type d -prune -exec rm -rf "{}" \;

staging:
	@test -d bin ||mkdir bin
	@test -d build ||mkdir build

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

.PHONY: all clean deps fmt vet test run testrun staging
.PHONY: binaries linux64 first
