// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import ()

// Chainer represents a printer which forwards its output to another
// printer. Generally used to define a processing stage in a printer pipeline.
type Chainer interface {
	Printer
	Chain(next Printer)
}

// Chained implements the Chain interface and contains the boilerplate required
// to setup a printer processing stage.
type Chained struct{ Next Printer }

// Chain sets printer to be called by PrintNext.
func (chained *Chained) Chain(next Printer) { chained.Next = next }

// PrintNext forwards the line to the next printer in the pipeline.
func (chained *Chained) PrintNext(line *Line) {
	if chained.Next != nil {
		chained.Next.Print(line)
	}
}

// Chain chains the given printer to the given chained printer and returns the
// chained printer.
func Chain(printer Chainer, next Printer) Printer {
	printer.Chain(next)
	return printer
}

// Fork duplicates all received lines to multiple printers.
func Fork(printers ...Printer) Printer {
	return PrinterFunc(func(line *Line) {
		for _, printer := range printers {
			printer.Print(line)
		}
	})
}
