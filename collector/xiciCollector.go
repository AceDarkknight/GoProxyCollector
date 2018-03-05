package collector

import (
	"fmt"
	"net/http"
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
		result *[]Result
		err    error
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
		liveTime := selection.Find("td:nth-child(9)").Text()
		if strings.Contains(liveTime, "天") {
			ip := selection.Find("td:nth-child(2)").Text()
			port := selection.Find("td:nth-child(3)").Text()
			location := selection.Find("td:nth-child(4) a").Text()
			source := "http://www.xicidaili.com"
			fmt.Println(ip + " " + port + " " + location + " " + source + " " + liveTime)
		}
	})

	return result, nil
}
