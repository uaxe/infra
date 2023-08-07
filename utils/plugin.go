package utils

import (
	"github.com/gin-gonic/gin"
	"strings"
)

type (
	IPlugin interface {
		Name() string
		RouterPrefix() string
		Register(public, private *gin.RouterGroup)
	}

	_plugins struct {
		r []IPlugin
	}

	defaultIPlugin struct {
		name, prefix string
		register     func(public, private *gin.RouterGroup)
	}
)

func (self *defaultIPlugin) Name() string {
	return self.name
}

func (self *defaultIPlugin) RouterPrefix() string {
	return self.prefix
}

func (self *defaultIPlugin) Register(public, private *gin.RouterGroup) {
	self.register(public, private)
}

func GetPlugin() *_plugins {
	return _defaultPlugin
}

func AddPlugin(name, prefix string, register func(public, private *gin.RouterGroup)) {
	i := defaultIPlugin{
		name:     name,
		prefix:   prefix,
		register: register,
	}
	_defaultPlugin.Add(&i)
}

var _defaultPlugin = new(_plugins)

func (self *_plugins) Add(r IPlugin) *_plugins {
	if self.r == nil {
		self.r = make([]IPlugin, 0)
	}
	self.r = append(self.r, r)
	return self
}

func (self *_plugins) Exists(name string) bool {
	return SliceFind(self.r, func(p IPlugin) bool {
		return strings.EqualFold(p.Name(), name)
	}) >= 0
}

func (self *_plugins) Register(public, private *gin.RouterGroup, skip func(IPlugin) bool) {
	SliceRange(self.r, func(p IPlugin) bool {
		if skip != nil && !skip(p) {
			prefix := p.RouterPrefix()
			p.Register(public.Group(prefix), private.Group(prefix))
		}
		return true
	})
}
