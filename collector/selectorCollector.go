package collector

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/AceDarkkinght/GoProxyCollector/result"
	"github.com/AceDarkkinght/GoProxyCollector/util"

	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
	"github.com/cihub/seelog"
	"github.com/parnurzeal/gorequest"
)

// SelectorCollector will use goquery(like jquery) to get the element we need.
type SelectorCollector struct {
	configuration *Config
	currentUrl    string
	currentIndex  int
	urls          []string
	selectorMap   map[string][]string
}

func NewSelectorCollector(config *Config) *SelectorCollector {
	if config == nil {
		return nil
	}

	if !config.Verify() || config.Type != COLLECTBYSELECTOR {
		seelog.Errorf("config name:%s is unavailable, please check your collectorConfig.xml", config.Name)
		return nil
	}

	selectorMap := make(map[string][]string)

	for _, value := range config.ValueRuleMap.Items {
		if value.Name == "table" {
			selectorMap[value.Name] = []string{value.Path}
		} else if value.Attr != "" {
			selectorMap[value.Name] = []string{value.Path, value.Attr}
		} else {
			selectorMap[value.Name] = []string{value.Path}
		}
	}

	// Most website appear their ip list as table, So table item is required.
	// For other situation, you can implement your own method.
	if v, ok := selectorMap["table"]; !ok || v[0] == "" {
		seelog.Errorf("config name:%s table selector's path should not be empty", config.Name)
		return nil
	}

	parameters := strings.Split(config.UrlParameters, ",")
	urls := util.MakeUrls(config.UrlFormat, parameters)
	return &SelectorCollector{
		configuration: config,
		urls:          urls,
		selectorMap:   selectorMap,
	}
}

func (c *SelectorCollector) Next() bool {
	if c.currentIndex >= len(c.urls) {
		return false
	}

	c.currentUrl = c.urls[c.currentIndex]
	c.currentIndex++

	seelog.Debugf("current url:%s", c.currentUrl)
	return true
}

func (c *SelectorCollector) Name() string {
	return c.configuration.Name
}

func (c *SelectorCollector) Collect(ch chan<- *result.Result) {
	// To avoid deadlock, channel must be closed.
	defer close(ch)

	response, _, errs := gorequest.New().Get(c.currentUrl).Set("User-Agent", util.RandomUA()).End()
	if len(errs) > 0 {
		seelog.Errorf("%+v", errs)
		return
	}

	if response.StatusCode != 200 {
		seelog.Errorf("GET %s failed, status code:%s", c.currentUrl, http.StatusText(response.StatusCode))
		return
	}

	defer response.Body.Close()

	// if the charset of website isn't utf-8, need to decode first.
	var decoder mahonia.Decoder
	if c.configuration.Charset != "utf-8" {
		decoder = mahonia.NewDecoder(c.configuration.Charset)
	}

	// Use goquery to find elements.
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		seelog.Errorf("parse %s error:%v", c.currentUrl, err)
		return
	}

	selection := doc.Find(c.selectorMap["table"][0])
	selection.Each(func(i int, sel *goquery.Selection) {
		var (
			ip       string
			port     int
			speed    float64
			location string
		)

		// Find value we need and store in a map.
		nameValue := make(map[string]string)
		for key, value := range c.selectorMap {
			if key != "table" {
				var temp string
				if len(value) == 1 {
					temp = sel.Find(value[0]).Text()
				} else if len(value) == 2 {
					temp, _ = sel.Find(value[0]).Attr(value[1])
				}

				// Decode.
				if temp != "" {
					if decoder != nil {
						temp = decoder.ConvertString(temp)
					}

					nameValue[key] = temp
				}
			}
		}

		if tempIp, ok := nameValue["ip"]; ok && util.IsIp(tempIp) {
			ip = tempIp
		}

		if tempPort, ok := nameValue["port"]; ok {
			port, _ = strconv.Atoi(tempPort)
		}

		if tempSpeed, ok := nameValue["speed"]; ok {
			reg := regexp.MustCompile(`^[1-9]\d*\.*\d*|0\.\d*[1-9]\d*`)
			if strings.Contains(tempSpeed, "秒") {
				speed, _ = strconv.ParseFloat(reg.FindString(tempSpeed), 64)
			}
		}

		if tempLocation, ok := nameValue["location"]; ok {
			location = tempLocation
		}

		// Speed must less than 3s.
		if ip != "" && port > 0 && speed >= 0 && speed < 3 {
			r := &result.Result{
				Ip:       ip,
				Port:     port,
				Location: location,
				Speed:    speed,
				Source:   c.currentUrl}

			seelog.Debugf("%v", r)
			ch <- r
		}
	})

	seelog.Debugf("finish collect url:%s", c.currentUrl)
}
