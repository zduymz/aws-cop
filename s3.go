package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)


func (d *Downloader) Get(bucket, key string) (*CloudTrailEvents, error) {
	buff := &aws.WriteAtBuffer{}
	events := &CloudTrailEvents{}

	if _, err := d.Download(buff,
		&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(key),
		}); err != nil {
		return nil, fmt.Errorf("unable to download item %q. Reason: %v", key, err)
	}

	gzipReader, err := gzip.NewReader(bytes.NewBuffer(buff.Bytes()))
	if err != nil {
		return nil, err
	}

	var data bytes.Buffer
	if _, err = data.ReadFrom(gzipReader); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data.Bytes(), events); err != nil {
		return nil, fmt.Errorf("cloudtrail unmarshal failed. Reason: %v", err)
	}

	//fmt.Println(events)

	return events, nil
}
