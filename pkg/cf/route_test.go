package cf_test

import (
	"github.com/4-cube/cf-shared-apigwv2/pkg/cf"
	"github.com/4-cube/cf-shared-apigwv2/x"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_RouteBuilder(t *testing.T) {
	type args struct {
		fnName  string
		intName string
		event   *cf.HttpApiEvent
		route   *cf.Route
	}
	tests := []struct {
		name string
		args args
		expectedName string
		expectedJSON []byte
	}{
		{
			name: "proxy-integration-route.json",
			args: args{
				fnName: "HelloWorldFunction",
				intName: "HelloWorldFunctionProxyRouteIntegration",
				event:  &cf.HttpApiEvent{
					Name:    "ProxyRoute",
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
			expectedName: "HelloWorldFunctionProxyRouteRoute",
			expectedJSON: x.LoadTestFile("proxy-integration-route.json"),
		},
		{
			name: "simple-path-route.json",
			args: args{
				fnName: "HelloWorldFunction",
				intName: "HelloWorldFunctionSimplePathIntegration",
				event:  &cf.HttpApiEvent{
					Name:    "SimplePath",
					HttpApi: cf.HttpApi{
						ImportApiId:          "ImportApiId",
						Auth:                 cf.Auth{},
						Method:               "GET",
						Path:                 "/hello",
						PayloadFormatVersion: "2.0",
						TimeoutInMillis:      5000,
					},
				},
			},
			expectedName: "HelloWorldFunctionSimplePathRoute",
			expectedJSON: x.LoadTestFile("simple-path-route.json"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rb := cf.NewRouteBuilder(tt.args.fnName, tt.args.intName, tt.args.event)
			j, err := rb.JSON()
			assert.NoError(t, err)
			require.NotNil(t, j)
			assert.JSONEq(t, string(tt.expectedJSON), string(j))
			assert.Equal(t, tt.expectedName, rb.Name())
		})
	}
}