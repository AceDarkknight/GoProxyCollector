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

	// Load database.
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

	// Sync DB every 5min.
	syncTicker := time.NewTicker(time.Minute * 5)
	go func() {
		for _ = range syncTicker.C {
			verifier.VerifyAndDelete(database)
			database.SyncKeys()
			seelog.Debug("verify and sync database.")
		}
	}()

	configs := collector.NewCollectorConfig("collectorConfig.xml")
	for {
		var wg sync.WaitGroup

		for _, c := range configs.Configs {
			wg.Add(1)
			go func(c *collector.Config) {
				defer wg.Done()

				// Panic handle must define fist.
				defer func() {
					if r := recover(); r != nil {
						seelog.Critical(r)
					}
				}()

				col := c.Collector()
				done := make(chan bool, 1)

				go func() {
					scheduler.RunCollector(col, database)
					done <- true
				}()

				// Set timeout to avoid deadlock.
				select {
				case <-done:
				case <-time.After(7 * time.Minute):
					seelog.Errorf("collector %s time out.", reflect.ValueOf(col).Type().String())
				}

			}(c)
		}

		wg.Wait()
		seelog.Debug("finish once, sleep 10 minutes.")
		time.Sleep(time.Minute * 10)
	}
}
