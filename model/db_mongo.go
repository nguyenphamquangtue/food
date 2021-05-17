package model

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Db struct {
	context context.Context
	client  *mongo.Client
}

var db = &Db{}

func (m *Db) ConnectDatabase() {
	DbUrl := os.Getenv("MONGO_URL")

	if m.context == nil || m.client == nil {
		var context = context.Background()
		var client, err = mongo.Connect(context, options.Client().ApplyURI(DbUrl))
		if err != nil {
			log.Fatal(err)
		}
		err = client.Ping(context, nil)
		if err != nil {
			log.Fatal(err)
		}
		m.client = client
		m.context = context
	}
}

func (m *Db) Disconnect() {
	_ = m.client.Disconnect(m.context)
}

func ConnectMongoDb() {
	configMongoUrl()
	db.ConnectDatabase()
}

func DisconnectMongoDb() {
	db.Disconnect()
}
