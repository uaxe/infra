package zconf

import (
	"encoding/json"
)

type JSON struct {
	name, indent string
	marshaler    func(any) ([]byte, error)
	unmarshaler  func([]byte, any) error
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

func (self *JSON) Name() string {
	return self.name
}

func (self *JSON) Marshal(in any) ([]byte, error) {
	return self.marshaler(in)
}

func (self *JSON) Unmarshal(in []byte, out any) error {
	return self.unmarshaler(in, out)
}
