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
		if strings.HasPrefix(kv, "EVENTSOURCE_") ||
			strings.HasPrefix(kv, "EVENTNAME_") ||
			strings.HasPrefix(kv, "IGNOREARN_") {
			os.Unsetenv(strings.Split(kv, "=")[0])
		}
	}
}

func TestReadEnvFull(t *testing.T) {
	setupTest()
	os.Setenv("EVENTSOURCE_EC2", "ec2.amazonaws.com")
	os.Setenv("EVENTNAME_EC2_1", "^RunInstances$")
	os.Setenv("EVENTNAME_EC2_2", "^TerminateInstances$")
	os.Setenv("IGNOREARN_EC2_1", "arn:aws:iam::1234:user/dmai")
	os.Setenv("IGNOREARN_EC2_2", "arn:aws:sts::1234:assumed-role/AWSServiceRoleForAutoScaling/AutoScaling")

	expect := map[string]*ConfigRule{
		"EC2": {
			EventSource: "ec2.amazonaws.com",
			EventNames: []*regexp.Regexp{
				regexp.MustCompile("^RunInstances$"),
				regexp.MustCompile("^TerminateInstances$")},
			IgnoreARNs: []string{"arn:aws:iam::1234:user/dmai", "arn:aws:sts::1234:assumed-role/AWSServiceRoleForAutoScaling/AutoScaling"},
		},
	}

	actual := ReadEnv()
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("env parse failed.\nExpected: %v.\nActual: %v", expect["EC2"], actual["EC2"])
	}
}

func TestReadEnvMissIgnoreARN(t *testing.T) {
	setupTest()
	os.Setenv("EVENTSOURCE_EC2", "ec2.amazonaws.com")
	os.Setenv("EVENTNAME_EC2_1", "^RunInstances$")
	os.Setenv("EVENTNAME_EC2_2", "^TerminateInstances$")
	expect := map[string]*ConfigRule{
		"EC2": {
			EventSource: "ec2.amazonaws.com",
			EventNames: []*regexp.Regexp{
				regexp.MustCompile("^RunInstances$"),
				regexp.MustCompile("^TerminateInstances$")},
			IgnoreARNs: []string{},
		},
	}

	actual := ReadEnv()
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("env parse failed.\nExpected: %v.\nActual: %v", expect["EC2"], actual["EC2"])
	}
}

func TestReadEnvMissEventName(t *testing.T) {
	setupTest()
	os.Setenv("EVENTSOURCE_EC2", "ec2.amazonaws.com")
	os.Setenv("IGNOREARN_EC2_1", "arn:aws:iam::1234:user/dmai")
	expect := map[string]*ConfigRule{
		"EC2": {
			EventSource: "ec2.amazonaws.com",
			EventNames:  []*regexp.Regexp{},
			IgnoreARNs:  []string{"arn:aws:iam::1234:user/dmai"},
		},
	}

	actual := ReadEnv()
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("env parse failed.\nExpected: %v.\nActual: %v", expect["EC2"], actual["EC2"])
	}
}

func TestReadEnvMissEventSource(t *testing.T) {
	setupTest()
	os.Setenv("EVENTNAME_EC2_1", "^RunInstances$")
	os.Setenv("IGNOREARN_EC2_1", "arn:aws:iam::1234:user/dmai")

	expect := map[string]*ConfigRule{
		"EC2": {
			EventSource: "",
			EventNames: []*regexp.Regexp{regexp.MustCompile("^RunInstances$")},
			IgnoreARNs: []string{"arn:aws:iam::1234:user/dmai"},
		},
	}

	actual := ReadEnv()
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("env parse failed.\nExpected: %v.\nActual: %v", expect["EC2"], actual["EC2"])
	}
}
