package main

import (
	"cicd"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
)

func main() {
	app := awscdk.NewApp(nil)

	cicd.NewCicdStack(app, "CicdStack", &cicd.CicdStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {

	
	return &awscdk.Environment{
	 Account: aws.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	 Region:  aws.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
