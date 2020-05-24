package cf

import (
	"encoding/json"
	"fmt"
	"strings"
)

type AuthConfig struct {
	AuthorizerName      string
	AuthorizationScopes []string
	JwtConfiguration    JwtConfiguration
	IdentitySource      []string
	DefaultAuthorizer   string
}

/*
  Oauth2Authorizer:
    Type: AWS::ApiGatewayV2::Authorizer
    Properties:
      ApiId: !Ref HttpApi
      AuthorizerType: "JWT"
      IdentitySource:
        - "$request.header.Authorization"
      JwtConfiguration:
        Audience:
          - audience
        Issuer: "https://cognito-idp.{region}.amazonaws.com/{userPoolId}"
      Name: Oauth2Authorizer

"Oauth2Authorizer": {
      "Type": "AWS::ApiGatewayV2::Authorizer",
      "Properties": {
        "ApiId": {
          "Ref": "HttpApi"
        },
        "IdentitySource": [
          "$request.header.Authorization"
        ],
        "Name": "Oauth2Authorizer",
        "JwtConfiguration": {
          "Audience": [
            "audience"
          ],
          "Issuer": "https://cognito-idp.{region}.amazonaws.com/{userPoolId}"
        },
        "AuthorizerType": "JWT"
      }
    },


"AuthorizationScopes": {
      "Export": {
        "Name": {
          "Fn::Sub": "${AWS::StackName}-DefaultAuthorizationScopes"
        }
      },
      "Description": "API Gateway Default OauthAuthorizer AuthorizationScopes",
      "Value": "openid,email,profile,phone"
    },
*/

type Authorizer struct {
	Type       string               `json:"Type"`
	Properties AuthorizerProperties `json:"Properties"`
}

type AuthorizerProperties struct {
	ApiId json.RawMessage `json:"ApiId"`

	// For HTTP APIs, specify JWT to use JSON Web Tokens.
	AuthorizerType string `json:"AuthorizerType"`

	// For JWT, a single entry that specifies where to extract the JSON Web Token (JWT) from inbound requests.
	// Currently only header-based and query parameter-based selections are supported, for example "$request.header.Authorization".
	IdentitySource []string `json:"IdentitySource"`

	// The name of the authorizer.
	Name string `json:"Name"`

	// The JWTConfiguration property specifies the configuration of a JWT authorizer.
	// Required for the JWT authorizer type. Supported only for HTTP APIs.
	JwtConfiguration JwtConfiguration `json:"JwtConfiguration"`
}

type JwtConfiguration struct {
	// A list of the intended recipients of the JWT. A valid JWT must provide an aud that matches at least one entry in this list.
	// See RFC 7519. Required for the JWT authorizer type. Supported only for HTTP APIs.
	Audience []string `json:"Audience"`

	// The base domain of the identity provider that issues JSON Web Tokens.
	// For example, an Amazon Cognito user pool has the following format: https://cognito-idp.{region}.amazonaws.com/{userPoolId} .
	// Required for the JWT authorizer type. Supported only for HTTP APIs.
	Issuer string
}

type AuthorizerScopes struct {
	Export      json.RawMessage `json:"Export"`
	Description string          `json:"Description"`
	Value       string          `json:"Value"`
}

type AuthorizerExport struct {
	ExportName string
	Export     json.RawMessage
}

type authorizerBuilder struct {
	apiId      string
	authConfig *AuthConfig
	Authorizer *Authorizer
}

func NewAuthorizerBuilder(apiId string, config *AuthConfig) ResourceBuilder {
	return &authorizerBuilder{
		apiId:      apiId,
		authConfig: config,
		Authorizer: &Authorizer{
			Type: "AWS::ApiGatewayV2::Authorizer",
			Properties: AuthorizerProperties{
				ApiId:            apiIdRef(apiId),
				AuthorizerType:   "JWT",
				IdentitySource:   config.IdentitySource,
				Name:             config.AuthorizerName,
				JwtConfiguration: config.JwtConfiguration,
			},
		},
	}
}

func (a *authorizerBuilder) JSON() ([]byte, error) {
	return json.Marshal(a.Authorizer)
}

func (a *authorizerBuilder) Name() string {
	return a.authConfig.AuthorizerName
}

func apiIdRef(apiId string) json.RawMessage {
	return []byte(fmt.Sprintf(`{"Ref": "%s"}`, apiId))
}

type authorizerExportBuilder struct {
	apiId  string
	config *AuthConfig
	export *AuthorizerExport
}

func NewAuthorizerExportBuilder(apiId string, config *AuthConfig) ResourceBuilder {
	return &authorizerExportBuilder{
		apiId:  apiId,
		config: config,
		export: &AuthorizerExport{
			ExportName: config.AuthorizerName,
			Export:     buildExport(config.AuthorizerName),
		},
	}
}

func (a *authorizerExportBuilder) JSON() ([]byte, error) {
	return a.export.Export, nil
}

func (a *authorizerExportBuilder) Name() string {
	return a.export.ExportName
}

func buildExport(name string) json.RawMessage {
	return []byte(fmt.Sprintf(`
	{
		"Description": "API Gateway OauthAuthorizer",
		"Value": {
			"Ref": "%s"
		},
		"Export": {
			"Name": {
				"Fn::Sub": "${AWS::StackName}-%s"
			}
		}
	}
	`, name, name))
}

type authExportScopesExportBuilder struct {
	config *AuthConfig
	export *AuthorizerExport
}

func NewAuthorizerScopesExportBuilder(config *AuthConfig) ResourceBuilder {
	return &authExportScopesExportBuilder{
		config: config,
		export: &AuthorizerExport{
			ExportName: config.AuthorizerName + "Scopes",
			Export:     buildScopesExport(config.AuthorizerName, config.AuthorizationScopes),
		},
	}
}

func buildScopesExport(name string, scopes []string) json.RawMessage {
	return []byte(fmt.Sprintf(`
	{
		"Description": "API Gateway OauthAuthorizer Scopes",
		"Value": "%s",
		"Export": {
			"Name": {
				"Fn::Sub": "${AWS::StackName}-%sScopes"
			}
		}
	}
	`, strings.Join(scopes, ","), name))
}

func (a *authExportScopesExportBuilder) JSON() ([]byte, error) {
	return a.export.Export, nil
}

func (a *authExportScopesExportBuilder) Name() string {
	return a.export.ExportName
}
