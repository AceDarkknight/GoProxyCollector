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
	firstIndex   int
	lastIndex    int
	baseUrl      string
	currentUrl   string
}

func NewCoderbusyCollector() *CoderbusyCollector {
	return &CoderbusyCollector{
		baseUrl:    "https://proxy.coderbusy.com/classical/https-ready.aspx?page=",
		firstIndex: 1,
		lastIndex:  15,
	}
}

func (c *CoderbusyCollector) Next() bool {
	if c.currentIndex >= c.lastIndex {
		return false
	}

	c.currentIndex++
	c.currentUrl = c.baseUrl + strconv.Itoa(c.currentIndex)

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
