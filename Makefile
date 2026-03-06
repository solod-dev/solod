CFLAGS = -g -std=gnu11 -Wall -Wextra -Werror -Wno-shadow -fsanitize=address -fsanitize=undefined -fstack-protector-all -lm

CLANG       = clang
GCC_NATIVE  = gcc-15
GCC_DOCKER  = docker run --rm -v "$(shell pwd)":/src -w /src gcc:15.2.0

compiler ?= $(СС)
RUN_CMD = ./build/main

ifeq ($(compiler), clang)
    CC = $(CLANG)
else ifeq ($(compiler), gcc)
    CC = $(GCC_NATIVE)
else ifeq ($(compiler), docker)
    CC = $(GCC_DOCKER) gcc
    RUN_CMD = $(GCC_DOCKER) ./build/main
endif

# --- Targets ---

inspect:
	go run ./cmd/inspect -- $(path)

test:
	@go test ./so/...
	@go test ./internal/...

dist:
	@rm -rf dist
	@mkdir -p dist/solod/bin
	@go build -o dist/solod/bin/so ./cmd/so
	@tar -czf dist/solod.tar.gz -C dist solod
	@echo "Created dist/solod.tar.gz"

run-cases:
	@failed=0; \
	for dir in testdata/*/; do \
		name=$$(basename $$dir); \
		if make run-case name=$$name > /dev/null 2>&1; then \
			echo "PASS $$name"; \
		else \
			echo "FAIL $$name"; \
			failed=1; \
		fi; \
	done; \
	if [ $$failed -eq 0 ]; then \
		echo "PASS"; \
	else \
		echo "FAIL"; \
		exit 1; \
	fi

run-case:
	@rm -rf generated/$(name)
	@mkdir -p generated/$(name)
	@cp testdata/$(name)/dst/*.ext.[ch] generated/$(name)/ 2>/dev/null || true
	@go run ./cmd/so translate -o generated/$(name) testdata/$(name)/src
	@make run-c path=generated/$(name)

run-c:
	@mkdir -p build
	@$(CC) $(CFLAGS) -I$(path) -o build/main $(shell find $(path) -name "*.c")
	@$(RUN_CMD)
	@rm -f build/main
