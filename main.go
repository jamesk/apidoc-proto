package main

import (
	"bufio"
	"os"

	"github.com/emicklei/proto"
	"github.com/jamesk/apidoc-proto/apidoc"
	"strings"
	"errors"
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

			//Apidoc types: boolean, date-iso8601, date-time-iso8601, decimal, double, integer, long, object, string, unit, uuid

			pMessage.Elements = append(
				pMessage.Elements,
				&proto. {Name:aValue.Name, Integer:i},
			)
		}
		pFile.Elements = append(pFile.Elements, &pMessage)
	}

	//Unions

	//Resources

	return pFile
}

func getProtoFieldFromApidoc(aField apidoc.Field, sequenceNumber int) (proto.Visitee, error) {
	//Apidoc types: boolean, date-iso8601, date-time-iso8601, decimal, double, integer, long, object, string, unit, uuid
	if strings.HasPrefix(aField.FieldType, "map[") {
		return getMapProtoFieldFromApidoc(aField, sequenceNumber)
	}

	getNormalProtoFieldFromApidoc(aField, sequenceNumber)
}

func getNormalProtoFieldFromApidoc(aField apidoc.Field, sequenceNumber int) (proto.NormalField, error) {
	pField := proto.NormalField{}

	fieldType := aField.FieldType
	if strings.HasPrefix(fieldType, "[") {
		if !strings.HasSuffix(fieldType, "]") {
			return pField, errors.New("Invalid type, starts with a [ but does not end with one")
		}

		pField.Repeated = true
		fieldType = fieldType[1:len(fieldType)-1]
	}

	pType, err := getProtoTypeFromBasicApidocType(fieldType)
	if err != nil {
		return pField, err
	}
	pField.Type = pType

	pField.Sequence = sequenceNumber

	return pField, nil
}

func getMapProtoFieldFromApidoc(aField apidoc.Field) (proto.MapField, error) {
	return proto.MapField{}, errors.New("Map fields are unsupported at the moment")
}

func getProtoTypeFromBasicApidocType(basicType string) (string, error) {
	switch basicType {
	case "boolean":
		return "bool", nil
	case "date-iso8601":
		return "string", nil
	case "date-time-iso8601":
		return "string", nil
	case "decimal":
		return "", errors.New("Cannot translate decimal field types to proto type")
	case "double":
		return "double", nil
	case "integer":
		return "int32", nil
	case "long":
		return "int64", nil
	case "object":
		return "", errors.New("Cannot translate object field types to proto type")
	case "string":
		return "string", nil
	case "unit":
		return "", errors.New("Cannot translate unit field types to proto type")
	case "uuid":
		return "string", nil
	default:
		//Custom type or wrong type
		return basicType, nil
	}
}