package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/spf13/viper"
	"github.com/whalebone/go-dockerhub-ci/model"
)

var SlackError = errors.New("non-ok response returned from Slack")

func sendSlackNotification(msg string) error {
	slackBody, _ := json.Marshal(model.CreateSlackMessage(msg))
	req, err := http.NewRequest(http.MethodPost, viper.GetString("WEBHOOK"), bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: model.DefaultTimeout}
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
