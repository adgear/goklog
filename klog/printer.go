// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import (
	"fmt"
	"log"
)

type Printer interface {
	Print(*Line)
}

type PrinterFunc func(*Line)

func (fn PrinterFunc) Print(line *Line) { fn(line) }

var NilPrinter = PrinterFunc(func(line *Line) {})

var DefaultPrinter = PrinterFunc(LogPrinter)

func LogPrinter(line *Line) { log.Printf("<%s> %s", line.Key, line.Value) }

func Keyf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
