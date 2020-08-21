package main

import (
	"regexp"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-event-reference-user-identity.html
type UserIdentity struct {
	// Type can be following values:
	// Root, IAMUser, AssumedRole, FederatedUser, AWSAccount, AWSService
	Type           string         `json:"type"`
	PrincipalId    string         `json:"principalId"`
	ARN            string         `json:"arn"`
	AccountId      string         `json:"accountId"`
	AccesskeyId    string         `json:"accesskeyId"`
	UserName       string         `json:"username"`
	SessionContext SessionContext `json:"sessionContext"`
	InvokedBy      string         `json:"invokedBy"`
}

/*
   "sessionContext": {
     "sessionIssuer": {
       "type": "Role",
       "principalId": "ABC",
       "arn": "arn:aws:iam::123:role/roleName",
       "accountId": "12345",
       "userName": "roleName"
     },
     "webIdFederationData": {},
     "attributes": {
       "mfaAuthenticated": "false",
       "creationDate": "2020-08-20T23:29:53Z"
     },
     "ec2RoleDelivery": "2.0"
   }
*/
type SessionContext struct {
	SessionIssuer SessionIssuer `json:"sessionIssuer"`
}

type SessionIssuer struct {
	Type        string `json:"type"`
	PrincipalId string `json:"principalId"`
	ARN         string `json:"arn"`
	AccountId   string `json:"accountId"`
	UserName    string `json:"userName"`
}

type CloudtrailEvent struct {
	UserIdentity UserIdentity `json:"userIdentity"`
	EventTime    string       `json:"eventTime"`
	EventSource  string       `json:"eventSource"`
	EventName    string       `json:"eventName"`
	AWSRegion    string       `json:"awsRegion"`
	EventType    string       `json:"eventType"`
	ReadOnly     bool         `json:"readOnly"`
	//requestParameters string
	//responseElements  string
}

type CloudTrailEvents struct {
	Events []CloudtrailEvent `json:"Records"`
}

type ConfigRule struct {
	EventSource string
	EventNames  []*regexp.Regexp
	IgnoreARNs  []string
}

type Downloader struct {
	*s3manager.Downloader
}
