package proxyPool

import (
	"errors"
	"fmt"
	"strconv"
	"sync"

	"github.com/AceDarkkinght/GoProxyCollector/storage"
)

type IpStatus uint

const (
	AVAILABLE = iota
	UNAVAILABLE
)

var mutex sync.Mutex

type IpItem struct {
	Ip     string
	Status IpStatus
}

type ProxyPool struct {
	Items []*IpItem
}

// NewProxyPool will return a ProxyPool whose length equal to 0 and capacity equal to given capacity.
// If given capacity parameter is less than 1, the default capacity set to 10000. For the capacity, larger is better.
func NewProxyPool(capacity int) *ProxyPool {
	if capacity <= 0 {
		capacity = 10000
	}

	return &ProxyPool{make([]*IpItem, 0, capacity)}
}

// Add will add item to ProxyPool.
// If ProxyPool' s length is less than its capacity, append the item. Otherwise, add item to empty position or item is unavailable.
// TODO: Handle the case when add a item to ProxyPool who has no empty position.
func (p *ProxyPool) Add(item *IpItem) {
	mutex.Lock()
	defer mutex.Unlock()

	if item.Ip == "" {
		return
	}

	if len(p.Items) < cap(p.Items) {
		p.Items = append(p.Items, item)
	} else {
		for i := 0; i < len(p.Items); i++ {
			if p.Items[i] == nil || p.Items[i].Status == UNAVAILABLE {
				p.Items[i] = item
			}
		}
	}
}

// Get will return the item in given index. If index is less than 0 or larger than length, return nil.
func (p *ProxyPool) Get(index int) *IpItem {
	mutex.Lock()
	defer mutex.Unlock()

	if index < 0 || index > len(p.Items) {
		return nil
	}

	return p.Items[index]
}

// Delete will delete the item in given index. The trick is from https://github.com/golang/go/wiki/SliceTricks
func (p *ProxyPool) Delete(index int) {
	mutex.Lock()
	defer mutex.Unlock()

	if index >= 0 && index < len(p.Items) {
		p.Items[index] = nil
		copy(p.Items[index:], p.Items[index+1:])
		p.Items[len(p.Items)-1] = nil
		p.Items = p.Items[:len(p.Items)-1]
	}
}

func (p *ProxyPool) SetStatus(index int, status IpStatus) error {
	mutex.Lock()
	defer mutex.Unlock()

	if index >= 0 && index < len(p.Items) {
		p.Items[index].Status = status
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

// Sync will sync with DB. After that, the ProxyPool will be refresh.
func (p *ProxyPool) Sync(storage storage.Storage) error {
	if storage == nil {
		return errors.New("invalid parameter, storage is null")
	}

	items := storage.GetAll()
	fmt.Printf("%v\n", items)
	if len(items) == 0 {
		return nil
	}

	mutex.Lock()
	p.Items = p.Items[0:0:cap(p.Items)]
	mutex.Unlock()

	for k := range items {
		if k != "" {
			p.Add(&IpItem{Ip: k, Status: AVAILABLE})
		}
	}

	return nil
}
