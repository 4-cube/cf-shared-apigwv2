package cf

import (
	"encoding/json"
	"fmt"
)

/*
{
  "Type" : "AWS::ApiGatewayV2::Integration",
  "Properties" : {
      "ImportApiId" : String,
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

// The AWS::ApiGatewayV2::Integration resource creates a route for an API.
type IntegrationProperties struct {
	ApiId json.RawMessage `json:"ApiId"`

	// Specifies the integration's HTTP method type.
	// For some reason it should be POST
	IntegrationMethod string `json:"IntegrationMethod"` //POST ???

	// AWS_PROXY: for integrating the route or method request with the Lambda function-invoking action with the client
	// request passed through as-is. This integration is also referred to as Lambda proxy integration.
	IntegrationType string `json:"IntegrationType"`

	// For a Lambda integration, specify the URI of a Lambda function.
	IntegrationUri json.RawMessage `json:"IntegrationUri"`

	// Specifies the format of the payload sent to an integration. Required for HTTP APIs.
	// For HTTP APIs, supported values for Lambda proxy integrations are 1.0 and 2.0.
	// For all other integrations, 1.0 is the only supported value
	PayloadFormatVersion string `json:"PayloadFormatVersion"`

	// Custom timeout between 50 and 29,000 milliseconds for WebSocket APIs and between 50 and 30,000 milliseconds for HTTP APIs.
	// The default timeout is 29 seconds for WebSocket APIs and 30 seconds for HTTP APIs.
	TimeoutInMillis int64 `json:"TimeoutInMillis,omitempty"`
}

type Integration struct {
	Type       string                `json:"Type"` //AWS::ApiGatewayV2::Integration
	Properties IntegrationProperties `json:"Properties"`
}

type integrationBuilder struct {
	fnName      string
	event       *HttpApiEvent
	integration *Integration
}

func NewIntegrationBuilder(fnName string, event *HttpApiEvent) ResourceBuilder {
	return &integrationBuilder{
		fnName: fnName,
		event:  event,
		integration: &Integration{
			Type: "AWS::ApiGatewayV2::Integration",
			Properties: IntegrationProperties{
				ApiId:                event.FnImportApiId(),
				IntegrationMethod:    "POST",
				IntegrationType:      "AWS_PROXY",
				IntegrationUri:       IntegrationUri(fnName),
				PayloadFormatVersion: "2.0",
				TimeoutInMillis:      event.HttpApi.TimeoutInMillis,
			},
		},
	}
}

// Marshal Integration ResourceKey to JSON
func (i *integrationBuilder) JSON() ([]byte, error) {
	return json.Marshal(i.integration)
}

// Get Integration ResourceKey Name
func (i *integrationBuilder) Name() string {
	return fmt.Sprintf("%s%sIntegration", i.fnName, i.event.Name)
}

func IntegrationUri(funName string) json.RawMessage {
	return []byte(fmt.Sprintf(`
	{
		"Fn::Sub": "arn:${AWS::Partition}:apigateway:${AWS::Region}:lambda:path/2015-03-31/functions/${%s.Arn}/invocations"
	}
	`, funName))
}
