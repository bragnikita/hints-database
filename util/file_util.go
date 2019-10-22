package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func TryUnmarshall(filename string, target interface{}) error {
	f, e := os.Open(filename)
	if e != nil {
		return e
	}
	defer f.Close()
	b, e := ioutil.ReadAll(f)
	if e != nil {
		return e
	}
	return json.Unmarshal(b, target)
}
