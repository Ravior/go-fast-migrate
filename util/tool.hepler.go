package util

import (
	"fmt"
	"reflect"
)

var ToolHelper = &toolHelper{}

type toolHelper struct{}

func (t *toolHelper) GetStrParam(key interface{}, _default string) string {
	if key == nil {
		return _default
	}

	if reflect.TypeOf(key).Kind() == reflect.String {
		return fmt.Sprintf("%v", key)
	}

	return ""
}
