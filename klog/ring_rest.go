// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import (
	"github.com/datacratic/gorest/rest"
)

type RestRing struct {
	*Ring
	PathPrefix string
}

func NewRestRing(path string, size int) *RestRing {
	ring := &RestRing{Ring: NewRing(size), PathPrefix: path}
	rest.AddService(ring)
	return ring
}

func (ring *RestRing) RESTRoutes() rest.Routes {
	prefix := ring.PathPrefix
	if len(prefix) == 0 {
		prefix = DefaultRestPath + "/ring"
	}

	return []*rest.Route{
		rest.NewRoute(prefix, "GET", ring.GetAll),
		rest.NewRoute(prefix+"/key/:key", "GET", ring.GetKey),
		rest.NewRoute(prefix+"/prefix/:prefix", "GET", ring.GetPrefix),
		rest.NewRoute(prefix+"/suffix/:suffix", "GET", ring.GetSuffix),
	}
}
