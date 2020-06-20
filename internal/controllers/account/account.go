package account

import (
	"fmt"
	"net/http"

	models "github.com/MomentsFromEarth/api/internal/models/user"
	"github.com/MomentsFromEarth/api/internal/services/user"
	"github.com/gin-gonic/gin"
)

// Profile is an entrypoint of controller
func Profile(c *gin.Context) {
	username, _ := c.Params.Get("username")
	u, err := user.Profile(username)
	if err != nil {
		errResponse(c, http.StatusBadGateway, err.Error())
	}
	c.JSON(http.StatusOK, u)
	c.Next()
}

// Create is an entrypoint of controller
func Create(c *gin.Context) {
	newUser := &models.NewUser{}
	err := c.Bind(newUser)
	if err != nil {
		errResponse(c, http.StatusBadRequest, err.Error())
	}
	email, _ := c.Get("email")
	newUser.Email = email.(string)
	u, err := user.Create(newUser)
	if err != nil {
		errResponse(c, http.StatusBadGateway, err.Error())
	}
	c.JSON(http.StatusCreated, u)
	c.Next()
}

// Read is an entrypoint of controller
func Read(c *gin.Context) {
	email, _ := c.Get("email")
	u, err := user.Read(email.(string))
	if err != nil {
		errResponse(c, http.StatusBadGateway, err.Error())
	}
	c.JSON(http.StatusOK, u)
	c.Next()
}

// Update is an entrypoint of controller
func Update(c *gin.Context) {

	// todo: fetch existing user, update that one instead

	updatedUser := &models.User{}
	err := c.Bind(updatedUser)
	if err != nil {
		errResponse(c, http.StatusBadRequest, err.Error())
	}
	email, _ := c.Get("email")
	updatedUser.Email = email.(string)
	u, err := user.Update(updatedUser)
	if err != nil {
		errResponse(c, http.StatusBadGateway, err.Error())
	}
	c.JSON(http.StatusOK, u)
	c.Next()
}

// Delete is an entrypoint of controller
func Delete(c *gin.Context) {
	email, _ := c.Get("email")
	u, err := user.Delete(email.(string))
	if err != nil {
		errResponse(c, http.StatusBadGateway, err.Error())
	}
	c.JSON(http.StatusOK, u)
	c.Next()
}

func errResponse(c *gin.Context, status int, message string) {
	fmt.Printf("[AccountError] %v", message)
	c.AbortWithStatusJSON(status, gin.H{"message": message})
}
