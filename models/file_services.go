package models

import (
	"github.com/juju/fslock"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

type (
	FileServices struct {
		RootFolder string
	}
)

func (f *FileServices) Init() error {
	return os.MkdirAll(f.RootFolder, 0766)
}

func (f *FileServices) NextNodeId() string {
	return time.Now().Format("20060102150405")
}

func (f *FileServices) fullPath(filename string) string {
	return path.Join(f.RootFolder, filename)
}

func (f *FileServices) Upsert(filename string, content string) error {
	return f.RunInLock(filename, func(fp string) error {
		f, err := os.OpenFile(fp, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0700)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(f, strings.NewReader(content))
		return err
	})
}

func (f *FileServices) Delete(filename string) error {
	return f.RunInLock(filename, func(s string) error {
		return os.Remove(s)
	})
}

func (f *FileServices) Get(filename string, call func(reader io.Reader) error) error {
	return f.RunInLock(filename, func(s string) error {
		f, err := os.OpenFile(s, os.O_CREATE|os.O_RDONLY, 0700)
		if err != nil {
			return err
		}
		defer f.Close()
		call(f)
		return nil
	})
}

func (f *FileServices) List() ([]string, error) {
	root, err := os.Open(f.RootFolder)
	if err != nil {
		return nil, err
	}
	return root.Readdirnames(-1)
}

func (f *FileServices) RunInLock(filename string, callable func(string) error) error {
	ff := f.fullPath(filename)
	lock := fslock.New(ff)
	defer lock.Unlock()
	err := lock.LockWithTimeout(10 * time.Second)
	if err != nil {
		return err
	}
	return callable(ff)
}
