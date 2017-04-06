# Apidoc.me to .proto conveter

Converts from an `apidoc` spec to a `.proto` file. This is still very rough, I would not even call it Alpha.

# Usage

Currently hard coded to use the example `test_service.json`. Simply run the program and it will generate a `.proto` file from that spec.
From the proto spec you can then generate code in the language of your choice.

# Example

To run just the conversion:

`go build && ./apidoc-proto test_service.json test_service.proto`

To get generated gRPC code:

`protoc test_service.proto --go_out=plugins=grpc:/tmp`

You should then find your generated go code in `/tmp/test_service.pb.go`

(See http://www.grpc.io/docs/quickstart/ for `protoc` instructions)

