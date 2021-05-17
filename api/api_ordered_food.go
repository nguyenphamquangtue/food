package api

import (
	"food/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func getFoodOrdersThisWeek(c *gin.Context) {
	ip := c.Request.Header.Get("ip")
	order, err := model.GetFoodOrdersThisWeek(ip)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := model.GenerateToken(order.User.Ip, order.End)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"food_orders": order,
		"token":       token,
	})
}

func getFoodOrdersNextWeek(c *gin.Context) {
	ip := c.Request.Header.Get("ip")
	order, err := model.GetFoodOrdersNextWeek(ip)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := model.GenerateToken(order.User.Ip, order.End)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"food_orders": order,
		"token":       token,
	})
}

func insertFoodOrders(c *gin.Context) {
	var order model.FoodOrders
	err := c.BindJSON(&order)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = model.InsertFoodOrders(order)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := model.GenerateToken(order.User.Ip, order.End)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}

func checkoutFoodOrders(c *gin.Context) {
	requestIp := c.Request.Header.Get("ip")
	token := c.Query("token")
	claims, err := model.VerifyToken(token)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		if listWsConnection[requestIp] != nil {
			_ = listWsConnection[requestIp].WriteMessage(websocket.TextMessage, []byte(err.Error()))
		}
		return
	}

	ip := claims.Ip
	food, user, isEated, err := model.CheckoutFoodOrders(ip)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		if listWsConnection[requestIp] != nil {
			_ = listWsConnection[requestIp].WriteMessage(websocket.TextMessage, []byte(err.Error()))
		}
		return
	}

	if listWsConnection[requestIp] != nil {
		_ = listWsConnection[requestIp].WriteJSON(gin.H{
			"id":       food.Id,
			"name":     food.Name,
			"img":      food.Img,
			"is_eated": isEated,
			"user":     user,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       food.Id,
		"name":     food.Name,
		"img":      food.Img,
		"is_eated": isEated,
		"user":     user,
	})
}

func GetNonFood(c *gin.Context) {
	id := c.Query("id")
	foodOrder, err := model.GetNonFood(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, foodOrder)
}

func GetFoodOrderDetail(c *gin.Context) {
	menuID := c.Query("menuId")
	username := c.Query("username")
	foodOrder, err := model.GetFoodOrderDetail(username, menuID)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, foodOrder)
}

func deleteFoodOrder(c *gin.Context) {
	id := c.Query("id")
	err := model.DeleteFoodOrder(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "delete success",
	})
}
