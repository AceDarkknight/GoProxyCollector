package collector

import (
	"net/http"
	"strconv"

	"github.com/AceDarkkinght/GoProxyCollector/result"
	"github.com/AceDarkkinght/GoProxyCollector/util"
	"github.com/PuerkitoBio/goquery"
	"github.com/cihub/seelog"
)

type kuaidailiCollector struct {
	currentIndex int
	firstIndex   int
	lastIndex    int
	baseUrl      string
	currentUrl   string
}

func NewKuaidailiCollector() *kuaidailiCollector {
	return &kuaidailiCollector{
		currentIndex: 0,
		firstIndex:   1,
		lastIndex:    5,
		baseUrl:      "https://www.kuaidaili.com/free/inha/"}
}

func (c *kuaidailiCollector) Next() bool {
	if c.currentIndex >= c.lastIndex {
		return false
	}

	c.currentIndex++
	c.currentUrl = c.baseUrl + strconv.Itoa(c.currentIndex)

	seelog.Debugf("current url:%s", c.currentUrl)
	return true
}

func (c *kuaidailiCollector) Collect(ch chan<- *result.Result) {
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
}
