package zhttp

import (
	"bytes"
	"encoding/json"
	"github.com/uaxe/infra/binder"
	"io"
	"net/http"
)

type ResponseParse interface {
	Parse(*http.Response, any) error
	ParseHeader(*http.Response, any) error
	ParseBody(*http.Response, any) error
}

var (
	_                ResponseParse = (*defaultParse)(nil)
	DefaultRespParse               = &defaultParse{}
)

type defaultParse struct{}

func (self *defaultParse) Parse(r *http.Response, obj any) error {
	if err := self.ParseBody(r, obj); err != nil {
		return err
	}

	if err := self.ParseHeader(r, obj); err != nil {
		return err
	}

	return nil
}

func (self *defaultParse) ParseHeader(r *http.Response, obj any) error {
	return binder.Header.Binding(r.Header, obj)
}

func (self *defaultParse) ParseBody(r *http.Response, obj any) error {
	switch r.Header.Get("Content-Type") {
	case "application/json":
		raw, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		r.Body = io.NopCloser(bytes.NewReader(raw))
		return json.Unmarshal(raw, obj)
	}
	return nil
}
