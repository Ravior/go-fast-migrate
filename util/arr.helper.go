package util

import "strings"

var ArrHelper = &arrHelper{}

type arrHelper struct{}

func (a *arrHelper) Get(arr []interface{}, index int) interface{} {
	if len(arr) >= index+1 {
		return arr[index]
	}
	return nil
}

// Check if slice contain a string
func (a *arrHelper) StrArrContain(s []string, str string) bool {
	for _, v := range s {
		if strings.Compare(v, str) == 0 {
			return true
		}
	}
	return false
}
