#!/usr/bin/env bash
set -x

cleanup() {
    rm -f overseer
}

trap cleanup EXIT

go build .
if [[ $? -eq 0 ]]; then
    ./overseer \
	-logtostderr
fi
