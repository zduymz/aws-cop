package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var (
	downloader Downloader
	slack      Slack
	//db         *dynamodb.DynamoDB
	//dynamodbTable string
)

func HandleRequest(ctx context.Context, e events.S3Event) {
	config := ReadEnv()
	lc, _ := lambdacontext.FromContext(ctx)
	log.Printf("Handler start: %s - records len: %d", lc.AwsRequestID, len(e.Records))
	for _, r := range e.Records {

		//lock, _ := dynamodbattribute.MarshalMap(DynamodbLock{
		//	Key: r.S3.Object.Key,
		//})
		//if _, err := db.PutItem(&dynamodb.PutItemInput{
		//	ConditionExpression: aws.String("attribute_not_exists(id)"),
		//	TableName:           aws.String("lambda-cloudtrail-watcher"),
		//	Item:                lock,
		//}); err != nil {
		//	if err.(awserr.Error).Code() == dynamodb.ErrCodeConditionalCheckFailedException {
		//		log.Println("passed")
		//		continue
		//	}
		//	log.Printf("failed to write to dynamodb. Reason: %v", err)
		//	continue
		//}

		log.Printf("begin processing: %s/%s", r.S3.Bucket.Name, r.S3.Object.Key)
		var tickets []string
		trail, err := downloader.Get(r.S3.Bucket.Name, r.S3.Object.Key)
		if err != nil {
			log.Printf("failed on s3. %v", err)
			continue
		}

		for _, event := range trail.Events {
			if config.checkWhitelistEvent(event.EventSource, event.EventName) {
				continue
			}

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
			case "AWSAccount":
				actor = event.UserIdentity.AccountId
			case "AWSService":
				actor = event.UserIdentity.InvokedBy
				// ingore api call by aws service
				continue
			default:
				actor = "unknown"
			}

			if config.checkWhiteListUserIdentity(event.EventSource, actor) {
				continue
			}

			if config.checkBlacklistUseragent(event.EventSource, event.UserAgent) {
				tickets = append(tickets, fmt.Sprintf("[%s] %s - %s - %s - %s", event.EventTime, actor, event.EventName, strings.TrimSuffix(event.EventSource, ".amazonaws.com"), event.EventId))
			}

			// log to cloudwatch for improving
			//log.Printf("%s - %s - %s - %s", actor, event.EventName, event.EventSource, event.UserAgent)
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
			if err := slack.Write(strings.Join(tickets[begin:], "\n")); err != nil {
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

	//db = dynamodb.New(sess)

	//if _, ok := os.LookupEnv("DEBUG"); !ok {
	//	log.SetOutput(ioutil.Discard)
	//}
}

func main() {
	lambda.Start(HandleRequest)
}
