#!/usr/bin/env bash

set -ex

source lib.sh

main() {
    build-client $1
}


main $@