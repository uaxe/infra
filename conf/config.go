package conf

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	dirvers = map[string]Driver{
		".json": &JSONDirver,
		".toml": &TOMLDirver,
		".yaml": &YAMLDirver,
		".yml":  &YAMLDirver,
	}
)

func Load(fpath string, val any, opts ...OptionFunc) error {
	ext := filepath.Ext(fpath)
	dirver, ok := dirvers[ext]
	if !ok {
		return fmt.Errorf("%s dirver not support", ext)
	}
	fi, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer fi.Close()
	if _, err = fi.Stat(); err != nil {
		return err
	}
	raw, err := io.ReadAll(fi)
	if err != nil {
		return err
	}
	if err = dirver.Unmarshal(raw, val); err != nil {
		return err
	}
	return nil
}

func LoadJSONBytes(raw []byte, val any, opts ...OptionFunc) error {
	return JSONDirver.Unmarshal(raw, val)
}

func LoadTOMLBytes(raw []byte, val any, opts ...OptionFunc) error {
	return TOMLDirver.Unmarshal(raw, val)
}

func LoadYAMLBytes(raw []byte, val any, opts ...OptionFunc) error {
	return YAMLDirver.Unmarshal(raw, val)
}
