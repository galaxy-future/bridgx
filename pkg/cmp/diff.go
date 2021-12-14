package cmp

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/galaxy-future/BridgX/pkg/utils"
	jsoniter "github.com/json-iterator/go"
)

const (
	columnTag           = "diff"
	ignoreFieldTagValue = "-"
	editOperation       = "edit"
	createOperation     = "create"
)

var errDifferenTypeCompare = errors.New("different types cannot be compared")

type DiffResult struct {
	Fields []FieldInfo
}

type FieldInfo struct {
	TagValue           string
	Name               string
	ValueNew, ValueOld reflect.Value
}

func (f FieldInfo) beautiful() (string, string) {
	return jsonDisplay(f.ValueOld), jsonDisplay(f.ValueNew)
}

func jsonDisplay(v reflect.Value) string {
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if v.Type().Name() == "Time" {
		t, ok := v.Interface().(time.Time)
		if ok {
			return utils.FormatTime(t)
		}
	}
	return fmt.Sprintf("%v", v)
}

// Diff will compare the `old` and `new` one then return diff result.Call DiffResult.Display function
// will return display map with new filed json value.
func Diff(o, n interface{}) (res DiffResult, err error) {
	defer func() {
		if r := recover(); r != nil {
			if e, ok := r.(error); ok {
				err = e
			}
			return
		}
	}()
	t1, t2 := reflect.TypeOf(o), reflect.TypeOf(n)
	if t1 != t2 {
		return DiffResult{}, errDifferenTypeCompare
	}
	v1, v2 := reflect.ValueOf(o), reflect.ValueOf(n)
	for i := 0; i < v1.NumField(); i++ {
		field := t1.Field(i)
		tagV := field.Tag.Get(columnTag)
		if tagV == ignoreFieldTagValue || !field.IsExported() {
			continue
		}
		f1Value := v1.Field(i)
		f2Value := v2.Field(i)
		if !compare(f1Value, f2Value) {
			res.Fields = append(res.Fields, FieldInfo{
				TagValue: tagV,
				ValueOld: f1Value,
				ValueNew: f2Value,
			})
		}
	}
	return res, nil
}

func (r DiffResult) Beautiful() ([]map[string]string, error) {
	output := make([]map[string]string, 0, len(r.Fields))
	if len(r.Fields) == 0 {
		return output, nil
	}
	for _, field := range r.Fields {
		operation := editOperation
		if field.ValueOld.IsZero() && !field.ValueNew.IsZero() {
			operation = createOperation
		}
		target := field.TagValue
		if target == "" {
			target = field.Name
		}

		o, n := field.beautiful()

		output = append(output, map[string]string{
			"operation": operation,
			"target":    target,
			"old":       o,
			"new":       n,
		})
	}
	return output, nil
}

func compare(v1, v2 reflect.Value) bool {
	k1, k2 := v1.Kind(), v2.Kind()
	if k1 != k2 {
		return false
	}
	switch k1 {
	case reflect.Bool:
		return v1.Bool() == v2.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v1.Int() == v2.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v1.Uint() == v2.Uint()
	case reflect.Float32, reflect.Float64:
		return v1.Float() == v2.Float()
	case reflect.Complex64, reflect.Complex128:
		return v1.Complex() == v2.Complex()
	case reflect.String:
		return v1.String() == v2.String()
	case reflect.Slice, reflect.Array, reflect.Map:
		b1, _ := jsoniter.Marshal(v1.Interface())
		b2, _ := jsoniter.Marshal(v2.Interface())
		return reflect.DeepEqual(b1, b2)
	case reflect.Ptr, reflect.Interface:
		return reflect.DeepEqual(v1.Elem(), v2.Elem())
	case reflect.Func, reflect.Struct, reflect.Chan, reflect.UnsafePointer:
		return true

	default:
		return true
	}
}
