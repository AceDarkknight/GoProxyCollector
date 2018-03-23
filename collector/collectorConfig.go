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
	Name          string `xml:"name,attr"`
	UrlFormat     string `xml:"urlFormat"`
	UrlParameters string `xml:"urlParameters"`
	Type          Type   `xml:"collectType"`
	Charset       string `xml:"charset"`
	ValueRuleMap  struct {
		Items []struct {
			Name string `xml:"name,attr"`
			Rule string `xml:"rule,attr"`
			Attr string `xml:"attribute,attr"`
		} `xml:"item"`
	} `xml:"valueNameRuleMap"`
}

// NewCollectorConfig will read collector configuration xml file and parse the xml.
func NewCollectorConfig(fileName string) *Configs {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	defer file.Close()
	var configXml Configs
	decoder := xml.NewDecoder(file)
	// To void panic when there is '&' in xml.
	// This comes from https://stackoverflow.com/questions/35191202/unmarshal-xml-with-unescaped-character-inside.
	decoder.Strict = false
	err = decoder.Decode(&configXml)
	if err != nil {
		panic(err)
	}

	return &configXml
}

// Verify will check the item parse from xml.
func (c *Config) Verify() bool {
	if c.UrlFormat == "" {
		return false
	}

	if c.Charset == "" {
		c.Charset = "utf-8"
	} else {
		c.Charset = strings.ToLower(c.Charset)
	}

	return true
}

// Collector will generate a collector by the CollectType value in config xml.
func (c *Config) Collector() Collector {
	switch c.Type {
	case COLLECTBYSELECTOR:
		return NewSelectorCollector(c)
	case COLLECTBYREGEX:
		return NewRegexCollector(c)
	default:
		return nil
	}
}
