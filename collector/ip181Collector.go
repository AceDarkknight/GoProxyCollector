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

type Ip181Collector struct {
	baseUrl    string
	currentUrl string
}

func NewIp181Collector() *Ip181Collector {
	return &Ip181Collector{
		baseUrl: "http://www.ip181.com/",
	}
}

func (c *Ip181Collector) Next() bool {
	if c.currentUrl != "" {
		return false
	}

	c.currentUrl = c.baseUrl
	return true
}

func (c *Ip181Collector) Collect(ch chan<- *result.Result) {
	response, _, errs := gorequest.New().Get(c.currentUrl).Set("User-Agent", util.RandomUA()).End()
	if len(errs) > 0 {
		seelog.Errorf("%+v", errs)
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

	selection := doc.Find(".ctable tr:not(:first-child)")
	selection.Each(func(i int, sel *goquery.Selection) {
		var (
			port  int
			speed float64
		)

		ip := sel.Find("td:nth-child(1)").Text()
		portString := sel.Find("td:nth-child(2)").Text()
		speedString := sel.Find("td:nth-child(5)").Text()
		location := sel.Find("td:nth-child(6)").Text()

		speedString = decoder.ConvertString(speedString)
		location = decoder.ConvertString(location)

		if !util.IsIp(ip) {
			ip = ""
		}

		port, _ = strconv.Atoi(portString)

		reg := regexp.MustCompile(`^[1-9]\d*\.\d*|0\.\d*[1-9]\d*`)
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

			seelog.Debugf("%v", r)
			ch <- r
		}
	})

	seelog.Debugf("finish collect url:%s", c.currentUrl)
	defer close(ch)
}
