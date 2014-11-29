// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import (
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	DefaultRingSize = 10 * 1000
)

type Ring struct {
	Size int

	initialize sync.Once

	count uint32
	ring  []unsafe.Pointer
}

func NewRing(size int) *Ring { return &Ring{Size: size} }

func (ring *Ring) Init() {
	ring.initialize.Do(ring.init)
}

func (ring *Ring) init() {
	if ring.Size == 0 {
		ring.Size = DefaultRingSize
	}

	ring.ring = make([]unsafe.Pointer, ring.Size)
}

func (ring *Ring) GetAll() []*Line {
	ring.Init()
	return ring.get(func(*Line) bool { return true })
}

func (ring *Ring) GetKey(key string) []*Line {
	ring.Init()
	return ring.get(func(line *Line) bool { return line.Key == key })
}

func (ring *Ring) GetPrefix(prefix string) []*Line {
	ring.Init()
	return ring.get(func(line *Line) bool {
		return strings.HasPrefix(line.Key, prefix)
	})
}

func (ring *Ring) GetSuffix(suffix string) []*Line {
	ring.Init()
	return ring.get(func(line *Line) bool {
		return strings.HasSuffix(line.Key, suffix)
	})
}

func (ring *Ring) Print(line *Line) {
	ring.Init()

	pos := int(atomic.AddUint32(&ring.count, 1)-1) % len(ring.ring)
	atomic.StorePointer(&ring.ring[pos], unsafe.Pointer(line))
}

func (ring *Ring) get(filter func(*Line) bool) (result []*Line) {
	var lines []*Line

	for i := 0; i < len(ring.ring); i++ {
		line := (*Line)(atomic.LoadPointer(&ring.ring[i]))
		if line != nil && filter(line) {
			lines = append(lines, line)
		}
	}

	sort.Sort(lineArray(lines))

	for _, line := range lines {
		result = append(result, line)
	}

	return
}
