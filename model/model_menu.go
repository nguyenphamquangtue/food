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

type MenuStatus int

const (
	StatusActive MenuStatus = 1
	StatusOpen   MenuStatus = 0
	StatusClose  MenuStatus = -1
)

type DailyMenu struct {
	Breakfast []Food `json:"breakfast,omitempty"`
	Lunch     []Food `json:"lunch,omitempty"`
}

type WeeklyMenu struct {
	Id     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Mon    DailyMenu          `bson:"mon,omitempty" json:"mon,omitempty"`
	Tue    DailyMenu          `bson:"tue,omitempty" json:"tue,omitempty"`
	Wed    DailyMenu          `bson:"wed,omitempty" json:"wed,omitempty"`
	Thu    DailyMenu          `bson:"thu,omitempty" json:"thu,omitempty"`
	Fri    DailyMenu          `bson:"fri,omitempty" json:"fri,omitempty"`
	Start  int64              `json:"start,omitempty" validate:"required"`
	End    int64              `json:"end,omitempty" validate:"required"`
	Status MenuStatus         `json:"status"`
	Next   bool               `json:"next"`
	Day    map[int]string     `json:"day,omitempty"`
}

func GetDayBeetween(timeStart, timeEnd time.Time) map[int]string {
	var listDay = make(map[int]string)
	for d := timeStart; !d.After(timeEnd); d = d.AddDate(0, 0, 1) {
		switch dayTime := d.Weekday().String(); dayTime {
		case "Monday":
			listDay[2] = d.Format("2006-01-02")
		case "Tuesday":
			listDay[3] = d.Format("2006-01-02")
		case "Wednesday":
			listDay[4] = d.Format("2006-01-02")
		case "Thursday":
			listDay[5] = d.Format("2006-01-02")
		case "Friday":
			listDay[6] = d.Format("2006-01-02")
		}
	}

	return listDay
}

func InsertWeeklyMenu(menu WeeklyMenu) error {
	err := validator.New().Struct(menu)
	if err != nil {
		return err
	}
	collection := db.client.Database(DB).Collection(CollectionWeeklyMenu)

	startDay := time.Unix(menu.Start/1000, 0)
	endDay := time.Unix(menu.End/1000, 0)

	// so sanh startDay va endDay trong cung 1 tuan

	startWeek := now.With(startDay).BeginningOfWeek()
	endWeek := now.With(startDay).EndOfWeek()

	if endDay.After(endWeek) {
		return errors.New("End day not include this week!")
	}

	filter := bson.M{
		"start": bson.M{"$gte": startWeek.Unix() * 1000},
		"end":   bson.M{"$lte": endWeek.Unix() * 1000},
	}
	checkExist := collection.FindOne(db.context, filter)
	if checkExist.Err() == nil {
		return errors.New("Food Order existed!")
	}

	Day := GetDayBeetween(startDay, endDay)
	menu.Day = Day

	_, err = collection.InsertOne(db.context, menu)
	if err != nil {
		return err
	}
	return nil
}

func GetNextMenu() (*WeeklyMenu, error) {

	var menu WeeklyMenu
	filter := bson.M{"next": true}

	collection := db.client.Database(DB).Collection(CollectionWeeklyMenu)
	result := collection.FindOne(db.context, filter)

	err := result.Decode(&menu)
	if err != nil {
		return nil, err
	}

	return &menu, nil
}

func GetAllMenu() (*[]WeeklyMenu, error) {
	var allMenu []WeeklyMenu

	findOptions := options.Find()
	findOptions.SetSort(bson.M{"start": -1})
	collection := db.client.Database(DB).Collection(CollectionWeeklyMenu)
	cursor, err := collection.Find(db.context, bson.M{}, findOptions)
	if err != nil {
		return nil, err
	}

	err = cursor.All(db.context, &allMenu)

	if err != nil {
		return nil, err
	}

	return &allMenu, nil
}

func GetMenu(id string) (*WeeklyMenu, error) {
	var detailMenu WeeklyMenu
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{
		"_id": objID,
	}
	collection := db.client.Database(DB).Collection(CollectionWeeklyMenu)
	cursor := collection.FindOne(db.context, filter)
	err := cursor.Decode(&detailMenu)

	if err != nil {
		return nil, err
	}
	return &detailMenu, nil
}

func ActiveMenu(id string) error {

	// reset all to not next week and status close
	collection := db.client.Database(DB).Collection(CollectionWeeklyMenu)
	updateValue := bson.M{"$set": bson.M{"next": false, "status": StatusClose}}
	_, err := collection.UpdateMany(db.context, bson.M{"next": true}, updateValue)
	if err != nil {
		return err
	}

	// active menu with id
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{
		"_id": objID,
	}
	updateValue = bson.M{"$set": bson.M{"next": true, "status": StatusActive}}
	_, err = collection.UpdateOne(db.context, filter, updateValue)
	if err != nil {
		return err
	}

	return nil
}

func DeactiveMenu(id string) error {

	collection := db.client.Database(DB).Collection(CollectionWeeklyMenu)

	// active menu with id
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{
		"_id": objID,
	}
	updateValue := bson.M{"$set": bson.M{"status": StatusClose}}
	_, err := collection.UpdateOne(db.context, filter, updateValue)
	if err != nil {
		return err
	}

	return nil
}
