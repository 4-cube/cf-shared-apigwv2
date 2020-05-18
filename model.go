package main

import "encoding/json"

type MacroInput struct {
	Region      string
	AccountId   string
	TransformId string
	RequestId   string
	Fragment    json.RawMessage
	Params      json.RawMessage
	Values      map[string]json.RawMessage `json:"templateParameterValues"`
}

type MacroOutput struct {
	Status    string          `json:"status"`
	RequestId string          `json:"requestId"`
	Fragment  json.RawMessage `json:"fragment"`
}

const MacroOutputStatusSuccess = "success"
const MacroOutputStatusFailure = "failure"

type HttpApiFunctionAuth struct {
	AuthorizationScopes []string
	Authorizer          string
}

type SharedHttpApi struct {
	ApiId                json.RawMessage
	Auth                 HttpApiFunctionAuth
	Method               string
	Path                 string
	PayloadFormatVersion string
	TimeoutInMillis      int64
}

type FunctionEvent struct {
	Name          string
	SharedHttpApi SharedHttpApi
}
