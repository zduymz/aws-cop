package main

import (
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// https://docs.aws.amazon.com/awscloudtrail/latest/userguide/cloudtrail-event-reference-record-contents.html
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
	UserAgent    string       `json:"userAgent"`
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

type Downloader struct {
	*s3manager.Downloader
}

type ConfigUserIdentity struct {
	Ignores map[string][]string
}

type ConfigUserAgent struct {
	MustHave      map[string][]string
	MustHaveRegex map[string][]*regexp.Regexp
}

type ConfigEventName struct {
	Ignores map[string][]string
}

type ConfigRule struct {
	// event name of api request
	// map[aws-service]map[event-name-specific] true
	EventName ConfigEventName

	// useragent in api request.
	// map[aws-service]map[user-agent]
	UserAgent ConfigUserAgent

	//
	UserIdentity ConfigUserIdentity
}

/*
	return true if event name is good to ingore
*/
func (c *ConfigRule) checkWhitelistEvent(esource, ename string) bool {
	esource = strings.ReplaceAll(esource, ".", "")
	for _, apiAction := range c.EventName.Ignores["default"] {
		if strings.HasPrefix(strings.ToLower(ename), strings.ToLower(apiAction)) {
			return true
		}
	}

	if apiActions, ok := c.EventName.Ignores[esource]; ok {
		for _, apiAction := range apiActions {
			if strings.HasPrefix(strings.ToLower(ename), strings.ToLower(apiAction)) {
				return true
			}
		}
	}
	return false
}

/*
	return true if useragent is in a watch list
*/
func (c *ConfigRule) checkBlacklistUseragent(esource, useragent string) bool {
	esource = strings.ReplaceAll(esource, ".", "")
	for _, ua := range c.UserAgent.MustHave["default"] {
		if ua == useragent {
			return true
		}
	}

	if uas, ok := c.UserAgent.MustHave[esource]; ok {
		for _, ua := range uas {
			if ua == useragent {
				return true
			}
		}
	}

	for _, reg := range c.UserAgent.MustHaveRegex["default"] {
		if reg.MatchString(useragent) {
			return true
		}
	}

	if regs, ok := c.UserAgent.MustHaveRegex[esource]; ok {
		for _, reg := range regs {
			if reg.MatchString(useragent) {
				return true
			}
		}
	}

	return false
}

/*
	return true if user is good to ignore
*/
func (c *ConfigRule) checkWhiteListUserIdentity(esource, user string) bool {
	esource = strings.ReplaceAll(esource, ".", "")
	for _, u := range c.UserIdentity.Ignores["default"] {
		if u == user {
			return true
		}
	}

	if us, ok := c.UserIdentity.Ignores[esource]; ok {
		for _, u := range us {
			if u == user {
				return true
			}
		}
	}
	return false
}
