package controllers

import (
	"github.com/daryanka/api-stress-tester/api/domains/user"
	"github.com/gin-gonic/gin"
)

type NoError struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Type    string `json:"type,omitempty"`
}

func GetAuthUser(c *gin.Context) user.User {
	u, exist := c.Get("user")
	if !exist {
		return user.User{}
	}
	val, ok := u.(user.User)
	if !ok {
		return user.User{}
	}
	return val
}