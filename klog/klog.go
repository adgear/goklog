// Copyright (c) 2014 Datacratic. All rights reserved.

package klog

// DefaultBufferC is the default channel size for stages that defer printing to
// background go-routines. A well tuned parameter can reduce the overhead of
// going through a channel while introducing a slight delay in the output of the
// log line.
const DefaultBufferC = 8

// DefaultPathREST is the default REST path prefix used if none is provided
// explicitly.
const DefaultPathREST = "/debug/klog"
