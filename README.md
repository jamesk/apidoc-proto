# Apidoc.me to .proto conveter

Converts from an `apidoc` spec to a `.proto` file. This is still very rough, I would not even call it Alpha.

# Usage

Currently hard coded to use the example `test_service.json`. Simply run the program and it will generate a `.proto` file from that spec.
From the proto spec you can then generate code in the language of your choice.

# Example

`go build && ./apidoc-proto && protoc test-apidoc-proto-ticket.proto --go_out=plugins=grpc:/tmp`

You should find your generated go code in `/tmp/test-apidoc-proto-ticket.pb.go`