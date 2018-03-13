package verifier

import (
	"encoding/json"
	"sync"

	//"github.com/AceDarkkinght/GoProxyCollector/collector"
	"github.com/AceDarkkinght/GoProxyCollector/result"
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

		go func(v []byte) {
			var result result.Result
			json.Unmarshal(v, &result)
			if !util.VerifyProxyIp(ip, result.Port) {
				storage.Delete(ip)
			}

			defer wg.Done()
		}(value)
	}

	wg.Wait()
}
