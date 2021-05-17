package api

import (
	"food/model"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getNextWeeklyMenu(c *gin.Context) {
	menu, err := model.GetNextMenu()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, menu)
}

func insertWeeklyMenu(c *gin.Context) {

	var menu model.WeeklyMenu

	err := c.BindJSON(&menu)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = model.InsertWeeklyMenu(menu)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func getAllWeeklyMenu(c *gin.Context) {

	listMenu, err := model.GetAllMenu()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, listMenu)

}

func getDetailWeeklyMenu(c *gin.Context) {
	id := c.Query("id")
	menu, err := model.GetMenu(id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, menu)
}

func activeWeeklyMenu(c *gin.Context) {
	id := c.Query("id")
	err := model.ActiveMenu(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func deactiveWeeklyMenu(c *gin.Context) {
	id := c.Query("id")
	err := model.DeactiveMenu(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}
