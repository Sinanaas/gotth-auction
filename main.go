package main

import (
	"log"

	"github.com/Sinanaas/gotth-auction/controllers"
	"github.com/Sinanaas/gotth-auction/initializers"
	"github.com/Sinanaas/gotth-auction/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var (
	server              *gin.Engine
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouterController
	BasicRouteController routes.BasicRouterController
	BasicController	 controllers.BasicController
)

func init() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)
	AuthController = controllers.NewAuthController(initializers.DB)
	AuthRouteController = routes.NewAuthRouterController(AuthController)
	BasicController = controllers.NewBasicController(initializers.DB)
	BasicRouteController = routes.NewBasicRouterController(BasicController)

	server = gin.Default()
	store := cookie.NewStore([]byte(config.SessionSecretKey))
	server.Use(sessions.Sessions("mysession", store))
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("? Could not load environment variables", err)
	}

	initializers.ConnectDB(&config)

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", config.ClientOrigin}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/")
	// auth router
	AuthRouteController.AuthRoute(router)
	// basic router
	BasicRouteController.BasicRoute(router)
	log.Fatal(server.Run(":" + config.ServerPort))
}
