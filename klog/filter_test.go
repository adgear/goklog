// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import (
	"strconv"
	"testing"
)

func DoFilterPrints(printer Printer) {
	printer.Print(L("a", "x"))
	printer.Print(L("a.b", "x"))
	printer.Print(L("a.c", "x"))
	printer.Print(L("a.b.c", "x"))
	printer.Print(L("a.c.b", "x"))

	printer.Print(L("b", "x"))
	printer.Print(L("b.a", "x"))
	printer.Print(L("b.c", "x"))
	printer.Print(L("b.a.c", "x"))
	printer.Print(L("b.c.a", "x"))

	printer.Print(L("c", "x"))
	printer.Print(L("c.a", "x"))
	printer.Print(L("c.b", "x"))
	printer.Print(L("c.a.b", "x"))
	printer.Print(L("c.b.a", "x"))
}

func TestFilter_Out(t *testing.T) {
	out := &TestPrinter{T: t}
	filter := &Filter{Type: FilterOut, Next: out}

	DoFilterPrints(filter)
	out.ExpectOrdered(
		"<a> x",
		"<a.b> x",
		"<a.c> x",
		"<a.b.c> x",
		"<a.c.b> x",

		"<b> x",
		"<b.a> x",
		"<b.c> x",
		"<b.a.c> x",
		"<b.c.a> x",

		"<c> x",
		"<c.a> x",
		"<c.b> x",
		"<c.a.b> x",
		"<c.b.a> x",
	)

	filter.Add("a.b.c")
	filter.Add("c.b.a")
	filter.Add("d")

	filter.AddPrefix("b")
	filter.AddPrefix("d")

	filter.AddSuffix("c")
	filter.AddSuffix("d")

	WaitForPropagation()

	DoFilterPrints(filter)
	out.ExpectOrdered(
		"<a> x",
		"<a.b> x",
		"<a.c.b> x",

		"<c.a> x",
		"<c.b> x",
		"<c.a.b> x",
	)

	filter.Remove("a.b.c")
	filter.Remove("c.b.a")
	filter.Remove("e")

	filter.AddPrefix("c")
	filter.RemovePrefix("b")
	filter.RemovePrefix("e")

	filter.AddSuffix("b")
	filter.RemoveSuffix("c")
	filter.RemoveSuffix("e")

	WaitForPropagation()

	DoFilterPrints(filter)
	out.ExpectOrdered(
		"<a> x",
		"<a.c> x",
		"<a.b.c> x",

		"<b.a> x",
		"<b.c> x",
		"<b.a.c> x",
		"<b.c.a> x",
	)
}

func TestFilter_In(t *testing.T) {
	out := &TestPrinter{T: t}
	filter := &Filter{Type: FilterIn, Next: out}

	DoFilterPrints(filter)
	out.ExpectOrdered()

	filter.Add("a.b.c")
	filter.Add("c.b.a")
	filter.Add("d")

	filter.AddPrefix("b")
	filter.AddPrefix("d")

	filter.AddSuffix("c")
	filter.AddSuffix("d")

	WaitForPropagation()

	DoFilterPrints(filter)
	out.ExpectOrdered(
		"<a.c> x",
		"<a.b.c> x",

		"<b> x",
		"<b.a> x",
		"<b.c> x",
		"<b.a.c> x",
		"<b.c.a> x",

		"<c> x",
		"<c.b.a> x",
	)

	filter.Remove("a.b.c")
	filter.Remove("c.b.a")
	filter.Remove("e")

	filter.AddPrefix("c")
	filter.RemovePrefix("b")
	filter.RemovePrefix("e")

	filter.AddSuffix("b")
	filter.RemoveSuffix("c")
	filter.RemoveSuffix("e")

	WaitForPropagation()

	DoFilterPrints(filter)
	out.ExpectOrdered(
		"<a.b> x",
		"<a.c.b> x",

		"<b> x",

		"<c> x",
		"<c.a> x",
		"<c.b> x",
		"<c.a.b> x",
		"<c.b.a> x",
	)
}

func BenchFilterKey(b *testing.B, n int) {
	filter := &Filter{Type: FilterIn, Next: NilPrinter}
	l := L("a", "x")

	for i := 0; i < n; i++ {
		filter.Add(strconv.Itoa(i))
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		filter.Print(l)
	}
}

func BenchmarkFilter_Key0(b *testing.B)  { BenchFilterKey(b, 0) }
func BenchmarkFilter_Key1(b *testing.B)  { BenchFilterKey(b, 1) }
func BenchmarkFilter_Key8(b *testing.B)  { BenchFilterKey(b, 8) }
func BenchmarkFilter_Key32(b *testing.B) { BenchFilterKey(b, 32) }

const (
	Miss = iota
	Early
	Late
)

func BenchFilterPrefix(b *testing.B, n int, hit int) {
	filter := &Filter{Type: FilterOut, Next: NilPrinter}
	l := L("a", "x")

	if hit != Miss {
		n--
	}

	if hit == Early {
		filter.AddPrefix("a")
	}

	for i := 0; i < n; i++ {
		filter.AddPrefix(strconv.Itoa(i))
	}

	if hit == Late {
		filter.AddPrefix("a")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		filter.Print(l)
	}
}

func BenchmarkFilter_PrefixMiss0(b *testing.B)  { BenchFilterPrefix(b, 0, Miss) }
func BenchmarkFilter_PrefixMiss1(b *testing.B)  { BenchFilterPrefix(b, 1, Miss) }
func BenchmarkFilter_PrefixMiss8(b *testing.B)  { BenchFilterPrefix(b, 8, Miss) }
func BenchmarkFilter_PrefixMiss16(b *testing.B) { BenchFilterPrefix(b, 16, Miss) }
func BenchmarkFilter_PrefixMiss32(b *testing.B) { BenchFilterPrefix(b, 32, Miss) }

func BenchmarkFilter_PrefixHitEarly0(b *testing.B)  { BenchFilterPrefix(b, 0, Early) }
func BenchmarkFilter_PrefixHitEarly1(b *testing.B)  { BenchFilterPrefix(b, 1, Early) }
func BenchmarkFilter_PrefixHitEarly8(b *testing.B)  { BenchFilterPrefix(b, 8, Early) }
func BenchmarkFilter_PrefixHitEarly16(b *testing.B) { BenchFilterPrefix(b, 16, Early) }
func BenchmarkFilter_PrefixHitEarly32(b *testing.B) { BenchFilterPrefix(b, 32, Early) }

func BenchmarkFilter_PrefixHitLate0(b *testing.B)  { BenchFilterPrefix(b, 0, Late) }
func BenchmarkFilter_PrefixHitLate1(b *testing.B)  { BenchFilterPrefix(b, 1, Late) }
func BenchmarkFilter_PrefixHitLate8(b *testing.B)  { BenchFilterPrefix(b, 8, Late) }
func BenchmarkFilter_PrefixHitLate16(b *testing.B) { BenchFilterPrefix(b, 16, Late) }
func BenchmarkFilter_PrefixHitLate32(b *testing.B) { BenchFilterPrefix(b, 32, Late) }
