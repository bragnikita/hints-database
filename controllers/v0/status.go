package v0

import (
	"github.com/labstack/echo"
	"net/http"
)

func SetRoutes(e *echo.Echo) {

	e.GET(path(""), GetStatus)

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

func path(path string) string {
	return "/v0/status" + path
}
