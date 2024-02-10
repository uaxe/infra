package zhttp

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/uaxe/infra/binder"
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

func (p *defaultParse) Parse(r *http.Response, obj any) error {
	if err := p.ParseBody(r, obj); err != nil {
		return err
	}
	return p.ParseHeader(r, obj)
}

func (p *defaultParse) ParseHeader(r *http.Response, obj any) error {
	return binder.Header.Binding(r.Header, obj)
}

func (p *defaultParse) ParseBody(r *http.Response, obj any) error {
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
