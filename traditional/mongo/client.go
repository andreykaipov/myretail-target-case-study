package mongo

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var URI string
var Client *mongo.Client

func init() {
	URI = os.Getenv("MONGO_URI")

	if URI == "" {
		log.Println("MONGO_URI is required for pricing information.")
		os.Exit(1)
	}

	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(URI))
	if err != nil {
		log.Fatal(err)
	}

	Client = client
}
