package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/whalebone/go-dockerhub-ci/model"
)

func DockerhubConfirmer(c echo.Context) error {
	payload := &model.DockerhubPayload{}
	err := c.Bind(payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.JSONError(err))
	}

	pushedAtTime := time.Unix(int64(payload.PushData.PushedAt), 0)

	slackData := fmt.Sprintf(
		"*%s* pushed *%s:%s* to Dockerhub at <!date^%d^{date_num} {time_secs}|%s>\n> *%s:%s*",
		payload.PushData.Pusher,
		payload.Repository.Name,
		payload.PushData.Tag,
		int64(payload.PushData.PushedAt),
		pushedAtTime,
		payload.Repository.Name,
		payload.PushData.Tag,
	)

	c.Logger().Debug(slackData)
	err = sendSlackNotification(slackData)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.JSONError(err))
	}

	response := &model.SuccessResponse{
		State:       "success",
		Description: "Slack notified",
		Context:     "A OK",
		TargetURL:   "",
	}

	err = confirmDockerhub(payload.CallbackURL, response)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.JSONError(err))
	}

	return c.JSON(http.StatusOK, response)
}

func confirmDockerhub(url string, payload *model.SuccessResponse) error {
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

	client := &http.Client{Timeout: model.DefaultTimeout}
	_, err = client.Do(req)

	return err
}
