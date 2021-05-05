package controllers

import "github.com/gin-gonic/gin"

type UserControllerI interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

type userController struct{}

var UserController UserControllerI = &userController{}


func (u *userController) Login(c *gin.Context) {

}

func (u *userController) Register(c *gin.Context) {

}