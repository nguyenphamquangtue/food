package model

import (
	"errors"
	"math/big"
	"net"
	"time"
)

var mapWeek = make(map[time.Weekday]string)
var mapWeekValue = make(map[string]int)

func GenerateMapWeek() {
	mapWeek[time.Monday] = "mon"
	mapWeek[time.Tuesday] = "tue"
	mapWeek[time.Wednesday] = "wed"
	mapWeek[time.Thursday] = "thu"
	mapWeek[time.Friday] = "fri"

	mapWeekValue["mon"] = 2
	mapWeekValue["tue"] = 3
	mapWeekValue["wed"] = 4
	mapWeekValue["thu"] = 5
	mapWeekValue["fri"] = 6
}

func getCurrentSession() (string, error) {
	now := time.Now().Hour()
	if now >= 7 && now < 11 {
		return "breakfast", nil
	}
	if now >= 11 && now <= 19 {
		return "lunch", nil
	}
	return "unknow", errors.New("Invalid time")
}

func getCurrentFood(orders FoodOrders, session string) *Food {
	var food Food
	day := time.Now().Weekday()
	switch day {
	case time.Monday:
		if session == "breakfast" {
			food = orders.Mon.Breakfast
		} else {
			food = orders.Mon.Lunch
		}
	case time.Tuesday:
		if session == "breakfast" {
			food = orders.Tue.Breakfast
		} else {
			food = orders.Tue.Lunch
		}
	case time.Wednesday:
		if session == "breakfast" {
			food = orders.Wed.Breakfast
		} else {
			food = orders.Wed.Lunch
		}
	case time.Thursday:
		if session == "breakfast" {
			food = orders.Thu.Breakfast
		} else {
			food = orders.Thu.Lunch
		}
	case time.Friday:
		if session == "breakfast" {
			food = orders.Fri.Breakfast
		} else {
			food = orders.Fri.Lunch
		}
	default:
		return nil
	}
	return &food
}

func IP4toInt(IPv4Addr string) int64 {
	IPv4Int := big.NewInt(0)
	IPv4Int.SetBytes(net.ParseIP(IPv4Addr).To4())
	return IPv4Int.Int64()
}
