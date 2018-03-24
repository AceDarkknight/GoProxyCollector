package collector

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/AceDarkkinght/GoProxyCollector/result"
	"github.com/AceDarkkinght/GoProxyCollector/util"

	"github.com/cihub/seelog"
	"github.com/parnurzeal/gorequest"
)

type RegexCollector struct {
	configuration *Config
	currentUrl    string
	currentIndex  int
	urls          []string
	selectorMap   map[string]string
}

// NewRegexCollector will create a collector who using regular expression to get item.
func NewRegexCollector(config *Config) *RegexCollector {
	if config == nil {
		return nil
	}

	if !config.Verify() || config.Type != COLLECTBYREGEX || len(config.ValueRuleMap.Items) < 1 {
		seelog.Errorf("config name:%s is unavailable, please check your collectorConfig.xml", config.Name)
		return nil
	}

	selectorMap := make(map[string]string)

	for _, value := range config.ValueRuleMap.Items {
		if value.Name == "" || value.Rule == "" {
			seelog.Errorf("config name:%s contains valueRuleMap item with empty name or rule, this item will be ignored.", config.Name)
			continue
		}

		selectorMap[value.Name] = value.Rule
	}

	parameters := strings.Split(config.UrlParameters, ",")
	urls := util.MakeUrls(config.UrlFormat, parameters)
	return &RegexCollector{
		configuration: config,
		urls:          urls,
		selectorMap:   selectorMap,
	}
}

func (c *RegexCollector) Next() bool {
	if c.currentIndex >= len(c.urls) {
		return false
	}

	c.currentUrl = c.urls[c.currentIndex]
	c.currentIndex++

	seelog.Debugf("current url:%s", c.currentUrl)
	return true
}

func (c *RegexCollector) Name() string {
	return c.configuration.Name
}

// TODO: Support to more websites.
func (c *RegexCollector) Collect(ch chan<- *result.Result) []error {
	// To avoid deadlock, channel must be closed.
	defer close(ch)

	response, bodyString, errs := gorequest.New().Get(c.currentUrl).Set("User-Agent", util.RandomUA()).End()
	if len(errs) > 0 {
		seelog.Errorf("%+v", errs)
		return errs
	}

	if response.StatusCode != 200 {
		errorMessage := fmt.Sprintf("GET %s failed, status code:%s", c.currentUrl, http.StatusText(response.StatusCode))
		seelog.Error(errorMessage)
		return []error{errors.New(errorMessage)}
	}

	if bodyString == "" {
		errorMessage := fmt.Sprintf("parse %s failed, can not parse response body.", c.currentUrl)
		seelog.Error(errorMessage)
		return []error{errors.New(errorMessage)}
	}

	regex := regexp.MustCompile(c.selectorMap["ip"])
	ipAddresses := regex.FindAllString(bodyString, -1)
	if len(ipAddresses) == 0 {
		errorMessage := fmt.Sprintf("can not found correct format ip address in url:%s", c.currentUrl)
		seelog.Error(errorMessage)
		return []error{errors.New(errorMessage)}
	}

	for _, ipAddress := range ipAddresses {
		temp := strings.Split(ipAddress, ":")
		if len(temp) == 2 {
			port, _ := strconv.Atoi(temp[1])
			if port <= 0 {
				continue
			}

			r := &result.Result{
				Ip:     temp[0],
				Port:   port,
				Source: c.currentUrl,
			}

			//seelog.Debugf("%v", r)
			ch <- r
		}
	}

	seelog.Debugf("finish collect url:%s", c.currentUrl)
	return nil
}
