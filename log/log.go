package logger

import (
	"fmt"
	"time"
)

type Logger struct {
	label string
}

func New(label string) *Logger {
	return &Logger{
		label: label,
	}
}

func (l *Logger) Fatal(args ...interface{}) {
	fmt.Println(args...)
	panic(args)
}

func (l *Logger) getPrefix() string {
	return "[" + l.label + "] [" + time.Now().Format(time.RFC3339) + "]"
}

func (l *Logger) Info(args ...interface{}) {
	newArgs := make([]any, len(args)+1)
	newArgs[0] = l.getPrefix()
	for i, arg := range args {
		newArgs[i+1] = arg
	}
	fmt.Println(newArgs...)
}
