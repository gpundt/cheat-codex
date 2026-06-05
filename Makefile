SHELL := /bin/bash
CURRENT_DIR := $(dir $(realpath $(lastword $(MAKEFILE_LIST))))
RUST_DIR := $(CURRENT_DIR)core/
RUST_BUILD_OUTPUT_DIR := $(RUST_DIR)target/release/
RUST_BUILD_OUTPUT_FILE := $(RUST_BUILD_OUTPUT_DIR)cheat-codex-core
RUST_DST_FILE := cheat-codex-core

TUI_DIR := $(CURRENT_DIR)tui/
TUI_BUILD_OUTPUT_DIR := $(TUI_DIR)build/
TUI_BUILD_OUTPUT_FILE := $(TUI_BUILD_OUTPUT_DIR)cheat-codex
TUI_DST_FILE := /usr/bin/cheat-codex

## Colors ##
RED     := \033[0;31m
GREEN   := \033[0;32m
YELLOW  := \033[0;33m
BLUE    := \033[0;34m
CYAN	:= \033[0;36m
RESET   := \033[0m
define start_step_message
	@echo -e "\n$(CYAN)[*] $(1) [*]$(RESET)"
endef
define error_message
	@echo "$(RED)ERROR$(RESET): $(1)"
endef
define successful
	@echo -e "\t - $(GREEN)*Successful*$(RESET)\n"
endef

all: prep_build_dirs build_tui

prep_build_dirs:
	$(call start_step_message,"Creating Build Output Dirs")
	@mkdir -p $(TUI_BUILD_OUTPUT_DIR)
	$(call successful)

build_memory_bins:
	$(call start_step_message,"Building Core Memory Operation Binaries '$(RUST_DIR)'")
	@cd $(RUST_DIR) && cargo build --release
	$(call successful)

build_tui: build_memory_bins
	$(call start_step_message,"Building TUI '$(TUI_DIR)'")
	@cd $(TUI_DIR) && \
	go mod tidy && go mod vendor && \
	go build -mod=vendor -ldflags="-s -w" -o $(TUI_BUILD_OUTPUT_FILE) cmd 

help:							## Displays available make targets
	@egrep -h '\s##\s' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "$(BLUE)  %-30s$(RESET) %s\n", $$1,