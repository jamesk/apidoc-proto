# Apidoc.me to .proto conveter

Converts from an `apidoc` spec to a `.proto` file. This is still very rough, I would not even call it Alpha.

# Usage

Currently hard coded to use the example `test_service.json`. Simply run the program and it will generate a `.proto` file from that spec.
From the proto spec you can then generate code in the language of your choice.

# Example

To run just the conversion:

`cd cli && go build && ./cli ../test_service.json ../test_service.proto && cd ..`

To get generated gRPC code:

`protoc test_service.proto --go_out=plugins=grpc:/tmp`

You should then find your generated go code in `/tmp/test_service.pb.go`

(See http://www.grpc.io/docs/quickstart/ for `protoc` instructions)


# Noted limitations
Note: This is not an exhaustive list
 * Resources are very opinionated and rough at the moment
 * Imports are not handled
 * Headers are not handled
 * Decimal is not supported as there is not a clear mapping to protobuf types
 * Object is converted to a map of string -> Any as there is no generic object type in proto
 * You cannot use array with Object as proto3's map is used and map cannot be repeated
 * Only service names are currently validated, all others are used raw which could lead to invalid .proto files
 * Currently no sequence number attribute in apidoc spec so user has to maintain order and number of items
 * Can't have nested map/array, don't think this is supported by apidoc, there is no explicit check for this at the moment though
