package cf_test

import (
	"github.com/4-cube/cf-shared-apigwv2/pkg/cf"
	"github.com/4-cube/cf-shared-apigwv2/x"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_PermissionBuilder(t *testing.T) {
	type args struct {
		fnName string
		event  *cf.HttpApiEvent
	}
	tests := []struct {
		name string
		args args
		expectedName string
		expectedJSON []byte
	}{
		{
			name: "proxy-integration-permission.json",
			args: args{
				fnName: "HelloWorldFunction",
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
			expectedName: "HelloWorldFunctionProxyRoutePermission",
			expectedJSON: x.LoadTestFile("proxy-integration-permission.json"),
		},
		{
			name: "simple-path-permission.json",
			args: args{
				fnName: "HelloWorldFunction",
				event:  &cf.HttpApiEvent{
					Name:    "SimpleRoute",
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
			expectedName: "HelloWorldFunctionSimpleRoutePermission",
			expectedJSON: x.LoadTestFile("simple-path-permission.json"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pb := cf.NewPermissionBuilder(tt.args.fnName, tt.args.event)
			j, err := pb.JSON()
			assert.NoError(t, err)
			require.NotNil(t, j)
			assert.JSONEq(t, string(tt.expectedJSON), string(j))
			assert.Equal(t, tt.expectedName, pb.Name())
		})
	}
}