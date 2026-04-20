package main

import (
	"log/slog"
	"net/http"

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

//	@title			Swagger Example API
//	@version		1.0
//	@description	This is a sample server celler server.
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/

//	@externalDocs.description	OpenAPI
//	@externalDocs.url			https://swagger.io/resources/open-api/

// main godoc
//
//	@Summary		Ping service
//	@Description	get pong
//	@Produce		json
//	@Router			/ping [get]
func main() {
	// Set up logger
	logger := slog.New(slog.Default().Handler())
	slog.SetDefault(logger)

	// Set up configs
	config := bootstrap.LoadConfig(logger)

	// Set up db connection
	db := bootstrap.NewDbConnection(config, logger)
	defer db.Close()

	// Set up web app
	r := gin.New()
	r.Use(bootstrap.SlogMiddleware(logger))
	r.Use(gin.Recovery())

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Pong",
		})
	})

	// Set up dependencies
	tr := repository.NewTrackRepository(db, logger)
	tsvc := service.NewTrackService(tr, logger)
	th := handlers.NewTrackHandler(tsvc, logger)

	trackGroup := r.Group("/track")
	{
		trackGroup.GET("/", th.GetTracks)
		trackGroup.GET("/:id", th.GetTrack)
		trackGroup.POST("/", th.CreateTracks)
		trackGroup.PUT("/", th.UpdateTrack)
		trackGroup.POST("/populate/:count", th.PopulateTracks)
	}

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// TODO: Graceful shutdown
	if config.Address != "" {
		r.Run(config.Address)
	} else {
		r.Run()
	}
}
