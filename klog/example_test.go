// Copyright (c) 2014 Datacratic. All rights reserved.

package klog_test

import (
	"github.com/datacratic/goklog/klog"

	"fmt"
	"log"
	"os"
)

// Barebone basic usage example of klog.
func Example_Simple() {

	// Default klog printer is a light wrapper around the standard log package
	// so we need to modify it a bit to make it work for our example.
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	log.SetPrefix("klog ")

	// KPrint and KPrintf are the only supported print functions where the first
	// argument is a key identifying the source of the log message. The key will
	// be used in future example to manipulate the log lines as they flow
	// through klog.
	klog.KPrint("a.b.c", "hello")
	klog.KPrintf("c.b.a", "%s", "world")

	// Output:
	// klog <a.b.c> hello
	// klog <c.b.a> world
}

// Here we'll setup a klog pipeline which will be used to cleanup our output
// stream.
func Example_Pipeline() {

	// It's pretty easy to create a printer so, because this is go, let's log to
	// a channel.
	linesC := make(chan string, 100)
	chanPrinter := func(line *klog.Line) {
		linesC <- fmt.Sprintf("<%s> %s", line.Key, line.Value)
	}

	// Next up we'll create a filter stage for our pipeline. FilterIn indicates
	// that all matching patterns will be removed from the log stream. Filters
	// are applied on the keys and can test for a full, prefix or suffix match.
	filter := klog.NewFilter(klog.FilterOut)
	filter.AddSuffix("debug")

	// Dedup is another useful stage that can be used to avoid spamming the logs
	// with duplicated messages. Deduping is done seperately for each key-stream
	// and is periodically flushed every second by default.
	dedup := klog.NewDedup()

	// Now we'll setup our pipeline using the handy klog.Chain function. klog
	// also includes a klog.Fork function which can be used to duplicate the
	// log-stream down multiple pipelines.
	klog.SetPrinter(
		klog.Chain(filter,
			klog.Chain(dedup,
				klog.PrinterFunc(chanPrinter))))

	// Time to print stuff out.
	klog.KPrint("test.info", "hello")
	klog.KPrint("test.info", "hello")
	klog.KPrint("test.debug", "whoops")
	klog.KPrint("test.info", "hello")
	klog.KPrint("test.info", "world")

	// Alright, let's read our output.
	for i := 0; i < 3; i++ {
		fmt.Println(<-linesC)
	}

	// Output:
	// <test.info> hello
	// <test.info> hello [2 times]
	// <test.info> world
}
