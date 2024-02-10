package zconf

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	drivers = map[string]Driver{
		".json": &JSONDirver,
		".toml": &TOMLDirver,
		".yaml": &YAMLDirver,
		".yml":  &YAMLDirver,
	}
)

func Load(fpath string, val any) error {
	ext := filepath.Ext(fpath)
	dirver, ok := drivers[ext]
	if !ok {
		return fmt.Errorf("%s dirver not support", ext)
	}
	fi, err := os.Open(fpath)
	if err != nil {
		return err
	}
	defer func() { _ = fi.Close() }()
	if _, err = fi.Stat(); err != nil {
		return err
	}
	raw, err := io.ReadAll(fi)
	if err != nil {
		return err
	}
	return dirver.Unmarshal(raw, val)
}

func LoadJSONBytes(raw []byte, val any) error {
	return JSONDirver.Unmarshal(raw, val)
}

func LoadTOMLBytes(raw []byte, val any) error {
	return TOMLDirver.Unmarshal(raw, val)
}

func LoadYAMLBytes(raw []byte, val any) error {
	return YAMLDirver.Unmarshal(raw, val)
}
