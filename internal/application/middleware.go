package application

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/dwnGnL/testWork/lib/goerrors"
	"github.com/gin-gonic/gin"
)

const UserIDKey = "user_id"

func CheckBaererToken(app Core) gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerString := c.GetHeader("Authorization")
		if bearerString == "" {
			goerrors.Log().Warnln("header Authorization not found")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		tokenStr, err := parseBearerToken(bearerString)
		if err != nil {
			goerrors.Log().WithError(err).Warnln("parse bearer token err")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		userID, err := app.GetAuth().CheckToken(tokenStr)
		if err != nil {
			goerrors.Log().WithError(err).Warnln("check token err")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set(UserIDKey, userID)
		c.Next()
	}
}

func parseBearerToken(header string) (string, error) {
	splitToken := strings.Split(header, "Bearer ")
	if len(splitToken) != 2 {
		return "", fmt.Errorf("invalid Bearer token format")
	}
	return splitToken[1], nil
}
