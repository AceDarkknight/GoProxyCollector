package proxyPool

import (
	"errors"
	"strconv"
	"sync"
)

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
	Items []*IpItem
}

func NewProxyPool(capacity int) *ProxyPool {
	if capacity <= 0 {
		capacity = 10000
	}

	return &ProxyPool{make([]*IpItem, 0, capacity)}
}

func (p *ProxyPool) Add(item *IpItem) {
	mutex.Lock()
	defer mutex.Unlock()

	if item.Ip == "" {
		return
	}

	if len(p.Items) < cap(p.Items) {
		p.Items = append(p.Items, item)
	} else {
		for _, v := range p.Items {
			if v == nil || v.Status == UNAVAILABLE {
				v = item
				break
			}
		}
	}
}

func (p *ProxyPool) Get(index int) *IpItem {
	mutex.Lock()
	defer mutex.Unlock()

	if index < 0 {
		return nil
	}

	return p.Items[index]
}

func (p *ProxyPool) Delete(index int) {
	mutex.Lock()
	defer mutex.Unlock()

	if index >= 0 && index < len(p.Items) {
		p.Items[index] = nil
	}
}

func (p *ProxyPool) SetStatus(index int, status IpStatus) error {
	mutex.Lock()
	defer mutex.Unlock()

	if index >= 0 && index < len(p.Items) {
		p.Items[index] = nil
		return nil
	}

	return errors.New("invalid index value: " + strconv.Itoa(index))
}

func (p *ProxyPool) ClearUnavailableItems() {
	mutex.Lock()
	defer mutex.Unlock()

	for _, v := range p.Items {
		if v != nil && v.Status == UNAVAILABLE {
			v = nil
		}
	}
}
