// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import (
	"testing"
	"time"
)

func TestDedup(t *testing.T) {
	out := &TestPrinter{T: t}
	dedup := Dedup{Rate: 10 * time.Millisecond}
	dedup.Chain(out)

	dedup.Print(L("a", "x"))
	dedup.Print(L("a", "x"))
	dedup.Print(L("a", "x"))
	dedup.Print(L("a", "y"))

	dedup.Print(L("b", "y"))
	dedup.Print(L("b", "x"))
	dedup.Print(L("b", "x"))
	dedup.Print(L("b", "x"))

	dedup.Print(L("c", "x"))
	dedup.Print(L("c", "y"))
	dedup.Print(L("c", "x"))
	dedup.Print(L("c", "y"))

	dedup.Print(L("d", "x"))
	dedup.Print(L("d", "y"))
	dedup.Print(L("d", "y"))
	dedup.Print(L("d", "x"))

	dedup.Print(L("e", "x"))
	dedup.Print(L("f", "y"))
	dedup.Print(L("e", "y"))
	dedup.Print(L("f", "y"))
	dedup.Print(L("e", "y"))
	dedup.Print(L("f", "y"))
	dedup.Print(L("e", "y"))
	dedup.Print(L("f", "x"))

	time.Sleep(10 * time.Millisecond)

	dedup.Print(L("a", "y"))
	dedup.Print(L("a", "y"))
	dedup.Print(L("a", "y"))
	dedup.Print(L("a", "x"))

	dedup.Print(L("b", "x"))
	dedup.Print(L("b", "x"))
	dedup.Print(L("b", "x"))
	dedup.Print(L("b", "y"))

	dedup.Print(L("e", "x"))
	dedup.Print(L("e", "y"))
	dedup.Print(L("e", "y"))
	dedup.Print(L("e", "y"))

	time.Sleep(10 * time.Millisecond)

	out.ExpectOrdered(
		"<a> x",
		"<a> x [2 times]",
		"<a> y",

		"<b> y",
		"<b> x",

		"<c> x",
		"<c> y",
		"<c> x",
		"<c> y",

		"<d> x",
		"<d> y",
		"<d> y",
		"<d> x",

		"<e> x",
		"<f> y",
		"<e> y",
		"<f> y [2 times]",
		"<f> x",
	)

	out.ExpectUnordered(
		"<b> x [2 times]",
		"<e> y [2 times]",
	)

	out.ExpectOrdered(
		"<a> y [3 times]",
		"<a> x",

		"<b> x [3 times]",
		"<b> y",

		"<e> x",
		"<e> y",
		"<e> y [2 times]",
	)
}

func BenchmarkDedup_Const(b *testing.B) {
	dedup := Dedup{Rate: 10 * time.Millisecond}
	l := L("a", "x")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dedup.Print(l)
	}
}

func BenchmarkDedup_NoDup(b *testing.B) {
	dedup := Dedup{Rate: 10 * time.Millisecond}
	l0, l1 := L("a", "x"), L("a", "y")

	b.ResetTimer()
	for i := 0; i < b.N/2; i++ {
		dedup.Print(l0)
		dedup.Print(l1)
	}
}
