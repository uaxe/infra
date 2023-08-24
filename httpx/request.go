package httpx

import (
	"context"
	"net/http"
)

type Request interface {
	Do(*http.Request) (*http.Response, error)
	ContextDo(context.Context, *http.Request) (*http.Response, error)
}

type DefaultRequest struct{}

func (self *DefaultRequest) Do(r *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(r)
}

func (self *DefaultRequest) ContextDo(ctx context.Context, r *http.Request) (*http.Response, error) {
	return http.DefaultClient.Do(r)
}
