package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func setGroundRules(e *echo.Echo, lvl log.Lvl) *echo.Echo {
	e.HideBanner = true
	e.HidePort = true

	DefaultLoggerConfig := middleware.LoggerConfig{
		Skipper: middleware.DefaultSkipper,
		Format: `{"time":"${time_rfc3339_nano}","level":"INFO","remote_ip":"${remote_ip}",` +
			`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
			`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
			`,"bytes_in":${bytes_in},"bytes_out":${bytes_out}}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05.00000",
	}

	// set level
	e.Logger.SetLevel(lvl)

	// e.Use(middleware.Logger())
	e.Use(middleware.LoggerWithConfig(DefaultLoggerConfig))
	e.Use(middleware.Recover())

	return e
}

// Serve the main function running the echo server
func Serve(port string, debug bool) {
	e := echo.New()

	var logLvl log.Lvl = log.INFO
	if debug {
		logLvl = log.DEBUG
	}

	e = setGroundRules(e, logLvl)

	e.Logger.Infof("starting server on port %s", port)

	e.GET("/", echo.WrapHandler(IndexFileServer()))
	e.GET("/version", version)
	e.PUT("/paste", pastePutHandler)
	e.POST("/paste", pastePostHandler)
	e.GET("/:hashVal", fetchHandler)

	// Start server
	go func() {
		if err := e.Start(":" + port); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	e.Logger.Debugf("got signal to quit %v", quit)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	e.Logger.Info("closing down server")
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
