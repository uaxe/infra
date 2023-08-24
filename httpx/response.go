package httpx

import (
	"net/http"

	"github.com/uaxe/infra/mapping"
)

type ResponseParse interface {
	Parse(*http.Response, any) error
	ParseHeader(*http.Response, any) error
	ParseBody(*http.Response, any) error
}

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
	return mapping.MapHeader(obj, r.Header)
}

func (self *DefaultParse) ParseBody(r *http.Response, obj any) error {
	return mapping.DecodeJSON(r.Body, obj)
}
