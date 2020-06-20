package moment

import (
	"fmt"
	"net/http"

	auth "github.com/MomentsFromEarth/api/internal/auth"
	models "github.com/MomentsFromEarth/api/internal/models/moment"

	"github.com/MomentsFromEarth/api/internal/services/moment"
	"github.com/MomentsFromEarth/api/internal/services/user"
	"github.com/gin-gonic/gin"
)

// Create is an entrypoint of controller
func Create(c *gin.Context) {
	email, _ := c.Get("email")
	u, err := user.Read(email.(string))
	if err != nil {
		errResponse(c, http.StatusNotFound, err.Error())
		return
	}
	newMoment := &models.NewMoment{}
	err = c.Bind(newMoment)
	if err != nil {
		errResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	newMoment.Creator = u.UserID
	moment, err := moment.Create(newMoment)
	if err != nil {
		errResponse(c, http.StatusBadGateway, err.Error())
		return
	}
	c.JSON(http.StatusOK, moment)
	c.Next()
}

// Read is an entrypoint of controller
func Read(c *gin.Context) {
	momentID, _ := c.Params.Get("moment_id")
	m, err := moment.Read(momentID)
	if err != nil {
		errResponse(c, http.StatusBadGateway, err.Error())
		return
	}
	c.JSON(http.StatusOK, m)
	c.Next()
}

// Update is an entrypoint of controller
func Update(c *gin.Context) {
	email, _ := c.Get("email")
	u, err := user.Read(email.(string))
	if err != nil {
		errResponse(c, http.StatusNotFound, err.Error())
		return
	}
	momentID, _ := c.Params.Get("moment_id")
	existing, err := moment.Read(momentID)
	if err != nil {
		errResponse(c, http.StatusBadGateway, err.Error())
		return
	}
	updatedMoment := &models.Moment{}
	err = c.Bind(updatedMoment)
	if err != nil {
		errResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if u.UserID != existing.Creator || updatedMoment.MomentID != existing.MomentID {
		errResponse(c, http.StatusUnauthorized, "Only the creator may update a moment")
		return
	}
	m, err := moment.Update(updatedMoment)
	if err != nil {
		errResponse(c, http.StatusBadGateway, err.Error())
		return
	}
	c.JSON(http.StatusOK, m)
	c.Next()
}

// Delete is an entrypoint of controller
func Delete(c *gin.Context) {
	email, _ := c.Get("email")
	u, err := user.Read(email.(string))
	if err != nil {
		errResponse(c, http.StatusNotFound, err.Error())
		return
	}
	momentID, _ := c.Params.Get("moment_id")
	existing, err := moment.Read(momentID)
	if err != nil {
		errResponse(c, http.StatusBadGateway, err.Error())
		return
	}
	if u.UserID != existing.Creator {
		errResponse(c, http.StatusUnauthorized, "Only the creator may delete a moment")
		return
	}
	m, err := moment.Delete(momentID)
	if err != nil {
		errResponse(c, http.StatusBadGateway, err.Error())
		return
	}
	c.JSON(http.StatusOK, m)
	c.Next()
}

// Callback is an entrypoint of controller
func Callback(c *gin.Context) {
	apiKey := c.Query("api_key")
	authorized := auth.IsValidKey(apiKey)
	if authorized != true {
		errResponse(c, http.StatusUnauthorized, "API Key invalid")
		return
	}
	momentID, _ := c.Params.Get("moment_id")
	existing, err := moment.Read(momentID)
	if err != nil {
		errResponse(c, http.StatusBadGateway, err.Error())
		return
	}
	updatedMoment := &models.Moment{}
	err = c.Bind(updatedMoment)
	if err != nil {
		errResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	if updatedMoment.MomentID != existing.MomentID {
		errResponse(c, http.StatusUnauthorized, "Moment ID does not match")
		return
	}
	m, err := moment.Update(updatedMoment)
	if err != nil {
		errResponse(c, http.StatusBadGateway, err.Error())
		return
	}
	c.JSON(http.StatusOK, m)
	c.Next()
}

func errResponse(c *gin.Context, status int, message string) {
	fmt.Printf("[MomentError] %v", message)
	c.AbortWithStatusJSON(status, gin.H{"message": message})
}
