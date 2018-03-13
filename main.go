package main

import (
	"time"

	"github.com/AceDarkkinght/GoProxyCollector/collector"
	"github.com/AceDarkkinght/GoProxyCollector/proxyPool"
	"github.com/AceDarkkinght/GoProxyCollector/scheduler"
	//"github.com/AceDarkkinght/GoProxyCollector/server"
	"github.com/AceDarkkinght/GoProxyCollector/storage"
	"github.com/AceDarkkinght/GoProxyCollector/verifier"
)

func main() {
	//go server.NewServer()

	database, err := storage.NewBoltDbStorage("proxy.db", "IpList")
	if err != nil {
		panic(err)
	}

	pool := proxyPool.NewProxyPool(10000)
	err = pool.Sync(database)
	if err != nil {
		panic(err)
	}

	// Sync ProxyPool with DB every 2min.
	syncTicker := time.NewTicker(time.Minute * 2)
	go func() {
		for _ = range syncTicker.C {
			verifier.VerifyAll(database)
			pool.Sync(database)
		}
	}()

	for {
		xiciCollector := collector.NewXiciCollector()
		scheduler.Run(xiciCollector, database)
	}

	defer database.Close()
}
