package tag

import (
	"fmt"
	"net/http"

	models "github.com/MomentsFromEarth/api/internal/models/tag"
	"github.com/MomentsFromEarth/api/internal/services/tag"
	"github.com/gin-gonic/gin"
)

// Create is an entrypoint of controller
func Create(c *gin.Context) {
	newTag := &models.NewTag{}
	err := c.Bind(newTag)
	if err != nil {
		errResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	tag, err := tag.Create(newTag)
	if err != nil {
		errResponse(c, http.StatusBadGateway, err.Error())
		return
	}
	c.JSON(http.StatusOK, tag)
	c.Next()
}

// Read is an entrypoint of controller
func Read(c *gin.Context) {
	tagID, _ := c.Params.Get("tag_id")
	t, err := tag.Read(tagID)
	if err != nil {
		errResponse(c, http.StatusBadGateway, err.Error())
		return
	}
	c.JSON(http.StatusOK, t)
	c.Next()
}

// Delete is an entrypoint of controller
func Delete(c *gin.Context) {
	tagID, _ := c.Params.Get("tag_id")
	t, err := tag.Delete(tagID)
	if err != nil {
		errResponse(c, http.StatusBadGateway, err.Error())
		return
	}
	c.JSON(http.StatusOK, t)
	c.Next()
}

func errResponse(c *gin.Context, status int, message string) {
	fmt.Printf("[TagError] %v", message)
	c.AbortWithStatusJSON(status, gin.H{"message": message})
}
