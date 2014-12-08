// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import (
	"github.com/datacratic/gorest/rest"
)

type RestFilter struct {
	*Filter
	PathPrefix string
}

func NewRestFilter(path string, def int) *RestFilter {
	filter := &RestFilter{Filter: NewFilter(def), PathPrefix: path}
	rest.AddService(filter)
	return filter
}

func (filter *RestFilter) RESTRoutes() rest.Routes {
	prefix := filter.PathPrefix
	if len(prefix) == 0 {
		prefix = DefaultRestPath + "/filter"
	}

	return []*rest.Route{
		rest.NewRoute(prefix, "GET", filter.Get),

		rest.NewRoute(prefix+"/key/:key", "PUT", filter.Add),
		rest.NewRoute(prefix+"/key/:key", "DELETE", filter.Remove),

		rest.NewRoute(prefix+"/prefix/:prefix", "PUT", filter.AddPrefix),
		rest.NewRoute(prefix+"/prefix/:prefix", "DELETE", filter.RemovePrefix),

		rest.NewRoute(prefix+"/suffix/:suffix", "PUT", filter.AddSuffix),
		rest.NewRoute(prefix+"/suffix/:suffix", "DELETE", filter.RemoveSuffix),
	}
}
