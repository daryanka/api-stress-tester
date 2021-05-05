package middleware

import (
	"fmt"
	"github.com/daryanka/api-stress-tester/api/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

const bearer = "bearer"

var authorized = utils.NewUnAuthorized("Unauthorized", "EXPIRED_TOKEN")

func ValidateAuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, authorized)
			c.Abort()
			return
		}

		if len(authHeader) < len(bearer) {
			c.JSON(http.StatusUnauthorized, authorized)
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader[len(bearer):], " ", "", -1)

		// Check if the token has expired
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil {
			c.JSON(http.StatusUnauthorized, authorized)
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// TODO get user
			var user users.User

			err := clients.Client.Table("users").
				Select("id", "email", "role", "name").
				Where("id", "=", claims["id"]).
				First(&user)

			if err != nil {
				c.JSON(http.StatusUnauthorized, authorized)
				c.Abort()
				return
			}

			c.Set("user", user)
			c.Next()
			return
		} else {
			c.JSON(http.StatusUnauthorized, authorized)
			c.Abort()
			return
		}
	}
}
