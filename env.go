package main

import (
	"os"
	"regexp"
	"strings"
)

func ReadEnv() ConfigRule {
	/*
		EVENTNAME_DEFAULT_IGNORES=Get,Describe,List,Head,ConsoleLogin
		USERAGENT_DEFAULT_MUSTHAVE=console.amazonaws.com
		USERAGENT_DEFAULT_MUSTHAVEREGEX=console.*.amazonaws.com,aws-internal*AWSLambdaConsole/*
		USERIDENTIFY_DEFAULT_IGNORES=AWS Internal
	*/

	config := ConfigRule{
		EventName: ConfigEventName{
			Ignores: map[string][]string{}},
		UserAgent: ConfigUserAgent{
			MustHave:      map[string][]string{},
			MustHaveRegex: map[string][]*regexp.Regexp{},
		},
		UserIdentity: ConfigUserIdentity{
			Ignores: map[string][]string{},
		},
	}

	for _, kv := range os.Environ() {
		key := strings.Split(kv, "=")[0]
		values := strings.Split(strings.Split(kv, "=")[1], ",")
		if strings.HasPrefix(kv, "EVENTNAME_") {

			if raw := strings.Split(key, "_"); len(raw) == 3 && len(values) > 0 {
				svc := strings.ToLower(raw[1])
				config.EventName.Ignores[svc] = values
			}
		}

		if strings.HasPrefix(kv, "USERIDENTITY_") {
			if raw := strings.Split(key, "_"); len(raw) == 3 && len(values) > 0 {
				svc := strings.ToLower(raw[1])
				config.UserIdentity.Ignores[svc] = values
			}
		}

		if strings.HasPrefix(kv, "USERAGENT_") {
			key := strings.Split(kv, "=")[0]
			if raw := strings.Split(key, "_"); len(raw) == 3 && len(values) > 0 {
				svc := strings.ToLower(raw[1])
				if strings.HasSuffix(strings.ToLower(raw[2]), "regex") {
					for _, v := range values {
						if v != "" {
						config.UserAgent.MustHaveRegex[svc] = append(config.UserAgent.MustHaveRegex[svc], regexp.MustCompile(v))
						}
					}
				} else {
					config.UserAgent.MustHave[svc] = values
				}
			}
		}
	}

	return config
}
