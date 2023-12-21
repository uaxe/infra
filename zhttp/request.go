package zhttp

import (
	"context"
	"net/http"
)

type Request interface {
	Do(*http.Request) (*http.Response, error)
	ContextDo(context.Context, *http.Request) (*http.Response, error)
}

var _ Request = (*DefaultRequest)(nil)

type DefaultRequest struct{}

type RequestOption func(r *DefaultRequest)

func NewRequest(opts ...RequestOption) Request {
	r := &DefaultRequest{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func (self *DefaultRequest) Do(r *http.Request) (*http.Response, error) {
	return DefaultClient.Do(r)
}

func (self *DefaultRequest) ContextDo(ctx context.Context, r *http.Request) (*http.Response, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	return DefaultClient.Do(r)
}
