package utils

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strings"
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

func (self *_auths) Add(r IAuth) *_auths {
	if self.r == nil {
		self.r = make([]IAuth, 0)
	}
	self.r = append(self.r, r)
	return self
}

func (self *_auths) Authenticator(c *gin.Context, req AuthReq) (*AuthResp, error) {
	for i := range self.r {
		if strings.EqualFold(self.r[i].Name(), req.Typ) {
			return self.r[i].Authenticator(c, req)
		}
	}
	return nil, ErrFailedAuthentication
}

func ClaimsValue(c *gin.Context, key string) any {
	data := ExtractClaims(c)
	if data[key] != nil {
		return data[key]
	}
	return nil
}

func ClaimsUid(c *gin.Context) string {
	val := ClaimsValue(c, "identity")
	if val != nil {
		return val.(string)
	}
	return ""
}

func ClaimsRoleId(c *gin.Context) int {
	val := ClaimsValue(c, "roleId")
	if val != nil {
		return int(val.(float64))
	}
	return -1
}
