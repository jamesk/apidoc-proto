package apidoc

/*
{
"name": string,
"apidoc": JSON Object of Apidoc (optional),
"info": JSON Object of Info (optional),
"namespace": string (optional),
"base_url": string (optional),
"description": string (optional),
"imports": JSON Array of Import (optional),
"headers": JSON Array of Header (optional),
"enums": JSON Object of Enum (optional),
"models": JSON Object of Model (optional),
"unions": JSON Object of Union (optional),
"resources": JSON Object of Resource (optional),
"attributes": JSON Array of Attribute (optional)
}
*/

type Spec struct {
	Name       string      `json:"name"`
	Enums      []Enum      `json:"enums"`
	Unions     []Union     `json:"unions"`
	Models     []Model     `json:"models"`
	Resources  []Resource  `json:"resources"`
	Attributes []Attribute `json:"attributes"`
}

/*
{
  "name": string,
  "plural": string (optional),
  "description": string (optional),
  "values": JSON Array of EnumValue,
  "attributes": JSON Array of Attribute (optional),
  "deprecation": JSON Object of Deprecation (optional)
}
*/
type Enum struct {
	Name   string      `json:"name"`
	Values []EnumValue `json:"values"`
}

/*
{
  "name": string,
  "description": string (optional),
  "attributes": JSON Array of Attribute (optional),
  "deprecation": JSON Object of Deprecation (optional)
}
*/
type EnumValue struct {
	Name string `json:"name"`
}

/*
{
  "name": string,
  "description": string (optional),
  "plural": string (optional),
  "fields": JSON Array of Field,
  "attributes": JSON Array of Attribute (optional),
  "deprecation": JSON Object of Deprecation (optional)
}
*/
type Model struct {
	Name   string  `json:"name"`
	Fields []Field `json:"fields"`
}

/*
{
  "name": string,
  "type": string,
  "description": string (optional),
  "required": boolean (optional, true by default),
  "default": value (optional),
  "example": string (optional),
  "minimum": long (optional),
  "maximum": long (optional),
  "attributes": JSON Array of Attribute (optional),
  "deprecation": JSON Object of Deprecation (optional)
}
*/
type Field struct {
	Name         string      `json:"name"`
	FieldType    string      `json:"type"`
	Required     bool        `json:"required"`
	DefaultValue interface{} `json:"default"`
}

/*
{
  "name": string,
  "plural": string (optional),
  "discriminator": string (optional),
  "description": string (optional),
  "types": JSON Array of UnionType,
  "attributes": JSON Array of Attribute (optional),
  "deprecation": JSON Object of Deprecation (optional)
}
*/
type Union struct {
	Name          string           `json:"name"`
	Discriminator string           `json:"discriminator"`
	UnionTypes    []UnionValueType `json:"types"`
}

/*
{
      "type": string,
      "description": string (optional),
      "attributes": JSON Array of Attribute (optional),
      "deprecation": JSON Object of Deprecation (optional)
    }
*/
type UnionValueType struct {
	TypeValue string `json:"type"`
}

/*
{
  "type": string,
  "path": string (optional),
  "description": string (optional),
  "operations": JSON Array of Operation,
  "attributes": JSON Array of Attribute (optional),
  "deprecation": JSON Object of Deprecation (optional)
}
*/
type Resource struct {
	ResourceType string      `json:"type"`
	Path         string      `json:"path,omitempty"`
	Operations   []Operation `json:"operations"`
}

/*
{
  "method": string,
  "path": string (optional),
  "description": string (optional),
  "body": JSON Object of Body (optional),
  "parameters": JSON Array of Parameter (optional),
  "responses": JSON Object of Response (optional),
  "attributes": JSON Array of Attribute (optional),
  "deprecation": JSON Object of Deprecation (optional)
}
*/
type Operation struct {
	Method     string      `json:"method"`
	Path       string      `json:"path,omitempty"`
	Body       Body        `json:"body"`
	Parameters []Parameter `json:"parameters"`
	Responses  []Response  `json:"responses"`
	Attributes []Attribute `json:"attributes"`
}

/*
{
  "type": string,
  "description": string (optional),
  "attributes": JSON Array of Attribute (optional),
  "deprecation": JSON Object of Deprecation (optional),
}
*/
type Body struct {
	BodyType   string      `json:"type"`
	Attributes []Attribute `json:"attributes"`
}

/*
{
  "name": string,
  "type": string,
  "location": string (optional),
  "description": string (optional),
  "deprecation": JSON Object of Deprecation (optional),
  "required": boolean (optional, true by default),
  "default": value (optional),
  "example": string (optional),
  "minimum": long (optional),
  "maximum": long (optional)
}
*/
type Parameter struct {
	Name          string      `json:"name"`
	ParameterType string      `json:"type"`
	Location      string      `json:"location"`
	Required      bool        `json:"required"`
	Default       interface{} `json:"default"`
}

/*
SEEMINGLY OUTDATED SPEC:
{
"http_status_code": {
"type": type,
"headers": JSON Array of Header (optional),
"description": string (optional),
"deprecation": JSON Object of Deprecation (optional)
}
}

ACTUAL EXAMPLE:
"responses": [
            {
              "code": {
                "integer": {
                  "value": 200
                }
              },
              "type": "ticket"
            },
*/
type Response struct {
	Code         ResponseCode `json:"code"`
	ResponseType string       `json:"type"`
	Attributes   []Attribute  `json:"attributes"`
}
type ResponseCode struct {
	Integer ResponseCodeValue `json:"integer"`
}
type ResponseCodeValue struct {
	Value int `json:"value"`
}

/*
{
  "name": string,
  "value": JSON Object
}
*/
type Attribute struct {
	Name  string                 `json:"name"`
	Value map[string]interface{} `json:"value"`
}
