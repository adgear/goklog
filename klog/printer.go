// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

// Printer represents a stage in the printing pipeline.
type Printer interface {
	Print(*Line)
}

// PrinterFunc implements the Printer interface for functions.
type PrinterFunc func(*Line)

// Print passes the line to the function.
func (fn PrinterFunc) Print(line *Line) { fn(line) }

// NilPrinter is a noop printer.
var NilPrinter = PrinterFunc(func(line *Line) {})

// DefaultPrinter is the default printer used by the kprint functions in klog.
var DefaultPrinter = PrinterFunc(LogPrinter)

// DefaultPrinter is the default printer used by the kfatal and kpanic functions
// in klog.
var DefaultFatalPrinter = PrinterFunc(LogPrinter)

// LogPrinter is forwards all lines to the golang standard log library.
func LogPrinter(line *Line) { log.Printf("<%s> %s", line.Key, line.Value) }

// Keyf is a utility formatting functions for key and is a light wrapper around
// fmt.Sprintf.
func Keyf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

// JsonPrinter forwards all lines to the golang standard log library
// in a json format
func JsonPrinter(line *Line) {
	split := strings.Split(line.Key, ".")
	var level string
	if len(split) > 0 {
		level = split[len(split)-1]
	}
	line.Key = strings.Join(split[:len(split)-1], ".")

	sLine := struct {
		*Line
		Level string `json:"level"`
	}{
		Line:  line,
		Level: level,
	}

	if js, err := json.Marshal(sLine); err != nil {
		log.Printf("line json marshal error: %s", err)
	} else {
		log.Printf("@cee: %s", js)
	}
}

// Structured printer if a JsonPrinter
var StructuredPrinter = PrinterFunc(JsonPrinter)
