package log

import "github.com/fatih/color"

var (
	infoColor = color.Cyan
	errorColor = color.Red
	debugColor = color.Yellow
)

type colored struct {
}

func (colored) Info(message string, args ...interface{}) {
	infoColor(message, args...)
}

func (colored) Error(message string, args ...interface{}) {
	errorColor(message, args...)
}

func (colored) Debug(message string, args ...interface{}) {
	debugColor(message, args...)
}

func NewColoredLogger() Logger {
	return &colored{}
}
