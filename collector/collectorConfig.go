package collector

import "os"

//github.com/antchfx/xmlquery
//github.com/beevik/etree
//type Configs struct {
//	Config []struct {
//	} `xml:"Configs>Config"`
//}

type Config struct {
	Name          string   `xml:"name,attr"`
	urlFormat     string   `xml:"urlFormat"`
	urlParameters []string `xml:"urlParameters"`
	currentUrl    string
	currentIndex  int
	ValueRuleMap  map[string]string
	Type          CollectType
}

type ConfigXml struct {
	Name          string   `xml:"name,attr"`
	UrlFormat     string   `xml:"urlFormat"`
	UrlParameters []string `xml:"urlParameters"`
	ValueRuleMap  struct {
		Items []struct {
			Name string `xml:"name,attr"`
			Path string `xml:"path,attr"`
			Attr string `xml:"attr"`
		} `xml:"item"`
	} `xml:"valueNamePathMap"`
}

type CollectType uint8

const (
	COLLECTBYSELECTOR = iota
	COLLECTBYREGEX
)

func NewCollectorConfig() *Config {
	file, err := os.Open("config.xml")
	if err != nil {
		panic(err)
	}

	defer file.Close()
}
