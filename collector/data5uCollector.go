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

type Data5uCollector struct {
	currentIndex int
	pendingPart  []string
	baseUrl      string
	currentUrl   string
}

func NewData5uCollector() *Data5uCollector {
	return &Data5uCollector{
		currentIndex: -1,
		pendingPart:  []string{"index.shtml", "gngn/index.shtml", "gnpt/index.shtml", "gwgn/index.shtml", "gwpt/index.shtml"},
		baseUrl:      "http://www.data5u.com/free/",
	}
}

func (c *Data5uCollector) Next() bool {
	if c.currentIndex >= len(c.pendingPart)-1 {
		return false
	}

	c.currentIndex++
	c.currentUrl = c.baseUrl + c.pendingPart[c.currentIndex]

	seelog.Debugf("current url:%s", c.currentUrl)
	return true
}

func (c *Data5uCollector) Collect(ch chan<- *result.Result) {
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

	selection := doc.Find(".l2")
	selection.Each(func(i int, sel *goquery.Selection) {
		var (
			port  int
			speed float64
		)

		ip := sel.Find("span:nth-child(1) li").Text()
		portString := sel.Find("span:nth-child(2) li").Text()
		location := sel.Find("span:nth-child(6) a:nth-child(1)").Text() +
			sel.Find("span:nth-child(6) a:nth-child(2)").Text()

		speedString := sel.Find("span:nth-child(8) li").Text()
		if !util.IsIp(ip) {
			ip = ""
		}

		port, _ = strconv.Atoi(portString)

		reg := regexp.MustCompile(`^[1-9]\d*\.\d*|0\.\d*[1-9]\d*`)
		if strings.Contains(speedString, "秒") {
			speed, _ = strconv.ParseFloat(reg.FindString(speedString), 64)
		}

		// Speed must less than 2s and live time must larger than 1 day.
		if ip != "" && port > 0 && speed > 0 && speed < 2 {
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
