package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gorillarock/granitex/db"
	"github.com/Gorillarock/granitex/model"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGetRxHandler(t *testing.T) {

	// Init
	dbMock := db.NewDBInteractorMock()

	rx := model.ResponsePayloadRxHandler{
		Error:   nil,
		Status:  http.StatusOK,
		Deleted: false,
		Emsg:    "test_emsg",
	}

	// Mock functions
	dbMock.EXPECT().GetRx(mock.AnythingOfType("context.backgroundCtx"), mock.AnythingOfType("model.GetRxParams")).Return(rx)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/handler/rx?i=test_id&v=test_verify&a=test_answer", nil)

	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	ctx := e.NewContext(req, rec)

	err := GetRxHandler(ctx)

	require.NoError(t, err)

	actualResponse := model.ResponsePayloadRxHandler{}
	respBody := rec.Body.String()
	err = json.Unmarshal([]byte(respBody), &actualResponse)
	require.NoError(t, err)
	require.True(t, actualResponse.Emsg == "test_emsg")
	require.Empty(t, actualResponse.Error)

	require.Equal(t, http.StatusOK, rec.Code)
}
