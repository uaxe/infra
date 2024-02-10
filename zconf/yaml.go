package zconf

import (
	yaml "gopkg.in/yaml.v3"
)

type Yaml struct {
	name        string
	marshaler   func(any) ([]byte, error)
	unmarshaler func([]byte, any) error
}

var YAMLDirver = Yaml{
	name: YamlName,
	marshaler: func(in any) ([]byte, error) {
		return yaml.Marshal(in)
	},
	unmarshaler: func(in []byte, out any) error {
		return yaml.Unmarshal(in, out)
	},
}

var _ Driver = (*Yaml)(nil)

func (y *Yaml) Name() string {
	return y.name
}

func (y *Yaml) Marshal(in any) ([]byte, error) {
	return y.marshaler(in)
}

func (y *Yaml) Unmarshal(in []byte, out any) error {
	return y.unmarshaler(in, out)
}
