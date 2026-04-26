package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "majestic-gondola/docs"
	"majestic-gondola/internal/handlers"
	"majestic-gondola/internal/repository"
	"majestic-gondola/internal/service"

	"majestic-gondola/bootstrap"
)

// gin-swagger middleware

//	@title			Majestic gondola API
//	@version		1.0
//	@description	API docs for the golang project for learning.

//	@host		localhost:8080
//	@BasePath	/

// @externalDocs.description	OpenAPI
// @externalDocs.url			https://swagger.io/resources/open-api/
func main() {
	// Set up configs
	config := bootstrap.LoadConfig()

	// Set logger level from config
	var logLevel slog.Level

	if err := logLevel.UnmarshalText([]byte(config.LogLevel)); err != nil {
		logLevel = slog.LevelInfo
	}

	// Set up logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: &logLevel}))
	slog.SetDefault(logger)
	logger.Info("Log level is set", slog.Any("Log level", logLevel))

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

	trackGroup := r.Group("/track")
	{
		trackGroup.GET("/", th.GetTracks)
		trackGroup.GET("/:id", th.GetTrack)
		trackGroup.POST("/", th.CreateTracks)
		trackGroup.PUT("/:id", th.UpdateTrack)
		trackGroup.POST("/populate/:count", th.PopulateTracks)
	}

	r.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// TODO: Graceful shutdown
	if config.Address != "" {
		r.Run(config.Address)
	} else {
		r.Run()
	}
}
