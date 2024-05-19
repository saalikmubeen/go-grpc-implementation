//go:build tools
// +build tools

package tools

//  go:build tools

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)

// List of blank imports of the protoc plugins. The reason we
// are doing this is because we are not using them directly in
// our code. We are using them as plugins for the protoc compiler.
// We just want to install them to our local machine, so that
// protoc can use them to generate the gRPC gateway code for us.
// We Use a tool dependency package to track the versions of the above executable packages:

// Next step, we run this 'go install' command to install all binaries
// of the plugins in the bin folder of the go path. Those binaries will
// be used by protoc to generate the gRPC gateway code for us:

// go install \
//     github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
//     github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
//     google.golang.org/protobuf/cmd/protoc-gen-go \
//     google.golang.org/grpc/cmd/protoc-gen-go-grpc

// Full documentation and usage:
// https://github.com/grpc-ecosystem/grpc-gateway#usage
