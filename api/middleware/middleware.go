package middleware

import (
	"fmt"
	"github.com/daryanka/api-stress-tester/api/domains/user"
	"github.com/daryanka/api-stress-tester/api/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"os"
	"strings"
)

const bearer = "bearer"

var unauthorized = utils.NewUnAuthorized("Unauthorized", "EXPIRED_TOKEN")

func ValidateAuthToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizationHeader := c.GetHeader("authorization")
		websocketHeader := c.GetHeader("Sec-WebSocket-Protocol")

		var authHeader string

		if  authorizationHeader != "" {
			authHeader = authorizationHeader
		} else if websocketHeader != "" {
			authHeader = "Bearer " + websocketHeader
		} else {
			authHeader = "Bearer " + c.Query("auth")
		}

		if authHeader == "" {
			c.JSON(unauthorized.Code(), unauthorized)
			c.Abort()
			return
		}

		if len(authHeader) < len(bearer) {
			c.JSON(unauthorized.Code(), unauthorized)
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
			c.JSON(unauthorized.Code(), unauthorized)
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			id, ok := claims["id"].(float64)
			if !ok {
				utils.Logger.Error("id from JWT not able to type assert as int64 ", claims["id"])
			}

			u, err := user.UserDao.Find(int64(id))
			if err != nil {
				c.JSON(unauthorized.Code(), unauthorized)
				return
			}

			c.Set("user", u)
			c.Next()
			return
		} else {
			c.JSON(unauthorized.Code(), unauthorized)
			c.Abort()
			return
		}
	}
}
