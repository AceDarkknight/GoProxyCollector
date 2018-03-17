package collector

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/AceDarkkinght/GoProxyCollector/result"
	"github.com/AceDarkkinght/GoProxyCollector/util"
	"github.com/cihub/seelog"
	"github.com/parnurzeal/gorequest"
)

type LiunianpCollector struct {
	baseUrl    string
	currentUrl string
}

func NewLiunianpCollector() *LiunianpCollector {
	return &LiunianpCollector{
		baseUrl: "http://www.89ip.cn/tiqv.php?sxb=&tqsl=20&ports=&ktip=&xl=on&submit=%CC%E1++%C8%A1",
	}
}

func (c *LiunianpCollector) Next() bool {
	if c.currentUrl != "" {
		return false
	}

	c.currentUrl = c.baseUrl
	return true
}

func (c *LiunianpCollector) Collect(ch chan<- *result.Result) {
	response, bodyString, errs := gorequest.New().Get(c.currentUrl).Set("User-Agent", util.RandomUA()).End()
	if len(errs) > 0 {
		seelog.Errorf("%+v", errs)
	}

	if response.StatusCode != 200 {
		seelog.Errorf("GET %s failed, status code:%s", c.currentUrl, http.StatusText(response.StatusCode))
		return
	}

	if bodyString == "" {
		seelog.Errorf("parse %s failed, can not find body", c.currentUrl)
		return
	}

	regex := regexp.MustCompile(`((?:(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))):[1-9]\d*`)
	ipAddresses := regex.FindAllString(bodyString, -1)
	if len(ipAddresses) <= 0 {
		seelog.Errorf("can not found correct format ip address in url:%s", c.currentUrl)
		return
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

			ch <- r
		}
	}

	seelog.Debugf("finish collect url:%s", c.currentUrl)
	defer close(ch)
}
