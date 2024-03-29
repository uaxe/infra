package zconf_test

import (
	"io/fs"
	"path/filepath"
	"strings"
	"testing"

	"github.com/uaxe/infra/zconf"
)

type Config struct {
	Settings struct {
		Database struct {
			Driver          string `yaml:"driver"`
			Source          string `yaml:"source"`
			ConnMaxIdleTime int    `yaml:"connMaxIdleTime"`
			ConnMaxLifeTime int    `yaml:"connMaxLifeTime"`
			MaxIdleConns    int    `yaml:"monnMaxLifeTime"`
			MaxOpenConns    int    `yaml:"maxOpenConns"`
		} `yaml:"database"`
	} `yaml:"settings"`
}

func TestLoad(t *testing.T) {
	err := filepath.Walk("./testdata", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		var cfg Config
		err = zconf.Load(path, &cfg)
		if err != nil {
			t.Fatalf("%v", err)
		}
		if !strings.EqualFold(cfg.Settings.Database.Driver, "mysql") {
			t.Fatalf("%v", cfg.Settings.Database)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("%v", err)
	}
}
