package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

func dockerhubConfirmer(c echo.Context) error {

	payload := &dockerhubPayload{}
	err := c.Bind(payload)
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONError(err))
	}

	pushedAtTime := time.Unix(int64(payload.PushData.PushedAt), 0)

	slackData := fmt.Sprintf("*%s* pushed *%s:%s* at %s", payload.PushData.Pusher, payload.Repository.Name, payload.PushData.Tag, pushedAtTime)
	c.Logger().Debug(string(slackData))
	err = sendSlackNotification(string(slackData))
	if err != nil {
		return c.JSON(http.StatusBadRequest, JSONError(err))
	}

	response := &successResponse{
		State:       "success",
		Description: "Slack notified",
		Context:     "A OK",
		TargetURL:   "",
	}

	confirmDockerhub(payload.CallbackURL, response)

	return c.JSON(http.StatusOK, response)
}

func confirmDockerhub(url string, payload *successResponse) error {
	responseBody, _ := json.Marshal(payload)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(responseBody))

	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	_, err = client.Do(req)

	return err
}

func sendSlackNotification(msg string) error {
	slackBody, _ := json.Marshal(slackRequestBody{Text: msg})
	req, err := http.NewRequest(http.MethodPost, viper.GetString("WEBHOOK"), bytes.NewBuffer(slackBody))
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	if buf.String() != "ok" {
		return errors.New("Non-ok response returned from Slack")
	}
	return nil
}

func main() {
	viper.AutomaticEnv()
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("PATH_PREFIX", "/")

	prefix := viper.GetString("PATH_PREFIX")
	fmt.Println("Path prefix: ", prefix)

	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	loggerConfig := middleware.DefaultLoggerConfig
	loggerConfig.Skipper = func(c echo.Context) bool {
		if c.Request().Method == http.MethodOptions {
			return true
		}
		return false
	}

	e.Use(middleware.LoggerWithConfig(loggerConfig))
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://*.dockerhub.com", "https://dockerhub.com"},
		AllowMethods: []string{http.MethodPost, http.MethodOptions},
	}))

	e.POST(prefix, dockerhubConfirmer)

	// Start server
	go func() {
		if err := e.Start(":" + viper.GetString("PORT")); err != nil {
			e.Logger.Info("shutting down the server")
			e.Logger.Error(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}

}
