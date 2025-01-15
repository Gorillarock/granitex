package db

import (
	"encoding/json"
	"errors"
	"granitex/model"
	"os"

	"github.com/google/uuid"
)

func MockReadFromDB() (output model.DocumentEntry, err error) {
	// read body contents from file
	docJson, err := os.ReadFile("mock_db.json")
	if err != nil {
		err = errors.New("Internal Server Error: " + err.Error())
		return
	}

	err = json.Unmarshal(docJson, &output)
	if err != nil {
		err = errors.New("Internal Server Error: " + err.Error())
		return
	}

	return
}

func MockSaveToDB(doc model.DocumentEntry) (id string, err error) {
	// generate UUID to simulate MongoDB's ObjectID
	id = uuid.New().String()
	doc.ID = id

	docJson, err := json.Marshal(doc)

	// write body contents to file
	if err = os.WriteFile("mock_db.json", docJson, 0644); err != nil {
		err = errors.New("Internal Server Error: " + err.Error())
		return
	}

	return id, nil
}
