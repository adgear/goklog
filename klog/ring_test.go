// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import (
	"strconv"
	"testing"
)

func TestRing_Get(t *testing.T) {
	ring := NewRing(10)

	ExpectOrdered(t, Simplify(ring.GetAll()))

	ring.Print(L("a", "x"))
	ring.Print(L("a.b.c", "x"))
	ring.Print(L("a", "y"))
	ring.Print(L("b", "x"))
	ring.Print(L("a", "z"))

	ExpectOrdered(t, Simplify(ring.GetAll()),
		"<a> x",
		"<a.b.c> x",
		"<a> y",
		"<b> x",
		"<a> z",
	)

	ExpectOrdered(t, Simplify(ring.GetKey("a")),
		"<a> x",
		"<a> y",
		"<a> z",
	)

	ExpectOrdered(t, Simplify(ring.GetPrefix("a")),
		"<a> x",
		"<a.b.c> x",
		"<a> y",
		"<a> z",
	)
}

func TestRing_Rollover(t *testing.T) {
	ring := NewRing(4)

	for i := 0; i < 6; i++ {
		ring.Print(L("a", strconv.Itoa(i)))
	}

	ExpectOrdered(t, Simplify(ring.GetAll()),
		"<a> 2",
		"<a> 3",
		"<a> 4",
		"<a> 5",
	)
}

func BenchmarkRing(b *testing.B) {
	ring := NewRing(100)
	l := L("a", "x")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			ring.Print(l)
		}
	})
}
