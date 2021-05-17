package main

import (
	"fmt"
	"food/api"
	"food/model"
	"log"
)

const port = ":7002"

func main() {
	// connect database
	model.GenerateMapWeek()
	model.ConnectMongoDb()
	defer model.DisconnectMongoDb()

	// model.FakeInserMenu()

	router := api.NewRouter()
	fmt.Println("Start server port 7002")
	log.Fatal(router.Run(port))
}
