package util

import (
	"github.com/fatih/color"
	"os"
)

var SysHelper = &sysHelper{}

type sysHelper struct{}

func (s *sysHelper) Exit(msg string, a ...interface{}) {
	color.Red(msg, a...)
	os.Exit(2)
}

func (s *sysHelper) CheckErr(err error) {
	if err != nil {
		s.Exit("Error: %v", err)
	}
}
