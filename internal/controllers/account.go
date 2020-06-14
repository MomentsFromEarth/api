package account

import (
	"net/http"

	"github.com/MomentsFromEarth/api/internal/services/user"
	"github.com/gin-gonic/gin"
)

// Create is an entrypoint of controller
func Create(c *gin.Context) {
}

// Read is an entrypoint of controller
func Read(c *gin.Context) {
	email, _ := c.Get("email")
	u := user.FromEmail(email.(string))
	c.JSON(http.StatusOK, u)
	c.Next()
}

// Update is an entrypoint of controller
func Update(c *gin.Context) {
}

// Delete is an entrypoint of controller
func Delete(c *gin.Context) {

}
