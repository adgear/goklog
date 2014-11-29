// Copyright (c) 2014 Datacratic. All rights reserved.

package klogr_test

import (
	"github.com/datacratic/goklog/klog"
	"github.com/datacratic/goklog/klog/rest"

	"fmt"
)

// The klogr package provides wrapppers for printer stages which makes them
// accessible from a rest interface using the gorest package.
func Example_Rest() {

	// Here we create a REST enabled filter which allows us to modify the
	// filtering rules of the stage while the program is running.
	filter := klogr.NewRestFilter("/path/to/filter", klog.FilterIn)

	// The ring stage prints to a ring buffer containing the last N elements
	// printed. Wrapping this stage in a REST interface allows us to query the
	// most recent log elements remotely.
	ring := klogr.NewRestRing("/path/to/ring", 1000)

	// REST enabled stages act like regular stages so they can be used directly
	// in a klog pipeline.
	klog.SetPrinter(
		klog.Chain(filter,
			klog.Fork(ring, klog.GetPrinter())))

	klog.KPrint("a.b.c", "hello world")
}
