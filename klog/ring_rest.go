// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import (
	"github.com/datacratic/gorest/rest"
)

// RingREST provides the REST interface for the Ring printer.
type RingREST struct {
	*Ring

	// PathPrefix will be preprended to all the REST paths. Defaults to
	// DefaultPathREST.
	PathPrefix string
}

// NewRingREST creates a new REST enabled Ring printer at the specified path
// with the given size. If path is empty then DefaultPathREST will be used
// instead.
func NewRingREST(path string, size int) *RingREST {
	ring := &RingREST{Ring: NewRing(size), PathPrefix: path}
	rest.AddService(ring)
	return ring
}

// RESTRoutes returns the set of gorest routes used to manipulate the Ring
// printer.
func (ring *RingREST) RESTRoutes() rest.Routes {
	prefix := ring.PathPrefix
	if len(prefix) == 0 {
		prefix = DefaultPathREST + "/ring"
	}

	return []*rest.Route{
		rest.NewRoute(prefix, "GET", ring.GetAll),
		rest.NewRoute(prefix+"/key/:key", "GET", ring.GetKey),
		rest.NewRoute(prefix+"/prefix/:prefix", "GET", ring.GetPrefix),
		rest.NewRoute(prefix+"/suffix/:suffix", "GET", ring.GetSuffix),
	}
}
