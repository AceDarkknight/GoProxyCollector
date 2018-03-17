package main

import (
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
	logger, err := seelog.LoggerFromConfigAsFile("logConfig.xml")
	if err != nil {
		panic(err)
	}

	seelog.ReplaceLogger(logger)
	seelog.Info("log initialize finish.")
	defer seelog.Flush()

	database, err := storage.NewBoltDbStorage("proxy.db", "IpList")
	if err != nil {
		seelog.Critical(err)
		panic(err)
	}

	seelog.Infof("database initialize finish.")

	// Sync data.
	database.SyncKeys()

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
		c := collector.NewCollector(collector.IP181)
		scheduler.Run(c, database)
	}

	defer database.Close()
}
