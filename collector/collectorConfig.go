package collector

type Config struct {
	Name          string
	urlFormat     string
	urlParameters []string
	currentUrl    string
	currentIndex  int
	ValueRuleMap  map[string]string
	Type          CollectType
}

type CollectType uint8

const (
	COLLECTBYSELECTOR = iota
	COLLECTBYREGEX
)

func NewCollectorConfig() *Config {

}
