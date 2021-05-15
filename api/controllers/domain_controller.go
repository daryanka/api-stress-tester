package controllers

import "github.com/gin-gonic/gin"

type DomainControllerI interface {

}

type domainController struct{}

var DomainController DomainControllerI = &domainController{}

func (i *domainController) All(c *gin.Context) {

}

func (i *domainController) New(c *gin.Context) {

}

func (i *domainController) Confirm(c *gin.Context) {

}

func (i *domainController) Remove(c *gin.Context) {

}
