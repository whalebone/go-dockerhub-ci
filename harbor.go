package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

func harborHandler(c echo.Context) error {

	payload := &harborPayload{}
	err := c.Bind(payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONError(err))
	}

	if payload.EventType != "pushImage" {
		return c.JSON(http.StatusOK, &successResponse{})
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
			return c.JSON(http.StatusBadRequest, JSONError(err))
		}
	}


	response := &successResponse{
		State:       "success",
		Description: "Slack notified",
		Context:     "A OK",
		TargetURL:   "",
	}

	return c.JSON(http.StatusOK, response)
}
