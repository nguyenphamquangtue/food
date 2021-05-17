package model

import "os"

const DB = "food"
const CollectionUser = "users"
const CollectionFood = "foods"
const CollectionWeeklyMenu = "weeklymenu"
const CollectionFoodOrders = "foodorders"

const (
	TokenAccess = "Access token"
)

const CronKey = "QagyKd5zEtLYlRo6"

func configMongoUrl() {
	if os.Getenv("APP_ENV") == "docker" {
		os.Setenv("MONGO_URL", "mongodb://dockerFood:dockerFood@mongo:27017")
	} else {
		os.Setenv("MONGO_URL", "mongodb://localhost:27017")
	}
}
