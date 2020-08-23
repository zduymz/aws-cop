package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

const SlackMessageSizeLimit = 2000 // actual is 100.000

type SlackRequestBody struct {
	Text string `json:"text"`
}

type Slack struct {
	WebHookUrl string
	DryRun     bool
}

func (s *Slack) Write(msg string) error {
	if s.DryRun {
		log.Println(msg)
		return nil
	}

	slackBody, _ := json.Marshal(SlackRequestBody{Text: fmt.Sprintf("```%s```", msg)})
	req, err := http.NewRequest(http.MethodPost, s.WebHookUrl, bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return err
	}

	if buf.String() != "ok" {
		return errors.New("Non-OK response returned from Slack")
	}

	return nil
}
