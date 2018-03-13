package verifier

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/AceDarkkinght/GoProxyCollector/result"
	"github.com/AceDarkkinght/GoProxyCollector/storage"
	"github.com/AceDarkkinght/GoProxyCollector/util"
)

// VerifyAndSave existing Ips to check it's available or not. Delete the unavailable Ips.
func VerifyAndDelete(storage storage.Storage) {
	if storage == nil {
		return
	}

	var wg sync.WaitGroup

	items := storage.GetAll()
	for ip, value := range items {
		wg.Add(1)

		go func(v []byte) {
			var r result.Result
			json.Unmarshal(v, &r)
			if !util.VerifyProxyIp(ip, r.Port) {
				storage.Delete(ip)
			}

			defer wg.Done()
		}(value)
	}

	wg.Wait()
}

// Verify Ips in channel. Save the available Ips.
func VerifyAndSave(resultChan <-chan *result.Result, storage storage.Storage) {
	var wg sync.WaitGroup
	for r := range resultChan {
		wg.Add(1)
		go func(r *result.Result) {
			if util.VerifyProxyIp(r.Ip, r.Port) {
				fmt.Printf("address %p,Ip:%s\n", r, r.Ip)
				storage.AddOrUpdate(r.Ip, r)
			}

			defer wg.Done()
		}(r)
	}

	wg.Wait()
}
