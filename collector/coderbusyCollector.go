package collector

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/AceDarkkinght/GoProxyCollector/result"
	"github.com/AceDarkkinght/GoProxyCollector/util"
	"github.com/PuerkitoBio/goquery"
	"github.com/cihub/seelog"
)

type CoderbusyCollector struct {
	currentIndex int
	currentUrl   string
	urls         []string
	urlParameter []string
}

func NewCoderbusyCollector() *CoderbusyCollector {
	parameter := []string{
		"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15",
	}

	urls := util.MakeUrls("https://proxy.coderbusy.com/classical/https-ready.aspx?page=%s", parameter)

	return &CoderbusyCollector{
		urls:         urls,
		urlParameter: parameter,
	}
}

func (c *CoderbusyCollector) Next() bool {
	if c.currentIndex >= len(c.urls) {
		return false
	}

	c.currentUrl = c.urls[c.currentIndex]
	c.currentIndex++

	seelog.Debugf("current url:%s", c.currentUrl)
	return true
}

func (c *CoderbusyCollector) Collect(ch chan<- *result.Result) {
	request, err := http.NewRequest("GET", c.currentUrl, nil)
	if err != nil {
		seelog.Errorf("make request to call %s error:%v", c.currentUrl, err)
		return
	}

	request.Header.Add("User-Agent", util.RandomUA())
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		seelog.Errorf("GET %s error:%v", c.currentUrl, err)
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

	selection := doc.Find(".table tr:not(:first-child)")
	selection.Each(func(i int, sel *goquery.Selection) {
		var (
			port     int
			speed    float64
			liveTime int
		)

		ip, _ := sel.Find("td:nth-child(2)").Attr("data-ip")
		portString := sel.Find(".port-box").Text()
		location := sel.Find("td:nth-child(3)").Text()
		speedString := sel.Find("td:nth-child(10)").Text()
		liveTimeString := sel.Find("td:nth-child(11)").Text()

		if !util.IsIp(ip) {
			ip = ""
		}

		port, _ = strconv.Atoi(portString)

		reg := regexp.MustCompile(`^[1-9]\d*\.\d*|0\.\d*[1-9]\d*`)
		if strings.Contains(speedString, "秒") {
			speed, _ = strconv.ParseFloat(reg.FindString(speedString), 64)
		}

		reg = regexp.MustCompile(`^[1-9]\d*`)
		if strings.Contains(liveTimeString, "天") {
			liveTime, _ = strconv.Atoi(reg.FindString(liveTimeString))
		}

		// Speed must less than 2s and live time must larger than 1 day.
		if ip != "" && port > 0 && speed > 0 && speed < 2 && liveTime > 0 {
			r := &result.Result{Ip: ip,
				Port:     port,
				Location: location,
				Speed:    speed,
				LiveTime: liveTime,
				Source:   c.currentUrl}

			ch <- r
		}
	})

	seelog.Debugf("finish collect url:%s", c.currentUrl)
	defer close(ch)
}
