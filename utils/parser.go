package utils

import (
	"errors"
	"net/http"
	"strings"
)

type (
	Parser interface {
		Name() string
		Parse(src any, dst any) error
	}
)

var (
	_            Parser = (*headerParser)(nil)
	HeaderParser        = &headerParser{}
)

type headerParser struct{}

func (h *headerParser) Name() string {
	return "header"
}

func (h *headerParser) Parse(src any, dst any) error {
	switch x := src.(type) {
	case http.Header:
		return MapWithTagSetValue(HttpHeaderMap(x), dst, h.Name())
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
