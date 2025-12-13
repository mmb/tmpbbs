// Package proto contains protobuf and grpc code generated from .proto files.
package proto

//go:generate protoc --go-grpc_opt=paths=source_relative --go-grpc_out=. --go_opt=paths=source_relative --go_out=.  post_sync.proto
