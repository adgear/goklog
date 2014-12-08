// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import (
	"github.com/datacratic/goset"

	"log"
	"strings"
	"sync"
)

const (
	// FilterOut indicates that by default all keys are allowed and that any keys
	// that hit on the filter will be discarded.
	FilterOut = 1

	// FilterIn indicates that by default all keys are discarded and that only
	// keys that hit on the filter will be kept.
	FilterIn = 2
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

// Filter can filter a line stream on the key of each line and discard any
// undesired lines. Filters can either be a full, prefix or suffix string match.
type Filter struct {
	Chained

	// Type is either FilterOut or FilterIn.
	Type int

	// Keys is the initial set of keys that will be used for full-key matches.
	Keys []string

	// Prefixes is the initial set of patterns use for prefix matches.
	Prefixes []string

	// Suffixes is the initial set of patterns use for suffix matches.
	Suffixes []string

	initialize sync.Once

	keys     set.String
	prefixes set.String
	suffixes set.String

	printC chan *Line
	opC    chan filterOp
	getC   chan chan map[string][]string
}

// NewFilter creates a new Filter configured to either FilterIn or FilterOut.
func NewFilter(def int) *Filter { return &Filter{Type: def} }

// Init initializes the filter. Calling this is optional since the object is
// lazily initialized as needed.
func (filter *Filter) Init() {
	filter.initialize.Do(filter.init)
}

func (filter *Filter) init() {
	if filter.Type != FilterOut && filter.Type != FilterIn {
		log.Panicf("invalid filter default '%d'", filter.Type)
	}

	filter.keys = set.NewString(filter.Keys...)
	filter.prefixes = set.NewString(filter.Prefixes...)
	filter.suffixes = set.NewString(filter.Suffixes...)

	filter.printC = make(chan *Line, DefaultBufferC)
	filter.opC = make(chan filterOp)
	filter.getC = make(chan chan map[string][]string)

	go filter.run()
}

// Add adds the given pattern to be used as a full-key match.
func (filter *Filter) Add(values ...string) *Filter {
	filter.Init()

	for _, value := range values {
		filter.opC <- filterOp{filterAdd, value}
	}

	return filter
}

// Remove removes the given pattern to be used as a full-key match.
func (filter *Filter) Remove(values ...string) *Filter {
	filter.Init()

	for _, value := range values {
		filter.opC <- filterOp{filterRemove, value}
	}

	return filter
}

// AddPrefix adds the given pattern to be used as a prefix match.
func (filter *Filter) AddPrefix(prefixes ...string) *Filter {
	filter.Init()

	for _, prefix := range prefixes {
		filter.opC <- filterOp{filterAddPrefix, prefix}
	}

	return filter
}

// RemovePrefix removes the given pattern to be used as a prefix match.
func (filter *Filter) RemovePrefix(prefixes ...string) *Filter {
	filter.Init()

	for _, prefix := range prefixes {
		filter.opC <- filterOp{filterRemovePrefix, prefix}
	}

	return filter
}

// AddSuffix adds the given pattern to be used as a suffix match.
func (filter *Filter) AddSuffix(suffixes ...string) *Filter {
	filter.Init()

	for _, suffix := range suffixes {
		filter.opC <- filterOp{filterAddSuffix, suffix}
	}

	return filter
}

// RemoveSuffix removes the given pattern to be used as a suffix match.
func (filter *Filter) RemoveSuffix(suffixes ...string) *Filter {
	filter.Init()

	for _, suffix := range suffixes {
		filter.opC <- filterOp{filterRemoveSuffix, suffix}
	}

	return filter
}

// Get returns the list of active filters.
func (filter *Filter) Get() map[string][]string {
	filter.Init()

	resultC := make(chan map[string][]string)
	filter.getC <- resultC
	return <-resultC
}

// Print forwards the line to the next printer if the filter is of type FilterIn
// and at least one of the patterns match the key or if the filter is of type
// FilterOut and none of the patterns match the key.
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
		filter.PrintNext(line)
	} else if filter.Type == FilterIn && hit {
		filter.PrintNext(line)
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
