// Package proto contains protobuf and grpc code generated from .proto files.
package proto

//go:generate go install google.golang.org/grpc/cmd/protoc-gen-go-grpc
//go:generate go install google.golang.org/protobuf/cmd/protoc-gen-go
//go:generate protoc --go-grpc_opt=paths=source_relative --go-grpc_out=. --go_opt=paths=source_relative --go_out=. post_sync.proto
