package models

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"sync"
)

type (
	DeskFile struct {
		Desks []*Desk `json:"items"`
	}

	Desk struct {
		Id      string   `json:"id"`
		Title   string   `json:"title"`
		NoteIds []string `json:"notes"`
	}

	DeskService struct {
		persistence   *FileServices
		desksFilePath string
		desks         []*Desk
		cacheMutex    *sync.RWMutex
	}
)

var Desks = *NewDeskService("desks")

func NewDeskFile() *DeskFile {
	return &DeskFile{}
}

func NewDeskService(root string) *DeskService {
	return &DeskService{
		persistence: &FileServices{
			RootFolder: root,
		},
		desksFilePath: "index.json",
		cacheMutex:    &sync.RWMutex{},
	}
}

func (d *DeskFile) serialize() (*string, error) {
	b, e := json.Marshal(d)
	if e != nil {
		return nil, e
	}
	s := string(b)
	return &s, nil
}

func (d *DeskFile) deserialize(b []byte) error {
	e := json.Unmarshal(b, d)
	if e != nil {
		return e
	}
	return nil
}

func (s *DeskService) Init() error {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	return s.reload()
}

func (s *DeskService) reload() error {
	err := s.persistence.Init()
	if err != nil {
		return err
	}
	err = s.persistence.Get(s.desksFilePath, func(reader io.Reader) error {
		b, e := ioutil.ReadAll(reader)
		if e != nil {
			return e
		}
		if len(b) == 0 {
			s.desks = []*Desk{}
			return nil
		}
		file := DeskFile{}
		e = file.deserialize(b)
		if e == nil {
			s.desks = file.Desks
		}
		return e
	})
	if os.IsNotExist(err) {
		s.desks = []*Desk{}
		return nil
	}
	return err
}

func (s *DeskService) Find(id string) *Desk {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	for _, d := range s.desks {
		if d.Id == id {
			return d
		}
	}
	return nil
}

func (s *DeskService) Upsert(desk *Desk) *Desk {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	isNew := desk.Id == ""
	if isNew {
		desk.Id = s.persistence.NextNodeId()
	}

	if isNew {
		s.desks = append(s.desks, desk)
	} else {
		index := -1
		for i, v := range s.desks {
			if v.Id == desk.Id {
				index = i
				break
			}
		}
		if index == -1 {
			s.desks = append(s.desks, desk)
		} else {
			s.desks[index] = desk
		}
	}

	if err := s.persist(); err != nil {
		_ = s.reload()
		panic(err)
	}

	return desk
}

func (s *DeskService) Delete(id string) *Desk {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	newIndex := make([]*Desk, 0, len(s.desks))
	var deleted *Desk
	for _, desk := range s.desks {
		if desk.Id != id {
			newIndex = append(newIndex, desk)
		} else {
			deleted = desk
		}
	}
	s.desks = newIndex

	if err := s.persist(); err != nil {
		_ = s.reload()
		panic(err)
	}
	return deleted
}

func (s *DeskService) persist() error {
	file := &DeskFile{
		Desks: s.desks,
	}
	str, err := file.serialize()
	if err != nil {
		return err
	}
	err = s.persistence.Upsert(s.desksFilePath, *str)
	if err != nil {
		return err
	}
	return nil
}
