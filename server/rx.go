// handles rx side of server

package server

import (
	"fmt"
	"granitex/db"
	"granitex/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetRxHandler(c echo.Context) error {
	// get query parameters
	var input model.GetRxParams
	if err := c.Bind(&input); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	fmt.Printf("Params: %+v\n", input)

	// read contents of file mock_db.json TODO: use MongoDB
	doc, err := db.MockReadFromDB()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	// check if ID exists
	foundMatch := doc.VerifyInput(input)
	if !foundMatch {
		return c.String(http.StatusNotFound, "Not Found")
	}

	// check if answer is correct
	match, mustDelete := doc.CheckAnswer(input.Answer)

	// delete if must
	// todo: implement
	_ = mustDelete

	if !match {
		return c.String(http.StatusUnauthorized, "Unauthorized.")
	}

	// Return eMsg
	return c.String(http.StatusOK, doc.EmsgResponse())
}
