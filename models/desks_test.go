package models

import (
	"github.com/bragnikita/hints-database/util"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestDeskService_Init(t *testing.T) {
	Clear("tmp/desks")

	service := NewDeskService("tmp/desks")
	service.Init()

	if assert.NoError(t, service.Init()) {
		assert.NotNil(t, service.desks)
	}
}

func TestDeskService_Upsert(t *testing.T) {
	Clear("tmp/desks")

	service := NewDeskService("tmp/desks")
	service.Init()

	desk := Desk{
		Id:      "",
		Title:   "Title",
		NoteIds: []string{"01", "03"},
	}
	service.Upsert(&desk)

	assert.NotEmpty(t, desk.Id)
	assert.Len(t, service.desks, 1)
	if assert.FileExists(t, "tmp/desks/index.json") {
		df := loadDesk("tmp/desks/index.json")
		assert.Len(t, df.Desks, 1)

		n := df.Desks[0]
		assert.Equal(t, n.Id, desk.Id)
		assert.Equal(t, n.Title, desk.Title)
		assert.True(t, assert.ObjectsAreEqualValues(desk.NoteIds, n.NoteIds))
	}

}

func TestDeskService_Delete(t *testing.T) {
	Clear("tmp/desks")
	err := util.CopyDir("test_data/desks", "tmp/desks")
	if err != nil {
		t.Fatal(err)
	}

	service := NewDeskService("tmp/desks")
	service.Init()

	deleted := service.Delete("20191020103656")

	if assert.NotNil(t, deleted) {
		assert.Nil(t, service.Find("20191020103656"))
		desk := loadDesk("tmp/desks/index.json")
		assert.Len(t, desk.Desks, 2)
	}
}

func loadDesk(path string) *DeskFile {
	f, e := os.Open(path)
	if e != nil {
		panic(e)
	}
	defer f.Close()
	b, e := ioutil.ReadAll(f)
	if e != nil {
		panic(e)
	}
	df := NewDeskFile()
	e = df.deserialize(b)
	if e != nil {
		panic(e)
	}
	return df
}
