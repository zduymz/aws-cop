package main

import (
	"os"
	"regexp"
	"strings"
)

/*
services:
  - eventsource: ec2.amazonaws.com
    eventnames:
      - ^RunInstances$
      - ^TerminateInstances$
    ignoredARNs:
      - arn:aws:iam::1234:user/dmai
      - arn:aws:sts::1234:assumed-role/AWSServiceRoleForAutoScaling/AutoScaling
*/

func ReadEnv() map[string]*ConfigRule {
	// EVENTSOURCE_EC2=ec2.amazonaws.com
	// EVENTNAME_EC2_1=^RunInstances$
	// EVENTNAME_EC2_2=^TerminateInstances$
	// IGNOREARN_EC2_1=arn:aws:iam::1234:user/dmai
	// IGNOREARN_EC2_2=arn:aws:sts::1234:assumed-role/AWSServiceRoleForAutoScaling/AutoScaling

	config := map[string]*ConfigRule{}

	for _, kv := range os.Environ() {
		if strings.HasPrefix(kv, "EVENTSOURCE_") {
			raw := strings.Split(kv, "=")
			if len(raw) != 2 {
				continue
			}
			key := raw[0]
			value := raw[1]

			raw = strings.Split(key, "_")
			if _, ok := config[raw[1]]; !ok {

				config[raw[1]] = &ConfigRule{
					EventSource: value,
					EventNames:  []*regexp.Regexp{},
					IgnoreARNs:  []string{},
				}
			} else {
				config[raw[1]].EventSource = value
			}
		}

		if strings.HasPrefix(kv, "EVENTNAME_") {
			raw := strings.Split(kv, "=")
			if len(raw) != 2 {
				continue
			}
			key := raw[0]
			value := raw[1]

			raw = strings.Split(key, "_")
			if _, ok := config[raw[1]]; !ok {
				config[raw[1]] = &ConfigRule{
					EventNames: []*regexp.Regexp{regexp.MustCompile(value)},
					IgnoreARNs: []string{},
				}
			} else {
				config[raw[1]].EventNames = append(config[raw[1]].EventNames, regexp.MustCompile(value))
			}
		}

		if strings.HasPrefix(kv, "IGNOREARN_") {
			raw := strings.Split(kv, "=")
			if len(raw) != 2 {
				continue
			}
			key := raw[0]
			value := raw[1]

			raw = strings.Split(key, "_")
			if _, ok := config[raw[1]]; !ok {
				config[raw[1]] = &ConfigRule{
					EventNames: []*regexp.Regexp{},
					IgnoreARNs: []string{value},
				}
			} else {
				config[raw[1]].IgnoreARNs = append(config[raw[1]].IgnoreARNs, value)
			}
		}
	}

	return config
}
