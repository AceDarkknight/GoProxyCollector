package proxyPool

import "sync"

type IpStatus uint

const (
	AVAILABL = iota
	UNAVAILABLE
	VERIFYING
)

var mutex sync.Mutex

type IpItem struct {
	Ip     string
	Status IpStatus
}

type ProxyPool struct {
	Items []IpItem
}

func NewProxyPool() *ProxyPool {
	return &ProxyPool{make([]IpItem, 1000)}
}

func (p *ProxyPool) Add(item IpItem) {
	mutex.Lock()
	defer mutex.Unlock()

	if item.Ip == "" {
		return
	}

	p.Items = append(p.Items, item)
}

func (p *ProxyPool) Get(index int) IpItem {
	mutex.Lock()
	defer mutex.Unlock()

	if index < 0 {
		return IpItem{Ip: "", Status: UNAVAILABLE}
	}

	return p.Items[index]
}
