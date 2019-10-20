package v0

import (
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetStatusPlainText(t *testing.T) {
	// Setuo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v0/status", nil)
	req.Header.Set(echo.HeaderAccept, echo.MIMETextPlain)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, GetStatus(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, echo.MIMETextPlain, rec.Header().Get(echo.HeaderContentType))
		assert.Equal(t, "OK", rec.Body.String())
	}
}

func TestGetStatusJSON(t *testing.T) {
	// Setuo
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v0/status", nil)
	req.Header.Set(echo.HeaderAccept, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, GetStatus(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, echo.MIMEApplicationJSON, rec.Header().Get(echo.HeaderContentType))

		var response map[string]interface{}
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response)) {
			assert.Equal(t, "OK", response["message"])
		}
	}
}
