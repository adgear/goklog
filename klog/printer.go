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

type Chainer interface {
	Printer
	Chain(next Printer)
}

func Chain(printer Chainer, next Printer) Printer {
	printer.Chain(next)
	return printer
}

func Fork(printers ...Printer) Printer {
	return PrinterFunc(func(line *Line) {
		for _, printer := range printers {
			printer.Print(line)
		}
	})
}

func Keyf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
