package zreflect

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

type Struct2MapOption func(v reflect.Value, f reflect.StructField) (string, bool)

func Struct2MapOptionWithTag(tag string) Struct2MapOption {
	return func(v reflect.Value, f reflect.StructField) (string, bool) {
		j, ok := f.Tag.Lookup(tag)
		if !ok || j == "-" {
			return j, false
		}
		return j, true
	}
}

func Struct2Map(x any, opts ...Struct2MapOption) (map[string]any, error) {
	t, v := ReflectTypeAndValue(x)
	if v.Kind() != reflect.Struct {
		return nil, errors.New("struct not found!")
	}
	if len(opts) == 0 {
		opts = append(opts, Struct2MapOptionWithTag("json"))
	}
	num := v.NumField()
	m := make(map[string]any, num)
	for i := 0; i < num; i++ {
		f := v.Field(i)
		if !f.CanInterface() {
			continue
		}
		for _, opt := range opts {
			if tag, ok := opt(f, t.Field(i)); ok {
				m[tag] = f.Interface()
				break
			}
		}
	}
	return m, nil
}

type MergeStructOption func(xf, yf reflect.StructField, xv, yv reflect.Value) error

func MergerStructIsZero() MergeStructOption {
	return func(xt, yt reflect.StructField, xv, yv reflect.Value) error {
		if xv.FieldByName(xt.Name).IsZero() && !yv.FieldByName(xt.Name).IsZero() {
			xv.FieldByName(xt.Name).Set(yv.FieldByName(yt.Name))
		}
		return nil
	}
}

func MergeSameStruct(x, y any, opts ...MergeStructOption) error {
	xt, xv := ReflectTypeAndValue(x)
	yt, yv := ReflectTypeAndValue(y)

	if xv.NumField() != yv.NumField() {
		return fmt.Errorf("not same num field %d,%d", xv.NumField(), yv.NumField())
	}

	if len(opts) == 0 {
		opts = append(opts, MergerStructIsZero())
	}

	for i := 0; i < xv.NumField(); i++ {
		if !xv.Field(i).CanInterface() || !yv.Field(i).CanInterface() {
			continue
		}
		if !strings.EqualFold(xt.Field(i).Name, yt.Field(i).Name) {
			return fmt.Errorf("name diff %s,%s", xt.Field(i).Name, yt.Field(i).Name)
		}
		for _, opt := range opts {
			if err := opt(xt.Field(i), yt.Field(i), xv, yv); err != nil {
				return err
			}
		}
	}

	return nil
}

func ReflectTypeAndValue(src any) (reflect.Type, reflect.Value) {
	xt, xv := reflect.TypeOf(src), reflect.ValueOf(src)
	switch xt.Kind() {
	case reflect.Ptr:
		xv, xt = xv.Elem(), xt.Elem()
	default:
	}
	return xt, xv
}

func StructWithTagSetValue(src, dst any, tag string) error {
	xt, xv := ReflectTypeAndValue(src)
	if xv.Kind() != reflect.Struct {
		return errors.New("src not struct")
	}
	yt, yv := ReflectTypeAndValue(dst)
	if yv.Kind() != reflect.Struct {
		return errors.New("dst not struct")
	}
	for i := 0; i < yv.NumField(); i++ {
		if !yv.Field(i).CanInterface() {
			continue
		}
		if j, ok := yt.Field(i).Tag.Lookup(tag); !ok || j == "-" {
			continue
		}
		if field, exists := xt.FieldByName(yt.Field(i).Name); exists {
			yv.FieldByName(field.Name).Set(xv.FieldByName(field.Name))
		}
	}
	return nil
}

func MapWithTagSetValue(src map[string]any, dst any, tag string) error {
	yt, yv := ReflectTypeAndValue(dst)
	if yv.Kind() != reflect.Struct {
		return errors.New("dst not struct")
	}
	for i := 0; i < yv.NumField(); i++ {
		if !yv.Field(i).CanInterface() {
			continue
		}
		if j, ok := yt.Field(i).Tag.Lookup(tag); ok && j != "-" {
			if field, exists := src[strings.ToLower(j)]; exists {
				yv.FieldByName(yt.Field(i).Name).Set(reflect.ValueOf(field))
			}
		}
	}
	return nil
}
