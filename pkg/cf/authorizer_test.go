package cf_test

import (
	"github.com/4-cube/cf-shared-apigwv2/pkg/cf"
	"github.com/4-cube/cf-shared-apigwv2/x"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_AuthorizerBuilder(t *testing.T) {
	type args struct {
		apiId  string
		config *cf.AuthConfig
	}
	tests := []struct {
		name         string
		args         args
		expectedName string
		expectedJSON []byte
	}{
		{
			name: "authorizer.json",
			args: args{
				apiId: "apiId",
				config: &cf.AuthConfig{
					AuthorizerName:      "Oauth2Authorizer",
					AuthorizationScopes: []string{"openid,email,profile,phone"},
					JwtConfiguration: cf.JwtConfiguration{
						Audience: []string{"audience"},
						Issuer:   "https://cognito-idp.{region}.amazonaws.com/{userPoolId}",
					},
					IdentitySource:    []string{"$request.header.Authorization"},
					DefaultAuthorizer: "Oauth2Authorizer",
				},
			},
			expectedName: "Oauth2Authorizer",
			expectedJSON: x.LoadTestFile("authorizer.json"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ab := cf.NewAuthorizerBuilder(tt.args.apiId, tt.args.config)
			j, err := ab.JSON()
			assert.NoError(t, err)
			require.NotNil(t, j)
			assert.JSONEq(t, string(tt.expectedJSON), string(j))
			assert.Equal(t, tt.expectedName, ab.Name())
		})
	}
}

func Test_AuthorizerExportBuilder(t *testing.T) {
	type args struct {
		apiId  string
		config *cf.AuthConfig
	}
	tests := []struct {
		name         string
		args         args
		expectedName string
		expectedJSON []byte
	}{
		{
			name: "authorizer-export.json",
			args: args{
				apiId: "apiId",
				config: &cf.AuthConfig{
					AuthorizerName:      "Oauth2Authorizer",
					AuthorizationScopes: []string{"openid,email,profile,phone"},
					JwtConfiguration: cf.JwtConfiguration{
						Audience: []string{"audience"},
						Issuer:   "https://cognito-idp.{region}.amazonaws.com/{userPoolId}",
					},
					IdentitySource:    []string{"$request.header.Authorization"},
					DefaultAuthorizer: "Oauth2Authorizer",
				},
			},
			expectedName: "Oauth2Authorizer",
			expectedJSON: x.LoadTestFile("authorizer-export.json"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ae := cf.NewAuthorizerExportBuilder(tt.args.apiId, tt.args.config)
			j, err := ae.JSON()
			assert.NoError(t, err)
			require.NotNil(t, j)
			assert.JSONEq(t, string(tt.expectedJSON), string(j))
			assert.Equal(t, tt.expectedName, ae.Name())
		})
	}
}

func Test_AuthorizerScopesExportBuilder(t *testing.T) {
	type args struct {
		apiId  string
		config *cf.AuthConfig
	}
	tests := []struct {
		name         string
		args         args
		expectedName string
		expectedJSON []byte
	}{
		{
			name: "authorizer-scopes-export.json",
			args: args{
				apiId: "apiId",
				config: &cf.AuthConfig{
					AuthorizerName:      "Oauth2Authorizer",
					AuthorizationScopes: []string{"openid,email,profile,phone"},
					JwtConfiguration: cf.JwtConfiguration{
						Audience: []string{"audience"},
						Issuer:   "https://cognito-idp.{region}.amazonaws.com/{userPoolId}",
					},
					IdentitySource:    []string{"$request.header.Authorization"},
					DefaultAuthorizer: "Oauth2Authorizer",
				},
			},
			expectedName: "Oauth2AuthorizerScopes",
			expectedJSON: x.LoadTestFile("authorizer-scopes-export.json"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ae := cf.NewAuthorizerScopesExportBuilder(tt.args.config)
			j, err := ae.JSON()
			assert.NoError(t, err)
			require.NotNil(t, j)
			assert.JSONEq(t, string(tt.expectedJSON), string(j))
			assert.Equal(t, tt.expectedName, ae.Name())
		})
	}
}
