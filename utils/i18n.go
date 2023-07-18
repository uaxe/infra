package utils

import (
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"strings"
	"sync"

	"github.com/uaxe/infra/core/conf"
)

const (
	LANG_DEFAULT = "zh-CN"
	ACCEPT_LANG  = "Accept-Language"
)

type I18n interface {
	HttpValue(req *http.Request, key, defaultVal string, params ...string) string
	LangValue(lang, key, defaultVal string, params ...string) string
	Append(i18nDir string) error
}

var _ I18n = (*defaultI18n)(nil)

var (
	locker sync.RWMutex
)

type defaultI18n struct {
	i18nVals map[string]string
}

func InitI18n(i18nDir string) (I18n, error) {
	i18n := &defaultI18n{i18nVals: make(map[string]string)}
	err := loadI18n(i18n.i18nVals, i18nDir)
	if err != nil {
		return nil, err
	}
	return i18n, nil
}

func HttpLanguage(req *http.Request) string {
	lang := req.Header.Get(ACCEPT_LANG)
	if len(lang) <= 0 {
		lang = LANG_DEFAULT
	}
	lang = lang[:strings.Index(lang, ",")]
	return lang
}

func langKey(lang, key string) string {
	return fmt.Sprintf("%s.%s", lang, strings.ToLower(key))
}

func (self *defaultI18n) HttpValue(req *http.Request, key, defaultVal string, params ...string) string {
	return self.LangValue(HttpLanguage(req), key, defaultVal, params...)
}

func (self *defaultI18n) LangValue(lang, key, defaultVal string, params ...string) string {
	locker.RLock()
	defer locker.RUnlock()

	lkey := langKey(lang, key)
	val, ok := self.i18nVals[lkey]
	paramsVal := make([]string, 0, len(params))
	for _, v := range params {
		pval, pok := self.i18nVals[lkey]
		if pok {
			paramsVal = append(paramsVal, langKey(lang, pval))
		} else {
			paramsVal = append(paramsVal, v)
		}
	}
	if ok {
		if len(paramsVal) > 0 {
			return fmt.Sprintf(val, paramsVal)
		}
		return val
	}
	return defaultVal
}

func (self *defaultI18n) Append(i18nDir string) error {
	if self.i18nVals == nil {
		self.i18nVals = make(map[string]string)
	}
	return loadI18n(self.i18nVals, i18nDir)
}

func loadI18n(i18nVals map[string]string, i18nDir string) error {
	locker.Lock()
	defer locker.Unlock()

	return filepath.WalkDir(i18nDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			info, er := d.Info()
			if er != nil {
				return er
			}
			fiPath := fmt.Sprintf("%s/%s", i18nDir, info.Name())
			values := make(map[string]string)
			err = conf.Load(fiPath, &values)
			if err != nil {
				return err
			}
			lang := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
			for k, v := range values {
				i18nVals[langKey(lang, k)] = v
			}
		}
		return nil
	})
}
