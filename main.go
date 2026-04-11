package main

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "majestic-gondola/docs"

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
	logger := slog.New(slog.Default().Handler())
	slog.SetDefault(logger)

	config := bootstrap.LoadConfig(logger)
	db := bootstrap.GetDbConnection(config, logger)
	defer db.Close()

	// r := gin.New()
	// r.Use(bootstrap.SlogMiddleware(logger))
	// r.Use(gin.Recovery())

	r := gin.Default()

	r.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "Pong",
		})
	})

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
