package controllers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartRouter() *gin.Engine {
	r := gin.New()

	corsDefault := cors.DefaultConfig()
	corsDefault.AllowOrigins = []string{"http://localhost:5000"}
	r.Use(cors.New(corsDefault))

	v1 := r.Group("/v1")
	{
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/login", UserController.Login)
			authRoutes.POST("/register", UserController.Register)
			authRoutes.POST("/verify", UserController.VerifyEmail)
		}

		v1.GET("/ws", WebSocketController.Connect)
	}


	return r
}
