package main

import (
	macro2 "github.com/4-cube/cf-shared-apigwv2/macro"
	l "github.com/4-cube/cf-shared-apigwv2/pkg/lambda"
	"github.com/4-cube/cf-shared-apigwv2/x"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(r l.MacroRequest) (*l.MacroResponse, error) {
	macro := macro2.NewMacro(r.Fragment, x.NewLogger())
	fragment, err := macro.ProcessFragment()

	if err != nil {
		return nil, err
	}

	return &l.MacroResponse{
		Status:    l.MacroOutputStatusSuccess,
		RequestId: r.RequestId,
		Fragment:  fragment,
	}, nil
}

func main() {
	lambda.Start(handler)
}
