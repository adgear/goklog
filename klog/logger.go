// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import (
	"fmt"
	"time"
)

// Logger is used to inject newly created lines into a printer pipeline.
type Logger struct{ Chained }

// New creates a new Logger which outputs to the given printer.
func New(next Printer) *Logger {
	return &Logger{Chained: Chained{Next: next}}
}

// KPrint is similar to fmt.Print but accepts a key as it's first parameter.
func (logger *Logger) KPrint(key string, v ...interface{}) {
	line := &Line{time.Now(), key, fmt.Sprint(v...)}
	logger.PrintNext(line)
}

// KPrintf is similar to fmt.Printf but accepts a key as it's first parameter.
func (logger *Logger) KPrintf(key, format string, v ...interface{}) {
	line := &Line{time.Now(), key, fmt.Sprintf(format, v...)}
	logger.PrintNext(line)
}

// logger is the default printer used by the global klog print functions.
var logger = New(DefaultPrinter)

// KPrint is similar to fmt.Print but accepts a key as it's first parameter.
func KPrint(key string, v ...interface{}) { logger.KPrint(key, v...) }

// KPrintf is similar to fmt.Printf but accepts a key as it's first parameter.
func KPrintf(key, format string, v ...interface{}) { logger.KPrintf(key, format, v...) }

// GetPrinter returns the global printer used by the global KPrint and KPrintf
// function.
func GetPrinter() Printer { return logger.Next }

// SetPrinter changes the global printer used by the global KPrint and KPrintf
// function.
func SetPrinter(next Printer) { logger.Chain(next) }
