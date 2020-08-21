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
	os.Setenv("EVENTSOURCE_EC2", "ec2.amazonaws.com")
	os.Setenv("EVENTNAME_EC2_1", "^RunInstances$")
	os.Setenv("EVENTNAME_EC2_2", "^StartInstances$")
	os.Setenv("EVENTNAME_EC2_3", "^StartInstances$")
	os.Setenv("EVENTNAME_EC2_4", "^StopInstances")
	os.Setenv("EVENTNAME_EC2_5", "^Modify.*")
	os.Setenv("EVENTNAME_EC2_6", "^Delete.*")
	os.Setenv("EVENTNAME_EC2_7", "^Update.*")
	os.Setenv("EVENTNAME_EC2_8", "^Terminate.*")
	os.Setenv("IGNOREARN_EC2_1", "arn:aws:iam::088921318242:user/dmai")
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
