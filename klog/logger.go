// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import (
	"fmt"
	"os"
	"time"
)

// Logger is used to inject newly created lines into a printer pipeline.
type Logger struct {
	Chained
	Fatal Printer
}

// New creates a new Logger which outputs to the given printer.
func New(next Printer, fatal Printer) *Logger {
	return &Logger{Chained: Chained{Next: next}}
}

func (logger *Logger) kprint(key, value string) {
	logger.PrintNext(&Line{time.Now(), key, value})
}

// KPrint is similar to log.Print but accepts a key as it's first parameter.
func (logger *Logger) KPrint(key string, v ...interface{}) {
	logger.kprint(key, fmt.Sprint(v...))
}

// KPrintf is similar to log.Printf but accepts a key as it's first parameter.
func (logger *Logger) KPrintf(key, format string, v ...interface{}) {
	logger.kprint(key, fmt.Sprintf(format, v...))
}

func (logger *Logger) kfatal(key, value string) {
	line := &Line{time.Now(), key, value}
	logger.Fatal.Print(line)
	os.Exit(1)
}

// KFatal is similar to log.Fatal but accepts a key as it's first parameter.
func (logger *Logger) KFatal(key string, v ...interface{}) {
	logger.kfatal(key, fmt.Sprint(v...))
}

// KFatalf is similar to log.Fatalf but accepts a key as it's first parameter.
func (logger *Logger) KFatalf(key, format string, v ...interface{}) {
	logger.kfatal(key, fmt.Sprintf(format, v...))
}

func (logger *Logger) kpanic(key, value string) {
	line := &Line{time.Now(), key, value}
	logger.Fatal.Print(line)
	panic(line.String())
}

// KPanic is similar to log.Panic but accepts a key as it's first parameter.
func (logger *Logger) KPanic(key string, v ...interface{}) {
	logger.kpanic(key, fmt.Sprint(v...))
}

// KPanicf is similar to log.Panicf but accepts a key as it's first parameter.
func (logger *Logger) KPanicf(key, format string, v ...interface{}) {
	logger.kpanic(key, fmt.Sprintf(format, v...))
}

// logger is the default printer used by the global klog print functions.
var logger = New(DefaultPrinter, DefaultPrinter)

// KPrint is similar to fmt.Print but accepts a key as it's first parameter.
func KPrint(key string, v ...interface{}) { logger.KPrint(key, v...) }

// KPrintf is similar to fmt.Printf but accepts a key as it's first parameter.
func KPrintf(key, format string, v ...interface{}) { logger.KPrintf(key, format, v...) }

// KFatal is similar to fmt.Fatal but accepts a key as it's first parameter.
func KFatal(key string, v ...interface{}) { logger.KFatal(key, v...) }

// KFatalf is similar to fmt.Fatalf but accepts a key as it's first parameter.
func KFatalf(key, format string, v ...interface{}) { logger.KFatalf(key, format, v...) }

// KPanic is similar to fmt.Panic but accepts a key as it's first parameter.
func KPanic(key string, v ...interface{}) { logger.KPanic(key, v...) }

// KPanicf is similar to fmt.Panicf but accepts a key as it's first parameter.
func KPanicf(key, format string, v ...interface{}) { logger.KPanicf(key, format, v...) }

// GetPrinter returns the global printer used by the global KPrint and KPrintf
// function.
func GetPrinter() Printer { return logger.Next }

// SetPrinter changes the global printer used by the global KPrint and KPrintf
// function.
func SetPrinter(next Printer) { logger.Chain(next) }

// SetFatalPrinter changes the global printer used by the global KFatal and
// KPanic funcions. Since the program is about to go down after calling these
// functions, the printer should be short and sweet and not defer work to a
// background goroutine.
func SetFatalPrinter(fatal Printer) { logger.Fatal = fatal }
