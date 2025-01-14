package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
)

func main() {
	fmt.Println("Hello, World!")

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	//e.GET("/", handler)
	e.Static("/v1/src", "client/src")
	e.File("/v1/tx", "client/tx_client.html")
	e.File("/v1/rx", "client/rx_client.html")
	e.POST("/v1/handler/tx", postTxHandler)
	e.GET("/v1/handler/rx", getRxHandler)

	e.Logger.Fatal(e.Start(":1323"))
}

type getRxParams struct {
	ID     string `query:"i" json:"id"`
	Verify string `query:"v" json:"verify"`
	Answer string `query:"a" json:"answer"`
}

func getRxHandler(c echo.Context) error {
	// get query parameters
	var input getRxParams
	if err := c.Bind(&input); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	fmt.Printf("Params: %+v\n", input)

	return c.String(http.StatusOK, "{\"status\": \"OK\"}")
}

type postTxBody struct {
	ID       string `json:"id"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Emsg     string `json:"emsg"`
	Try      int    `json:"try"`
}

type postTxResponse struct {
	ID string `json:"id"`
}

func postTxHandler(c echo.Context) error {
	// write contents of request body to file mock_db.json TODO: use MongoDB
	// get body contents
	var input postTxBody
	if err := c.Bind(&input); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	log.Info(input)

	id, err := mockSaveToDB(input)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	resp := postTxResponse{ID: id}
	respJson, err := json.Marshal(resp)
	return c.String(http.StatusOK, string(respJson))
}

func mockSaveToDB(input postTxBody) (id string, err error) {
	// generate UUID to simulate MongoDB's ObjectID
	id = uuid.New().String()
	input.ID = id

	inputJson, err := json.Marshal(input)

	// write body contents to file
	if err = os.WriteFile("mock_db.json", inputJson, 0644); err != nil {
		err = errors.New("Internal Server Error: " + err.Error())
		return
	}

	return id, nil
}
