// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import (
	"github.com/datacratic/gorest/rest"
)

type RingREST struct {
	*Ring
	PathPrefix string
}

func NewRingREST(path string, size int) *RingREST {
	ring := &RingREST{Ring: NewRing(size), PathPrefix: path}
	rest.AddService(ring)
	return ring
}

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
