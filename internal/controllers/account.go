package account

import (
	"github.com/gin-gonic/gin"
)

// Create is an entrypoint of controller
func Create(c *gin.Context) {
}

// Read is an entrypoint of controller
func Read(c *gin.Context) {
	email, _ := c.Get("email")
	c.JSON(200, gin.H{
		"email": email,
	})
	c.Next()
}

// Update is an entrypoint of controller
func Update(c *gin.Context) {
}

// Delete is an entrypoint of controller
func Delete(c *gin.Context) {

}
