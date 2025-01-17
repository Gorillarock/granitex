package db

import (
	"encoding/json"
	"errors"
	"granitex/model"
	"os"
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
