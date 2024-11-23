#!/bin/bash

set -eu

TAG=$1

export RELEASE_DIR=release-$TAG
mkdir -p "$RELEASE_DIR"

BUILD_ARGS=(-ldflags '-s -w')

export GOOS=darwin
export GOARCH=amd64
go build "${BUILD_ARGS[@]}" -o "$RELEASE_DIR/tmpbbs-$GOOS-$GOARCH"

export GOOS=linux
export GOARCH=amd64
go build "${BUILD_ARGS[@]}" -o "$RELEASE_DIR/tmpbbs-$GOOS-$GOARCH"

export GOARCH=arm
go build "${BUILD_ARGS[@]}" -o "$RELEASE_DIR/tmpbbs-$GOOS-$GOARCH"

export GOARCH=arm64
go build "${BUILD_ARGS[@]}" -o "$RELEASE_DIR/tmpbbs-$GOOS-$GOARCH"

export GOARCH=mips GOMIPS=softfloat
go build "${BUILD_ARGS[@]}" -o "$RELEASE_DIR/tmpbbs-$GOOS-$GOARCH-$GOMIPS"

export GOOS=windows
export GOARCH=amd64
go build "${BUILD_ARGS[@]}" -o "$RELEASE_DIR/tmpbbs-$GOOS-$GOARCH"
