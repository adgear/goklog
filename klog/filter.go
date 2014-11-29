// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import (
	"github.com/datacratic/goset"

	"log"
	"strings"
	"sync"
)

const (
	FilterOut = 1
	FilterIn  = 2
)

const (
	filterAdd = iota
	filterRemove
	filterAddPrefix
	filterRemovePrefix
	filterAddSuffix
	filterRemoveSuffix
)

type filterOp struct {
	Op    int
	Value string
}

type Filter struct {
	Type int

	Keys     []string
	Prefixes []string
	Suffixes []string

	Next Printer

	initialize sync.Once

	keys     set.String
	prefixes set.String
	suffixes set.String

	printC chan *Line
	opC    chan filterOp
	getC   chan chan map[string][]string
}

func NewFilter(def int) *Filter           { return &Filter{Type: def} }
func (filter *Filter) Chain(next Printer) { filter.Next = next }

func (filter *Filter) Init() {
	filter.initialize.Do(filter.init)
}

func (filter *Filter) init() {
	if filter.Type != FilterOut && filter.Type != FilterIn {
		log.Panicf("invalid filter default '%d'", filter.Type)
	}

	if filter.Next == nil {
		filter.Next = NilPrinter
	}

	filter.keys = set.NewString(filter.Keys...)
	filter.prefixes = set.NewString(filter.Prefixes...)
	filter.suffixes = set.NewString(filter.Suffixes...)

	filter.printC = make(chan *Line, DefaultBufferC)
	filter.opC = make(chan filterOp)
	filter.getC = make(chan chan map[string][]string)

	go filter.run()
}

func (filter *Filter) Add(value string) {
	filter.Init()
	filter.opC <- filterOp{filterAdd, value}
}

func (filter *Filter) Remove(value string) {
	filter.Init()
	filter.opC <- filterOp{filterRemove, value}
}

func (filter *Filter) AddPrefix(prefix string) {
	filter.Init()
	filter.opC <- filterOp{filterAddPrefix, prefix}
}

func (filter *Filter) RemovePrefix(prefix string) {
	filter.Init()
	filter.opC <- filterOp{filterRemovePrefix, prefix}
}

func (filter *Filter) AddSuffix(suffix string) {
	filter.Init()
	filter.opC <- filterOp{filterAddSuffix, suffix}
}

func (filter *Filter) RemoveSuffix(suffix string) {
	filter.Init()
	filter.opC <- filterOp{filterRemoveSuffix, suffix}
}

func (filter *Filter) Get() map[string][]string {
	filter.Init()

	resultC := make(chan map[string][]string)
	filter.getC <- resultC
	return <-resultC
}

func (filter *Filter) Print(line *Line) {
	filter.Init()
	filter.printC <- line
}

func (filter *Filter) print(line *Line) {
	hit := filter.keys.Test(line.Key)

	if !hit {
		for prefix := range filter.prefixes {
			if hit = strings.HasPrefix(line.Key, prefix); hit {
				break
			}
		}
	}

	if !hit {
		for suffix := range filter.suffixes {
			if hit = strings.HasSuffix(line.Key, suffix); hit {
				break
			}
		}
	}

	if filter.Type == FilterOut && !hit {
		filter.Next.Print(line)
	} else if filter.Type == FilterIn && hit {
		filter.Next.Print(line)
	}
}

func (filter *Filter) op(op int, value string) {
	switch op {

	case filterAdd:
		filter.keys.Put(value)
	case filterRemove:
		filter.keys.Del(value)

	case filterAddPrefix:
		filter.prefixes.Put(value)
	case filterRemovePrefix:
		filter.prefixes.Del(value)

	case filterAddSuffix:
		filter.suffixes.Put(value)
	case filterRemoveSuffix:
		filter.suffixes.Del(value)

	default:
		log.Panicf("unknown filter op type '%d'", op)
	}
}

func (filter *Filter) get(resultC chan map[string][]string) {
	resultC <- map[string][]string{
		"keys":     filter.keys.Array(),
		"prefixes": filter.prefixes.Array(),
		"suffixes": filter.suffixes.Array(),
	}
}

func (filter *Filter) run() {
	for {
		select {
		case line := <-filter.printC:
			filter.print(line)
		case op := <-filter.opC:
			filter.op(op.Op, op.Value)
		case c := <-filter.getC:
			filter.get(c)
		}
	}
}
