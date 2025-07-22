package pkg

import "C"
import (
	"backend/pkg/service"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func UserIdentity(c *gin.Context, userService service.UserService) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "empty auth header"})
		c.Abort()
		return
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "invalid auth header"})
		c.Abort()
		return
	}

	if len(headerParts[1]) == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": "empty token"})
		c.Abort()
		return
	}

	userId, err := userService.ParseToken(headerParts[1])
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"Error": err.Error()})
		c.Abort()
		return
	}

	c.Set(userCtx, userId)

}

func getUserId(c *gin.Context) (string, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		return "", errors.New("user id not found")
	}

	if !ok {
		return "", errors.New("user id is of invalid type")
	}
	strId := id.(string)
	return strId, nil
}
