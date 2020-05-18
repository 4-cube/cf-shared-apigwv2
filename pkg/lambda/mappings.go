package lambda

import (
	"encoding/json"
)

// Macro Function Request Event Mapping
// see: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-macros.html
type MacroRequest struct {
	// The region in which the macro resides.
	Region string `json:"region"`

	// The account ID of the account from which the macro is invoking the Lambda function.
	AccountId string `json:"accountId"`

	// The name of the macro invoking this function.
	TransformId string `json:"transformId"`

	// The ID of the request invoking this function.
	RequestId string `json:"requestId"`

	// The template content available for custom processing, in JSON format.
	// For macros included in the Transform template section, this is the entire template except for the Transform section.
	// For macros included in an Fn::Transform intrinsic function call, this includes all sibling nodes (and their children)
	// based on the location of the intrinsic function within the template except for the Fn::Transform function.
	// For more information, see AWS CloudFormation Macro Scope.
	Fragment json.RawMessage `json:"fragment"`

	// For Fn::Transform function calls, any specified parameters for the function.
	// AWS CloudFormation does not evaluate these parameters before passing them to the function.
	// For macros included in the Transform template section, this section is empty.
	Params json.RawMessage `json:"params"`

	// Any parameters specified in the Parameters section of the template.
	// AWS CloudFormation evaluates these parameters before passing them to the function.
	ParamValues map[string]json.RawMessage `json:"templateParameterValues"`
}

// Macro Function Response Mapping
// see: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-macros.html
type MacroResponse struct {
	// The status of the request (case-insensitive). Should be set to "success".
	// AWS CloudFormation treats any other response as a failure.
	Status string `json:"status"`

	// The ID of the request invoking this function.
	// This must match the request ID provided by AWS CloudFormation when invoking the function.
	RequestId string `json:"requestId"`

	// The processed template content for AWS CloudFormation to include in the processed template, including siblings.
	// AWS CloudFormation replaces the template content that is passed to the Lambda function with the template fragment
	// it receives in the Lambda response.
	//
	// The processed template content must be valid JSON, and its inclusion in the processed template must result in a valid template.
	//
	// If your function doesn't actually change the template content that AWS CloudFormation passes to it,
	// but you still need to include that content in the processed template, your function needs to return that template
	// content to AWS CloudFormation in its response.
	Fragment json.RawMessage `json:"fragment"`
}

const MacroOutputStatusSuccess = "success"
