package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"

	"majestic-gondola/docs"
	"majestic-gondola/internal/handlers"
	"majestic-gondola/internal/repository"
	"majestic-gondola/internal/service"

	"majestic-gondola/bootstrap"
)

// gin-swagger middleware

//	@title			Majestic gondola API
//	@version		1.0
//	@description	API docs for the golang project for learning.

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	// Set up configs
	config := bootstrap.LoadConfig()

	addr := fmt.Sprintf("%s:%d", config.Host, config.Port)
	httpAddr := "http://" + addr

	// Set dynamic swagger info
	docs.SwaggerInfo.Host = addr
	docs.SwaggerInfo.BasePath = "/"

	// Set logger level from config
	var logLevel slog.Level

	if err := logLevel.UnmarshalText([]byte(config.LogLevel)); err != nil {
		logLevel = slog.LevelInfo
	}

	// Set up logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: &logLevel}))
	slog.SetDefault(logger)
	logger.Info("Log level is set", slog.String("Log level", logLevel.String()))

	// Set up db connection
	db := bootstrap.NewDbConnection(config, logger)
	defer db.Close()

	// Set up web app
	r := gin.New()
	r.SetTrustedProxies(nil)
	r.Use(bootstrap.SlogMiddleware(logger))
	r.Use(gin.Recovery())

	// Set up dependencies
	tr := repository.NewTrackRepository(db, logger)
	tsvc := service.NewTrackService(tr, logger)
	th := handlers.NewTrackHandler(tsvc, logger)

	// Set up endpoints
	trackGroup := r.Group("/track")
	{
		trackGroup.GET("/", th.GetTracks)
		trackGroup.GET("/:id", th.GetTrack)
		trackGroup.POST("/", th.CreateTracks)
		trackGroup.PUT("/:id", th.UpdateTrack)
		trackGroup.POST("/populate/:count", th.PopulateTracks)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})

	// Set up graceful shutdown
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Listen error", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	logger.Info("Starting server", slog.String("address", httpAddr), slog.String("swagger", httpAddr+"/swagger"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Forced shutdown", slog.Any("error", err))
	}
}
