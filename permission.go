package main

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

func NewPermission(funName string, event FunctionEvent) (string, []byte, error) {
	permission := Permission{
		Type: "AWS::Lambda::Permission",
		Properties: PermissionProperties{
			Action:       "lambda:InvokeFunction",
			FunctionName: jsonRef(funName),
			Principal:    "apigateway.amazonaws.com",
			SourceArn:    sourceArn(event),
		},
	}
	json, err := json.Marshal(permission)
	return getPermissionName(funName, event), json, err
}

func getPermissionName(funName string, event FunctionEvent) string {
	return fmt.Sprintf("%s%sPermission", funName, event.Name)
}

func jsonRef(name string) json.RawMessage {
	return []byte(fmt.Sprintf(`{"Ref": "%s"}`, name))
}

func sourceArn(event FunctionEvent) json.RawMessage {
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

func arn(event FunctionEvent) string {
	method := parseHttpMethod(event)
	path := parsePath(event)
	return fmt.Sprintf("arn:${AWS::Partition}:execute-api:${AWS::Region}:${AWS::AccountId}:${__ApiId__}/${__Stage__}/%s%s", method, path)
}

func parsePath(event FunctionEvent) string {
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

func parseHttpMethod(event FunctionEvent) string {
	switch strings.ToUpper(event.SharedHttpApi.Method) {
	case "":
		fallthrough
	case "ANY":
		return "*"
	default:
		return strings.ToUpper(event.SharedHttpApi.Method)
	}
}
