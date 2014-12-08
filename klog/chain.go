// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import ()

type Chainer interface {
	Printer
	Chain(next Printer)
}

type Chained struct{ Next Printer }

func (chained *Chained) Chain(next Printer) { chained.Next = next }

func (chained *Chained) PrintNext(line *Line) {
	if chained.Next != nil {
		chained.Next.Print(line)
	}
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
