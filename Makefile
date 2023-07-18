NAME=ipsw
BINARY=packer-plugin-$(NAME)

PROJECT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))

.DEFAULT_GOAL = build

# Generate

.PHONY: generate
generate:
	@go generate ./...

# Build & Install

.PHONY: build
build:
	@go build -o $(BINARY)

.PHONY: install
install: build
	@mkdir -p ~/.packer.d/plugins/
	@mv $(BINARY) ~/.packer.d/plugins/$(BINARY)

# Test

test: check
check: plugin-check acceptance-test

acceptance-test: export PKR_VAR_appledb_test_path = $(PROJECT_DIR)/datasource/test-fixtures/
acceptance-test: export PACKER_ACC = 1
acceptance-test: export PACKER_PLUGIN_PATH = $(PROJECT_DIR)
acceptance-test: build
	@go test -count 1 -v $(shell find . | grep acc_test) -timeout=120m

PACKER_SDC := go run github.com/hashicorp/packer-plugin-sdk/cmd/packer-sdc
plugin-check: build
	@$(PACKER_SDC) plugin-check $(BINARY)

# Clean

.PHONY: clean
clean:
	@rm -f $(BINARY) crash.log
