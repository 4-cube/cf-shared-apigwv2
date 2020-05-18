package cf

import (
	"encoding/json"
	"fmt"
	"strings"
)

/*
"HelloWorldFunctionCatchAllPermission": {
      "Type": "AWS::Lambda::PermissionProperties",
      "Properties": {
        "Action": "lambda:InvokeFunction",
        "Principal": "apigateway.amazonaws.com",
        "FunctionName": {
          "Ref": "HelloWorldFunction"
        },
        "SourceArn": {
          "Fn::Sub": [
            "arn:${AWS::Partition}:execute-api:${AWS::Region}:${AWS::AccountId}:${__ApiId__}/${__Stage__}/GET/hello",
            {
              "__Stage__": "*",
              "__ApiId__": {
                "Ref": "ServerlessHttpApi"
              }
            }
          ]
        }
      }
    },
*/

type PermissionProperties struct {

	// Action AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-permission.html#cfn-lambda-permission-action
	Action string `json:"Action,omitempty"`

	// FunctionName AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-permission.html#cfn-lambda-permission-functionname
	FunctionName json.RawMessage `json:"FunctionName,omitempty"`

	// Principal AWS CloudFormation Property
	// Required: true
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-permission.html#cfn-lambda-permission-principal
	Principal string `json:"Principal,omitempty"`

	// SourceArn AWS CloudFormation Property
	// Required: false
	// See: http://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-resource-lambda-permission.html#cfn-lambda-permission-sourcearn
	SourceArn json.RawMessage `json:"SourceArn,omitempty"`
}

type Permission struct {
	Type       string               `json:"Type"`
	Properties PermissionProperties `json:"Properties"`
}

type permissionBuilder struct {
	fnName     string
	event      *SharedHttpApiEvent
	permission *Permission
}

func NewPermissionBuilder(fnName string, event *SharedHttpApiEvent) ResourceBuilder {
	return &permissionBuilder{
		fnName: fnName,
		event:  event,
		permission: &Permission{
			Type: "AWS::Lambda::Permission",
			Properties: PermissionProperties{
				Action:       "lambda:InvokeFunction",
				FunctionName: FunctionNameRef(fnName),
				Principal:    "apigateway.amazonaws.com",
				SourceArn:    SourceArn(event),
			},
		},
	}
}

// Marshal Permission ResourceKey to JSON
func (p *permissionBuilder) JSON() ([]byte, error) {
	return json.Marshal(p.permission)
}

// Get Integration ResourceKey Name
func (p *permissionBuilder) Name() string {
	return fmt.Sprintf("%s%sPermission", p.fnName, p.event.Name)
}

func FunctionNameRef(name string) json.RawMessage {
	return []byte(fmt.Sprintf(`{"Ref": "%s"}`, name))
}

func SourceArn(event *SharedHttpApiEvent) json.RawMessage {
	return []byte(fmt.Sprintf(`
	{
		"Fn::Sub": [
	       "%s",
	       {
				"__Stage__": "*",
				"__ApiId__": %s
	       }
		]
	}
	`, arn(event), event.SharedHttpApi.ApiId))
}

func arn(event *SharedHttpApiEvent) string {
	method := parseHttpMethod(event)
	path := parsePath(event)
	return fmt.Sprintf("arn:${AWS::Partition}:execute-api:${AWS::Region}:${AWS::AccountId}:${__ApiId__}/${__Stage__}/%s%s", method, path)
}

func parsePath(event *SharedHttpApiEvent) string {
	path := strings.ToLower(event.SharedHttpApi.Path)
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}

	if strings.ContainsAny(path, "{proxy}") {
		path = strings.ReplaceAll(path, "{proxy}", "*")
	}

	if strings.ContainsAny(path, "{proxy+}") {
		path = strings.ReplaceAll(path, "{proxy+}", "*")
	}
	return path
}

func parseHttpMethod(event *SharedHttpApiEvent) string {
	switch strings.ToUpper(event.SharedHttpApi.Method) {
	case "":
		fallthrough
	case "ANY":
		return "*"
	default:
		return strings.ToUpper(event.SharedHttpApi.Method)
	}
}
