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

func (p *defaultIPlugin) Name() string {
	return p.name
}

func (p *defaultIPlugin) RouterPrefix() string {
	return p.prefix
}

func (p *defaultIPlugin) Register(public, private *gin.RouterGroup) {
	p.register(public, private)
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

func (p *_plugins) Add(r IPlugin) *_plugins {
	if p.r == nil {
		p.r = make([]IPlugin, 0)
	}
	p.r = append(p.r, r)
	return p
}

func (p *_plugins) Exists(name string) bool {
	return SliceFind(p.r, func(p IPlugin) bool {
		return strings.EqualFold(p.Name(), name)
	}) >= 0
}

func (p *_plugins) Register(public, private *gin.RouterGroup, skip func(IPlugin) bool) {
	SliceRange(p.r, func(p IPlugin) bool {
		if skip != nil && !skip(p) {
			prefix := p.RouterPrefix()
			p.Register(public.Group(prefix), private.Group(prefix))
		}
		return true
	})
}
