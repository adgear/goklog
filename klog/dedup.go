// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import (
	"fmt"
	"sync"
	"time"
)

const DefaultDedupRate = 1 * time.Second

type dedupLine struct {
	Value string
	Count int
}

type Dedup struct {
	Rate time.Duration
	Next Printer

	initialize sync.Once

	lines  map[string]*dedupLine
	printC chan *Line
}

func NewDedup() *Dedup                  { return new(Dedup) }
func (dedup *Dedup) Chain(next Printer) { dedup.Next = next }

func (dedup *Dedup) Init() {
	dedup.initialize.Do(dedup.init)
}

func (dedup *Dedup) init() {
	if dedup.Rate == 0 {
		dedup.Rate = DefaultDedupRate
	}

	if dedup.Next == nil {
		dedup.Next = NilPrinter
	}

	dedup.lines = make(map[string]*dedupLine)
	dedup.printC = make(chan *Line, DefaultBufferC)

	go dedup.run()
}

func (dedup *Dedup) Print(line *Line) {
	dedup.Init()
	dedup.printC <- line
}

func (dedup *Dedup) print(line *Line) {
	counter, ok := dedup.lines[line.Key]
	if !ok {
		counter = new(dedupLine)
		dedup.lines[line.Key] = counter
	}

	if counter.Value != line.Value {
		dedup.send(line.Key, counter)

		dedup.Next.Print(line)
		counter.Count = 0
		counter.Value = line.Value

	} else {
		counter.Count++
	}
}

func (dedup *Dedup) flush() {
	for key, counter := range dedup.lines {
		dedup.send(key, counter)
		counter.Count = 0
	}
}

func (dedup *Dedup) send(key string, counter *dedupLine) {
	if counter.Count == 0 {
		return
	}

	var value string

	if counter.Count == 1 {
		value = counter.Value
	} else {
		value = fmt.Sprintf("%s [%d times]", counter.Value, counter.Count)
	}

	dedup.Next.Print(&Line{Timestamp: time.Now(), Key: key, Value: value})
}

func (dedup *Dedup) run() {
	tickC := time.Tick(dedup.Rate)

	for {
		select {
		case line := <-dedup.printC:
			dedup.print(line)

		case <-tickC:
			dedup.flush()
		}
	}
}
