NAME=ipsw
BINARY=packer-plugin-$(NAME)

PROJECT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
.DEFAULT_GOAL = build
PACKER_SDC := go run github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc

# Generate

HCL2_SOURCES := $(shell grep -r -l "go:generate.*mapstructure-to-hcl2" **/*.go)
HCL2_GENERATED = $(HCL2_SOURCES:.go=.hcl2spec.go)
$(HCL2_GENERATED): %.hcl2spec.go : %.go
	@go generate -run="-command|hcl2" $<

# Build

.PHONY: build
build: $(HCL2_GENERATED)
	@go build -o $(BINARY)

PACKER_PLUGIN_PATH ?= $(HOME)/.config/packer/plugins
PLUGIN_INSTALL_PATH = $(PACKER_PLUGIN_PATH)/github.com/torarnv/$(NAME)
PLUGIN_VERSION=$(shell git describe --tags HEAD | sed 's/v.*-g.*/v1337.0.0/')
PLUGIN_INSTALL_FILEPATH = $(PLUGIN_INSTALL_PATH)/$(BINARY)_$(PLUGIN_VERSION)_$(API_VERSION)_$(GOOS)_$(GOARCH)

.PHONY: install
install: export GOOS = darwin
install: export GOARCH = $(shell arch)
install: export API_VERSION = 5.0
install: build
	@mkdir -p $(PLUGIN_INSTALL_PATH)
	@cp -c $(BINARY) $(PLUGIN_INSTALL_FILEPATH)
	@shasum -a 256 $(PLUGIN_INSTALL_FILEPATH) | \
		cut -d ' ' -f 1 > $(PLUGIN_INSTALL_FILEPATH)_SHA256SUM

# Test

test: check
check: plugin-check acceptance-test

acceptance-test: export PKR_VAR_appledb_test_path = $(PROJECT_DIR)/datasource/test-fixtures/
acceptance-test: export PACKER_ACC = 1
acceptance-test: export PACKER_PLUGIN_PATH = $(PROJECT_DIR)
acceptance-test: build
	@go test -count 1 -v $(shell find . | grep _test.go) -timeout=120m

plugin-check: build
	@$(PACKER_SDC) plugin-check $(BINARY)

# Clean

.PHONY: clean
clean:
	@rm -Rf $(BINARY) $(HCL2_GENERATED) $(DOC_GENERATED) crash.log build/ dist/ docs.zip

# Docs

DOC_SOURCES := $(shell grep -r -l "go:generate.*struct-markdown" **/*.go)
DOC_GENERATED = docs-partials
$(DOC_GENERATED): $(DOC_SOURCES)
	@go generate -run="-command|markdown" $(DOC_SOURCES)

.PHONY: docs
docs: $(DOC_GENERATED)
	@rm -Rf build/docs/
	@$(PACKER_SDC) renderdocs -src docs/ -partials docs-partials/ -dst build/docs/
	@cp README.md build/docs/
	@./.web-docs/scripts/compile-to-webdocs.sh "." "build/docs/" ".web-docs" "<orgname>"
