package util

import (
	"fmt"
	"os/exec"
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

func (t *toolHelper) RunGoCmd(path string, params ...string) (string, error)  {
	cmd := exec.Command("go", "run", path)
	cmd.Args = append(cmd.Args, params...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
