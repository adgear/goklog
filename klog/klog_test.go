// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import (
	"github.com/datacratic/goset"

	"fmt"
	"sync"
	"testing"
	"time"
)

func WaitForPropagation() {
	time.Sleep(10 * time.Millisecond)
}

func L(key, value string) *Line {
	return &Line{time.Now(), key, value}
}

func Simplify(lines []*Line) (result []string) {
	for _, line := range lines {
		result = append(result, fmt.Sprintf("<%s> %s", line.Key, line.Value))
	}
	return
}

func ExpectOrdered(t *testing.T, lines []string, exp ...string) {
	err := false
	for i := 0; i < len(lines); i++ {
		a := "<missing>"
		if len(lines) > i {
			a = lines[i]
		} else {
			err = true
		}

		b := "<missing>"
		if len(exp) > i {
			b = exp[i]
		} else {
			err = true
		}

		cmp := ""
		if len(lines) > i && len(exp) > i {
			cmp = "=="
			if lines[i] != exp[i] {
				cmp = "!="
				err = true
			}
		}

		fmt.Printf("%d: '%s' %s '%s'\n", i, a, cmp, b)
	}

	if err {
		t.Fail()
	}
}

func ExpectUnordered(t *testing.T, lines []string, exp ...string) {
	a := set.NewString(lines...)
	b := set.NewString(exp...)

	if diff := a.Difference(b); len(diff) != 0 {
		t.Errorf("FAIL: extra lines %s", diff)
	}

	if diff := b.Difference(a); len(diff) != 0 {
		t.Errorf("FAIL: extra lines %s", diff)
	}
}

type TestPrinter struct {
	T *testing.T

	initialize sync.Once
	linesC     chan *Line
}

func (printer *TestPrinter) Init() {
	printer.initialize.Do(printer.init)
}

func (printer *TestPrinter) init() {
	printer.linesC = make(chan *Line, 100)
}

func (printer *TestPrinter) Print(line *Line) {
	printer.Init()
	printer.linesC <- line
}

func (printer *TestPrinter) GetLines(n int) (result []*Line) {
	printer.Init()
	timeoutC := time.After(100 * time.Millisecond)

	for {
		select {
		case line := <-printer.linesC:
			if result = append(result, line); len(result) == n {
				return
			}
		case <-timeoutC:
			return
		}
	}
}

func (printer *TestPrinter) ExpectOrdered(exp ...string) {
	lines := printer.GetLines(len(exp))
	ExpectOrdered(printer.T, Simplify(lines), exp...)
}

func (printer *TestPrinter) ExpectUnordered(exp ...string) {
	lines := printer.GetLines(len(exp))
	ExpectUnordered(printer.T, Simplify(lines), exp...)
}
