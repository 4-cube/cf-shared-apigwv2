package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(input MacroInput) (*MacroOutput, error) {
	macro := NewMacro(input.Fragment)
	fragment, err := macro.ProcessFragment()

	if err != nil {
		return nil, err
	}

	return &MacroOutput{
		Status:    MacroOutputStatusSuccess,
		RequestId: input.RequestId,
		Fragment:  fragment,
	}, nil
}

func main() {

	//json, err := ioutil.ReadFile("./test-templates/simple-shared-http-sam-import-value.json")
	//if err != nil {
	//	panic("Can't load JSON template")
	//}
	//
	//macro := NewMacro(json)
	//
	//macro.ProcessFragment()

	lambda.Start(handler)
}
