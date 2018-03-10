package scheduler

import (
	"math/rand"
	"time"

	"github.com/AceDarkkinght/GoProxyCollector/collector"
	"github.com/AceDarkkinght/GoProxyCollector/storage"
	"github.com/AceDarkkinght/GoProxyCollector/util"
)

func Start(collector collector.Collector, storage storage.Storage) {
	if collector == nil || storage == nil {
		return
	}

	for {
		if !collector.Next() {
			break
		}

		// Collect.
		results, err := collector.Collect()
		if err == nil && len(results) > 0 {
			// Verify.
			for _, r := range results {
				if util.VerifyHTTP(r.Ip, r.Port) {
					storage.AddOrUpdate(r.Ip, r)
				}
			}
		}

		// Wait at least 2s.
		t := rand.New(rand.NewSource(time.Now().Unix())).Intn(10) + 2
		time.Sleep(time.Duration(t) * time.Second)
	}
}
