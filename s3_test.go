package main

import (
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func TestS3Download(t *testing.T) {
	bucket := "apixio-cloudtrail-logs"
	key := "AWSLogs/088921318242/CloudTrail/us-west-2/2020/08/21/088921318242_CloudTrail_us-west-2_20200821T0010Z_5QHZfpRjEuckvf9N.json.gz"

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		t.Errorf("failed to create aws session. Reason: %v", err)
	}

	downloader = Downloader{
		s3manager.NewDownloader(sess),
	}

	events, err := downloader.Get(bucket, key)
	if err != nil {
		t.Error(err)
	}

	if len(events.Events) == 0 {
		t.Error("failed to parse data")
	}
}
