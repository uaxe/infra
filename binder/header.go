package binder

import (
	"errors"
	"net/http"
	"strings"

	"github.com/uaxe/infra/zreflect"
)

var (
	_      Binder = (*header)(nil)
	Header        = &header{}
)

type header struct{}

func (h *header) Name() string {
	return "header"
}

func (h *header) Binding(src any, dst any) error {
	switch x := src.(type) {
	case http.Header:
		return zreflect.MapBindStruct(HttpHeaderMap(x), dst, h.Name())
	case map[string]any:
		return zreflect.MapBindStruct(x, dst, h.Name())
	default:
		return errors.New(`src not  http header`)
	}
}

func HttpHeaderMap(h http.Header) map[string]any {
	m := make(map[string]any, len(h))
	for k := range h {
		m[strings.ToLower(k)] = h.Get(k)
	}
	return m
}
