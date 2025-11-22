#!/bin/bash

set -eu

VERSION=$1
COMMIT=$2

export RELEASE_DIR=release
mkdir -p "$RELEASE_DIR"

export CGO_ENABLED=0

BUILD_ARGS=(-ldflags "-s -w -X github.com/mmb/tmpbbs/internal/tmpbbs.Version=$VERSION -X github.com/mmb/tmpbbs/internal/tmpbbs.Commit=$COMMIT")

build() {
  export GOOS=$1 GOARCH=$2

  go build "${BUILD_ARGS[@]}" -o "$RELEASE_DIR/tmpbbs-$VERSION-$GOOS-$GOARCH"
}

build android arm64
build darwin amd64
build darwin arm64
build linux 386
build linux amd64
build linux arm
build linux arm64
build windows 386
build windows amd64

export GOOS=linux GOARCH=mips GOMIPS=softfloat
go build "${BUILD_ARGS[@]}" -o "$RELEASE_DIR/tmpbbs-$VERSION-$GOOS-$GOARCH-$GOMIPS"
