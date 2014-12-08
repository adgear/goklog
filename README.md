# goklog #

Keyed logger used to define processing pipelines on a log stream.

## Installation ##

You can download the code via the usual go utilities:

```
go get github.com/datacratic/goklog
```

To build the code and run the test suite along with several static analysis
tools, use the provided Makefile:

```
make test
```

Note that the usual go utilities will work just fine but we require that all
commits pass the full suite of tests and static analysis tools.

## Why Another Logging Library? ##

While klog does offer nifty features, it was generally designed to compliment
existing log framework. It does so by defining a very generic set of
abstractions (lines, printers and chained printers) which can be used to
manipulate the log output prior to handing off to an existing log
framework.

## Examples ##

Examples are available in the following test suites:

* [**Pipeline**](klog/example_test.go): Basic library usage along with how to
  setup a pipeline.
* [**REST pipeline**](klog/example_rest_test.go): How to setup a REST enabled
pipeline.

## Stages ##

In this section we'll give a quick overview of the various klog printers.

### Filter ###

Filter can filter out/in lines based on a full-text, prefix or suffix matching
rules. Note that this implementation is generally more expensive then filtering
strategies used other logging libraries but its performance are generally
acceptable and the flexbility is always useful.

Filters are also available through a REST interface which allows for the live
manipulation of logs in real-time.

| Path | Method | Description |
| --- | --- | --- |
| `/debug/klog/filter` | `GET` | Returns the list of active patterns |
| `/debug/klog/filter/key/:key` | `PUT` | Adds the given full-text pattern |
| `/debug/klog/filter/key/:key` | `DELETE` | Removes the given full-test pattern |
| `/debug/klog/filter/prefix/:prefix` | `PUT` | Adds the given prefix pattern |
| `/debug/klog/filter/prefix/:prefix` | `DELETE` | Removes the given prefix pattern |
| `/debug/klog/filter/suffix/:suffix` | `PUT` | Adds the given suffix pattern |
| `/debug/klog/filter/suffix/:suffix` | `DELETE` | Removes the given suffix pattern |

### Dedup ###

Dedup is used to aggregate the consecutive identical lines for a given key into
a single line. This is useful to avoid flooding the logs with endless identical
messages.

### Ring ###

Ring logs all the received lines into a fixed size ring buffer in a lock-free
manner. It's mostly useful when coupled with its REST interface which enables
the logs of a given service to be queried remotely.

| Path | Method | Description |
| --- | --- | --- |
| `/debug/klog/ring` | `GET` | Returns all the lines currently in the ring buffer |
| `/debug/klog/ring/key/:key` | `GET` | Returns all the lines associated with the given key |
| `/debug/klog/ring/prefix/:prefix` | `GET` | Returns all the lines associated with the given prefix |
| `/debug/klog/ring/suffix/:suffix` | `GET` | Returns all the lines associated with the given suffix |

## License ##

The source code is available under the Apache License. See the LICENSE file for
more details.
