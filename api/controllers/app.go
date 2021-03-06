package controllers

import (
	"github.com/daryanka/api-stress-tester/api/clients"
	"github.com/daryanka/api-stress-tester/api/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func StartRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())

	corsDefault := cors.DefaultConfig()
	if os.Getenv("MODE") == "PROD" {
		corsDefault.AllowOrigins = []string{"https://api-tester.daryanamin.co.uk"}
	} else {
		corsDefault.AllowOrigins = []string{"http://localhost:5000"}
	}

	corsDefault.AddAllowHeaders("authorization")
	r.Use(cors.New(corsDefault))

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid route",
			"code":  http.StatusNotFound,
		})
	})

	r.GET("/health", func(c *gin.Context) {
		dbOk := "OK"
		if err := clients.DB.Ping(); err != nil {
			dbOk = "ERROR"
		}
		c.JSON(http.StatusOK, gin.H{
			"status":   "OK",
			"database": dbOk,
		})
	})

	v1 := r.Group("/v1")
	{
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/login", UserController.Login)
			authRoutes.POST("/register", UserController.Register)
			authRoutes.POST("/verify", UserController.VerifyEmail)
			authRoutes.GET("/me", middleware.ValidateAuthToken(), UserController.Me)
		}

		domainRoutes := v1.Group("/domains", middleware.ValidateAuthToken())
		{
			domainRoutes.GET("/all", DomainController.All)
			domainRoutes.POST("/create", DomainController.Create)
			domainRoutes.POST("/confirm", DomainController.Confirm)
			domainRoutes.DELETE("/remove/:id", DomainController.Remove)
		}

		requestRoutes := v1.Group("/requests", middleware.ValidateAuthToken())
		{
			requestRoutes.GET("/all", RequestOverviewController.All)
			requestRoutes.GET("/individual/:id", RequestOverviewController.Single)
			requestRoutes.POST("/create", RequestOverviewController.Create)
			requestRoutes.DELETE("/remove/:id", RequestOverviewController.Remove)
		}

		// Websocket Connection
		v1.GET("/ws", middleware.ValidateAuthToken(), WebSocketController.Connect)
	}

	return r
}
