CFLAGS ?= -O1 -g -std=gnu11 -Wall -Wextra -Werror -Wno-shadow -fsanitize=address -fsanitize=undefined -fstack-protector-all -fno-omit-frame-pointer
LDLIBS ?= -lm

CLANG       = clang
GCC_NATIVE  = gcc-15
GCC_DOCKER  = docker run --rm -v "$(shell pwd)":/src -w /src gcc:15.2.0
RISCV64 = docker run --rm --platform linux/riscv64 -v "$(shell pwd)":/src -w /src solod/riscv64
I386 = docker run --rm --platform linux/i386 -v "$(shell pwd)":/src -w /src solod/i386

compiler =
RUN_CMD = ./build/main

# Set CC and CFLAGS based on the selected compiler.
ifeq ($(compiler), clang)
    CC = $(CLANG)
else ifeq ($(compiler), gcc)
    CC = $(GCC_NATIVE)
	CFLAGS += -fanalyzer -D_FORTIFY_SOURCE=2
else ifeq ($(compiler), docker)
    CC = $(GCC_DOCKER) gcc
	CFLAGS += -fanalyzer -D_FORTIFY_SOURCE=2
    RUN_CMD = $(GCC_DOCKER) ./build/main
else ifeq ($(compiler), riscv64)
	CC = $(RISCV64) gcc
	CFLAGS = -O1 -g -std=gnu11 -Wall -Wextra -Werror -Wno-shadow
	RUN_CMD = $(RISCV64) ./build/main
else ifeq ($(compiler), i386)
	CC = $(I386) gcc
	CFLAGS = -O1 -g -std=gnu11 -Wall -Wextra -Werror -Wno-shadow
	RUN_CMD = $(I386) ./build/main
endif

# Preload mimalloc if available.
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
    MIMALLOC_LIB := $(shell ls /opt/homebrew/lib/libmimalloc.dylib /usr/local/lib/libmimalloc.dylib 2>/dev/null | head -1)
    ifneq ($(MIMALLOC_LIB),)
        MIMALLOC_PRELOAD := DYLD_INSERT_LIBRARIES=$(MIMALLOC_LIB)
    endif
else ifeq ($(UNAME_S),Linux)
    MIMALLOC_LIB := $(shell ls /usr/lib/libmimalloc.so /usr/local/lib/libmimalloc.so 2>/dev/null | head -1)
    ifneq ($(MIMALLOC_LIB),)
        MIMALLOC_PRELOAD := LD_PRELOAD=$(MIMALLOC_LIB)
    endif
endif

inspect:
	go run ./cmd/inspect -- $(path)

test:
	@go test ./so/...
	@go test ./internal/...

prepare-riscv64:
	@printf 'FROM alpine:edge\nRUN apk add --no-cache gcc musl-dev\n' \
		| docker build --platform=linux/riscv64 -t solod/riscv64 -
	@docker run --rm -it --platform=linux/riscv64 -v $(shell pwd):/src solod/riscv64 uname -m

prepare-i386:
	@printf 'FROM alpine:edge\nRUN apk add --no-cache gcc musl-dev\n' \
		| docker build --platform=linux/i386 -t solod/i386 -
	@docker run --rm -it --platform=linux/i386 -v $(shell pwd):/src solod/i386 uname -m

update-dst:
	make run-case name=$(name)
	cp generated/$(name)/main.* testdata/$(name)/dst
	go test -run TestTranslate/$(name) ./internal/compiler

run-cases:
	@make run-cases-by pattern="testdata/lang/*/ testdata/std/*/"

run-cases-windows:
	@make run-cases-by CFLAGS="-g -std=gnu11 -Wall -Wextra -Werror -Wno-shadow -lm" pattern="testdata/lang/*/"

run-cases-by:
	@failed=0; \
	for dir in $(pattern); do \
		name=$${dir#testdata/}; \
		name=$${name%/}; \
		if make run-case name=$$name > /tmp/so_test_out.txt 2>&1; then \
			echo "PASS $$name"; \
		else \
			echo "FAIL $$name"; \
			cat /tmp/so_test_out.txt; \
			failed=1; \
		fi; \
	done; \
	rm -f /tmp/so_test_out.txt; \
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
	@$(CC) $(CFLAGS) \
		-I$(path) \
		-o build/main \
		$(shell find $(path) -name "*.c") \
		$(LDLIBS)
	@$(RUN_CMD)
	@rm -f build/main

.PHONY: bench
bench:
	@cd bench/$(name) && go test -bench=.
	@CFLAGS="-Ofast -march=native -flto -funroll-loops -DNDEBUG" \
	$(MIMALLOC_PRELOAD) \
	go run ./cmd/so run ./bench/$(name)
