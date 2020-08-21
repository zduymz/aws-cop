package main

import (
	"os"
	"reflect"
	"regexp"
	"testing"
)

func TestReadEnvFull(t *testing.T) {
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
		t.Errorf("env parse failed. Expected: %v. Actual: %v", expect, actual)
	}
}

func TestReadEnvMissIgnoreARN(t *testing.T) {
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
		t.Errorf("env parse failed. Expected: %v. Actual: %v", expect, actual)
	}
}

func TestReadEnvMissEventName(t *testing.T) {
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
		t.Errorf("env parse failed. Expected: %v. Actual: %v", expect, actual)
	}
}

func TestReadEnvMissEventSource(t *testing.T) {
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
		t.Errorf("env parse failed. Expected: %v. Actual: %v", expect, actual)
	}
}
