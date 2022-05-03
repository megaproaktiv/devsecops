package cicd

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscodebuild"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscodecommit"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/constructs-go/constructs/v10"
)

type CicdStackProps struct {
	awscdk.StackProps
}

func NewCicdStack(scope constructs.Construct, id string, props *CicdStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	this := awscdk.NewStack(scope, &id, &sprops)

	policy1 := awsiam.NewPolicyDocument(&awsiam.PolicyDocumentProps{
		Minimize: aws.Bool(true),
		Statements: &[]awsiam.PolicyStatement{
			awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
				Actions: &[]*string{
					aws.String("logs:CreateLogGroup"),
					aws.String("logs:CreateLogStream"),
					aws.String("logs:PutLogEvents"),
				},
				Resources: &[]*string{
					aws.String("*"),
				},
			}),
			},
		})
	policy2 := awsiam.NewPolicyDocument(&awsiam.PolicyDocumentProps{
		Minimize: aws.Bool(true),
		Statements: &[]awsiam.PolicyStatement{
			awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
				Actions: &[]*string{
					aws.String("codecommit:GitPull"),
				},
				Resources: &[]*string{
					aws.String("*"),
				},
			}),
			},
		})
	policy3 := awsiam.NewPolicyDocument(&awsiam.PolicyDocumentProps{
		Minimize: aws.Bool(true),
		Statements: &[]awsiam.PolicyStatement{
			awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
				Actions: &[]*string{
					aws.String("codecommit:GitPull"),
					aws.String("codebuild:CreateReportGroup"),
					aws.String("codebuild:CreateReport"),
					aws.String("codebuild:UpdateReport"),
					aws.String("codebuild:BatchPutTestCases"),
					aws.String("codebuild:BatchPutCodeCoverages"),
				},
				Resources: &[]*string{
					aws.String("*"),
				},
			}),
			},
		})
	policy4 := awsiam.NewPolicyDocument(&awsiam.PolicyDocumentProps{
		Minimize: aws.Bool(true),
		Statements: &[]awsiam.PolicyStatement{
			awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
				Actions: &[]*string{
					aws.String("s3:PutObject"),
					aws.String("s3:GetObject"),
					aws.String("s3:GetObjectVersion"),
					aws.String("s3:GetBucketAcl"),
					aws.String("s3:GetBucketLocation"),
				},
				Resources: &[]*string{
					aws.String("*"),
				},
			}),
			},
		})
	policy5 := awsiam.NewPolicyDocument(&awsiam.PolicyDocumentProps{
		Minimize: aws.Bool(true),
		Statements: &[]awsiam.PolicyStatement{
			awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
				Actions: &[]*string{
					aws.String("iam:PassRole"),
				},
				Resources: &[]*string{
					aws.String("*"),
				},
			}),
			},
		})
	policy6 := awsiam.NewPolicyDocument(&awsiam.PolicyDocumentProps{
		Minimize: aws.Bool(true),
		Statements: &[]awsiam.PolicyStatement{
			awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
				Actions: &[]*string{
					aws.String("cloudformation:*"),
				},
				Resources: &[]*string{
					aws.String("*"),
				},
			}),
			},
		})
		
	policy7 := awsiam.NewPolicyDocument(&awsiam.PolicyDocumentProps{
		Minimize: aws.Bool(true),
		Statements: &[]awsiam.PolicyStatement{
			awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
				Actions: &[]*string{
					aws.String("dynamodb:DeleteItem"),
				},
				Resources: &[]*string{
					aws.String("*"),
				},
			}),
			},
		})



	// The code that defines your stack goes here
	servicerole := awsiam.NewRole(this, aws.String("servicerole"), &awsiam.RoleProps{
		AssumedBy:           awsiam.NewServicePrincipal(aws.String("codebuild.amazonaws.com"), nil),
		Description:         aws.String("ServiceRole for codebuild devsecops"),		
		InlinePolicies:      &map[string]awsiam.PolicyDocument{
			  "logs": policy1,
			  "git": policy2,
			  "report": policy3,
			  "s3": policy4,
			  "iam": policy5,
			  "cfn": policy6,
			  "dynamodb": policy7,
		},
		Path:                aws.String("/letsbuild/"),
		RoleName:            aws.String("devsecops-codebuild-service-role"),
	})
	servicerole.AddManagedPolicy(awsiam.ManagedPolicy_FromManagedPolicyArn(this, aws.String("readall"), aws.String("arn:aws:iam::aws:policy/ReadOnlyAccess")))
	
	awscodebuild.NewProject(this, aws.String("buildproject"), &awscodebuild.ProjectProps{
		BuildSpec:                           awscodebuild.BuildSpec_FromSourceFilename(aws.String("buildspec_sec.yml")),
		CheckSecretsInPlainTextEnvVariables: aws.Bool(true),
		ConcurrentBuildLimit:                aws.Float64(1),
		Description:                         aws.String("Dev Secops - build with steampipe AWS foundational_security"),
		Environment:                         &awscodebuild.BuildEnvironment{
			BuildImage:           awscodebuild.LinuxBuildImage_STANDARD_5_0(),
			ComputeType:          awscodebuild.ComputeType_SMALL,
			EnvironmentVariables: &map[string]*awscodebuild.BuildEnvironmentVariable{},
			Privileged:           aws.Bool(false),
		},
		EnvironmentVariables:                &map[string]*awscodebuild.BuildEnvironmentVariable{},
		Logging:                             &awscodebuild.LoggingOptions{
			CloudWatch: &awscodebuild.CloudWatchLoggingOptions{
				Enabled:  aws.Bool(true),
				LogGroup: awslogs.NewLogGroup(this, aws.String("builddevsecops"), &awslogs.LogGroupProps{
					LogGroupName:  aws.String("builddevsecops"),
					Retention:     awslogs.RetentionDays_THREE_MONTHS,
					RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
				}),
				Prefix:   aws.String("codebuild"),
			},
		},
		ProjectName:                         aws.String("devsecops"),
		QueuedTimeout:                       awscdk.Duration_Minutes(aws.Float64(30)),
		Role:                                servicerole,
		Timeout:                             awscdk.Duration_Minutes(aws.Float64(30)),
		Artifacts:                           nil,
		Source:                              awscodebuild.Source_CodeCommit( &awscodebuild.CodeCommitSourceProps{
			Identifier:      aws.String("source"),
			Repository:      awscodecommit.Repository_FromRepositoryName(this, aws.String("coderepo"), aws.String("devsecops")),
			BranchOrRef:     aws.String("refs/heads/main"),
			CloneDepth:      aws.Float64(1),
			FetchSubmodules: aws.Bool(false),
		}),
	})

	return this
}
