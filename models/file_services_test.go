package models

import (
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

func TestFsInit(t *testing.T) {
	defer Clear("tmp/notes")
	f := FileServices{
		RootFolder: "tmp/notes",
	}
	assert.NoError(t, f.Init())

}

func TestFsCreate(t *testing.T) {
	file := "tmp/notes/note_1.txt"
	defer Clear("tmp/notes")

	f := FileServices{
		RootFolder: "tmp/notes",
	}
	f.Init()

	if assert.NoError(t, f.Upsert("note_1.txt", "Data")) {

		if assert.FileExists(t, file) {
			assert.Equal(t, "Data", *ReadString(file))
		}
	}
}

func TestFsDelete(t *testing.T) {
	defer Clear("tmp/notes")
	f := FileServices{
		RootFolder: "tmp/notes",
	}
	f.Init()
	f.Upsert("somefile", "somedata")

	if assert.NoError(t, f.Delete("somefile")) {
		info, _ := os.Stat("tmp/notes/somefile")
		assert.Nil(t, info)
	}
}

func TestFsGet(t *testing.T) {
	defer Clear("tmp/notes")
	f := FileServices{
		RootFolder: "tmp/notes",
	}
	f.Init()
	err := f.Upsert("somefile", "somedata")
	if err != nil {
		t.Fatal(err)
	}

	var content string
	readerFn := func(reader io.Reader) error {
		b, e := ioutil.ReadAll(reader)
		if e != nil {
			return e
		}
		content = string(b)
		return nil
	}
	if assert.NoError(t, f.Get("somefile", readerFn)) {
		assert.Equal(t, "somedata", content)
	}
}

func ReadString(filename string) *string {
	f, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	if err != nil {
		return nil
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil
	}
	str := string(b)
	return &str
}

func Clear(filename string) {
	_ = os.RemoveAll(filename)
}
