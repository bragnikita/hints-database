package main

import (
	v0 "github.com/bragnikita/hints-database/controllers/v0"
	"github.com/bragnikita/hints-database/models"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
)
import "github.com/labstack/echo/middleware"

func main() {
	if err := models.Notes.Init(); err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	v0.SetRoutes(e)
	v0.SetNotesRoutes(e)

	e.Logger.Fatal(e.Start(":3001"))
}
