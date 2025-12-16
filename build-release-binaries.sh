#!/bin/bash

set -eu

VERSION=$1
COMMIT=$2
DATE=$(date -u +%Y-%m-%dT%H:%M:%SZ)

export RELEASE_DIR=release
mkdir -p "$RELEASE_DIR"

export CGO_ENABLED=0

BUILD_ARGS=(-ldflags "-s -w -X main.version=$VERSION -X main.commit=$COMMIT -X main.date=$DATE")

build() {
  export GOOS=$1 GOARCH=$2

  go build "${BUILD_ARGS[@]}" -o "$RELEASE_DIR/tmpbbs-$VERSION-$GOOS-$GOARCH"
}

build_386() {
  export GOOS=$1 GOARCH=386 GO386=$2

  go build "${BUILD_ARGS[@]}" -o "$RELEASE_DIR/tmpbbs-$VERSION-$GOOS-$GOARCH-$GO386"
}

build_arm() {
  export GOOS=$1 GOARCH=arm GOARM=$2

  go build "${BUILD_ARGS[@]}" -o "$RELEASE_DIR/tmpbbs-$VERSION-$GOOS-$GOARCH-$GOARM"
}

build_mips() {
  export GOOS=$1 GOARCH=$2 GOMIPS=$3

  go build "${BUILD_ARGS[@]}" -o "$RELEASE_DIR/tmpbbs-$VERSION-$GOOS-$GOARCH-$GOMIPS"
}

build android arm64
build darwin amd64
build darwin arm64
build linux amd64
build linux arm64
build windows amd64
build windows arm64
build_386 linux softfloat
build_386 linux sse2
build_386 windows softfloat
build_386 windows sse2
build_arm linux 5
build_arm linux 6
build_arm linux 7
build_mips linux mips hardfloat
build_mips linux mips softfloat
build_mips linux mips64 hardfloat
build_mips linux mips64 softfloat
build_mips linux mips64le hardfloat
build_mips linux mips64le softfloat
build_mips linux mipsle hardfloat
build_mips linux mipsle softfloat
