package v0

import (
	"encoding/json"
	"fmt"
	"github.com/bragnikita/hints-database/util"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
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

func TestAuth(t *testing.T) {

	os.Clearenv()
	os.Setenv("HD_USERNAME", "nikita")
	os.Setenv("HD_PASSWORD", "password")
	os.Setenv("HD_JWTSECRET", "jwtsecret")
	util.MustNotError(util.InitConfig())

	reqJson := map[string]string{
		"username": "nikita",
		"password": "password",
	}
	reqJsonStr, _ := json.Marshal(reqJson)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v0/auth", strings.NewReader(string(reqJsonStr)))
	req.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	if assert.NoError(t, Auth(c)) {
		if !assert.Equal(t, rec.Code, http.StatusOK) {
			t.FailNow()
		}
		var m map[string]interface{}
		if assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &m)) {
			fmt.Println(m["token"])
			assert.NotEmpty(t, m["token"])
		}
	}
}
