package cf_test

import (
	"github.com/4-cube/cf-shared-apigwv2/pkg/cf"
	"github.com/stretchr/testify/assert"
	"testing"
)

var event = &cf.HttpApiEvent{
	Name: "MyHttpApiEvent",
	HttpApi: cf.HttpApi{
		ImportApiId:          "MyApiId",
		Auth:                 cf.Auth{},
		Method:               "GET",
		Path:                 "/",
		PayloadFormatVersion: "2.0",
		TimeoutInMillis:      5000,
	},
}

func TestHttpApiEvent_FnImportApiId(t *testing.T) {
	expected := `{"Fn::ImportValue": "MyApiId"}`

	t.Run("FnImportApiId", func(t *testing.T) {
		assert.JSONEq(t, expected, string(event.FnImportApiId()))
	})
}

func TestHttpApiEvent_RouteKey(t *testing.T) {
	expected := "GET /"

	t.Run("RouteKey", func(t *testing.T) {
		assert.Equal(t, expected, event.RouteKey())
	})
}
