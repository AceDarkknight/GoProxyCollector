package collector

import (
	"strings"

	"github.com/AceDarkkinght/GoProxyCollector/util"
	"github.com/cihub/seelog"
)

type RegexCollector struct {
	configuration *Config
	currentUrl    string
	currentIndex  int
	urls          []string
	selectorMap   map[string][]string
}

func NewRegexCollector(config *Config) *RegexCollector {
	if config == nil {
		return nil
	}

	if !config.Verify() || config.Type != COLLECTBYREGEX || len(config.ValueRuleMap.Items) < 1 {
		seelog.Errorf("config name:%s is unavailable, please check your collectorConfig.xml", config.Name)
		return nil
	}

	parameters := strings.Split(config.UrlParameters, ",")
	urls := util.MakeUrls(config.UrlFormat, parameters)
	return &RegexCollector{
		configuration: config,
		urls:          urls,
		//selectorMap:   selectorMap,
	}
}
