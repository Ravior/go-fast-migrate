package util

import (
	"github.com/fatih/color"
	"github.com/golang-module/carbon"
)

var LogHelper = &logHelper{}

type logHelper struct{}

func (l *logHelper) Debug(msg string, a ...interface{}) {
	color.White(carbon.Now().Format("Y/m/d H:i:s")+": "+msg, a...)
}

func (l *logHelper) Info(msg string, a ...interface{}) {
	color.Black(carbon.Now().Format("Y/m/d H:i:s")+": "+msg, a...)
}

func (l *logHelper) Warn(msg string, a ...interface{}) {
	color.Yellow(carbon.Now().Format("Y/m/d H:i:s")+": "+msg, a...)
}

func (l *logHelper) Error(msg string, a ...interface{}) {
	color.Red(carbon.Now().Format("Y/m/d H:i:s")+": "+msg, a...)
}
