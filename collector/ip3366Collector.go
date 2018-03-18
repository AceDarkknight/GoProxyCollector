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

type Ip3366Collector struct {
	currentIndex int
	firstIndex   int
	lastIndex    int
	baseUrl      string
	currentUrl   string
}

func NewIp3366Collector() *Ip3366Collector {
	return &Ip3366Collector{
		firstIndex: 1,
		lastIndex:  10,
		baseUrl:    "http://www.ip3366.net/?stype=1&page=",
	}
}

func (c *Ip3366Collector) Next() bool {
	if c.currentIndex >= c.lastIndex {
		return false
	}

	c.currentIndex++
	c.currentUrl = c.baseUrl + strconv.Itoa(c.currentIndex)

	seelog.Debugf("current url:%s", c.currentUrl)
	return true
}

func (c *Ip3366Collector) Collect(ch chan<- *result.Result) {
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

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		seelog.Errorf("parse %s error:%v", c.currentUrl, err)
		return
	}

	// Because of the charset is gbk, need to decode first.
	decoder := mahonia.NewDecoder("gbk")

	selection := doc.Find("#list tr:not(:first-child)")
	selection.Each(func(i int, sel *goquery.Selection) {
		var (
			port  int
			speed float64
		)

		ip := sel.Find("td:nth-child(1)").Text()
		portString := sel.Find("td:nth-child(2)").Text()
		location := sel.Find("td:nth-child(6)").Text()
		speedString := sel.Find("td:nth-child(7)").Text()

		speedString = decoder.ConvertString(speedString)
		location = decoder.ConvertString(location)

		if !util.IsIp(ip) {
			ip = ""
		}

		port, _ = strconv.Atoi(portString)

		reg := regexp.MustCompile(`^[1-9]\d*\.*\d*|0\.\d*[1-9]\d*`)
		if strings.Contains(speedString, "秒") {
			speed, _ = strconv.ParseFloat(reg.FindString(speedString), 64)
		}

		// Speed must less than 3s.
		if ip != "" && port > 0 && speed >= 0 && speed < 3 {
			r := &result.Result{Ip: ip,
				Port:     port,
				Location: location,
				Speed:    speed,
				Source:   c.currentUrl}

			ch <- r
		}
	})

	seelog.Debugf("finish collect url:%s", c.currentUrl)
	defer close(ch)
}
