# AWS SAM Shared Api Gateway V2 (HttpApi) macro

The [AWS Serverless Application Model][sam] (SAM) defines a simplified framework
on top of CloudFormation for building serverless apps.
 
There are several issues that are IMHO inconvenient:
* Not supported: [Attach functions to externally defined API Gateway][gh-issue-apgw]
* Not supported: [Lambda as ALB target][gh-issue-alb] 
    * community solution can be found [sam-alb][sam-alb]
* Really annoying: [HttpApi event type does not work with local start-api][gh-issue-start-local]

We created this project in order to enable using externally defined API gateway.
Currently, we only supports ApiGatewayV2 (HttpApi). 
You can learn more about ClodFormation Macros [here][macro].
  
## Installation
In order to be able to use this Macro functions, you have to install it in every region which is used for SAM Lambda
projects.
```shell script
make build
make deploy
```

## Usage

Usage is pretty simple. You can see a complete example in [`demo.yaml`](/demo.yaml).
There are two parts to it. First, you need to add a reference to the macro in
the list of transforms in your template. 

Replace this:

    Transform: "AWS::Serverless-2016-10-31"
    
With this:

    Transform: ["SharedApiGatewayV2", "AWS::Serverless-2016-10-31"]

Note that the order matters - `SharedApiGatewayV2` needs to come first. Next, you can now
use a new event type. Here's what that looks like:

```yaml
  Function:
    Type: AWS::Serverless::Function
    Properties:
      Handler: index.handler
      Runtime: python3.7
      Events:
        hello:
          Type: SharedHttpApi
          Properties:
            Type: SharedHttpApi
            Properties:
              Path: /hello/{proxy+}
              Method: ANY
              ApiId: !ImportValue http-apigw-HttpApi #Replace http-apigw-HttpApi with your own
```

The complete set of properties for the `SharedHttpApi` event type are:
ApiId                json.RawMessage
	Auth                 HttpApiFunctionAuth
	Method               string
	Path                 string
	PayloadFormatVersion string
	TimeoutInMillis      int64

| Property name | Description                                                  |
| ------------- | ------------------------------------------------------------ |
| `ApiId`       | Type: String. **Required**. Should be the import value of an shared ApiGateway (Fn:ImportValue |
| `Method`      | Type: String. One of: GET, POST, PUT, ANY |
| `Path`        | Type: |
| `PayloadFormatVersion`  | Type: String.  |
| `TimeoutInMillis` | Type: Integer |
| `Auth` | Type 


[sam]: https://github.com/awslabs/serverless-application-model
[gh-issue-alb]: https://github.com/awslabs/serverless-application-model/issues/721
[gh-issue-apgw]: https://github.com/awslabs/serverless-application-model/issues/149
[gh-issue-start-local]: https://github.com/awslabs/aws-sam-cli/issues/1641
[sam-alb]: https://github.com/glassechidna/sam-alb
[macro]: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-macros.html

