package collector

import "github.com/AceDarkkinght/GoProxyCollector/result"

type Collector interface {
	Next() bool
	Collect(chan<- *result.Result)
}

func NewCollector(name string) Collector {

}

type CollectorType uint
