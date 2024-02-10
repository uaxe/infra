package zhttp

import (
	"net/http"
)

var DefaultClient = &http.Client{
	Transport: Transport(false, nil),
}
