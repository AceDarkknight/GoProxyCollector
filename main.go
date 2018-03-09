package main

import (
	"github.com/AceDarkkinght/GoProxyCollector/collector"
	"github.com/AceDarkkinght/GoProxyCollector/scheduler"
	"github.com/AceDarkkinght/GoProxyCollector/server"
	"github.com/AceDarkkinght/GoProxyCollector/storage"
)

func main() {
	go server.NewServer()

	boltDb, err := storage.NewBoltDbStorage("proxy.db", "IpList")
	if err != nil {
		panic(err)
	}

	xiciCollector := collector.NewXiciCollector()
	scheduler.Start(xiciCollector, boltDb)
}
