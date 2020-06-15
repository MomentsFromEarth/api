package account

import (
	"fmt"
	"net/http"

	"github.com/MomentsFromEarth/api/internal/models"
	"github.com/MomentsFromEarth/api/internal/services/user"
	"github.com/gin-gonic/gin"
)

// Create is an entrypoint of controller
func Create(c *gin.Context) {
	newUser := &models.NewUser{}
	err := c.Bind(newUser)
	if err != nil {
		badRequest(c, err.Error())
	}
	u := user.Create(newUser)
	c.JSON(http.StatusCreated, u)
	c.Next()
}

// Read is an entrypoint of controller
func Read(c *gin.Context) {
	email, _ := c.Get("email")
	u := user.Get(email.(string))
	c.JSON(http.StatusOK, u)
	c.Next()
}

// Update is an entrypoint of controller
func Update(c *gin.Context) {
}

// Delete is an entrypoint of controller
func Delete(c *gin.Context) {

}

func badRequest(c *gin.Context, message string) {
	fmt.Printf("[AccountError] %v", message)
	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": message})
}
