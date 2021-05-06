package controllers

import "github.com/gin-gonic/gin"

func StartRouter() *gin.Engine {
	r := gin.New()

	v1 := r.Group("/v1")


	authRoutes := v1.Group("/auth")
	{
		authRoutes.POST("/login", UserController.Login)
		authRoutes.POST("/register", UserController.Register)
	}

	return r
}
