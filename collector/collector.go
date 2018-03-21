package collector

import (
	"github.com/AceDarkkinght/GoProxyCollector/result"
)

type Collector interface {
	Next() bool
	Name() string
	Collect(chan<- *result.Result)
}

type Type uint8

const (
	XICI Type = iota
	KUAIDAILI
	DATA5U
	CODERBUSY
	LIUNIAN
	IP181
	IP3366
	KXDAILI
)

func NewCollector(t Type) Collector {
	switch t {
	//case XICI:
	//	return NewXiciCollector()
	//case KUAIDAILI:
	//	return NewKuaidailiCollector()
	//case DATA5U:
	//	return NewData5uCollector()
	case CODERBUSY:
		return NewCoderbusyCollector()
	//case LIUNIAN:
	//	return NewLiunianpCollector()
	//case IP181:
	//	return NewIp181Collector()
	//case IP3366:
	//	return NewIp3366Collector()
	//case KXDAILI:
	//	return NewKxdailiCollector()
	default:
		return nil
	}
}

func AllType() []Type {
	return []Type{
		XICI,
		KUAIDAILI,
		DATA5U,
		CODERBUSY,
		LIUNIAN,
		IP181,
		IP3366,
		KXDAILI,
	}
}
