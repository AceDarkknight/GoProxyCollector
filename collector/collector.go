package collector

import (
	"github.com/AceDarkkinght/GoProxyCollector/result"
)

type Collector interface {
	Next() bool
	Name() string
	Collect(chan<- *result.Result)
}

type Type uint8

const (
	COLLECTBYSELECTOR Type = iota
	COLLECTBYREGEX
)
