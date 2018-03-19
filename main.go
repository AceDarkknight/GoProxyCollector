package main

import (
	"reflect"
	"sync"
	"time"

	"github.com/AceDarkkinght/GoProxyCollector/collector"
	"github.com/AceDarkkinght/GoProxyCollector/scheduler"
	"github.com/AceDarkkinght/GoProxyCollector/server"
	"github.com/AceDarkkinght/GoProxyCollector/storage"
	"github.com/AceDarkkinght/GoProxyCollector/verifier"

	"github.com/cihub/seelog"
)

func main() {
	// Load log.
	scheduler.SetLogger("logConfig.xml")
	defer seelog.Flush()

	database, err := storage.NewBoltDbStorage("proxy.db", "IpList")
	if err != nil {
		seelog.Critical(err)
		panic(err)
	}

	// Sync data.
	database.SyncKeys()
	seelog.Infof("database initialize finish.")
	defer database.Close()

	// Start server
	go server.NewServer(database)

	// Sync DB every 2min.
	syncTicker := time.NewTicker(time.Minute * 5)
	go func() {
		for _ = range syncTicker.C {
			verifier.VerifyAndDelete(database)
			database.SyncKeys()
			seelog.Debug("verify and sync database.")
		}
	}()

	for {
		pendingTypes := collector.AllType()

		var wg sync.WaitGroup
		for _, pendingType := range pendingTypes {
			wg.Add(1)
			go func(t collector.Type) {
				c := collector.NewCollector(t)
				done := make(chan bool, 1)

				go func() {
					scheduler.RunCollector(c, database)
					done <- true
				}()

				// Set timeout to avoid deadlock.
				select {
				case <-done:
				case <-time.After(7 * time.Minute):
					seelog.Errorf("collector %s time out.", reflect.ValueOf(c).Type().String())
				}

				defer func() {
					if r := recover(); r != nil {
						seelog.Critical(r)
					}
				}()

				defer wg.Done()
			}(pendingType)
		}

		wg.Wait()
		seelog.Debug("finish once, sleep 10 minutes.")
		time.Sleep(time.Minute * 10)
	}
}
