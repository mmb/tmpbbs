#!/bin/bash

set -eu

VERSION=$1
COMMIT=$2

export RELEASE_DIR=release
mkdir -p "$RELEASE_DIR"

BUILD_ARGS=(-ldflags "-s -w -X github.com/mmb/tmpbbs/internal/tmpbbs.Version=$VERSION -X github.com/mmb/tmpbbs/internal/tmpbbs.Commit=$COMMIT")

build() {
  export GOOS=$1 GOARCH=$2 CGO_ENABLED=$3

  go build "${BUILD_ARGS[@]}" -o "$RELEASE_DIR/tmpbbs-$VERSION-$GOOS-$GOARCH"
}

build android arm64 0
build darwin amd64 0
build darwin arm64 0
build linux 386 0
build linux amd64 0
build linux arm 0
build linux arm64 0
build windows 386 0
build windows amd64 0

export GOOS=linux GOARCH=mips GOMIPS=softfloat CGO_ENABLED=0
go build "${BUILD_ARGS[@]}" -o "$RELEASE_DIR/tmpbbs-$VERSION-$GOOS-$GOARCH-$GOMIPS"
