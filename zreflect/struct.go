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
	t, v := TypeAndValue(x)
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
	return func(x, y reflect.StructField, xv, yv reflect.Value) error {
		val1, val2 := xv.FieldByName(x.Name), yv.FieldByName(y.Name)
		if !val1.IsZero() && val2.IsZero() {
			if !val2.CanAddr() {
				return fmt.Errorf("dst field [%s] not can addr", y.Name)
			}
			val2.Set(val1)
		}
		return nil
	}
}

func StructWithTag(src, dst any, tag string) error {
	xt, xv := TypeAndValue(src)
	if xv.Kind() != reflect.Struct {
		return errors.New("src not struct")
	}
	yt, yv := TypeAndValue(dst)
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

func MapBindStruct(src map[string]any, dst any, tag string) error {
	yt, yv := TypeAndValue(dst)
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

func MergeStruct(src, dst any, opts ...MergeStructOption) error {
	xt, xv := TypeAndValue(src)
	yt, yv := TypeAndValue(dst)

	if len(opts) == 0 {
		opts = append(opts, MergerStructIsZero())
	}
	xlen, ylen := xv.NumField(), yv.NumField()
	for i := 0; i < xlen; i++ {
		if i > ylen || !xv.Field(i).CanInterface() || !yv.Field(i).CanInterface() {
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
