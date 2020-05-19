package cf

import (
	"encoding/json"
	"fmt"
)

// Configures authorization at the event level.
// Configure Auth for a specific API + Path + Method
type Auth struct {
	// The Authorizer for a specific Function
	// If you have specified a Global Authorizer on the API and want to make a specific Function public, override by setting Authorizer to NONE.
	Authorizer string `json:"Authorizer, omitempty"`

	// Authorization scopes to apply to this API + Path + Method.
	// Scopes listed here will override any scopes applied by the DefaultAuthorizer if one exists.
	AuthorizationScopes []string `json:"AuthorizationScopes, omitempty"`
}

// The object describing an event source with type HttpApi.
// see: https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/sam-property-function-httpapi.html
// Only difference between SAM HttpApi is that we are expecting ImportApiId to be Fn::ImportValue intrinsic function
type HttpApi struct {
	// Shared API Gateway ID - expected to be: ExportedApiId
	ImportApiId string `json:"ImportApiId"`

	// Auth configuration for this specific Api+Path+Method.
	// Useful for overriding the API's DefaultAuthorizer or setting auth config on an individual path when no DefaultAuthorizer is specified.
	Auth Auth `json:"Auth,omitempty"`

	// HTTP method for which this function is invoked.
	Method string `json:"Method,omitempty"`

	// Uri path for which this function is invoked. Must start with /.
	Path string `json:"Path"`

	// Specifies the format of the payload sent to an integration. Default: 2.0
	PayloadFormatVersion string `json:"PayloadFormatVersion,omitempty"`

	// Custom timeout between 50 and 29,000 milliseconds.
	// Default: 5000
	TimeoutInMillis int64 `json:"TimeoutInMillis,omitempty"`
}

type HttpApiEvent struct {
	Name    string
	HttpApi HttpApi
}

func (e *HttpApiEvent) RouteKey() string {
	return fmt.Sprintf("%s %s", e.HttpApi.Method, e.HttpApi.Path)
}

func (e *HttpApiEvent) FnImportApiId() json.RawMessage {
	return []byte(fmt.Sprintf(`{"Fn::ImportValue": "%s"}`, e.HttpApi.ImportApiId))
}
