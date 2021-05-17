package api

import (
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Route struct {
	method  string
	path    string
	ip      gin.HandlerFunc
	handler gin.HandlerFunc
}

type AuthenRoute struct {
	method  string
	path    string
	ip      gin.HandlerFunc
	authen  gin.HandlerFunc
	handler gin.HandlerFunc
}

var listRoute = []Route{
	// user
	{http.MethodGet, "/api/me", ip, getUserInfo},

	// menu
	{http.MethodGet, "/api/weekly-menu/next", ip, getNextWeeklyMenu},

	// food-orders
	{http.MethodPost, "/api/food-orders/new", ip, insertFoodOrders},
	{http.MethodGet, "/api/food-orders/current", ip, getFoodOrdersThisWeek},
	{http.MethodGet, "/api/food-orders/next", ip, getFoodOrdersNextWeek},
	{http.MethodPut, "/api/food-orders/checkout", ip, checkoutFoodOrders},

	{http.MethodPut, "/api/user", ip, scheduleUpdateUser},
	{http.MethodGet, "/ws/food", ip, WsFood},
}

var listAuthenRoute = []AuthenRoute{
	{http.MethodGet, "api/user/detail", ip, authen, getUserByUsername},

	{http.MethodPost, "/api/weekly-menu/new", ip, authen, insertWeeklyMenu},
	{http.MethodGet, "/api/weekly-menu/all", ip, authen, getAllWeeklyMenu},
	{http.MethodGet, "/api/weekly-menu/detail", ip, authen, getDetailWeeklyMenu},
	{http.MethodPut, "/api/weekly-menu/active", ip, authen, activeWeeklyMenu},
	{http.MethodPut, "/api/weekly-menu/deactive", ip, authen, deactiveWeeklyMenu},
	{http.MethodGet, "/api/weekly-menu/not-eat", ip, authen, GetNonFood},
	{http.MethodGet, "/api/weekly-menu/total-daily", ip, authen, TotalDaily},
	{http.MethodGet, "/api/weekly-menu/total-week", ip, authen, TotalWeek},

	{http.MethodGet, "/api/food-orders/detail", ip, authen, GetFoodOrderDetail},
	{http.MethodDelete, "/api/food-orders/delete", ip, authen, deleteFoodOrder},

	{http.MethodPost, "/api/food/new", ip, authen, insertFood},
	{http.MethodGet, "/api/food/all", ip, authen, getAllFood},
	{http.MethodPut, "/api/food/update", ip, authen, updateFood},
	{http.MethodDelete, "/api/food/delete", ip, authen, deleteFood},

	{http.MethodPost, "/api/upload", ip, authen, upload},
}

func NewRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
	}))

	for _, r := range listRoute {
		router.Handle(r.method, r.path, r.ip, r.handler)
	}

	for _, r := range listAuthenRoute {
		router.Handle(r.method, r.path, r.ip, r.authen, r.handler)
	}
	router.StaticFS("/file", http.Dir("public"))

	return router
}
