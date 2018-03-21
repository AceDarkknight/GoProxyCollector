package scheduler

import (
	"math/rand"
	"os"
	"time"

	"github.com/AceDarkkinght/GoProxyCollector/collector"
	"github.com/AceDarkkinght/GoProxyCollector/result"
	"github.com/AceDarkkinght/GoProxyCollector/storage"
	"github.com/AceDarkkinght/GoProxyCollector/verifier"
	"github.com/cihub/seelog"
)

// RunCollector will start to run a collector and save records to storage.
func RunCollector(collector collector.Collector, storage storage.Storage) {
	if collector == nil || storage == nil {
		return
	}

	seelog.Debugf("start to run collector:%s", collector.Name())
	for {
		resultChan := make(chan *result.Result, 100)
		if !collector.Next() {
			break
		}

		// Collect.
		go collector.Collect(resultChan)

		// Verify.
		verifier.VerifyAndSave(resultChan, storage)

		// Wait at least 5s to avoid the website block our IP.
		t := rand.New(rand.NewSource(time.Now().Unix())).Intn(10) + 5
		seelog.Debugf("sleep %d second", t)
		time.Sleep(time.Duration(t) * time.Second)
	}

	seelog.Debugf("finish to run collector:%s finish", collector.Name())
}

// NewLogger will load the seelog's configuration file.
// If file name is not supplied, it will use default configuration.
func SetLogger(fileName string) {
	if _, err := os.Stat(fileName); err == nil {
		logger, err := seelog.LoggerFromConfigAsFile(fileName)
		if err != nil {
			panic(err)
		}

		seelog.ReplaceLogger(logger)
	} else {
		configString := `<seelog>
                        <outputs formatid="main">
                            <filter levels="info,error,critical">
                                <rollingfile type="date" filename="log/AppLog.log" namemode="prefix" datepattern="02.01.2006"/>
                            </filter>
                            <console/>
                        </outputs>
                        <formats>
                            <format id="main" format="%Date %Time [%LEVEL] %Msg%n"/>
                        </formats>
                        </seelog>`
		logger, err := seelog.LoggerFromConfigAsString(configString)
		if err != nil {
			panic(err)
		}

		seelog.ReplaceLogger(logger)
	}

	seelog.Info("log initialize finish.")
}
