.DEFAULT_GOAL := build

################################################################################
# PHONY Tasks
################################################################################
.PHONY: build retool rotate_tls_keys setup test

build: $(BUILD_DEPS)
	$(info building binary bin/goraphql_server...)
	@go build -o bin/goraphql_server ./cmd/goraphql_server

retool:
	@if test -z $(shell which retool); then \
		go get github.com/twitchtv/retool; \
		retool add github.com/golang/dep/cmd/dep v0.4.1; \
	fi

setup: retool rotate_tls_keys
	$(info Synching dev tools and dependencies...)
	@retool sync
	@retool do dep ensure

test: $(BUILD_DEPS)
	@go test -race -cover -v ./pkg/...
