package utils

import (
	"reflect"
)

func GetTypeName(t interface{}) string {
	return reflect.TypeOf(t).String()
}
