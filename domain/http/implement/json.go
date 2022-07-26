package implement

import "encoding/json"

type Json struct {
	data map[string]any
}

func (j *Json) Put(key string, value any) {
	j.data[key] = value
}

func (j *Json) Remove(key string) {
	delete(j.data, key)
}

func (j *Json) Marshal() ([]byte, error) {
	return json.Marshal(j.data)
}

func NewJsonMap() Json {
	return Json{data: map[string]any{}}
}
