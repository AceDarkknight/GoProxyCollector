package collector

import "github.com/AceDarkkinght/GoProxyCollector/result"

type Collector interface {
	Next() bool
	Collect(chan<- *result.Result)
}

func NewCollector(t Type) Collector {
	switch t {
	case XICI:
		return NewXiciCollector()
	case KUAIDAILI:
		return NewKuaidailiCollector()
	case DATA5U:
		return NewData5uCollector()
	case CODERBUSY:
		return NewCoderbusyCollector()
	case LIUNIAN:
		return NewLiunianpCollector()
	case IP181:
		return NewIp181Collector()
	default:
		return nil
	}
}

type Type uint8

const (
	XICI Type = iota
	KUAIDAILI
	DATA5U
	CODERBUSY
	LIUNIAN
	IP181
)
