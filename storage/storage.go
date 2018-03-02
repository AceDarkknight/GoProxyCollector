package storage

type Storage interface {
	Exist(string) bool
	Get(string) []byte
	Delete(string) bool
	Update(string, interface{}) error
}
