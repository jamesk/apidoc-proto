package main

import (
	"errors"

	"encoding/json"

	"github.com/emicklei/proto"
	"github.com/jamesk/apidoc-proto/apidoc"
)

type RpcParameter struct {
	IsStream bool
	RpcType  string
	Message  proto.Message
}

func getRpcParameter(aType string, attributes []apidoc.Attribute) (RpcParameter, error) {
	var param RpcParameter

	var arrayAsStream ArrayAsStreamAttribute
	for _, attr := range attributes {
		if attr.Name == ArrayAsStreamAttributeName {
			//TODO: inefficient, any actual performance impact?
			d, err := json.Marshal(attr.Value)
			if err != nil {
				panic(err)
			}

			json.Unmarshal(d, &arrayAsStream)
		}
	}

	if len(aType) != 0 {
		_, isMap, err := getMapType(aType)
		if err != nil {
			panic(err)
		}

		if isMap {
			panic("Map bodies is unsupported at the moment")
		}

		arrayType, isArray, err := getArrayType(aType)
		if err != nil {
			panic(err)
		}

		if isArray {
			if arrayAsStream.Value {
				param.IsStream = true
				param.RpcType = getProtoTypeFromBasicApidocType(arrayType)
			} else {
				//Implicitly create an array message type for the given type
				//TODO: handle name clashes?
				arrayMessageName := arrayType + "Array"
				pMessage := proto.Message{Name: arrayMessageName}
				pField := proto.NormalField{Field: &proto.Field{}}
				pField.Name = "value"
				pField.Repeated = true
				pField.Sequence = 1
				pField.Type = arrayType

				pMessage.Elements = append(
					pMessage.Elements,
					&pField,
				)
				param.Message = pMessage
				param.RpcType = getProtoTypeFromBasicApidocType(arrayMessageName)
			}
		} else {
			if arrayAsStream.Value {
				return param, errors.New("Array as stream attribute set but body is not an array")
			}

			param.RpcType = getProtoTypeFromBasicApidocType(aType)
		}
	} else {
		param.RpcType = "google.protobuf.Empty"
	}

	return param, nil
}
