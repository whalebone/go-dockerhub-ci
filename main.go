package main

import (
	"context"
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

const DEFAULT_TIMEOUT = 10 * time.Second

func main() {
	viper.AutomaticEnv()
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("PATH_PREFIX", "/")

	prefix := ensureAbsolute(viper.GetString("PATH_PREFIX"))
	fmt.Println("Path prefix: ", prefix)

	e := echo.New()
	e.Logger.SetLevel(log.DEBUG)

	loggerConfig := middleware.DefaultLoggerConfig
	loggerConfig.Skipper = func(c echo.Context) bool {
		return c.Request().Method == http.MethodOptions
	}

	e.Use(middleware.LoggerWithConfig(loggerConfig))
	e.Use(middleware.Recover())

	// e.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	// 	e.Logger.Debug(string(reqBody))
	// }))

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"https://*.dockerhub.com", "https://dockerhub.com", "https://harbor.whalebone.io"},
		AllowMethods: []string{http.MethodPost, http.MethodOptions},
	}))

	e.POST(prefix, dockerhubConfirmer)
	e.POST(appendPath(prefix,"harbor"), harborHandler)

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

	ctx, cancel := context.WithTimeout(context.Background(), DEFAULT_TIMEOUT)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}

func ensureAbsolute(path string) string {
	if path[0] == '/' {
		return path
	}
	return "/" + path
}

func appendPath(path string, second string) string {
	if path[len(path)-1] == '/' {
		return path + second
	}

	return path + "/" + second
}
