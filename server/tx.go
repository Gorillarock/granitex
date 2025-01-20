// Handles tx side of server

package server

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"

	"github.com/Gorillarock/granitex/db"
	"github.com/Gorillarock/granitex/model"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

func PostTxHandler(c echo.Context) error {
	var err error
	// write contents of request body to file mock_db.json TODO: use MongoDB
	// get body contents
	var input model.PostTxBody
	if err := c.Bind(&input); err != nil {
		return c.String(http.StatusBadRequest, "Bad Request")
	}

	log.Info(input)
	// convert to documentEntry
	doc := input.ToDocumentEntry()
	doc.Verify, err = generateVerify()
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	id, err := db.DBHandler.InsertTx(c.Request().Context(), doc)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	resp := model.PostTxResponse{ID: id, Question: input.Question, Verify: doc.Verify}

	return c.String(http.StatusOK, resp.Response())
}

// generate pseudorandom number to attach to document for verification
func generateVerify() (string, error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%d", nBig.Int64()), nil
}
