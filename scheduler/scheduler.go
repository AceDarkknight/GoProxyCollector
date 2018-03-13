package scheduler

import (
	"fmt"
	_ "math/rand"
	"sync"
	"time"

	"github.com/AceDarkkinght/GoProxyCollector/collector"
	"github.com/AceDarkkinght/GoProxyCollector/result"

	//"github.com/AceDarkkinght/GoProxyCollector/result"
	"github.com/AceDarkkinght/GoProxyCollector/storage"
	"github.com/AceDarkkinght/GoProxyCollector/util"
)

func Run(collector collector.Collector, storage storage.Storage) {
	if collector == nil || storage == nil {
		return
	}

	//var wg sync.WaitGroup

	for {
		resultChan := make(chan *result.Result, 1)
		if !collector.Next() {
			break
		}

		// Collect.
		go collector.Collect(resultChan)

		// Verify.
		go Verify(resultChan, storage)
		//for r := range resultChan {
		//	wg.Add(1)
		//
		//	go func() {
		//		if util.VerifyProxyIp(r.Ip, r.Port) {
		//			storage.AddOrUpdate(r.Ip, r)
		//		}
		//
		//		defer wg.Done()
		//	}()
		//
		//	wg.Wait()
		//}

		//if err == nil && len(results) > 0 {
		//	// Verify.
		//	for _, result := range results {
		//		wg.Add(1)
		//
		//		r := result
		//		go func() {
		//			if util.VerifyProxyIp(r.Ip, r.Port) {
		//				storage.AddOrUpdate(r.Ip, r)
		//			}
		//
		//			defer wg.Done()
		//		}()
		//	}
		//
		//	wg.Wait()
		//}

		//if err == nil && results != nil && len(results) > 0 {
		//	for i := 0; i < len(results); i++ {
		//		wg.Add(1)
		//
		//		go func(r *result.Result) {
		//			if util.VerifyProxyIp(r.Ip, r.Port) {
		//				storage.AddOrUpdate(r.Ip, r)
		//			}
		//
		//			defer wg.Done()
		//		}(results[i])
		//	}
		//
		//	wg.Wait()
		//}

		// Wait at least 2s to avoid the website block our IP.
		//t := rand.New(rand.NewSource(time.Now().Unix())).Intn(10) + 2
		time.Sleep(30 * time.Second)
	}
}

func Verify(resultChan <-chan *result.Result, storage storage.Storage) {
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
