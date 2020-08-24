package main

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func TestHandleRequest(t *testing.T) {
	os.Setenv("EVENTNAME_DEFAULT_IGNORES", "Get,Describe,List,Head,ConsoleLogin")
	os.Setenv("USERAGENT_DEFAULT_MUSTHAVE", "signin.amazonaws.com,console.amazonaws.com")
	os.Setenv("USERAGENT_DEFAULT_MUSTHAVEREGEX", ".*")
	os.Setenv("USERIDENTITY_DEFAULT_IGNORES", "ecs-tasks.amazonaws.com,ec2.amazonaws.com,monitoring.rds.amazonaws.com,lambda.amazonaws.com,AWSServiceRoleForEC2SpotFleet")
	os.Setenv("EVENTNAME_ec2amazonawscom_IGNORES","CreateTags")
	os.Setenv("EVENTNAME_ssmamazonawscom_IGNORES","Update,Put")
	os.Setenv("EVENTNAME_stsamazonawscom_IGNORES","AssumeRole")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		t.Errorf("failed to create aws session. Reason: %v", err)
	}

	downloader = Downloader{
		s3manager.NewDownloader(sess),
	}
	s3Event := events.S3Event{Records: []events.S3EventRecord{
		{
			S3: events.S3Entity{
				Bucket: events.S3Bucket{
					Name: "apixio-cloudtrail-logs",
				},
				Object: events.S3Object{
					Key: "AWSLogs/088921318242/CloudTrail/us-west-2/2020/08/21/088921318242_CloudTrail_us-west-2_20200821T0020Z_e3TWhqa1AP8ku6Pz.json.gz",
				},
			},
		},
	}}

	HandleRequest(context.Background(), s3Event)
}
