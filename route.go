package main

import (
	"encoding/json"
	"fmt"
)

/*
{
  "Type" : "AWS::ApiGatewayV2::Route",
  "Properties" : {
      "ApiId" : String,
      "ApiKeyRequired" : Boolean,
      "AuthorizationScopes" : [ String, ... ],
      "AuthorizationType" : String,
      "AuthorizerId" : String,
      "ModelSelectionExpression" : String,
      "OperationName" : String,
      "RequestModels" : Json,
      "RequestParameters" : Json,
      "RouteKey" : String,
      "RouteResponseSelectionExpression" : String,
      "Target" : String
    }
}
*/

type RouteProperties struct {
	ApiId             json.RawMessage `json:"ApiId,omitempty"`
	AuthorizationType string          `json:"AuthorizationType,omitempty"` //For HTTP APIs, valid values are NONE for open access, or JWT for using JSON Web Tokens.
	RouteKey          string          `json:"RouteKey"`                    //POST /signup
	Target            json.RawMessage `json:"Target"`                      //Integration target
}

type Route struct {
	Type       string          `json:"Type"`      //AWS::ApiGatewayV2::Route
	DependsOn  []string        `json:"DependsOn"` //We are depending on Integration
	Properties RouteProperties `json:"Properties"`
}

func NewRoute(funName string, integrationName string, event FunctionEvent) (string, []byte, error) {
	route := Route{
		Type:      "AWS::ApiGatewayV2::Route",
		DependsOn: []string{integrationName},
		Properties: RouteProperties{
			ApiId:    event.SharedHttpApi.ApiId,
			RouteKey: getRouteKey(event),
			Target:   getIntegrationTarget(integrationName),
		},
	}
	j, err := json.Marshal(route)
	return getRouteName(funName, event), j, err
}

func getRouteName(funName string, event FunctionEvent) string {
	return fmt.Sprintf("%s%sRoute", funName, event.Name)
}

func getRouteKey(event FunctionEvent) string {
	return fmt.Sprintf("%s %s", event.SharedHttpApi.Method, event.SharedHttpApi.Path)
}

func getIntegrationTarget(integrationName string) json.RawMessage {
	return []byte(fmt.Sprintf(`
	{
		"Fn::Join": [
		"/",
		[
			"integrations",
		{
			"Ref": "%s"
		}
		]
		]
	}
	`, integrationName))
}
