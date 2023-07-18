package utils

import (
	"github.com/gin-gonic/gin"
)

type (
	IRouter interface {
		Register(public, private *gin.RouterGroup)
	}

	_router struct {
		r []IRouter
	}
)

func GetRouter() *_router {
	return _defaultRouter
}

var _defaultRouter = new(_router)

func (self *_router) Add(r IRouter) *_router {
	if self.r == nil {
		self.r = make([]IRouter, 0)
	}
	self.r = append(self.r, r)
	return self
}

func (self *_router) Register(public, private *gin.RouterGroup) {
	for i := range self.r {
		self.r[i].Register(public, private)
	}
}
