#SHELL=/bin/bash -o pipefail

.PHONY: deps clean test build deploy bucket

S3_BUCKET=4-cube-serverless-projects

deps:
	go get -u ./...

clean: 
	rm -rf ./cf-shared-apigwv2

test:
	go test ./...

build: deps test
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ./cf-shared-apigwv2 .
	sam build

bucket:
	if ! aws s3api head-bucket --bucket ${S3_BUCKET} 2>/dev/null; then \
	aws s3api create-bucket --bucket ${S3_BUCKET} --region eu-west-1 --create-bucket-configuration LocationConstraint=eu-west-1; \
	fi

deploy: bucket
	sam deploy --stack-name cf-shared-apigwv2-macro \
 		--no-confirm-changeset \
 		--capabilities CAPABILITY_IAM \
 		--region eu-west-1 \
 		--s3-bucket ${S3_BUCKET} \
 		--s3-prefix cf-shared-apigwv2

