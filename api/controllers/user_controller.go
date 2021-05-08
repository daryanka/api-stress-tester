package controllers

import (
	"database/sql"
	"github.com/daryanka/api-stress-tester/api/clients"
	"github.com/daryanka/api-stress-tester/api/domains/users"
	"github.com/daryanka/api-stress-tester/api/utils"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
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

	user, err := users.UserDao.FindByEmail(reqBody.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			e := utils.NewUnAuthorized("Invalid username/password combination")
			c.JSON(e.Code(), e)
			return
		}
		e := utils.StandardInternalServerError()
		c.JSON(e.Code(), e)
		return
	}

	// Check Password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(reqBody.Password))
	if err != nil {
		e := utils.NewUnAuthorized("Invalid username/password combination")
		c.JSON(e.Code(), e)
		return
	}

	token, err := utils.CreateAuthToken(user.ID)
	if err != nil {
		e := utils.StandardInternalServerError()
		c.JSON(e.Code(), e)
		return
	}

	c.JSON(http.StatusOK, token)
}

func (u *userController) Register(c *gin.Context) {
	var reqBody users.ReqRegister

	if ok := utils.GinShouldPassAll(c,
		utils.GinShouldBindJSON(&reqBody),
		utils.GinShouldValidate(&reqBody),
	); !ok {
		return
	}

	if reqBody.PasswordConfirmation != reqBody.Password {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"password": "The passwords need to match",
		})
		return
	}

	inUse, err := users.UserDao.EmailInUse(reqBody.Email)

	if err != nil {
		e := utils.StandardInternalServerError()
		c.JSON(e.Code(), e)
		return
	}

	if inUse {
		e := utils.NewBadRequest("This email is already associated to an account", "EMAIL_IN_USE")
		c.JSON(e.Code(), e)
		return
	}

	// Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), 11)
	if err != nil {
		e := utils.StandardInternalServerError()
		c.JSON(e.Code(), e)
		return
	}

	// Generate Email Token
	EmailToken := utils.RandStringRunes(64)

	_, err = users.UserDao.Create(&users.User{
		Name:       reqBody.Name,
		Email:      reqBody.Email,
		EmailToken: EmailToken,
		Password:   string(hashedPassword),
	})

	if err != nil {
		e := utils.StandardInternalServerError()
		c.JSON(e.Code(), e)
		return
	}

	go clients.SendMail("register.html", reqBody.Email, "Welcome Email", struct {
		Token string
	}{EmailToken})

	c.JSON(http.StatusOK, NoError{
		Error:   false,
		Message: "Successfully created account check inbox for confirmation email",
	})
}


func VerifyEmail(c *gin.Context) {
	token := c.Query("token")

	valid, err := users.UserDao.VerifyEmail(token)
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