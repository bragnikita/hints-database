package main

import (
	v0 "github.com/bragnikita/hints-database/controllers/v0"
	"github.com/bragnikita/hints-database/middlewares"
	"github.com/bragnikita/hints-database/models"
	. "github.com/bragnikita/hints-database/util"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)
import "github.com/labstack/echo/middleware"

func main() {
	HandleError(InitConfig(), log.Fatal)
	HandleError(models.Notes.Init(), log.Fatal)
	HandleError(models.Desks.Init(), log.Fatal)

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	public := e.Group("/auth")
	v0.SetStatusRoutes(public)

	secured := e.Group("/")
	secured.Use(middlewares.BuildJwtMiddleware(AppConfig.JwtSecret))
	v0.SetNotesRoutes(secured)

	e.Logger.Fatal(e.Start(":3001"))
}
