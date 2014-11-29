// Copyright (c) 2014 Datacratic. All rights reserved.

package klogr

import (
	"github.com/datacratic/goklog/klog"
	"github.com/datacratic/gorest/rest"
)

type RestRing struct {
	*klog.Ring
	PathPrefix string
}

func NewRestRing(path string, size int) *RestRing {
	ring := &RestRing{Ring: klog.NewRing(size), PathPrefix: path}
	rest.AddService(ring)
	return ring
}

func (ring *RestRing) RESTRoutes() rest.Routes {
	prefix := ring.PathPrefix
	if len(prefix) == 0 {
		prefix = DefaultPath + "/ring"
	}

	return []*rest.Route{
		rest.NewRoute(prefix, "GET", ring.GetAll),
		rest.NewRoute(prefix+"/key/:key", "GET", ring.GetKey),
		rest.NewRoute(prefix+"/prefix/:prefix", "GET", ring.GetPrefix),
		rest.NewRoute(prefix+"/suffix/:suffix", "GET", ring.GetSuffix),
	}
}
