package zconf

import (
	"bytes"
	toml "github.com/BurntSushi/toml"
)

type Toml struct {
	name, indent string
	marshaler    func(any) ([]byte, error)
	unmarshaler  func([]byte, any) error
}

var TOMLDirver = JSON{
	name: TomlName,
	marshaler: func(in any) ([]byte, error) {
		var buff bytes.Buffer
		err := toml.NewEncoder(&buff).Encode(in)
		return buff.Bytes(), err
	},
	unmarshaler: func(in []byte, out any) error {
		return toml.Unmarshal(in, out)
	},
}

var _ Driver = (*Toml)(nil)

func (self *Toml) Name() string {
	return self.name
}

func (self *Toml) Marshal(in any) ([]byte, error) {
	return self.marshaler(in)
}

func (self *Toml) Unmarshal(in []byte, out any) error {
	return self.unmarshaler(in, out)
}
