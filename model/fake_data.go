package model

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// fake insert user to database
func FakeInsertUser() {
	var listUser []interface{}
	// for i := 0; i < 300000; i++ {
	// 	username := fmt.Sprintf("user %v", i)
	// 	ip := "192.167.1.89"
	// 	listUser = append(listUser, User{Username: username, Ip: ip})
	// }
	listUser = append(listUser, User{User: "lebron", Ip: "172.16.41.15"})
	listUser = append(listUser, User{User: "andy", Ip: "172.16.41.14"})
	listUser = append(listUser, User{User: "kul", Ip: "172.16.41.4"})

	collection := db.client.Database(DB).Collection(CollectionUser)
	name, err := collection.Indexes().CreateOne(db.context, mongo.IndexModel{
		Keys:    bson.M{"ip": 1},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(name)
	}

	_, err = collection.InsertMany(db.context, listUser)
	if err != nil {
		fmt.Println(err.Error())
	}
}

/**
fake insert weekly menu to database
*/
func FakeInserMenu() {
	// now.WeekStartDay = time.Monday
	// fmt.Println(now.BeginningOfWeek())

	dailyMenu := DailyMenu{
		Breakfast: []Food{
			{Name: "trung kho", Img: "http://www.img.com/img.png"},
			{Name: "trung luon", Img: "http://www.img.com/img.png"},
			{Name: "trung hap", Img: "http://www.img.com/img.png"},
		},
		Lunch: []Food{
			{Name: "trung kho", Img: "http://www.img.com/img.png"},
			{Name: "trung luon", Img: "http://www.img.com/img.png"},
			{Name: "trung hap", Img: "http://www.img.com/img.png"},
		},
	}
	weekMenu := WeeklyMenu{
		Mon:    dailyMenu,
		Start:  time.Date(2020, 8, 3, 0, 0, 0, 0, time.Local).Unix() * 1000,
		End:    time.Date(2020, 8, 7, 23, 59, 59, 0, time.Local).Unix() * 1000,
		Status: StatusClose,
		Next:   false,
	}

	err := InsertWeeklyMenu(weekMenu)
	if err != nil {
		fmt.Println(err.Error())
	}

	weekMenu.Mon.Breakfast[0].Name = "trung xao"
	err = InsertWeeklyMenu(weekMenu)
	if err != nil {
		fmt.Println(err.Error())
	}

	weekMenu.Mon.Breakfast[0].Name = "trung hap"
	err = InsertWeeklyMenu(weekMenu)
	if err != nil {
		fmt.Println(err.Error())
	}

	err = InsertWeeklyMenu(weekMenu)
	if err != nil {
		fmt.Println(err.Error())
	}

}
