package utils

import (
	"gorm.io/gorm"
)

type (
	_tables struct {
		db *gorm.DB
		t  []any
	}
)

func GetMigrate() *_tables {
	return _defaultMigrate
}

var _defaultMigrate = new(_tables)

func (self *_tables) SetDB(db *gorm.DB) *_tables {
	self.db = db
	return self
}

func (self *_tables) Add(v any) *_tables {
	if self.t == nil {
		self.t = make([]any, 0)
	}
	self.t = append(self.t, v)
	return self
}

func (self *_tables) AutoMigrate() error {
	if self.t == nil {
		return nil
	}
	return self.db.AutoMigrate(self.t...)
}
