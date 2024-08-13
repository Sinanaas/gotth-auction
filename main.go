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
	server               *gin.Engine
	AuthController       controllers.AuthController
	AuthRouteController  routes.AuthRouterController
	BasicRouteController routes.BasicRouterController
	BasicController      controllers.BasicController
	WebsocketController  controllers.WebsocketController
	WebsocketRouteController routes.WebsocketRouterController
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
	WebsocketController = controllers.NewWebsocketController(initializers.DB)
	WebsocketRouteController = routes.NewWebsocketRouterController(WebsocketController)

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

	// cors config
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:8000", config.ClientOrigin}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowWebSockets = true
	corsConfig.AllowCredentials = true

	// serve static file
	server.Static("/static", "./static")
	server.Static("/uploads", "./uploads")

	server.Use(cors.New(corsConfig))

	router := server.Group("/")

	// auth router
	AuthRouteController.AuthRoute(router)
	// basic router
	BasicRouteController.BasicRoute(router)
	// websocket router
	WebsocketRouteController.WebsocketRoute(router)

	log.Fatal(server.Run(":" + config.ServerPort))
}
