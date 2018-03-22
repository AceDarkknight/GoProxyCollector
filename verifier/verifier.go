package verifier

import (
	"encoding/json"
	"strconv"
	"sync"
	"time"

	"github.com/AceDarkkinght/GoProxyCollector/result"
	"github.com/AceDarkkinght/GoProxyCollector/storage"
	"github.com/AceDarkkinght/GoProxyCollector/util"
	"github.com/cihub/seelog"
	"github.com/parnurzeal/gorequest"
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

type Verifier struct {
	pool sync.Pool
}

func NewVerifier() *Verifier {
	pool := sync.Pool{
		New: func() interface{} {
			return gorequest.New().Timeout(time.Second * 5).Get("http://httpbin.org/get")
		},
	}

	return &Verifier{pool: pool}
}

func (v *Verifier) Verify(ips <-chan *result.Result, availableIps chan<- *result.Result) {
	var wg sync.WaitGroup
	for ip := range ips {
		proxy := "http://" + ip.Ip + ":" + strconv.Itoa(ip.Port)
		wg.Add(1)
		go func(r *result.Result) {
			defer wg.Done()
			superAgent := v.pool.Get().(*gorequest.SuperAgent)
			resp, _, errs := superAgent.Proxy(proxy).End()

			v.pool.Put(superAgent)

			if errs != nil {
				return
			}

			if resp.StatusCode != 200 {
				return
			}

			availableIps <- r
		}(ip)
	}

	wg.Wait()
	defer close(availableIps)
}
