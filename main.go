package main

import (
	"github.com/AceDarkkinght/GoProxyCollector/collector"
	"github.com/AceDarkkinght/GoProxyCollector/scheduler"
	"github.com/AceDarkkinght/GoProxyCollector/storage"
)

func main() {
	boltDb, err := storage.NewBoltDbStorage("proxy.db")
	if err != nil {
		panic(err)
	}

	xiciCollector := collector.NewXiciCollector()
	scheduler.Start(xiciCollector, boltDb)
}
