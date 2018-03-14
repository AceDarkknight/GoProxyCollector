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

type XiciCollector struct {
	currentIndex int
	firstIndex   int
	lastIndex    int
	baseUrl      string
	currentUrl   string
}

// NewXiciCollector will return a new collector of http://www.xicidaili.com.
// The will get first two page of http://www.xicidaili.com by default.
func NewXiciCollector() *XiciCollector {
	return &XiciCollector{
		firstIndex:   1,
		lastIndex:    3,
		currentIndex: 0,
		baseUrl:      "http://www.xicidaili.com/nn/",
		currentUrl:   "http://www.xicidaili.com/nn/"}
}

// Next will return the next page.
func (c *XiciCollector) Next() bool {
	if c.currentIndex >= c.lastIndex {
		return false
	}

	c.currentIndex++
	c.currentUrl = c.baseUrl + strconv.Itoa(c.currentIndex)

	seelog.Debugf("current url:%s", c.currentUrl)
	return true
}

// Collect will collect the ip and port and other information of the page.
func (c *XiciCollector) Collect(ch chan<- *result.Result) {
	if !strings.HasPrefix(c.currentUrl, "http://www.xicidaili.com") {
		return
	}

	request, err := http.NewRequest("GET", c.currentUrl, nil)
	if err != nil {
		//return nil, err
		return
	}

	request.Header.Add("User-Agent", util.RandomUA())
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		//return nil, err
		return
	}

	if response.StatusCode != 200 {
		//return nil, errors.New(http.StatusText(response.StatusCode))
		return
	}

	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		//return nil, err
		return
	}

	selection := doc.Find("#ip_list tr:not(:first-child)")
	selection.Each(func(i int, selection *goquery.Selection) {
		var (
			port     int
			speed    float64
			liveTime int
		)

		ip := selection.Find("td:nth-child(2)").Text()
		portString := selection.Find("td:nth-child(3)").Text()
		location := selection.Find("td:nth-child(4) a").Text()
		speedString, _ := selection.Find("td:nth-child(7) div").Attr("title")
		liveTimeString := selection.Find("td:nth-child(9)").Text()

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
				Source:   c.baseUrl}

			ch <- r
		}
	})

	seelog.Debugf("finish collect url%s", c.currentUrl)
	defer close(ch)
}
