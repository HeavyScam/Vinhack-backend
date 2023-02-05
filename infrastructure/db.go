package infrastructure

import (
	"context"
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Database struct
type Database struct {
	MongoClient *mongo.Client
}

func NewDatabase() Database {

	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_DB_URL")))
	if err != nil {
		panic(err.Error())
	}
	err = client.Connect(context.TODO())
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Database connection established")

	return Database{
		MongoClient: client,
	}

}
