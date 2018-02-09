PKGS := $(shell go list ./... | grep -v /vendor)

BIN_DIR := $(GOPATH)/bin

BINARY := presla

VERSION ?= vlatest
PLATFORMS := windows linux darwin
os = $(word 1, $@)

all: deps bindata release done

format:
	go fmt gitlab.com/3stadt/...

run:
	@echo "=> Starting Server..."
	go run main.go

deps:
	@echo "=> Installing dependencies..."
	@go get -u github.com/golang/dep/cmd/dep
	@$(BIN_DIR)/dep ensure

done:
	@echo "=> Done"

bindata:
	@echo "=> Generating binary data..."
	@go get -u github.com/jteeuwen/go-bindata/...
	@rm -f src/Handlers/bindata.go
	@go-bindata -o src/Handlers/bindata.go -ignore=.*-inkscape\.svg -pkg Handlers templates/... static/...

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	@echo "=> Creating release for $(os)..."
	@mkdir -p release/
	@GOOS=$(os) GOARCH=amd64 go build -o release/$(BINARY)-$(VERSION)-$(os)-amd64
	@if [ "$(os)" == "windows" ]; then mv release/"$(BINARY)"-"$(VERSION)"-"$(os)"-amd64 release/"$(BINARY)"-"$(VERSION)"-"$(os)"-amd64.exe; fi;
	@echo "> Created $(os)-$(VERSION).tar.gz"

.PHONY: release
release: windows linux darwin