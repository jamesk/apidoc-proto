package apidoc_test

import (
	"testing"

	"github.com/jamesk/apidoc-proto/apidoc"
	"github.com/stretchr/testify/assert"
)

var testServiceExpected = apidoc.Spec{Name: "test-apidoc-proto-ticket", Enums: []apidoc.Enum{{Name: "ticketType", Values: []apidoc.EnumValue{{Name: "feature"}, {Name: "bug"}}}, {Name: "size", Values: []apidoc.EnumValue{{Name: "small"}, {Name: "medium"}, {Name: "large"}}}}, Unions: []apidoc.Union{{Name: "ticket", Discriminator: "", UnionTypes: []apidoc.UnionValueType{{TypeValue: "feature"}, {TypeValue: "bug"}}}}, Models: []apidoc.Model{{Name: "bug", Fields: []apidoc.Field{{Name: "id", FieldType: "uuid", Required: true, DefaultValue: interface{}(nil)}, {Name: "title", FieldType: "string", Required: true, DefaultValue: interface{}(nil)}, {Name: "problem", FieldType: "string", Required: false, DefaultValue: interface{}(nil)}, {Name: "affectedVersion", FieldType: "integer", Required: true, DefaultValue: interface{}(nil)}}}, {Name: "feature", Fields: []apidoc.Field{{Name: "id", FieldType: "uuid", Required: true, DefaultValue: interface{}(nil)}, {Name: "title", FieldType: "string", Required: true, DefaultValue: interface{}(nil)}, {Name: "story", FieldType: "string", Required: false, DefaultValue: interface{}(nil)}, {Name: "acceptanceCriteria", FieldType: "string", Required: false, DefaultValue: interface{}(nil)}, {Name: "size", FieldType: "size", Required: true, DefaultValue: "medium"}}}}, Resources: []apidoc.Resource{{ResourceType: "ticket", Path: "", Operations: []apidoc.Operation{{Method: "GET", Path: "", Body: apidoc.Body{BodyType: ""}, Parameters: []apidoc.Parameter{}, Responses: []apidoc.Response{{Code: apidoc.ResponseCode{Integer: apidoc.ResponseCodeValue{Value: 200}}, ResponseType: "[ticket]"}}}, {Method: "GET", Path: "/bugs", Body: apidoc.Body{BodyType: ""}, Parameters: []apidoc.Parameter{}, Responses: []apidoc.Response{{Code: apidoc.ResponseCode{Integer: apidoc.ResponseCodeValue{Value: 200}}, ResponseType: "[bug]"}}}, {Method: "GET", Path: "/features", Body: apidoc.Body{BodyType: ""}, Parameters: []apidoc.Parameter{}, Responses: []apidoc.Response{{Code: apidoc.ResponseCode{Integer: apidoc.ResponseCodeValue{Value: 200}}, ResponseType: "[feature]"}}}, {Method: "GET", Path: "/:id", Body: apidoc.Body{BodyType: ""}, Parameters: []apidoc.Parameter{{Name: "id", ParameterType: "uuid", Location: "Path", Required: true, Default: interface{}(nil)}}, Responses: []apidoc.Response{{Code: apidoc.ResponseCode{Integer: apidoc.ResponseCodeValue{Value: 200}}, ResponseType: "ticket"}, {Code: apidoc.ResponseCode{Integer: apidoc.ResponseCodeValue{Value: 404}}, ResponseType: "unit"}}}, {Method: "POST", Path: "", Body: apidoc.Body{BodyType: "ticket"}, Parameters: []apidoc.Parameter{}, Responses: []apidoc.Response{{Code: apidoc.ResponseCode{Integer: apidoc.ResponseCodeValue{Value: 201}}, ResponseType: "ticket"}}}}}}}

func TestUnmarshal(t *testing.T) {
	testSpec, err := apidoc.GetSpecFromFile("test_service.json")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, testServiceExpected, testSpec)
}