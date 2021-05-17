package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type roleType int

const (
	RoleAdmin roleType = 1
	RoleStaff roleType = 0
	RoleScan  roleType = -1
)

type User struct {
	Id   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	User string             `validate:"required" json:"user,omitempty"`
	Ip   string             `validate:"required" json:"ip,omitempty"`
	Dept string             `validate:"required" json:"dept,omitempty"`
	Team string             `validate:"required" json:"team,omitempty"`
	Role roleType           `json:"role"`
}

// get user information
func GetUserByIp(ip string) (*User, error) {
	//exception
	if ip == "172.16.112.11" || ip == "172.16.112.12" {
		user := User{User: "food_scanner", Role: RoleScan}
		return &user, nil
	}

	//create filter
	filter := bson.M{
		"ip": ip,
	}

	//find and decode user
	collection := db.client.Database(DB).Collection(CollectionUser)
	result := collection.FindOne(db.context, filter)

	var user User
	err := result.Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUserByUsername(username string) (*User, error) {
	filter := bson.M{
		"user": username,
		"ip": primitive.Regex{Pattern: "172.16", Options: ""},
	}

	collection := db.client.Database(DB).Collection(CollectionUser)
	result := collection.FindOne(db.context, filter)

	var user User
	err := result.Decode(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// get all user
func GetAllUser() ([]User, error) {

	collection := db.client.Database(DB).Collection(CollectionUser)
	cursor, err := collection.Find(db.context, bson.M{})
	if err != nil {
		return nil, err
	}

	var listUser []User
	err = cursor.All(db.context, &listUser)
	if err != nil {
		return nil, err
	}

	return listUser, nil
}

// insert user to database
func InsertUser(user User) error {
	collection := db.client.Database(DB).Collection(CollectionUser)
	_, err := collection.InsertOne(db.context, user)
	return err
}

func InsertManyUser(body []User) error {
	collection := db.client.Database(DB).Collection(CollectionUser)
	_, _ = collection.Indexes().CreateOne(db.context, mongo.IndexModel{
		Keys: bson.M{"ip": 1},
		// Options: options.Index().SetUnique(true),
	})

	getAllUser, err := collection.Find(db.context, bson.M{})
	if err != nil {
		return err
	}

	var user []User
	err = getAllUser.All(db.context, &user)
	if err != nil {
		return err
	}

	if !(len(user) > 0) {
		var newUser []interface{}
		for _, value := range body {
			newUser = append(newUser, value)
		}
		_, err = collection.InsertMany(db.context, newUser)

		if err != nil {
			return err
		}

	} else {
		var arrayIP []string
		for _, ip := range user {
			arrayIP = append(arrayIP, ip.Ip)
		}

		var newBody []interface{}
		for _, value := range body {
			if !isContain(arrayIP, value.Ip) {
				newBody = append(newBody, value)
			}
		}
		if len(newBody) > 0 {
			_, err = collection.InsertMany(db.context, newBody)
		}
	}

	return err
}

func isContain(arr []string, value string) bool {
	for _, v := range arr {
		if value == v {
			return true
		}
	}
	return false
}

func SetRoleAdmin() error {
	listAdmin := []string{
		"172.16.5.9",
		"172.16.5.2",
		"172.16.5.13",
		"172.16.41.4",
	}
	collection := db.client.Database(DB).Collection(CollectionUser)
	for _, admin := range listAdmin {
		filter := bson.M{
			"ip": admin,
		}
		updateRole := bson.M{"$set": bson.M{"role": 1}}
		_, err := collection.UpdateOne(db.context, filter, updateRole)
		if err != nil {
			return err
		}
	}
	return nil
}

func DeleteManyUser() error {
	collection := db.client.Database(DB).Collection(CollectionUser)
	_, err := collection.DeleteMany(db.context, bson.M{})
	return err
}
