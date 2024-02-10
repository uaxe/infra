package zconf

import (
	"encoding/json"
)

type JSON struct {
	name        string
	marshaler   func(any) ([]byte, error)
	unmarshaler func([]byte, any) error
}

var JSONDirver = JSON{
	name: JsonName,
	marshaler: func(in any) ([]byte, error) {
		return json.Marshal(in)
	},
	unmarshaler: func(in []byte, out any) error {
		return json.Unmarshal(in, out)
	},
}

var _ Driver = (*JSON)(nil)

func (j *JSON) Name() string {
	return j.name
}

func (j *JSON) Marshal(in any) ([]byte, error) {
	return j.marshaler(in)
}

func (j *JSON) Unmarshal(in []byte, out any) error {
	return j.unmarshaler(in, out)
}
