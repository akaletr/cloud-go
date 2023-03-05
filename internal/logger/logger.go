package logger

import "fmt"

type logger struct {
}

func New() Logger {
	return &logger{}
}

func (logger *logger) Debug(err error) {
	fmt.Println(err)
}
