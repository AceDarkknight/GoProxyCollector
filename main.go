package main

import (
	"github.com/AceDarkkinght/GoProxyCollector/collector"
	"github.com/AceDarkkinght/GoProxyCollector/scheduler"
)

func main() {
	xiciCollector := collector.NewXiciCollector()
	scheduler.Start(xiciCollector)
}
