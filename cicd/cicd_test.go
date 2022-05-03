package cicd_test

import (
	"cicd"
	"encoding/json"
	"os"
	"testing"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestCicdStack(t *testing.T) {
	// GIVEN
	app := awscdk.NewApp(nil)

	// WHEN
	stack := cicd.NewCicdStack(app, "MyStack",
		&cicd.CicdStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	// THEN
	bytes, err := json.Marshal(app.Synth(nil).GetStackArtifact(stack.ArtifactId()).Template())
	if err != nil {
		t.Error(err)
	}

	template := gjson.ParseBytes(bytes)
	name := template.Get("Resources.buildproject87EEBE72.Properties.Name").String()
	assert.Equal(t, "devsecops", name)
}

func env() *awscdk.Environment {

	
	return &awscdk.Environment{
	 Account: aws.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	 Region:  aws.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
