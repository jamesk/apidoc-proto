package aproto

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/emicklei/proto"
	"github.com/jamesk/apidoc-proto/apidoc"
)

const (
	syntaxVersion = "proto3"
)

func GetProtoFromAPISpec(spec apidoc.Spec) proto.Proto {
	pFile := proto.Proto{}

	pFile.Elements = append(
		pFile.Elements,
		&proto.Syntax{Value: syntaxVersion},
	)

	//TODO:hack for rpc, not always needed
	pFile.Elements = append(
		pFile.Elements,
		&proto.Import{Filename: "google/protobuf/empty.proto"},
	)
	//TODO: hack for Object types, not always needed
	pFile.Elements = append(
		pFile.Elements,
		&proto.Import{Filename: "google/protobuf/any.proto"},
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
			pType := getProtoTypeFromBasicApidocType(aUnionType.TypeValue)
			if len(pType) == 0 {
				panic(fmt.Sprintf("Couldn't find a proto type from union type [%v]", aUnionType))
			}

			//TODO: OPINION: I decided to use the type as the name here, should convert to underscore case? (as per https://developers.google.com/protocol-buffers/docs/style#message-and-field-names)
			//TODO: how to handle sequence number changes
			field := proto.Field{Name: pType, Type: pType, Sequence: i + 1}

			oneOff.Elements = append(
				oneOff.Elements,
				&proto.OneOfField{Field: &field},
			)
		}
		pMessage.Elements = append(pMessage.Elements, &oneOff)
		pFile.Elements = append(pFile.Elements, &pMessage)
	}

	//Resources
	var nameReplace ServiceNameReplaceAttribute
	for _, attr := range spec.Attributes {
		if attr.Name == ServiceNameReplaceAttributeName {
			//TODO: inefficient, any actual performance impact?
			d, err := json.Marshal(attr.Value)
			if err != nil {
				panic(err)
			}

			json.Unmarshal(d, &nameReplace)
		}
	}

	if len(spec.Resources) == 0 {
		return pFile
	}

	service := proto.Service{Name: getSafeServiceName(spec.Name, nameReplace)}
	pFile.Elements = append(pFile.Elements, &service)

	for _, resource := range spec.Resources {
		//TODO: default to only taking first operation, make attribute to show that all operations should be considered.
		for _, operation := range resource.Operations {
			rpc := proto.RPC{}

			rpc.Name = getProtoMethodName(resource, operation)
			if rpc.Name == "" {
				continue //TODO: fail fast? Leaving in to handle /:id/ type paths
			}

			//TODO: Handle path and query params? Attributes to handle behaviour?
			request, err := getRpcParameter(operation.Body.BodyType, operation.Body.Attributes)
			if err != nil {
				panic(err)
			}
			rpc.StreamsRequest = request.IsStream
			rpc.RequestType = request.RpcType
			if len(request.Message.Name) > 0 {
				pFile.Elements = append(pFile.Elements, &request.Message)
			}

			if len(operation.Responses) == 0 {
				rpc.ReturnsType = "google.protobuf.Empty"
			} else {
				//TODO: handle multiple responses (only when an attribute is set)
				responseType := operation.Responses[0].ResponseType
				response, err := getRpcParameter(responseType, []apidoc.Attribute{})
				if err != nil {
					panic(err)
				}
				rpc.StreamsReturns = response.IsStream
				rpc.ReturnsType = response.RpcType
				if len(response.Message.Name) > 0 {
					pFile.Elements = append(pFile.Elements, &response.Message)
				}
			}

			service.Elements = append(service.Elements, &rpc)
		}
	}

	return pFile
}

var ProtoIdentifierRegex = regexp.MustCompile(`[A-Za-z_][\w_]*`)

const (
	ArrayAsStreamAttributeName      = "aproto:array-as-stream"
	ServiceNameReplaceAttributeName = "aproto:service-name-replace"
)

/*
{
	"name": "aproto:array-as-stream",
	"value": {
	  "value" : true
	}
}
*/
type ArrayAsStreamAttribute struct {
	Value bool `json:"value"`
}

/*
{
      "name": "aproto:service-name-replace",
      "value": {
        "regexMaps" : [{"regex": "[^\\w_]", "replace": "_"}]
      }
    }
*/
type ServiceNameReplaceAttribute struct {
	RegexMaps []ServiceNameReplaceAttributeRegexMap `json:"regexMaps"`
}

type ServiceNameReplaceAttributeRegexMap struct {
	Regex   string `json:"regex"`   //The regex to search for
	Replace string `json:"replace"` //The replace string, see regexp.Regex.Expand for details
}

//Gets a safe service name, proto expects /[A-Za-z_][\w_]*/ for service names
func getSafeServiceName(name string, replaceAttribute ServiceNameReplaceAttribute) string {
	replacedName := name
	for _, m := range replaceAttribute.RegexMaps {
		r, err := regexp.Compile(m.Regex)
		if err != nil {
			panic(err)
		}

		//TODO: logging fmt.Printf("Replacing, name currently: [%v], regex [%v]\n", replacedName, m.Regex)
		replacedName = r.ReplaceAllString(replacedName, m.Replace)
		//TODO: logging fmt.Printf("Done replacing, name currently: [%v]\n", replacedName)
	}

	pName := convertPartsToPascalCase(strings.Split(replacedName, " "))
	if len(ProtoIdentifierRegex.FindString(pName)) != len(pName) {
		panic(fmt.Sprintf(
			"Invalid proto name for service, processed name was [%v], original name was [%v] ([%v] after replacing) but protobuf identifier must follow the regex [%v]",
			pName,
			name,
			replacedName,
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
	path := resource.Path
	if len(path) == 0 {
		path = resource.ResourceType
	}

	parts := []string{}
	//TODO: attribute to ignore method?
	method := strings.ToLower(operation.Method)
	if len(method) != 0 {
		method = strings.ToUpper(method[:1]) + method[1:]
		parts = append(parts, method)
	}
	parts = append(parts, strings.Split(path, "/")...)
	parts = append(parts, strings.Split(operation.Path, "/")...)

	return convertPathSegmentsToProtoName(parts)
}

func convertPathSegmentsToProtoName(parts []string) string {
	name := ""
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		if part[0] == ':' {
			return ""
			//TODO: panic("Cannot handle path variables")
		}

		name += strings.ToUpper(part[:1]) + part[1:]
	}

	return name
}

func getProtoFieldFromApidoc(aField apidoc.Field, sequenceNumber int) (proto.Visitee, error) {
	mapType, isMap, err := getMapType(aField.FieldType)
	if err != nil {
		return nil, err
	}

	if isMap {
		return getMapProtoField(aField.Name, mapType, sequenceNumber)
	}
	return getNormalProtoFieldFromApidoc(aField, sequenceNumber)
}

func getMapType(value string) (string, bool, error) {
	mapType := value
	if strings.HasPrefix(mapType, "map[") {
		if !strings.HasSuffix(mapType, "]") {
			return value, false, fmt.Errorf("Invalid type, starts with a [ but does not end with one, type was [%v]", value)
		}

		return mapType[4 : len(mapType)-1], true, nil
	}

	return value, false, nil
}

func getArrayType(value string) (string, bool, error) {
	arrayType := value
	if strings.HasPrefix(arrayType, "[") {
		if !strings.HasSuffix(arrayType, "]") {
			return value, false, fmt.Errorf("Invalid type, starts with a [ but does not end with one, type was [%v]", value)
		}

		return arrayType[1 : len(arrayType)-1], true, nil
	}

	return value, false, nil
}

func getNormalProtoFieldFromApidoc(aField apidoc.Field, sequenceNumber int) (proto.Visitee, error) {
	fieldType, isArray, err := getArrayType(aField.FieldType)
	if err != nil {
		return nil, err
	}

	pType := getProtoTypeFromBasicApidocType(fieldType)
	if len(pType) == 0 {
		if fieldType == "object" {
			if isArray {
				return nil, createUnsupportedError(aField.Name, aField.FieldType)
			}

			return getMapProtoField(aField.Name, "google.protobuf.Any", sequenceNumber)
		}

		return nil, createUnsupportedError(aField.Name, pType)
	}

	pField := proto.NormalField{Field: &proto.Field{}}
	pField.Repeated = isArray
	pField.Type = pType
	pField.Sequence = sequenceNumber
	//TODO: Translate names as per proto styles? https://developers.google.com/protocol-buffers/docs/style#message-and-field-names
	pField.Name = aField.Name

	return &pField, nil
}

func getMapProtoField(name string, mapType string, sequenceNumber int) (proto.Visitee, error) {
	pField := proto.MapField{Field: &proto.Field{}}

	pField.Sequence = sequenceNumber
	pField.Name = name
	pType := getProtoTypeFromBasicApidocType(mapType)
	if len(pType) == 0 {
		return nil, createUnsupportedError(name, pType)
	}
	pField.Type = pType
	pField.KeyType = "string"

	return &pField, nil
}

//Input is a "basic" type i.e. not a map or array
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
		return "google.protobuf.Empty"
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
