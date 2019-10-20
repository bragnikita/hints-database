package models

import (
	"github.com/bragnikita/hints-database/controllers/util"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestInitNotesIndex(t *testing.T) {
	service := NewNotesService("test_data/notes")
	if assert.NoError(t, service.Init()) {
		assert.Len(t, service.Index, 3)
		for _, rec := range service.Index {
			assert.NotEmpty(t, rec.Id)
			assert.NotEmpty(t, rec.Title)
			assert.NotEmpty(t, rec.Created)
			assert.NotEmpty(t, rec.Description)
		}
	}
}

func TestGetFullListing(t *testing.T) {
	service := NewNotesService("test_data/notes")
	service.Init()

	res := service.GetIndex(NotesFilter{
		Content: true,
	})
	if assert.Len(t, res, 3) {
		assert.NotEmpty(t, res[0].Content)
		assert.Equal(t, res[0].Id, "01")
	}

}

func TestGetFiltered(t *testing.T) {
	service := NewNotesService("test_data/notes")
	service.Init()
	assert.Len(t, service.GetIndex(NotesFilter{
		Ids: []string{"01", "03"},
	}), 2)
}

func TestNotesService_Save(t *testing.T) {
	Clear("tmp/notes")
	check(util.CopyDir("test_data/notes", "tmp/notes"))

	service := NewNotesService("tmp/notes")
	check(service.Init())

	note := Note{
		NotesIndexRecord: NotesIndexRecord{Title: "New note"},
		Content:          "New note content",
	}
	if assert.NoError(t, service.Save(&note)) {
		fn := "tmp/notes/" + note.Id + ".json"
		assert.FileExists(t, fn)

		var saved Note
		if assert.NoError(t, util.TryUnmarshall(fn, &saved)) {
			assert.Equal(t, saved.Title, "New note")
			assert.Equal(t, saved.Content, "New note content")
		}
	}
}

func TestNotesService_Delete(t *testing.T) {
	Clear("tmp/notes")
	check(util.CopyDir("test_data/notes", "tmp/notes"))

	service := NewNotesService("tmp/notes")
	check(service.Init())

	deleted := service.Delete("01")
	if assert.NotNil(t, deleted) {
		fn := "tmp/notes/01.json"
		_, err := os.Stat(fn)
		assert.Error(t, err)
		assert.Len(t, service.Index, 2)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
