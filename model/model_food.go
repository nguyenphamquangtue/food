package model

import (
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Food struct {
	Id   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name string             `json:"name,omitempty"`
	Img  string             `json:"img,omitempty"`
}

// insert food to database
func InsertFood(food Food) error {
	collection := db.client.Database(DB).Collection(CollectionFood)
	_, err := collection.InsertOne(db.context, food)
	if err != nil {
		return err
	}
	return nil
}

// get all foods in database
func GetAllFood() (*[]Food, error) {

	collection := db.client.Database(DB).Collection(CollectionFood)
	cursor, err := collection.Find(db.context, bson.M{})
	if err != nil {
		return nil, err
	}

	var listFoods []Food
	err = cursor.All(db.context, &listFoods)
	if err != nil {
		return nil, err
	}

	return &listFoods, nil
}

func UpdateFood(food Food) error {
	filter := bson.M{
		"_id": food.Id,
	}

	collection := db.client.Database(DB).Collection(CollectionFood)
	result, err := collection.ReplaceOne(db.context, filter, food)
	if err != nil {
		return err
	}
	if result.ModifiedCount == 0 {
		return errors.New("No food update")
	}

	return nil
}

func DeleteFood(id string) error {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{
		"_id": objID,
	}
	collection := db.client.Database(DB).Collection(CollectionFood)
	_, err := collection.DeleteOne(db.context, filter)
	if err != nil {
		return err
	}
	return nil
}
