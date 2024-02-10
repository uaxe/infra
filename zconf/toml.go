package zconf

import (
	"bytes"

	toml "github.com/BurntSushi/toml"
)

type Toml struct {
	name        string
	marshaler   func(any) ([]byte, error)
	unmarshaler func([]byte, any) error
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

func (t *Toml) Name() string {
	return t.name
}

func (t *Toml) Marshal(in any) ([]byte, error) {
	return t.marshaler(in)
}

func (t *Toml) Unmarshal(in []byte, out any) error {
	return t.unmarshaler(in, out)
}
