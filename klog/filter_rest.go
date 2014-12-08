// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

import (
	"github.com/datacratic/gorest/rest"
)

// FilterREST provides the REST interface for the Filter chained printer.
type FilterREST struct {
	*Filter

	// PathPrefix will be pre-pended to all the REST paths. Defaults to
	// DefaultPathREST.
	PathPrefix string
}

// NewFilterREST creates a new REST enabled Filter chained printer at the
// specified path. The def parameter should be either FilterIn or FilterOut.
func NewFilterREST(path string, def int) *FilterREST {
	filter := &FilterREST{Filter: NewFilter(def), PathPrefix: path}
	rest.AddService(filter)
	return filter
}

// RESTRoutes returns the set of gorest routes used to manipulate the Filter
// chained printer.
func (filter *FilterREST) RESTRoutes() rest.Routes {
	prefix := filter.PathPrefix
	if len(prefix) == 0 {
		prefix = DefaultPathREST + "/filter"
	}

	return []*rest.Route{
		rest.NewRoute(prefix, "GET", filter.Get),

		rest.NewRoute(prefix+"/key/:key", "PUT", filter.add),
		rest.NewRoute(prefix+"/key/:key", "DELETE", filter.remove),

		rest.NewRoute(prefix+"/prefix/:prefix", "PUT", filter.addPrefix),
		rest.NewRoute(prefix+"/prefix/:prefix", "DELETE", filter.removePrefix),

		rest.NewRoute(prefix+"/suffix/:suffix", "PUT", filter.addSuffix),
		rest.NewRoute(prefix+"/suffix/:suffix", "DELETE", filter.removeSuffix),
	}
}

func (filter *FilterREST) add(value string)          { filter.Add(value) }
func (filter *FilterREST) remove(value string)       { filter.Remove(value) }
func (filter *FilterREST) addPrefix(value string)    { filter.AddPrefix(value) }
func (filter *FilterREST) removePrefix(value string) { filter.RemovePrefix(value) }
func (filter *FilterREST) addSuffix(value string)    { filter.AddSuffix(value) }
func (filter *FilterREST) removeSuffix(value string) { filter.RemoveSuffix(value) }
