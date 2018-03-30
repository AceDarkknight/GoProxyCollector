package storage

type Storage interface {
	Exist(string) bool
	Get(string) []byte
	Delete(string) bool
	AddOrUpdate(string, interface{}) error
	GetAll() map[string][]byte
	Close()
	GetRandomOne() []byte
}

func NewStorage() (Storage, error) {
	return NewBoltDbStorage("proxy.db", "IpList")
}
