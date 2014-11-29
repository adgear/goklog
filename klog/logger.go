// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import (
	"fmt"
	"time"
)

type Logger struct{ Next Printer }

func New(next Printer) *Logger {
	return &Logger{Next: next}
}

func (logger *Logger) KPrint(key string, v ...interface{}) {
	line := &Line{time.Now(), key, fmt.Sprint(v...)}
	logger.Next.Print(line)
}

func (logger *Logger) KPrintf(key, format string, v ...interface{}) {
	line := &Line{time.Now(), key, fmt.Sprintf(format, v...)}
	logger.Next.Print(line)
}

func (logger *Logger) GetPrinter() Printer     { return logger.Next }
func (logger *Logger) SetPrinter(next Printer) { logger.Next = next }

var logger = New(DefaultPrinter)

func KPrint(key string, v ...interface{})          { logger.KPrint(key, v...) }
func KPrintf(key, format string, v ...interface{}) { logger.KPrintf(key, format, v...) }

func GetPrinter() Printer     { return logger.GetPrinter() }
func SetPrinter(next Printer) { logger.SetPrinter(next) }
