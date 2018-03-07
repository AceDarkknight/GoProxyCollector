package collector

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/AceDarkkinght/GoProxyCollector/util"
	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

type XiciCollector struct {
	currentIndex int
	firstIndex   int
	lastIndex    int
	baseUrl      string
}

// NewXiciCollector will return a new collector of http://www.xicidaili.com.
// The will get first two page of http://www.xicidaili.com by default.
func NewXiciCollector() *XiciCollector {
	return &XiciCollector{
		firstIndex:   1,
		lastIndex:    2,
		currentIndex: 0,
		baseUrl:      "http://www.xicidaili.com/nn/"}
}

// Next will return the next page.
func (c *XiciCollector) Next() string {
	if c.currentIndex >= c.lastIndex {
		return ""
	}

	c.currentIndex++
	return c.baseUrl + strconv.Itoa(c.currentIndex)
}

// Collect will collect the ip and port and other information of the page.
func (c *XiciCollector) Collect(url string) ([]Result, error) {
	if !strings.HasPrefix(url, "http://www.xicidaili.com") {
		return nil, errors.New(fmt.Sprintf("incorrect url:%s\n", url))
	}

	results := make([]Result, 0)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	request.Header.Add("User-Agent", util.RandomUA())
	client := http.DefaultClient
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New(http.StatusText(response.StatusCode))
	}

	defer response.Body.Close()

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
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

		reg := regexp.MustCompile("((?:(?:25[0-5]|2[0-4]\\d|((1\\d{2})|([1-9]?\\d)))\\.){3}(?:25[0-5]|2[0-4]\\d|((1\\d{2})|([1-9]?\\d))))")
		if !reg.Match([]byte(ip)) {
			ip = ""
		}

		port, _ = strconv.Atoi(portString)

		reg = regexp.MustCompile(`^[1-9]\d*\.\d*|0\.\d*[1-9]\d*`)
		if strings.Contains(speedString, "秒") {
			speed, _ = strconv.ParseFloat(reg.FindString(speedString), 64)
		}

		reg = regexp.MustCompile(`^[1-9]\d*`)
		if strings.Contains(liveTimeString, "天") {
			liveTime, _ = strconv.Atoi(reg.FindString(liveTimeString))
		}

		// Speed must less than 1s and live time must larger than 1 day.
		if ip != "" && port > 0 && speed > 0 && speed < 1 && liveTime > 0 {
			results = append(results,
				Result{Ip: ip,
					Port:     port,
					Location: location,
					Speed:    speed,
					LiveTime: liveTime,
					Source:   url})
		}
	})

	return results, nil
}
