package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Gorillarock/granitex/db"
	"github.com/Gorillarock/granitex/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestPostTxHandler(t *testing.T) {

	// Init
	dbMock := db.NewDBInteractorMock()

	// Mock functions
	dbMock.EXPECT().InsertTx(mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("model.DocumentEntry")).Return("test_id", nil)

	params := model.PostTxBody{
		Question: "test_question",
		Answer:   "test_answer",
		Emsg:     "test_emsg",
	}

	paramsJson, err := json.Marshal(params)
	require.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/v1/handler/tx", strings.NewReader(string(paramsJson)))

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	err = PostTxHandler(ctx)
	require.NoError(t, err)

	respBody := rec.Body.String()
	require.Contains(t, respBody, "i=test_id")
	require.Contains(t, respBody, "q=test_question")

	require.Equal(t, http.StatusOK, rec.Code)
}
