PKGS := $(shell go list ./... | grep -v /vendor)

BIN_DIR := $(GOPATH)/bin

BINARY := presla

VERSION ?= vlatest
PLATFORMS := windows linux darwin
os = $(word 1, $@)

all: deps vet test bindata release compress bindata-debug done

vet:
	@echo "=> Running go vet, please check if there is output..."
	@go vet ./...

compress:
	@echo "=> Compressing binaries for version $(VERSION)"
	@find ./release -iname "*$(VERSION)*" -exec upx {} -o {}-compressed \; && mv release/presla-"$(VERSION)"-windows-amd64.exe-compressed release/presla-"$(VERSION)"-windows-amd64-compressed.exe

format:
	@echo "=> Running go fmt..."
	@go fmt gitlab.com/3stadt/...

test: bindata-debug
	@echo "=> Running tests..."
	@go test -v ./...

run: bindata-debug
	@echo "=> Starting Server..."
	@go run main.go

deps:
	@echo "=> Installing dependencies, this may take a while..."
	@go get -u github.com/golang/dep/cmd/dep
	@"$(BIN_DIR)"/dep version
	@"$(BIN_DIR)"/dep ensure

done:
	@echo "=> Done"

bindata-debug:
	@echo "=> Generating binary data for development..."
	@go get -u github.com/jteeuwen/go-bindata/...
	@rm -f src/Handlers/bindata.go
	@"$(BIN_DIR)"/go-bindata -o src/Handlers/bindata.go -debug -ignore=.*-inkscape\.svg -pkg Handlers templates/... static/... executors/...

bindata:
	@echo "=> Generating binary data for production..."
	@go get -u github.com/jteeuwen/go-bindata/...
	@rm -f src/Handlers/bindata.go
	@"$(BIN_DIR)"/go-bindata -o src/Handlers/bindata.go -ignore=.*-inkscape\.svg -pkg Handlers templates/... static/... executors/...

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	@echo "=> Creating release for $(os)..."
	@mkdir -p release/
	@GOOS=$(os) GOARCH=amd64 go build -o release/$(BINARY)-$(VERSION)-$(os)-amd64
	@if [ "$(os)" == "windows" ]; then mv release/"$(BINARY)"-"$(VERSION)"-"$(os)"-amd64 release/"$(BINARY)"-"$(VERSION)"-"$(os)"-amd64.exe; fi;
	@echo "> Created $(os)-$(VERSION).tar.gz"

.PHONY: release
release: windows linux darwin