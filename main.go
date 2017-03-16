package main

import (
	"fmt"
	"io/ioutil"

	"github.com/golang/protobuf/proto"
	proto2 "github.com/golang/protobuf/protoc-gen-go/descriptor"
	"github.com/jamesk/apidoc-proto/apidoc"
)

func main() {
	aSpec, err := apidoc.GetSpecFromFile("test_service.json")
	if err != nil {
		panic(err)
	}

	pFile, err := getProtoFromAPISpec(aSpec)
	if err != nil {
		panic(err)
	}

	fmt.Println(proto.MarshalTextString(&pFile))

	ioutil.WriteFile("test_service.proto", []byte(proto.MarshalTextString(&pFile)), 0666)
}

func getProtoFromAPISpec(spec apidoc.Spec) (proto2.FileDescriptorProto, error) {
	pFile := proto2.FileDescriptorProto{}

	syntax := "proto3"
	pFile.Syntax = &syntax
	pFile.Name = &spec.Name
	for _, aEnum := range spec.Enums {
		pEnum := proto2.EnumDescriptorProto{}

		pEnum.Name = &aEnum.Name
		for i, aValue := range aEnum.Values {
			pEnumValue := proto2.EnumValueDescriptorProto{}
			pEnumValue.Name = &aValue.Name

			enumIndex := int32(i)
			pEnumValue.Number = &enumIndex

			pEnum.Value = append(pEnum.Value, &pEnumValue)
		}


		pFile.EnumType = append(pFile.EnumType, &pEnum)
	}
//rpc SendTest (TestRequest) returns (TestReply) {}
	m := proto2.MethodDescriptorProto{}
	mName := "aFunc"
	mInputType := "int32"
	mOutputType := "int32"
	m.Name = &mName
	m.InputType = &mInputType
	m.OutputType = &mOutputType

	s := proto2.ServiceDescriptorProto{}
	sName := "MyService"
	s.Name = &sName
	s.Method = []*proto2.MethodDescriptorProto{&m}
	pFile.Service = []*proto2.ServiceDescriptorProto{&s}

	return pFile, nil
}

