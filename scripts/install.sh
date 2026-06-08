#!/bin/bash
source ./_helpers.sh

ARCH = $(uname -m | sed -e 's/x86_64/amd64/' -e 's/aarch64/arm64/')

CURRENT_DIR = $(pwd)
MAKEFILE_OUTPUT_DIR = "${CURRENT_DIR}/../build/"

CHEAT_CODEX_ETC = /etc/cheat-codex/
CHEAT_CODEX_OPT = /opt/cheat-codex/
CHEAT_CODEX_BIN = "${CHEAT_CODEX_OPT}bin/"
DST_CORE_BIN = "${CHEAT_CODEX_BIN}cheat-codex-core"
DST_TUI_BIN = "${CHEAT_CODEX_BIN}cheat-codex-tui"
CHEAT_CODEX_SRC = "${CHEAT_CODEX_OPT}src/"
CHEAT_CODEX_MAPS = "${CHEAT_CODEX_OPT}maps/"

RUST_SRC_DIR = "${CURRENT_DIR}/../core/"
SRC_CORE_BIN = "${MAKEFILE_OUTPUT_DIR}cheat-codex-core"

UI_SRC_DIR = "${CURRENT_DIR}/../ui/"
SRC_TUI_BIN = "${MAKEFILE_OUTPUT_DIR}cheat-codex-tui_linux-${ARCH}"

SRC_MAPS = "${CURRENT_DIR}/../maps/*"

### Core Functonality ###
function prep_cheat_codex_dirs() {
    start_step_message "Prepping cheat-codex Directories"
    mkdir -p $CHEAT_CODEX_ETC
    mkdir -p $CHEAT_CODEX_SRC
    mkdir -p $CHEAT_CODEX_BIN
    mkdir -p $CHEAT_CODEX_MAPS
    successful
}

function build_binaries() {
    if ! cd .. && make; then
        error_message "Failed to build cheat-codex binaries"
    fi
}

function move_file() {
    start_step_message "Moving cheat-codex Files"
    cp -r $RUST_SRC_DIR "${CHEAT_CODEX_SRC}core"
    cp -r $UI_SRC_DIR "${CHEAT_CODEX_SRC}ui"
    cp $SRC_CORE_BIN $DST_CORE_BIN
    cp $SRC_TUI_BIN $DST_TUI_BIN
    cp -r $SRC_MAPS $CHEAT_CODEX_MAPS
    successful
}

function main() {
    prep_cheat_codex_dirs
    build_binaries
    move_files
}
main

