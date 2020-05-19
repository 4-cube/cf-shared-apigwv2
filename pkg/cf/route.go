package cf

import (
	"encoding/json"
	"fmt"
)

/*
{
  "Type" : "AWS::ApiGatewayV2::Route",
  "Properties" : {
      "ImportApiId" : String,
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

// The AWS::ApiGatewayV2::Route resource creates a route for an API.
type Route struct {
	Type       string          `json:"Type"`      //AWS::ApiGatewayV2::Route
	DependsOn  []string        `json:"DependsOn"` //We are depending on Integration
	Properties RouteProperties `json:"Properties"`
}

type RouteProperties struct {
	ApiId json.RawMessage `json:"ImportApiId,omitempty"`

	//For HTTP APIs, valid values are NONE for open access, or JWT for using JSON Web Tokens.
	AuthorizationType string `json:"AuthorizationType,omitempty"`

	// The route key for the route.
	// eg: "POST /hello"
	RouteKey string `json:"RouteKey"`

	// The target for the route. !Ref to AWS::ApiGatewayV2::Integration
	Target json.RawMessage `json:"Target"`
}

type routeBuilder struct {
	fnName  string
	intName string
	event   *HttpApiEvent
	route   *Route
}

func NewRouteBuilder(fnName string, intName string, event *HttpApiEvent) ResourceBuilder {
	return &routeBuilder{
		fnName:  fnName,
		intName: intName,
		event:   event,
		route: &Route{
			Type:      "AWS::ApiGatewayV2::Route",
			DependsOn: []string{intName},
			Properties: RouteProperties{
				ApiId:    event.FnImportApiId(),
				RouteKey: event.RouteKey(),
				Target:   IntegrationTarget(intName),
			},
		},
	}
}

// Marshal Route ResourceKey to JSON
func (r *routeBuilder) JSON() ([]byte, error) {
	return json.Marshal(r.route)
}

// Get Route ResourceKey Name
func (r *routeBuilder) Name() string {
	return fmt.Sprintf("%s%sRoute", r.fnName, r.event.Name)
}

func IntegrationTarget(intName string) json.RawMessage {
	return []byte(fmt.Sprintf(`{"Fn::Join": ["/", ["integrations", {"Ref": "%s"}]]}`, intName))
}
