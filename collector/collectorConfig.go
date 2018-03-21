package collector

import (
	"encoding/xml"
	"os"
	"strings"
)

type Configs struct {
	Configs []Config `xml:"config"`
}

type Config struct {
	Name          string      `xml:"name,attr"`
	UrlFormat     string      `xml:"urlFormat"`
	UrlParameters string      `xml:"urlParameters"`
	Type          CollectType `xml:"collectType"`
	Charset       string      `xml:"charset"`
	ValueRuleMap  struct {
		Items []struct {
			Name string `xml:"name,attr"`
			Path string `xml:"path,attr"`
			Attr string `xml:"attribute,attr"`
		} `xml:"item"`
	} `xml:"valueNamePathMap"`
}

type CollectType uint8

const (
	COLLECTBYSELECTOR CollectType = iota
	COLLECTBYREGEX
)

// NewCollectorConfig will read collector configuration xml file and parse the xml.
func NewCollectorConfig(fileName string) *Configs {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	var configXml Configs
	err = xml.NewDecoder(file).Decode(&configXml)
	if err != nil {
		panic(err)
	}

	return &configXml
}

// Verify will check the item parse from xml.
func (c *Config) Verify() bool {
	if c.UrlFormat == "" || len(c.ValueRuleMap.Items) < 3 {
		return false
	}

	if c.Charset == "" {
		c.Charset = "utf-8"
	} else {
		c.Charset = strings.ToLower(c.Charset)
	}

	return true
}

func (c *Config) Collector() Collector {
	switch c.Type {
	case COLLECTBYSELECTOR:
		return NewSelectorCollector(c)
	case COLLECTBYREGEX:
		return nil
	default:
		return nil
	}
}
