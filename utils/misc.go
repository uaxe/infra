package utils

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	_ = 1.0 << (10 * iota)
	KB
	MB
	GB
	TB
	PB
	EB
)

type Struct2MapOptionFunc func(v reflect.Value, f reflect.StructField) (string, bool)

func OptionSkipTags(tags ...string) Struct2MapOptionFunc {
	m := Slice2MapStruct(tags)
	return func(v reflect.Value, f reflect.StructField) (string, bool) {
		j := f.Tag.Get("json")
		g := f.Tag.Get("gorm")
		if j == "-" || g == "-" {
			return "", true
		}
		if !strings.EqualFold(j, g) {
			if strings.HasPrefix(g, "column:") {
				j = strings.TrimPrefix(g, "column:")
			}
		}
		if j == "" {
			return j, true
		}
		if _, ok := m[j]; ok {
			return j, true
		}
		return j, false
	}
}

func Struct2Map(x any, opts ...Struct2MapOptionFunc) (map[string]any, error) {
	t := reflect.TypeOf(x)
	v := reflect.ValueOf(x)
	if t.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	if v.Kind() != reflect.Struct {
		return nil, errors.New("struct not found!")
	}
	if len(opts) == 0 {
		opts = append(opts, OptionSkipTags())
	}
	num := v.NumField()
	m := make(map[string]any, num)
loop:
	for i := 0; i < num; i++ {
		f := v.Field(i)
		if !f.CanInterface() {
			continue
		}
		for _, opt := range opts {
			if tag, ok := opt(f, t.Field(i)); ok {
				continue loop
			} else {
				m[tag] = f.Interface()
			}
		}
	}
	return m, nil
}

type MergeStructOptionFunc func(xf, yf reflect.StructField, xv, yv reflect.Value) (bool, error)

func MergerStructIsZero() MergeStructOptionFunc {
	return func(xt, yt reflect.StructField, xv, yv reflect.Value) (bool, error) {
		if xv.FieldByName(xt.Name).IsZero() && !yv.FieldByName(xt.Name).IsZero() {
			xv.FieldByName(xt.Name).Set(yv.FieldByName(yt.Name))
			return true, nil
		}
		return false, nil
	}
}

func MergeStruct(x, y any, opts ...MergeStructOptionFunc) error {
	xt := reflect.TypeOf(x)
	xv := reflect.ValueOf(x)
	if xt.Kind() == reflect.Ptr {
		xv = xv.Elem()
		xt = xt.Elem()
	}

	yt := reflect.TypeOf(y)
	yv := reflect.ValueOf(y)
	if yt.Kind() == reflect.Ptr {
		yv = yv.Elem()
		yt = yt.Elem()
	}

	if xv.NumField() != yv.NumField() {
		return fmt.Errorf("not same num field %d,%d", xv.NumField(), yv.NumField())
	}
	if len(opts) == 0 {
		opts = append(opts, MergerStructIsZero())
	}
loop:
	for i := 0; i < xv.NumField(); i++ {
		if !xv.Field(i).CanInterface() || !yv.Field(i).CanInterface() {
			continue
		}
		if !strings.EqualFold(xt.Field(i).Name, yt.Field(i).Name) {
			return fmt.Errorf("name diff %s,%s", xt.Field(i).Name, yt.Field(i).Name)
		}
		for _, opt := range opts {
			if ok, err := opt(xt.Field(i), yt.Field(i), xv, yv); err != nil {
				return err
			} else if ok {
				continue loop
			}
		}
	}

	return nil
}
