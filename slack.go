package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/spf13/viper"
	"net/http"
)

var SlackError = errors.New("non-ok response returned from Slack")

func sendSlackNotification(msg string) error {
	slackBody, _ := json.Marshal(slackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, viper.GetString("WEBHOOK"), bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: DEFAULT_TIMEOUT}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return SlackError
	}

	return nil
}
