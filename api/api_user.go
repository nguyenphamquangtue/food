package api

import (
	"encoding/base64"
	"food/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getUserInfo(c *gin.Context) {
	ip := c.Request.Header.Get("ip")

	user, err := model.GetUserByIp(ip)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	if user.Role == model.RoleAdmin {
		ipEnc := base64.StdEncoding.EncodeToString([]byte(ip))
		c.SetCookie("_iz", ipEnc, 24*3600, "/", "localhost", false, false)
	}
	c.JSON(http.StatusOK, user)
}

func getUserByUsername(c *gin.Context) {
	username := c.Query("username")

	user, err := model.GetUserByUsername(username)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, user)
}

func scheduleUpdateUser(c *gin.Context) {
	cronKey := c.Query("cronKey")
	if cronKey != model.CronKey {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"message": "Error!",
		})
		return
	}
	err := model.CronUpdateUser()
	if err != nil {
		c.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
