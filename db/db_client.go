package db

import (
	"context"
	"fmt"
	"granitex/model"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DBHandler DBClient

const (
	DATABASE            = "granitex"
	COLLECTION_MESSAGES = "messages"
)

type DBClient struct {
	Client *mongo.Client
	user   string
	pass   string
	uri    string
}

func (d *DBClient) init() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	d.user = os.Getenv("DB_USERNAME")
	d.pass = os.Getenv("DB_PASSWORD")
	if d.user == "" || d.pass == "" {
		log.Fatal("DB_USER or DB_PASS not found in .env file")
	}
	d.uri = "mongodb://" + d.user + ":" + d.pass + "@localhost:27017"
	return err
}

func (d *DBClient) connect(ctx context.Context) error {
	clientOptions := options.Client().ApplyURI(d.uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return err
	}
	d.Client = client
	log.Println("Connected to MongoDB")
	return nil
}

func (d *DBClient) disconnect(ctx context.Context) error {
	return d.Client.Disconnect(ctx)
}

func NewDBClient() error {
	DBHandler = DBClient{}
	return DBHandler.init()
}

// Inserts a TX Post request into the database
func (d *DBClient) InsertTx(ctx context.Context, document model.DocumentEntry) (string, error) {
	d.connect(ctx)
	defer d.disconnect(ctx)
	db := d.Client.Database(DATABASE)
	collection := db.Collection(COLLECTION_MESSAGES)

	result, err := collection.InsertOne(ctx, document)
	if err != nil {
		return "", err
	}

	oID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("Failed to convert InsertedID to ObjectID")
	}

	id := oID.Hex()

	return id, nil
}

func (d *DBClient) GetRx(ctx context.Context, params model.GetRxParams) model.ResponsePayloadRxHandler {
	functionName := "GetRx"
	resp := model.ResponsePayloadRxHandler{Emsg: "", Status: http.StatusOK}
	d.connect(ctx)
	defer d.disconnect(ctx)
	db := d.Client.Database(DATABASE)
	collection := db.Collection(COLLECTION_MESSAGES)
	id, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		resp.Status = http.StatusBadRequest
		resp.Error = err
		return resp
	}

	filter := bson.M{"_id": id, "verify": params.Verify}

	var result model.DocumentEntry
	err = collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		resp.Status = http.StatusNotFound
		resp.Error = err
		return resp
	}

	answeredCorrectly := verifyAnswer(params.Answer, result.Answer)
	if !answeredCorrectly {
		fmt.Printf("%s: Unauthorized Access Attempt", functionName)
		resp.Status = http.StatusUnauthorized
		// Increment try
		result.Try++
		if result.Try >= 3 {
			// Delete document from DB
			deleteResult, err := collection.DeleteOne(ctx, filter)
			if err != nil {
				fmt.Printf("Error in %s: %v", functionName, err)
				// Do not add the error to the response, as user is not authorized
				return resp
			}
			if deleteResult.DeletedCount == 0 {
				// Do not add the error to the response, as user is not authorized
				fmt.Printf("Error in %s: Failed to delete document", functionName)
				return resp
			}
			resp.Deleted = true
		} else {
			// Update result.Try in DB
			update := bson.M{"$set": bson.M{"try": result.Try}}
			timedCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()

			err = collection.FindOneAndUpdate(timedCtx, filter, update).Decode(&result)
			if err != nil {
				// Do not add the error to the response, as user is not authorized
				fmt.Printf("Error in %s: %v", functionName, err)
			}
			return resp
		}
	} else {
		resp.Emsg = result.Emsg
		// Delete document from DB
		deleteResult, err := collection.DeleteOne(ctx, filter)
		if err != nil {
			// Do not add the error to the response, as this is an internal error
			fmt.Printf("Error in %s: %v", functionName, err)
			return resp
		}
		if deleteResult.DeletedCount == 0 {
			// Do not add the error to the response, as this is an internal error
			fmt.Printf("Error in %s: Failed to delete document", functionName)
		}
		resp.Deleted = true
	}

	if resp.Status == http.StatusOK && resp.Emsg == "" {
		resp.Status = http.StatusInternalServerError
		resp.Error = model.ERROR_UNKNOWN
	}

	return resp
}

func verifyAnswer(givenAnswer string, expectedAnswer string) bool {
	return givenAnswer == expectedAnswer
}
