// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import (
	"fmt"
	"sync"
	"time"
)

// DefaultDedupRate is used when Dedup.Rate is left empty.
const DefaultDedupRate = 1 * time.Second

type dedupLine struct {
	Value string
	Count int
}

// Dedup aggregates multiple consecutive identical lines into a single line with
// the number of time it was seen appended at the end. The first line
// encountered is always printed right away and only subsequent identical lines
// are held back prior to being printed. Held back line are printed at a set
// configurable rate.
type Dedup struct {
	Chained

	// Rate determines the interval at which duplicated lines are dumped.
	Rate time.Duration

	initialize sync.Once

	lines  map[string]*dedupLine
	printC chan *Line
}

// NewDedup creates a new Dedup printer.
func NewDedup() *Dedup { return new(Dedup) }

// Init initializes the object. Note that calling this is optional since the
// object is lazily initialized as needed.
func (dedup *Dedup) Init() {
	dedup.initialize.Do(dedup.init)
}

func (dedup *Dedup) init() {
	if dedup.Rate == 0 {
		dedup.Rate = DefaultDedupRate
	}

	dedup.lines = make(map[string]*dedupLine)
	dedup.printC = make(chan *Line, DefaultBufferC)

	go dedup.run()
}

// Print checks the line checking for duplicates. If the line was never seen
// before it is passed to the chained printer right away otherwise it is held
// back and counted.
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

		dedup.PrintNext(line)
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

	dedup.PrintNext(&Line{Timestamp: time.Now(), Key: key, Value: value})
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
