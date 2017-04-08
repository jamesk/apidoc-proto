package main

import (
	"encoding/json"

	"errors"

	"github.com/emicklei/proto"
	"github.com/jamesk/apidoc-proto/apidoc"
)

type RpcParameter struct {
	IsStream bool
	RpcType  string
	Message  proto.Message
}

func getRpcParameter(aType string, attributes []apidoc.Attribute) (RpcParameter, error) {
	if len(aType) == 0 {
		param := RpcParameter{}
		param.RpcType = "google.protobuf.Empty"
		return param, nil
	}

	if param, ok, err := getRpcMapParam(aType); err != nil {
		return param, err
	} else if ok {
		return param, err
	}

	if param, ok, err := getRPCArrayParam(aType, attributes); err != nil {
		return param, err
	} else if ok {
		return param, err
	}

	param := RpcParameter{}
	param.RpcType = getProtoTypeFromBasicApidocType(aType)
	return param, nil
}

func getRpcMapParam(aType string) (RpcParameter, bool, error) {
	param := RpcParameter{}
	mapType, isMap, err := getMapType(aType)
	if err != nil {
		return param, false, err
	}
	if !isMap {
		return param, false, nil
	}

	//TODO: name clashes?
	messageName := mapType + "Map"
	pMessage := proto.Message{Name: messageName}
	pField, err := getMapProtoField("value", mapType, 1)
	if err != nil {
		return param, false, err
	}
	pMessage.Elements = append(pMessage.Elements, pField)

	param.RpcType = messageName
	param.Message = pMessage

	return param, true, nil
}

func getRPCArrayParam(aType string, attributes []apidoc.Attribute) (RpcParameter, bool, error) {
	param := RpcParameter{}

	var arrayAsStream ArrayAsStreamAttribute
	for _, attr := range attributes {
		if attr.Name == ArrayAsStreamAttributeName {
			//TODO: inefficient, any actual performance impact?
			d, err := json.Marshal(attr.Value)
			if err != nil {
				return param, false, err
			}

			json.Unmarshal(d, &arrayAsStream)
		}
	}

	arrayType, isArray, err := getArrayType(aType)
	if err != nil {
		return param, false, err
	}
	if !isArray {
		if arrayAsStream.Value {
			return param, false, errors.New("Array as stream attribute set but body is not an array")
		}
		return param, false, nil
	}

	if arrayAsStream.Value {
		param.IsStream = true
		param.RpcType = getProtoTypeFromBasicApidocType(arrayType)

		return param, true, nil
	}

	//Implicitly create an array message type for the given type
	//TODO: handle name clashes?
	arrayMessageName := arrayType + "Array"
	pMessage := proto.Message{Name: arrayMessageName}
	pField := proto.NormalField{Field: &proto.Field{}}
	pField.Name = "value"
	pField.Repeated = true
	pField.Sequence = 1
	pField.Type = arrayType

	pMessage.Elements = append(pMessage.Elements, &pField)
	param.Message = pMessage
	param.RpcType = getProtoTypeFromBasicApidocType(arrayMessageName)

	return param, true, nil

}
