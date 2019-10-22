package v0

import (
	"github.com/bragnikita/hints-database/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"net/http"
	"time"
)

type (
	RequestAuth struct {
		Username string `json:"username",required:"true"`
		Password string `json:"password",required:"true"`
	}
)

func SetStatusRoutes(e *echo.Group) {
	e.GET(path("/status"), GetStatus)
	e.POST(path("/auth"), Auth)
}

func GetStatus(c echo.Context) error {
	contentType := c.Request().Header.Get("Accept")
	if contentType == "" || contentType == "*/*" {
		contentType = echo.MIMETextPlain
	}
	c.Response().Header().Set(echo.HeaderContentType, contentType)
	switch contentType {
	case echo.MIMETextPlain:
		return c.String(http.StatusOK, "OK")
	case echo.MIMEApplicationJSON:
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "OK",
		})
	default:
		return c.NoContent(http.StatusOK)
	}
}

func Auth(e echo.Context) (err error) {
	var req RequestAuth
	if err = e.Bind(&req); err != nil {
		return err
	}
	if util.AppConfig.Username == req.Username &&
		util.AppConfig.Password == req.Password {

		claims := jwt.StandardClaims{
			Id:        req.Username,
			ExpiresAt: time.Now().Add(time.Hour * 24 * 365 * 10).Unix(),
		}

		token := jwt.New(jwt.SigningMethodHS256)
		token.Claims = claims

		tokenString, err := token.SignedString([]byte(util.AppConfig.JwtSecret))
		util.MustNotError(err)

		return e.JSON(http.StatusOK, map[string]interface{}{
			"token": tokenString,
		})
	}
	return e.NoContent(http.StatusUnauthorized)
}

func path(path string) string {
	return "/v0" + path
}
