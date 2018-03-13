package scheduler

import (
	"math/rand"
	"time"

	"github.com/AceDarkkinght/GoProxyCollector/collector"
	"github.com/AceDarkkinght/GoProxyCollector/result"
	"github.com/AceDarkkinght/GoProxyCollector/storage"
	"github.com/AceDarkkinght/GoProxyCollector/verifier"
)

func Run(collector collector.Collector, storage storage.Storage) {
	if collector == nil || storage == nil {
		return
	}

	for {
		resultChan := make(chan *result.Result, 10)
		if !collector.Next() {
			break
		}

		// Collect.
		go collector.Collect(resultChan)

		// Verify.
		verifier.VerifyAndSave(resultChan, storage)

		// Wait at least 2s to avoid the website block our IP.
		t := rand.New(rand.NewSource(time.Now().Unix())).Intn(10) + 2
		time.Sleep(time.Duration(t) * time.Second)
	}
}
