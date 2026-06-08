#!/bin/bash
source ./_helpers.sh

ARCH=$(uname -m | sed -e 's/x86_64/amd64/' -e 's/aarch64/arm64/')

MAKEFILE_OUTPUT_DIR=../build/

CHEAT_CODEX_ETC=/etc/cheat-codex
CHEAT_CODEX_OPT=/opt/cheat-codex
CHEAT_CODEX_BIN=$CHEAT_CODEX_OPT/bin


CHEAT_CODEX_SRC=$CHEAT_CODEX_OPT/src/
CHEAT_CODEX_MAPS=$CHEAT_CODEX_OPT/maps/

RUST_SRC_DIR=../core/
SRC_CORE_BIN=$MAKEFILE_OUTPUT_DIR/cheat-codex-core
DST_CORE_BIN=$CHEAT_CODEX_BIN/cheat-codex-core
CORE_BIN_LINK=/usr/local/bin/cheat-codex-core

UI_SRC_DIR=../ui/
SRC_TUI_BIN=$MAKEFILE_OUTPUT_DIR/cheat-codex-tui_linux-$ARCH
DST_TUI_BIN=$CHEAT_CODEX_BIN/cheat-codex-tui
TUI_BIN_LINK=/usr/local/bin/cheat-codex-tui

SRC_MAPS=../maps/

### Core Functonality ###
function remove_previous_installation() {
    sudo rm -rf $CHEAT_CODEX_SRC
    sudo rm -rf $CHEAT_CODEX_MAPS
    sudo rm -rf $CHEAT_CODEX_BIN
    sudo rm -rf $CORE_BIN_LINK
    sudo rm -rf $TUI_BIN_LINK
}

function prep_cheat_codex_dirs() {
    start_step_message "Prepping cheat-codex Directories"
    _create_dir $CHEAT_CODEX_ETC
    _create_dir $CHEAT_CODEX_SRC
    _create_dir $CHEAT_CODEX_BIN
    _create_dir $CHEAT_CODEX_MAPS
    successful
}
function _create_dir() {
    if [ ! -d "$1" ]; then
        start_step_message "$1" "substep"
        if ! sudo mkdir -p "$1"; then
            error_message "Failed to create directory '$1'"
        fi
    fi
}

function build_binaries() {
    cd .. 
    if ! make; then
        error_message "Failed to build cheat-codex binaries"
    fi
    cd scripts
}

function move_files() {
    start_step_message "Moving cheat-codex Files"
    _move_file $RUST_SRC_DIR $CHEAT_CODEX_SRC/core
    _move_file $UI_SRC_DIR $CHEAT_CODEX_SRC/ui
    _move_file $SRC_CORE_BIN $DST_CORE_BIN
    sudo chmod +x $DST_CORE_BIN

    _move_file $SRC_TUI_BIN $DST_TUI_BIN
    sudo chmod +x $DST_TUI_BIN

    for FILE in $(ls $SRC_MAPS); do
        _move_file $SRC_MAPS$FILE $CHEAT_CODEX_MAPS
    done
    successful
}
function _move_file() {
    start_step_message "$1 -> $2" "substep"
    if ! sudo cp -r $1 $2; then
        error_message "Failed to move $1 to $2"
    fi
}

function create_links() {
    start_step_message "Creating Links"
    if ! sudo ln -s $DST_CORE_BIN $CORE_BIN_LINK; then
        error_message "Failed to 'ln -s ${DST_CORE_BIN} ${CORE_BIN_LINK}'"
    fi

    if ! sudo ln -s $DST_TUI_BIN $TUI_BIN_LINK; then
        error_message "Failed to 'ln -s ${DST_TUI_BIN} ${TUI_BIN_LINK}'"
    fi
    successful
}

function main() {
    remove_previous_installation
    prep_cheat_codex_dirs
    build_binaries
    move_files
    create_links
}
main

