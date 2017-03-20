package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
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
				&proto.EnumField{Name: aValue.Name, Integer: i},
			)
		}
		pFile.Elements = append(pFile.Elements, &pEnum)
	}

	//Models
	for _, aModel := range spec.Models {
		pMessage := proto.Message{Name: aModel.Name}
		for i, aField := range aModel.Fields {
			field, err := getProtoFieldFromApidoc(aField, i+1) //TODO: how to handle sequence number changes
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
	for _, aUnion := range spec.Unions {
		pMessage := proto.Message{Name: aUnion.Name}
		oneOff := proto.Oneof{Name: aUnion.Discriminator}
		if len(oneOff.Name) == 0 {
			oneOff.Name = "union"
		}

		for i, aUnionType := range aUnion.UnionTypes {
			pType := getProtoTypeFromBasicApidocType(aUnionType.TypeValue) //TODO: how to handle sequence number changes
			if len(pType) == 0 {
				panic(fmt.Sprintf("Couldn't find a proto type from union type [%v]", aUnionType))
			}

			//TODO: OPINION: I decided to use the type name here, should convert to underscore case
			field := proto.Field{Name: pType, Type: pType, Sequence: i + 1}

			pMessage.Elements = append(
				pMessage.Elements,
				&proto.OneOfField{Field: &field},
			)
		}
		pFile.Elements = append(pFile.Elements, &pMessage)
	}

	//Resources
	service := proto.Service{Name: getSafeServiceName(spec.Name)}
	pFile.Elements = append(pFile.Elements, &service)

	for _, resource := range spec.Resources {
		for _, operation := range resource.Operations {
			rpc := proto.RPC{}

			rpc.Name = getProtoMethodName(resource, operation)
			//TODO: Handle gathering path, query and body params into one argument (or don't bother??)
			service.Elements = append(service.Elements, &rpc)
		}
	}

	return pFile
}

var ProtoIdentifierRegex = regexp.MustCompile(`[A-Za-z_][\w_]*`)

//Gets a safe service name, proto expects /[A-Za-z_][\w_]*/ for service names
func getSafeServiceName(name string) string {
	pName := convertPartsToPascalCase(strings.Split(name, " "))
	if len(ProtoIdentifierRegex.FindString(pName)) != len(pName) {
		//TODO: use special aproto attribute
		panic(fmt.Sprintf(
			"Invalid proto name for service, name was [%v] but protobuf identifier must follow the regex [%v]",
			name,
			ProtoIdentifierRegex.String(),
		))
	}

	return pName
}

func convertPartsToPascalCase(parts []string) string {
	name := ""
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}

		name += strings.ToUpper(part[:1]) + part[1:]
	}

	return name
}

func getProtoMethodName(resource apidoc.Resource, operation apidoc.Operation) string {
	if len(operation.Path) == 0 {
		return convertPathToProtoName(resource.Path)
	}

	parts := append(strings.Split(resource.Path, "/"), strings.Split(operation.Path, "/")...)

	return convertPathSegmentsToProtoName(parts)
}

func convertPathToProtoName(path string) string {
	parts := strings.Split(strings.TrimSuffix(path, "/"), "/")

	return convertPathSegmentsToProtoName(parts)
}

func convertPathSegmentsToProtoName(parts []string) string {
	name := ""
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		if part[0] == ':' {
			continue
		}

		name += strings.ToUpper(part[:1]) + part[1:]
	}

	return name
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
		fieldType = fieldType[1 : len(fieldType)-1]
	}

	pType := getProtoTypeFromBasicApidocType(fieldType)
	if len(pType) == 0 {
		return nil, createUnsupportedError(aField.Name, pType)
	}
	pField.Type = pType

	pField.Sequence = sequenceNumber
	//TODO: Translate names as per proto styles? https://developers.google.com/protocol-buffers/docs/style#message-and-field-names
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
