.PHONY: build install test clean

BINARY_NAME=terraform-provider-circleci-server
VERSION=1.0.0
OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)
INSTALL_DIR=$(HOME)/.terraform.d/plugins/registry.terraform.io/anduril/circleci-server/$(VERSION)/$(OS_ARCH)

build:
	go build -o $(BINARY_NAME)

install: build
	@echo "Installing provider to $(INSTALL_DIR)"
	@mkdir -p $(INSTALL_DIR)
	@cp $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Provider installed successfully!"
	@echo "You can now use it in your Terraform configs"

test:
	go test ./... -v

clean:
	rm -f $(BINARY_NAME)
	rm -rf $(INSTALL_DIR)

fmt:
	go fmt ./...

.DEFAULT_GOAL := build
