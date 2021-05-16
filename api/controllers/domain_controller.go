package controllers

import (
	"github.com/daryanka/api-stress-tester/api/domains/domains"
	"github.com/daryanka/api-stress-tester/api/services"
	"github.com/daryanka/api-stress-tester/api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type DomainControllerI interface {
	All(c *gin.Context)
	Create(c *gin.Context)
	Confirm(c *gin.Context)
	Remove(c *gin.Context)
}

type domainController struct{}

var DomainController DomainControllerI = &domainController{}

func (i *domainController) All(c *gin.Context) {
	user := GetAuthUser(c)

	res, err := services.DomainService.All(user.ID)

	if err != nil {
		c.JSON(err.Code(), err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (i *domainController) Create(c *gin.Context) {
	var reqBody domains.CreateDomain

	if ok := utils.GinShouldPassAll(c,
		utils.GinShouldBindJSON(&reqBody),
		utils.GinShouldValidate(&reqBody),
	); !ok {
		return
	}

	user := GetAuthUser(c)
	reqBody.UserID = user.ID

	id, err := services.DomainService.Create(reqBody)

	if err != nil {
		c.JSON(err.Code(), err)
		return
	}

	c.JSON(http.StatusOK, NoError{
		ID: id,
	})
}

func (i *domainController) Confirm(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)

	err := services.DomainService.Verify(idInt)

	if err != nil {
		c.JSON(err.Code(), err)
		return
	}

	c.JSON(http.StatusOK, NoError{})
}

func (i *domainController) Remove(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)

	user := GetAuthUser(c)

	err := services.DomainService.Delete(user.ID, idInt)

	if err != nil {
		c.JSON(err.Code(), err)
		return
	}

	c.JSON(http.StatusOK, NoError{})
}
