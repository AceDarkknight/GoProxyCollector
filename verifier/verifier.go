package verifier

import (
	"encoding/json"
	"sync"

	"github.com/AceDarkknight/GoProxyCollector/result"
	"github.com/AceDarkknight/GoProxyCollector/storage"
	"github.com/AceDarkknight/GoProxyCollector/util"

	"github.com/cihub/seelog"
)

// VerifyAndSave existing Ips to check it's available or not. Delete the unavailable Ips.
func VerifyAndDelete(storage storage.Storage) {
	if storage == nil {
		return
	}

	var wg sync.WaitGroup

	items := storage.GetAll()
	for _, value := range items {
		wg.Add(1)

		go func(v []byte) {
			var r result.Result
			json.Unmarshal(v, &r)
			if !util.VerifyProxyIp(r.Ip, r.Port) {
				storage.Delete(r.Ip)
			}

			defer wg.Done()
		}(value)
	}

	wg.Wait()
}

// Verify ips in channel. Save the available ips.
func VerifyAndSave(resultChan <-chan *result.Result, storage storage.Storage) {
	var wg sync.WaitGroup
	for r := range resultChan {
		wg.Add(1)
		go func(r *result.Result) {
			if util.VerifyProxyIp(r.Ip, r.Port) {
				storage.AddOrUpdate(r.Ip, r)
				seelog.Debugf("insert %s to DB", r.Ip)
			}

			defer wg.Done()
		}(r)
	}

	wg.Wait()
}
