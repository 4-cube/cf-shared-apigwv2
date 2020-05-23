package cf_test

import (
	"github.com/4-cube/cf-shared-apigwv2/pkg/cf"
	"github.com/4-cube/cf-shared-apigwv2/x"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_IntegrationBuilder(t *testing.T) {
	type args struct {
		fnName string
		event  *cf.HttpApiEvent
		route  *cf.Route
	}
	tests := []struct {
		name         string
		args         args
		expectedName string
		expectedJSON []byte
	}{
		{
			name: "simple-integration.json",
			args: args{
				fnName: "HelloWorldFunction",
				event: &cf.HttpApiEvent{
					Name: "ProxyRoute",
					HttpApi: cf.HttpApi{
						ImportApiId:          "ImportApiId",
						Auth:                 cf.Auth{},
						Method:               "ANY",
						Path:                 "/{proxy+}",
						PayloadFormatVersion: "2.0",
						TimeoutInMillis:      5000,
					},
				},
			},
			expectedName: "HelloWorldFunctionProxyRouteIntegration",
			expectedJSON: x.LoadTestFile("simple-integration.json"),
		},
		{
			name: "payload-version-1-integration.json",
			args: args{
				fnName: "HelloWorldFunction",
				event: &cf.HttpApiEvent{
					Name: "ProxyRoute",
					HttpApi: cf.HttpApi{
						ImportApiId:          "ImportApiId",
						Auth:                 cf.Auth{},
						Method:               "ANY",
						Path:                 "/{proxy+}",
						PayloadFormatVersion: "1.0",
						TimeoutInMillis:      5000,
					},
				},
			},
			expectedName: "HelloWorldFunctionProxyRouteIntegration",
			expectedJSON: x.LoadTestFile("payload-version-1-integration.json"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := cf.NewIntegrationBuilder(tt.args.fnName, tt.args.event)
			j, err := rb.JSON()
			assert.NoError(t, err)
			require.NotNil(t, j)
			assert.JSONEq(t, string(tt.expectedJSON), string(j))
			assert.Equal(t, tt.expectedName, rb.Name())
		})
	}
}
