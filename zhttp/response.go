package zhttp

import (
	"net/http"
)

type ResponseParse interface {
	Parse(*http.Response, any) error
	ParseHeader(*http.Response, any) error
	ParseBody(*http.Response, any) error
}

var _ ResponseParse = (*DefaultParse)(nil)

type DefaultParse struct{}

func (self *DefaultParse) Parse(r *http.Response, obj any) error {
	if err := self.ParseHeader(r, obj); err != nil {
		return err
	}
	if err := self.ParseBody(r, obj); err != nil {
		return err
	}
	return nil
}

func (self *DefaultParse) ParseHeader(r *http.Response, obj any) error {
	return nil
}

func (self *DefaultParse) ParseBody(r *http.Response, obj any) error {
	return nil
}
