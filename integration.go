package main

import (
	"encoding/json"
	"fmt"
)

/*
{
  "Type" : "AWS::ApiGatewayV2::Integration",
  "Properties" : {
      "ApiId" : String,
      "ConnectionId" : String,
      "ConnectionType" : String,
      "ContentHandlingStrategy" : String,
      "CredentialsArn" : String,
      "Description" : String,
      "IntegrationMethod" : String,
      "IntegrationType" : String,
      "IntegrationUri" : String,
      "PassthroughBehavior" : String,
      "PayloadFormatVersion" : String,
      "RequestParameters" : Json,
      "RequestTemplates" : Json,
      "TemplateSelectionExpression" : String,
      "TimeoutInMillis" : Integer,
      "TlsConfig" : TlsConfig
    }
}
*/

type IntegrationProperties struct {
	ApiId                json.RawMessage `json:"ApiId,omitempty"`
	IntegrationMethod    string          `json:"IntegrationMethod"`    //POST ???
	IntegrationType      string          `json:"IntegrationType"`      //AWS_PROXY
	IntegrationUri       json.RawMessage `json:"IntegrationUri"`       //For a Lambda integration, specify the URI of a Lambda function.
	PayloadFormatVersion string          `json:"PayloadFormatVersion"` //Lambda proxy integrations are 1.0 and 2.0. For all other integrations, 1.0 is the only supported value.
	TimeoutInMillis      int64           `json:"TimeoutInMillis,omitempty"`
}

type Integration struct {
	Type       string                `json:"Type"` //AWS::ApiGatewayV2::Integration
	Properties IntegrationProperties `json:"Properties"`
}

func NewIntegration(funName string, event FunctionEvent) (string, []byte, error) {
	integration := Integration{
		Type: "AWS::ApiGatewayV2::Integration",
		Properties: IntegrationProperties{
			ApiId:                event.SharedHttpApi.ApiId,
			IntegrationMethod:    "POST",
			IntegrationType:      "AWS_PROXY",
			IntegrationUri:       getIntegrationUri(funName),
			PayloadFormatVersion: "2.0",
			TimeoutInMillis:      event.SharedHttpApi.TimeoutInMillis,
		},
	}
	j, err := json.Marshal(integration)
	return getIntegrationName(funName, event), j, err
}

func getIntegrationName(funName string, event FunctionEvent) string {
	return fmt.Sprintf("%s%sIntegration", funName, event.Name)
}

func getIntegrationUri(funName string) json.RawMessage {
	return []byte(fmt.Sprintf(`
	{
		"Fn::Sub": "arn:${AWS::Partition}:apigateway:${AWS::Region}:lambda:path/2015-03-31/functions/${%s.Arn}/invocations"
	}
	`, funName))
}
