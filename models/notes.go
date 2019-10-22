package models

import (
	json2 "encoding/json"
	"github.com/bragnikita/hints-database/util"
	"io"
	"io/ioutil"
	"sync"
)

type (
	NotesIndexRecord struct {
		Id          string   `json:"id"`
		Title       string   `json:"title"`
		Description string   `json:"description"`
		Created     string   `json:"created"`
		Tags        []string `json:"tags"`
	}
	Note struct {
		NotesIndexRecord `json:",inline"`
		Content          string `json:"content"`
		Meta             string `json:"meta"`
	}

	NotesFilter struct {
		Ids     []string
		Content bool
	}

	NotesService struct {
		Persistance *FileServices
		Index       []NotesIndexRecord
		cacheMutex  *sync.RWMutex
	}
)

func NewNotesService(root string) *NotesService {
	return &NotesService{
		Persistance: &FileServices{
			RootFolder: root,
		},
		cacheMutex: &sync.RWMutex{},
	}
}

var Notes = NewNotesService("data/notes")

func (n *Note) IsNew() bool {
	return n.Id == ""
}

func (s *NotesService) Init() error {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()
	util.MustDo(s.Persistance.Init())

	list, err := s.Persistance.List()
	if err != nil {
		return err
	}

	index := make([]NotesIndexRecord, 0, len(list))

	updater := func(r io.Reader) error {
		note, err := s.unmarshall(r)
		if err != nil {
			return err
		}
		index = append(index, NotesIndexRecord{
			Id:          note.Id,
			Created:     note.Created,
			Tags:        note.Tags,
			Title:       note.Title,
			Description: note.Description,
		})
		return nil
	}

	for _, notefile := range list {
		err := s.Persistance.Get(notefile, updater)
		if err != nil {
			return err
		}
	}

	s.Index = index

	return nil
}

func (s *NotesService) Save(note *Note) error {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	isNew := note.IsNew()
	if isNew {
		note.Id = s.Persistance.NextNodeId()
	}
	bytes, err := json2.Marshal(note)
	if err != nil {
		return err
	}
	str := string(bytes)
	err = s.Persistance.Upsert(s.filename(note.Id), str)
	if err == nil {
		record := NotesIndexRecord{
			Id:          note.Id,
			Created:     note.Created,
			Tags:        note.Tags,
			Title:       note.Title,
			Description: note.Description,
		}
		if isNew {
			s.Index = append(s.Index, record)
		} else {
			i := -1
			for y, rec := range s.Index {
				if rec.Id == note.Id {
					i = y
					break
				}
			}
			if i == -1 {
				s.Index = append(s.Index, record)
			} else {
				s.Index[i] = record
			}
		}
	}
	return err
}

func (s *NotesService) GetIndex(filter NotesFilter) []*Note {
	s.cacheMutex.RLock()
	defer s.cacheMutex.RUnlock()

	var result []NotesIndexRecord
	if filter.Ids != nil {
		result = make([]NotesIndexRecord, 0, len(filter.Ids))
		for _, record := range s.Index {
			for _, id := range filter.Ids {
				if record.Id == id {
					result = append(result, record)
				}
			}
		}
	}
	if result == nil {
		result = s.Index
	}
	var resultFull []*Note
	resultFull = make([]*Note, 0, len(result))
	for _, res := range result {
		err := s.Persistance.Get(s.filename(res.Id), func(reader io.Reader) error {
			n, e := s.unmarshall(reader)
			if e != nil {
				return e
			}
			resultFull = append(resultFull, n)
			return nil
		})
		if err != nil {
			panic(err)
		}
	}

	return resultFull
}

func (s *NotesService) Delete(id string) *NotesIndexRecord {
	s.cacheMutex.Lock()
	defer s.cacheMutex.Unlock()

	result := make([]NotesIndexRecord, 0, len(s.Index))
	var idx int = -1
	for i, rec := range s.Index {
		if rec.Id != id {
			result = append(result, rec)
		} else {
			idx = i
		}
	}
	if idx < 0 {
		return nil
	}
	deleted := s.Index[idx]

	err := s.Persistance.Delete(s.filename(deleted.Id))
	if err != nil {
		panic(err)
	}
	s.Index = result
	return &deleted
}

func (s *NotesService) unmarshall(reader io.Reader) (*Note, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var note Note
	err = json2.Unmarshal(b, &note)
	return &note, err
}

func (s *NotesService) filename(id string) string {
	return id + ".json"
}
