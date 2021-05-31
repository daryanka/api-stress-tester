package controllers

import (
	"github.com/daryanka/api-stress-tester/api/domains/request_overviews"
	"github.com/daryanka/api-stress-tester/api/services"
	"github.com/daryanka/api-stress-tester/api/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type RequestOverviewControllerI interface {
	All(c *gin.Context)
	Remove(c *gin.Context)
	Single(c *gin.Context)
	Create(c *gin.Context)
}

type requestOverviewController struct{}

var RequestOverviewController RequestOverviewControllerI = &requestOverviewController{}

func (i *requestOverviewController) All(c *gin.Context) {
	user := GetAuthUser(c)

	res, err := services.RequestOverviewService.All(user.ID)

	if err != nil {
		c.JSON(err.Code(), err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (i *requestOverviewController) Remove(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)

	user := GetAuthUser(c)

	err := services.RequestOverviewService.Delete(idInt, user.ID)

	if err != nil {
		c.JSON(err.Code(), err)
		return
	}

	c.JSON(http.StatusOK, NoError{})
}

func (i *requestOverviewController) Single(c *gin.Context) {
	id := c.Param("id")
	idInt, _ := strconv.ParseInt(id, 10, 64)

	user := GetAuthUser(c)

	res, err := services.RequestOverviewService.GetSingle(idInt, user.ID)

	if err != nil {
		c.JSON(err.Code(), err)
		return
	}

	c.JSON(http.StatusOK, res)
}

func (i *requestOverviewController) Create(c *gin.Context) {
	var reqBody request_overviews.NewRequest

	if ok := utils.GinShouldPassAll(c,
		utils.GinShouldBindJSON(&reqBody),
		utils.GinShouldValidate(&reqBody),
	); !ok {
		return
	}

	reqBody.UserID = GetAuthUser(c).ID

	id, err := services.RequestOverviewService.Create(reqBody)

	if err != nil {
		c.JSON(err.Code(), err)
		return
	}

	c.JSON(http.StatusOK, NoError{
		Message: "Started stress test",
		ID:      id,
	})
}
