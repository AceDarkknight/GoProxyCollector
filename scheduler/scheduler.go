package scheduler

import (
	"math/rand"
	"sync"
	"time"

	"github.com/AceDarkkinght/GoProxyCollector/collector"
	"github.com/AceDarkkinght/GoProxyCollector/storage"
	"github.com/AceDarkkinght/GoProxyCollector/util"
)

func Start(collector collector.Collector, storage storage.Storage) {
	if collector == nil || storage == nil {
		return
	}

	var wg sync.WaitGroup

	for {
		if !collector.Next() {
			break
		}

		// Collect.
		results, err := collector.Collect()
		if err == nil && len(results) > 0 {
			// Verify.
			for _, r := range results {
				wg.Add(1)

				go func(ip string, port int, result collector.Result) {
					if util.VerifyProxyIp(ip, port) {
						storage.AddOrUpdate(ip, r)
					}

					defer wg.Done()
				}(r.Ip, r.Port, r)
			}

			wg.Wait()
		}

		// Wait at least 5s to avoid the website block our IP.
		t := rand.New(rand.NewSource(time.Now().Unix())).Intn(15) + 5
		time.Sleep(time.Duration(t) * time.Second)
	}
}
