//go:build !solution

package structtags

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)

type fieldInfo struct {
	name  string
	field reflect.StructField
}

func prepareFieldMap(typ reflect.Type) map[string]fieldInfo {
	fields := make(map[string]fieldInfo, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		tag := field.Tag.Get("http")
		name := tag
		if name == "" {
			name = strings.ToLower(field.Name)
		}
		fields[name] = fieldInfo{name: name, field: field}
	}
	return fields
}

func Unpack(req *http.Request, ptr interface{}) error {
	if err := req.ParseForm(); err != nil {
		return err
	}
	v := reflect.ValueOf(ptr).Elem()
	typ := v.Type()
	fieldMap := prepareFieldMap(typ)
	for name, values := range req.Form {
		fieldInfo, ok := fieldMap[name]
		if !ok {
			continue
		}
		fieldValue := v.FieldByName(fieldInfo.field.Name)
		for _, value := range values {
			if fieldValue.Kind() == reflect.Slice {
				elem := reflect.New(fieldValue.Type().Elem()).Elem()
				if err := populate(elem, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
				fieldValue.Set(reflect.Append(fieldValue, elem))
			} else {
				if err := populate(fieldValue, value); err != nil {
					return fmt.Errorf("%s: %v", name, err)
				}
			}
		}
	}
	return nil
}

func populate(v reflect.Value, value string) error {
	switch v.Kind() {
	case reflect.String:
		v.SetString(value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)

	case reflect.Bool:
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		v.SetBool(b)

	default:
		return fmt.Errorf("unsupported kind: %s", v.Type().Kind())
	}
	return nil
}
