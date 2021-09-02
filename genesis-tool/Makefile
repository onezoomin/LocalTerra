export GO111MODULE = on

BUILDDIR ?= $(CURDIR)/build

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILDDIR)/
GOARCH=amd64 $(MAKE) build

$(BUILD_TARGETS): go.sum $(BUILDDIR)/
	go $@ -mod=readonly $(BUILD_FLAGS) $(BUILD_ARGS) ./...

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

.PHONY: build 