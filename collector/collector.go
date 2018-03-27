package collector

import (
	"github.com/AceDarkknight/GoProxyCollector/result"
)

type Collector interface {
	Next() bool
	Name() string
	Collect(chan<- *result.Result) []error
}

type Type uint8

const (
	COLLECTBYSELECTOR Type = iota
	COLLECTBYREGEX
)
