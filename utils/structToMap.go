package utils

import (
	"reflect"
	"strings"
)

func StructToMap(data interface{}) map[string]interface{} {

	fields := make(map[string]interface{})

	v := reflect.ValueOf(data)
	size := v.NumField()

	for i := 0; i < size; i++ {
		field := v.Field(i)

		if !field.CanInterface() {
			continue
		}

		value := field.Interface()
		key := strings.ToLower(v.Type().Field(i).Name)

		if !reflect.DeepEqual(value, reflect.Zero(field.Type()).Interface()) {
			fields[key] = value
		}
	}

	return fields
}