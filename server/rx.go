// handles rx side of server

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"granitex/db"
	"granitex/model"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

var (
	ERROR_UNAUTHORIZED = errors.New("Unauthorized")
)

func GetRxHandler(c echo.Context) error {
	resp := model.ResponsePayloadRxHandler{}
	// get query parameters
	var input model.GetRxParams
	if err := c.Bind(&input); err != nil {
		resp.Error = fmt.Errorf("Bad Request: %v", err)
		resp.Status = http.StatusBadRequest
		return c.String(resp.Status, getJson(resp))
	}

	resp = db.DBHandler.GetRx(c.Request().Context(), input)
	if resp.Error != nil {
		log.Error(resp.Error)
	}

	if resp.Deleted && resp.Status == http.StatusUnauthorized {
		resp.Error = fmt.Errorf("Unauthorized Access Attempts Exceeded: Message Destroyed")
	}

	return c.String(resp.Status, getJson(resp))
}

func getJson(resp model.ResponsePayloadRxHandler) string {
	respJson, err := json.Marshal(resp)
	if err != nil {
		log.Error(err)
		return "{\"Error\": \"Internal Server Error\"}"
	}
	return string(respJson)
}
