// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import (
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"
)

// DefaultRingSize is used if Ring.Size is set to 0.
const DefaultRingSize = 1000

// Ring adds all printed lines to a fixed size ring buffer which is written and
// read atomically. Lines in the ring are read-back all at once and can be
// filtered as needed.
type Ring struct {

	// Size indicates the size of the ring used to log lines to. If 0 then
	// DefaultRingSize is used instead.
	Size int

	initialize sync.Once

	count uint32
	ring  []unsafe.Pointer
}

// NewRing creates a new Ring printer of the given size. If size is 0 then
// DefaultRingSize is used instead.
func NewRing(size int) *Ring { return &Ring{Size: size} }

// Init initializes the object. Calling this is optional since the object will
// lazily initialize itself when needed.
func (ring *Ring) Init() {
	ring.initialize.Do(ring.init)
}

func (ring *Ring) init() {
	if ring.Size == 0 {
		ring.Size = DefaultRingSize
	}

	ring.ring = make([]unsafe.Pointer, ring.Size)
}

// GetAll returns all the lines in the ring sorted by their timestamp.
func (ring *Ring) GetAll() []*Line {
	ring.Init()
	return ring.get(func(*Line) bool { return true })
}

// GetKey returns all the lines in the ring with the given key sorted by their
// timestamp.
func (ring *Ring) GetKey(key string) []*Line {
	ring.Init()
	return ring.get(func(line *Line) bool { return line.Key == key })
}

// GetPrefix returns all the lines in the ring with the given prefix sorted by
// their timestamp.
func (ring *Ring) GetPrefix(prefix string) []*Line {
	ring.Init()
	return ring.get(func(line *Line) bool {
		return strings.HasPrefix(line.Key, prefix)
	})
}

// GetSuffix returns all the lines in the ring with the given suffix sorted by
// their timestamp.
func (ring *Ring) GetSuffix(suffix string) []*Line {
	ring.Init()
	return ring.get(func(line *Line) bool {
		return strings.HasSuffix(line.Key, suffix)
	})
}

// Print adds the given line to the ring overwritting any older line present.
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
