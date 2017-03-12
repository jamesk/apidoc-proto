package main

import (
	"io/ioutil"

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

	ioutil.WriteFile("test_service.proto", []byte(pFile.String()), 0)
}

func getProtoFromAPISpec(spec apidoc.Spec) (proto2.FileDescriptorProto, error) {
	pFile := proto2.FileDescriptorProto{}

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

	return pFile, nil
}

