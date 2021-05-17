package model

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/now"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type StatusEat int

const (
	StatusEated    StatusEat = 1
	StatusNotEated StatusEat = 0
)

// type DailyFood struct {
// 	Food   Food      `json:"food,omitempty"`
// 	Status StatusEat `json:"status"`
// }

type DailyFoodOrders struct {
	Breakfast Food `json:"breakfast,omitempty"`
	Lunch     Food `json:"lunch,omitempty"`
}

type FoodOrders struct {
	Id                primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	User              User               `json:"user,omitempty" validate:"required"`
	Mon               DailyFoodOrders    `json:"mon,omitempty"`
	Tue               DailyFoodOrders    `json:"tue,omitempty"`
	Wed               DailyFoodOrders    `json:"wed,omitempty"`
	Thu               DailyFoodOrders    `json:"thu,omitempty"`
	Fri               DailyFoodOrders    `json:"fri,omitempty"`
	Start             int64              `json:"start,omitempty" validate:"required"`
	End               int64              `json:"end,omitempty" validate:"required"`
	StatusNonBreakast []int              `json:"status_non_breakfast,omitempty" bson:"status_non_breakfast"`
	StatusNonLunch    []int              `json:"status_non_lunch,omitempty" bson:"status_non_lunch"`
}

func InsertFoodOrders(order FoodOrders) error {
	err := validator.New().Struct(order)
	if err != nil {
		return err
	}

	collection := db.client.Database(DB).Collection(CollectionFoodOrders)

	filter := bson.M{
		"start":   order.Start,
		"end":     order.End,
		"user.ip": order.User.Ip,
	}

	checkExist := collection.FindOne(db.context, filter)
	if checkExist.Err() == nil {
		return errors.New("Food Order existed!")
	}
	_, err = collection.InsertOne(db.context, order)
	if err != nil {
		return err
	}
	return nil
}

func GetFoodOrdersThisWeek(ip string) (*FoodOrders, error) {

	startWeekTime := now.BeginningOfWeek().Unix() * 1000
	endWeekTime := now.EndOfWeek().Unix() * 1000

	filter := bson.M{
		"user.ip": ip,
		"start":   bson.M{"$gte": startWeekTime},
		"end":     bson.M{"$lte": endWeekTime},
	}

	var foodOrders FoodOrders
	collection := db.client.Database(DB).Collection(CollectionFoodOrders)
	result := collection.FindOne(db.context, filter)

	err := result.Decode(&foodOrders)
	if err != nil {
		return nil, err
	}
	return &foodOrders, nil
}

func GetFoodOrdersNextWeek(ip string) (*FoodOrders, error) {
	menu, err := GetNextMenu()
	if err != nil {
		return nil, err
	}

	start := menu.Start
	end := menu.End
	filter := bson.M{
		"user.ip": ip,
		"start":   start,
		"end":     end,
	}

	var foodOrder FoodOrders
	collection := db.client.Database(DB).Collection(CollectionFoodOrders)
	result := collection.FindOne(db.context, filter)
	err = result.Decode(&foodOrder)
	if err != nil {
		return nil, err
	}

	return &foodOrder, nil
}

func GetAllFoodOrders() (*[]FoodOrders, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.M{"start": -1})
	collection := db.client.Database(DB).Collection(CollectionFoodOrders)

	cursor, err := collection.Find(db.context, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}

	var listFoodOrders []FoodOrders
	err = cursor.All(db.context, &listFoodOrders)
	if err != nil {
		return nil, err
	}

	return &listFoodOrders, nil
}

/** this function use to verify whether user use food or not */
func CheckoutFoodOrders(ip string) (*Food, *User, bool, error) {

	// get today
	_now := time.Now()
	day := mapWeek[_now.Weekday()]

	// check valid date and session
	if day == "" {
		return nil, nil, false, errors.New("Invalid date time")
	}

	// get session : breakfast or lunch
	session, err := getCurrentSession()
	if err != nil {
		return nil, nil, false, err
	}

	startWeekTime := now.BeginningOfWeek().Unix() * 1000
	endWeekTime := now.EndOfWeek().Unix() * 1000
	filter := bson.M{
		"user.ip": ip,
		"start":   bson.M{"$gte": startWeekTime},
		"end":     bson.M{"$lte": endWeekTime},
	}

	// get field need update in database
	var updateField string
	if session == "breakfast" {
		updateField = "status_non_breakfast"
	} else {
		updateField = "status_non_lunch"
	}

	// get remove value in list status then remove it (database)
	removeValue := mapWeekValue[day]
	updateValue := bson.M{"$pull": bson.M{updateField: removeValue}}

	// update in database
	collection := db.client.Database(DB).Collection(CollectionFoodOrders)
	result := collection.FindOneAndUpdate(db.context, filter, updateValue)

	if result.Err() != nil {
		return nil, nil, false, result.Err()
	}

	var orders FoodOrders
	err = result.Decode(&orders)
	if err != nil {
		return nil, nil, false, err
	}

	isEated := false
	if session == "breakfast" {
		isEated = !isExist(orders.StatusNonBreakast, removeValue)
	} else {
		isEated = !isExist(orders.StatusNonLunch, removeValue)
	}
	food := getCurrentFood(orders, session)

	return food, &orders.User, isEated, nil
}

func isExist(arr []int, value int) bool {
	for _, v := range arr {
		if value == v {
			return true
		}
	}
	return false
}

type NewNonOrder struct {
	User              User        `json:"user,omitempty"`
	StatusNonBreakast []int       `json:"status_non_breakfast,omitempty" bson:"status_non_breakfast"`
	StatusNonLunch    []int       `json:"status_non_lunch,omitempty" bson:"status_non_lunch"`
	GetDay            interface{} `json:"get_day,omitempty"`
}

func GetNonFood(id string) (*[]NewNonOrder, error) {
	menu, err := GetMenu(id)
	if err != nil {
		return nil, err
	}
	start := menu.Start
	end := menu.End

	filter := bson.M{
		"start": start,
		"end":   end,
	}
	collection := db.client.Database(DB).Collection(CollectionFoodOrders)
	result, err := collection.Find(db.context, filter)
	if err != nil {
		return nil, err
	}
	var foodOrder []FoodOrders
	err = result.All(db.context, &foodOrder)
	if err != nil {
		return nil, err
	}

	var newNonOrder []NewNonOrder
	for _, food := range foodOrder {
		newNonOrder = append(newNonOrder, NewNonOrder{
			User:              food.User,
			StatusNonBreakast: food.StatusNonBreakast,
			StatusNonLunch:    food.StatusNonLunch,
		})
	}

	return &newNonOrder, nil
}

func GetFoodOrderDetail(name string, id string) (*FoodOrders, error) {
	menu, err := GetMenu(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{
		"user.user": name,
		"start":     menu.Start,
		"end":       menu.End,
	}
	collection := db.client.Database(DB).Collection(CollectionFoodOrders)
	result := collection.FindOne(db.context, filter)
	if result.Err() != nil {
		return nil, result.Err()
	}
	var orders FoodOrders
	err = result.Decode(&orders)
	if err != nil {
		return nil, err
	}
	return &orders, nil
}

func DeleteFoodOrder(id string) error {
	objectId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objectId}
	collection := db.client.Database(DB).Collection(CollectionFoodOrders)
	result, err := collection.DeleteOne(db.context, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("no document was found")
	}

	return nil
}
