![](https://github.com/4-cube/cf-shared-apigwv2/workflows/Build/badge.svg)
![](https://goreportcard.com/badge/github.com/4-cube/cf-shared-apigwv2)

# AWS SAM Shared Api Gateway V2 (HttpApi) macro

The [AWS Serverless Application Model][sam] (SAM) defines a simplified framework
on top of CloudFormation for building serverless apps.
 
There are several issues that are IMHO inconvenient:
* Not supported: [Attach functions to externally defined API Gateway][gh-issue-apgw]
* Not supported: [Lambda as ALB target][gh-issue-alb] 
    * community solution can be found [sam-alb][sam-alb]
* Really annoying: [HttpApi event type does not work with local start-api][gh-issue-start-local]

We created this project in order to enable using externally defined API gateway.
Currently, we only support ApiGatewayV2 (HttpApi). 
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

Second, you need to add a reference to the macro in the list of transforms in your template. 

Replace this:
```yaml
Transform: "AWS::Serverless-2016-10-31"
```

With this:
```yaml
Transform: ["SharedApiGatewayV2", "AWS::Serverless-2016-10-31"]
```

Note that the order matters - `SharedApiGatewayV2` needs to come first. 

Finally, you can use a modified event type (`HttpApi`) type, but you must configure value of `ImportApiId`. 

```yaml
  Events:
    CatchAll: #this is just event name
      Type: HttpApi
      Properties:
        Path: /{proxy+}
        Method: ANY
        ImportApiId: http-apigw-HttpApi #Reference shared Api Gateway - exact exported value from parent stack
```

You can see a complete example in [`demo.yaml`](/demo.yaml).

The complete set of properties for the modified `HttpApi` event type are:

| Property name             | Description                                                                   |
| ------------------------- | ----------------------------------------------------------------------------- |
| `ImportApiId`             | Type: String. **Required**. Should be `shared-api-gateway-ref`                |
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

We want to thank [aidansteele](https://github.com/aidansteele) for his work on [sam-alb](https://github.com/glassechidna/sam-alb)
since we used sam-alb project as an inspiration.
 