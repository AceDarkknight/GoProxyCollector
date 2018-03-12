package verifier

import (
	"encoding/json"
	"sync"

	"github.com/AceDarkkinght/GoProxyCollector/collector"
	"github.com/AceDarkkinght/GoProxyCollector/storage"
	"github.com/AceDarkkinght/GoProxyCollector/util"
)

func VerifyAll(storage storage.Storage) {
	if storage == nil {
		return
	}

	var wg sync.WaitGroup

	items := storage.GetAll()
	for ip, value := range items {
		wg.Add(1)

		go func() {
			var result collector.Result
			json.Unmarshal(value, &result)
			if !util.VerifyProxyIp(ip, result.Port) {
				storage.Delete(ip)
			}

			defer wg.Done()
		}()
	}
}
