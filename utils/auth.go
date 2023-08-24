package utils

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
)

type (
	IAuth interface {
		Name() string
		Authenticator(*gin.Context, AuthReq) (*AuthResp, error)
	}

	_auths struct {
		r []IAuth
	}

	AuthReq struct {
		Typ  string          `json:"typ"`
		Data json.RawMessage `json:"data"`
	}

	AuthResp struct {
		Typ  string  `json:"typ"`
		User SysUser `json:"user"`
	}

	SysUser struct {
		UserId string `json:"userId"`
		RoleId int64  `json:"roleId"`
	}
)

func GetAuth() *_auths {
	return _defaultAuths
}

var _defaultAuths = new(_auths)

func (auth *_auths) Add(r IAuth) *_auths {
	if auth.r == nil {
		auth.r = make([]IAuth, 0)
	}
	auth.r = append(auth.r, r)
	return auth
}

var (
	ErrFailedAuthentication = errors.New("incorrect Username or Password")
)

func (auth *_auths) Authenticator(c *gin.Context, req AuthReq) (*AuthResp, error) {
	for i := range auth.r {
		if strings.EqualFold(auth.r[i].Name(), req.Typ) {
			return auth.r[i].Authenticator(c, req)
		}
	}
	return nil, ErrFailedAuthentication
}
