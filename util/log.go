package util

import (
	"fmt"
)

type TinyLogger struct {
	isVerboseEnabled bool
}

func NewTinyLogger(isVerboseEnabled bool) *TinyLogger {
	return &TinyLogger{isVerboseEnabled}
}

func (l *TinyLogger) Verboseln(v ...interface{}) {
	if l.isVerboseEnabled {
		v = append([]interface{}{"[verbose]"}, v...)
		fmt.Println(v...)
	}
}
func (l *TinyLogger) Verbosef(format string, v ...interface{}) {
	if l.isVerboseEnabled {
		fmt.Printf("[verbose] "+format, v...)
	}
}

func (l *TinyLogger) Writeln(v ...interface{}) {
	fmt.Println(v...)
}

func (l *TinyLogger) Writef(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}
