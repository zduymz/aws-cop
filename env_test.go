package main

import (
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func setupTest() {
	for _, kv := range os.Environ() {
		if strings.HasPrefix(kv, "USERAGENT_") ||
			strings.HasPrefix(kv, "EVENTNAME_") ||
			strings.HasPrefix(kv, "USERIDENTIFY_") {
			os.Unsetenv(strings.Split(kv, "=")[0])
		}
	}
}

func TestReadEnvFull(t *testing.T) {
	setupTest()
	os.Setenv("EVENTNAME_DEFAULT_IGNORES", "Get,Describe,List,Head,ConsoleLogin")
	os.Setenv("USERIDENTITY_DEFAULT_IGNORES", "ec2.amazonaws.com")
	os.Setenv("USERAGENT_DEFAULT_MUSTHAVE", "console.amazonaws.com")
	os.Setenv("USERAGENT_DEFAULT_MUSTHAVEREGEX", "console.*.amazonaws.com,aws-internal*AWSLambdaConsole/*")

	expect := ConfigRule{
		EventName: ConfigEventName{
			Ignores: map[string][]string{
				"default": {"Get", "Describe", "List", "Head", "ConsoleLogin"},
			}},
		UserAgent: ConfigUserAgent{
			MustHave: map[string][]string{
				"default": {"console.amazonaws.com"},
			},
			MustHaveRegex: map[string][]*regexp.Regexp{
				"default": {regexp.MustCompile("console.*.amazonaws.com"), regexp.MustCompile("aws-internal*AWSLambdaConsole/*")},
			},
		},
		UserIdentity: ConfigUserIdentity{
			Ignores: map[string][]string{
				"default": {"ec2.amazonaws.com"},
			},
		},
	}

	actual := ReadEnv()
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("env parse failed.\nExpected: %+v.\nActual: %+v", expect, actual)
	}
}