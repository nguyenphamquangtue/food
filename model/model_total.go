package model

import (
	"encoding/json"

	"go.mongodb.org/mongo-driver/bson"
)

type Times struct {
	Breakfast interface{} `json:"breakfast,omitempty"`
	Lunch     interface{} `json:"lunch,omitempty"`
}

type TotalFood struct {
	Mon Times `json:"mon,omitempty"`
	Tue Times `json:"tue,omitempty"`
	Wed Times `json:"wed,omitempty"`
	Thu Times `json:"thu,omitempty"`
	Fri Times `json:"fri,omitempty"`
}

func printUniqueValue(arr []string) map[string]int {
	dict := make(map[string]int)
	for _, num := range arr {
		if num != "" {
			dict[num] = dict[num] + 1
		}
	}
	return dict
}

func TotalDaily(id string) (*TotalFood, error) {
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
	cur, err := collection.Find(db.context, filter)
	if err != nil {
		return nil, err
	}

	var foodOrders []FoodOrders
	err = cur.All(db.context, &foodOrders)
	if err != nil {
		return nil, err
	}

	var breakMon, breakTue, breakWed, breakThu, breakFri []string
	var lunchMon, lunchTue, lunchWed, lunchThu, lunchFri []string

	for _, food := range foodOrders {
		breakMon = append(breakMon, food.Mon.Breakfast.Name)
		breakTue = append(breakTue, food.Tue.Breakfast.Name)
		breakWed = append(breakWed, food.Wed.Breakfast.Name)
		breakThu = append(breakThu, food.Thu.Breakfast.Name)
		breakFri = append(breakFri, food.Fri.Breakfast.Name)

		lunchMon = append(lunchMon, food.Mon.Lunch.Name)
		lunchTue = append(lunchTue, food.Tue.Lunch.Name)
		lunchWed = append(lunchWed, food.Wed.Lunch.Name)
		lunchThu = append(lunchThu, food.Thu.Lunch.Name)
		lunchFri = append(lunchFri, food.Fri.Lunch.Name)
	}

	getbreakMon := printUniqueValue(breakMon)
	getbreakTue := printUniqueValue(breakTue)
	getbreakWed := printUniqueValue(breakWed)
	getbreakThu := printUniqueValue(breakThu)
	getbreakFri := printUniqueValue(breakFri)

	getlunchMon := printUniqueValue(lunchMon)
	getlunchTue := printUniqueValue(lunchTue)
	getlunchWed := printUniqueValue(lunchWed)
	getlunchThu := printUniqueValue(lunchThu)
	getlunchFri := printUniqueValue(lunchFri)

	var totalFood TotalFood

	totalFood.Mon.Breakfast = getbreakMon
	totalFood.Tue.Breakfast = getbreakTue
	totalFood.Wed.Breakfast = getbreakWed
	totalFood.Thu.Breakfast = getbreakThu
	totalFood.Fri.Breakfast = getbreakFri

	totalFood.Mon.Lunch = getlunchMon
	totalFood.Tue.Lunch = getlunchTue
	totalFood.Wed.Lunch = getlunchWed
	totalFood.Thu.Lunch = getlunchThu
	totalFood.Fri.Lunch = getlunchFri

	return &totalFood, nil
}

func TotalWeek(id string) (*interface{}, error) {

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
	cur, err := collection.Find(db.context, filter)
	if err != nil {
		return nil, err
	}

	var foodOrders []FoodOrders
	err = cur.All(db.context, &foodOrders)
	if err != nil {
		return nil, err
	}

	var listfood []string

	for _, food := range foodOrders {
		listfood = append(listfood, food.Mon.Breakfast.Name)
		listfood = append(listfood, food.Tue.Breakfast.Name)
		listfood = append(listfood, food.Wed.Breakfast.Name)
		listfood = append(listfood, food.Thu.Breakfast.Name)
		listfood = append(listfood, food.Fri.Breakfast.Name)

		listfood = append(listfood, food.Mon.Lunch.Name)
		listfood = append(listfood, food.Tue.Lunch.Name)
		listfood = append(listfood, food.Wed.Lunch.Name)
		listfood = append(listfood, food.Thu.Lunch.Name)
		listfood = append(listfood, food.Fri.Lunch.Name)
	}

	getfood := printUniqueValue(listfood)
	data, err := json.Marshal(getfood)
	if err != nil {
		return nil, err
	}
	var result interface{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, err
}
