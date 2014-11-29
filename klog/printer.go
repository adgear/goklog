// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import (
	"fmt"
	"log"
	"time"
)

type Line struct {
	Timestamp time.Time `json:"ts"`
	Key       string    `json:"key"`
	Value     string    `json:"val"`
}

func (line *Line) String() string {
	return fmt.Sprintf("%s <%s> %s", line.Timestamp, line.Key, line.Value)
}

type lineArray []*Line

func (array lineArray) Len() int           { return len(array) }
func (array lineArray) Swap(i, j int)      { array[i], array[j] = array[j], array[i] }
func (array lineArray) Less(i, j int) bool { return array[i].Timestamp.Before(array[j].Timestamp) }

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
