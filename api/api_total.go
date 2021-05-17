package api

import (
	"food/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func TotalDaily(c *gin.Context) {
	id := c.Query("id")
	order, err := model.TotalDaily(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, order)
}

func TotalWeek(c *gin.Context) {
	id := c.Query("id")
	order, err := model.TotalWeek(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, order)
}
