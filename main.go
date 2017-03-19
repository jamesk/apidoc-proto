package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/emicklei/proto"
	"github.com/jamesk/apidoc-proto/apidoc"
)

func main() {
	aSpec, err := apidoc.GetSpecFromFile("test_service.json")
	if err != nil {
		panic(err)
	}

	pFile := getProtoFromAPISpec(aSpec)

	//TODO: safe name?
	f, err := os.Create(aSpec.Name + ".proto")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	proto.NewFormatter(w, "  ").Format(&pFile)
	w.Flush()
}

const (
	syntaxVersion = "proto3"
)

func getProtoFromAPISpec(spec apidoc.Spec) proto.Proto {
	pFile := proto.Proto{}

	pFile.Elements = append(
		pFile.Elements,
		&proto.Syntax{Value: syntaxVersion},
	)

	for _, aEnum := range spec.Enums {
		pEnum := proto.Enum{Name: aEnum.Name}
		for i, aValue := range aEnum.Values {

			pEnum.Elements = append(
				pEnum.Elements,
				&proto.EnumField{Name:aValue.Name, Integer:i},
			)
		}
		pFile.Elements = append(pFile.Elements, &pEnum)
	}

	//Models
	for _, aModel := range spec.Models {
		pMessage := proto.Message{Name: aModel.Name}
		for i, aField := range aModel.Fields {
			field, err := getProtoFieldFromApidoc(aField, i+1)
			if err != nil {
				panic(err)
			}

			pMessage.Elements = append(
				pMessage.Elements,
				field,
			)
		}
		pFile.Elements = append(pFile.Elements, &pMessage)
	}

	//Unions

	//Resources

	return pFile
}

func getProtoFieldFromApidoc(aField apidoc.Field, sequenceNumber int) (proto.Visitee, error) {
	if strings.HasPrefix(aField.FieldType, "map[") {
		return getMapProtoFieldFromApidoc(aField, sequenceNumber)
	}

	return getNormalProtoFieldFromApidoc(aField, sequenceNumber)
}

func getNormalProtoFieldFromApidoc(aField apidoc.Field, sequenceNumber int) (*proto.NormalField, error) {
	pField := proto.NormalField{Field: &proto.Field{}}

	fieldType := aField.FieldType
	if strings.HasPrefix(fieldType, "[") {
		if !strings.HasSuffix(fieldType, "]") {
			return nil, errors.New("Invalid type, starts with a [ but does not end with one")
		}

		pField.Repeated = true
		fieldType = fieldType[1:len(fieldType)-1]
	}

	pType := getProtoTypeFromBasicApidocType(fieldType)
	if len(pType) == 0 {
		return nil, createUnsupportedError(aField.Name, pType)
	}
	pField.Type = pType

	pField.Sequence = sequenceNumber
	pField.Name = aField.Name

	return &pField, nil
}

func getMapProtoFieldFromApidoc(aField apidoc.Field, sequenceNumber int) (*proto.MapField, error) {
	return &proto.MapField{}, createUnsupportedError(aField.Name, "map")
}

func getProtoTypeFromBasicApidocType(basicType string) string {
	switch basicType {
	case "boolean":
		return "bool"
	case "date-iso8601":
		return "string"
	case "date-time-iso8601":
		return "string"
	case "decimal":
		return ""
	case "double":
		return "double"
	case "integer":
		return "int32"
	case "long":
		return "int64"
	case "object":
		return ""
	case "string":
		return "string"
	case "unit":
		return ""
	case "uuid":
		return "string"
	default:
		//Custom type or wrong type, can't tell here
		return basicType
	}
}


type UnsupportedTypeError error
func createUnsupportedError(fieldName string, fieldType string) error {
	return UnsupportedTypeError(errors.New(fmt.Sprintf("Cannot translate field [%s], field type: [%s] is unsupported", fieldName, fieldType)))
}