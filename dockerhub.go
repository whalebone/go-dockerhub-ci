package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func dockerhubConfirmer(c echo.Context) error {
	payload := &dockerhubPayload{}
	err := c.Bind(payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONError(err))
	}

	pushedAtTime := time.Unix(int64(payload.PushData.PushedAt), 0)

	slackData := fmt.Sprintf(
		"*%s* pushed *%s:%s* at %s",
		payload.PushData.Pusher,
		payload.Repository.Name,
		payload.PushData.Tag,
		pushedAtTime,
	)

	c.Logger().Debug(slackData)
	err = sendSlackNotification(slackData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONError(err))
	}

	response := &successResponse{
		State:       "success",
		Description: "Slack notified",
		Context:     "A OK",
		TargetURL:   "",
	}

	err = confirmDockerhub(payload.CallbackURL, response)
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONError(err))
	}

	return c.JSON(http.StatusOK, response)
}

func confirmDockerhub(url string, payload *successResponse) error {
	responseBody, _ := json.Marshal(payload)

	req, err := http.NewRequestWithContext(
		context.Background(),
		http.MethodPost,
		url,
		bytes.NewBuffer(responseBody),
	)

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: DEFAULT_TIMEOUT}
	_, err = client.Do(req)

	return err
}
