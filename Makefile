VERSION=$(shell ./bin/contra -q -version)
BINARY=contra
RPMDIR=rpm
DEBDIR=deb
all: test binaries

binaries: staging clean linux64

linux64:
	@echo -----linux64-----
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/$(BINARY) contra.go

release: binaries compress rpm64 deb64

compress:
	@echo -----compress-----
	upx --brute bin/$(BINARY)

packages: binaries rpm64 deb64

rpm64: binaries
	@echo -----rpm64-----
	mkdir -p build/$(RPMDIR)/contra/usr/local/bin
	mkdir -p build/$(RPMDIR)/contra/etc/systemd/system
	mkdir -p build/$(RPMDIR)/contra/opt/contra/workspace
	cp bin/$(BINARY) build/$(RPMDIR)/contra/usr/local/bin/
	cp contra.example.conf build/$(RPMDIR)/contra/etc/contra.conf.dist
	cp files/rpm/contra.service build/$(RPMDIR)/contra/etc/systemd/system/contra.service
	fpm --description "Configuration Tracking for Network Devices" --url "https://github.com/aja-video/contra" \
		--license "mit" -m "it@aja.com" -d git -p bin/ -s dir -t rpm -n contra -a x86_64 --epoch 0 -v $(VERSION) \
		--before-install files/rpm/pre-install.sh --before-remove files/rpm/pre-remove.sh --after-install \
		files/rpm/post-install.sh --after-remove files/rpm/post-remove.sh --after-upgrade files/rpm/after-upgrade.sh \
		-C build/$(RPMDIR)/$(BINARY) .

deb64: binaries
	@echo -----deb64-----
	mkdir -p build/$(DEBDIR)/contra/usr/local/bin
	cp bin/$(BINARY) build/$(DEBDIR)/contra/usr/local/bin/
	mkdir -p build/$(DEBDIR)/contra/etc
	cp contra.example.conf build/$(DEBDIR)/contra/etc/contra.conf.dist
	fpm --description "Configuration Tracking for Network Devices" --url "https://github.com/aja-video/contra" \
		--license "mit" -m "it@aja.com" -d git -p bin/ -s dir -t deb -n contra -a amd64 -v $(VERSION) -C build/$(DEBDIR)/$(BINARY) .

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

.PHONY: all clean fmt vet test run testrun staging
.PHONY: binaries linux64
