package context

import (
	"github.com/akinoccc/hysaif/api/models"

	"github.com/gin-gonic/gin"
)

func GetCurrentUser(c *gin.Context) *models.User {
	userInterface, ok := c.Get("user")
	if !ok {
		return nil
	}
	user, ok := userInterface.(models.User)
	if !ok {
		return nil
	}
	return &user
}
