package collector

import "github.com/AceDarkkinght/GoProxyCollector/result"

type Collector interface {
	Next() bool
	Collect(chan<- *result.Result)
}
