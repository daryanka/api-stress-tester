package controllers

import (
	"github.com/daryanka/api-stress-tester/api/domains/users"
	"github.com/daryanka/api-stress-tester/api/utils"
	"github.com/gin-gonic/gin"
)

type UserControllerI interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
}

type userController struct{}

var UserController UserControllerI = &userController{}


func (u *userController) Login(c *gin.Context) {
	var reqBody users.ReqLogin

	if ok := utils.GinShouldPassAll(c,
		utils.GinShouldBindJSON(&reqBody),
		utils.GinShouldValidate(&reqBody),
	); !ok {
		return
	}
}

func (u *userController) Register(c *gin.Context) {

}