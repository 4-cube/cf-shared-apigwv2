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
###### There are some hardcoded values in Makefile - please review Makefile before running it.


## Usage
First, you need to have already provisioned ApiGatewayV2 instance.
To be able to reference this Api Gateway instance, you `must export Resource Reference`.
 
```yaml
Outputs:
  HttpApi:
    Description: "API Gateway Resource"
    Value: !Ref HttpApi
    Export:
      Name: !Sub "${AWS::StackName}-HttpApi"
```
You can find example in [`demo-apigwv2.yaml`](/demo-apigwv2.yaml). 

Secondly, you need to add a reference to the macro in the list of transforms in your template. 

Replace this:
```yaml
Transform: "AWS::Serverless-2016-10-31"
```

With this:
```yaml
Transform: ["SharedApiGatewayV2", "AWS::Serverless-2016-10-31"]
```

Note that the order matters - `SharedApiGatewayV2` needs to come first. 

Finally, you can use a new event (`SharedHttpApi`) type. 

```yaml
  Events:
    CatchAll: #this is just event name
      Type: SharedHttpApi #We specify new event type
      Properties:
        Path: /{proxy+}
        Method: ANY
        ApiId: !ImportValue http-apigw-HttpApi #Reference shared Api Gateway
```

You can see a complete example in [`demo.yaml`](/demo.yaml).

The complete set of properties for the `SharedHttpApi` event type are:
ApiId                json.RawMessage
	Auth                 HttpApiFunctionAuth
	Method               string
	Path                 string
	PayloadFormatVersion string
	TimeoutInMillis      int64

| Property name             | Description                                                                   |
| ------------------------- | ----------------------------------------------------------------------------- |
| `ApiId`                   | Type: String. **Required**. Should be `Fn:ImportValue shared-api-gateway-ref` |
| `Method`                  | Type: String. If empty or omitted `ANY` will be used                          |
| `Path`                    | Type: String.                                                                 |
| `PayloadFormatVersion`    | Type: String. Default value: `2.0`                                            |
| `TimeoutInMillis`         | Type: Integer. Default: `5000`                                                |
| `Auth` | Type: Auth. **Currently not implemented.**                                                       |

[sam]: https://github.com/awslabs/serverless-application-model
[gh-issue-alb]: https://github.com/awslabs/serverless-application-model/issues/721
[gh-issue-apgw]: https://github.com/awslabs/serverless-application-model/issues/149
[gh-issue-start-local]: https://github.com/awslabs/aws-sam-cli/issues/1641
[sam-alb]: https://github.com/glassechidna/sam-alb
[macro]: https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/template-macros.html

