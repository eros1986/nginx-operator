GODIR = $(shell go list ./... | grep -v /vendor/)
PKG := github.com/eros1986/nginx-operator
GOARCH := amd64
GOOS := linux
BUILD_IMAGE ?= golang:1.9.0-alpine

pre-build:
	@echo "pre build"
	@echo "clean all flycheck files"
	@find . -name "flycheck*" | xargs rm -f
.PHONY: pre-build

build-dirs: pre-build
	@mkdir -p .go/src/$(PKG) ./go/bin
	@mkdir -p release
.PHONY: build-dirs

build-operator: build-dirs
	@docker run                                                            \
	    --rm                                                               \
	    -ti                                                                \
	    -u $$(id -u):$$(id -g)                                             \
	    -v $$(pwd)/.go:/go                                                 \
	    -v $$(pwd):/go/src/$(PKG)                                          \
	    -v $$(pwd)/release:/go/bin                                         \
	    -e GOOS=$(GOOS)                                                    \
	    -e GOARCH=$(GOARCH)                                                \
	    -e CGO_ENABLED=0                                                   \
	    -w /go/src/$(PKG)                                                  \
	    $(BUILD_IMAGE)                                                     \
	    go install -v -pkgdir /go/pkg ./cmd/operator
.PHONY: build-operator

image-operator: build-operator
	@sh build/operator.sh
.PHONY: image-operator
