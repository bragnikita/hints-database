package v0

import (
	"github.com/bragnikita/hints-database/models"
	"github.com/labstack/echo"
	"net/http"
)

type (
	NotesApi struct {
	}
)

func SetNotesRoutes(e *echo.Echo) {

	api := NotesApi{}

	e.GET(notesPath(""), api.List)
	e.GET(notesPath("/full"), api.ListFull)
	e.GET(notesPath("/:id"), api.Get)
	e.POST(notesPath("/:id"), api.Upsert)
	e.DELETE(notesPath("/:id"), api.Delete)
}

func (api *NotesApi) Load(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

func (api *NotesApi) ListFull(c echo.Context) error {
	return api.list(c, true)
}

func (api *NotesApi) List(c echo.Context) error {
	return api.list(c, false)
}

func (api *NotesApi) list(c echo.Context, full bool) error {
	filter := models.NotesFilter{
		Content: full,
	}
	deskId := c.QueryParam("desk")
	if deskId != "" {
		deck := models.Desks.Find(deskId)
		if deck == nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"message": "Desk with id " + deskId + " not found",
			})
		}
		filter.Ids = deck.NoteIds
	}

	list := models.Notes.GetIndex(filter)

	return c.JSON(http.StatusOK, map[string]interface{}{
		"list": list,
	})
}

func (api *NotesApi) Get(c echo.Context) error {
	filter := models.NotesFilter{
		Ids:     []string{c.Param("id")},
		Content: true,
	}
	list := models.Notes.GetIndex(filter)
	if len(list) == 0 {
		return c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"item": list[0],
	})
}

func (api *NotesApi) Upsert(c echo.Context) error {

	var note models.Note
	if err := c.Bind(&note); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	e := models.Notes.Save(&note)
	if e != nil {
		return c.String(http.StatusInternalServerError, e.Error())
	}
	return c.NoContent(http.StatusNotImplemented)
}

func (api *NotesApi) Delete(c echo.Context) error {
	id := c.Param("id")
	deleted := models.Notes.Delete(id)
	if deleted == nil {
		return c.NoContent(http.StatusNotFound)
	}
	return c.NoContent(http.StatusOK)
}

func notesPath(path string) string {
	return "/v0/nodes" + path
}
