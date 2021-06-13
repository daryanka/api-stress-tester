package controllers

import (
	"github.com/daryanka/api-stress-tester/api/domains/user"
	"github.com/daryanka/api-stress-tester/api/services"
	"github.com/daryanka/api-stress-tester/api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserControllerI interface {
	Login(c *gin.Context)
	Register(c *gin.Context)
	VerifyEmail(c *gin.Context)
	Me(c *gin.Context)
}

type userController struct{}

var UserController UserControllerI = &userController{}

func (u *userController) Login(c *gin.Context) {
	var reqBody user.ReqLogin

	if ok := utils.GinShouldPassAll(c,
		utils.GinShouldBindJSON(&reqBody),
		utils.GinShouldValidate(&reqBody),
	); !ok {
		return
	}

	token, e := services.UserService.Login(&reqBody)
	if e != nil {
		c.JSON(e.Code(), e)
		return
	}

	c.JSON(http.StatusOK, token)
}

func (u *userController) Register(c *gin.Context) {
	var reqBody user.ReqRegister

	if ok := utils.GinShouldPassAll(c,
		utils.GinShouldBindJSON(&reqBody),
		utils.GinShouldValidate(&reqBody),
	); !ok {
		return
	}

	if reqBody.PasswordConfirmation != reqBody.Password {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"password": "Passwords do not match",
		})
		return
	}

	e := services.UserService.Register(&reqBody)

	if e != nil {
		c.JSON(e.Code(), e)
		return
	}

	c.JSON(http.StatusOK, NoError{
		Error:   false,
		Message: "Successfully created account check inbox for confirmation email",
	})
}

func (u *userController) VerifyEmail(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		e := utils.NewBadRequest("Unable to verify email")
		c.JSON(e.Code(), e)
		return
	}

	valid, err := user.UserDao.VerifyEmail(token)
	if err != nil || !valid {
		e := utils.StandardInternalServerError()
		c.JSON(e.Code(), e)
		return
	}

	c.JSON(http.StatusOK, NoError{
		Error:   false,
		Message: "Successfully activated account",
	})
}

func (u *userController) Me(c *gin.Context) {
	c.JSON(http.StatusOK, GetAuthUser(c))
}
