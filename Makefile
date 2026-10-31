VERSION := $(shell grep '^version=' plugin/module.prop | cut -d= -f2)
BIN     := plugin/system/bin
GOFLAGS := -trimpath -ldflags="-s -w"

.PHONY: all host arm64 arm zip test clean
all: zip

host:   ; go build -o dist/axbridge ./server
arm64:  ; GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build $(GOFLAGS) -o $(BIN)/axbridge-arm64 ./server
arm:    ; GOOS=linux GOARCH=arm GOARM=7 CGO_ENABLED=0 go build $(GOFLAGS) -o $(BIN)/axbridge-arm ./server
test:   ; go vet ./... && go test ./...
zip: arm64 arm ; sh scripts/package.sh $(VERSION)
clean:  ; rm -rf dist $(BIN)/axbridge-arm $(BIN)/axbridge-arm64
