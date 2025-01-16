package db

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DBHandler DBClient

type DBClient struct {
	Client *mongo.Client
	user   string
	pass   string
	uri    string
}

func (d *DBClient) init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}

	d.user = os.Getenv("DB_USERNAME")
	d.pass = os.Getenv("DB_PASSWORD")
	if d.user == "" || d.pass == "" {
		log.Fatal("DB_USER or DB_PASS not found in .env file")
	}
	d.uri = "mongodb://localhost:27017"
}

func (d *DBClient) Connect(ctx context.Context) error {
	if d.user == "" || d.pass == "" {
		d.init()
	}
	clientOptions := options.Client().ApplyURI(d.uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	d.Client = client
	log.Println("Connected to MongoDB")
	return nil
}
