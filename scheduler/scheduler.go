package scheduler

import (
	"github.com/AceDarkkinght/GoProxyCollector/collector"
)

func Start(collector collector.Collector) {
	if collector == nil {
		return
	}

	resultChan := make(chan collector.Result, 100)
	for {
		url := collector.Next()
		if url == "" {
			break
		}

		results, err := collector.Collect(url)
		if err != nil {
			return
		}

		if len(results) == 0 {
			return
		}

		for _, result := range results {
			resultChan <- result
		}
	}
}
