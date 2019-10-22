package middlewares

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBuildJwtMiddleware(t *testing.T) {

	fn := BuildJwtMiddleware("jwtsecret")

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4ODcwNzg5OTQsImp0aSI6Im5pa2l0YSJ9.tTFQsdX3uOCFuocYE4Ts8y3w0nBr7CfLRbqzJkq4BAU"

	pass := false
	valid := false
	user := ""

	h := func(ctx echo.Context) error {
		pass = true
		token := ctx.Get("user").(*jwt.Token)
		claims := token.Claims.(*jwt.StandardClaims)
		user = claims.Id
		valid = token.Valid
		return nil
	}

	wrapped := fn(h)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add(echo.HeaderAuthorization, "Bearer "+token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if assert.NoError(t, wrapped(c)) {
		assert.True(t, pass)
		assert.True(t, valid)
		assert.Equal(t, user, "nikita")
	}

}

func TestBuildJwtMiddlewareWrongSecret(t *testing.T) {

	fn := BuildJwtMiddleware("other_secret")

	token := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE4ODcwNzg5OTQsImp0aSI6Im5pa2l0YSJ9.tTFQsdX3uOCFuocYE4Ts8y3w0nBr7CfLRbqzJkq4BAU"

	pass := false
	h := func(ctx echo.Context) error {
		pass = true
		return nil
	}

	wrapped := fn(h)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add(echo.HeaderAuthorization, "Bearer "+token)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := wrapped(c)
	if assert.Error(t, err) {
		httpErr := err.(*echo.HTTPError)
		assert.Equal(t, httpErr.Code, 401)
		assert.False(t, pass)
	}

}
