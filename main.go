package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func inSlice(x string, xs []string) bool {
	for _, v := range xs {
		if x == v {
			return true
		}
	}
	return false
}

var (
	downloader Downloader
	slack      Slack
)

func HandleRequest(ctx context.Context, e events.S3Event) {
	configs := ReadEnv()

	log.Printf("process: %d", len(e.Records))
	for _, r := range e.Records {
		var tickets []string
		trail, err := downloader.Get(r.S3.Bucket.Name, r.S3.Object.Key)
		if err != nil {
			log.Printf("failed on s3. %v", err)
			continue
		}

		for _, event := range trail.Events {
			for _, rule := range configs {
				// check eventsource is in check list
				if rule.EventSource == event.EventSource {
					//check useridentity.arn is in ignore list
					if !inSlice(event.UserIdentity.ARN, rule.IgnoreARNs) {
						// check eventname is in check list
						for _, r := range rule.EventNames {
							if r.MatchString(event.EventName) {
								var actor string
								switch event.UserIdentity.Type {
								case "IAMUser":
									actor = event.UserIdentity.UserName
								case "AssumedRole":
									actor = event.UserIdentity.SessionContext.SessionIssuer.UserName
								case "Root":
									actor = "root"
								case "FederatedUser":
									actor = event.UserIdentity.SessionContext.SessionIssuer.UserName
								case "AWSService":
									actor = event.UserIdentity.InvokedBy
								case "AWSAccount":
									actor = event.UserIdentity.AccountId
								}

								tickets = append(tickets, fmt.Sprintf("[%s] %s - %s - %s", event.EventTime, actor, event.EventName, event.EventSource))
							}
						}
					}
				}
			}
		}
		// ready to send
		count := 0
		begin := 0
		end := len(tickets)
		for i, ticket := range tickets {
			count = count + len(ticket)
			end = i
			if count > SlackMessageSizeLimit {
				if err := slack.Write(strings.Join(tickets[begin:end], "\n")); err != nil {
					log.Printf("failed to send slack. Reason: %v", err)
				}
				begin = i
				count = 0
			}
		}
		if begin < end {
			if err := slack.Write(strings.Join(tickets[begin:end], "\n")); err != nil {
				log.Printf("failed to send slack. Reason: %v", err)
			}
		}
	}
}

func init() {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	if err != nil {
		log.Fatalf("failed to create aws session. Reason: %v", err)
	}

	downloader = Downloader{
		s3manager.NewDownloader(sess),
	}

	url, ok := os.LookupEnv("SLACK_WEBHOOK")
	slack = Slack{
		WebHookUrl: url,
		DryRun:     !ok,
	}

}

func main() {
	lambda.Start(HandleRequest)
}
