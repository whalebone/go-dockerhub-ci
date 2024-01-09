package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/whalebone/go-dockerhub-ci/model"
)

func HarborHandler(c echo.Context) error {
	payload := &model.HarborPayload{}
	err := c.Bind(payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, model.JSONError(err))
	}

	if payload.EventType != "pushImage" && payload.EventType != "PUSH_ARTIFACT" {
		return c.JSON(http.StatusOK, &model.SuccessResponse{})
	}

	pushedAtTime := time.Unix(int64(payload.Time), 0)

	for _, resource := range payload.EventData.Resources {
		slackData := fmt.Sprintf(
			"*%s* pushed *%s:%s* to Harbor at <!date^%d^{date_num} {time_secs}|%s>\n> *%s*",
			payload.User,
			payload.EventData.Repository.FullRepoName,
			resource.Tag,
			int64(payload.Time),
			pushedAtTime,
			resource.Resource,
		)
		err = sendSlackNotification(slackData)
		if err != nil {
			return c.JSON(http.StatusBadRequest, model.JSONError(err))
		}
	}

	response := &model.SuccessResponse{
		State:       "success",
		Description: "Slack notified",
		Context:     "A OK",
		TargetURL:   "",
	}

	return c.JSON(http.StatusOK, response)
}
