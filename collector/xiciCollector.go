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

func CollectXici(url string) (*[]Result, error) {
	if !strings.HasPrefix(url, "http://www.xicidaili.com") {
		return nil, errors.New(fmt.Sprintf("incorrect url:%s\n", url))
	}

	var (
		results []Result
		err     error
	)

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
		speedString, existSpeed := selection.Find("td:nth-child(7) div").Attr("title")
		liveTimeString := selection.Find("td:nth-child(9)").Text()

		port, _ = strconv.Atoi(portString)
		reg := regexp.MustCompile(`^[1-9]\d*\.\d*|0\.\d*[1-9]\d*$`)
		if strings.Contains(speedString, "秒") {
			speed, err = strconv.ParseFloat(reg.FindString(speedString), 64)
			if err != nil {

			}
		}

		reg = regexp.MustCompile(`^[1-9]\d*$`)
		if strings.Contains(liveTimeString, "天") {
			liveTime, err = strconv.Atoi(reg.FindString(liveTimeString))
			if err != nil {

			}
		}

		if ip != "" && port > 0 && existSpeed && speed > 0 && liveTime > 0 {
			results = append(results,
				Result{Ip: ip,
					Port:     port,
					Location: location,
					Speed:    speed,
					LiveTime: liveTime,
					Source:   "http://www.xicidaili.com"})
		}
	})

	return &results, nil
}
